import React, { memo, useEffect, useState, useCallback } from 'react';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import Grid2 from '@mui/material/Grid2';
import { getIconSrc } from './Utils';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import { useQueryStore } from './QueryState';
import LabelOp from './LabelOp';
import { useBackend } from './BackendProvider';
import { IconButton } from '@mui/material';
import HighlightOffIcon from '@mui/icons-material/HighlightOff';
import Picker from './Picker';


export default memo(({ link }) => {
  const links = useQueryStore((state) => state.links);
  const updateLink = useQueryStore((state) => state.updateLink);
  const removeLink = useQueryStore((state) => state.removeLink);
  const [fromOptions, setFromOptions] = useState([]);
  const [toOptions, setToOptions] = useState([]);
  const updateVertex = useQueryStore((state) => state.updateVertex);
  const vertices = useQueryStore((state) => state.vertices);
  const backend = useBackend();

  const select = useCallback(() => {
    const valueToSet = !link.selected;
    vertices.forEach(v => updateVertex(v.id, {selected: link.fromVertexId === v.id || link.toVertexId === v.id ? valueToSet : false}));
    links.forEach(l => updateLink(l.id, {selected: l.id === link.id ? valueToSet : false}));
  }, [links, updateLink, link, updateVertex, vertices]);

  useEffect(() => {
    const loadAttributes = async () => {
      if (link.from === null || link.from === undefined)
        return;
      const attrs = await backend.attributes(link.from);
      setFromOptions(attrs.map(a => {return {value: a.path, description: a.description}}))
    };
    loadAttributes();
  }, [link.from, setFromOptions, backend]);

  useEffect(() => {
    const loadAttributes = async () => {
      if (link.to === null || link.to === undefined)
        return;
      const attrs = await backend.attributes(link.to);
      setToOptions(attrs.map(a => {return {value: a.path, description: a.description}}))
    };
    loadAttributes();
  }, [link.to, setToOptions, backend]);

  return (
    <Stack direction="row" sx={{border: link.selected ? "3px solid #aaaaff" : "1px solid gray", borderRadius: "10px"}} onClick={select}>
      <Box sx={{ flexGrow: 1 }}>
      <Grid2 container spacing={0} sx={{fontWeight: 600, justifyContent: "center", fontSize: "22px", textAlign: 'center'}}>
        <Grid2 size={4}>
            <Picker 
            value={link.from} 
            getValue={(v) => v.resourceType} getDescription={(v) => v.id} getIcon={(v) => getIconSrc(v)}
            updateData={(v) => updateLink(link.id, {from: v.resourceType, fromVertexId: v.id})} options={vertices}
            valuePlaceholder="From" popperPrompt="Select from vertex"/>
        </Grid2>
        <Grid2 size={4}>
            <ArrowForwardIcon/>
        </Grid2>
        <Grid2 size={4}>
            <Picker 
            value={link.to} 
            getValue={(v) => v.resourceType} getDescription={(v) => v.id} getIcon={(v) => getIconSrc(v)}
            updateData={(v) => updateLink(link.id, {to: v.resourceType, toVertexId: v.id})} options={vertices}
            valuePlaceholder="To" popperPrompt="Select to vertex"/>
        </Grid2>
        <Grid2 size={4} sx={{fontSize: "16px"}}>
            <Picker 
            value={link.fromAttr}
            getValue={(v) => v.value} getDescription={(v) => v.description} getIcon={null}
            updateData={(v) => updateLink(link.id, {fromAttr: v.value})} options={fromOptions}
            valuePlaceholder="From" popperPrompt="Select from attr"/>
        </Grid2>
        <Grid2 size={4} sx={{fontSize: "16px"}}>
          <LabelOp
            op={link.op ?? "="} 
            change={(val) => { updateLink(link.id, {op: val})}}
          />
        </Grid2>
        <Grid2 size={4} sx={{fontSize: "16px"}}>
            <Picker 
            value={link.toAttr}
            getValue={(v) => v.value} getDescription={(v) => v.description} getIcon={null}
            updateData={(v) => updateLink(link.id, {toAttr: v.value})} options={toOptions}
            valuePlaceholder="To" popperPrompt="Select to attr"/>
        </Grid2>
      </Grid2>
      </Box>
      <IconButton sx={{color: "gray", height: "12px", pt: "14px"}} onClick={() => removeLink(link.id)}>
        <HighlightOffIcon/>
      </IconButton>
    </Stack>
  );
});
