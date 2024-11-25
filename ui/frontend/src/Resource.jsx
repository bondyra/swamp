import React, { memo } from 'react';
import { JSONTree } from 'react-json-tree';
 
export default memo(({ data, isConnectable }) => {
  return (
    <div className="resource-node">
        <img src={data.icon} alt={data.resource_type} />
        <label>{data.module}.{data.resource_type}</label>
        <em>{data.id}</em>
        <JSONTree shouldExpandNodeInitially={() => false} data={data.content} 
              labelRenderer={([key]) => <strong>{key}</strong>}
              valueRenderer={(raw) => <em>{raw}</em>}
            />
    </div>
  );
});
