import App from './App';
import React from 'react';
import ReactDOM from 'react-dom/client';
 
const rootElement = document.getElementById('root') as HTMLElement;
const root = ReactDOM.createRoot(rootElement);

root.render(
  <React.StrictMode>
    <h1>¡Hola desde index.tsx!</h1>
    <App/>
  </React.StrictMode>
);