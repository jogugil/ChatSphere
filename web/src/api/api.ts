import axios from 'axios';

const apiUrl = 'http://localhost:4000'; // URL de tu servidor

// Login: obtiene un token JWT
export const login = async (nickname: string) => {
  const response = await axios.post(`${apiUrl}/login`, { nickname });
  return response.data;
};

// Enviar mensaje
export const sendMessage = async (token: string, roomId: string, message: string) => {
  await axios.post(
    `${apiUrl}/newmessage`,
    { tokenSession: token, idSala: roomId, mensaje: message },
    { headers: { Authorization: `Bearer ${token}` } }
  );
};

// Obtener mensajes
export const getMessages = async (token: string, roomId: string) => {
  const response = await axios.post(
    `${apiUrl}/getmessages`,
    { tokenSession: token, idSala: roomId },
    { headers: { Authorization: `Bearer ${token}` } }
  );
  return response.data.messages;
};

// Obtener usuarios activos
export const getActiveUsers = async (token: string, roomId: string) => {
  const response = await axios.post(
    `${apiUrl}/getusers`,
    { tokenSession: token, idSala: roomId },
    { headers: { Authorization: `Bearer ${token}` } }
  );
  return response.data.users;
};
