import React, { memo, useState } from 'react';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import EditIcon from '@mui/icons-material/Edit';
import TextField, {textFieldClasses} from '@mui/material/TextField';

const boxTheme = (selected) => ({
	'&:hover': {
	  backgroundColor: selected ? null : 'rgba(55,55,55,.5)',
	  borderColor: '#0062cc',
	  boxShadow: 'none',
	},
    backgroundColor: selected ?  'rgba(55,55,155,.5)' : null,
    paddingTop: "4px",
    paddingLeft: "5px", 
    paddingRight: "5px",
    fontFamily: "Monospace",
})

const themeFunction = (theme) => ({
	padding: 0,
  	color: "#ffffff",
	'& input': {
		color: "#ffffff",
		padding: "0px 0px 0px 5px",
		fontFamily: "Monospace",
		fontSize: "10px",
		height: "20px",
	},
	[`&.${textFieldClasses.root}`]: {
		width: "100px",
		minWidth:"10px",
	},
})

export default memo(({ name, selected, onSelect, onEditEnd }) => {
    const [renaming, setRenaming] = useState(false);
    const [editedValue, setEditedValue] = useState(name)

  return (
    <Box sx={boxTheme(selected)}
    onClick={() => onSelect(name)}
    >
        <Stack alignItems="center" height="100%" sx={{padding: "2px"}} direction="row">
        {
            selected && 
            <EditIcon onClick={() => setRenaming(true)}/>
        }
        {
            selected && renaming &&
            <TextField
            sx={themeFunction}
            variant="outlined"
            value={editedValue}
            onChange={(e) => setEditedValue(e.target.value)}
            onKeyDown={(e) => {
                if (e.key === 'Enter') {
                    setRenaming(false);
                    onEditEnd(editedValue);
                }
            }}
            />
        }
        {
            !renaming &&
            <Box sx={{display: "inline-block", verticalAlign: "middle"}}>{name}</Box>
        }
        </Stack>
    </Box>
  );
});
