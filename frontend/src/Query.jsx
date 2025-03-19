import React, { memo } from 'react';
import { useCallback, useEffect, useState } from 'react';
import { Handle, Position, useReactFlow } from '@xyflow/react';

import {JSONPath} from 'jsonpath-plus';

import QueryWizard from './QueryWizard';
import { getAllJSONPaths } from './Utils';
import { useBackend } from './BackendProvider';

export default memo(({ id, data, isConnectable }) => {
  const reactFlow = useReactFlow();
    const [childPaths, setChildPaths] = useState([]);
    const [parentPaths, setParentPaths] = useState([]);
    const backend = useBackend();

  const addNewNodesAndEdges = (items) => {
    var newNodes = [];
    var newEdges = [];
    items.forEach(item => {
      const newNodeId = `${item.resourceType}.${item.result._id}`;
      newNodes.push({
        id: newNodeId,
        position: { x: 0, y: 0 },
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
      })
    }, [id, reactFlow]);

    const updateChildPath = useCallback(
      (newValue) => {
        reactFlow.updateNodeData(id, (node) => {
          return { 
            ...node.data,
            childPath: newValue
          };
        })
      }, [id, reactFlow]
    );
  
    const updateParentPath = useCallback(
      (newValue) => {
        reactFlow.updateNodeData(id, (node) => {
          return { 
            ...node.data,
            parentPath: newValue
          };
        })
      }, [id, reactFlow]
    );

  const setLabels = useCallback(
    (labels) => {
      reactFlow.updateNodeData(id, (node) => {
        return {...node.data, labels: labels};
      })
    }, [id, reactFlow]);

  useEffect(() => {
    const loadAttributes = async () => {
      if (data.resourceType === null || data.resourceType === undefined)
        return []
      const attributes = await backend.attributes(data.resourceType)
      setChildPaths(attributes.map(a=> {return {value: a.path, description: a.description}}))
    };
    loadAttributes();
  }, [data.resourceType, setChildPaths, backend]);

  useEffect(() => {
    setParentPaths(
      data.parent ? getAllJSONPaths(data.parent).map(p => {
        return {
          value: p,
          description: JSONPath({path: p, json: data.parent})
        }
      }) : []
    )
  }, [data.parent])

  return (
    <>
      <div className="wrapper">
        <div className="inner">
          <div className="body">
            <QueryWizard 
            nodeId={id} resourceType={data.resourceType} labels={data.labels} doSomethingWithResults={addNewNodesAndEdges} onResourceTypeUpdate={updateResourceType}
            setLabels={setLabels}
            join={data.parent !== undefined}
            childPath={data.childPath} childPaths={childPaths} onChildPathUpdate={updateChildPath}
            parentPath={data.parentPath} parentPaths={parentPaths} onParentPathUpdate={updateParentPath} parentResourceType={data.parentResourceType}
            getParentVal={(p) => JSONPath({path: p, json: data.parent})}
            />
            <Handle type="target" position={Position.Top} id="a" style={{opacity: 0}} />
            <Handle type="source" position={Position.Bottom} id="b" style={{opacity: 0}} />
          </div>
        </div>
      </div>
    </>
  );
});
