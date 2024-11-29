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

  const scan = useCallback(
  () => {
    var newNodes = [];
    var newEdges = [];
    const mod = 'aws';
    const promises = ['vpc', 'subnet', 'rtb', 'igw', 'sg', 'nat', 'eip', 'eni', 'nacl'].map(function(resource_type) {
      return fetch(`http://localhost:8000/${mod}/${resource_type}`)
        .then(response => response.json())
        .then(response => {
          response.results.forEach(function(result) {
            const newId = `${mod}.${resource_type}.${result.id}`;
            newNodes.push({
              id: newId,
              position: { x: 0, y: 0 },
              type: 'resource',
              data: {
                id: newId,
                module: mod,
                resource_type: resource_type,
                resource_id: result.id,
                obj: result.obj,
                icon: `./icons/aws/${resource_type}.svg` 
              },
            });
            result.parents.forEach(function(parent) {
              const parentId = `${parent.module}.${parent.resource_type}.${parent.id}`;
              newEdges.push({id: `${parentId}-${newId}`, source: parentId, target: newId, style: {strokeWidth: 5} });
            return newNodes;
          })
        })
      })
    });
    Promise.all(promises).then(() => {
      console.log(`Adding a total of ${newNodes.length} nodes and ${newEdges.length} edges`);
      console.log(newNodes);
      console.log(newEdges);
      setNodes(newNodes);
      setEdges(newEdges);
    });
  }, []
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
          <button onClick={() => scan()}>PLZ CLICK ME TO FETCH STUFF</button>
      </Panel>

      <Panel position="top-right">
        <button onClick={() => onLayout()}>PLZ CLICK ME TO ORDER STUFF</button>
      </Panel>
        <Panel position="top-left">
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
