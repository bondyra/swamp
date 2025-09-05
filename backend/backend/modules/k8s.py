from typing import AsyncGenerator, Dict, List

from kubernetes_asyncio import config
from kubernetes_asyncio.dynamic import DynamicClient
import requests

from backend.model import Attribute, Label, Provider, GenericQueryException


description = "Module for interacting with Kubernetes resources"


# TODO: should be schema for each context
_OPENAPI_SCHEMA = None


async def _generate_example_from_openapi_schema(openapi_path: str) -> Dict:
    global _OPENAPI_SCHEMA
    if not _OPENAPI_SCHEMA:
        _OPENAPI_SCHEMA = await _load_openapi_schema()
    definition = _OPENAPI_SCHEMA["definitions"][openapi_path]
    return _generate_example_rec(_OPENAPI_SCHEMA, definition, openapi_path)


async def _load_openapi_schema():
    async with await config.new_client_from_config() as api:
        client = await DynamicClient(api)
        # we cannot get schema with normal client, so making primitive request to API server's openapi endpoint
        cfg = client.configuration
        response = None
        if cfg.api_key:
            headers = {"authorization": cfg.api_key["BearerToken"] if "BearerToken" in cfg.api_key else cfg.api_key["authorization"]}
            response = requests.get(f"{cfg.host}/openapi/v2", headers=headers, verify=cfg.ssl_ca_cert)
        else:
            response = requests.get(f"{cfg.host}/openapi/v2", cert=(cfg.cert_file, cfg.key_file), verify=cfg.ssl_ca_cert)
        assert response.status_code // 100 == 2
        return response.json()


def _generate_example_rec(schema, definition, name) -> Dict:
    properties = definition.get("properties")
    if properties:
        return {
            key: _process_val(schema, key, val)
            for key, val in properties.items()
            if key not in {"apiVersion", "kind"}
        }
    else:
        return f"{name}_VALUE"


def _process_val(schema, key, val):
    if "$ref" in val:
        new_definition = schema["definitions"][val["$ref"].replace("#/definitions/", "")]
        return _generate_example_rec(schema, new_definition, key)
    elif val["type"] == "array":
        if "$ref" in val["items"]:
            new_definition = schema["definitions"][val["items"]["$ref"].replace("#/definitions/", "")]
            el = _generate_example_rec(schema, new_definition, key)
            return [el, el]
        else:
            return [f"{key}_VALUE", f"{key}_VALUE"]
    else:
        return f"{key}_VALUE"


class Kubernetes(Provider):
    @staticmethod
    def provider_name() -> str:
        return "k8s"

    @staticmethod
    def provider_description() -> str:
        return "Provider for interacting with K8S resources via kubernetes-asyncio"
    
    @staticmethod
    def resources() -> List[str]:
        return list(_resources.keys())

    @staticmethod
    def description(r: str) -> str:
        return _resources[r]["description"]
    
    @classmethod
    def icon(cls, r: str) -> str:
        return _resources[r].get("icon", "todo-needs-prefetch-cache-on-frontend")

    @classmethod
    async def attributes(cls, r: str) -> List[Attribute]:
        ctx = Attribute(path="_context", description="Kubernetes context to use", allowed_values=await _get_contexts())
        ns = Attribute(path="_namespace", description="Kubernetes namespace this resource sits in", depends_on="_context")
        return [ctx, ns] if _resources[r]["namespaced"] else [ctx]
    
    @classmethod
    async def attribute_values(cls, r: str, attribute: str, **kwargs) -> List[str]:
        # stupidly hardwired for now (and probably for a long time)
        if attribute == "_namespace":
            if "_context" not in kwargs:
                raise GenericQueryException(f"To get attribute values of namespace, we need context")
            return await _get_namespaces(kwargs["_context"])
        raise GenericQueryException(f"Not supported for attribute {attribute}")
    
    @classmethod
    async def example(cls, r: str) -> Dict:
        # TODO: selecting default context for now, but resource might not exist there! 
        # we should go over contexts and see which one has this resource
        result = await _generate_example_from_openapi_schema(_resources[r]["openapi_path"])
        return result

    @classmethod
    async def get(cls, r: str, labels: Dict[str, Label]) -> AsyncGenerator[Dict, None]:
        namespaced = _resources[r]["namespaced"]
        if "_context" not in labels:
            raise GenericQueryException("You need to provide _context value to query K8S resource")
        if namespaced and "_namespace" not in labels:
            raise GenericQueryException("You need to provide _namespace value to query K8S resource")
        context = labels["_context"].val
        namespace = labels["_namespace"].val if namespaced else None
        async with await config.new_client_from_config(context=context) as api:
            client = await DynamicClient(api)
            v1 = await client.resources.get(api_version=_resources[r]["api_version"], kind=_resources[r]["kind"])
            kwargs, extra_return_values = {}, {}
            if namespaced:
                kwargs["namespace"] = namespace
                extra_return_values["_namespace"] = namespace
            response = await v1.get(**kwargs)
            for item in response.items:
                yield {
                    "_id": item.metadata.name,
                    "_context": context,
                    **extra_return_values,
                    **item.to_dict()
                }


_CONTEXTS = []


async def _get_contexts():
    global _CONTEXTS
    if _CONTEXTS:
        return _CONTEXTS
    await config.load_kube_config()
    contexts, _ = config.list_kube_config_contexts()   
    _CONTEXTS = [ctx["name"] for ctx in contexts]
    return _CONTEXTS


_CONTEXT_TO_NAMESPACES = {}


async def _get_namespaces(context):
    if context in _CONTEXT_TO_NAMESPACES:
        return _CONTEXT_TO_NAMESPACES[context]
    try:
        async with await config.new_client_from_config(context=context) as api:
            client = await DynamicClient(api)
            v1 = await client.resources.get(api_version="v1", kind="Namespace")
            response = await v1.get()
            namespaces = [it.metadata.name for it in response.items]
            _CONTEXT_TO_NAMESPACES[context] = namespaces
    except config.config_exception.ConfigException:
        raise GenericQueryException(f"Invalid context {context}")
    return _CONTEXT_TO_NAMESPACES[context]


_resources = {
    "config_map": {
        "kind": "ConfigMap",
        "api_version": "v1",
        "openapi_path": "io.k8s.api.core.v1.ConfigMap",
        "description": "Config map",
        "namespaced": True,
    },
    "replica_set": {
        "kind": "ReplicaSet",
        "api_version": "v1",
        "openapi_path": "io.k8s.api.apps.v1.ReplicaSet",
        "description": "Replica set",
        "namespaced": True,
    },
    "deployment": {
        "kind": "Deployment",
        "api_version": "v1",
        "openapi_path": "io.k8s.api.apps.v1.Deployment",
        "description": "Deployment",
        "namespaced": True,
    },
    "pod": {
        "kind": "Pod",
        "api_version": "v1",
        "openapi_path": "io.k8s.api.core.v1.Pod",
        "description": "Pod",
        "namespaced": True,
    },
    "pvc": {
        "kind": "PersistentVolumeClaim",
        "api_version": "v1",
        "openapi_path": "io.k8s.api.core.v1.PersistentVolumeClaim",
        "description": "Persistent volume claim",
        "namespaced": True,
    },
    "secret": {
        "kind": "Secret",
        "api_version": "v1",
        "openapi_path": "io.k8s.api.core.v1.Secret",
        "description": "Secret",
        "namespaced": True,
    },
    "service_account": {
        "kind": "ServiceAccount",
        "api_version": "v1",
        "openapi_path": "io.k8s.api.core.v1.ServiceAccount",
        "description": "Service account",
        "namespaced": True,
    },
    "service": {
        "kind": "Service",
        "api_version": "v1",
        "openapi_path": "io.k8s.api.core.v1.Service",
        "description": "Service",
        "namespaced": True,
    },
    "event": {
        "kind": "Event",
        "api_version": "v1",
        "openapi_path": "io.k8s.api.core.v1.Event",
        "description": "Event",
        "namespaced": True,
    },
    "endpoints": {
        "kind": "Endpoints",
        "api_version": "v1",
        "openapi_path": "io.k8s.api.core.v1.Endpoints",
        "description": "Endpoints",
        "namespaced": True,
    },
    "node": {
        "kind": "Node",
        "api_version": "v1",
        "openapi_path": "io.k8s.api.core.v1.Node",
        "description": "Node",
        "namespaced": False,
    },
    "pv": {
        "kind": "PersistentVolume",
        "api_version": "v1",
        "openapi_path": "io.k8s.api.core.v1.PersistentVolume",
        "description": "Persistent volume",
        "namespaced": True,
    },
}