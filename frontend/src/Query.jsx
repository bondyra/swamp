import React, { memo, useEffect, useCallback, useState } from 'react';
import { Handle, Position, useReactFlow } from '@xyflow/react';

import QueryWizard from './QueryWizard';

//todo collapse
export default memo(({ id, data, isConnectable }) => {
  const reactFlow = useReactFlow();
    const [previousLabelVars, setPreviousLabelsVars] = useState(null);

  const addNewNodesAndEdges = (items) => {
    var newNodes = [];
    var newEdges = [];
    const currNode = reactFlow.getNode(id);
    // todo: move it to querywizard and add link label if it's a join query
    items.forEach(item => {
      const newNodeId = `${item.resourceType}.${item.result._id}`;
      newNodes.push({
        id: newNodeId,
        position: { 
          x: currNode.position.x ?? 0, 
          y: currNode.position.y + 100 ?? 0 
        },
        type: 'resource',
        data: {
          id: newNodeId,
          resourceType: item.resourceType,
          inline: {},
          result: item.result
        },
      });
      newEdges.push({id: `${id}-${newNodeId}`, source: id, target: newNodeId, style: {strokeWidth: 5} });
    });
    reactFlow.addNodes(newNodes);
    reactFlow.addEdges(newEdges);
  };

  const updateResourceType = useCallback(
    (newValue) => {
      reactFlow.updateNodeData(id, (node) => {
          return {
            ...node.data,
            resourceType: newValue
          };
      });
    }, [id, reactFlow]);

  const setLabels = useCallback(
    (labels) => {
      reactFlow.updateNodeData(id, (node) => {
        return {...node.data, labels: labels};
      })
    }, [id, reactFlow]);
  
  useEffect(() => {
    var result = new Map();
    reactFlow.getNodes().forEach(n => {
      (n.data.labels ?? []).forEach(l => {
        if(! result[l.key])
          result[l.key] = new Set()
        result[l.key].add(l.val)
      })
    })
    setPreviousLabelsVars(result)
  }, [setPreviousLabelsVars, reactFlow])

  return (
    <>
      <div className="wrapper">
        <div className="inner">
          <div className="body">
            <QueryWizard 
            nodeId={id} resourceType={data.resourceType} labels={data.labels} 
            doSomethingWithResults={addNewNodesAndEdges} onResourceTypeUpdate={updateResourceType}
            previousLabelVars={previousLabelVars}
            setLabels={setLabels}
            parent={data.parent} parentResourceType={data.parentResourceType}
            />
            <Handle type="target" position={Position.Top} id="a" style={{opacity: 0}} />
            <Handle type="source" position={Position.Bottom} id="b" style={{opacity: 0}} />
          </div>
        </div>
      </div>
    </>
  );
});
