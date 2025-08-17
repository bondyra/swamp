import React, { memo } from 'react';
import { Handle, Position, useReactFlow } from '@xyflow/react';
import Stack from '@mui/material/Stack';
import { Box } from '@mui/material';
import { getIconSrc } from './Utils';

export default memo(({ id, data, isConnectable }) => {
  const reactFlow = useReactFlow();

  return (
    <>
      <Handle type="target" position={Position.Top} id="a" style={{opacity: 0, height: "1px", width:"1px"}} isConnectable={false}/>
      <Stack direction="column" sx={{pt: "2px", height: 18, width: 24, alignItems: "center", background: data.selected ? '#aaaaff' : '#3e3e3e'}}>
        <Box
          component="img"
          sx={{height: 8, width: 8, flexShrink: 0, padding: "0", margin: "0", bgcolor: "white", borderRadius: '2px'}}
          src={getIconSrc(data.resourceType)} alt=""
        />
        <svg viewBox='0 0 24 10' width= "24" height="10" padding="0px" preserveAspectRatio="none">
          <text x="0" y="5" fontSize="8" fontFamily="monospace" textLength="24" lengthAdjust="spacingAndGlyphs" fill="white">
            {data.resourceType}
          </text>
        </svg>
      </Stack>
      <Handle type="source" style={{opacity: 0, height: "1px", width:"1px"}} position={Position.Bottom} id="b" isConnectable={false}/>
    </>
  );
});
