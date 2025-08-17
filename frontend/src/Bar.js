import '@xyflow/react/dist/base.css';

import React, { useCallback, useEffect, useState } from 'react';
import IconButton, {iconButtonClasses} from '@mui/material/IconButton';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import Divider from '@mui/material/Divider';
import { Button } from '@mui/material';
import Stack from '@mui/material/Stack';
import Box from '@mui/material/Box';
import GraphTab from './GraphTab';
import { useQueryStore } from './QueryState';

import { randomString } from './Utils'

const themeFunction = (theme) => ({
    padding: 0,
    color: "gray",
    '& input': {
        color: "#ffffff",
        padding: "0px 0px 0px 5px",
        fontFamily: "Monospace",
        fontSize: "10px",
        height: "20px"
    },
    [`&.${iconButtonClasses.root}`]: {
        width: "100px",
        minWidth:"10px",
    },
})

const tabPrefix = "__swamp_"

export function listTabs() {
    return Object.entries(localStorage).map(it => it[0]).filter(k => k.startsWith(tabPrefix)).map(k => k.replace(tabPrefix, ""));
}

export function loadTab(name) {
    return localStorage.getItem(`${tabPrefix}${name}`);
}

export function saveTab(name, content) {
    localStorage.setItem(`${tabPrefix}${name}`, content);
}

export function removeTab(name) {
    localStorage.removeItem(`${tabPrefix}${name}`);
}


const barTheme = (theme) => ({
    backgroundColor: "black",
	  padding: 0,
    height: "64px",
    color: "gray", 
    borderRight: "1px solid gray", 
    fontFamily: "monospace",
    ["& .MuiDrawer-paper"]: { width: 300, boxSizing: 'border-box', backgroundColor: "#141414" }
})


const Bar = () => {
  const [tabs, setTabs] = useState([]);
  const [currentTab, setCurrentTab] = useState("default");
  const vertices = useQueryStore((state) => state.vertices);
  const setVertices = useQueryStore((state) => state.setVertices);
  const links = useQueryStore((state) => state.links);
  const setLinks = useQueryStore((state) => state.setLinks);
  const fields = useQueryStore((state) => state.fields);
  const setFields = useQueryStore((state) => state.setFields);
  const setAlert = useQueryStore((state) => state.setAlert);

  // actual state manipulation functions
  const saveState = useCallback((name) => {
    var state = {
      vertices: vertices ?? [],
      links: links ?? [],
      fields: fields ?? [],
    };
    const content = JSON.stringify(state);
    saveTab(name, content)
    setAlert(`Graph saved to ${name}`);
  }, [vertices, links, fields, setAlert]);

  const loadState = useCallback((name) => {
    const content = loadTab(name);
    if(content){
      const state = JSON.parse(content);
      setVertices(state["vertices"] ?? []);
      setLinks(state["links"] ?? []);
      setFields(state["fields"] ?? []);
      setAlert(`Graph "${name}" loaded`);
    }
  }, [setVertices, setLinks, setFields, setAlert])

    // when current tab changes
    useEffect(() => {
      const tabs = listTabs();
      if(!tabs || tabs.length === 0) {
        saveState("default");
        setCurrentTab("default");
      } else {
        setTabs(tabs);
      }
    }, [currentTab, setTabs, setCurrentTab, saveState]);

    // load graph on start (if exists)
    useEffect(() => {
      if(currentTab) loadState(currentTab);
    }, [currentTab, loadState]);
    
    // obvious tab actions
    const addNewTab = useCallback(() => {
      const newTab = `0_${randomString(4)}`;
      saveState(newTab);
      setCurrentTab(newTab);
    }, [saveState]);

    const renameTab = useCallback((newName) => {
      saveState(newName);
      removeTab(currentTab);
      setCurrentTab(newName);
    }, [setCurrentTab, saveState, currentTab]);

  return (
    <Stack sx={barTheme} direction="row">
      <Stack direction="row" sx={{alignItems: "center"}}>
          <Box
              component="img"
              sx={{
              height: 36,
              flexShrink: 0,
              borderRadius: '3px',
              padding: "0px",
              mr: "5px",
              mt: "5px"
              }}
              src={"./asset.svg"} alt=""
          />
          <Box sx={{fontSize: "36px", fontWeight: 800, color: "gray", fontStyle: "bold"}}>Swamp</Box>
          <Box sx={{fontSize: "16px", fontWeight: 100, color: "gray",  mt: "8px", ml: "5px"}}>{process.env.REACT_APP_VERSION ?? "dev"}</Box>
          <Button href="https://github.com/bondyra/swamp">
            <Box component="img" sx={{height: 24, flexShrink: 0, mr: "5px"}} src={"./github.svg"} />
          </Button>
          <Divider orientation='vertical' sx={{background: "gray"}}/>
          <Stack sx={{background: "#000", padding: "0px", height: "100%"}} direction="row">
              <Stack direction="row">
              {
                  tabs.sort().map(t => {
                  return (
                      <GraphTab key={t} name={t} selected={t === currentTab} onSelect={setCurrentTab} onEditEnd={renameTab}/>
                  );
                  })
              }
              <IconButton key={`tab-add`} onClick={addNewTab} sx={themeFunction}>
                  <AddCircleOutlineIcon/>
              </IconButton>
              </Stack>
          </Stack>
      </Stack>
    </Stack>
  );
}

export default Bar;
