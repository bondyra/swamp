import React, { useState, useCallback } from 'react';
import {
  Background,
  BackgroundVariant,
  ReactFlow,
  addEdge,
  ConnectionLineType,
  useEdgesState,
  useNodesState,
  MiniMap,
  Controls,
  Panel,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import dagre from '@dagrejs/dagre';

import Resource from './Resource';

const dagreGraph = new dagre.graphlib.Graph().setDefaultEdgeLabel(() => ({}));

const getLayoutedElements = (nodes, edges, direction = 'TB') => {
  dagreGraph.setGraph({ rankdir: direction });
  nodes.forEach((node) => {
    dagreGraph.setNode(node.id, { width: node.measured.width, height: node.measured.height });
  });
 
  edges.forEach((edge) => {
    dagreGraph.setEdge(edge.source, edge.target);
  });
 
  dagre.layout(dagreGraph);

  const newNodes = nodes.map((node) => {
    const nodeWithPosition = dagreGraph.node(node.id);
    const newNode = {
      ...node,
      targetPosition: 'top',
      sourcePosition: 'bottom',
      position: {
        x: nodeWithPosition.x - node.measured.width / 2,
        y: nodeWithPosition.y - node.measured.height / 2,
      },
    };
 
    return newNode;
  });
 
  return { nodes: newNodes, edges };
};

const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(
  [],
  [],
);

const nodeTypes = {
  resource: Resource
};

function App() {
  const [nodes, setNodes, onNodesChange] = useNodesState(layoutedNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(layoutedEdges);
  const [colorMode, setColorMode] = useState('dark');
  const onConnect = useCallback(
    (params) =>
      setEdges((eds) =>
        addEdge(
          { ...params, type: ConnectionLineType.SmoothStep, animated: true },
          eds,
        ),
      ),
    [],
  );

  const onLayout = useCallback(
    (direction) => {
      const { nodes: layoutedNodes, edges: layoutedEdges } =
        getLayoutedElements(nodes, edges, direction);
 
      setNodes([...layoutedNodes]);
      setEdges([...layoutedEdges]);
    },
    [nodes, edges],
  );

  const query = useCallback(
    () => {
      var newNodes = [];
      var newEdges = [];
      const q = document.getElementById('main_query').value;
      const [mod, resource_type] = q.split('.');
      fetch(`http://localhost:8000/${mod}/${resource_type}?pattern=`)
        .then(response => response.json())
        .then(response => {
          console.log(response)
          response.results.forEach(function(result) {
            newNodes.push({
              id: `${mod}.${resource_type}.${result.id}`,
              position: { x: 0, y: 0 },
              type: 'resource',
              data: {
                id: `${mod}.${resource_type}.${result.id}`,
                module: mod,
                resource_type: resource_type,
                resource_id: result.id,
                obj: result.obj,
                icon: `./icons/aws/${resource_type}.svg` 
              },
            });
          })
        })
        .then(() => {
          setNodes(newNodes);
          setEdges(newEdges);
        });
    }
  );

  const onChange = (evt) => {
    setColorMode(evt.target.value);
  };

  return (
    <div style={{ width: '100vw', height: '100vh' }}>
      <ReactFlow 
      nodes={nodes}
      edges={edges}
      onNodesChange={onNodesChange}
      onEdgesChange={onEdgesChange}
      onConnect={onConnect}
      connectionLineType={ConnectionLineType.SmoothStep}
      nodeTypes={nodeTypes}
      colorMode={colorMode}
      fitView
      >

      <Panel position='top-center'>
          <textarea placeholder="Enter query" id="main_query" />
          <button onClick={() => query()}>QUERY</button>
      </Panel>

      <Panel position="top-left">
        <button onClick={() => onLayout()}>Refresh layout</button>
      </Panel>
        <Panel position="top-right">
          <select onChange={onChange} data-testid="colormode-select">
            <option value="dark">dark</option>
            <option value="light">light</option>
            <option value="system">system</option>
          </select>
        </Panel>
        <MiniMap />
        <Controls />
        <Background variant={BackgroundVariant.Lines} color='#222222'  />
      </ReactFlow>
    </div>
  );
}

export default App;
