import * as React from 'react';
import PropTypes from 'prop-types';
import { styled } from '@mui/material/styles';
import Popper from '@mui/material/Popper';
import ClickAwayListener from '@mui/material/ClickAwayListener';
import Autocomplete, { autocompleteClasses } from '@mui/material/Autocomplete';
import ButtonBase from '@mui/material/ButtonBase';
import InputBase from '@mui/material/InputBase';
import Box from '@mui/material/Box';
import CheckIcon from '@mui/icons-material/Check';
import Tooltip from '@mui/material/Tooltip';

import {JSONPath} from 'jsonpath-plus';

const StyledAutocompletePopper = styled('div')(({ theme }) => ({
  [`& .${autocompleteClasses.paper}`]: {
    boxShadow: 'none',
    margin: 0,
    color: 'inherit',
    fontSize: 13,
  },
  [`& .${autocompleteClasses.listbox}`]: {
    padding: 0,
    backgroundColor: '#fff',
    ...theme.applyStyles('dark', {
      backgroundColor: '#1c2128',
    }),
    [`& .${autocompleteClasses.option}`]: {
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
      [`&.${autocompleteClasses.focused}, &.${autocompleteClasses.focused}[aria-selected="true"]`]:
        {
          backgroundColor: theme.palette.action.hover,
        },
    },
  },
  [`&.${autocompleteClasses.popperDisablePortal}`]: {
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
  color: '#24292e',
  backgroundColor: '#fff',
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

const StyledInput = styled(InputBase)(({ theme }) => ({
  padding: 10,
  width: 'fit-content',
  borderBottom: '1px solid #eaecef',
  ...theme.applyStyles('dark', {
    borderBottom: '1px solid #30363d',
  }),
  '& input': {
    borderRadius: 4,
    padding: 8,
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    fontSize: 14,
    backgroundColor: '#fff',
    border: '1px solid #30363d',
    ...theme.applyStyles('dark', {
      backgroundColor: '#0d1117',
      border: '1px solid #eaecef',
    }),
    '&:focus': {
      boxShadow: '0px 0px 0px 3px rgba(3, 102, 214, 0.3)',
      borderColor: '#0366d6',
      ...theme.applyStyles('dark', {
        boxShadow: '0px 0px 0px 3px rgb(12, 45, 107)',
        borderColor: '#388bfd',
      }),
    },
  },
}));


function getAllOptions(obj, prefix = '') {
  let options = [];
  
  if (typeof obj === 'object' && obj !== null) {
      for (let key in obj) {
        if (Array.isArray(obj[key])) {
          for(var i = 0; i < obj[key].length; i++) {
            let newPrefix = prefix === '' ? `${key}[${i}]` : `${prefix}.${key}[${i}]`;
            options = options.concat(getAllOptions(obj[key][i], newPrefix));
          }
        }
        else if (typeof obj[key] === 'object'){
          let newPrefix = prefix === '' ? key : `${prefix}.${key}`;
          options = options.concat(getAllOptions(obj[key], newPrefix));
        } else {
          let newPrefix = prefix === '' ? key : `${prefix}.${key}`;
          options.push({path: newPrefix, data: JSONPath({path: newPrefix, json: obj})});
        }
      }
  }
  
  return options;
}

const Button = styled(ButtonBase)(({ theme }) => ({
  fontSize: 13,
  width: 'fit-content',
  textAlign: 'left',
  fontWeight: 600,
  color: '#586069',
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
    width: 'fit-content',
  },
  '& svg': {
    width: 16,
    height: 16,
  },
}));

export default function SingleLabelValPicker({labelVal, data, onFieldUpdate, descr, disabled}) {
  const allOptions = getAllOptions(data)
  const [anchorEl, setAnchorEl] = React.useState(null);
  const [options, setOptions] = React.useState([])

  const handleClick = (event) => {
    setOptions(allOptions);
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    if (anchorEl) {
      anchorEl.focus();
    }
    setAnchorEl(null);
  };

  const open = Boolean(anchorEl);
  const id = open ? 'github-label' : undefined;

  return (
    <React.Fragment>
    <Tooltip title={labelVal ? JSONPath({path: labelVal, json: data}) : ""}>
      <Box sx={{ width: '100%', fontSize: 13, paddingBottom: "8px" }}>
        <Button disableRipple aria-describedby={id} onClick={handleClick} disabled={disabled || false}>
          <span>{labelVal ? `this.${labelVal}` : "Choose field"}</span>
        </Button>
      </Box>
      </Tooltip>
      <StyledPopper id={id} open={open} anchorEl={anchorEl} placement="bottom-start">
        <ClickAwayListener onClickAway={handleClose}>
          <div>
            <Box
              sx={(t) => ({
                borderBottom: '1px solid #30363d',
                padding: '8px 10px',
                fontWeight: 600,
                ...t.applyStyles('light', {
                  borderBottom: '1px solid #eaecef',
                }),
              })}
            >
              {descr ?? "Choose one of the values"}
            </Box>
            <Autocomplete
              open
              onClose={(event, reason) => handleClose()}
              value={labelVal}
              onChange={(event, newValue, reason) => {onFieldUpdate(newValue.path)}}
              renderTags={() => null}
              noOptionsText="No labels"
              renderOption={(props, option, { selected }) => {
                const { key, ...optionProps } = props;
                return (
                  <li key={key} {...optionProps}>
                    <Box
                      sx={(t) => ({
                        flexGrow: 1,
                        '& span': {
                          color: '#8b949e',
                          ...t.applyStyles('light', {
                            color: '#586069',
                          }),
                        },
                      })}
                    >
                      {option.path}
                      <br />
                      <span>{option.data ?? "(null)"}</span>
                    </Box>
                    <Box
                      component={CheckIcon}
                      sx={{ opacity: 0.6, width: 18, height: 18 }}
                      style={{
                        visibility: selected ? 'visible' : 'hidden',
                      }}
                    />
                  </li>
                );
              }}
              options={[...options]}
              getOptionLabel={(o) => o.path ?? o}
              renderInput={(params) => (
                <StyledInput
                  ref={params.InputProps.ref}
                  inputProps={params.inputProps}
                  autoFocus
                  placeholder="Select path"
                />
              )}
              slots={{
                popper: PopperComponent,
              }}
            />
          </div>
        </ClickAwayListener>
      </StyledPopper>
    </React.Fragment>
  );
}
