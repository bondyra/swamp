import React, { memo } from 'react';
import { useCallback } from 'react';
import { Handle, Position, useReactFlow } from '@xyflow/react';

import QueryWizard from './QueryWizard';

let labelId = 1;
const newLabelId = () => labelId++;

export default memo(({ id, data, isConnectable }) => {
  const reactFlow = useReactFlow();

  const addNewNodesAndEdges = (results) => {
    var newNodes = [];
    var newEdges = [];
    results.forEach(result => {
      const newNodeId = `${result.provider}.${result.resource_type}.${result.data.__id}`;
      newNodes.push({
        id: newNodeId,
        position: { x: 0, y: 0 },
        type: 'resource',
        data: {
          id: newNodeId,
          inline: {},
          ...result
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
      })
    }, [id, reactFlow]);

  const _updateLabel = useCallback(
  ({labelId, newKey, newVal}) => {
    reactFlow.updateNodeData(id, (node) => {
      var labels = node.data.labels ?? []
      return { 
        ...node.data,
        labels: labels.map(l => {
        if (l.id === labelId){
          l.key = newKey ?? l.key
          l.val = newVal ?? l.val
        }
        return l;
        })
      };
    })
  }, [id, reactFlow]);

  const updateLabelKey = (labelId, newKey) => _updateLabel({labelId: labelId, newKey: newKey})
  const updateLabelVal = (labelId, newVal) => _updateLabel({labelId: labelId, newVal: newVal})
  
  const addLabel = useCallback(
    () => {
      reactFlow.updateNodeData(id, (node) => {
        var labels = node.data.labels ?? []
        return { ...node.data, labels: labels.concat({id: newLabelId(), key: "", val: ""}) } ;
      });
    }, [id, reactFlow]
  )
  const deleteLabel = useCallback(
    (labelId) => {
      reactFlow.updateNodeData(id, (node) => {
        var labels = node.data.labels ?? []
        return { ...node.data, labels: labels.filter(x => x.id !== labelId) } ;
      });
    }, [id, reactFlow]
  )
  const overwriteLabels = useCallback(
    (newLabels) => {
      reactFlow.updateNodeData(id, (node) => {
        return { ...node.data, labels: newLabels.map(l => {return {id: newLabelId(), ...l}})} ;
      });
    }, [id, reactFlow]
  )

  return (
    <>
      <div className="wrapper">
        <div className="inner">
          <div className="body">
            <QueryWizard 
            nodeId={id} resourceType={data.resourceType} labels={data.labels} doSomethingWithResults={addNewNodesAndEdges} onResourceTypeUpdate={updateResourceType}
            addLabel={addLabel} deleteLabel={deleteLabel} updateLabelKey={updateLabelKey} updateLabelVal={updateLabelVal} overwriteLabels={overwriteLabels}
            />
            <Handle type="target" position={Position.Top} id="a" style={{opacity: 0}} />
            <Handle type="source" position={Position.Bottom} id="b" style={{opacity: 0}} />
          </div>
        </div>
      </div>
    </>
  );
});
