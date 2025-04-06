

import React, { memo } from 'react';
import Button from '@mui/material/Button';
import Stack from '@mui/material/Stack';
import { Tooltip } from '@mui/material';


export default memo(({ op, change }) => {
  return (
    <Stack direction="row" sx={{padding: '0px 7px 0px 7px', fontWeight: 600, pb: "4px"}}>
        <Tooltip title="equals">
            <Button sx={{color: op === "eq" ? "white": "gray", height: "16px", width: "16px", minWidth: "16px"}} onClick={() => change("eq")}>
                =
            </Button>
        </Tooltip>
        <Tooltip title="contains">
            <Button sx={{color: op === "contains" ? "white": "gray", height: "16px", width: "16px", minWidth: "16px", transform: "rotate(180deg)"}} onClick={() => change("contains")}>
                ∈
            </Button>
        </Tooltip>
    </Stack>
  );
});
