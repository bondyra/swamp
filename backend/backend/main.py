from fastapi import FastAPI, Request
from fastapi.exceptions import HTTPException
from fastapi.middleware.cors import CORSMiddleware
import jsonpath_ng
import uvicorn

from backend.model import all_resource_types, GenericQueryException, handler
from backend.modules.aws import AWS


app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],
    allow_credentials=True,
    allow_methods=["GET"],
    allow_headers=["*"],
)


_cache = {}


@app.get("/resource-types")
async def resource_types():
    return all_resource_types()


@app.get("/attributes")
async def attributes(r: Request):
    validate(r)
    try:
        return handler(r.query_params["provider"], r.query_params["resource"]).attributes()
    except GenericQueryException as e:
        raise HTTPException(status_code=400, detail=str(e))


@app.get("/get")
async def get(r: Request):
    validate(r)
    try:
        results = await do_get(**r.query_params)
        return {"results": results}
    except GenericQueryException as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


def validate(request: Request):
    if "provider" not in request.query_params:
        raise HTTPException(status_code=400, detail='You must specify "provider"')
    if "resource" not in request.query_params:
        raise HTTPException(status_code=400, detail='You must specify "resource"')


async def do_get(provider, resource, **attrs):
    # some of the filters might not get used in handler, running this for the second time on actual results
    results = await _cached(provider, resource, **attrs)
    path_vals = [(jsonpath_ng.parse(key), val) for key, val in attrs.items()]
    return [r for r in results if _matches(r, path_vals)]


async def _cached(provider, resource, **attrs):
    if not _cache.get(provider):
        _cache[provider] = {}
    elif _cache[provider].get(resource) is not None:
        return _cache[provider][resource]
    results = [r async for r in handler(provider, resource).get(**attrs)]
    _cache[provider][resource] = results
    return results


def _matches(res: dict, path_vals) -> bool:
    return all(_get_val(res, path) == val for path, val in path_vals)


def _get_val(res, path):
    x = path.find(res)
    return x[0].value if x and len(x) > 0 else None


def start():
    uvicorn.run("backend.main:app", host="0.0.0.0", port=8000, reload=True)
