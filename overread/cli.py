# !/usr/bin/env python
# PYTHON_ARGCOMPLETE_OK

import asyncio
import argparse, argcomplete
from datetime import datetime, date
import json

from overread.model import modules, Result
from overread.render import render_ls, render_get


async def execute_ls(module, resource_type):
    return [Result(id, content, json.dumps(content, default=_serializer)) async for id, content in modules[module].ls(resource_type)]


async def execute_get(module, resource_type, id=None):
    id, content = await modules[module].get(resource_type, id)
    return Result(id, content, json.dumps(content, default=_serializer))


def _serializer(obj):
    if isinstance(obj, (datetime, date)):
        return obj.isoformat()
    raise TypeError(f"Type {type(obj)} not serializable")


def validate_resource_type(module_name: str, resource_type: str):
    if module_name not in modules:
        raise Exception(f"Unsupported module {module_name}. Available ones are: {modules.keys()}")
    if resource_type not in modules[module_name].resource_types():
        raise Exception(f"Resource type {resource_type} does not have any matches in configured module: {module_name}!")


def suggest_module(prefix, *args, **kwargs):
    return {name: mod.description  for name, mod in modules.items() if name.startswith(prefix)}


def suggest_resource_type(prefix, parsed_args, **kwargs):
    if parsed_args.module in modules:
        return {r: d for r, d in modules[parsed_args.module].resource_types().items()  if f"{r}".startswith(prefix)}
    return {}


def suggest_resource_id(prefix, parsed_args, **kwargs):
    try:
        validate_resource_type(parsed_args.module, parsed_args.resource_type)
    except Exception:
        return []
    data = asyncio.run(execute_ls(parsed_args.module, parsed_args.resource_type))
    return [result.id for result in data if result.id.startswith(prefix)]


def suggest_prop_ls(prefix, parsed_args, **kwargs):
    try:
        validate_resource_type(parsed_args.module, parsed_args.resource_type)
    except Exception as e:
        print("KURWA")
        print(e)
        return []
    return {p: d for (p, d) in modules[parsed_args.module].schema_ls(parsed_args.resource_type) if p.startswith(prefix)}


def suggest_prop_get(prefix, parsed_args, **kwargs):
    try:
        validate_resource_type(parsed_args.module, parsed_args.resource_type)
    except Exception:
        return []
    return {p: d for (p, d) in modules[parsed_args.module].schema_get(parsed_args.resource_type) if p.startswith(prefix)}


parser = argparse.ArgumentParser(prog="s", add_help=False, usage=argparse.SUPPRESS)
subparsers = parser.add_subparsers(dest='op')
get_parser = subparsers.add_parser('get', add_help=False, usage=argparse.SUPPRESS)
get_parser.add_argument("module").completer = suggest_module
get_parser.add_argument("resource_type").completer = suggest_resource_type
get_parser.add_argument("resource_id").completer = suggest_resource_id
get_parser.add_argument("props", nargs="*", default=[]).completer = suggest_prop_get
ls_parser = subparsers.add_parser('ls', add_help=False, usage=argparse.SUPPRESS)
ls_parser.add_argument("module").completer = suggest_module
ls_parser.add_argument("resource_type").completer = suggest_resource_type
ls_parser.add_argument("props", nargs="*", default=[]).completer = suggest_prop_ls
argcomplete.autocomplete(parser, append_space=False)


def run():
    args = parser.parse_args()
    validate_resource_type(args.module, args.resource_type)
    if args.op == "ls":
        data = asyncio.run(execute_ls(args.module, args.resource_type))
        render_ls(data, args.props)
    if args.op == "get":
        result = asyncio.run(execute_get(args.module, args.resource_type, args.resource_id))
        render_get(result, args.props)
