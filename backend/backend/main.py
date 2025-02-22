from backend.model import handler
from backend.modules.aws import AWS

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


# fixme add sane serialization
@app.get("/{module}/{resource_type}")
async def query(module: str, resource_type: str):
    results = await ls(module, resource_type)
    print(results)
    return {"results": results}


async def ls(module, resource_type):
    if not cache.get(module):
        cache[module] = {}
    elif cache[module].get(resource_type) is not None:
        return cache[module][resource_type]
    results = [{"id": r.id, "obj": r.obj} async for r in handler(module, resource_type).ls()]
    cache[module][resource_type] = results
    return results


def start():
    uvicorn.run("backend.main:app", host="0.0.0.0", port=8000, reload=True)
