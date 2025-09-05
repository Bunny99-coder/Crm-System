// in main.tsx
import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import App from './App.tsx';
import { AuthProvider } from './context/AuthContext';
import { MantineProvider } from '@mantine/core'; 
import { DataProvider } from './context/DataContext';// <-- Import
import './main.css'; // <-- Import our new CSS file

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <MantineProvider defaultColorScheme="dark">
      <BrowserRouter>
        <AuthProvider>
          <DataProvider> {/* <-- Wrap App with DataProvider */}
            <App />
          </DataProvider>
        </AuthProvider>
      </BrowserRouter>
    </MantineProvider>
  </React.StrictMode>
);