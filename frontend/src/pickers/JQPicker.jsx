import {Fragment, useEffect, useState} from 'react';
import PropTypes from 'prop-types';
import { styled } from '@mui/material/styles';
import Popper from '@mui/material/Popper';
import ClickAwayListener from '@mui/material/ClickAwayListener';
import ButtonBase from '@mui/material/ButtonBase';
import { TextField, textFieldClasses } from '@mui/material';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import * as jq from "jq-wasm";


const StyledAutocompletePopper = styled('div')(({ theme }) => ({
  [`& .${textFieldClasses.paper}`]: {
    boxShadow: 'none',
    margin: 0,
    color: 'inherit',
  },
  [`& .${textFieldClasses.listbox}`]: {
    padding: 0,
    backgroundColor: '#fff',
    ...theme.applyStyles('dark', {
      backgroundColor: '#1c2128',
    }),
    [`& .${textFieldClasses.option}`]: {
      minHeight: 'auto',
      alignItems: 'flex-start',
      padding: 8,
      borderBottom: '1px solid #eaecef',
      ...theme.applyStyles('dark', {
        borderBottom: '1px solid #30363d',
      }),
      '&[aria-selected="true"]': {
        backgroundColor: 'transparent',
      },
      [`&.${textFieldClasses.focused}, &.${textFieldClasses.focused}[aria-selected="true"]`]:
        {
          backgroundColor: theme.palette.action.hover,
        },
    },
  },
  [`&.${textFieldClasses.popperDisablePortal}`]: {
    position: 'relative',
  },
}));


function PopperComponent(props) {
  const { disablePortal, anchorEl, open, ...other } = props;
  return <StyledAutocompletePopper {...other} />;
}


PopperComponent.propTypes = {
  anchorEl: PropTypes.any,
  disablePortal: PropTypes.bool,
  open: PropTypes.bool.isRequired,
};


const StyledPopper = styled(Popper)(({ theme }) => ({
  border: '1px solid #e1e4e8',
  boxShadow: `0 8px 24px ${'rgba(149, 157, 165, 0.2)'}`,
  color: '#ffffff',
  backgroundColor: '#000000',
  borderRadius: 6,
  width: 300,
  zIndex: theme.zIndex.modal,
  fontSize: 13,
  ...theme.applyStyles('dark', {
    border: '1px solid #30363d',
    boxShadow: '0 8px 24px rgb(1, 4, 9)',
    color: '#c9d1d9',
    backgroundColor: '#1c2128',
  }),
}));


const Button = styled(ButtonBase)(({ theme }) => ({
  fontWeight: 600,
  color: '#ffffff',
  ...theme.applyStyles('dark', {
    color: '#8b949e',
  }),
  '&:hover,&:focus': {
    color: '#0366d6',
    ...theme.applyStyles('dark', {
      color: '#58a6ff',
    }),
  },
  '& span': {
    width: '100%',
  },
  '& svg': {
    width: 16,
    height: 16,
  },
}));


export default function JQPicker({value, updateData, example, disabled, handleCloseExt, buttonSx}) {
  const [anchorEl, setAnchorEl] = useState(null);
  const [result, setResult] = useState("");
  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = async () => {
    if (anchorEl) {
      anchorEl.focus();
    }
    setAnchorEl(null);
    updateData(value);
    if (handleCloseExt)
      await handleCloseExt();
  };

  useEffect(() => {
    const runJq = async () => {
      if (value === null || value === undefined || !value)
        return;
      const newVal = await jq.raw(example, value, ["-r", "-c"]);
      setResult(newVal.stdout);
    }
    runJq();
  }, [value])

  const open = Boolean(anchorEl);
  const id = open ? 'github-label' : undefined;

  return (
    <Fragment>
      <Box
      sx={
        buttonSx ?? {
          fontSize: "14px",
          mt: '0px',
          height: 20,
          padding: '.15em 2px',
          fontWeight: 600,
          lineHeight: '15px',
          borderRadius: '2px'
        }
      }>
        <Button 
        disabled={disabled || false} disableRipple aria-describedby={id} onClick={handleClick} 
        sx={{padding: "0", margin:"0"}}>
          <Box
            key={value}
            // style={{
            //   backgroundColor: "#fff",
            //   color: theme.palette.getContrastText("#fff"),
            // }}
          >
            {value || "jq query"}
          </Box>
        </Button>
      </Box>
      <StyledPopper id={id} open={open} anchorEl={anchorEl} placement="bottom-start">
        <ClickAwayListener onClickAway={handleClose}>
          <div>
            <Stack>
              <TextField sx={{
                width: "100%", background: "black",
                '& input': {
                  color: "#ffffff",
                  padding: "5px",
                  fontFamily: "Monospace",
                  fontSize: "12px",
                  height: "20px"
                },
              }} value={value} placeholder='Write JQ query here' onChange={async e => {
                  updateData(e.target.value);
              }}/>
            </Stack>
            <Stack>
              <Box sx={{fontFamily: "monospace", fontSize: "8px", fontWeight: 100, margin: "4px"}}>Example resource:</Box>
              <Box 
                component="section" 
                sx={{width: "100%", height: "100px", fontSize: "10px", fontWeight: 600, borderBottom: '1px solid #eaecef', padding: '8px 10px', overflow: 'scroll'}}
              >
                <pre><code>{JSON.stringify(example, null, 2)}</code></pre>
              </Box>
            </Stack>
            <Stack>
              <Box sx={{fontFamily: "monospace", fontSize: "8px", fontWeight: 100, margin: "4px"}}>What your JQ returns:</Box>
              <Box 
                component="section" 
                sx={{width: "100%", height: "auto", fontSize: "10px", borderBottom: '1px solid #eaecef', padding: '8px 10px', overflow: 'scroll'}}
              >
                <pre><code>{result}</code></pre>
              </Box>
            </Stack>
          </div>
        </ClickAwayListener>
      </StyledPopper>
    </Fragment>
  );
}
