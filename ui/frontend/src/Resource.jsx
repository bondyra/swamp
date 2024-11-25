import React, { memo, useState } from 'react';
import { JSONTree } from 'react-json-tree';
import { Handle, Position, useReactFlow } from "@xyflow/react";


 
export default memo(({ data, isConnectable }) => {
  const reactFlow = useReactFlow();
  return (
    <div className="resource-node">
        <img src={data.icon} alt={data.resource_type} />
        <label>{data.module}.{data.resource_type}</label>
        <em>{data.resource_id}</em>
        <JSONTree shouldExpandNodeInitially={() => false} data={data.obj} 
              labelRenderer={([key]) => <strong>{key}</strong>}
              valueRenderer={(raw) => <em>{raw}</em>}
            />
        <textarea placeholder='What to query?' onKeyDown={(evt) =>{
          if (evt.key === 'Enter') {
            const query = evt.target.value;
            console.log("CIPA")
            const [mod, resource_type] = query.split('.');
            const newNodes = [];
            const newEdges = [];
            fetch(`http://localhost:8000/${mod}/${resource_type}?pattern=${data.resource_id}`)
            .then(response => response.json())
            .then(response => {
              response.results.forEach(function(result) {
                const newId = `${mod}.${resource_type}.${result.id}`;
                newNodes.push({
                  id: newId,
                  position: { x: 0, y: 0 },
                  type: 'resource',
                  parentId: data.id,
                  data: {
                    id: newId,
                    module: mod,
                    resource_type: resource_type,
                    resource_id: result.id,
                    obj: result.obj,
                    icon: `./icons/aws/${resource_type}.svg` 
                  },
                });
                newEdges.push({id: `${data.id}-${newId}`, source: data.id, target: newId, style: {strokeWidth: 5} });
              })
            })
            .then(() => {
              reactFlow.addNodes(newNodes);
              reactFlow.addEdges(newEdges);
            });
          }
        }}/>
        <Handle type="target" position={Position.Top} id="a" />
        <Handle type="source" position={Position.Bottom} id="b" />
    </div>
  );
});
