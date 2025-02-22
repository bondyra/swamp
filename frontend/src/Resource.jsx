import React, { memo } from 'react';
import { Handle, Position, useReactFlow } from '@xyflow/react';
import { useState } from 'react';
import Button from '@mui/material/Button';
import Collapse from '@mui/material/Collapse';
import TextField from '@mui/material/TextField';
import ReactJsonView from '@microlink/react-json-view'
import {JSONPath} from 'jsonpath-plus';

import EmbedWizard from './EmbedWizard'

const themeFunction = (theme) => ({
  background: "#000000",
  color: "#ffffff",
  width: "100%",
  '& input': {
    color: "#aaaaaa",
    padding: "4px 4px",
    fontFamily: "Monospace"
  }
})

function getByJsonPointers(obj, pointers) {
  function getValueByPointer(obj, pointer) {
      if (pointer === "") return obj; // Root object case
      const parts = pointer.split("/").slice(1); // Remove initial empty string from split
      let current = obj;
      for (const part of parts) {
          if (current === undefined) return undefined;
          current = current[decodeURIComponent(part)];
      }
      return current;
  }

  const result = {};
  for (const pointer of pointers) {
      result[pointer] = getValueByPointer(obj, pointer);
  }
  return result;
}

// const themeFunction = (theme) => ({
//   width: 250,
//   color: "#ffffff",
//   padding: "2px",
//   '& textarea': {
//     color: "#ffffff",
//     fontFamily: "Monospace",
//     fontSize: "8px",
//     padding: "1px",
//     lineHeight: "1.3"
//   },
// })

export default memo(({ id, data, isConnectable }) => {
  const { updateNodeData } = useReactFlow();
  const [embedding, setEmbedding] = useState(false);
  const onEmbed = (_) => setEmbedding(!embedding)

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
                <div className="header-label-row"><label className="swamp-resource-id">{data.resource_id}</label></div>
              </div>
            </div>
            <div className="resource-row">
            <TextField 
              variant="outlined" 
              sx={themeFunction}
              value={data._dataFilter}
              onChange={(event) => {
                var pointers = JSONPath({path: event.target.value, json: data.obj, resultType: "pointer"});
                var newData = data.obj
                if (pointers !== undefined && pointers.length > 0){
                  newData = getByJsonPointers(data.obj, pointers)
                }
                updateNodeData(id, (node) => {
                  return { ...node.data, _dataFilter: event.target.value, _displayedData: newData } ;
                });
              }}
            />
            </div>
            <div className="resource-row">
            <ReactJsonView
              name={null}
              collapsed={1}
              displayDataTypes={false}
              displayObjectSize={false}
              displayArrayKey={false}
              theme="shapeshifter"
              src={data._displayedData}
            />
            </div>
            <div className="resource-row">
              <Button color="primary" size="small" variant="contained" fullWidth onClick={onEmbed}>Embed resources</Button>
            </div>
            <div className="resource-row">
              <Collapse in={embedding} timeout="auto" unmountOnExit>
                <EmbedWizard/>
              </Collapse>
            </div>
            <Handle type="target" position={Position.Top} id="a" style={{opacity: 0}} />
            <Handle type="source" position={Position.Bottom} id="b" style={{ background: "white" }} />
          </div>
        </div>
      </div>
    </>
  );
});
