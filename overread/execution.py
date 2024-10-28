from datetime import datetime, date
import json

from overread.model import Query, Result


async def execute(query: Query):
    async for id, content in query.module.get(query.thing_type):
        yield Result(id, content, json.dumps(content, default=_serializer))


def _serializer(obj):
    if isinstance(obj, (datetime, date)):
        return obj.isoformat()
    raise TypeError(f"Type {type(obj)} not serializable")
