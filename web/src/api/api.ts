import axios from 'axios';
import { UUID } from '../models/Message'; 
import { LoginResponse } from '../types/typesComm';
import {WebSocketManager} from '../comm/WebSocketManager';

 
const apiUrl = import.meta.env.VITE_API_URL; //'http://localhost:8081';

interface MessagesResponse {
  messages: string[];
}

interface UsersResponse {
  users: string[];
}
 

 

export const login = async (nickname: string): Promise<LoginResponse> => {

  console.log('login _ API URL:', apiUrl);

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
      roomid: data.roomid || '',
      roomname: data.roomname || '',
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
        roomid: '',
        roomname: '',
      } as LoginResponse;
    }

    // Manejar otros errores genéricos
    return {
      status: 'nok',
      message: 'Error durante el login. Inténtelo de nuevo más tarde.',
      token: '',
      nickname: nickname,
      roomid: '',
      roomname: '',
    } as LoginResponse;
  }
};
 
// Enviar mensajes al servidor GoChat 
export const sendMessage = async (nickName: string, token: string, roomId: string, roomName: string, message: string) => {
  const requestData = {
    nickname: nickName,           // Reemplaza con el nombre del usuario
    roomid: roomId,               // ID de la sala
    roomname: roomName,           // Nombre de la sala
    tokensession: token,          // El token de sesión
    message: message              // El mensaje que deseas enviar
  };

  try {
    console.log(`sendMessage:${apiUrl}:${message}`);
    // Realizar la solicitud POST
    const serializedData = JSON.stringify(requestData);
    const response = await axios.post(
      `${apiUrl}/newmessage`, // URL del servidor GoChat
      serializedData,            // Datos de la solicitud en formato JSON
      { headers: {
                Authorization: `Bearer ${token}`,       // Autenticación
                'Content-Type': 'application/json',             // Tipo de contenido
                'x-gochat': apiUrl                              // Cabecera personalizada
        }
      } // Encabezado de autenticación
    );
    console.log(`Valor de status: ${response.data.status}`);
    console.log('Respuesta $ {response.data.status} : completa:', JSON.stringify(response, null, 2));
    // Si el servidor responde correctamente, devolver el estado de éxito
    const responseData = response.data;     
    if (responseData && responseData.status === 'ok') {
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
    console.log(`sendMessage:response:${error}`);
    return JSON.stringify({
      status: 'error',
      message: 'El servidor está temporalmente fuera de servicio. Por favor, intenta nuevamente en un instante.'
    });
  }
};

/*
requestData := map[string]string{
  "operation":     "listmenssage",
  "nickname":      nickname,
  "roomid":        idsala,
  "tokensesion":   token,
  "lastmessageid": ultimoIdMensaje,
  "x_gochat":      "http://localhost:8081",
}
*/
// Obtener mensajes. Función que realiza la petición para obtener los mensajes
// Función para manejar la respuesta del WebSocket
// Definir la función callback para procesar la respuesta
const handleMessage = (event: MessageEvent): any => {
  const message = event.data;

  // Procesar el mensaje recibido
  try {
    const parsedData = JSON.parse(message); // Intentamos parsear el mensaje como JSON
    return parsedData;  // Devuelve el objeto procesado
  } catch (error) {
    console.error('Error al parsear el mensaje:', error);
    return message;  // Si no se puede parsear, devolvemos el mensaje como string
  }
};

export const getMessages = async (
  socketManager: WebSocketManager, 
  nickName: string, 
  token: string, 
  roomId: UUID, 
  roomName: string, 
  lastSeenMessageId: UUID | null | undefined, 
  x_gochat: string
): Promise<string> => {
  
  const requestData = {
    nickname: nickName,
    roomid: roomId,
    roomname: roomName,
    tokensession: token,
    lastmessageid: lastSeenMessageId,
    operation: "listmenssage",
    x_gochat: x_gochat,
  };

  try {
    // Obtener el WebSocket real desde WebSocketManager
    const ws = socketManager.getSocket();

    if (!ws) {
      throw new Error('WebSocket no disponible.');
    }
    socketManager.setOnMessageCallback(handleMessage);
    // Promise que se resuelve cuando la conexión WebSocket se abre
    const wsPromise = new Promise<WebSocket>((resolve, reject) => {
      ws.onopen = () => resolve(ws);
      ws.onerror = (err) => reject(new Error(`getMessages :::: WebSocket error: ${err}`));
    });

    await wsPromise;

    // Enviar los datos a través del WebSocket
    ws.send(JSON.stringify(requestData));

 
    // Esperar a que se reciba el mensaje y se procese con handleMessage
    return new Promise<any>((resolve, reject) => {
      ws.onmessage = (event) => {
        try {
          const result = handleMessage(event);  // Usamos el callback `handleMessage`
          resolve(result);  // Devolvemos el resultado procesado
        } catch (error) {
          reject(error);
        }
      };

      // Manejar errores durante la conexión WebSocket
      ws.onerror = (err) => {
        console.error("Error de WebSocket:", err);
        reject(new Error('Error de WebSocket'));
      };
    });
  } catch (error) {
    console.error('Error al conectar con WebSocket:', error);
    throw error;
  }
};


 

  /*
	requestData := struct {
		RoomId      string `json:"roomid"`
		TokenSesion string `json:"tokensesion"`
		Nickname    string `json:"nickname"`
		Operation   string `json:"operation"`
		X_GoChat    string `json:"x_gochat"`
	}{
		RoomId:      idsala,
		TokenSesion: token,
		Nickname:    nickname,
		Operation:   "listusers",
		X_GoChat:    "http://localhost:8081",
	}
  */
export const getAliveUsers = async (
  socketManager: WebSocketManager, 
  nickname: string, 
  token: string, 
  roomId: string, 
  x_gochat: string
): Promise<string> => {
  // Crear los datos que se enviarán en el WebSocket
  const requestData = {
    roomid: roomId,
    tokensesion: token,
    nickname: nickname,
    operation: "listusers",
    x_gochat: x_gochat
  };

  try {
    // Obtener el WebSocket real desde WebSocketManager
    const ws = socketManager.getSocket();

    if (!ws) {
      throw new Error('WebSocket no disponible.');
    }

    // Promise que se resuelve cuando la conexión WebSocket se abre
    const wsPromise = new Promise<WebSocket>((resolve, reject) => {
      ws.onopen = () => resolve(ws);
      ws.onerror = (err) => reject(new Error(`getAliveUsers ::: WebSocket error: ${err}`));
    });

    await wsPromise;

    // Enviar los datos a través del WebSocket
    ws.send(JSON.stringify(requestData));

    // Esperar y procesar la respuesta
    return new Promise<string>((resolve, reject) => {
      ws.onmessage = (event) => {
        const response = event.data;
        console.log('Respuesta completa:', response);
        resolve(response);
      };

      // Manejar errores durante la conexión WebSocket
      ws.onerror = (err) => {
        console.error("Error de WebSocket:", err);
        reject(new Error('Error de WebSocket'));
      };
    });
  } catch (error) {
    console.error('Error al conectar con WebSocket:', error);
    throw error;
  }
};