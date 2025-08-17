import { memo, useEffect, useState } from 'react';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import List from '@mui/material/List';
import Divider from '@mui/material/Divider';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import ListItem from '@mui/material/ListItem';
import Vertex from './Vertex';
import Link from './Link';
import { useBackend } from './BackendProvider';
import { useQueryStore } from './QueryState';
import PreviewFlow from './PreviewFlow';
import { NiceButton } from './NiceButton';

const menuTheme = (theme) => ({
    backgroundColor: "black",
	padding: 0,
    height: "100vh",
    width: "20%",
    color: "gray", 
    borderRight: "1px solid gray", 
    fontFamily: "monospace",
    ["& .MuiDrawer-paper"]: { width: 300, boxSizing: 'border-box', backgroundColor: "#141414" }
})

const listTheme = (theme) => ({
    fontSize: "28px"
})

const listItemTheme = (theme) => ({
    padding: "0px",
    pr: "2px",
    pb: "5px"
})

export default memo(() => {
    const vertices = useQueryStore((state) => state.vertices);
    const links = useQueryStore((state) => state.links);
    const addVertex = useQueryStore((state) => state.addVertex);
    const addLink = useQueryStore((state) => state.addLink);
    const backend = useBackend();
    const [resourceTypes, setResourceTypes] = useState([]);

    useEffect(() => {
        async function update() {
            const newResourceTypes = await backend.resourceTypes();
            setResourceTypes(newResourceTypes);
        }
        update();
    }, [backend, setResourceTypes]);

  return (
    <Stack sx={menuTheme}>
        <Divider sx={{background: "gray"}}/>
        {/* VERTICES */}
        <Box sx={{width: "100%", fontSize: "24px"}}>Resources</Box>
        <List sx={listTheme}>
            {
                vertices.map(
                    (vertex) => 
                    <ListItem key={vertex.id} sx={listItemTheme}>
                        <Vertex vertex={vertex} resourceTypes={resourceTypes}/>
                    </ListItem>
                )
            }
                    <ListItem key="add" sx={listItemTheme}>
                        <NiceButton variant="contained" aria-label="run" backgroundcolor="primary" sx={{height: "24px"}} onClick={addVertex}>
                            <AddCircleOutlineIcon sx={{mr: "5px"}}/>
                            <p>Add resources</p>
                        </NiceButton>
                    </ListItem>
        </List>
        <Divider sx={{background: "gray"}}/>
        {/* EDGES */}
        <Box sx={{width: "100%", fontSize: "24px"}}>Edges</Box>
        <List sx={listTheme}>
            {
                links.map(
                    (link) => 
                    <ListItem key={link.id} sx={listItemTheme}>
                        <Link link={link}/>
                    </ListItem>
                )
            }
                    <ListItem key="add" sx={listItemTheme}>
                        <NiceButton variant="contained" aria-label="run" backgroundcolor="primary" sx={{height: "24px"}} onClick={addLink}>
                            <AddCircleOutlineIcon sx={{mr: "5px"}}/>
                            <p>Add edge rule</p>
                        </NiceButton>
                    </ListItem>
        </List>
        <Divider sx={{background: "gray"}}/>
        {/* PREVIEW */}
        <Box sx={{width: "100%", fontSize: "16px", color:"gray"}}>Query preview</Box>
        <Box sx={{border: "1px dashed gray", mr: "5px"}}><PreviewFlow/></Box>
        <Divider sx={{background: "gray"}}/>
    </Stack>
  );
});
