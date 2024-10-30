# !/usr/bin/env python
# PYTHON_ARGCOMPLETE_OK

import asyncio
import argparse, argcomplete
from collections import namedtuple
from datetime import datetime, date
import json

from overread.render import render
from overread.builtin_modules import overread_aws, overread_k8s


Result = namedtuple("Result", "id content blob")


modules = {"aws": overread_aws, "k8s": overread_k8s}


async def execute(module, thing_type, id=None):
    async for id, content in module.get(thing_type, id):
        yield Result(id, content, json.dumps(content, default=_serializer))


def _serializer(obj):
    if isinstance(obj, (datetime, date)):
        return obj.isoformat()
    raise TypeError(f"Type {type(obj)} not serializable")


def parse_thing_type(thing: str):
    if "/" not in thing:
        raise Exception("Thing must be in format: module/thing_type")
    module_name, thing_type = thing.split("/", 1)
    if module_name not in modules:
        raise Exception(f"Unsupported module {module_name}. Available ones are: {modules.keys()}")
    module = modules[module_name]
    if thing_type not in module.thing_types():
        raise Exception(f"Thing type {thing_type} does not have any matches in configured module: {module_name}!")
    return module, thing_type


def suggest_thing_type(prefix, *args, **kwargs):
    return [f"{mod}/{thing}" for mod in modules for thing in modules[mod].thing_types() if f"{mod}/{thing}".startswith(prefix)]


def suggest_id(prefix, parsed_args, **kwargs):
    try:
        module, thing_type = parse_thing_type(parsed_args.thing_type)
    except Exception:
        return []
    data = execute(module, thing_type)
    c = asyncio.run(collect(data))
    return [r.id for r in c if r.id.startswith(prefix)]


async def collect(data):
    return [d async for d in data]


def suggest_prop(prefix, parsed_args, **kwargs):
    try:
        module, thing_type = parse_thing_type(parsed_args.thing_type)
    except Exception:
        return []
    schema = asyncio.run(module.schema(thing_type))
    return [p for p in schema if p.startswith(prefix)]


parser = argparse.ArgumentParser(prog="s", add_help=False)
subparsers = parser.add_subparsers(dest='op')
get_parser = subparsers.add_parser('get')
get_parser.add_argument("thing_type").completer = suggest_thing_type
get_parser.add_argument("id", nargs="?").completer = suggest_id
get_parser.add_argument("--props", "-p", nargs="*", default=[]).completer = suggest_prop
argcomplete.autocomplete(parser, append_space=False)


def run():
    args = parser.parse_args()
    if args.op == "get":
        module, thing_type = parse_thing_type(args.thing_type)
        data = execute(module, thing_type, args.id)
        # TODO: render appropriately when id != null
        default_props = module.default_props(thing_type)
        asyncio.run(render(data, args.props, default_props))
