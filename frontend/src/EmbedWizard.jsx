import React, { memo } from 'react';
import { useState } from 'react';

import Autocomplete from '@mui/material/Autocomplete';
import TextField from '@mui/material/TextField';


const whatProps = {
  options: [
    { name: 'aws.vpc' },
    { name: 'aws.subnet' }
  ],
  getOptionLabel: (option) => option.name,
}; 
const howProps = {
  options: [
    { name: 'dupa' },
    { name: 'cipa' }
  ],
  getOptionLabel: (option) => option.name,
};

const themeFunction = (theme) => ({
  width: 75,
  color: "#ffffff",
  '& input': {
    color: "#ffffff",
    fontFamily: "Monospace"
  },
  '& label': {
    color: "#aaaaaa",
    fontFamily: "Monospace"
  },
})


export default memo(({ props }) => {
  const [what, setWhat] = useState();
  const [how, setHow] = useState("");
  return (
    <>
    <Autocomplete
      {...whatProps}
      sx={themeFunction}
      value={what || null}
      onChange={(event, newWhat) => {
        setWhat(newWhat);
      }}
      autoComplete
      includeInputInList
      renderInput={(params) => (
        <TextField {...params} label="What?" variant="standard" />
      )}
    />
    <Autocomplete
    {...howProps}
    sx={themeFunction}
    value={how || null}
    onChange={(event, newHow) => {
      setHow(newHow);
    }}
    autoComplete
    includeInputInList
    renderInput={(params) => (
      <TextField {...params} label="How?" variant="standard" />
    )}
    />
    </>
  );
});
