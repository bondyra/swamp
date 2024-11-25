import React, { useState, useCallback } from 'react';
import {
  Background,
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
    const q = document.getElementById('queryy').value;
    const things = q.split('\n');
    var promises = [];
    var results = [];
    things.forEach(thing => {
      const [mod, resource_type] = thing.trim().split('.');
      promises.push(
        fetch('http://localhost:8000/' + mod + '/' + resource_type)
        .then(response => response.json())
        .then(r => {
          results.push({module: mod, resource_type: resource_type, results: r});
        })
      );
    })
    var newNodes = [];
    Promise.all(promises).then(() => {
      results.forEach(r => {

        Object.keys(r.results).forEach(function(key, index) {
          console.log(key, r.results[key]);
          newNodes.push({
            id: r.module + '.' + r.resource_type + '.' + key,
            position: { x: 0, y: 0 },
            type: 'resource',
            data: { 
              id: key, module: r.module, resource_type: r.resource_type, 
              content: r.results[key], 
              icon: "./icons/aws/" + r.resource_type + ".svg" 
            }
          }); 
        });
      });
      setNodes(newNodes);
    });
  });

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
          <textarea placeholder="Enter query" id="queryy" />
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
        <Background  />
      </ReactFlow>
    </div>
  );
}

export default App;
