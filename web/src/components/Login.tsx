import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { login } from '../api/api';
import { useAuth } from './AuthContext';
import { LoginResponse } from '../types/typesComm';
import '../styles/Login.css'; // Asegúrate de tener este archivo CSS en el mismo directorio
 
 

const Login: React.FC = () => {
  const [nickname, setNickname] = useState<string>('');

  const handleLogin = async () => {
    try {
      const response: LoginResponse = await login(nickname);
      setToken(response.token);
      setNickname(response.nickname);
      setRoomId(response.idsala);
      setRoomName(response.namesala);
    } catch (error) {
      console.error('Error during login:', error);
    }
  };

  return (
    <div className="login-container">
      <div className="left-column">
        <div className="logo">Gochat</div>
        <div className="banners">
          <p>Banner 1</p>
          <p>Banner 2</p>
        </div>
        <div className="footer">
          <p>&copy; 2024 Gochat</p>
          <p>{new Date().toLocaleDateString()}</p>
        </div>
      </div>
      <div className="right-column">
        <div className="login-box">
          <h2>Iniciar sesión</h2>
          <p>Bienvenido a Gochat, introduce tu nickname para comenzar.</p>
          <input
            type="text"
            placeholder="Introduce tu nickname"
            value={nickname}
            onChange={(e) => setNickname(e.target.value)}
          />
          <button onClick={handleLogin}>Iniciar sesión</button>
        </div>
      </div>
    </div>
  );
};

export default Login;

function setToken(token: string) {
  throw new Error('Function not implemented.');
}


function setRoomId(idsala: string) {
  throw new Error('Function not implemented.');
}


function setRoomName(namesala: string) {
  throw new Error('Function not implemented.');
}
