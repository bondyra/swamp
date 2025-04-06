import { createContext, useContext } from 'react';
import StateManager from './StateManager';

const StateManagerContext = createContext();

export default function StateManagerProvider({ children }) {
  const manager = new StateManager();

  return (
    <StateManagerContext.Provider value={client}>
      {children}
    </StateManagerContext.Provider>
  );
}

export function useStateManager() {
  return useContext(StateManager);
}
