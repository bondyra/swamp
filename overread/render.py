from jsonpath_ng import parse


async def render(results, props, default_props):
    exprs = [parse(p) for p in props or default_props]
    print_header(exprs, w=32)
    async for result in results:
        print_result(result, exprs, w=32)


def print_header(exprs, w):
    print("".join([f"{'id':<{w}}", *(f"{str(e):<{w}}" for e in exprs)]))


def print_result(result, exprs, w):
    values = [next((m.value for m in e.find(result.content)), "") for e in exprs]
    print("".join([f"{_trunc(result.id, w):<{w}}", *(f"{_trunc(str(v), w):<{w}}" for v in values)]))


def _trunc(s: str, length: int) -> str:
    return s[:max(0,length-5)] + "..." if len(s) > length else s
