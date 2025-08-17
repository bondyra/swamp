
import React, { createContext, useEffect, useState } from 'react';
import { 
  useNodesInitialized, useReactFlow
} from '@xyflow/react';
import {
  forceSimulation,
  forceLink,
  forceManyBody,
  forceX,
  forceY,
} from 'd3-force';
import { collide } from './collide.js';

import '@xyflow/react/dist/style.css';


const simulation = forceSimulation()
  .force('charge', forceManyBody().strength(-1000))
  .force('x', forceX().x(0).strength(0.05))
  .force('y', forceY().y(0).strength(0.05))
  .force('collide', collide())
  .alphaTarget(0.05)
  .stop();

export const D3ForceAsyncLayoutContext = createContext({ skipNextRecompute: () => {} });

export const D3ForceAsyncLayoutProvider = (props) => {
  const { skipInitialLayout, children } = props;
  const [shouldSkipNextRecompute, setShouldSkipNextRecompute] = useState(skipInitialLayout);
  const [recomputeTriggered, setRecomputeTriggered] = useState(false);
  const nodesInitialized = useNodesInitialized();
  const { getNodes, getEdges, setNodes, setEdges, fitBounds, fitView } = useReactFlow();
  const [running, setRunning] = useState(true);

  useEffect(() => {
    (async () => {
      await new Promise(res => setTimeout(res, 5000));
      setRunning(false)
    })();
  }, [running, setRunning]);

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
      //////////////////
        console.log("ADDADSDDASSDDASADSSDSDSDASDDDSDSDSD")
        console.log(running);
  let nodes = getNodes().map((node) => ({
    ...node,
    x: node.position.x,
    y: node.position.y,
  }));
        console.log(nodes);
  let edges = getEdges().map((edge) => edge);

  // If React Flow hasn't initialized our nodes with a width and height yet, or
  // if there are no nodes in the flow, then we can't run the simulation!
  if (!running || nodes.length === 0) return;

  simulation.nodes(nodes).force(
    'link',
    forceLink(edges)
      .id((d) => d.id)
      .strength(0.05)
      .distance(100),
  );
  // The tick function is called every animation frame while the simulation is
  // running and progresses the simulation one step forward each time.
  const tick = () => {
    simulation.tick();
    setNodes(
      nodes.map((node) => ({
        ...node,
        position: { x: node.x, y: node.y },
      })),
    );

    window.requestAnimationFrame(() => {
      // Give React and React Flow a chance to update and render the new node
      // positions before we fit the viewport to the new layout.
      fitView();
      // If the simulation hasn't been stopped, schedule another tick.
      if (running) tick();
    });
  };
  window.requestAnimationFrame(tick);

      //////////////////
      // const nodes = getLayoutedElements(getNodes(), getEdges());

      // setNodes(
      //   nodes.map(n => ({
      //     ...n,
      //     data: {
      //       ...n.data,
      //       isHidden: false
      //     }
      //   }))
      // );
      // fitBounds({ width: width ?? 0, height: height ?? 0, x: 0, y: 0 }, { duration: 100 });
    } catch {}
  }, [recomputeTriggered, nodesInitialized, shouldSkipNextRecompute, getNodes, getEdges, setNodes, setEdges, fitBounds, running, fitView]);

  return (
    <D3ForceAsyncLayoutContext.Provider value={{ skipNextRecompute: setShouldSkipNextRecompute }}>
      {children}
    </D3ForceAsyncLayoutContext.Provider>
  );
};
