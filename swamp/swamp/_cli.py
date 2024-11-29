# !/usr/bin/env python
# PYTHON_ARGCOMPLETE_OK

import asyncio
import argparse, argcomplete
from datetime import datetime, date
import json
from typing import List

from jsonpath_ng import parse, This

from swamp.model import handler, modules, resource_types, Result
from swamp.modules.aws import AWS  # noqa


async def execute_ls(module, resource_type):
    results = [r async for r in handler(module, resource_type).ls()]
    return results


async def execute_get(module, resource_type, id=None):
    result = await handler(module, resource_type).get(id)
    return result


def suggest_modules(prefix, *args, **kwargs):
    return {name: mod.description()  for name, mod in modules().items() if name.startswith(prefix)}


def suggest_resource_type(prefix, parsed_args, **kwargs):
    if parsed_args.module in modules():
        return {r: d for r, d in resource_types(parsed_args.module).items()  if f"{r}".startswith(prefix)}
    return {}


def suggest_resource_id(prefix, parsed_args, **kwargs):
    data = asyncio.run(execute_ls(parsed_args.module, parsed_args.resource_type))
    return [result.id for result in data if result.id.startswith(prefix)]


def suggest_prop_ls(prefix, parsed_args, **kwargs):
    return {p: d for (p, d) in handler(parsed_args.module, parsed_args.resource_type).schema_ls().items() if p.startswith(prefix)}


def suggest_prop_get(prefix, parsed_args, **kwargs):
    return {p: d for (p, d) in handler(parsed_args.module, parsed_args.resource_type).schema_get().items() if p.startswith(prefix)}


parser = argparse.ArgumentParser(prog="s", add_help=False, usage=argparse.SUPPRESS)
subparsers = parser.add_subparsers(dest='op')
get_parser = subparsers.add_parser('get', add_help=False, usage=argparse.SUPPRESS)
get_parser.add_argument("module").completer = suggest_modules
get_parser.add_argument("resource_type").completer = suggest_resource_type
get_parser.add_argument("resource_id").completer = suggest_resource_id
get_parser.add_argument("props", nargs="*", default=[]).completer = suggest_prop_get
ls_parser = subparsers.add_parser('ls', add_help=False, usage=argparse.SUPPRESS)
ls_parser.add_argument("module").completer = suggest_modules
ls_parser.add_argument("resource_type").completer = suggest_resource_type
ls_parser.add_argument("props", nargs="*", default=[]).completer = suggest_prop_ls
argcomplete.autocomplete(parser, append_space=False)



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


def run():
    args = parser.parse_args()
    if args.op == "ls":
        data = asyncio.run(execute_ls(args.module, args.resource_type))
        render_ls(data, args.props)
    if args.op == "get":
        result = asyncio.run(execute_get(args.module, args.resource_type, args.resource_id))
        render_get(result, args.props)
