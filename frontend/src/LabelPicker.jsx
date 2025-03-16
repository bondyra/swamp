import Box from '@mui/material/Box';
import AddCircleIcon from '@mui/icons-material/AddCircle';
import RemoveCircleIcon from '@mui/icons-material/RemoveCircle';
import IconButton, {iconButtonClasses} from '@mui/material/IconButton';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import Stack from '@mui/material/Stack';
import TextField, { textFieldClasses} from '@mui/material/TextField';
import React, { useEffect, useState } from 'react';

import { useBackend } from './BackendProvider';
import GenericPicker from './GenericPicker'
import SingleLabelValPicker from './SingleLabelValPicker'

let labelId = 1;
const newLabelId = () => labelId++;

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

export default function LabelPicker({ resourceType, labels, setLabels, disabled }) {
	const [attributes, setAttributes] = useState(new Map())
	const backend = useBackend();

	const labelsWithUpdatedLabel = (labels, updatedLabel) => {
		return labels.map(l => {
			if (l.id === updatedLabel.id)
				l = updatedLabel
			return l
		})
	}

	useEffect(() => {
		const loadAttributes = async () => {
			if (resourceType === null || resourceType === undefined)
				return []
			const attributes = await backend.attributes(resourceType)
			setAttributes(new Map(attributes.map(a=> [a.path, a])))
			setLabels(attributes.filter(a => a.query_required).map(a => {return {id: newLabelId(), key: a.path, val: "", required: true, allowedValues: a.allowed_values}}))
		};
        loadAttributes();
	}, [resourceType, setLabels, backend]);

	return (
		<>
			<Stack direction="row" sx={{fontStyle: "italic", fontSize: "10px"}}>
				<p>Where:</p>
			</Stack>
            <List dense={true} sx={{padding: "0px"}}>
			{
				labels.map(
					label => {
						return <ListItem key={`${label.id}-list-item`} sx={{padding: "0px"}}>		
							<ChevronRightIcon />
							<GenericPicker 
							key={`${label.id}-picker`} value={label.key} valuePlaceholder="Filter" 
							disabled = {label.required || disabled || false}
							updateData={(newKey) => setLabels(labelsWithUpdatedLabel(labels, {...label, key: newKey, allowedValues: (attributes.get(newKey) ?? Object()).allowed_values}))} 
							options={[...attributes.keys()]}/>
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
									(!label.allowedValues && !disabled) &&
									<TextField key={`${label.id}-val-inner`}
										variant="outlined"
										disabled={disabled || false}
										sx={themeFunction}
										value={label.val}
										fullWidth= {false}
										onChange={(event) => setLabels(labelsWithUpdatedLabel(labels, {...label, val: event.target.value}))}
									/>
								}
								{
									(!label.allowedValues && disabled) &&
									<Box sx={{ fontSize: 13, width:"fit-content", color: "#586069"}}>{label.val}</Box>
								}
								{
									label.allowedValues &&
									<SingleLabelValPicker key={`${label.id}-val-inner`}
										labelVal={label.val}
										options={label.allowedValues}
										onFieldUpdate={(newValue) =>  setLabels(labelsWithUpdatedLabel(labels, {...label, val: newValue}))}
										disabled={disabled || false}
									/>
								}
								{ (!label.required && !disabled) &&
								<IconButton key={`${label.id}-val-del-outer`} sx={{padding: "0px", ml: "5px"}} aria-label="delete"
								onClick={() => setLabels(labels.filter(x => x.id !== label.id))}>
									<RemoveCircleIcon key={`${label.id}-val-del-inner`} color="secondary" sx={{width: "16px", height: "16px"}}/>
								</IconButton>
								}
							</Stack>
						</ListItem>
					}
				)
			}
			{
				!disabled &&
				<ListItem sx={{padding: "0px"}}>
					<IconButton aria-label="add" sx={{padding: "0px", margin: "0px", width: "fit-content"}} size="small" disabled={disabled || false} 
					onClick={() => setLabels([...labels, {id: newLabelId(), key:  "", val: ""}])}>
						<AddCircleIcon color="primary" sx={{width: "16px", height: "16px", padding: "0px"}}/>
					</IconButton>
				</ListItem>
			}
           	</List>
		</>
	);
}
