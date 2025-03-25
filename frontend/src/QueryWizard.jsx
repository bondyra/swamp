import Box from '@mui/material/Box';
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress';
import DownloadForOfflineIcon from '@mui/icons-material/DownloadForOffline';
import ErrorIcon from '@mui/icons-material/Error';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import WarningIcon from '@mui/icons-material/Warning';
import Stack from '@mui/material/Stack';
import { styled } from '@mui/material/styles';
import { Tooltip } from '@mui/material';
import React, { useEffect } from 'react';
import { useCallback, useState } from 'react';

import LabelPicker from './LabelPicker';
import SingleFieldPicker from './SingleFieldPicker';
import { useBackend } from './BackendProvider';
import { getIconSrc } from './Utils';

const RunButton = styled(Button)({
  width: "auto",
  height: "auto",
	boxShadow: 'none',
	textTransform: 'none',
  backgroundColor: "gray",
	fontSize: 10,
	padding: '1px',
	lineHeight: 1,
	fontFamily: "monospace",
	'&:hover': {
	  backgroundColor: '#0069d9',
	  borderColor: '#0062cc',
	  boxShadow: 'none',
	},
	'&:active': {
	  boxShadow: 'none',
	  backgroundColor: '#0062cc',
	  borderColor: '#005cbf',
	},
	'&:focus': {
	  boxShadow: '0 0 0 0.2rem rgba(0,123,255,.5)',
	},
});


const leftTheme = (theme) => ({
	'&:hover': {
	  backgroundColor: 'rgba(15,15,15,.5)',
	  borderColor: '#0062cc',
	  boxShadow: 'none',
	},
})

const validateJoinLabel = (childPath, parentPath) => {
  if (!childPath) return "Left side of join cannot be empty!";
  if (!parentPath) return "Right side of join cannot be empty!";
  return null;
}

const validateLabels = (labels) => {
  const invalidAllowedValueLabels = labels.filter(l => l.allowedValues && !l.allowedValues.includes(l.val))
  if(invalidAllowedValueLabels && invalidAllowedValueLabels.length > 0){
    return `Following labels have not allowed values: ${invalidAllowedValueLabels.map(l => l.key).join(",")}`
  }
  return null
}

export default function QueryWizard({
  nodeId, resourceType, labels, doSomethingWithResults, onResourceTypeUpdate, setLabels, previousLabelVars,
  join,
  childPath, onChildPathUpdate, childPaths, parentPath, onParentPathUpdate, parentPaths, getParentVal, parentResourceType
}) {
  const [disabled, setDisabled] = useState(false);
  const [loading, setLoading] = useState(false);
  const [resourceTypes, setResourceTypes] = useState([])
  const [status, setStatus] = useState("initial")
  const [message, setMessage] = useState(null)

  const backend = useBackend();

  const setError = useCallback((msg) => {
    setStatus("failed");
    setMessage(msg);
    setLoading(false);
    setDisabled(false);
  }, [setStatus, setMessage, setLoading, setDisabled])

  const onClick = useCallback(async () => {
    setDisabled(true);
    setLoading(true);
    if (!resourceType) {
      setError("You must select a resource type!")
      return
    }
    var queryLabels = labels
    var err = null
    if (join) {
      err = validateJoinLabel(childPath, parentPath)
      if(err) {
        setError(err);
        return
      }
      const parentVal = getParentVal(parentPath)
      queryLabels = [...queryLabels, {key: childPath, val: parentVal}]
    }
    err = validateLabels(queryLabels)
    if(err) {
      setError(err);
      return
    }
    try {
      var results = await backend.query(resourceType, queryLabels)
      doSomethingWithResults(results)
    } catch (e) {
      setError(e.message)
      return
    }
    setLoading(false);
    if (results && results.length > 0){
      setStatus("success");
      setMessage(`Returned ${results.length} results`);
    } else {
      setStatus("warning");
      setMessage(`No results!`);
      setDisabled(false);
    }
  }, [join, doSomethingWithResults, resourceType, labels, childPath, parentPath, backend, getParentVal, setError])

  useEffect(() => {
    async function update() {
      const newResourceTypes = await backend.resourceTypes();
      setResourceTypes(newResourceTypes);
    }
    update();
  }, [backend, setResourceTypes]);

  return (
    <>
      <Stack direction="row">
        <Stack sx={leftTheme}>
          <Stack direction="row">
            <Box sx={{fontSize: "14px", fontWeight:"600", mr: "10px", fontFamily: "monospace"}}>
              <p>GET</p>
            </Box>
            <SingleFieldPicker disabled={disabled || false} value={resourceType} updateData={onResourceTypeUpdate} options={resourceTypes} getIconSrc={getIconSrc}
            valuePlaceholder="What?" popperPrompt={join ? "Select resource to join" : "Select resource to query"}/>
          </Stack>
          {
            join &&
            <Stack direction="row">
              <Box sx={{fontSize: "12px", fontWeight:"600", mr: "5px", fontFamily: "monospace"}}>
                <p>WHEN</p>
              </Box>
              <SingleFieldPicker disabled={disabled || false} value={childPath} valuePlaceholder="Child field" updateData={onChildPathUpdate} options={childPaths}
              popperPrompt={resourceType ? `Choose attribute of ${resourceType} to join on` : "Please select resource type to join first"}
              />
              <Box sx={{fontSize: "10px", fontWeight:"600", fontFamily: "monospace", padding: '0px 7px 0px 7px'}}>
                <p>=</p>
              </Box>
              <SingleFieldPicker disabled={disabled || false} value={parentPath} valuePlaceholder="Parent field" updateData={onParentPathUpdate} options={parentPaths}
              popperPrompt={`Choose attribute of parent (${parentResourceType}) to join on`}
              />
            </Stack>
          }
          <LabelPicker nodeId={nodeId} resourceType={resourceType} labels={labels} 
          hasRun={disabled || false} setLabels={setLabels} previousLabelVars={previousLabelVars} parentResourceType={parentResourceType}/>
        </Stack>
        <Stack>
          <Tooltip title={message || "Query status will be here"}>
            <Box sx={{ml: "5px", mb: "10px", height: "20%"}} justifyContent="center" display="flex" >
              {status === "failed" && <ErrorIcon sx={{color: "darkred"}} fontSize="small"/>}
              {status === "warning" && <WarningIcon sx={{color: "yellow"}} fontSize="small"/>}
              {status === "success" && <CheckCircleIcon sx={{color: "green"}} fontSize="small"/>}
            </Box>
          </Tooltip>
          <Tooltip title="Run or refresh">
            <RunButton variant="contained" aria-label="run"  onClick={onClick} backgroundcolor="primary" sx={{ml: "5px", height:"80%"}}>
              {loading && <CircularProgress color="white" size="20px"/>}
              {!loading && <DownloadForOfflineIcon/>}
            </RunButton>
          </Tooltip>
        </Stack>
      </Stack>
    </> 
  );
};
