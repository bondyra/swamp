import React, { memo, useCallback, useEffect, useState } from 'react';
import { styled } from '@mui/material/styles';
import { Handle, Position, useReactFlow } from '@xyflow/react';
import ButtonBase from '@mui/material/ButtonBase';
import Collapse from '@mui/material/Collapse';
import DatasetIcon from '@mui/icons-material/Dataset';
import DatasetLinkedIcon from '@mui/icons-material/DatasetLinked';
import SettingsEthernetIcon from '@mui/icons-material/SettingsEthernet';

import DataDisplay from './DataDisplay';
import QueryWizard from './QueryWizard';
import FieldPicker from './FieldPicker';
import { Box } from '@mui/material';


let linkId = 1;
const newLinkId = () => `link-${linkId++}`;

let inlineLabelId = 1;
const newInlineLabelId = () => inlineLabelId++;

// const themeFunction = (theme) => ({
//   background: "#000000",
//   color: "#ffffff",
//   width: "100%",
//   '& input': {
//     color: "#aaaaaa",
//     padding: "4px 4px",
//     fontFamily: "Monospace"
//   }
// })


function getAllPaths(obj, prefix = '') {
	let paths = [];
	
	if (typeof obj === 'object' && obj !== null) {
		for (let key in obj) {
		  if (Array.isArray(obj[key])) {
        for(var i = 0; i < obj[key].length; i++) {
          let newPrefix = prefix === '' ? `${key}[${i}]` : `${prefix}.${key}[${i}]`;
          paths = paths.concat(getAllPaths(obj[key][i], newPrefix));
        }
		  }
		  else if (typeof obj[key] === 'object'){
        let newPrefix = prefix === '' ? key : `${prefix}.${key}`;
        paths = paths.concat(getAllPaths(obj[key], newPrefix));
		  } else {
        let newPrefix = prefix === '' ? key : `${prefix}.${key}`;
        paths.push(newPrefix);
		  }
		}
	}
	
	return paths;
}


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
  const [inlineChildPaths, setInlineChildPaths] = useState([])
  const [embedding, setEmbedding] = useState(false);
  const onEmbed = (_) => setEmbedding(!embedding)

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
      const [provider, resource] = data.inline.resourceType.split(".")
      const attributes = await fetch(`http://localhost:8000/attributes?provider=${provider}&resource=${resource}`).then(response => response.json());
      setInlineChildPaths(attributes.map(a=> a.path))
    };
    loadAttributes();
  }, [data.inline.resourceType, setInlineChildPaths]);

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

  const updateInlineLabel = useCallback(
  (labelId, data) => {
    reactFlow.updateNodeData(id, (node) => {
      var labels = node.data.inline.labels ?? []
      return { 
        ...node.data,
        inline: {
          ...node.data.inline,
          labels: labels.map(l => {
            if (l.id === labelId){
              return {...l, ...data}
            }
            return l;
          })
        }
      };
    })
  }, [id, reactFlow]);

  const addInlineLabel = useCallback(
    () => {
      reactFlow.updateNodeData(id, (node) => {
        var labels = node.data.inline.labels ?? []
        return { 
          ...node.data,
          inline: {...node.data.inline, labels: labels.concat({id: newInlineLabelId(), key: "", val: ""}) }
        };
      });
    }, [id, reactFlow]
  )

  const deleteInlineLabel = useCallback(
    (labelId) => {
      reactFlow.updateNodeData(id, (node) => {
        var labels = node.data.inline.labels ?? []
        return {
          ...node.data,
          inline: {...node.data.inline, labels: labels.filter(x => x.id !== labelId) }
        };
      });
    }, [id, reactFlow]
  )

  const overwriteInlineLabels = useCallback(
    (newLabels) => {
      reactFlow.updateNodeData(id, (node) => {
        return {
          ...node.data,
          inline: {...node.data.inline, labels: newLabels.map(l => {return {id: newInlineLabelId(), ...l}}) }
        };
      });
    }, [id, reactFlow]
  )

  return (
    <>
      <div className="wrapper">
        <div className="inner">
          <div className="body">
            <div className="resource-row">
              <div className="header-icon">
                <img src={`./icons/aws/${data.resource_type}.svg`} alt={data.resource_type} />
              </div>
              <div className="header-label">
                <div className="header-label-row"><label className="swamp-resource">{data.provider}.{data.resource_type}</label></div>
                <div className="header-label-row"><label className="swamp-resource-id">{data.data.__id}</label></div>
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
            <Button disableRipple aria-describedby={id} onClick={onEmbed}>
              <DatasetIcon />
              <span>Inline</span>
            </Button>
            {/* to do - button to apply embeddings to other siblings */}
            <Button disableRipple aria-describedby={id} sx={{width:"auto"}}> 
              <SettingsEthernetIcon />
            </Button>
            </div>
            <div className="resource-row">
              <Collapse in={embedding} timeout="auto" unmountOnExit>
              <Box sx={{borderLeft: "1px solid gray", marginLeft: "12px", paddingLeft:"12px"}}>
                <QueryWizard 
                nodeId={id} resourceType={data.inline.resourceType} labels={data.inline.labels || []} doSomethingWithResults={inlineResources} onResourceTypeUpdate={updateInlineResourceType}
                addLabel={addInlineLabel} deleteLabel={deleteInlineLabel} updateLabel={updateInlineLabel} overwriteLabels={overwriteInlineLabels}
                join childPath={data.inline.childPath} childPaths={inlineChildPaths} onChildPathUpdate={updateInlineChildPath} parentPath={data.inline.parentPath} parentPaths={getAllPaths(data.data)} onParentPathUpdate={updateInlineParentPath}
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
