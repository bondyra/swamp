import React, { memo } from 'react';
import TextField, { textFieldClasses} from '@mui/material/TextField';
import { useQueryStore } from './QueryState';
import Stack from '@mui/material/Stack';
import { IconButton } from '@mui/material';
import HighlightOffIcon from '@mui/icons-material/HighlightOff';


const themeFunction = (theme) => ({
  color: "#ffffff",
'& input': {
  color: "#ffffff",
  padding: "0px 0px 0px 5px",
  fontFamily: "Monospace",
  fontSize: "10px",
  height: "20px"
},
[`&.${textFieldClasses.root}`]: {
  width: "50%",
  minWidth:"10px",
  border: "1px solid white"
},
})

export default memo(({ parameter }) => {
  const updateParameter = useQueryStore((state) => state.updateParameter);
  const removeParameter = useQueryStore((state) => state.removeParameter);
  return (
    <>
      <Stack direction="row" sx={{alignItems: "center"}}>
        <TextField key={`${parameter.id}-key`}
                            variant="outlined"
                            sx={themeFunction}
                            value={parameter.key}
                            fullWidth= {false}
                            onChange={(event) => {updateParameter(parameter.id, {key: event.target.value});}}
                          />
        <TextField key={`${parameter.id}-val`}
                            variant="outlined"
                            sx={themeFunction}
                            value={parameter.val}
                            fullWidth= {false}
                            onChange={(event) => {updateParameter(parameter.id, {val: event.target.value});}}
                          />
      </Stack>
      <IconButton sx={{color: "white", height: "12px"}} onClick={() => removeParameter(parameter.id)}>
        <HighlightOffIcon/>
      </IconButton>
    </>
  );
});