import React, { memo, useEffect, useState } from 'react';
import { Stack } from '@mui/material';
import { useQueryStore } from './state/QueryState';
import JQPicker from './pickers/JQPicker';
import {Box} from '@mui/material';
import * as jq from "jq-wasm"


export default memo(({ fieldId, value, data }) => {
  const updateField = useQueryStore((state) => state.updateField);
  const redisplay = useQueryStore((state) => state.redisplay);
  const setRedisplay = useQueryStore((state) => state.setRedisplay);
  const [text, setText] = useState("")

  useEffect(() => {
    async function update() {
      if(redisplay){
        const t = await jq.raw(data, value, ["-r", "-c"]);
        setText(t.stdout);
      }
    }
    update();
    setRedisplay(false);
  }, [redisplay, setRedisplay]);

  return (
    <Stack sx={{padding: "0px", margin: "0px"}}>
      <JQPicker 
        value={value}
        example={data}
        placeholder='Write JQ query here'
        updateData={(newVal) => {
          updateField(fieldId, {val: newVal});
        }}
        handleCloseExt={() => setRedisplay(true)}
        buttonSx={{
          fontSize: "8px",
          margin: '0',
          height: "10px",
          padding: '0',
          background: "#333333",
          borderRadius: "2px",
          fontWeight: 100,
          lineHeight: '10px',
          display: 'flex', 
          alignItems: 'center'
        }}
      />
      <Box sx={{margin: '0', padding: '0', height: "24px", display: 'flex', fontWeight: 997, fontSize: "22px", alignItems: 'center'}}>
        {text}
      </Box>
    </Stack>
  );
});
