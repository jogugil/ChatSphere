import React, { useState, useEffect, useContext, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { getMessages, sendMessage, getAliveUsers } from '../api/api';
import { useAuth } from './AuthContext';  // Importamos el contexto de autenticación
import { Room } from '../models/Room';
import { Message, UUID } from "../models/Message";
import { User } from "../models/User";
import { MessageResponse } from '../types/typesComm';
import { WErrorMessage, } from "./ErrorMessage"; // Componente para mostrar errores
import {ResponseUser, UsersAlive} from '../types/typesComm'
import '../styles/chat.css';
import BannerProgramming from './BannerProgramming';
import BannerCloud from './BannerCloud';
import prohibitedWords from "./prohibitedWords";
import { getClientInformation } from '../utils/ClientData'; // Ajusta la ruta según tu estructura de carpetas
import {WebSocketManager} from '../comm/WebSocketManager';

const Clock = () => {
  const [currentTime, setCurrentTime] = useState<string>('');

  // Function to update the current time
  const updateTime = () => {
    const now = new Date();
    const formattedTime = now.toLocaleString(); // You can customize the format as needed
    setCurrentTime(formattedTime);
  };

  useEffect(() => {
    // Update time every second
    const intervalId = setInterval(updateTime, 1000);

    // Set initial time
    updateTime();

    // Clear the interval on component unmount
    return () => clearInterval(intervalId);
  }, []);

  return (
    <p><strong>Fecha:</strong> {currentTime}</p>
  );
};

const Chat = () => {
  // Estructuras que define los usuarios  activos que vienen del Servidor Gochat
  interface  aliveUsers {
    Nickname: string;
    LastActionTime: string;
  };

  interface ResponseUser  {
    Status: string;
    Message: string;
    TokenSesion: string;
    Nickname: string;
    RoomId: string;
    AliveUsers: aliveUsers[];
  };
  const navigate = useNavigate(); //Navegar por la web de Gochat

  //WebSocket
  const [connectionError, setConnectionError] = useState<string | null>(null);
  const apiIP          = import.meta.env.VITE_IP_SERVER_GOCHAT;
  const apiPORT        = import.meta.env.VITE_PORT_SERVER_GOCHAT;
  const socketUrl      = `ws://${apiIP}:${apiPORT}/ws`;
  const socketMessages = useState<WebSocketManager | null>(null); 
  const socketUsers    = useState<WebSocketManager | null>(null); 

  const [isMessageSendable, setIsMessageSendable] = useState(false); // Para habilitar/deshabilitar input y botón
  const [messages, setMessages]         = useState<Message[]> ([]);        // Estado para los mensajes
  const [messageText, setMessageText]   = useState<string> ('');           // Estado para el texto del mensaje
  const [aliveUsers, setAliveUsers]     = useState<string[]> ([]);         // Estado para los usuarios activos
 
  
 
  //Ventana de error /información
  const [minimized, setMinimized]       = useState (false);
  const closeErrorMessage = () => {
    setShowError(false);
    setIsMessageSendable(true); //  habilitar el envío de mensajes
  };  
  const minimizeErrorMessage = () => {
    setMinimized(true);
  }; 
  const restoreErrorMessage  = () => setMinimized (false);  // Función para restaurar el mensaje de error
  // Función para mostrar el error
  const showErrorModal = (message: string) => {
    setErrorMessage(message);
    setShowError(true);
    setIsMessageSendable(false); // Deshabilitar el envío de mensajes
  };

  //Control coor fondo zona central
  const [isDarkMode, setIsDarkMode] = useState(false);
  const toggleTheme = () => {
    setIsDarkMode(!isDarkMode);
  };
  
  //Logica de usuario y chat
  const { token, nickName, roomId, roomName } = useAuth();  // Obtener el usuario y el token del contexto
  const [userChat, setUserChat]               = useState<User | null>(null);  // Estado para el objeto User
  const [room, setRoom]                       = useState<Room | null>(null);  // Estado para el objeto Room
  const [showError, setShowError]             = useState(false);
  const [errorMessage, setErrorMessage]       = useState('');
  const [isErrorActive, setIsErrorActive]     = useState(false);
  const [initialized, setInitialized]         = useState(false);
 
  // Timeout definido en el archivo de entorno (.env). Cada cuanto tiempo el polling realia la petición de mensajes/usuarios
  const timeout = parseInt(import.meta.env.VITE_TIMEOUT, 10000) || 50000;

  // Controlo input envío de mensjaes
  // Tipo explícito para las claves válidas
  type EscapeChar = "<" | ">" | "&" | "\"" | "'";

  // Mapa de caracteres a escapar
  const escapeMap: Record<EscapeChar, string> = {
    "<": "&lt;",
    ">": "&gt;",
    "&": "&amp;",
    "\"": "&quot;",
    "'": "&#39;",
  };

  // Función para escapar caracteres peligrosos
  const escapeHTML = (text: string): string => {
    return text.replace(/[<>&"']/g, (char) => escapeMap[char as EscapeChar] || char);
  };

  // Eliminamos cualquier elemento que no sea carácter alfanumérico
  const sanitizeMessage = (text: string): string => {
    return text.replace(/[^a-zA-Z0-9\s]/g, "");
  };

  // Verifica si el mensaje contiene palabras prohibidas
  const containsProhibitedWords = (text: string) => {
    return prohibitedWords.some((word) => text.toLowerCase().includes(word));
    };
  
  //envia el mensaje que el usuario pone en el input
  const handleSendMessage = async () => {
    if (!messageText.trim() || !isMessageSendable) return; // No enviar si el mensaje está vacío o el envío está deshabilitado
   
      try { 
        if (userChat && room) {
          const sanitizedMessage = sanitizeMessage(messageText.trim());
          const escapedMessage = escapeHTML(sanitizedMessage);

          if (!escapedMessage) {
            showErrorModal("No puedes enviar un mensaje vacío.");
 
            return;
          }

          if (containsProhibitedWords(escapedMessage)) {
            showErrorModal("El mensaje contiene lenguaje prohibido.");
 
            return;
          }

          const response = await sendMessage(userChat.nickname, userChat.token, room.roomId, room.roomName, messageText);
          setMessageText(''); // Limpiamos el input
          
        }
      } catch (err) {
        console.error("Error al enviar el mensaje:", err);
        showErrorModal('No se pudo enviar el mensaje. Intente de nuevo.');
        throw new Error('El servicio de chat no está disponible. Disculpe las molestias. Por favor, intente ingresar nuevamente más tarde.');
 
      }
   
};

// Función para manejar la carga periódica de datos
function startPeriodicUpdates (userChat: any) {
  // Verificar que el usuario esté disponible
  if (!userChat) {
    console.error("Usuario no disponible.");
    return;
  }

  // Función para realizar ambas tareas: Petición lsitado mensajes y petición listado de usaurios
  const updateData = () => {
    try {
      // Cargar mensajes del chat
      console.log("Cargando mensajes del chat...");
      loadMessages(userChat);

      // Cargar usuarios activos
      console.log("Cargando usuarios activos...");
      loadAliveUsers(userChat);
    } catch (error) {
      console.error("Error durante la actualización periódica:", error);
    }
  };

  // Ejecutar inmediatamente antes de iniciar el intervalo
  updateData ();

  // Configurar la actualización periódica
  const intervalId = setInterval (updateData, timeout);

  // Retornar el ID del intervalo para permitir detenerlo si es necesario
  return intervalId;
}

//Función que proicesa el mensaje JSON del servidor GoChat. lista de usaurios activos
const loadAliveUsers = async (userObject: User ) => {
    try {
      const datosCliente = await getClientInformation();
      console.log('Usuarios activos:datosCliente:', datosCliente);
      if ( socketUsers &&  socketUsers.isConnected) {
        const response = await getAliveUsers ( socketUsers, userObject.nickname, userObject.token, userObject.roomId, datosCliente );
      
        // Parsear la respuesta JSON
        const data: ResponseUser = JSON.parse(response);  // Asegúrate de que la respuesta es un JSON
        if (data.Status === 'success' && data.AliveUsers) {
          // Extraer los nicknames de los usuarios activos
          const nicknames = data.AliveUsers.map(user => user.Nickname);
          setAliveUsers(nicknames);  // Establecer el estado con los nicknames
          console.log('Usuarios activos:', nicknames);
        } else {
          showErrorModal('Error al obtener los usuarios activos');
   
          console.error('Error al obtener usuarios activos:', data.Message);
          return;
        }
      }     
    } catch (err) {
      showErrorModal('Error al obtener los usuarios activos: ' + err); 
      console.error('Error al obtener usuarios activos:', err);
      return;
    }
  };
  // Obtener mensajes históricos al cargar la página
  const loadMessages = async (userObject: any) => {
    // Generar un UUID vacío para el primer request
    const emptyUUID: UUID = "00000000-0000-0000-0000-000000000000";
      
    try {
        if (!room) {
          // Mostrar una ventana emergente de error       
          console.error('Error , no se creo el objeto room');
          throw new Error('El servicio de chat no está disponible. Disculpe las molestias. Por favor, intente ingresar nuevamente más tarde.');
        }
        
        const messageID = (room.lastIDMessageId ? room.lastIDMessageId : emptyUUID);
        
        // Hacer la   llamada a la API para obtener los mensajes
        const datosCliente = await getClientInformation();
        console.log('loadMessages:datosCliente:', datosCliente);
        if (socketMessages && socketMessages.isConnected) {
          const messageList  = await getMessages(socketMessages, userObject.NickName, userObject.token, userObject.idRoom, room.roomName, messageID, datosCliente);
          console.log ("loadMessages:",messageList);
          // Convertir los mensajes a objetos
          room.updateMessages(messageList);
  
          // Actualizar el estado de los mensajes
          setMessages(room.messageList);
  
          console.log('Mensajes cargados:', room.messageList);
        }

      } catch (err) {
        showErrorModal('El Servicio de Chat está temporalmente cerrado. Intente logarse más tarde');
        console.error('Error al cargar los mensajes:', err);
 
      }


  };
  // Función para manejar el cambio en el campo de mensaje
  const handleMessageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setMessageText(e.target.value);
    console.log('Cambio en el mensaje:', e.target.value);
    setIsMessageSendable(e.target.value.trim() !== ''); // Habilitar el botón solo si hay texto
  };
  //Función para salir del GoChat
  const logoutAndRedirect = (event: React.MouseEvent<HTMLButtonElement>) => {
    // Lógica para destruir los objetos Room, User, Messages
    localStorage.removeItem("Room");
    localStorage.removeItem("User");
    localStorage.removeItem("Messages");
  
    // Redirigir al login
    window.location.href = "/"; // Redirige a la página de login
  };
 

  const authenticateUser = () => {
    if (nickName && roomId && roomName && token) {
      const user = new User(nickName, 'Alive', roomId, roomName, token);
      setUserChat(user);
      setIsMessageSendable(true); //  habilitar el envío de mensajes
      console.log('Usuario autenticado:', user);
    } else {
      showErrorModal('Datos del usuario no válidos');
 
      navigate('/');
    }
  };
  
  const initializeRoom = (user: User | null) => {
    if (user) {
      const roomU = new Room(user.roomId, user.roomName);
      setRoom(roomU);
      console.log('Sala creada:', roomU);
    }  
  };
  
  // UseEffect para crear userChat
  useEffect(() => {
    if (nickName && token && roomId && roomName) {
      authenticateUser (); // Establecer userChat cuando todos los valores estén disponibles
    }
  }, [token, nickName, roomId, roomName]); // Depende de los valores del contexto

  useEffect(() => {
    async function authenticateAndInitialize() {
      await authenticateUser(); // Autenticación del usuario
      console.log ("   Se llama a authenticateUser")
      if (userChat && socketMessages?.isConnected && socketUsers?.isConnected) {
        console.log (" Los objetos userChat creados. Se llama a initializeRoom")

        // Solo se procede con la inicialización si el usuario está autenticado y el WebSocket está conectado
        initializeRoom(userChat); // Se crea la sala
      }
    }
    
    authenticateAndInitialize();
  }, [userChat]); // Depende solo de userChat para ejecutarse

  // UseEffect para crear la sala, solo cuando userChat está listo y no se ha inicializado antes
  useEffect(() => {      
      initializeRoom (userChat);
  }, [userChat, initialized]); // Este efecto depende de userChat y de initialized

  useEffect(() => {
    // Intentar conectar al WebSocket de forma asíncrona
    const connectWebSocket = async () => {
      try {
        let socketMessages = new WebSocketManager(socketUrl, "messagesSocket");
        let sockeUsers     = new WebSocketManager(socketUrl, "usersSocket");
      }catch (error: unknown) {
        if (error instanceof Error) {
          setConnectionError(error.message); // Ahora accedemos a error.message de forma segura
          console.error('Error al conectar al WebSocket:', error);
        } else {
          // En caso de que error no sea una instancia de Error
          setConnectionError('Error desconocido al conectar al WebSocket');
          console.error('Error desconocido al conectar al WebSocket', error);
        }
      }
    };

    connectWebSocket();

    return () => {
      if (wsManagerRef.current && typeof wsManagerRef.current.close === 'function') {
        wsManagerRef.current.close();
      }
    };
  }, [room]);

  // Este useEffect comienza los updates periódicos una vez que userChat, room y WebSocket estén listos
  useEffect(() => {
    // Solo ejecutar startPeriodicUpdates cuando userChat, room, y el WebSocket estén listos
    if (userChat && room && wsManagerRef.current?.isConnected && !initialized) {
      console.log("Se llaman los updates periódicos del chat");
      //startPeriodicUpdates(userChat); // Activa los mensajes periódicos
      setInitialized(true); // Marca como inicializado para evitar ejecuciones futuras
    }
  }, [userChat, room, wsManagerRef.current?.isConnected, initialized]); // Se ejecuta cuando cambian userChat, room o WebSocket

  // Asegúrate de que tanto el usuario como el WebSocket y la sala estén listos antes de mostrar el chat
  if (!userChat || !wsManagerRef.current?.isConnected || !room) {
    return <div>Cargando...</div>; // O cualquier otro indicador de que el chat no está listo
  }
  else{
    return (
      <div className="chat-container">
        <div className="chat-left-column">
          <div className="chat-title">ChatSphere</div>
          <div className="chat-subtitle">GoChat ZeroMQ</div>

          <div className="chat-logo-container">
            <div className="chat-logo">
              <img src="/images/logo.webp" alt="Logo" />
            </div>
          </div>

          <div className="chat-container">
            <div className="chat-banners">
              <div className="chat-banner-programming">
                <BannerProgramming  
                  titleSlogan="Desarrollos ágiles para tus aplicaciones"
                  subtitleSlogan="Escalabilidad y elasticidad eficientes"
                  imageUrl="/images/pattern.png"
                />
              </div>
            </div>
          </div>

          <div className="footer">
            <p>&copy; 2024 José Javier Gutiérrez Gil</p>
            <p className="email-style">&copy; jogugil@gmail.com // jogugi@posgrado.upv.es</p>
          </div>
        </div>

        <div className={`chat-room ${isDarkMode ? 'chat-room-dark' : 'chat-room-light'}`}>
          <div className="chat-content">
            <div className="messages-display">
              <div className="theme-toggle-btn" onClick={toggleTheme}>
                Cambiar Tema
              </div>
              <h3>Mensajes</h3>
              <ul>
                {messages.map((msg, index) => (
                  <li key={index}>
                    <strong>{msg.nickname}</strong>: {msg.message}
                  </li>
                ))}
              </ul>
            </div>

            <div className="input-section">
              <button className="send-btn"
                onClick={handleSendMessage}
                disabled={!isMessageSendable} // Deshabilitar el botón si no se puede enviar
              >
                Enviar
              </button>
              <input
                type="text"
                value={messageText}
                onChange={handleMessageChange}
                disabled={!isMessageSendable} // Deshabilitar si no es enviable
                placeholder="Escribe tu mensaje..."
              />
            </div>
          </div>
        </div>

        <div className="chat-right-column">
          <div className="chat-metricsbox">
            <Clock />
            <p><strong>Mensajes enviados:</strong> {messages.length}</p>
            <p><strong>Usuarios activos:</strong> {aliveUsers.length}</p>
          </div>
          <div className="chat-active-users-box">
            <div className="chat-active-users-header">
              <h3>Usuarios Activos</h3>
            </div>
            <ul>
              {aliveUsers.map((user, index) => (
                <li key={index}>{user}</li>  
              ))}
            </ul>
          </div>

          <div className="chat-banner-cloud">
            <BannerCloud 
              imageUrl="../../public/images/cloudcomm.png"
              titleSlogan="Inteligencia Aplicada en la Nube"
              subtitleSlogan="Soluciones avanzadas en Cloud, Clusters y Serverless para un futuro más eficiente"
            />
          </div>

          <div className="chat-error-message">
            <WErrorMessage
              message={errorMessage || ''}
              showError={showError}
              isDarkMode={isDarkMode}
              closeErrorMessage={closeErrorMessage}
              minimizeErrorMessage={minimizeErrorMessage}
              restoreErrorMessage={restoreErrorMessage}
              minimized={minimized}
              iconType="error" // O "info", dependiendo de tu caso
            />
          </div>
          <button id="logoutButton" className="logout-btn" onClick={logoutAndRedirect}>Salir</button>
        </div>
      </div>
    );
  }
};
export default Chat;



 