import React, { useCallback, useEffect, useRef, useState } from 'react';
import { useReactFlow } from '@xyflow/react';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button'
import Stack from '@mui/material/Stack';
import { Alert, Tooltip } from '@mui/material';
import { createTheme, ThemeProvider } from '@mui/material/styles';
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
import { DagreLayoutProvider } from './DagreLayoutProvider';

const version = "v0.0.1"
let dummyId = 1;
const newDummyId = () => `${dummyId++}`;

const theme = createTheme({
  palette: {
  },
});

const initialNodes = [
  {
    id: 'root-query',
    position: { x: 0, y: 0 },
    type: 'query',
    data: {labels: []},
  }
]

const nodeTypes = {
  resource: Resource,
  query: Query,
};

export const downloadJson = ( json ) => {
  const element = document.createElement( "a" );
  element.setAttribute( "href", "data:application/json;base64," + btoa(json) );
  element.setAttribute( "download", `graph-${ Date.now() }` );
  element.style.display = "none";
  document.body.appendChild( element );
  element.click();
  document.body.removeChild( element );
}

const SwampApp = () => {
  const reactFlowWrapper = useRef(null);
  const reactFlow = useReactFlow();
  const [addDummyNode, setAddDummyNode] = useState(false)
  const [delDummyNode, setDelDummyNode] = useState(false)
  // eslint-disable-next-line
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [rfInstance, setRfInstance] = useState(null);
  // eslint-disable-next-line
  const [graphName, setGraphName] = useState("default")
  const [alrt, setAlrt] = useState("")

  useEffect(() => {
    (async () => {
      await new Promise(res => setTimeout(res, 1000));
      setAlrt("")
    })();
  }, [alrt, setAlrt])

  const saveGraph = useCallback(() => {
    if (rfInstance) {
      const flow = rfInstance.toObject();
      localStorage.setItem(graphName, JSON.stringify(flow));
      setAlrt("Graph saved to local storage!")
    }
  }, [rfInstance, graphName])

  // load graph on mount if exists in local storage
  useEffect(() => {
    async function restoreFlow() {
      const flow = JSON.parse(localStorage.getItem(graphName));
      if (flow) {
        const { x = 0, y = 0, zoom = 1 } = flow.viewport;
        setNodes(flow.nodes || []);
        setEdges(flow.edges || []);
        reactFlow.setViewport({ x, y, zoom });
        setAlrt("Graph restored from local storage!");
      }
    };
    restoreFlow();
  }, [setNodes, setEdges, reactFlow, graphName])

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
    onNodesChange(changes);
    if (changes.some(c=> c.type === "dimensions")){
      // force layout by adding a dummy node
      setAddDummyNode(true)
    }
  }, [onNodesChange, setAddDummyNode])

  useEffect(() => {
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

  useEffect(() => {
    document.addEventListener("keydown", (event) => {
      if (event.key === "s" && (navigator.platform.match("Mac") ? event.metaKey : event.ctrlKey)) {
        event.preventDefault();
        saveGraph();
      }
    }, false);
  }, [saveGraph])

  return (
    <div className="wrapper" ref={reactFlowWrapper} style={{ width: '100vw', height: '100vh' }}>
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
            <Stack sx={{borderRadius: "10px", background: "#141414", padding: "0px"}} direction="row">
              <Stack direction="row" sx={{padding: "5px", border: "1px solid gray"}}>
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
                  <Box sx={{fontSize: "20px", fontWeight: 800, fontFamily: "monospace"}}>Swamp</Box>
                  <Box sx={{fontSize: "16px", fontWeight: 100, fontFamily: "monospace", mt: "4px", ml: "5px"}}>{version}</Box>
                </Stack>
                <Button href="https://github.com/bondyra/swamp">
                  <Box component="img" sx={{height: 24, flexShrink: 0, mr: "5px"}} src={"./github.svg"} />
                </Button>
              </Stack>
              <Tooltip title="THIS DOESNT DO ANYTHING ATM">
                <Stack direction="row" sx={{border: "1px solid gray", borderRight: "0px"}}>
                  <Box key="tab 1" sx={{borderRight: "2px solid gray", fontFamily: "monospace", paddingTop: "10px", paddingLeft: "5px", paddingRight: "5px"}}>MOCK TAB 1</Box>
                  <Box key="tab 2" sx={{borderRight: "1px solid gray", fontFamily: "monospace", paddingTop: "10px", paddingLeft: "5px", paddingRight: "5px"}}>MOCK TAB 2</Box>
                  <Box key="tab 3" sx={{borderRight: "1px solid gray", fontFamily: "monospace", paddingTop: "10px", paddingLeft: "5px", paddingRight: "5px"}}>MOCK TAB 3</Box>
                </Stack>
              </Tooltip>
            </Stack>
          </Panel>
          <Panel position="top-center">
            {alrt && <Alert variant="outlined" severity="success" sx={{color: "lightgreen"}}>{alrt}</Alert>}
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
      <DagreLayoutProvider skipInitialLayout>
        <BackendProvider>
          <SwampApp/>
        </BackendProvider>
      </DagreLayoutProvider>
    </ReactFlowProvider>
  );
}


export default App;
