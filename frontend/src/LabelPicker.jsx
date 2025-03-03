import Box from '@mui/material/Box';
import Alert from '@mui/material/Alert';
import AddCircleIcon from '@mui/icons-material/AddCircle';
import RemoveCircleIcon from '@mui/icons-material/RemoveCircle';
import IconButton, {iconButtonClasses} from '@mui/material/IconButton';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import Stack from '@mui/material/Stack';
import TextField, { textFieldClasses} from '@mui/material/TextField';
import React, { useCallback, useState } from 'react';

import GenericPicker from './GenericPicker'
import SingleLabelValPicker from './SingleLabelValPicker'

const themeFunction = (theme) => ({
	padding: 0,
  	color: "#ffffff",
	'& input': {
		color: "#ffffff",
		padding: "0px 0px 0px 5px",
		fontFamily: "Monospace",
		fontSize: "10px",
		height: "20px"
	},
	[`&.${textFieldClasses.root}`]: {
		width: "100px",
		minWidth:"10px",
	  },
	  [`&.${iconButtonClasses.root}`]: {
		  width: "100px",
		  minWidth:"10px",
		},
})

export default function LabelPicker({ resourceType, labels, sourceData, disabled, addLabel, deleteLabel, updateLabelKey, updateLabelVal }) {
	const [isErr, setErr] = useState(false)
	const getLabelKeyOptions = useCallback(
		() => {
			// todo: fetch labels from backend for a resource type
			return [
				`${resourceType}-1`,
				`${resourceType}-2`
			];
		}, [resourceType]
	);

	return (
		<>
			{
				isErr &&
				<Alert variant="filled" severity="error">
					Papiez pedal matke jebal
				</Alert>
			}
			<Stack direction="row" sx={{fontStyle: "italic", fontSize: "10px"}}>
				<p>Where:</p>
			</Stack>
            <List dense={true} sx={{padding: "0px"}}>
			{
				labels.map(
					label => {
						return <ListItem sx={{padding: "0px"}}>		
							<ChevronRightIcon />
							<GenericPicker key={`${label.id}-picker`} value={label.key} valuePlaceholder="Filter" updateData={(newKey) => updateLabelKey(label.id, newKey)} getOptions={getLabelKeyOptions}/>
							<Box key={`${label.id}-eq`}
								sx={{
									padding: '0px 10px 0px 10px',
									fontWeight: 600, mt: "2px"
								}}>
							=
							</Box>
							<Stack 
								key={`${label.id}-val-outer`}
								sx={{
								height: 20,
								padding: '0px',
								fontWeight: 600,
								lineHeight: '15px',
								borderRadius: '2px',
								}}
								direction="row"
							>
								{
									!sourceData &&
									<TextField key={`${label.id}-val-inner`}
										variant="outlined"
										disabled={disabled || false}
										sx={themeFunction}
										value={label.val}
										fullWidth= {false}
										onChange={(event) => {setErr(false); if (event.target.value==="dupa") setErr(true); updateLabelVal(label.id, event.target.value)}}
									/>
								}
								{
									sourceData &&
									<SingleLabelValPicker key={`${label.id}-val-inner`}
										labelVal={label.val}
										data={sourceData}
										onFieldUpdate={(newValue) => {updateLabelVal(label.id, newValue)}}
										disabled={disabled || false}
									/>
								}
								<IconButton key={`${label.id}-val-del-outer`} sx={{padding: "0px", ml: "5px"}} aria-label="delete" disabled={disabled || false} onClick={() => deleteLabel(label.id)}>
									<RemoveCircleIcon key={`${label.id}-val-del-inner`} color="secondary" sx={{width: "16px", height: "16px"}}/>
								</IconButton>
							</Stack>
						</ListItem>
					}
				)
			}
				<ListItem sx={{padding: "0px"}}>
					<IconButton aria-label="add" sx={{padding: "0px", margin: "0px", width: "fit-content"}} size="small" disabled={disabled || false} onClick={addLabel}>
						<AddCircleIcon color="primary" sx={{width: "16px", height: "16px", padding: "0px"}}/>
					</IconButton>
				</ListItem>
           	</List>
		</>
	);
}
