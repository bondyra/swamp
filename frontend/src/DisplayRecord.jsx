import { memo, useEffect, useState } from 'react';
import { Stack } from '@mui/material';
import { useQueryStore } from './state/QueryState';
import JQPicker from './pickers/JQPicker';
import {Box} from '@mui/material';
import RemoveCircleIcon from '@mui/icons-material/RemoveCircle';
import {Typography} from '@mui/material';
import IconButton, {iconButtonClasses} from '@mui/material/IconButton';
import * as jq from "jq-wasm"


export default memo(({ fieldId, value, data }) => {
  const updateField = useQueryStore((state) => state.updateField);
  const removeField = useQueryStore((state) => state.removeField);
  const redisplay = useQueryStore((state) => state.redisplay);
  const setRedisplay = useQueryStore((state) => state.setRedisplay);
  const [text, setText] = useState("");

  

  useEffect(() => {
    async function update() {
      if((redisplay && value) || (!text && value)){  // if needs to be redisplayed and there's some value OR it's a first load and there is some value already (on mount)
        const t = await jq.raw(data, value, ["-r", "-c"]);
        setText(t.stdout);
      }
    }
    update();
    setRedisplay(false);
  }, [redisplay, setRedisplay]);

  return (
    <Stack sx={{padding: "0px", margin: "0px", width: "100%"}}>
      <Stack direction="row" sx={{background: "#333333"}}>
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
          borderRadius: "2px",
          fontWeight: 100,
          width: "100%",
          lineHeight: '10px',
          alignItems: 'center'
        }}
      />
      <IconButton key={`${fieldId}-val-del-outer`} sx={{padding: "0px", ml: "0px", color: "gray"}} aria-label="delete" onClick={() => {removeField(fieldId);}}>
        <RemoveCircleIcon key={`${fieldId}-val-del-inner`} sx={{width: "16px", height: "10px", padding: "0px"}}/>
      </IconButton>
      </Stack>
      <Box sx={{margin: '0', padding: '0', display: 'flex', alignItems: 'center'}}>
        <Typography sx={{ fontWeight: 600, fontSize: "20px", fontFamily: "monospace", whiteSpace: 'pre-line' }}>
          {text}
        </Typography>
      </Box>
    </Stack>
  );
});
