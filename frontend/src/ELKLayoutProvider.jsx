
import { createContext, useEffect, useState } from 'react';
import { useNodesInitialized, useReactFlow } from '@xyflow/react';
import ELK from 'elkjs/lib/elk.bundled.js';

const elk = new ELK();

export const ELKLayoutContext = createContext({ skipNextRecompute: () => {} });

export const ELKLayoutProvider = (props) => {
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
  
    try{
      const layoutOptions = {
        'elk.algorithm': 'layered',
        'elk.layered.spacing.nodeNodeBetweenLayers': 100,
        'elk.spacing.nodeNode': 80,
        'elk.direction': 'DOWN'
      };
      const nnodes = getNodes();
      const existingNodeIds = nnodes.map(n => n.id);
      const eedges = getEdges().filter(e => existingNodeIds.includes(e.source) && existingNodeIds.includes(e.target));
      const graph = {
        id: 'root',
        layoutOptions: layoutOptions,
        children: nnodes.map((node) => ({
          ...node,
          width: node.measured.width,
          height: node.measured.height,
        })),
        edges: eedges,
      };
    
      elk.layout(graph).then(({ children }) => {
        const newNodes = nnodes.map((node) => {
          const nodeWithPosition = children.filter(c => c.id === node.id)[0]
          const newNode = {
            ...node,
            position: {
              x: nodeWithPosition.x,
              y: nodeWithPosition.y
            },
          };
       
          return newNode;
        });
        setNodes(newNodes);
      });
    } catch(e) {}
  }, [recomputeTriggered, nodesInitialized, shouldSkipNextRecompute, getNodes, getEdges, setNodes, setEdges, fitBounds]);

  return (
    <ELKLayoutContext.Provider value={{ skipNextRecompute: setShouldSkipNextRecompute }}>
      {children}
    </ELKLayoutContext.Provider>
  );
};
