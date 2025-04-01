const BASE_API_URL = "http://localhost:8000"

export default class Backend {
  constructor() {
    this.base_url =  BASE_API_URL;
  }

  async resourceTypes() {
    return fetch(`${this.base_url}/resource-types`)
    .then(response => response.json())
    .then(response => response.map(result => {
      return {
        value: `${result.provider}.${result.resource}`,
        description: result.description
      }
    }))
  }

  async attributes(resourceType) {
    const [provider, resource] = resourceType.split(".")
    return fetch(`${this.base_url}/attributes?_provider=${provider}&_resource=${resource}`)
    .then(response => response.json());
  }

  async attributeValues(resourceType, attribute, params) {
    const [provider, resource] = resourceType.split(".")
    const paramsQs = params.map(x => `${x.key}=${x.val}`).join("&")
    return fetch(`${this.base_url}/attribute-values?_provider=${provider}&_resource=${resource}&attribute=${attribute}&${paramsQs}`)
    .then(response => response.json());
  }

  async linkSuggestion(childResourceType, parentResourceType) {
    const [childProvider, childResource] = childResourceType.split(".")
    const [parentProvider, parentResource] = parentResourceType.split(".")
    return fetch(`${this.base_url}/link-suggestion?child_provider=${childProvider}&child_resource=${childResource}&parent_provider=${parentProvider}&parent_resource=${parentResource}`)
    .then(response => response.json());
  }

  async query(resourceType, labels) {
    const [provider, resource] = resourceType.split(".");
    const qs = (labels ?? []).map(l => `${encodeURIComponent(btoa(l.key))},${encodeURIComponent(btoa(l.op ?? "eq"))},${encodeURIComponent(btoa(l.val))}`).join("&");
    return fetch(`${this.base_url}/get?_provider=${provider}&_resource=${resource}&${qs}`)
      .then(response => {
        this.throwForStatus(response);
        return response.json();
      })
      .then(response => {
        return response.results.map(result => {
            return {
              resourceType: resourceType,
              result: result
            };
        })
      });
  }

  throwForStatus(response) {
    const firstDigit = Math.floor(response.status/100)
    switch (firstDigit) {
      case 4:
      case 5:
        throw new Error(`HTTP status ${response.status} on request to backend.`)
      default:
        break;
    }
  }
}
