

import { memo, useEffect, useCallback, useState } from 'react';
import Stack from '@mui/material/Stack';
import Box from '@mui/material/Box';
import { getIconSrc } from './Utils';
import SingleFieldPicker from './pickers/SingleFieldPicker';
import LabelPicker from './pickers/LabelPicker';

import { useBackend } from './BackendProvider';
import { useQueryStore } from './state/QueryState';

import { randomString } from './Utils';
import { IconButton } from '@mui/material';
import HighlightOffIcon from '@mui/icons-material/HighlightOff';


const vertexTheme = (theme) => ({
	fontSize: "22px", 
  fontWeight:"600",
  width: "100%",
  pl: "5px",
  fontFamily: "monospace",
  // '&:hover': {
	//   backgroundColor: 'rgba(15,15,15,.5)',
	//   borderColor: '#0062cc',
	//   boxShadow: 'none',
	// },
})


export default memo(({ vertex, resourceTypes }) => {
  const vertices = useQueryStore((state) => state.vertices);
  const updateVertex = useQueryStore((state) => state.updateVertex);
  const removeVertex = useQueryStore((state) => state.removeVertex);
  const setLabels = useQueryStore((state) => state.setLabels);
	const [attributes, setAttributes] = useState(new Map());
  const links = useQueryStore((state) => state.links);
  const updateLink = useQueryStore((state) => state.updateLink);
	const backend = useBackend();

  const select = useCallback(() => {
    const valueToSet = !vertex.selected;
    links.forEach(l => updateLink(l.id, {selected: l.fromVertexId === vertex.id || l.toVertexId === vertex.id ? valueToSet : false}));
    vertices.forEach(v => updateVertex(v.id, {selected: false}));
    updateVertex(vertex.id, {selected: valueToSet});
  }, [links, updateLink, vertex, updateVertex, vertices]);

  useEffect(() => {  // when resourceType changes
    const loadAttributes = async () => {
      if (vertex.resourceType === null || vertex.resourceType === undefined)
        return [];
      // re-load attributes
      const attrs = await backend.attributes(vertex.resourceType);
      setAttributes(new Map(attrs.map(a=> [a.path, a])));
      // refresh labels with potentially new required attributes
      const newRequiredLabels = attrs.map(a => {
        return {
          id: randomString(8), key: a.path, val: null, required: true, allowedValues: a.allowed_values, dependsOn: a.depends_on
        }
      });
      setLabels(vertex.id, newRequiredLabels)
    };
        loadAttributes();
  }, [vertex.id, vertex.resourceType, setLabels, backend]);

  return (
    <Stack direction="row" sx={{border: vertex.selected ? "3px solid #aaaaff" : "1px solid gray", borderRadius: "10px"}} onClick={select}>
      <Stack direction="column" sx={vertexTheme}>
          <Box sx={{fontSize: "8px", fontWeight: 100}}>{vertex.id}</Box>
          <SingleFieldPicker value={vertex.resourceType} updateData={(v) => {updateVertex(vertex.id, {resourceType: v})}} options={resourceTypes} getIconSrc={getIconSrc}
          valuePlaceholder="What?" popperPrompt="Select resource to query"/>
          <Box sx={{fontSize: "12px"}}>
            <LabelPicker resourceType={vertex.resourceType} labels={vertex.labels} setLabels={(ll) => updateVertex(vertex.id, {labels: ll})} attributes={attributes}/>
          </Box>
      </Stack>
      <IconButton sx={{color: "gray", height: "12px", pt: "14px"}} onClick={() => removeVertex(vertex.id)}>
        <HighlightOffIcon/>
      </IconButton>
    </Stack>
  );
});
