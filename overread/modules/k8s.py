from kubernetes_asyncio import config
from kubernetes_asyncio.dynamic import DynamicClient


description = "Module for interacting with Kubernetes resources"

_config  = {
    "cm": {
        "kind": "ConfigMap",
        "default_props": ["labels.*", "data"]
    },
    "ep": {
        "kind": "Endpoints"
    },
    "event": {
        "kind": "Event"
    },
    "node": {
        "kind": "Node"
    },
    "pvc": {
        "kind": "PersistentVolumeClaim"
    },
    "pv": {
        "kind": "PersistentVolume"
    },
    "pod": {
        "kind": "Pod"
    },
    "secret": {
        "kind": "Secret"
    },
    "sa": {
        "kind": "ServiceAccount"
    },
    "service": {
        "kind": "Service"
    }
}


async def get(resource_type, id):
    thing = _config[resource_type]
    async with await config.new_client_from_config() as api:
        client = await DynamicClient(api)
        v1 = await client.resources.get(api_version="v1", kind=thing["kind"])
        response = await v1.get()
        for item in response.items:
            yield item.metadata.name, item.to_dict()


def resource_types():
    return list(_config.keys())


def default_props(thing_type):
    return _config[thing_type].get("default_props", []) if thing_type in _config else []
