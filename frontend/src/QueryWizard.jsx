import Box from '@mui/material/Box';
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress';
import DownloadForOfflineIcon from '@mui/icons-material/DownloadForOffline';
import Stack from '@mui/material/Stack';
import { styled } from '@mui/material/styles';
import React from 'react';
import { useCallback, useState } from 'react';

import LabelPicker from './LabelPicker';
import GenericPicker from './GenericPicker';
  
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

const resourceTypes = await fetch(`http://localhost:8000/resource-types`).then(response => response.json());


export default function QueryWizard({
  nodeId, resourceType, labels, doSomethingWithResults, onResourceTypeUpdate, addLabel, deleteLabel, updateLabel, overwriteLabels,
  join, childPath, childPaths, onChildPathUpdate, parentPath, parentPaths, onParentPathUpdate
}) {
  const [disabled, setDisabled] = useState(false)
  const [loading, setLoading] = useState(false);

  const getIconSrc = useCallback((r) => r ? `./icons/${r.replace(".", "/")}.svg` : undefined, [])

  const query = useCallback(async () => {
    const [provider, resource_type] = resourceType.split(".")
    const qs = (labels ?? []).map(l=> `${l.key}=${l.val}`).join("&")
    // TODO: inject pre-request validation (e.g. label vals must not be empty)
    // TOOD: display errors
    return await fetch(`http://localhost:8000/get?provider=${provider}&resource=${resource_type}&${qs}`)
      .then(response => response.json())
      .then(response => {
        return response.results.map(result => {
          return {
            provider: provider,
            resource_type: resource_type,
            data: result
          };
        })
      })
  }, [resourceType, labels]);

  const onClick = useCallback(async () => {
    setDisabled(true);
    setLoading(true);
    var results = await query()
    doSomethingWithResults(results)
    setLoading(false);
  }, [doSomethingWithResults, query])

  return (
    <>
      <Stack direction="row">
        <Stack sx={leftTheme}>
          <Stack direction="row">
            <Box sx={{fontSize: "14px", fontWeight:"600", mr: "10px", fontFamily: "monospace"}}>
              <p>{join ? "JOIN" : "SELECT"}</p>
            </Box>
            <GenericPicker disabled={disabled || false} value={resourceType} valuePlaceholder="What?" updateData={onResourceTypeUpdate} options={resourceTypes} getIconSrc={getIconSrc}/>
          </Stack>
          {
            join &&
            <Stack direction="row">
              <Box sx={{fontSize: "14px", fontWeight:"600", mr: "10px", fontFamily: "monospace"}}>
                <p>ON</p>
              </Box>
              <GenericPicker disabled={disabled || false} value={childPath} valuePlaceholder="Child field" updateData={onChildPathUpdate} options={childPaths}/>
              <Box sx={{fontSize: "14px", fontWeight:"600", ml: "10px", mr: "10px", fontFamily: "monospace"}}>
                <p>=</p>
              </Box>
              <GenericPicker disabled={disabled || false} value={parentPath} valuePlaceholder="Parent field" updateData={onParentPathUpdate} options={parentPaths}/>
            </Stack>
          }
          <LabelPicker 
          nodeId={nodeId} resourceType={resourceType} labels={labels} disabled={disabled || false}
          addLabel={addLabel} deleteLabel={deleteLabel} updateLabel={updateLabel} overwriteLabels={overwriteLabels}
          />
        </Stack>
        <RunButton variant="contained" aria-label="run" disabled={loading}  onClick={onClick} backgroundcolor="primary" sx={{ml: "5px"}}>
          {loading && <CircularProgress color="white" size="20px"/>}
          {!loading && <DownloadForOfflineIcon/>}
        </RunButton>
      </Stack>
    </>
  );
};
