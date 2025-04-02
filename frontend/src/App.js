import React, { useCallback, useEffect, useRef, useState } from 'react';
import { useReactFlow } from '@xyflow/react';
import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import IconButton, {iconButtonClasses} from '@mui/material/IconButton';
import Button from '@mui/material/Button'
import Stack from '@mui/material/Stack';
import { Alert } from '@mui/material';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import AddIcon from '@mui/icons-material/Add';
import {
  ReactFlow,
  addEdge,
  ConnectionLineType,
  useEdgesState,
  useNodesState,
  MiniMap,
  Controls,
  Panel,
  ReactFlowProvider
} from '@xyflow/react';
import '@xyflow/react/dist/base.css';


import BackendProvider from './BackendProvider';
import Resource from './Resource';
import Query from './Query';
// import { D3ForceLayoutProvider } from './D3ForceLayoutProvider';
// import { DagreLayoutProvider } from './DagreLayoutProvider';
import { ELKLayoutProvider } from './ELKLayoutProvider';
import { randomString } from './Utils'
import GraphTab from './GraphTab';


const version = "v0.0.1"
const graphPrefix = "__graph_"
let dummyId = 1;
const newDummyId = () => `${dummyId++}`;

const theme = createTheme({
  palette: {
  },
});

const newRootQueryNode = () => {
  return {
    id: `root-query-${randomString(4)}`,
    position: { x: 0, y: 0 },
    type: 'query',
    data: {labels: []},
  }
}

const initialNodes = [
  newRootQueryNode()
]

const nodeTypes = {
  resource: Resource,
  query: Query,
};

const themeFunction = (theme) => ({
	padding: 0,
  	color: "#ffffff",
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

const SwampApp = () => {
  const reactFlowWrapper = useRef(null);
  const reactFlow = useReactFlow();
  const [addDummyNode, setAddDummyNode] = useState(false)
  const [delDummyNode, setDelDummyNode] = useState(false)
  // eslint-disable-next-line
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);

  const [rfInstance, setRfInstance] = useState(null);
  const [tabs, setTabs] = useState([]);
  const [currentTab, setCurrentTab] = useState("default");
  const [alrt, setAlrt] = useState("");

  // timed display of alert on its change
  useEffect(() => {
    (async () => {
      await new Promise(res => setTimeout(res, 1000));
      setAlrt("")
    })();
  }, [alrt, setAlrt]);

  const refreshTabs = useCallback(() => {
    const t = Object.entries(localStorage).map(it => it[0]).filter(k => k.startsWith(graphPrefix)).map(k => k.replace(graphPrefix, ""));
    setTabs(t ?? ["default"]);
  }, [setTabs]);

  // load tabs from local storage on mount
  useEffect(() => {
    refreshTabs();
  }, [refreshTabs]);

  // add new tab
  const addNewTab = useCallback(() => {
    const newTab = randomString(8);
    localStorage.setItem(`${graphPrefix}${newTab}`, JSON.stringify({nodes: initialNodes, edges: [], viewport: {}}));
    refreshTabs();
    setCurrentTab(newTab);
  }, [refreshTabs]);

  // rename tab
  const renameTab = useCallback((newName) => {
    const content = localStorage.getItem(`${graphPrefix}${currentTab}`);
    localStorage.setItem(`${graphPrefix}${newName}`, content);
    localStorage.removeItem(`${graphPrefix}${currentTab}`);
    refreshTabs();
    setCurrentTab(newName);
  }, [setCurrentTab, currentTab, refreshTabs]);

  // // switch tabs
  // const switchTab = useCallback((name) => {
  //   setCurrentTab(name);

  // }, []);

  // load graph on start (if exists)
  useEffect(() => {
    async function restoreFlow() {
      const flow = JSON.parse(localStorage.getItem(`${graphPrefix}${currentTab}`));
      if (flow) {
        const { x = 0, y = 0, zoom = 1 } = flow.viewport;
        setNodes(flow.nodes || []);
        setEdges(flow.edges || []);
        reactFlow.setViewport({ x, y, zoom });
        setAlrt("Graph restored from local storage!");
      }
    };
    restoreFlow();
  }, [setNodes, setEdges, reactFlow, currentTab]);

  // save graph
  const saveGraph = useCallback(() => {
    if (rfInstance) {
      const flow = rfInstance.toObject();
      localStorage.setItem(`${graphPrefix}${currentTab}`, JSON.stringify(flow));
      setAlrt(`Graph saved to local storage!`)
    }
  }, [rfInstance, currentTab]);

  // RF stuff
  const onConnect = useCallback(
    (params) =>
      setEdges((eds) =>
        addEdge(
          { ...params, type: ConnectionLineType.Step, animated: true },
          eds,
        ),
      ),
    [setEdges],
  );

  const onNodesChangeExt = useCallback((changes) => {
    console.log(changes);
    onNodesChange(changes);
    if (changes.some(c=> c.type === "dimensions")){ //dimensions for dagre/d3!!!!!!!!!
      // force layout by adding a dummy node
      setAddDummyNode(true)
    }
  }, [onNodesChange, setAddDummyNode])

  useEffect(() => {
    console.log(`del: ${delDummyNode} add: ${addDummyNode}`)
    if (delDummyNode) {
      reactFlow.deleteElements({nodes: nodes.filter(n => n.id.startsWith("__DUMMY__"))})
      setDelDummyNode(false)
    }
    if (addDummyNode) {
      reactFlow.addNodes([{id: `__DUMMY__${newDummyId()}`, position: { x: 0, y: 0 }}])
      setAddDummyNode(false)
      setDelDummyNode(true)  // immediately mark this node for deletion
    }
  }, [reactFlow, nodes, addDummyNode, delDummyNode, setAddDummyNode, setDelDummyNode])

  return (
    <div className="wrapper" ref={reactFlowWrapper} style={{ width: '100vw', height: '100vh' }} onKeyDown={evt => {
      if (evt.key === "s" && (navigator.platform.match("Mac") ? evt.metaKey : evt.ctrlKey)) {
        evt.preventDefault();
        saveGraph();
      }
    }}>
      <ThemeProvider theme={theme}>
        <ReactFlow 
        nodes={nodes}
        edges={edges}
        deleteKeyCode={null}
        onNodesChange={onNodesChangeExt}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        onInit={setRfInstance}
        connectionLineType={ConnectionLineType.Step}
        nodeTypes={nodeTypes}
        nodesDraggable={false}
        colorMode={"dark"}
        fitView
        >
          <Panel position="top" style={{ width: "100%" }}>
          <AppBar>
            <Stack sx={{background: "#141414", padding: "0px"}} direction="row">
              <Stack direction="row" sx={{padding: "5px"}}>
                <Box
                  component="img"
                  sx={{
                    height: 24,
                    flexShrink: 0,
                    borderRadius: '3px',
                    padding: "0px",
                    mr: "5px",
                    mt: "5px"
                  }}
                  src={"./asset.svg"} alt=""
                />
                <Stack direction="row">
                  <Box sx={{fontSize: "24px", fontWeight: 800, fontStyle: "bold", fontFamily: "monospace"}}>Swamp</Box>
                  <Box sx={{fontSize: "16px", fontWeight: 100, fontFamily: "monospace", mt: "8px", ml: "5px"}}>{version}</Box>
                </Stack>
                <Button href="https://github.com/bondyra/swamp">
                  <Box component="img" sx={{height: 24, flexShrink: 0, mr: "5px"}} src={"./github.svg"} />
                </Button>
              </Stack>
              <Stack direction="row">
                {
                  tabs.map(t => {
                    return (
                      <GraphTab name={t} selected={t === currentTab} onSelect={setCurrentTab} onEditEnd={renameTab}/>
                    );
                  })
                }
                <IconButton key={`tab-add}`} onClick={addNewTab} sx={themeFunction}>
                  <AddIcon/>
                </IconButton>
              </Stack>
            </Stack>
            </AppBar>
          </Panel>
          <Panel position="top-center">
            {alrt && <Alert variant="outlined" severity="success" sx={{color: "lightgreen"}}>{alrt}</Alert>}
          </Panel>
          <Panel position="top-right">
            <Button onClick={() => setNodes([...nodes, newRootQueryNode()])}>Add new root node</Button>
          </Panel>
          <MiniMap />
          <Controls />
        </ReactFlow>
      </ThemeProvider>
    </div>
  );
}

const App = () => {
  return (
    <ReactFlowProvider>
      <ELKLayoutProvider skipInitialLayout>
        <BackendProvider>
          <SwampApp/>
        </BackendProvider>
      </ELKLayoutProvider>
    </ReactFlowProvider>
  );
}


export default App;
