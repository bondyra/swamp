from typing import AsyncGenerator, Dict, Iterable, List, Optional
from enum import Enum
import jq
import re

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
    icon: str


class Op(Enum):
    EQUALS = "=="  # todo: change to eq
    NOT_EQUALS = "!="  # todo: change to ne
    CONTAINS = "contains"
    NOT_CONTAINS = "not contains"
    LIKE = "like"
    NOT_LIKE = "not like"


# todo tests
class Label(BaseModel):
    key: str
    val: str
    op: Op
    jq_key: Optional[str] = None

    def matches(self, data: Dict) -> bool:
        if not self.jq_key:
            self.jq_key = jq.compile(self.key if self.key.startswith(".") else f".{self.key}")

        results = [str(x) for x in self.jq_key.input_value(data).all() if x]

        if self.op == Op.EQUALS:
            return results and len(results) == 1 and results[0] == self.val
        if self.op == Op.NOT_EQUALS:
            return results and len(results) == 1 and results[0] != self.val
        if self.op == Op.CONTAINS:
            return results and len(results) > 0 and self.val in results
        if self.op == Op.NOT_CONTAINS:
            return results and len(results) > 0 and self.val not in results
        if self.op in {Op.LIKE, Op.NOT_LIKE}:
            val_regex = re.compile(self.val)
            if self.op == Op.LIKE:
                return results and len(results) == 1 and bool(val_regex.match(results[0]))
            if self.op == Op.NOT_LIKE:
                return results and len(results) == 1 and not val_regex.match(results[0])
        return False


_provider_registry = {}


def _init_providers():
    if not _provider_registry:
        from backend.modules.aws import AWS
        from backend.modules.k8s import Kubernetes 

def provider(name: str):
    _init_providers()
    return _provider_registry[name]


def iter_all_resource_types() -> Iterable[ResourceType]:
    _init_providers()
    for p in _provider_registry:
        for r in _provider_registry[p].resources():
            yield ResourceType(
                provider=p, 
                resource=r,
                description =_provider_registry[p].description(r).replace("<p>", "").replace("</p>", ""),
                icon=_provider_registry[p].icon(r)
            )


class GenericQueryException(Exception):
    pass


class _ProviderMeta(type):
    def __new__(cls, *args, **kwargs):
        inst = super().__new__(cls, *args)
        if inst.provider_name():
            if inst.provider_name() in _provider_registry:
                raise ValueError(f"Provider {inst.provider_name()} already registered")
            _provider_registry[inst.provider_name()] = inst
        return inst


class Provider(metaclass=_ProviderMeta):
    @staticmethod
    def provider_description() -> str:
        pass

    @staticmethod
    def provider_name() -> str:
        pass
    
    @staticmethod
    def resources() -> List[str]:
        pass

    @staticmethod
    def description(resource: str) -> str:
        pass

    @staticmethod
    def icon(resource: str) -> str:
        pass

    @classmethod
    async def get(cls, resource: str, labels: Dict[str, Label]) -> AsyncGenerator[Dict, None]:
        pass

    @classmethod
    async def attributes(cls, resource: str) -> List[Attribute]:
        pass
    
    @classmethod
    async def attribute_values(cls, resource: str, attribute: str, **kwargs) -> List[str]:
        raise GenericQueryException(f"Not supported for resource {resource}, attribute {attribute}")
    
    @classmethod
    def example(cls, resource: str) -> Dict:
        pass
