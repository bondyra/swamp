import React, { memo, useCallback, useState } from 'react';
import { styled } from '@mui/material/styles';
import { Handle, Position, useReactFlow } from '@xyflow/react';
import ButtonBase from '@mui/material/ButtonBase';
import Collapse from '@mui/material/Collapse';
import Stack from '@mui/material/Stack';
import PlaylistPlayIcon from '@mui/icons-material/PlaylistPlay';
import DatasetLinkedIcon from '@mui/icons-material/DatasetLinked';
import SettingsEthernetIcon from '@mui/icons-material/SettingsEthernet';
import { Tooltip } from '@mui/material';

import DataDisplay from './DataDisplay';
import QueryWizard from './QueryWizard';
import MultipleFieldPicker from './MultipleFieldPicker';
import { Box } from '@mui/material';
import { randomString, getIconSrc, getAllUniqueValues } from './Utils';


const Button = styled(ButtonBase)(({ theme }) => ({
  fontSize: 13,
  width: '100%',
  textAlign: 'left',
  paddingBottom: 8,
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
    width: '100%',
  },
  '& svg': {
    width: 16,
    height: 16,
  },
}));

export default memo(({ id, data, isConnectable }) => {
  const reactFlow = useReactFlow();
  const [isLinked, setIsLinked] = useState(false);
  const [inlineCollapsed, setInlineCollapsed] = useState(true);

  const addQuery = useCallback((event) => {
    const targetNodeId = `${id}-link-${randomString(8)}`
    // we need to remove the wrapper bounds, in order to get the correct position
    const { clientX, clientY } = 'changedTouches' in event ? event.changedTouches[0] : event;
    // todo: move it to querywizard and add link label if it's a join query
    const newNode = {
      id: targetNodeId,
      position: reactFlow.screenToFlowPosition({
        x: clientX,
        y: clientY,
      }),
      type: 'query',
      data: {
        resourceType: null,
        parentResourceType: data.resourceType, 
        parent: data.result,
        labels: [
          {
            id: "link", key: "", val: "", required: true, 
            allowedValues: getAllUniqueValues(data.result),
            dependsOn: null
          }
        ]
      },
      origin: [0.5, 0.0],
    };

    reactFlow.setNodes((nds) => nds.concat(newNode));
    reactFlow.setEdges((eds) =>
      eds.concat({ id: `${id}-${targetNodeId}`, source: id, target: targetNodeId, style: {strokeWidth: 5} }),
    );
    setIsLinked(true);  // make it useEffect 
  }, [id, reactFlow, data.result, data.resourceType]);

  const inlineResources = useCallback(
    (results) => {
      reactFlow.updateNodeData(id, (node) => {
        return {
          ...node.data,
          inline: {...node.data.inline, results: results.map(r => r.result)}
        };
      })
    }, [id, reactFlow]
  );

  const updateInlineResourceType = useCallback(
    (newValue) => {
      reactFlow.updateNodeData(id, (node) => {
        return { 
          ...node.data,
          inline: {...node.data.inline, resourceType: newValue}
        };
      })
    }, [id, reactFlow]);

  const updateSelectedFields = useCallback(
    (newValue) => {
      reactFlow.updateNodeData(id, (node) => {
        return { 
          ...node.data,
          selectedFields: newValue
        };
      })
    }, [id, reactFlow]
  )

  const updateInlineSelectedFields = useCallback(
    (newValue) => {
      reactFlow.updateNodeData(id, (node) => {
        return { 
          ...node.data,
          inline: {...node.data.inline, selectedFields: newValue}
        };
      })
    }, [id, reactFlow]
  );

  const setInlineLabels = useCallback(
    (labels) => {
      reactFlow.updateNodeData(id, (node) => {
        return { 
          ...node.data,
          inline: {...node.data.inline, labels: labels}
        };
      })
    }, [id, reactFlow]);

  const propagateFieldsToSiblings = useCallback(
    () => {
      reactFlow.setNodes(
        reactFlow.getNodes().map(n => {
          if (n.data.queryId !== data.queryId || n.id === id)
            return n;
          return {
            ...n,
            data: {
              ...n.data,
              selectedFields: data.selectedFields
            }
          }
        })
      )
    }, [reactFlow, data.queryId, data.selectedFields, id]
  )

  return (
    <>
      <div className="wrapper">
        <div className="inner">
          <div className="body">
          <Handle type="target" position={Position.Top} id="a" style={{opacity: 0}} />
            <Stack direction="row">
              <Box
                component="img"
                sx={{
                  height: 24,
                  flexShrink: 0,
                  borderRadius: '3px',
                  padding: "1px",
                  mr: "5px"
                }}
                src={getIconSrc(data.resourceType)} alt=""
              />
              <Stack direction="column">
                <div className="header-label-row"><label className="swamp-resource">{data.resourceType}</label></div>
                <div className="header-label-row"><label className="swamp-resource-id">{data.result._id}</label></div>
              </Stack>
            </Stack>
            <hr/>
            <Stack direction="row">
              <MultipleFieldPicker data={data.result} selectedFields={data.selectedFields || []} updateSelectedFields={updateSelectedFields}
              header={"Select fields"} descr={`Select fields of ${data.resourceType} to display`}/>
              <Tooltip title="Propagate selections to all siblings">
                <Button disableRipple aria-describedby={id} sx={{width: "auto", pb: "0px"}} onClick={propagateFieldsToSiblings}> 
                  <SettingsEthernetIcon/>
                </Button>
              </Tooltip>
            </Stack>
            <Box> 
              <DataDisplay nodeId={id} data={data.result} selectedFields={data.selectedFields || []}/>
            </Box>
            <hr/>
            <Box>
            <Tooltip title="Join some other resources and display it here for convenience">
              <Button disableRipple aria-describedby={id} onClick={(_) => setInlineCollapsed(!inlineCollapsed)} sx={{fontSize: "10px", pb: "0px"}}>
                  <PlaylistPlayIcon />
                  <span>Extra resources</span>
                </Button>
                </Tooltip>
              </Box>
            <Box>
              <Collapse in={!inlineCollapsed} timeout="auto" unmountOnExit>
              <Box sx={{borderLeft: "1px solid gray", ml: "12px", pl:"12px", mt:"0px", pt: "0px"}}>
                <QueryWizard 
                nodeId={id} resourceType={data.inline.resourceType} labels={data.inline.labels || [{id: "link", key: "", val: "", required: true, allowedValues: getAllUniqueValues(data.result), dependsOn: null}]}
                doSomethingWithResults={inlineResources} onResourceTypeUpdate={updateInlineResourceType}
                setLabels={setInlineLabels}
                parent={data.result} parentResourceType={data.resourceType}
                />
                <MultipleFieldPicker data={data.inline.results} selectedFields={data.inline.selectedFields || []} updateSelectedFields={updateInlineSelectedFields} header="Results"/>
                <DataDisplay multiple nodeId={id} data={data.inline.results || []} selectedFields={data.inline.selectedFields || []}/>
              </Box>
              </Collapse>
            </Box>
            {
              !isLinked && 
              <Box>
                <Tooltip title="Join some other resources as new nodes">
                  <Button disableRipple aria-describedby={id} onClick={addQuery} sx={{pb: "0px"}}>
                    <DatasetLinkedIcon/>
                  </Button>
                </Tooltip>
              </Box>
            }
            <Handle type="source" style={{opacity: isLinked ? 1 : 0, borderRadius: "10%", height:"8px", width:"8px", bottom: "4px"}} position={Position.Bottom} id="b"/>
          </div>
        </div>
      </div>
    </>
  );
});
