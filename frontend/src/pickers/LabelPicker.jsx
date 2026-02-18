import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import DoubleArrowIcon from '@mui/icons-material/DoubleArrow';
import RemoveCircleIcon from '@mui/icons-material/RemoveCircle';
import IconButton, {iconButtonClasses} from '@mui/material/IconButton';
import EmergencyIcon from '@mui/icons-material/Emergency';
import Stack from '@mui/material/Stack';
import TextField, { textFieldClasses} from '@mui/material/TextField';
import Box from '@mui/material/Box';

import JQPicker from './JQPicker'
import SingleLabelValPicker from './SingleLabelValPicker'
import { Tooltip } from '@mui/material';
import LabelOp from '../LabelOp';
import { randomString } from '../Utils';
import { useBackend } from '../BackendProvider';
import { useQueryStore } from '../state/QueryState';
import { useState, useEffect } from 'react';


const themeFunction = (theme) => ({
  	color: "#ffffff",
	'& input': {
		color: "#ffffff",
		padding: "0px 0px 0px 5px",
		fontFamily: "Monospace",
		fontSize: "14px",
  		fontWeight: 600,
		height: "20px"
	},
	[`&.${textFieldClasses.root}`]: {
		width: "auto",
		minWidth:"10px",
	},
	[`&.${iconButtonClasses.root}`]: {
		width: "100px",
		minWidth:"10px",
	},
})

export default function LabelPicker({ resourceType, labels, setLabels, attributes }) {
	const backend = useBackend();
	const savedLabels = useQueryStore((state) => state.savedLabels);
	const mergedLabels = (ll, extraLabels) => {
		return ll.map(l => {
			const ul = extraLabels.filter(x => l.id === x.id)
			return ul ? {...l, ...ul[0]} : l
		})
	}
	const [example, setExample] = useState({})

    useEffect(() => {
        async function update() {
			if (resourceType){
				const ex = await backend.example(resourceType);
				setExample(ex);
			}
        }
        update();
	}, [backend, resourceType, setExample])

	return (
		<>
            <Stack direction="row" sx={{padding: "2px", flexWrap: 'wrap'}}>
			{
				labels
				.sort((a, b) => b.required - a.required)  // meaning: put required labels first - descending order based on "required" flag (true=1, false=0)
				.map(
					label => {
						return <Box key={`${label.id}-list-item`} sx={{padding: "0", border: "1px solid gray", borderRadius: "10px"}}>
						<Stack direction="row" sx={{alignItems: "center"}}>
							{
								!label.required && <Tooltip title="User specified JQ query"><DoubleArrowIcon sx={{width: "16px", height: "16px", color: "yellow"}}/></Tooltip>
							}
							{
								label.required && <Tooltip title="Required by provider"><EmergencyIcon sx={{width: "16px", height: "16px", color: "darkred"}}/></Tooltip>
							}
							<JQPicker 
								key={`${label.id}-picker`} 
								value={label.key}
								example={example}
								disabled = {label.required || false}
								updateData={(newKey) => {
									var newContent = { key: newKey };
									if (!label.op && !label.val){
										const matchingSavedLabel = savedLabels.filter(s => s.key === newKey);
										newContent.op = matchingSavedLabel.length > 0 ? matchingSavedLabel[0].op: null
										newContent.val = matchingSavedLabel.length > 0 ? matchingSavedLabel[0].val: null
									}
									const newLabels = mergedLabels(
										labels, [{...label, ...newContent}]
									);
									setLabels(newLabels);
								}}
								// not needed V ?
								options={[...attributes.values().filter(v=> !v.query_required).map(v => {return {value: v.path, description: v.description}})]}
							/>
							<LabelOp op={label.op ?? "=="} 
										change={(val) => {
											const newLabels = mergedLabels(
												labels, [{...label, op: val}]
											)
											setLabels(newLabels);
										}}/>
							<Stack
								key={`${label.id}-val-outer`}
								sx={{
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
										fullWidth= {true}
										onChange={(event) => {
											setLabels(mergedLabels(
												labels, [{...label, val: event.target.value}]
											));
										}}
									/>
								}
								{
									(label.allowedValues || label.dependsOn) &&
									<SingleLabelValPicker key={`${label.id}-val-inner`}
										labelVal={label.val}
										options={label.allowedValues || []}
										onFieldUpdate={async (newValue) =>  {
											// update allowed values for each label for which its key dependsOn this key
											const val = typeof newValue === "object" ? newValue.value : newValue;
											const dependentLabels = await Promise.all(
												labels
												.filter(l => l.dependsOn === label.key)
												.map(async l => {
													var av = await backend.attributeValues(resourceType, l.key, [{key: label.key, val: val}]);
													return {...l, allowedValues: av}
												})
											)
											const newLabels = mergedLabels(
												labels, [...dependentLabels, {...label, val: val}]
											);
											setLabels(newLabels);
										}}
										descr={`Select one of allowed values for ${label.key}`}
									/>
								}
								{ !label.required &&
								<IconButton key={`${label.id}-val-del-outer`} sx={{padding: "0px", ml: "5px"}} aria-label="delete"
								onClick={() => {setLabels(labels.filter(x => x.id !== label.id));}}>
									<RemoveCircleIcon key={`${label.id}-val-del-inner`} color="secondary" sx={{width: "16px", height: "16px", padding: "0px"}}/>
								</IconButton>
								}
							</Stack>
						</Stack>
						</Box>
					}
				)
			}
			<Box sx={{padding: "0", pl: "2px", pb: "4px"}}>
				<IconButton aria-label="add" sx={{width: "fit-content", padding: "0"}} size="small"
				onClick={() => {
					setLabels([...labels, {id: randomString(8), key: "", val: "", op: "==", required: false}]);
				}}>
					<AddCircleOutlineIcon sx={{width: "16px", height: "16px", padding: "0px", color: "gray"}}/>
				</IconButton>
			</Box>
           	</Stack>
		</>
	);
}
