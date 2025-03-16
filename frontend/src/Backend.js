const BASE_API_URL = "http://localhost:8000"

export default class Backend {
  constructor() {
    this.base_url =  BASE_API_URL;
  }

  async resourceTypes() {
    return fetch(`${this.base_url}/resource-types`)
    .then(response => response.json());
  }

  async attributes(resourceType) {
    const [provider, resource] = resourceType.split(".")
    return fetch(`${this.base_url}/attributes?provider=${provider}&resource=${resource}`)
    .then(response => response.json());
  }

  async query(resourceType, labels) {
    const [provider, resource] = resourceType.split(".")
    const qs = (labels ?? []).map(l=> `${l.key}=${l.val}`).join("&")
    console.log(`${this.base_url}/get?provider=${provider}&resource=${resource}&${qs}`)
    return fetch(`${this.base_url}/get?provider=${provider}&resource=${resource}&${qs}`)
    .then(response => {
      this.throwForStatus(response);
      return response.json();
    })
    .then(response => {
      return response.results.map(result => {
          return {
            resourceType: resourceType,
            metadata: result.metadata,
            data: result.data
          };
      })
    })
  }

  throwForStatus(response) {
    const firstDigit = Math.floor(response.status/100)
    switch (firstDigit) {
      case 4:
      case 5:
        throw new Error(`HTTP status ${response.status} on request to backend.`)
    }
  }
}
