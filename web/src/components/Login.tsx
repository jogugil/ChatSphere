import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom'; // Usamos useNavigate en lugar de useHistory
import { login } from '../api/api';
import { setToken } from '../utils/jwtUtils';

const Login = () => {
  const [nickname, setNickname] = useState('');
  const navigate = useNavigate(); // Usamos useNavigate en lugar de useHistory

  const handleLogin = async () => {
    try {
      const response = await login(nickname);
      setToken(response.token); // Guardamos el token en localStorage o en un contexto
      navigate('/chat');     // Redirigimos al chat después del login
    } catch (error) {
      console.error('Error al iniciar sesión', error);
    }
  };

  return (
    <div>
      <h2>Iniciar sesión</h2>
      <input
        type="text"
        placeholder="Introduce tu nickname"
        value={nickname}
        onChange={(e) => setNickname(e.target.value)}
      />
      <button onClick={handleLogin}>Iniciar sesión</button>
    </div>
  );
};

export default Login;

