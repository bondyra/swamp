from jsonpath_ng import parse, This
import json


def render_ls(results, props, w=32):
    exprs = [parse(p) for p in props]
    print("".join([f"{'id':<{w}}", *(f"{str(e):<{w}}" for e in exprs)]))
    for result in results:
        values = [next((m.value for m in e.find(result.content)), "") for e in exprs]
        print("".join([f"{_trunc(result.id, w):<{w}}", *(f"{_trunc(str(v), w):<{w}}" for v in values)]))


def _trunc(s: str, length: int) -> str:
    return s[:max(0,length-5)] + "..." if len(s) > length else s


def render_get(result, props, indent=2):
    exprs = [parse(p) for p in props] or [parse('*')]
    values = {_concrete_path(m): m.value for e in exprs for m in e.find(result.content)}
    print(json.dumps(values, sort_keys=True, indent=indent).strip())


def _concrete_path(m):
    m0 = m
    results = []
    while not isinstance(m0.path, This):
        results.append(str(f".{m0.path}"))
        m0  = m0.context
    return "".join(reversed(results)).lstrip(".")
