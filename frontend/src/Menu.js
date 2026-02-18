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
import { useQueryStore } from './state/QueryState';
import PreviewFlow from './preview/PreviewFlow';
import { NiceButton } from './ui-elements/NiceButton';
import { TextField } from '@mui/material';
import { Button } from '@mui/material';
import YAML from "yaml";
import { toast } from 'react-toastify';

const menuTheme = (theme) => ({
    backgroundColor: "black",
	padding: 0,
    height: "100vh",
    width: "20%",
    color: "gray", 
    borderRight: "1px solid gray", 
    fontFamily: "monospace",
    overflow: "scroll",
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
    const setVertices = useQueryStore((state) => state.setVertices);
    const links = useQueryStore((state) => state.links);
    const setLinks = useQueryStore((state) => state.setLinks);
    const fields = useQueryStore((state) => state.fields);
    const setFields = useQueryStore((state) => state.setFields);
    const addVertex = useQueryStore((state) => state.addVertex);
    const addLink = useQueryStore((state) => state.addLink);
    const backend = useBackend();
    const [resourceTypes, setResourceTypes] = useState([]);
    const [queryInput, setQueryInput] = useState("");

    useEffect(() => {
        var state = {
            vertices: vertices ?? [],
            links: links ?? [],
            fields: fields ?? []
        };
        setQueryInput(YAML.stringify(state, null, 2));
    }, [vertices, links, fields, setQueryInput]);

    const applyYaml = () => {
        try {
            const state = queryInput.trim() === "" ? {} : YAML.parse(queryInput);
            setVertices(state.vertices ?? []);
            setLinks(state.links ?? []);
            setFields(state.fields ?? []);
            toast.success("Graph loaded", {className: "swamp-toast", bodyClassName: "swamp-toast-body"});
        } catch (e) {
            toast.error("You fucked up", {className: "swamp-toast", bodyClassName: "swamp-toast-body"})
        }
    };

    useEffect(() => {
        async function update() {
            const newResourceTypes = await backend.resourceTypes();
            setResourceTypes(newResourceTypes);
        }
        update();
    }, [backend, setResourceTypes]);

    const dupa = (e) => {
        const isEnter = e.key === "Enter";
        const isModifier = e.ctrlKey || e.metaKey;

        if (isEnter && isModifier) {
            e.preventDefault();
            applyYaml();
        }
    };

  return (
    <Stack sx={menuTheme}>
        <Divider sx={{background: "gray"}}/>
        <Stack direction="row">
            <Box sx={{fontSize: "28px", fontWeight: 800, color: "gray", fontStyle: "bold"}}>Swamp</Box>
            <Box sx={{fontSize: "12px", fontWeight: 100, color: "gray",  mt: "8px", ml: "5px"}}>{process.env.REACT_APP_VERSION ?? "dev"}</Box>
            <Button href="https://github.com/bondyra/swamp">
                <Box component="img" sx={{height: 24, flexShrink: 0, mr: "5px"}} src={"./github.svg"} />
            </Button>
        </Stack>
        <Divider sx={{background: "gray"}}/>
        {/* VERTICES */}
        <Box sx={{width: "100%", fontSize: "22px"}}>Resources</Box>
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
        <Box sx={{width: "100%", fontSize: "22px"}}>Edges</Box>
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
        {/* QUERY */}
        <Divider sx={{background: "gray"}}/>
        <Box sx={{width: "100%", fontSize: "16px", color:"gray"}}>YAML (load with Ctrl+Enter)</Box>
        <Stack direction="row">
            <TextField 
            onChange={(e) => {setQueryInput(e.target.value);}}
            multiline value={queryInput} sx={{textarea: { color: "#ffffff", fontFamily: "monospace"}, width: "100%"}}
            onKeyDown={dupa}
            />
        </Stack> 
        <Divider sx={{background: "gray"}}/>
    </Stack>
  );
});
