from typing import AsyncGenerator, Dict, List

from kubernetes_asyncio import config
from kubernetes_asyncio.dynamic import DynamicClient
import requests

from backend.model import Attribute, Handler, Label, Provider, GenericQueryException


description = "Module for interacting with Kubernetes resources"


class Kubernetes(Provider):
    @staticmethod
    def name() -> str:
        return "k8s"

    @staticmethod
    def description() -> str:
        return "Provider for interacting with K8S resources via kubernetes-asyncio"

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


class K8sHandler(Handler):
    @staticmethod
    def provider() -> str:
        return "k8s"
    
    @classmethod
    async def attribute_values(cls, attribute: str, **kwargs) -> List[str]:
        # stupidly hardwired for now (and probably for a long time)
        if attribute == "_namespace":
            if "_context" not in kwargs:
                raise GenericQueryException(f"To get attribute values of namespace, we need context")
            return await _get_namespaces(kwargs["_context"])
        raise GenericQueryException(f"Not supported for attribute {attribute}")
    
    @classmethod
    async def example(cls) -> Dict:
        # TODO: actual attributes might vary on context due to different kube versions! selecting default context for now
        result = await _generate_example_from_openapi_schema(cls.openapi_path())
        return result


class NamespacedK8sHandler(K8sHandler):
    @classmethod
    async def get(cls, labels: Dict[str, Label]) -> AsyncGenerator[Dict, None]:
        if "_context" not in labels:
            raise GenericQueryException("You need to provide _context value to query K8S resource")
        if "_namespace" not in labels:
            raise GenericQueryException("You need to provide _namespace value to query K8S resource")
        context, namespace = labels["_context"].val, labels["_namespace"].val
        async with await config.new_client_from_config(context=context) as api:
            client = await DynamicClient(api)
            v1 = await client.resources.get(api_version="v1", kind=cls.kind())
            response = await v1.get(namespace=namespace)
            for item in response.items:
                yield {
                    "_id": item.metadata.name,
                    "_context": context,
                    "_namespace": namespace,
                    **item.to_dict()
                }

    @staticmethod
    def kind() -> str:
        raise NotImplementedError()

    @classmethod
    async def attributes(cls) -> List[Attribute]:
        return [
            Attribute(path="_context", description="Kubernetes context to use", allowed_values=await _get_contexts()),
            Attribute(path="_namespace", description="Kubernetes namespace this resource sits in", depends_on="_context")
        ]


class ConfigMapHandler(NamespacedK8sHandler):
    @classmethod
    def description(cls) -> str:
        return "Config map"

    @staticmethod
    def kind() -> str:
        return "ConfigMap"

    @staticmethod
    def openapi_path() -> str:
        return "io.k8s.api.core.v1.ConfigMap"

    @staticmethod
    def resource() -> str:
        return "cm"
    

class ReplicaSetHandler(NamespacedK8sHandler):
    @classmethod
    def description(cls) -> str:
        return "Replica set"

    @staticmethod
    def kind() -> str:
        return "ReplicaSet"

    @staticmethod
    def openapi_path() -> str:
        return "io.k8s.api.apps.v1.ReplicaSet"

    @staticmethod
    def resource() -> str:
        return "rs"
    

class DeploymentHandler(NamespacedK8sHandler):
    @classmethod
    def description(cls) -> str:
        return "Deployment"

    @staticmethod
    def kind() -> str:
        return "Deployment"

    @staticmethod
    def openapi_path() -> str:
        return "io.k8s.api.apps.v1.Deployment"

    @staticmethod
    def resource() -> str:
        return "deployment"


class PodHandler(NamespacedK8sHandler):
    @classmethod
    def description(cls) -> str:
        return "Pod"

    @staticmethod
    def kind() -> str:
        return "Pod"

    @staticmethod
    def openapi_path() -> str:
        return "io.k8s.api.core.v1.Pod"

    @staticmethod
    def resource() -> str:
        return "pod"


class PersistentVolumeClaimHandler(NamespacedK8sHandler):
    @classmethod
    def description(cls) -> str:
        return "Persistent volume claim"

    @staticmethod
    def kind() -> str:
        return "PersistentVolumeClaim"

    @staticmethod
    def openapi_path() -> str:
        return "io.k8s.api.core.v1.PersistentVolumeClaim"

    @staticmethod
    def resource() -> str:
        return "pvc"


class SecretHandler(NamespacedK8sHandler):
    @classmethod
    def description(cls) -> str:
        return "Secret"

    @staticmethod
    def kind() -> str:
        return "Secret"

    @staticmethod
    def openapi_path() -> str:
        return "io.k8s.api.core.v1.Secret"

    @staticmethod
    def resource() -> str:
        return "secret"


class ServiceAccountHandler(NamespacedK8sHandler):
    @classmethod
    def description(cls) -> str:
        return "Service account"

    @staticmethod
    def kind() -> str:
        return "ServiceAccount"

    @staticmethod
    def openapi_path() -> str:
        return "io.k8s.api.core.v1.ServiceAccount"

    @staticmethod
    def resource() -> str:
        return "sa"


class ServiceHandler(NamespacedK8sHandler):
    @classmethod
    def description(cls) -> str:
        return "Service"

    @staticmethod
    def kind() -> str:
        return "Service"

    @staticmethod
    def openapi_path() -> str:
        return "io.k8s.api.core.v1.Service"

    @staticmethod
    def resource() -> str:
        return "service"


class EventHandler(NamespacedK8sHandler):
    @classmethod
    def description(cls) -> str:
        return "Event"

    @staticmethod
    def kind() -> str:
        return "Event"

    @staticmethod
    def openapi_path() -> str:
        return "io.k8s.api.core.v1.Event"

    @staticmethod
    def resource() -> str:
        return "event"


class EndpointsHandler(NamespacedK8sHandler):
    @classmethod
    def description(cls) -> str:
        return "Endpoints - what service manages"

    @staticmethod
    def kind() -> str:
        return "Endpoints"

    @staticmethod
    def openapi_path() -> str:
        return "io.k8s.api.core.v1.Endpoints"

    @staticmethod
    def resource() -> str:
        return "ep"


class GlobalK8sHandler(K8sHandler):
    @classmethod
    async def get(cls, labels: Dict[str, Label]) -> AsyncGenerator[Dict, None]:
        if "_context" not in labels:
            raise GenericQueryException("You need to provide _context value to query this K8S resource")
        context = labels["_context"].val
        async with await config.new_client_from_config() as api:
            client = await DynamicClient(api)
            v1 = await client.resources.get(api_version="v1", kind=cls.kind())
            response = await v1.get()
            for item in response.items:
                yield {
                    "_id": item.metadata.name,
                    "_context": context,
                    **item.to_dict()
                }


class NodeHandler(NamespacedK8sHandler):  # TODO: it shouldn't be namespaced
    @classmethod
    def description(cls) -> str:
        return "Kubernetes node"

    @staticmethod
    def kind() -> str:
        return "Node"

    @staticmethod
    def openapi_path() -> str:
        return "io.k8s.api.core.v1.Node"

    @staticmethod
    def resource() -> str:
        return "node"


class PersistentVolumeHandler(NamespacedK8sHandler):  # TODO: it shouldn't be namespaced
    @classmethod
    def description(cls) -> str:
        return "Persistent volume"

    @staticmethod
    def kind() -> str:
        return "PersistentVolume"

    @staticmethod
    def openapi_path() -> str:
        return "io.k8s.api.core.v1.PersistentVolume"

    @staticmethod
    def resource() -> str:
        return "pv"


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
