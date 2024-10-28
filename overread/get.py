# !/usr/bin/env python
# PYTHON_ARGCOMPLETE_OK

import argparse
from typing import Tuple

import argcomplete

from overread.execution import execute
from overread.modules import find_module, load_modules
from overread.render import render
from overread.model import Query


def suggest_thing(prefix, parsed_args, **kwargs):
    if '/' in prefix:
        module_name = prefix.split("/")[0]
        return [thing for thing in modules[module_name].thing_types() if thing.startswith(prefix)]
    else:
        return modules.keys()


get_parser = argparse.ArgumentParser(prog="get")
get_parser.add_argument("thing").completer = suggest_thing
get_parser.add_argument("props", nargs="*", help="Props to display")
get_parser.add_argument("--quiet", "-q", action="store_true", help="Display ID only")
argcomplete.autocomplete(get_parser)

modules = load_modules()


async def get(args):
    argcomplete.autocomplete(get_parser)
    args = get_parser.parse_args(args)
    module_name, module, thing_type = parse_thing(args.thing, modules)
    query = Query(module_name, module, thing_type, args.props, args.quiet)
    data = execute(query)
    await render(query, data)


schema_parser = argparse.ArgumentParser(prog="schema")
schema_parser.add_argument("thing")


async def schema(args):
    args = schema_parser.parse_args(args)
    _, module, thing_type = parse_thing(args.thing, modules)
    schema = await module.schema(thing_type)
    print(schema)


def parse_thing(thing: str, available_modules) -> Tuple[str, str, str]:
    parts = thing.split("/")
    if len(parts) == 1:
        thing_type = parts[0]
        module_name = find_module(thing_type, available_modules)
    else:
        module_name, thing_type = parts[0], parts[1]
        if module_name not in available_modules:
            raise Exception(f"Unsupported module {module_name}. Available ones are: {available_modules.keys()}")
    module = available_modules[module_name]
    if thing_type not in module.thing_types():
        raise Exception(f"Thing type {thing_type} does not have any matches in configured module: {module_name}!")
    return module_name, module, thing_type
