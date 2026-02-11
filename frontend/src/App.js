import '@xyflow/react/dist/base.css';

import { useState } from 'react';
import Menu from './Menu';
import Flow from './Flow';
import BackendProvider from './BackendProvider';
import CircularProgress from '@mui/material/CircularProgress';
import PlayCircleIcon from '@mui/icons-material/PlayCircle';
import Stack from '@mui/material/Stack';
import Box from '@mui/material/Box';
import { useQueryStore } from './state/QueryState';
import { NiceVerticalButton } from './ui-elements/NiceButton';
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";


const App = () => {
  const [loading, setLoading] = useState(false);
  const setTriggered = useQueryStore((state) => state.setTriggered);

  return (
    <div className="wrapper" style={{ width: '100vw', height: '100vh' }}>
      <BackendProvider>
        <Stack>
          <ToastContainer position="top-center" autoClose={2000} />
          {/* <Bar/> */}
          <Stack direction="row">
            <Menu/>
            <Box sx={{background:"black"}}>
              <NiceVerticalButton variant="contained" onClick={() => {setLoading(true); setTriggered(true); setLoading(false);}} disabled={loading}>
                  {!loading && <PlayCircleIcon/>}
                  {loading && <CircularProgress color="white" size="20px" sx={{color: "#ffffff"}}/>}
              </NiceVerticalButton>
            </Box>
            <Box sx={{ width: "80%", height: "100vh"}}>
              <Flow/>
            </Box>
          </Stack>
        </Stack>
      </BackendProvider>
    </div>
  );
}

export default App;
