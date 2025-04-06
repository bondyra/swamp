import Box from '@mui/material/Box';
import Button from '@mui/material/Button'
import IconButton from '@mui/material/IconButton';
import OpenInFullIcon from '@mui/icons-material/OpenInFull';
import CloseFullscreenIcon from '@mui/icons-material/CloseFullscreen';
import CircularProgress from '@mui/material/CircularProgress';
import DownloadForOfflineIcon from '@mui/icons-material/DownloadForOffline';
import InfoIcon from '@mui/icons-material/Info';
import ErrorIcon from '@mui/icons-material/Error';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import WarningIcon from '@mui/icons-material/Warning';
import Stack from '@mui/material/Stack';
import { styled } from '@mui/material/styles';
import { Tooltip } from '@mui/material';
import React, { useEffect } from 'react';
import { useCallback, useState } from 'react';
import { useReactFlow } from '@xyflow/react';

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

const validateLabels = (labels) => {
  const invalidAllowedValueLabels = labels.filter(l => l.allowedValues && !l.allowedValues.includes(l.val))
  if(invalidAllowedValueLabels && invalidAllowedValueLabels.length > 0){
    return `Following labels have not allowed values: ${invalidAllowedValueLabels.map(l => l.key).join(",")}`
  }
  return null
}

export default function QueryWizard({
  nodeId, resourceType, labels, doSomethingWithResults, onResourceTypeUpdate, setLabels, parent, parentResourceType
}) {
  const reactFlow = useReactFlow();
  const [loading, setLoading] = useState(false);
  const [resourceTypes, setResourceTypes] = useState([])
  const [status, setStatus] = useState("initial")
  const [message, setMessage] = useState(null)
  const [fullSize, setFullSize] = useState(true);
  const [previousLabelVars, setPreviousLabelsVars] = useState(null);
    
  useEffect(() => {
      var result = new Map();
      reactFlow.getNodes().forEach(n => {
        (n.data.labels ?? []).forEach(l => {
          if(! result[l.key])
            result[l.key] = new Set()
          result[l.key].add(l.val)
        })
      })
      setPreviousLabelsVars(result)
  }, [setPreviousLabelsVars, reactFlow])

  const backend = useBackend();

  const setError = useCallback((msg) => {
    setStatus("failed");
    setMessage(msg);
    setLoading(false);
  }, [setStatus, setMessage, setLoading])

  const onClick = useCallback(async () => {
    setLoading(true);
    if (!resourceType) {
      setError("You must select a resource type!")
      return
    }
    var queryLabels = labels
    var err = null
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
    }
  }, [doSomethingWithResults, resourceType, labels, backend, setError])

  useEffect(() => {
    async function update() {
      const newResourceTypes = await backend.resourceTypes();
      setResourceTypes(newResourceTypes);
    }
    update();
  }, [backend, setResourceTypes]);

  return (
    <>
    {
      fullSize &&
      <Stack direction="row">
        <Stack sx={leftTheme}>
          <Stack direction="row">
            <Box sx={{fontSize: "14px", fontWeight:"600", mr: "10px", fontFamily: "monospace"}}>
              <p>GET</p>
            </Box>
            <SingleFieldPicker value={resourceType} updateData={onResourceTypeUpdate} options={resourceTypes} getIconSrc={getIconSrc}
            valuePlaceholder="What?" popperPrompt={parent !== null ? "Select resource to join" : "Select resource to query"}/>
          </Stack>
          <LabelPicker nodeId={nodeId} resourceType={resourceType} labels={labels} 
          setLabels={setLabels} previousLabelVars={previousLabelVars} parent={parent} parentResourceType={parentResourceType}/>
        </Stack>
        <Stack>
          <Stack direction="row" sx={{justifyContent: "flex-end", height: "30%"}}>
            <Tooltip title={message || "Query status will be here"}>
              <Box sx={{ml: "5px", mb: "10px"}} justifyContent="center" display="flex" >
                {status === "initial" && <InfoIcon sx={{color: "lightblue"}} fontSize="small"/>}
                {status === "failed" && <ErrorIcon sx={{color: "darkred"}} fontSize="small"/>}
                {status === "warning" && <WarningIcon sx={{color: "yellow"}} fontSize="small"/>}
                {status === "success" && <CheckCircleIcon sx={{color: "green"}} fontSize="small"/>}
              </Box>
            </Tooltip>
            <Tooltip title="Minimize">
              <IconButton onClick={() => setFullSize(false)} sx={{ml: "5px", mb: "10px"}} display="flex">
                    <CloseFullscreenIcon sx={{color: "gray"}} fontSize="small"/>
              </IconButton>
            </Tooltip>
          </Stack>
          <Tooltip title="Run or refresh">
            <RunButton variant="contained" aria-label="run"  onClick={onClick} backgroundcolor="primary" sx={{ml: "5px", height:"80%"}}>
              {loading && <CircularProgress color="white" size="20px"/>}
              {!loading && <DownloadForOfflineIcon/>}
            </RunButton>
          </Tooltip>
        </Stack>
      </Stack>
    }
    {
      !fullSize &&
      <Stack>
        <Stack direction="row" sx={{justifyContent: "center"}}>
          <Tooltip title={message || "Query status will be here"}>
            <Box justifyContent="center" display="flex" >
              {status === "initial" && <InfoIcon sx={{color: "lightblue"}} fontSize="small"/>}
              {status === "failed" && <ErrorIcon sx={{color: "darkred"}} fontSize="small"/>}
              {status === "warning" && <WarningIcon sx={{color: "yellow"}} fontSize="small"/>}
              {status === "success" && <CheckCircleIcon sx={{color: "green"}} fontSize="small"/>}
            </Box>
          </Tooltip>
          <Tooltip title="Maximize">
            <IconButton onClick={() => setFullSize(true)} sx={{padding: "0px"}}>
              <OpenInFullIcon sx={{color: "gray"}} fontSize="small"/>
            </IconButton>
          </Tooltip>
        </Stack>
        <Box sx={{fontSize: "14px", fontWeight:"600", fontFamily: "monospace", justifyContent: "center", display: "flex"}}>
          <p>{resourceType ?? "<who knows what>"}</p>
        </Box>
        <Tooltip title="Run or refresh">
          <RunButton variant="contained" aria-label="run"  onClick={onClick} backgroundcolor="primary">
            {loading && <CircularProgress color="white" size="20px"/>}
            {!loading && <DownloadForOfflineIcon/>}
          </RunButton>
        </Tooltip>
      </Stack>
    }
    </> 
  );
};
