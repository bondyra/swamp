from typing import AsyncGenerator, Dict, List

from kubernetes_asyncio import config
from kubernetes_asyncio.dynamic import DynamicClient
import requests

from backend.model import Attribute, Handler, Label, Provider, GenericQueryException, LinkInfo


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


async def _load_attributes_from_openapi_schema(openapi_path: str) -> List[Attribute]:
    global _OPENAPI_SCHEMA
    if not _OPENAPI_SCHEMA:
        _OPENAPI_SCHEMA = await _load_openapi_schema()
    definition = _OPENAPI_SCHEMA["definitions"][openapi_path]
    return _attributes_rec(_OPENAPI_SCHEMA, definition, path="")


async def _load_openapi_schema():
    async with await config.new_client_from_config() as api:
        client = await DynamicClient(api)
        # we cannot get schema with normal client, so making primitive request to API server's openapi endpoint
        cfg = client.configuration
        response = requests.get(f"{cfg.host}/openapi/v2", cert=(cfg.cert_file, cfg.key_file), verify=cfg.ssl_ca_cert)
        assert response.status_code == 200
        return response.json()


def _attributes_rec(schema, definition, path=""):
    properties = definition.get("properties")
    if properties:
        for key, val in properties.items():
            if key in {"apiVersion", "kind"}:
                continue
            if "$ref" in val:
                new_definition = schema["definitions"][val["$ref"].replace("#/definitions/", "")]
                yield from _attributes_rec(schema, new_definition, path=f"{path}.{key}" if path else key)
            elif val["type"] == "array":  # TODO: don't really know what to do with lists for now
                continue
            else:
                yield Attribute(path=f"{path}.{key}" if path else key, description=val["description"], query_required=False)
    else:
        yield Attribute(path=path, description=definition["description"], query_required=False)


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
        # TODO: allowed_values of namespace depends on context, which must be chosen!
        # TODO: actual attributes might vary on context due to different kube versions! selecting default context for now
        resource_attributes = await _load_attributes_from_openapi_schema(cls.openapi_path())
        
        return [
            Attribute(path="_context", description="Kubernetes context to use", query_required=True, allowed_values=await _get_contexts()),
            Attribute(path="_namespace", description="Kubernetes namespace this resource sits in", query_required=True, depends_on="_context"),
            *resource_attributes
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

    @classmethod
    async def links(cls) -> List[LinkInfo]:
        return [
        ]
    

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

    @classmethod
    async def links(cls) -> List[LinkInfo]:
        return [
            LinkInfo(path="metadata.ownerReferences[*].name", parent_provider= "k8s", parent_resource="deployment", parent_path="metadata.name")
            # link to sa todo
            # link to secret todo
            # link to cm todo
        ]
    

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
    
    @classmethod
    async def links(cls) -> List[LinkInfo]:
        return [
            # link to sa todo
            # link to secret todo
            # link to cm todo
        ]


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
    
    @classmethod
    async def links(cls) -> List[LinkInfo]:
        return [
            LinkInfo(path="metadata.ownerReferences[*].name", parent_provider= "k8s", parent_resource="rs", parent_path="metadata.name")
            # link to sa todo
            # link to secret todo
            # link to cm todo
            # link to service todo
            # link to ep todo
        ]


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
    
    @classmethod
    async def links(cls) -> List[LinkInfo]:
        return [
            # link to PV todo
        ]


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
    
    @classmethod
    async def links(cls) -> List[LinkInfo]:
        return [
        ]


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
    
    @classmethod
    async def links(cls) -> List[LinkInfo]:
        return [
        ]


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
    
    @classmethod
    async def links(cls) -> List[LinkInfo]:
        return [
        ]


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
    
    @classmethod
    async def links(cls) -> List[LinkInfo]:
        return [
        ]


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
    
    @classmethod
    async def links(cls) -> List[LinkInfo]:
        return [
            # link to service todo
        ]


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


class NodeHandler(NamespacedK8sHandler):
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
    
    @classmethod
    async def links(cls) -> List[LinkInfo]:
        return [
        ]


class PersistentVolumeHandler(NamespacedK8sHandler):
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
    
    @classmethod
    async def links(cls) -> List[LinkInfo]:
        return [
        ]


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
