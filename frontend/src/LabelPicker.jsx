import AddIcon from '@mui/icons-material/Add';
import Autocomplete from '@mui/material/Autocomplete';
import CloseIcon from '@mui/icons-material/Close';
import IconButton from '@mui/material/IconButton';
import CircularProgress from '@mui/material/CircularProgress';
import Stack from '@mui/material/Stack';
import TextField from '@mui/material/TextField';
import React, { useCallback } from 'react';
import { useReactFlow } from '@xyflow/react';


let labelId = 1;
const newLabelId = () => labelId++;

const themeFunction = (theme) => ({
  width: 200,
	padding: 0,
  color: "#ffffff",
  '& input': {
    color: "#ffffff",
    padding: "4px 4px",
    fontFamily: "Monospace"
  },
	'.MuiAutocomplete-root': {
    padding: 0
	},
	'.MuiAutocomplete-input': {
    padding: 0
	},
	'.MuiOutlinedInput-root': {
    padding: 0
	}
})

export default function LabelPicker({ nodeId, resource, labels, sourceData, disabled }) {
	const { updateNodeData } = useReactFlow();

	const updateLabel = useCallback(
		(label, fun) => {
			updateNodeData(nodeId, (node) => {
				return { 
					...node.data,
					labels: labels.map(l => {
						if (l.id === label.id)
							l = fun(l)
						return l
					})
				};
			})
		}, [labels, nodeId, updateNodeData]);

  const handleOpen = (label) => {
		updateLabel(label, l => {l.open = true; return l;});
    (async () => {
			updateLabel(label, l => {l.optionLoading = true; return l;});
      await new Promise((resolve) => {setTimeout(() => {resolve();}, 1e3);});  // test
			updateLabel(label, l => {l.optionLoading = false; return l;});
      const wwhat = resource ? resource.name : "notselected"
			var options = [
				{ name: `${wwhat}-1` },
				{ name: `${wwhat}-2` }
			]
			updateLabel(label, l => {l.options = options; return l;});
    })();
  };

  const handleClose = (label) => {
		updateLabel(label, l => {l.open = false; return l;});
		updateLabel(label, l => {l.options = []; return l;});
  };

  return (
    <>
		{
			labels.map(
				label => {
					return <Stack direction="row">
						{/* key */}
						<Autocomplete
						sx={themeFunction}
						disabled={disabled || false}
						open={label.open || false}
						onOpen={() => handleOpen(label)}
						onClose={() => handleClose(label)}
						options={label.options || []}
						loading={label.optionLoading || false}
						getOptionLabel={(option) => option.name}
						value={label.key || null}
						onChange={(_, newValue) => {updateLabel(label, (l) => {l.key = newValue; return l})}}
						autoComplete
						includeInputInList
						renderInput={(params) => (
							<TextField
								{...params}
								label=""
								slotProps={{
									input: {
										...params.InputProps,
										endAdornment: (
											<React.Fragment>
												{(label.optionLoading || false) ? <CircularProgress color="secondary" size={20} /> : null}
												{params.InputProps.endAdornment}
											</React.Fragment>
										),
									},
								}}
							/>
						)}
						/>
						=
						{/* value */}
						<TextField
							variant="outlined"
							disabled={disabled || false}
							sx={themeFunction}
							value={label.val}
							onChange={(event) => {updateLabel(label, (l) => {l.val = event.target.value; return l})}}
						/>
						<IconButton aria-label="delete" disabled={disabled || false} onClick={() => {
							updateNodeData(nodeId, (node) => {
								return { ...node.data, labels: labels.filter(x => x.id !== label.id) } ;
							});
						}}>
							<CloseIcon color="secondary"/>
						</IconButton>
					</Stack>
				}
			)
		}
			<IconButton aria-label="delete" disabled={disabled || false} onClick={(_) => {
				updateNodeData(nodeId, (node) => {
					return { ...node.data, labels: labels.concat({id: newLabelId(), key: "", val: ""}) } ;
				});
			}}>
        <AddIcon color="primary"/>
      </IconButton>
    </>
  );
}
