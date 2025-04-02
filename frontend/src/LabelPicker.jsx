import AddCircleIcon from '@mui/icons-material/AddCircle';
import RemoveCircleIcon from '@mui/icons-material/RemoveCircle';
import IconButton, {iconButtonClasses} from '@mui/material/IconButton';
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight';
import LinkIcon from '@mui/icons-material/Link';
import EmergencyIcon from '@mui/icons-material/Emergency';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import Stack from '@mui/material/Stack';
import TextField, { textFieldClasses} from '@mui/material/TextField';
import React, { useEffect, useState } from 'react';

import {JSONPath} from 'jsonpath-plus';

import { useBackend } from './BackendProvider';
import SingleFieldPicker from './SingleFieldPicker'
import SingleLabelValPicker from './SingleLabelValPicker'
import { Tooltip } from '@mui/material';
import LabelOp from './LabelOp';

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

export default function LabelPicker({ resourceType, labels, setLabels, previousLabelVars, parent, parentResourceType }) {
	const [attributes, setAttributes] = useState(new Map())
	const [localLabels, setLocalLabels] = useState(labels);
	const backend = useBackend();

	useEffect(() => {setLabels(localLabels);}, [localLabels, setLabels]);

	const labelsWithUpdatedLabelsById = (ll, extraLabels) => {
		return ll.map(l => {
			const ul = extraLabels.filter(x => l.id === x.id)
			return ul ? {...l, ...ul[0]} : l
		})
	}

	const labelsWithNewLabels = (ll, extraLabels) => {
		const extraLabelKeys = extraLabels.map(l => l.key)

		return [
			...ll.filter(l => !extraLabelKeys.includes(l.key)),
			...extraLabels
		]
	}

	useEffect(() => {  // when resourceType changes
		const loadAttributes = async () => {
			if (resourceType === null || resourceType === undefined)
				return [];
			// re-load attributes
			const attributes = await backend.attributes(resourceType);
			const requiredAttributes = attributes.filter(a => a.query_required);
			setAttributes(new Map(attributes.map(a=> [a.path, a])));
			// refresh labels with potentially new required attributes
			const newRequiredAttributes = requiredAttributes.map(a => {
				const val = [...(previousLabelVars ?? new Map())[a.path] ?? [""]][0]
				return {
					id: newLabelId(), key: a.path, val: val, keyImmutable: true, 
					required: true, allowedValues: a.allowed_values, dependsOn: a.depends_on
				}
			});
			// if there's a link label, try to get default link key val
			var linkKeyVal = null;
			if (parentResourceType) {
				const resp = await backend.linkSuggestion(resourceType, parentResourceType)
				linkKeyVal = {key: resp.key, val: resp.val, op: resp.op}
			}

			const labelsWithLinkSuggestion = (ll, linkKeyVal) => {
				if(!linkKeyVal || (linkKeyVal.key === "" && linkKeyVal.val === ""))
					return ll
				const actualValues = JSONPath({path: linkKeyVal.val, json: parent})
				if (!actualValues || actualValues.length !== 1)
					return ll
				
				return ll.map(l => {
					if (l.id === "link")
						return {...l, ...{key: linkKeyVal.key, op: linkKeyVal.op, val: actualValues[0]}};
					return l;
				});
			}
			// set the labels (local only to prevent re-render loop)
			setLocalLabels(oldLabels => {
				return labelsWithLinkSuggestion(
					labelsWithNewLabels(oldLabels, newRequiredAttributes), 
					linkKeyVal
				)
			});
		};
        loadAttributes();
	}, [resourceType, parentResourceType, setLocalLabels, backend, previousLabelVars, parent]);

	return (
		<>
            <List dense={true} sx={{padding: "0px", margin: "0px"}}>
			{
				localLabels
				.sort((a, b) => (b.id === "link") - (a.id === "link"))
				.sort((a, b) => b.required - a.required)  // meaning: put required labels first - descending order based on "required" flag (true=1, false=0)
				.map(
					label => {
						return <ListItem key={`${label.id}-list-item`} sx={{padding: "0px", margin: "0px", height:"16px"}}>		
							{
								label.id === "link" && <Tooltip title="Specify how to link this?"><LinkIcon sx={{width: "8px", height: "8px", padding: "0px", color: "yellow"}}/></Tooltip>
							}
							{
								label.id !== "link" && label.required && <Tooltip title="Required by provider"><EmergencyIcon sx={{width: "8px", height: "8px", padding: "0px", color: "darkred"}}/></Tooltip>
							}
							{
								label.id !== "link" && !label.required && <KeyboardArrowRightIcon sx={{width: "8px", height: "8px", padding: "0px"}}/>
							}
							<SingleFieldPicker 
							freeSolo={true}
							key={`${label.id}-picker`} value={label.key} valuePlaceholder="Filter" 
							disabled = {label.keyImmutable || false}
							updateData={(newKey) => {
								const newLabels = labelsWithUpdatedLabelsById(
									localLabels,
									[{
										...label, 
										key: newKey, 
										allowedValues: (attributes.get(newKey) ?? Object()).allowed_values, 
										required: (attributes.get(newKey) ?? Object()).query_required ?? false,
										dependsOn: (attributes.get(newKey) ?? Object()).depends_on ?? null,
									}]
								);
								setLocalLabels(newLabels);
							}} 
							options={[...attributes.values().filter(v=> !v.query_required).map(v => {return {value: v.path, description: v.description}})]}/>
							<LabelOp op={label.op ?? "eq"} 
										change={(val) => {
											const newLabels = labelsWithUpdatedLabelsById(
												localLabels, [{...label, op: val}]
											)
											setLocalLabels(newLabels);
										}}/>
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
									(!label.allowedValues && !label.dependsOn) &&
									<TextField key={`${label.id}-val-inner`}
										variant="outlined"
										sx={themeFunction}
										value={label.val}
										fullWidth= {false}
										onChange={(event) => {
											const newLabels = labelsWithUpdatedLabelsById(
												localLabels, [{...label, val: event.target.value}]
											)
											setLocalLabels(newLabels);
										}}
									/>
								}
								{
									(label.allowedValues || label.dependsOn) &&
									<SingleLabelValPicker key={`${label.id}-val-inner`}
										labelVal={label.val}
										options={label.allowedValues || []}
										onFieldUpdate={async (newValue) =>  {
											// update allowed values for each label for which it key dependsOn this key
											const val = typeof newValue === "object" ? newValue.value : newValue;  // TODO: label val should object
											const dependentLabels = await Promise.all(
												localLabels
												.filter(l => l.dependsOn === label.key)
												.map(async l => {
													return {...l, allowedValues: await backend.attributeValues(resourceType, l.key, [{key: label.key, val: val}])}
												})
											)
											const newLabels = labelsWithUpdatedLabelsById(
												localLabels, [...dependentLabels, {...label, val: val}]
											);
											setLocalLabels(newLabels);
										}}
										descr={`Select one of allowed values for ${label.key}`}
									/>
								}
								{ !label.required &&
								<IconButton key={`${label.id}-val-del-outer`} sx={{padding: "0px", ml: "5px"}} aria-label="delete"
								onClick={() => {
									const newLabels = localLabels.filter(x => x.id !== label.id);
									setLocalLabels(newLabels);
								}}>
									<RemoveCircleIcon key={`${label.id}-val-del-inner`} color="secondary" sx={{width: "16px", height: "16px", padding: "0px"}}/>
								</IconButton>
								}
							</Stack>
						</ListItem>
					}
				)
			}
			<ListItem sx={{padding: "0px"}}>
				<IconButton aria-label="add" sx={{padding: "0px", margin: "0px", width: "fit-content"}} size="small"
				onClick={() => {
					const newLabels = [...localLabels, {id: newLabelId(), key:  "", val: ""}];
					setLocalLabels(newLabels);
				}}>
					<AddCircleIcon color="primary" sx={{width: "16px", height: "16px", padding: "0px"}}/>
				</IconButton>
			</ListItem>
           	</List>
		</>
	);
}
