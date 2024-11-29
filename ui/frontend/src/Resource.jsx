import React, { memo } from 'react';
import { JSONTree } from 'react-json-tree';
import { Handle, Position } from "@xyflow/react";
 
export default memo(({ data, isConnectable }) => {
  return (
    <div className="resource-node">
        <img src={data.icon} alt={data.resource_type} />
        <label>{data.module}.{data.resource_type}</label>
        <em>{data.resource_id}</em>
        <JSONTree shouldExpandNodeInitially={() => false} data={data.obj} 
              labelRenderer={([key]) => <strong>{key}</strong>}
              valueRenderer={(raw) => <em>{raw}</em>}
            />
        <Handle type="target" position={Position.Top} id="a" />
        <Handle type="source" position={Position.Bottom} id="b" />
    </div>
  );
});
