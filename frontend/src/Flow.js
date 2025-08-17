import React, { useCallback, useEffect, useState } from 'react';
import { Panel, useReactFlow } from '@xyflow/react';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import { Button } from '@mui/material';
import {
  ReactFlow,
  addEdge,
  ConnectionLineType,
  useEdgesState,
  useNodesState,
  MiniMap,
  Controls,
  ReactFlowProvider
} from '@xyflow/react';
import '@xyflow/react/dist/base.css';

import Resource from './Resource';
import Query from './Query';
import { DagreLayoutProvider } from './DagreLayoutProvider';
import { randomString } from './Utils'
import { useBackend } from './BackendProvider';
import { useQueryStore } from './QueryState';
import { JSONPath } from 'jsonpath-plus';

const theme = createTheme({
  palette: {},
});

const initialNodes = [
  {id: "1", position: {x: 0, y: 0}, type: 'resource', data: {id: "1", vertexId: "a", resourceType: "twoja.stara1", result: {a: 1, b: 2}}},
  {id: "2", position: {x: 0, y: 0}, type: 'resource', data: {id: "2", vertexId: "b", resourceType: "twoja.stara2", result: {a: 1, b: 2}}},
  {id: "3", position: {x: 0, y: 0}, type: 'resource', data: {id: "3", vertexId: "c", resourceType: "twoja.stara3", result: {a: 1, b: 2}}},
  {id: "4", position: {x: 0, y: 0}, type: 'resource', data: {id: "4", vertexId: "d", resourceType: "twoja.stara4", result: {a: 1, b: 2}}},
]

const initialEdges = [
  {id: `1-2`, source: "1", target: "2", style: {strokeWidth: 5} },
  {id: `2-3`, source: "2", target: "3", style: {strokeWidth: 5} },
  {id: `2-4`, source: "2", target: "4", style: {strokeWidth: 5} },
]

const nodeTypes = {
  resource: Resource,
  query: Query,
};

const linksToMap = (links, keyfun) => {
  const res = {};
  for (const l of links) {
    const key = keyfun(l);
    if (!key)  // ignore not configured links
      continue;
    if (!res[key]) {
      res[key] = [];
    }
    res[key].push(l);
  }
  return res
}

const SwampFlow = () => {
  const backend = useBackend();
  const reactFlow = useReactFlow();
  const [addDummyNode, setAddDummyNode] = useState(false)
  const [delDummyNode, setDelDummyNode] = useState(false)
  // eslint-disable-next-line
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);

  const vertices = useQueryStore((state) => state.vertices);
  const updateVertex = useQueryStore((state) => state.updateVertex);
  const links = useQueryStore((state) => state.links);
  const triggered = useQueryStore((state) => state.triggered);
  const setTriggered = useQueryStore((state) => state.setTriggered);

  useEffect(() => {
    async function update() {
      if (!triggered) {
        return;
      }
      var allNodes = [...nodes];
      var linkFromMap = linksToMap(links, (l) => l.fromVertexId);
      var linkToMap = linksToMap(links, (l) => l.toVertexId);
      for await (const item of backend.query(vertices)) {
        const id = `${item.resourceType}.${item.result._id}`
        const newNode = {
          id: id,
          position: {
            x: 0,
            y: 0
          },
          type: 'resource',
          data: {
            id: id,
            vertexId: item.vertexId,
            resourceType: item.resourceType,
            result: item.result
          },
        };
        allNodes.push(newNode);
        setNodes(oldNodes => [...oldNodes, newNode]);
        if (item.vertexId in linkFromMap){
          const linksTo = linkFromMap[item.vertexId];
          for (const lt of linksTo) {
            const potentialToNodes = allNodes.filter(n => lt.toVertexId === n.data.vertexId);
            if (potentialToNodes){
              const fromValue = JSONPath({path: lt.fromAttr, json: item.result});
              // todo: op, for now it's hardcoded to "eq"
              potentialToNodes.filter(n => JSONPath({path: lt.toAttr, json: n.data.result}) === fromValue).forEach(n => {
                setEdges(oe => [...oe, {id: `${id}-${n.id}`, source: id, target: n.id, style: {strokeWidth: 5} }]);
              });
            }
          }
        }
        if (item.vertexId in linkToMap){
          const linksFrom = linkToMap[item.vertexId];
          for (const lf of linksFrom) {
            const potentialFromNodes = allNodes.filter(n => lf.fromVertexId === n.data.vertexId);
            if (potentialFromNodes){
              const toValue = JSONPath({path: lf.toAttr, json: item.result})[0];
              // todo: op, for now it's hardcoded to "eq"
              potentialFromNodes.filter(n => JSONPath({path: lf.fromAttr, json: n.data.result})[0] === toValue).forEach(n => {
                setEdges(oe => [...oe, {id: `${n.id}-${id}`, source: n.id, target: id, style: {strokeWidth: 5} }]);
              });
            }
          }
        }
      }
    }
    setTriggered(false);
    update();
  }, [backend, triggered, setTriggered]);

  // RF stuff
  const onConnect = useCallback(
    (params) =>
      setEdges((eds) =>
        addEdge(
          { ...params, type: ConnectionLineType.Step, animated: true },
          eds,
        ),
      ),
    [setEdges],
  );

  const onNodesChangeExt = useCallback((changes) => {
    onNodesChange(changes);
    if (changes.some(c=> c.type === "dimensions")){
      // force layout by adding a dummy node
      setAddDummyNode(true)
    }
    changes.filter(c=> c.type === "select").forEach(c => {
      const node = nodes.filter(n => n.id === c.id)[0];
      updateVertex(node.data.vertexId, {selected: c.selected});
    })
  }, [onNodesChange, setAddDummyNode, updateVertex, nodes]);

  useEffect(() => {
    if (delDummyNode) {
      reactFlow.deleteElements({nodes: nodes.filter(n => n.id.startsWith("__DUMMY__"))})
      setDelDummyNode(false)
    }
    if (addDummyNode) {
      reactFlow.addNodes([{id: `__DUMMY__${randomString(8)}`, position: { x: 0, y: 0 }}])
      setAddDummyNode(false)
      setDelDummyNode(true)  // immediately mark this node for deletion
    }
  }, [reactFlow, nodes, addDummyNode, delDummyNode, setAddDummyNode, setDelDummyNode])

  return (
    <ThemeProvider theme={theme}>
        <ReactFlow
        nodes={nodes}
        edges={edges}
        deleteKeyCode={null}
        onNodesChange={onNodesChangeExt}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        connectionLineType={ConnectionLineType.Step}
        nodeTypes={nodeTypes}
        nodesDraggable={false}
        colorMode={"dark"}
        fitView
        proOptions={{ hideAttribution: true }}
        >
            <MiniMap />
            <Controls />
          <Panel position='top'>
            <Button onClick={() => setAddDummyNode(true)}>aaaaaaa</Button>
          </Panel>
        </ReactFlow>
    </ThemeProvider>
  );
}

const Flow = () => {
  return (
    <ReactFlowProvider>
      <DagreLayoutProvider skipInitialLayout>
        <SwampFlow/>
      </DagreLayoutProvider>
    </ReactFlowProvider>
  );
}

export default Flow;
