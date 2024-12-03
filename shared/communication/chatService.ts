import axios from 'axios';

const API_BASE = 'http://localhost:8080/api';

export const login = async (nickname: string) => {
  const response = await axios.post(`${API_BASE}/login`, { nickname });
  return response.data;
};

export const sendMessage = async (nickname: string, message: string, token: string, roomId: string) => {
  const response = await axios.post(`${API_BASE}/sendMessage`, {
    nickname,
    message,
    token,
    roomId
  });
  return response.data;
};

export const getMessages = async () => {
  const response = await axios.get(`${API_BASE}/getMessages`);
  return response.data;
};

export const getUsers = async () => {
  const response = await axios.get(`${API_BASE}/getUsers`);
  return response.data;
};
