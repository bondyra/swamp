import React, { memo } from 'react';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import { useQueryStore } from './state/QueryState';
import DisplayRecord from './DisplayRecord';
import { IconButton } from '@mui/material';
import AddBoxIcon from '@mui/icons-material/AddBox';

export default memo(({ vertexId, data }) => {
  const fields = useQueryStore(state => state.fields);
  const addField = useQueryStore(state => state.addField);
  return (
    <Stack>
      {
        fields.filter(f => f.vertexId === vertexId).map(f => {
          return <DisplayRecord fieldId={f.id} value={f.val} data={data} />
        })
      }
			<Box sx={{padding: "0", pl: "2px", pb: "4px"}}>
				<IconButton aria-label="add" sx={{width: "fit-content", padding: "0"}} size="small"
				onClick={() => addField(vertexId)}>
					<AddBoxIcon sx={{width: "16px", height: "16px", padding: "0px", color: "gray"}}/>
				</IconButton>
			</Box>
    </Stack>
  );
});
