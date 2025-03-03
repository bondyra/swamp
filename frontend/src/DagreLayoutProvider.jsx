
import { createContext, useEffect, useState } from 'react';
import { useNodesInitialized, useReactFlow } from '@xyflow/react';
import dagre from '@dagrejs/dagre';


const dagreGraph = new dagre.graphlib.Graph().setDefaultEdgeLabel(() => ({}));


const getLayoutedElements = (nodes, edges) => {
  dagreGraph.setGraph({ rankdir: 'TB' });
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
 
  return { nodes: newNodes, edges: edges, height: dagreGraph.graph().height, width: dagreGraph.graph().width };
};

export const DagreLayoutContext = createContext({ skipNextRecompute: () => {} });

export const DagreLayoutProvider = (props) => {
  const { skipInitialLayout, children } = props;
  const [shouldSkipNextRecompute, setShouldSkipNextRecompute] = useState(skipInitialLayout);
  const [recomputeTriggered, setRecomputeTriggered] = useState(false);
  const nodesInitialized = useNodesInitialized();
  const { getNodes, getEdges, setNodes, setEdges, fitBounds } = useReactFlow();

  // Trigger recompute when nodes uninitialized
  useEffect(() => {
    if (!nodesInitialized) {
      setRecomputeTriggered(true);
    }
  }, [nodesInitialized, setRecomputeTriggered]);

  // If this recompute skipped, reset the flag for the next recompute
  useEffect(() => {
    if (shouldSkipNextRecompute && recomputeTriggered && nodesInitialized) {
      setShouldSkipNextRecompute(false);
      setRecomputeTriggered(false);
    }
  }, [shouldSkipNextRecompute, nodesInitialized, recomputeTriggered, setShouldSkipNextRecompute, setRecomputeTriggered]);

  // If this recompute was not skipped, run the recomputation
  useEffect(() => {
    if (!(recomputeTriggered && nodesInitialized && !shouldSkipNextRecompute)) {
      return;
    }

    setRecomputeTriggered(false);

    // getLayoutedElements calls dagre to compute the layout. height and width correspond to 
    // dagreGraph.graph().height and width, which allows us to fit the viewport to the layout
    const { nodes, edges, height, width } = getLayoutedElements(getNodes(), getEdges());

    setEdges(edges.map(e => ({ ...e, hidden: false })));
    setNodes(
      nodes.map(n => ({
        ...n,
        data: {
          ...n.data,
          isHidden: false
        }
      }))
    );
    console.log(`fitBounds width: ${width} height: ${height} x: 0 y: 0`)
    fitBounds({ width: width ?? 0, height: height ?? 0, x: 0, y: 0 }, { duration: 100 });
  }, [recomputeTriggered, nodesInitialized, shouldSkipNextRecompute, getNodes, getEdges, setNodes, setEdges, fitBounds]);

  return (
    <DagreLayoutContext.Provider value={{ skipNextRecompute: setShouldSkipNextRecompute }}>
      {children}
    </DagreLayoutContext.Provider>
  );
};
