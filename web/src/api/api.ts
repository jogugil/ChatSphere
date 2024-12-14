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

// Enviar mensajes al servidor GoChat 
export const sendMessage = async (nickName: string, token: string, roomId: string, roomName: string, message: string) => {
  const requestData = {
    nickname: nickName,           // Reemplaza con el nombre del usuario
    idSala: roomId,               // ID de la sala
    nameSala: roomName,           // Nombre de la sala
    tokenSession: token,          // El token de sesión
    mensaje: message              // El mensaje que deseas enviar
  };

  try {
    // Realizar la solicitud POST
    const response = await axios.post(
      `${apiUrl}/newmessage`, // URL del servidor GoChat
      requestData,            // Datos de la solicitud en formato JSON
      { headers: { Authorization: `Bearer ${token}` } } // Encabezado de autenticación
    );

    // Si el servidor responde correctamente, devolver el estado de éxito
    if (response.data && response.data.status === 'success') {
      console.log('Mensaje enviado correctamente:', response.data.message);
      return JSON.stringify({
        status: 'success',
        message: 'El mensaje fue enviado correctamente.'
      });
    } else {
      // Si el estado no es success, devolver un error
      return JSON.stringify({
        status: 'error',
        message: 'No se pudo enviar el mensaje. Intenta nuevamente más tarde.'
      });
    }
  } catch (error) {
    // Si hay un error en la solicitud (por ejemplo, error de conexión o servidor fuera de servicio)
    console.error('Error en la solicitud:', error);

    return JSON.stringify({
      status: 'error',
      message: 'El servidor está temporalmente fuera de servicio. Por favor, intenta nuevamente en un instante.'
    });
  }
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
export const getAliveUsers = async (token: string, roomId: string): Promise<string> => {
  const response = await axios.post(
    `${apiUrl}/getusers`,
    JSON.stringify({ tokenSession: token, idSala: roomId }),
    {
      headers: { 
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      }
    }
  );

  const data = response.data as string;
  return data;
};