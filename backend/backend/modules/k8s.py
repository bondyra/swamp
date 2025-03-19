from typing import AsyncGenerator, Dict, List

from kubernetes_asyncio import config
from kubernetes_asyncio.dynamic import DynamicClient
import requests

from backend.model import Attribute, Handler, Provider, GenericQueryException


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


class NamespacedK8sHandler(K8sHandler):
    @classmethod
    async def get(cls, **required_attrs) -> AsyncGenerator[Dict, None]:
        if "_context" not in required_attrs:
            raise GenericQueryException("You need to provide _context value to query K8S resource")
        if "_namespace" not in required_attrs:
            raise GenericQueryException("You need to provide _namespace value to query K8S resource")
        context, namespace = required_attrs["_context"], required_attrs["_namespace"]
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
            Attribute(path="_namespace", description="Kubernetes namespace this resource sits in", query_required=True, allowed_values=["default", "kube-system", "kube-public", "kube-node-lease"]),
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
    async def get(cls, **required_attrs) -> AsyncGenerator[Dict, None]:
        if "_context" not in required_attrs:
            raise GenericQueryException("You need to provide _context value to query this K8S resource")
        context = required_attrs["_context"]
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
    async with await config.new_client_from_config(context=context) as api:
        client = await DynamicClient(api)
        v1 = await client.resources.get(api_version="v1", kind="Namespace")
        response = await v1.get()
        namespaces = [it.metadata.name for it in response.items]
        _CONTEXT_TO_NAMESPACES[context] = namespaces
    return _CONTEXT_TO_NAMESPACES[context]
