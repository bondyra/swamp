// jqLoader.js
import jqModule from "jq-web";

let jqInstance = null;

export async function getJq() {
  if (!jqInstance) {
    jqInstance = await jqModule;
  }
  return jqInstance;
}
