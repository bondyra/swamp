from collections import namedtuple
from typing import AsyncGenerator, Dict


Result = namedtuple("Result", "id content blob")

modules = {}


class ResourceMeta(type):
    def __new__(cls, *args, **kwargs):
        inst = super().__new__(cls, *args, **kwargs)
        if "name" not in kwargs:
            raise ValueError("Resource must have a name")
        setattr(inst, "__ov_name", kwargs["name"])
        if "module" not in kwargs:
            raise ValueError("Resource must have a module")
        kwargs["module"].__ov_resource_types.append(inst)
        return inst


class Resource(metaclass=ResourceMeta):
    def description(cls) -> str:
        pass

    async def ls(cls) -> AsyncGenerator[Result]:
        pass

    async def get(cls, resource_id) -> Result:
        pass

    def schema_ls(cls) -> Dict[str, str]:
        pass

    def schema_get(cls) -> Dict[str, str]:
        pass


class ModuleMeta(type):
    def __new__(cls, *args, **kwargs):
        if "name" not in kwargs:
            raise ValueError("Module must have a name")
        inst = super().__new__(cls, *args, **kwargs)
        setattr(inst, "__ov_name", kwargs["name"])
        setattr(inst, "__ov_resource_types", [])
        modules[kwargs["name"]] = inst
        return inst


class Module(metaclass=ModuleMeta):
    def resource_types(self) -> Dict[str, Resource]:
        return {r.__ov_name: r for r in self.__ov_resource_types}
