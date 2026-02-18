import base64
from fastapi import FastAPI, Request
from fastapi.exceptions import HTTPException
from fastapi.middleware.cors import CORSMiddleware
import uvicorn
from typing import Iterable, Tuple
import logging

from backend.model import Label, GenericQueryException, iter_all_resource_types, provider


app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],
    allow_credentials=True,
    allow_methods=["GET"],
    allow_headers=["*"],
)


# todo: cache?


@app.get("/resource-types")
async def resource_types():
    return list(iter_all_resource_types())


@app.get("/attributes")
async def attributes(r: Request):
    p, resource = extract_provider_and_resource(r)
    try:
        result = await provider(p).attributes(resource)
        return result
    except GenericQueryException as e:
        raise HTTPException(status_code=400, detail=str(e))


@app.get("/attribute-values")
async def attribute_values(r: Request):
    p, resource = extract_provider_and_resource(r)
    if "attribute" not in r.query_params:  # todo - label
        raise HTTPException(status_code=400, detail='You must specify "attribute"')
    try:
        result = await provider(p).attribute_values(resource, **r.query_params)
        return result
    except GenericQueryException as e:
        raise HTTPException(status_code=400, detail=str(e))


@app.get("/example")
async def example(r: Request):
    p, resource = extract_provider_and_resource(r)
    try:
        result = await provider(p).example(resource)
        return result
    except GenericQueryException as e:
        raise HTTPException(status_code=400, detail=str(e))


@app.get("/get")
async def get(r: Request):
    try:
        p, resource = extract_provider_and_resource(r)
        labels = {k: v for k, v in iter_request_labels(r)}
        results = await do_get(p, resource, labels)
        return {"results": results}
    except GenericQueryException as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        logging.exception(e)
        raise HTTPException(status_code=500, detail=str(e))


def extract_provider_and_resource(request: Request):
    qp = request.query_params
    if "_provider" not in qp:
        raise HTTPException(status_code=400, detail='You must specify "_provider"')
    if "_resource" not in qp:
        raise HTTPException(status_code=400, detail='You must specify "_resource"')
    return qp["_provider"], qp["_resource"]


def iter_request_labels(request: Request) -> Iterable[Tuple[str, Label]]:
    for key, val in request.query_params.items():
        if val:
            continue  # thats not a label - everything must be in key
        try:
            label_key, label_op, label_val = key.split(",")
            k = base64.b64decode(label_key).decode()
            l = Label(key=k, val=base64.b64decode(label_val).decode(), op=base64.b64decode(label_op).decode())
            yield k, l
        except ValueError:
            raise HTTPException(status_code=400, detail='Invalid label provided - it must be in form KEY,OP,VAL')
        except UnicodeDecodeError:
            raise HTTPException(status_code=400, detail='Invalid label provided - all keys, ops and vals must be base64 encoded!')


async def do_get(p, resource, labels):
    # some of the filters might not get used, running this for the second time on actual results
    results = [r async for r in provider(p).get(resource, labels)]
    return [r for r in results if all(l.matches(r) for l in labels.values())]


def start():
    uvicorn.run("backend.main:app", host="0.0.0.0", port=8000, reload=True)
