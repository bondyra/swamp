import asyncio
import sys

from overread.get import get, schema


def main():
    if len(sys.argv) == 1:
        print("usage")
    else:
        cmd = sys.argv[1]
        if cmd == "get":
            asyncio.run(get(sys.argv[2:]))
        if cmd == "schema":
            asyncio.run(schema(sys.argv[2:]))
    return
