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

  async example(resourceType) {
    const [provider, resource] = resourceType.split(".")
    return fetch(`${this.base_url}/example?_provider=${provider}&_resource=${resource}`)
    .then(response => response.json());
  }

  async attributeValues(resourceType, attribute, params) {
    const [provider, resource] = resourceType.split(".")
    const paramsQs = params.map(x => `${x.key}=${x.val}`).join("&")
    return fetch(`${this.base_url}/attribute-values?_provider=${provider}&_resource=${resource}&attribute=${attribute}&${paramsQs}`)
    .then(response => response.json());
  }

  async suggestion(from, to) {
    if (from === "aws.vpc" && to === "aws.subnet")
      return {fromAttr: ".VpcId", op: "=", toAttr: ".VpcId"}
    if (from === "aws.vpc" && to === "aws.route_table")
      return {fromAttr: ".VpcId", op: "=", toAttr: ".VpcId"}
    if (from === "aws.vpc" && to === "aws.network_acl")
      return {fromAttr: ".VpcId", op: "=", toAttr: ".VpcId"}
    if (from === "aws.network_acl" && to === "aws.subnet")
      return {fromAttr: ".Associations[].SubnetId", op: "contains", toAttr: ".SubnetId"}
    if (from === "aws.subnet" && to === "aws.network_acl")
      return {fromAttr: ".SubnetId", op: "contains", toAttr: ".Associations[].SubnetId"}
    return null;
  }

  async* query(vertices) {
    // todo: error handling
    const promises = vertices.filter(v => v.resourceType).map(v => {
      return this.queryOne(v.resourceType, v.labels, v.id);
    });
    yield* this.asCompleted(promises);
  }

  async queryOne(resourceType, labels, vertexId) {
    const [provider, resource] = resourceType.split(".");
    const qs = (labels ?? []).map(l => `${encodeURIComponent(btoa(l.key))},${encodeURIComponent(btoa(l.op ?? "=="))},${encodeURIComponent(btoa(l.val))}`).join("&");
    return fetch(`${this.base_url}/get?_provider=${provider}&_resource=${resource}&${qs}`)
      .then(async response => {
        await this.throwForStatus(response);
        return response.json();
      }) 
      .then(response => {
        return response.results.map(result => {
            return {
              resourceType: resourceType,
              result: result,
              vertexId: vertexId,
            };
        })
      });
  }

  async* asCompleted(promises) {
    const pending = new Set(promises.map(async (p, i) => {
      const items = await p;
      return { items, index: i };
    }));

    while (pending.size > 0) {
      const { items, index } = await Promise.race(pending);
      
      for (const p of pending) {
        p.then(result => {
          if (result.index === index) {
            pending.delete(p);
          }
        }).catch(() => {
          pending.delete(p);
        });
      }

      for (const item of items) {
        yield item;
      }
    }
  }

  async throwForStatus(response) {
    const firstDigit = Math.floor(response.status/100)
    switch (firstDigit) {
      case 4:
      case 5:
        const rr = await response.json();
        throw new Error(`HTTP status ${response.status} on request to backend: ${rr.detail}.`);
      default:
        break;
    }
  }
}
