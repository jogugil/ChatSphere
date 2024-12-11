import React from 'react';
import { useHistory } from 'react-router-dom';

const LogoutButton = () => {
  const history = useHistory();

  const handleLogout = () => {
    // Eliminar el token y redirigir al login
    localStorage.removeItem('token');
    history.push('/login');
  };

  return <button onClick={handleLogout}>Salir</button>;
};
