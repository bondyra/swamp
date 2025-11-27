

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
  const mapLabels = useQueryStore((state) => state.mapLabels);
	const [attributes, setAttributes] = useState(new Map());
  const links = useQueryStore((state) => state.links);
  const updateLink = useQueryStore((state) => state.updateLink);
  const savedLabels = useQueryStore((state) => state.savedLabels);
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
        const matchingSavedLabel = savedLabels.filter(s => s.key === a.path);
        const op = matchingSavedLabel.length > 0 ? matchingSavedLabel[0].op: "==";
        const val = matchingSavedLabel.length > 0 ? matchingSavedLabel[0].val: null;
        return {
          id: randomString(8), key: a.path, op: op, val: val, required: true, allowedValues: a.allowed_values, dependsOn: a.depends_on
        };
      });
      setLabels(vertex.id, newRequiredLabels)
    };
        loadAttributes();
  }, [vertex.id, vertex.resourceType, setLabels, savedLabels, backend]);

  useEffect(() => { // when savedLabels change
    console.log(savedLabels);
    mapLabels(vertex.id, (l) => {
      if (!l.op && !l.val && savedLabels.some(s => s.key === l.key)){
        const sl = savedLabels.filter(s => s.key === s.val)[0];
        return {...l, op: sl.op, val: sl.val}
      }
      return l;
    })
  }, [vertex.id, mapLabels, savedLabels])

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
