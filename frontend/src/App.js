import React, { useCallback, useEffect, useRef, useState } from 'react';
import { useReactFlow } from '@xyflow/react';
import Alert from '@mui/material/Alert';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button'
import Stack from '@mui/material/Stack';
import { Tooltip } from '@mui/material';
import DownloadIcon from '@mui/icons-material/Download';
import UploadIcon from '@mui/icons-material/Upload';
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
  const [alertMessage, setAlertMessage] = useState("")
  const [rfInstance, setRfInstance] = useState(null);

  const saveGraph = useCallback(() => {
    if (rfInstance) {
      const flow = rfInstance.toObject();
      const jsn = JSON.stringify(flow)
      downloadJson(jsn)
    }
  }, [rfInstance])

  const loadGraph = useCallback(async (evt) => {
    if (!evt.target.files || evt.target.files.length === 0){
      setAlertMessage("Cannot load, no file selected")
      return
    }
    const graphFile = evt.target.files[0]
    const reader = new FileReader()
    reader.onload = () => {
      const flow = JSON.parse(reader.result);
      if (flow) {
        const { x = 0, y = 0, zoom = 1 } = flow.viewport;
        setNodes(flow.nodes || []);
        setEdges(flow.edges || []);
        reactFlow.setViewport({ x, y, zoom });
      }
      setAlertMessage("")
    }
    reader.readAsText(graphFile)
  }, [setNodes, setEdges, reactFlow])

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

  return (
    <div className="wrapper" ref={reactFlowWrapper} style={{ width: '100vw', height: '100vh' }}>
      <ThemeProvider theme={theme}>
        <ReactFlow 
        nodes={nodes}
        edges={edges}
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
          <Panel position="top-left">
            <Stack sx={{borderRadius: "10px", background: "#141414", padding: "5px"}}>
              <Stack direction="row">
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
                  <Box sx={{fontSize: "20px", fontWeight: 800}}>Swamp</Box>
                  <Box sx={{fontSize: "14px", fontWeight: 100, mt: "6px", ml: "5px"}}>{version}</Box>
                </Stack>
              </Stack>
              <Button href="https://github.com/bondyra/swamp">
                <Box component="img" sx={{height: 16, flexShrink: 0, mr: "5px"}} src={"./github.svg"} />
                <span>GitHub</span>
              </Button>
              <Button onClick={saveGraph}>
                <UploadIcon/>
                <span>Save graph</span>
              </Button>
              <Button>
                <Tooltip title="Click text on the right, not the icon">
                  <DownloadIcon/>
                </Tooltip>
                <label htmlFor="loadGraph">Load graph</label>
                <input id="loadGraph" accept=".json" type="file" style={{display:"none"}} onChange={loadGraph}/>
              </Button>
              {alertMessage && <Alert severity="error">{alertMessage}</Alert>}
            </Stack>
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
