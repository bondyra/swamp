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

@app.get("/{module}/{resource_type}")
async def read_root(module: str, resource_type: str):
    resources = {resource_id: content async for resource_id, content in handler(module, resource_type).ls()}
    return resources

def start():
    uvicorn.run("backend.main:app", host="0.0.0.0", port=8000, reload=True)
