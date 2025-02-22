import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Stack from '@mui/material/Stack';
import React, { memo } from 'react';
import { useCallback, useState } from 'react';
import { Handle, Position, useReactFlow } from '@xyflow/react';

import LabelPicker from './LabelPicker';
import SwampPicker from './SwampPicker'

export default memo(({ id, data, isConnectable }) => {
  const [disabled, setDisabled] = useState(false)
  const [loading, setLoading] = useState(false);
  const reactFlow = useReactFlow();

  const query = useCallback((event) => {
    setLoading(true)
    var newNodes = [];
    var newEdges = [];

    const [provider, resource_type] = data.resourceType.split(".")

    const promise = fetch(`http://localhost:8000/${provider}/${resource_type}`)
      .then(response => response.json())
      .then(response => {
        response.results.forEach(function(result) {
          const newNodeId = `${provider}.${resource_type}.${result.id}`;
          newNodes.push({
            id: newNodeId,
            position: { x: 0, y: 0 },
            type: 'resource',
            data: {
              id: newNodeId,
              provider: provider,
              resource_type: resource_type,
              resource_id: result.id,
              obj: result.obj,
              _dataFilter: "",
              _displayedData: result.obj,
            },
          });
          newEdges.push({id: `${id}-${newNodeId}`, source: id, target: newNodeId, style: {strokeWidth: 5} });
        })
      });
    Promise.all([promise]).then(() => {
      reactFlow.addNodes(newNodes);
      reactFlow.addEdges(newEdges);
    });
    setLoading(false)
  }, [reactFlow, data, id]);

  const updateResourceType = useCallback(
    (newResourceType) => {
      reactFlow.updateNodeData(id, (node) => {
        return { 
          ...node.data,
          resourceType: newResourceType
        };
      })
    }, [id, reactFlow]);

  const getResourceTypeOptions = useCallback(
    () => {
      return [
      {
        name: 'aws.vpc',
        description: 'Twoja matka',
      },
      {
        name: 'aws.subnet',
        description: 'To dziwka',
      }
    ];
  }, []);

  // todo focus shift
  return (
    <>
      <div className="wrapper">
        <div className="inner">
          <div className="body">
            <Stack>
              <SwampPicker value={data.resourceType} updateData={updateResourceType} getOptions={getResourceTypeOptions}/>
              <LabelPicker nodeId={id} resource={data.resourceType} labels={data.labels} sourceData={{}} disabled={disabled} />
              <Button color="primary" size="small" variant="contained" disabled={loading || disabled} onClick={() => {setDisabled(true); query()}}>Run</Button>
              {loading && <CircularProgress/>}
            </Stack>
            <Handle type="target" position={Position.Top} id="a" style={{opacity: 0}} />
            <Handle type="source" position={Position.Bottom} id="b" style={{opacity: 0}} />
          </div>
        </div>
      </div>
    </>
  );
});
