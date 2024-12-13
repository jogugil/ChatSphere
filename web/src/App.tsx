import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './components/AuthContext'; // Asegúrate de que el path sea correcto

import Login from './components/Login'; // Asegúrate de que la ruta sea correcta
import Chat from './components/Chat';  // Asegúrate de que la ruta sea correcta

const App: React.FC = () => {
  return (
    <AuthProvider>   
 
      <Router>
        <Routes>
          <Route path="/" element={<Login />} />
          <Route path="/chat" element={<Chat />} />
        </Routes>
      </Router>
    </AuthProvider>
);
};

export default App;

