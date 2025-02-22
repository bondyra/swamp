import React, { useCallback, useRef } from 'react';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import {
  ReactFlow,
  addEdge,
  ConnectionLineType,
  useEdgesState,
  useNodesState,
  useReactFlow,
  MiniMap,
  Controls,
  ReactFlowProvider
} from '@xyflow/react';
import '@xyflow/react/dist/base.css';

import Resource from './Resource';
import Query from './Query';
import { DagreLayoutProvider } from './DagreLayoutProvider';

const theme = createTheme({
  palette: {
  },
});

let queryId = 1;
const newQueryId = () => `query-${queryId++}`;

const initialNodes = [
  {
    id: newQueryId(),
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
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const { screenToFlowPosition } = useReactFlow();
  // const [colorMode, setColorMode] = useState('dark');
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

  const onConnectEnd = useCallback(
    (event, connectionState) => {
      // when a connection is dropped on the pane it's not valid
      if (!connectionState.isValid) {
        const sourceNodeId = connectionState.fromNode.id
        const sourceNode = nodes.filter(n=> n.id === sourceNodeId)[0]
        const targetNodeId = newQueryId()
        // we need to remove the wrapper bounds, in order to get the correct position
        const { clientX, clientY } = 'changedTouches' in event ? event.changedTouches[0] : event;
        
        const newNode = {
          id: targetNodeId,
          position: screenToFlowPosition({
            x: clientX,
            y: clientY,
          }),
          type: 'query',
          data: {resourceType: null, labels: [], nodeData: sourceNode.data},
          origin: [0.5, 0.0],
        };
 
        setNodes((nds) => nds.concat(newNode));
        setEdges((eds) =>
          eds.concat({ id: `${sourceNodeId}-${targetNodeId}`, source: sourceNodeId, target: targetNodeId, style: {strokeWidth: 5} }),
        );
      }
    },
    [screenToFlowPosition, setNodes, setEdges, nodes],
  );

  // const onChange = (evt) => {
  //   setColorMode(evt.target.value);
  // };

  return (
    <div className="wrapper" ref={reactFlowWrapper} style={{ width: '100vw', height: '100vh' }}>
    <ThemeProvider theme={theme}>
      <ReactFlow 
      nodes={nodes}
      edges={edges}
      onNodesChange={onNodesChange}
      onEdgesChange={onEdgesChange}
      onConnect={onConnect}
      onConnectEnd={onConnectEnd}
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

const App = () => (
  <ReactFlowProvider>
    <DagreLayoutProvider skipInitialLayout>
      <SwampApp />
    </DagreLayoutProvider>
  </ReactFlowProvider>
);

export default App;
