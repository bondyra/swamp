from swamp.model import handler
from swamp.modules.aws import AWS

import uvicorn
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],
    allow_credentials=True,
    allow_methods=["GET"],
    allow_headers=["*"],
)

cache = {}


@app.get("/{module}/{resource_type}")
async def query(module: str, resource_type: str, pattern: str = None):
    all_results = await ls(module, resource_type)
    print(all_results)
    return {"results": [{"id": r.id, "obj": r.obj} for r in all_results if not pattern or pattern in str(r.obj)]}


async def ls(module, resource_type):
    if cache.get(module) is None:
        cache[module] = {}
    elif cache[module].get(resource_type) is not None:
        return cache[module][resource_type]
    results = [r async for r in handler(module, resource_type).ls()]
    cache[module][resource_type] = results
    return results


def start():
    uvicorn.run("backend.main:app", host="0.0.0.0", port=8000, reload=True)
