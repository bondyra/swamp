from datetime import datetime, date
import json
from typing import List

from jsonpath_ng import parse, This

from overread.model import Result


def _serializer(obj):
    if isinstance(obj, (datetime, date)):
        return obj.isoformat()
    raise TypeError(f"Type {type(obj)} not serializable")


def render_ls(results: List[Result], props: List[str], w: int = 32):
    exprs = [parse(p) for p in props]
    print("".join([f"{'id':<{w}}", *(f"{str(e):<{w}}" for e in exprs)]))
    for result in results:
        values = [next((m.value for m in e.find(result.obj)), "") for e in exprs]
        print("".join([f"{_trunc(result.id, w):<{w}}", *(f"{_trunc(str(v), w):<{w}}" for v in values)]))


def _trunc(s: str, length: int) -> str:
    return s[:max(0,length-5)] + "..." if len(s) > length else s


def render_get(result: Result, props: List[str], indent: int = 2):
    exprs = [parse(p) for p in props] or [parse('*')]
    values = {_concrete_path(m): m.value for e in exprs for m in e.find(result.obj)}
    print(json.dumps(values, default=_serializer, sort_keys=True, indent=indent))


def _concrete_path(m):
    m0 = m
    results = []
    while not isinstance(m0.path, This):
        results.append(str(f".{m0.path}"))
        m0  = m0.context
    return "".join(reversed(results)).lstrip(".")
