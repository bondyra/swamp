from jsonpath_ng import parse
import yaml


def render_results(results, props, default_props, w=32):
    exprs = [parse(p) for p in props or default_props]
    print("".join([f"{'id':<{w}}", *(f"{str(e):<{w}}" for e in exprs)]))
    for result in results:
        values = [next((m.value for m in e.find(result.content)), "") for e in exprs]
        print("".join([f"{_trunc(result.id, w):<{w}}", *(f"{_trunc(str(v), w):<{w}}" for v in values)]))


def _trunc(s: str, length: int) -> str:
    return s[:max(0,length-5)] + "..." if len(s) > length else s


def render_result(result, props, default_props, indent=4):
    exprs = [parse(p) for p in props or default_props]
    print(exprs)
    content = {}
    for e in exprs:
        dupaa = e.filter(lambda d: True, result.content)
        print(dupaa)
        print('*')
        content.update(dupaa)
    print(yaml.dump(content, sort_keys=False, indent=indent, width=0).strip())
