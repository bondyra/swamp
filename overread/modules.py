import glob
import importlib
import os
import sys
from typing import Dict, Iterable


def load_modules() -> Dict:
    def _load_modules():
        sys.path += [os.path.join(os.path.dirname(os.path.realpath(__file__)), "builtin_modules")]
        if os.getenv("OVERREAD_MODULE_PATHS"):
            module_paths = os.environ["OVERREAD_MODULE_PATHS"].split(",")
            sys.path += module_paths
        ov_modules = {m for m in _discover()}
        for m in ov_modules:
            mod = importlib.import_module(m)
            yield m.split("overread_", 1)[-1], mod

    return dict(_load_modules())


def _discover() -> Iterable[str]:
    for dir in sys.path:
        yield from {
            os.path.splitext(os.path.basename(full_path))[0]
            for full_path in glob.glob(os.path.join(dir, "overread_*.py"))
        }

def find_module(thing: str, available_modules) -> str:
    matches = set()
    for name, mod in available_modules.items():
        if thing in mod.thing_types():
            matches.add(name)
    if len(matches) > 1:
        # warning
        pass
    if not matches:
        raise Exception(f"Cannot find thing {thing} in the modules: {available_modules.keys()}!")
    return matches.pop()
