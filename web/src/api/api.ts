import axios from 'axios';
import { UUID } from '../models/Message';
import { LoginResponse } from '../types/typesComm';

const apiUrl = 'http://localhost:4000'; // URL de tu servidor

interface MessagesResponse {
  messages: string[];
}

interface UsersResponse {
  users: string[];
}
 

 

export const login = async (nickname: string): Promise<LoginResponse> => {
  const apiUrl = import.meta.env.VITE_API_URL;

  try {
    const response = await axios.post(
      `${apiUrl}/login`,
      JSON.stringify({ nickname }),
      {
        headers: {
          'Content-Type': 'application/json',
          'x-gochat': apiUrl,
        },
      }
    );

    const data = response.data;

    // Asegurar un LoginResponse consistente
    return {
      status: data.status || 'nok',
      message: data.message || 'Error desconocido',
      token: data.token || '',
      nickname: data.nickname || nickname,
      idsala: data.idsala || '',
      namesala: data.namesala || '',
    } as LoginResponse;
  } catch (error: any) {
    console.error('Error during login:', error);

    // Detectar errores específicos
    if (error.code === 'ERR_NETWORK') {
      return {
        status: 'nok',
        message: 'El servidor GoChat no está disponible. Disculpe las molestias.',
        token: '',
        nickname: nickname,
        idsala: '',
        namesala: '',
      } as LoginResponse;
    }

    // Manejar otros errores genéricos
    return {
      status: 'nok',
      message: 'Error durante el login. Inténtelo de nuevo más tarde.',
      token: '',
      nickname: nickname,
      idsala: '',
      namesala: '',
    } as LoginResponse;
  }
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
// Función que realiza la petición para obtener los mensajes
export const getMessages = async (token: string, roomId: any, lastSeenMessageId: UUID | null | undefined): Promise<string> => {
  try {
    const response = await fetch('/api/getMessages', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify({
        roomId: roomId,
        lastSeenMessageId: lastSeenMessageId,
      }),
    });

    if (!response.ok) {
      throw new Error('Failed to fetch messages');
    }

    const data = await response.text(); // Cambiado a .text() para obtener el string sin procesar
    return data;
  } catch (error) {
    console.error('Error fetching messages:', error);
    throw error; // Asegúrate de lanzar el error para que pueda ser manejado por el llamador
  }
};
export const getActiveUsers = async (token: string, roomId: string): Promise<string> => {
  const response = await axios.post(
    `${apiUrl}/getusers`,
    JSON.stringify({ tokenSession: token, idSala: roomId }), // Pasar el JSON como string
    { headers: { 
        'Content-Type': 'application/json', // Asegurarse de que el tipo de contenido sea JSON
        Authorization: `Bearer ${token}` 
      } 
    }
  );

  const data = response.data as string; // Asegurarse de que la respuesta se maneje como string
  return data;
};