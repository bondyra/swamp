import React, { useCallback, useEffect, useRef } from 'react';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import {
  ReactFlow,
  addEdge,
  ConnectionLineType,
  useEdgesState,
  useNodesState
} from '@xyflow/react';
import '@xyflow/react/dist/base.css';
import { ReactFlowProvider, useReactFlow, MarkerType } from '@xyflow/react';
import { useQueryStore } from './QueryState';

import PreviewNode from './PreviewNode';
import { DagreLayoutProvider } from './DagreLayoutProvider';

const theme = createTheme({
  palette: {
  },
});

const nodeTypes = {
  previewNode: PreviewNode
};

const SwampPreview = () => {
  const reactFlowWrapper = useRef(null);
  const reactFlow = useReactFlow();
  const vertices = useQueryStore((state) => state.vertices);
  const links = useQueryStore((state) => state.links);
  // eslint-disable-next-line
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);

  // refresh preview
  useEffect(() => {
      reactFlow.setNodes(vertices.map(v => {
          return {
              id: v.id,
              position: {x: 0, y: 0},
              type: 'previewNode',
              data: {resourceType: v.resourceType ?? "null", selected: v.selected}
          };
      }));
      reactFlow.setEdges(links.filter(l => l.from && l.to).map(l => {
          return {
              type: 'straight',
              id: `${l.from}-${l.to}`,
              source: l.from,
              target: l.to, 
              style: {strokeWidth: 2, stroke: l.selected ? '#aaaaff' : '#3e3e3e'},
              markerEnd: {
                type: MarkerType.Arrow,
                strokeWidth: 2,
                width: 8,
                height: 8,
                color: l.selected ? '#aaaaff' : '#3e3e3e'
              },
          };
      }));
  }, [reactFlow, vertices, links]);

  // RF stuff
  const onConnect = useCallback(
    (params) =>
      setEdges((eds) =>
        addEdge(
          { ...params, type: ConnectionLineType.Straight, animated: true },
          eds,
        ),
      ),
    [setEdges],
  );

  return (
    <div className="wrapper" ref={reactFlowWrapper} style={{ width: 'full', height: '128px' }}>
      <ThemeProvider theme={theme}>
        <ReactFlow
        nodes={nodes}
        edges={edges}
        deleteKeyCode={null}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        connectionLineType={ConnectionLineType.Straight}
        nodeTypes={nodeTypes}
        edgesFocusable={false}
        nodesDraggable={false}
        nodesConnectable={false}
        nodesFocusable={false}
        elementsSelectable={false}
        zoomOnDoubleClick={false}
        colorMode={"dark"}
        fitView
        proOptions={{ hideAttribution: true }}
        >
        </ReactFlow>
      </ThemeProvider>
    </div>
  );
}

const PreviewFlow = () => {
  return (
    <ReactFlowProvider>
      <DagreLayoutProvider skipInitialLayout>
        <SwampPreview/>
      </DagreLayoutProvider>
    </ReactFlowProvider>
  );
}

export default PreviewFlow;
