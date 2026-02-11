import { memo } from 'react';
import { Handle, Position, useReactFlow } from '@xyflow/react';
import Stack from '@mui/material/Stack';

import DataDisplay from './DataDisplay';
import { Box } from '@mui/material';
import { getIconSrc } from './Utils';


// const Button = styled(ButtonBase)(({ theme }) => ({
//   fontSize: 13,
//   width: '100%',
//   textAlign: 'left',
//   paddingBottom: 8,
//   fontWeight: 600,
//   color: '#586069',
//   ...theme.applyStyles('dark', {
//     color: '#8b949e',
//   }),
//   '&:hover,&:focus': {
//     color: '#0366d6',
//     ...theme.applyStyles('dark', {
//       color: '#58a6ff',
//     }),
//   },
//   '& span': {
//     width: '100%',
//   },
//   '& svg': {
//     width: 16,
//     height: 16,
//   },
// }));

export default memo(({ id, data, isConnectable }) => {
  useReactFlow();

  // const updateSelectedFields = useCallback(
  //   (newValue) => {
  //     reactFlow.updateNodeData(id, (node) => {
  //       return { 
  //         ...node.data,
  //         selectedFields: newValue
  //       };
  //     })
  //   }, [id, reactFlow]
  // )

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
            <DataDisplay vertexId={data.vertexId} data={data.result}/>
            <Handle type="source" style={{opacity: 1, height:"1px", width:"1px"}} position={Position.Top} id="b"/>
          </div>
        </div>
      </div>
    </>
  );
});
