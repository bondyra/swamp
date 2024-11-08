# !/usr/bin/env python
# PYTHON_ARGCOMPLETE_OK

import asyncio
import argparse, argcomplete
from collections import namedtuple
from datetime import datetime, date
import json

from overread.render import render_results, render_result
from overread.builtin_modules import overread_aws, overread_k8s


Result = namedtuple("Result", "id content blob")


modules = {"aws": overread_aws, "k8s": overread_k8s}


async def execute_ls(module, thing_type):
    return [Result(id, content, json.dumps(content, default=_serializer)) async for id, content in module.ls(thing_type)]


async def execute_get(module, thing_type, id=None):
    id, content = await module.get(thing_type, id)
    return Result(id, content, json.dumps(content, default=_serializer))


def _serializer(obj):
    if isinstance(obj, (datetime, date)):
        return obj.isoformat()
    raise TypeError(f"Type {type(obj)} not serializable")


def parse_thing(module_name: str, thing_type: str):
    if module_name not in modules:
        raise Exception(f"Unsupported module {module_name}. Available ones are: {modules.keys()}")
    module = modules[module_name]
    if thing_type not in module.thing_types():
        raise Exception(f"Thing type {thing_type} does not have any matches in configured module: {module_name}!")
    return module, thing_type


def suggest_module(prefix, *args, **kwargs):
    return [m for m in modules if m.startswith(prefix)]


def suggest_thing(prefix, parsed_args, **kwargs):
    if parsed_args.module in modules:
        return [thing for mod in modules for thing in modules[parsed_args.module].thing_types() if f"{thing}".startswith(prefix)]
    return []


def suggest_id(prefix, parsed_args, **kwargs):
    try:
        module, thing_type = parse_thing(parsed_args.module, parsed_args.thing)
    except Exception:
        return []
    data = asyncio.run(execute_ls(module, thing_type))
    return [result.id for result in data if result.id.startswith(prefix)]


def suggest_prop(prefix, parsed_args, **kwargs):
    try:
        module, thing_type = parse_thing(parsed_args.module, parsed_args.thing)
    except Exception:
        return []
    schema = asyncio.run(module.schema_ls(thing_type))
    return [p for p in schema if p.startswith(prefix)][:5]


parser = argparse.ArgumentParser(prog="s", add_help=False, usage=argparse.SUPPRESS)
subparsers = parser.add_subparsers(dest='op')
get_parser = subparsers.add_parser('get', add_help=False, usage=argparse.SUPPRESS)
get_parser.add_argument("module").completer = suggest_module
get_parser.add_argument("thing").completer = suggest_thing
get_parser.add_argument("id").completer = suggest_id
get_parser.add_argument("props", nargs="*", default=[]).completer = suggest_prop
ls_parser = subparsers.add_parser('ls', add_help=False, usage=argparse.SUPPRESS)
ls_parser.add_argument("module").completer = suggest_module
ls_parser.add_argument("thing").completer = suggest_thing
ls_parser.add_argument("props", nargs="*", default=[]).completer = suggest_prop
argcomplete.autocomplete(parser, append_space=False)


def run():
    args = parser.parse_args()
    module, thing_type = parse_thing(args.module, args.thing)
    if args.op == "ls":
        data = asyncio.run(execute_ls(module, thing_type))
        default_props = module.default_props(thing_type)  # todo
        render_results(data, args.props, default_props)
    if args.op == "get":
        result = asyncio.run(execute_get(module, thing_type, args.id))
        default_props = module.default_props(thing_type)  # todo
        render_result(result, args.props, default_props)
