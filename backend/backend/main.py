import base64
from fastapi import FastAPI, Request
from fastapi.exceptions import HTTPException
from fastapi.middleware.cors import CORSMiddleware
import jsonpath_ng
import uvicorn
from typing import Iterable, Tuple

from backend.model import Label, GenericQueryException, handler, iter_all_resource_types
from backend.modules.aws import AWS
from backend.modules.k8s import Kubernetes


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


@app.get("/link-suggestion")
async def link_suggestion(child_provider: str, child_resource: str, parent_provider: str, parent_resource: str):
    try:
        links = await handler(child_provider, child_resource).links()
        matching_links = [
            l for l in links
            if l.parent_provider == parent_provider and l.parent_resource == parent_resource
        ]
        if not matching_links:
            return {"key": "", "val": "", "op": "eq"}
        m = matching_links[0]
        return {"key": m.path, "val": m.parent_path, "op": m.op}
    except Exception as e:
        return {"key": "", "val": "", "op": "eq"}


@app.get("/attributes")
async def attributes(r: Request):
    provider, resource = extract_provider_and_resource(r)
    try:
        result = await handler(provider, resource).attributes()
        return result
    except GenericQueryException as e:
        raise HTTPException(status_code=400, detail=str(e))


@app.get("/attribute-values")
async def attribute_values(r: Request):
    provider, resource = extract_provider_and_resource(r)
    if "attribute" not in r.query_params:  # todo - label
        raise HTTPException(status_code=400, detail='You must specify "attribute"')
    try:
        result = await handler(provider, resource).attribute_values(**r.query_params)
        return result
    except GenericQueryException as e:
        raise HTTPException(status_code=400, detail=str(e))


@app.get("/get")
async def get(r: Request):
    provider, resource = extract_provider_and_resource(r)
    labels = {k: v for k, v in iter_request_labels(r)}
    try:
        results = await do_get(provider, resource, labels)
        return {"results": results}
    except GenericQueryException as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
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
            l = Label(key_jsonpath=jsonpath_ng.parse(k), val=base64.b64decode(label_val).decode(), op=base64.b64decode(label_op).decode())
            yield k, l
        except ValueError:
            raise HTTPException(status_code=400, detail='Invalid label provided - it must be in form KEY,OP,VAL')
        except UnicodeDecodeError:
            raise HTTPException(status_code=400, detail='Invalid label provided - all keys, ops and vals must be base64 encoded!')


async def do_get(provider, resource, labels):
    # some of the filters might not get used in handler, running this for the second time on actual results
    results = [r async for r in handler(provider, resource).get(labels)]
    print(results)
    return [r for r in results if all(l.matches(r) for l in labels.values())]


def start():
    uvicorn.run("backend.main:app", host="0.0.0.0", port=8000, reload=True)
