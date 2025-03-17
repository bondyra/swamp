from typing import AsyncGenerator, Dict, List

from kubernetes_asyncio import config
from kubernetes_asyncio.dynamic import DynamicClient

from backend.model import Attribute, Handler, Provider, GenericQueryException


description = "Module for interacting with Kubernetes resources"


class Kubernetes(Provider):
    @staticmethod
    def name() -> str:
        return "k8s"

    @staticmethod
    def description() -> str:
        return "Provider for interacting with K8S resources via kubernetes-asyncio"



class K8sHandler(Handler):
    @staticmethod
    def provider() -> str:
        return "k8s"


class NamespacedK8sHandler(K8sHandler):
    @classmethod
    async def get(cls, **required_attrs) -> AsyncGenerator[Dict, None]:
        if "metadata.context" not in required_attrs:
            raise GenericQueryException("You need to provide metadata.context value to query K8S resource")
        if "metadata.namespace" not in required_attrs:
            raise GenericQueryException("You need to provide metadata.namespace value to query K8S resource")
        context, namespace = required_attrs["metadata.context"], required_attrs["metadata.namespace"]
        async with await config.new_client_from_config() as api:
            client = await DynamicClient(api)
            v1 = await client.resources.get(api_version="v1", kind=cls.kind())
            response = await v1.get(namespace=namespace)
            for item in response.items:
                yield {
                    "metadata": {
                        "id": item.metadata.name,
                        "context": context,
                        "namespace": namespace
                    },
                    "data": item.to_dict()
                }

    @staticmethod
    def kind() -> str:
        raise NotImplementedError()

    @classmethod
    def attributes(cls) -> List[Attribute]:
        return [
            Attribute(path="metadata.context", description="Kubernetes context to use", query_required=True, allowed_values=_get_contexts()),
            Attribute(path="metadata.namespace", description="AWS region", query_required=True, allowed_values=["default", "kube-system"]),
            # TODO
        ]


class ConfigMapHandler(NamespacedK8sHandler):
    @classmethod
    def description(cls) -> str:
        return "Config map"

    @staticmethod
    def kind() -> str:
        return "ConfigMap"

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
    def resource() -> str:
        return "secret"


class Handler(NamespacedK8sHandler):
    @classmethod
    def description(cls) -> str:
        return "Service account"

    @staticmethod
    def kind() -> str:
        return "ServiceAccount"

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
    def resource() -> str:
        return "ep"

# TODO good luck
# "node": {
#     "kind": "Node"
# },
# "pv" {
#     "kind": "PersistentVolume"
# }

def _get_contexts():
    return ["minikube"]  # TODO
