# !/usr/bin/env python
# PYTHON_ARGCOMPLETE_OK

import asyncio
import argparse, argcomplete

from overread.model import handler, modules, resource_types
from overread.modules.aws import AWS  # noqa
from overread.render import render_ls, render_get


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


def run():
    args = parser.parse_args()
    if args.op == "ls":
        data = asyncio.run(execute_ls(args.module, args.resource_type))
        render_ls(data, args.props)
    if args.op == "get":
        result = asyncio.run(execute_get(args.module, args.resource_type, args.resource_id))
        render_get(result, args.props)
