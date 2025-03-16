import React, { memo, useCallback, useEffect, useState } from 'react';
import { styled } from '@mui/material/styles';
import { Handle, Position, useReactFlow } from '@xyflow/react';
import ButtonBase from '@mui/material/ButtonBase';
import Collapse from '@mui/material/Collapse';
import DatasetIcon from '@mui/icons-material/Dataset';
import DatasetLinkedIcon from '@mui/icons-material/DatasetLinked';
import SettingsEthernetIcon from '@mui/icons-material/SettingsEthernet';

import {JSONPath} from 'jsonpath-plus';

import DataDisplay from './DataDisplay';
import QueryWizard from './QueryWizard';
import FieldPicker from './FieldPicker';
import { Box } from '@mui/material';
import { useBackend } from './BackendProvider';
import { getAllJSONPaths, getIconSrc } from './Utils';


let linkId = 1;
const newLinkId = () => `link-${linkId++}`;


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
  const [inlineChildPaths, setInlineChildPaths] = useState([]);
  const [inlineParentPaths, setInlineParentPaths] = useState([]);
  const [inlineCollapsed, setInlineCollapsed] = useState(true);
  const backend = useBackend();

  const addQuery = useCallback((event) => {
    const sourceNodeId = id
    const sourceNode = reactFlow.getNodes().filter(n=> n.id === sourceNodeId)[0]
    const targetNodeId = `${id}-${newLinkId()}`
    // we need to remove the wrapper bounds, in order to get the correct position
    const { clientX, clientY } = 'changedTouches' in event ? event.changedTouches[0] : event;
    
    const newNode = {
      id: targetNodeId,
      position: reactFlow.screenToFlowPosition({
        x: clientX,
        y: clientY,
      }),
      type: 'query',
      data: {resourceType: null, labels: [], nodeData: sourceNode.data},
      origin: [0.5, 0.0],
    };

    reactFlow.setNodes((nds) => nds.concat(newNode));
    reactFlow.setEdges((eds) =>
      eds.concat({ id: `${sourceNodeId}-${targetNodeId}`, source: sourceNodeId, target: targetNodeId, style: {strokeWidth: 5} }),
    );
  }, [id, reactFlow]);

  const inlineResources = useCallback(
    (results) => {
      reactFlow.updateNodeData(id, (node) => {
        return {
          ...node.data,
          inline: {...node.data.inline, results: results}
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

  useEffect(() => {
    const loadAttributes = async () => {
      if (data.inline.resourceType === null || data.inline.resourceType === undefined)
        return []
			const attributes = await backend.attributes(data.inline.resourceType)
      setInlineChildPaths(attributes.map(a=> a.path))
    };
    loadAttributes();
  }, [data.inline.resourceType, setInlineChildPaths, backend]);

  useEffect(() => {setInlineParentPaths(getAllJSONPaths(data.data))}, [data.data])

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

  const updateInlineChildPath = useCallback(
    (newValue) => {
      reactFlow.updateNodeData(id, (node) => {
        return { 
          ...node.data,
          inline: {...node.data.inline, childPath: newValue}
        };
      })
    }, [id, reactFlow]
  );

  const updateInlineParentPath = useCallback(
    (newValue) => {
      reactFlow.updateNodeData(id, (node) => {
        return { 
          ...node.data,
          inline: {...node.data.inline, parentPath: newValue}
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

  return (
    <>
      <div className="wrapper">
        <div className="inner">
          <div className="body">
            <div className="resource-row">
              <div className="header-icon">
                <img src={getIconSrc(data.resourceType)}/>
              </div>
              <div className="header-label">
                <div className="header-label-row"><label className="swamp-resource">{data.resourceType}</label></div>
                <div className="header-label-row"><label className="swamp-resource-id">{data.metadata.id}</label></div>
              </div>
            </div>
            <hr/>
            <div className="resource-row">
              <FieldPicker data={data.data} selectedFields={data.selectedFields || []} updateSelectedFields={updateSelectedFields}/>
              {/* to do - button to apply jsonpaths to other siblings */}
              <Button disableRipple aria-describedby={id} sx={{width: "auto"}}> 
                <SettingsEthernetIcon/>
              </Button>
            </div>
            <div className="resource-row">
              <DataDisplay nodeId={id} data={data.data} selectedFields={data.selectedFields || []}/>
            </div>
            <div className="resource-row">
            <Button disableRipple aria-describedby={id} onClick={(_) => setInlineCollapsed(!inlineCollapsed)}>
              <DatasetIcon />
              <span>Inline</span>
            </Button>
            {/* to do - button to apply inline to other siblings */}
            <Button disableRipple aria-describedby={id} sx={{width:"auto"}}> 
              <SettingsEthernetIcon />
            </Button>
            </div>
            <div className="resource-row">
              <Collapse in={!inlineCollapsed} timeout="auto" unmountOnExit>
              <Box sx={{borderLeft: "1px solid gray", marginLeft: "12px", paddingLeft:"12px"}}>
                <QueryWizard 
                nodeId={id} resourceType={data.inline.resourceType} labels={data.inline.labels || []} doSomethingWithResults={inlineResources} onResourceTypeUpdate={updateInlineResourceType}
                setLabels={setInlineLabels}
                join={true}
                childPath={data.inline.childPath} childPaths={inlineChildPaths} onChildPathUpdate={updateInlineChildPath} 
                parentPath={data.inline.parentPath} parentPaths={inlineParentPaths} onParentPathUpdate={updateInlineParentPath}
                getParentVal={(p) => JSONPath({path: p, json: data.data})}
                />
                <FieldPicker data={data.inline.results} selectedFields={data.inline.selectedFields || []} updateSelectedFields={updateInlineSelectedFields} header="Results"/>
                <DataDisplay multiple nodeId={id} data={data.inline.results || []} selectedFields={data.inline.selectedFields || []}/>
              </Box>
              </Collapse>
            </div>
            <div className="resource-row">
            {/* to do - replace edge on drop from App with this button */}
            <Button disableRipple aria-describedby={id} onClick={addQuery}>
              <DatasetLinkedIcon/>
              <span>Link</span>
            </Button>
            </div>
            <Handle type="source" position={Position.Bottom} id="b"/>
            <Handle type="target" position={Position.Top} id="a" style={{opacity: 0}} />
          </div>
        </div>
      </div>
    </>
  );
});
