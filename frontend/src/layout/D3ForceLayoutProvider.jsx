
import { createContext, useEffect, useState } from 'react';
import { useNodesInitialized, useReactFlow } from '@xyflow/react';
import {
  forceSimulation,
  forceLink,
  forceManyBody,
  forceX,
  forceY,
} from 'd3-force';
import { collide } from '../collide.js';

const simulation = forceSimulation()
  .force('charge', forceManyBody().strength(-1000))
  .force('x', forceX().x(0).strength(0.05))
  .force('y', forceY().y(0).strength(0.05))
  .force('collide', collide())
  .alphaTarget(0.01)
  .stop();


const getLayoutedElements = (nodes, edges) => {
  let nds = nodes.map((node) => ({
    ...node,
    x: node.position.x,
    y: node.position.y,
  }));
  let eds = edges.map((edge) => edge);

  simulation.nodes(nds).force(
    'link',
    forceLink(eds)
      .id((d) => d.id)
      .strength(0.05)
      .distance(100)
  ).tick(1000).stop();
  
  const newNodes =  nds.map((node) => ({
    ...node,
    position: { x: node.fx ?? node.x, y: node.fy ?? node.y },
  }));
  return newNodes;
};

export const D3ForceLayoutContext = createContext({ skipNextRecompute: () => {} });

export const D3ForceLayoutProvider = (props) => {
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
    try{
      const nodes = getLayoutedElements(getNodes(), getEdges());
      console.log("RE-LAYOUT")
      setNodes(
        nodes.map(n => ({
          ...n,
          data: {
            ...n.data,
            isHidden: false
          }
        }))
      );
      // fitBounds({ width: width ?? 0, height: height ?? 0, x: 0, y: 0 }, { duration: 100 });
    } catch {}
  }, [recomputeTriggered, nodesInitialized, shouldSkipNextRecompute, getNodes, getEdges, setNodes, setEdges, fitBounds]);

  return (
    <D3ForceLayoutContext.Provider value={{ skipNextRecompute: setShouldSkipNextRecompute }}>
      {children}
    </D3ForceLayoutContext.Provider>
  );
};
