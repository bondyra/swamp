import base64
import json
import urllib
import requests
import sys

provider = sys.argv[1]
resource = sys.argv[2]

def _transform(l):
    k,o,v = l.split("+")
    bk = base64.b64encode(k.encode()).decode()
    bo = base64.b64encode(o.encode()).decode()
    bv = base64.b64encode(v.encode()).decode()
    sk = urllib.parse.quote(bk)
    so = urllib.parse.quote(bo)
    sv = urllib.parse.quote(bv)
    return f"{sk},{so},{sv}"

labels = [_transform(l) for l in sys.argv[3:]] if len(sys.argv) > 2 else []

print(json.dumps(requests.get(f"http://localhost:8000/get?_provider={provider}&_resource={resource}&{'&'.join(labels)}").json()))
