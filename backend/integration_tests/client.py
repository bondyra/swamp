import requests
import base64
from urllib.parse import quote


class TestClient:
    def __init__(self, port):
        self._url = f"http://localhost:{port}"
    
    def _response(self, r):
        response = ""
        try:
            response = r.json()
        except:
            pass
        return {
            "status": r.status_code,
            "response": response
        }

    def resource_types(self):
        return self._response(requests.get(f"{self._url}/resource-types"))

    def attributes(self, provider: str, resource: str):
        return self._response(requests.get(f"{self._url}/attributes?_provider={provider}&_resource={resource}"))

    def attribute_values(self, provider: str, resource: str, **params):
        qs = "&".join((f"{k}={v}" for k,v in {"_provider": provider, "_resource": resource, **params}.items()))
        qs = f"&{qs}" if qs else ""
        return self._response(requests.get(f"{self._url}/attribute-values?{qs}"))

    def example(self, provider: str, resource: str):
        return self._response(requests.get(f"{self._url}/example?_provider={provider}&_resource={resource}"))

    def get(self, provider: str, resource: str, *params, **pparams):
        blob = "&".join((
            f"{quote(base64.b64encode(p[0].encode("utf-8")).decode())},{quote(base64.b64encode(p[1].encode("utf-8")).decode())},{quote(base64.b64encode(p[2].encode("utf-8")).decode())}=" for p in params
        ))
        qs = "&".join((f"{k}={v}" for k,v in pparams.items()))
        qs = f"&{qs}&{blob}" if qs and blob else f"&{qs}" if qs and not blob else f"&{blob}" if not qs and blob else ""
        return self._response(requests.get(f"{self._url}/get?_provider={provider}&_resource={resource}{qs}"))
