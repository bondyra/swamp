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
  mapLabels:(vertexId, f) => set((state) => ({ vertices: [
    ...state.vertices.map(v => {
        if (v.id !== vertexId)
            return v;
        return {...v, labels: v.labels.map(x => f(x))}
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
  triggered: false, // used to actually requery backend
  setTriggered: (val) => set((state) => ({triggered: val})),
  redisplay: false,  // used to trigger JQ queries on display
  setRedisplay: (val) => set((state) => ({redisplay: val})),
  fields: [],
  removeField: (fieldId) => set((state) => ({fields: state.fields.filter(f => f.id !== fieldId)})),
  updateField: (fieldId, data) => set((state) => ({fields: state.fields.map(f => {if (f.id !== fieldId) return f; else return {...f, ...data}})})),
  addField: (vertexId) => set((state) => ({fields: [...state.fields, {vertexId: vertexId, id: randomString(16), val: ""}]})),
  setFields: (ff) => set((state) => ({fields: ff}) ),
  savedLabels: [],
  saveLabel: (label) => set((state) => ({savedLabels: label.op && label.val ? [...state.savedLabels.filter(l => l.key !== label.key), label] : state.savedLabels}))
}));
