from collections import namedtuple
from typing import AsyncGenerator, Dict, List


__all__ = ["Module", "Result"]


ResourceDesc = namedtuple("ResourceDesc", "module resource_type id")
Result = namedtuple("Result", "id obj")

_module_registry = {}
_handler_registry = {}


def modules():
    return _module_registry


def resource_types(module: str):
    return {r: rr.description() for r, rr in _handler_registry[module].items()}


def handler(module: str, resource_type: str):
    return _handler_registry[module][resource_type]


class _ModuleMeta(type):
    def __new__(cls, *args, **kwargs):
        inst = super().__new__(cls, *args)
        if inst.name():
            if inst.name() in _module_registry:
                raise ValueError(f"Module name {inst.name()} already registered")
            _module_registry[inst.name()] = inst
        return inst


class Module(metaclass=_ModuleMeta):
    @staticmethod
    def description() -> str:
        pass

    @staticmethod
    def name() -> str:
        pass


class _HandlerMeta(type):
    def __new__(cls, name, bases, *args, **kwargs):
        inst = super().__new__(cls, name, bases, *args)
        if inst.module() and inst.resource_type():
            _handler_registry.setdefault(inst.module(), {})
            _handler_registry[inst.module()][inst.resource_type()] = inst
        return inst


class Handler(metaclass=_HandlerMeta):
    @staticmethod
    def module() -> str:
        pass
    
    @staticmethod
    def resource_type() -> str:
        pass

    @staticmethod
    def description() -> str:
        pass

    @classmethod
    async def ls(cls) -> AsyncGenerator[Result, None]:
        pass

    @classmethod
    async def get(cls, resource_id: str) -> Result:
        pass

    @classmethod
    def parents(cls, obj: Dict) -> List[ResourceDesc]:
        pass

    @classmethod
    def children(cls, obj: Dict) -> List[ResourceDesc]:
        pass

    @classmethod
    def schema_ls(cls) -> Dict[str, str]:
        pass

    @classmethod
    def schema_get(cls) -> Dict[str, str]:
        pass
