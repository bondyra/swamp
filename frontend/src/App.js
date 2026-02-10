import '@xyflow/react/dist/base.css';

import { useEffect, useState } from 'react';
import { Alert } from '@mui/material';
import Bar from './Bar';
import Menu from './Menu';
import Flow from './Flow';
import BackendProvider from './BackendProvider';
import CircularProgress from '@mui/material/CircularProgress';
import PlayCircleIcon from '@mui/icons-material/PlayCircle';
import Stack from '@mui/material/Stack';
import Box from '@mui/material/Box';
import { useQueryStore } from './state/QueryState';
import { NiceVerticalButton } from './ui-elements/NiceButton';


const App = () => {
  const [loading, setLoading] = useState(false);
  const setTriggered = useQueryStore((state) => state.setTriggered);
  const alert = useQueryStore((state) => state.alert);
  const alertType = useQueryStore((state) => state.alertType);
  const setAlert = useQueryStore((state) => state.setAlert);

  // timed display of alert on its change
  useEffect(() => {
    (async () => {
        if (alert === "")
          return;
        const timer = setTimeout(() => {
          setAlert("");
        }, 5000);
        return (() => clearTimeout(timer));
    })();
  }, [alert, setAlert]);

  return (
    <div className="wrapper" style={{ width: '100vw', height: '100vh' }}>
      <BackendProvider>
        <Stack>
          {
            alert && 
            <Stack sx={{background:"black", position: "absolute", alignSelf: "center", width: "auto"}}>
              <Alert variant="outlined" severity={alertType} 
              sx={{color: alertType === "error" ? "red": "lightblue", fontFamily: "monospace", fontWeight: 600, fontSize: "16px"}}>
                {alert}
              </Alert>
            </Stack>
          }
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
