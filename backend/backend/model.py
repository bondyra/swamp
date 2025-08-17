from typing import Any, AsyncGenerator, Dict, Iterable, List, Optional
import jq

from pydantic import BaseModel


class Attribute(BaseModel):
    path: str
    description: str
    allowed_values: Optional[List[str]] = None
    depends_on: Optional[str] = None


class ResourceType(BaseModel):
    provider: str
    resource: str
    description: str


class Label(BaseModel):
    key: str
    val: str
    op: str
    jq_key: Optional[str] = None

    def matches(self, data: Dict) -> bool:
        if not self.jq_key:
            self.jq_key = jq.compile(self.key if self.key.startswith(".") else f".{self.key}")

        results = [str(x) for x in self.jq_key.input_value(data).all() if x]
        if self.op == "==":
            return results and len(results) == 1 and results[0] == self.val
        if self.op == "!=":
            return results and len(results) == 1 and results[0] != self.val
        if self.op == "~~":
            return results and len(results) > 0 and self.val in results
        if self.op == "!~":
            return results and len(results) > 0 and self.val not in results
        return False


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
    async def get(cls, labels: Dict[str, Label]) -> AsyncGenerator[Dict, None]:
        pass

    @classmethod
    async def attributes(cls) -> List[Attribute]:
        pass
    
    @classmethod
    async def attribute_values(cls, attribute: str, **kwargs) -> List[str]:
        raise GenericQueryException(f"Not supported for attribute {attribute}")
    
    @classmethod
    def example(cls) -> Dict:
        pass


def handler(provider: str, resource: str) -> Handler:
    return _handler_registry[provider][resource]
