import { createContext, useContext } from 'react';
import Backend from './Backend';

const BackendContext = createContext();

export default function BackendProvider({ children }) {
  const client = new Backend();

  return (
    <BackendContext.Provider value={client}>
      {children}
    </BackendContext.Provider>
  );
}

export function useBackend() {
  return useContext(BackendContext);
}
