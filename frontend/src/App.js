import React, { useCallback, useEffect, useRef, useState } from 'react';
import { useReactFlow } from '@xyflow/react';
import { createTheme, ThemeProvider } from '@mui/material/styles';
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


let dummyId = 1;
const newDummyId = () => `${dummyId++}`;

const theme = createTheme({
  palette: {
  },
});

const initialNodes = [
  {
    id: 'root-query',
    position: { x: 0, y: 0 },
    type: 'query',
    data: {resourceType: null, resourceOpen: true, labels: [], nodeData: null},
  }
]

const nodeTypes = {
  resource: Resource,
  query: Query,
};

const SwampApp = () => {
  const reactFlowWrapper = useRef(null);
  const reactFlow = useReactFlow();
  const [addDummyNode, setAddDummyNode] = useState(false)
  const [delDummyNode, setDelDummyNode] = useState(false)
  // eslint-disable-next-line
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
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
  }, [onNodesChange, setAddDummyNode])

  useEffect(() => {
    if (delDummyNode) {
      reactFlow.deleteElements({nodes: nodes.filter(n => n.id.startsWith("__DUMMY__"))})
      setDelDummyNode(false)
    }
    if (addDummyNode) {
      reactFlow.addNodes([{id: `__DUMMY__${newDummyId()}`, position: { x: 0, y: 0 }}])
      setAddDummyNode(false)
      setDelDummyNode(true)  // immediately mark this node for deletion
    }
  }, [reactFlow, nodes, addDummyNode, delDummyNode, setAddDummyNode, setDelDummyNode])

  return (
    <div className="wrapper" ref={reactFlowWrapper} style={{ width: '100vw', height: '100vh' }}>
      <ThemeProvider theme={theme}>
        <ReactFlow 
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChangeExt}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        connectionLineType={ConnectionLineType.Step}
        nodeTypes={nodeTypes}
        nodesDraggable={false}
        colorMode={"dark"}
        fitView
        >
          {/* <Panel position="top-left">
            <select onChange={onChange} data-testid="colormode-select">
              <option value="dark">dark</option>
              <option value="light">light</option>
              <option value="system">system</option>
            </select>
          </Panel> */}
          <MiniMap />
          <Controls />
        </ReactFlow>
      </ThemeProvider>
    </div>
  );
}

const App = () => {
  return (
    <ReactFlowProvider>
      <DagreLayoutProvider skipInitialLayout>
        <SwampApp/>
      </DagreLayoutProvider>
    </ReactFlowProvider>
  );
}


export default App;
