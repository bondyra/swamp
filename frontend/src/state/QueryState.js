import { create } from 'zustand';

import { randomString } from '../Utils';

export const useQueryStore = create((set) => ({
  vertices: [],
  addVertex: () => set((state) => ({ vertices: [...state.vertices, {id: randomString(8), labels: []}] })),
  removeVertex: (vertexId) => set((state) => ({ vertices: state.vertices.filter(v => v.id !== vertexId) })),
  updateVertex: (vertexId, data) => set((state) => ({ vertices: state.vertices.map(v => { if (v.id !== vertexId) return v; return {...v, ...data}; })})),
  addLabel: (vertexId) => set((state) => ({ vertices: [
    ...state.vertices.map(v => {
        if (v.id !== vertexId)
            return v;
        return {...v, labels: [...v.labels, {id: randomString(8), key:  "", val: ""}]}
    })
  ]})),
  setLabels:(vertexId, labels) => set((state) => ({ vertices: [
    ...state.vertices.map(v => {
        if (v.id !== vertexId)
            return v;
        return {...v, labels: labels}
    })
  ]})),
  removeLabel: (vertexId, labelId) => set((state) => ({ vertices: [
    ...state.vertices.map(v => {
        if (v.id !== vertexId)
            return v;
        return {...v, labels: v.labels.filter(l => l.id !== labelId)}
    })
  ]})),
  setVertices: (vv) => set((state) => ({vertices: vv}) ),
  links: [],
  addLink: () => set((state) => ({ links: [...state.links, {id: randomString(8)}] })),
  removeLink: (linkId) => set((state) => ({ links: state.links.filter(l => l.id !== linkId) })),
  updateLink: (linkId, data) => set((state) => ({ links: state.links.map(l => { if (l.id !== linkId) return l; return {...l, ...data}; })})),
  setLinks: (ll) => set((state) => ({links: ll}) ),
  parameters: [],
  addParameter: () => set((state) => ({ parameters: [...state.parameters, {id: randomString(8)}] })),
  removeParameter: (parameterId) => set((state) => ({ parameters: state.parameters.filter(p => p.id !== parameterId) })),
  updateParameter: (parameterId, data) => set((state) => ({ parameters: state.parameters.map(p => { if (p.id !== parameterId) return p; return {...p, ...data}; })})),
  triggered: false, // used to actually requery backend
  setTriggered: (val) => set((state) => ({triggered: val})),
  redisplay: false,  // used to trigger JQ queries on display
  setRedisplay: (val) => set((state) => ({redisplay: val})),
  fields: [],
  removeField: (fieldId) => set((state) => ({fields: state.fields.filter(f => f.id !== fieldId)})),
  updateField: (fieldId, data) => set((state) => ({fields: state.fields.map(f => {if (f.id !== fieldId) return f; else return {...f, ...data}})})),
  addField: (vertexId) => set((state) => ({fields: [...state.fields, {vertexId: vertexId, id: randomString(16)}]})),
  setFields: (ff) => set((state) => ({fields: ff}) ),
  // alerts
  alert: "Hello",
  setAlert: (alrt) => set((state) => ({alert: alrt}))
}));
