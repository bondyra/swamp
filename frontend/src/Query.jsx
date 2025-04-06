import React, { memo, useCallback } from 'react';
import { Handle, Position, useReactFlow } from '@xyflow/react';

import QueryWizard from './QueryWizard';

export default memo(({ id, data, isConnectable }) => {
  const reactFlow = useReactFlow();

  function* getItemIdsToDelete (nodeIds, parentId) {
    // 1 traverse graph to identify which nodes/edges to delete
    // - when nodeId has only one parent (parentId), we need to delete it + getItemIdsToDelete for its child nodes
    // - when nodeId has multiple parents, we delete the edge parentId->nodeId only
    const nodeIdAndItsParentEdges = nodeIds.map(ni => {return {nodeId: ni, parentEdges: reactFlow.getEdges().filter(e => e.target === ni)}});
    const nodeIdsToDelete = nodeIdAndItsParentEdges.filter(x => x.parentEdges.length <= 1).map(x => x.nodeId);
    const edgeIdsToDelete = nodeIdAndItsParentEdges.filter(x => x.parentEdges.length > 1).map(x => x.parentEdges.filter(e => e.source === parentId).map(e => e.id)[0]);
    for (const nodeId of nodeIds) {
      const childrenNodeIds = reactFlow.getEdges().filter(e => e.source === nodeId).map(e => e.target);
      yield* getItemIdsToDelete(childrenNodeIds, nodeId);
    }
    yield* nodeIdsToDelete.map(i => {return {type: "node", id: i}});
    yield* edgeIdsToDelete.map(i => {return {type: "edge", id: i}});
  }

  const addNewNodesAndEdges = (items) => {
    const thisQueryNode = reactFlow.getNode(id);
    const allChildrenNodesOfThisQueryNode = reactFlow.getNodes().filter(n => n.data.queryId === id);
    const allChildrenNodeIds = allChildrenNodesOfThisQueryNode.map(n => n.id);
    const newItems = items.map(item => { return {...item, nodeId: `${item.resourceType}.${item.result._id}`}});
    const newIds = newItems.map(ni => ni.nodeId)
    const childrenNodeIdsToDelete = allChildrenNodeIds.filter(i => !newIds.includes(i))
    const itemIdsToDelete = [...getItemIdsToDelete(childrenNodeIdsToDelete, id)];
    const nodeIdsToDelete = itemIdsToDelete.filter(x => x.type === "node").map(x => x.id);
    const edgeIdsToDelete = itemIdsToDelete.filter(x => x.type === "edge").map(x => x.id);
    var newNodes = [];
    var newEdges = [];
    newItems.forEach(item => {
      if (allChildrenNodeIds.includes(item.nodeId))  // if node already exists, don't overwrite it
        return
      newNodes.push({
        id: item.nodeId,
        position: { 
          x: thisQueryNode.position.x ?? 0, 
          y: thisQueryNode.position.y + 100 ?? 0 
        },
        type: 'resource',
        data: {
          id: item.nodeId,
          resourceType: item.resourceType,
          inline: {},
          result: item.result,
          queryId: id
        },
      });
      newEdges.push({id: `${id}-${item.nodeId}`, source: id, target: item.nodeId, style: {strokeWidth: 5} });
    });
    reactFlow.setNodes([
      ...reactFlow.getNodes().filter(n => !nodeIdsToDelete.includes(n.id)),
      ...newNodes
    ])
    reactFlow.setEdges([
      ...reactFlow.getEdges().filter(e => !edgeIdsToDelete.includes(e.id)),
      ...newEdges
    ])
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

  return (
    <>
      <div className="wrapper">
        <div className="inner">
          <div className="body">
            <QueryWizard 
              nodeId={id} resourceType={data.resourceType} labels={data.labels} 
              doSomethingWithResults={addNewNodesAndEdges} onResourceTypeUpdate={updateResourceType}
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
