from typing import AsyncGenerator, Dict, Iterable, List, Optional

from pydantic import BaseModel


class Attribute(BaseModel):
    path: str
    description: str
    query_required: bool
    allowed_values: Optional[List[str]] = None


class ResourceType(BaseModel):
    provider: str
    resource: str
    description: str


_provider_registry = {}
_handler_registry = {}


def iter_all_resource_types() -> Iterable[ResourceType]:
    for p in _provider_registry:
        for r in _handler_registry[p]:
            yield ResourceType(
                provider=p, 
                resource=r,
                description =_handler_registry[p][r].description().replace("<p>", "").replace("</p>", "")
            )


class GenericQueryException(Exception):
    pass


class _ProviderMeta(type):
    def __new__(cls, *args, **kwargs):
        inst = super().__new__(cls, *args)
        if inst.name():
            if inst.name() in _provider_registry:
                raise ValueError(f"Provider {inst.name()} already registered")
            _provider_registry[inst.name()] = inst
        return inst


class Provider(metaclass=_ProviderMeta):
    @staticmethod
    def description() -> str:
        pass

    @staticmethod
    def name() -> str:
        pass


class _HandlerMeta(type):
    def __new__(cls, name, bases, *args, **kwargs):
        inst = super().__new__(cls, name, bases, *args)
        if inst.provider() and inst.resource():
            _handler_registry.setdefault(inst.provider(), {})
            _handler_registry[inst.provider()][inst.resource()] = inst
        return inst


class Handler(metaclass=_HandlerMeta):
    @staticmethod
    def provider() -> str:
        pass
    
    @staticmethod
    def resource() -> str:
        pass

    @staticmethod
    def description() -> str:
        pass

    @classmethod
    async def get(cls, **required_attrs) -> AsyncGenerator[Dict, None]:
        pass

    @classmethod
    def attributes(cls) -> List[Attribute]:
        pass


def handler(provider: str, resource: str) -> Handler:
    return _handler_registry[provider][resource]
