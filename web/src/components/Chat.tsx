import React, { useState, useEffect, useContext } from 'react';
import { useNavigate } from 'react-router-dom';
import { getMessages, sendMessage, getAliveUsers } from '../api/api';
import { useAuth } from './AuthContext';  // Importamos el contexto de autenticación
import { Room } from '../models/Room';
import { Message, UUID } from "../models/Message";
import { User } from "../models/User";
import { MessageResponse } from '../types/typesComm';
import { WErrorMessage } from "./ErrorMessage"; // Componente para mostrar errores
import {ResponseUser, UsersAlive} from '../types/typesComm'
import '../styles/chat.css';
import BannerProgramming from './BannerProgramming';
import BannerCloud from './BannerCloud';


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
    IdSala: string;
    AliveUsers: aliveUsers[];
  };

  const [messages, setMessages]         = useState<Message[]> ([]);        // Estado para los mensajes
  const [messageText, setMessageText]   = useState<string> ('');           // Estado para el texto del mensaje
  const [aliveUsers, setAliveUsers]     = useState<string[]> ([]);         // Estado para los usuarios activos
  const [error, setError]               = useState<string | null> (null);  // Estado para errores
  
  const [errorMessage, setErrorMessage] = useState<string | null> (null);
  const [showError, setShowError]       = useState (false);
  const [minimized, setMinimized]       = useState (false);

  const navigate = useNavigate();
  
  const closeErrorMessage    = () => setShowError (false);  // Función para cerrar el mensaje de error
  const minimizeErrorMessage = () => setMinimized (true);   // Función para minimizar el mensaje de error
  const restoreErrorMessage  = () => setMinimized (false);  // Función para restaurar el mensaje de error

  const [isDarkMode, setIsDarkMode] = useState(false);
  const toggleTheme = () => {
    setIsDarkMode(!isDarkMode);
  };
  
  const { token, nickName, roomId, roomName } = useAuth();  // Obtener el usuario y el token del contexto
  console.log ("Caht.tsx desoues de Authcontx:", token, nickName, roomId, roomName);
  const [userChat, setUserChat] = useState<User | null>(null);  // Estado para el objeto User
  useEffect(() => {
    console.log ("Caht.tsx useEffect :", token, nickName, roomId, roomName);
    if (nickName && roomId && roomName && token) {
      const user = new User(nickName, 'Alive', roomId, roomName, token);
      setUserChat(user); // Asignar el usuario a un estado
      console.log('Usuario autenticado:', userChat);
    } else {
      setError('Datos del usuario no válidos');
      setErrorMessage('Datos del usuario no válidos');
      setShowError(true);
      navigate('/Login');  // Redirigir al login si no hay datos del usuario
    }
  }, [userChat, navigate]);


  console.log('Usuario autenticado:', userChat);
  console.log("userChat:", userChat);
  console.log("userChat.token:", userChat?.token);
  console.log("userChat.roomId:", userChat?.roomId);
  
  const [room, setRoom] = useState<Room | null>(null);  // Estado para el objeto Room
  useEffect(() => {
    if (nickName && roomId && roomName && token) {
      const roomU = new Room(roomId, roomName);
      setRoom(roomU); // Asignar el usuario a un estado
    } else {
      setError('Datos del usuario no válidos');
      setErrorMessage('Datos del usuario no válidos');
      setShowError(true);
      navigate('/Login');  // Redirigir al login si no hay datos del usuario
    }
  }, [room, navigate]);
  
  console.log('Sala creada:' );
  console.log("room:", room);
  console.log("room.roomId:", room?.roomId);
  console.log("room.roomName:", room?.roomName);

//Función que proicesa el mensaje JSON del servidor GoChat. lista de usaurios activos
const loadAliveUsers = async (userObject: User ) => {
  try {
    const response = await getAliveUsers(userObject.token, userObject.roomId);
    
    // Parsear la respuesta JSON
    const data: ResponseUser = JSON.parse(response);  // Asegúrate de que la respuesta es un JSON
    if (data.Status === 'success' && data.AliveUsers) {
      // Extraer los nicknames de los usuarios activos
      const nicknames = data.AliveUsers.map(user => user.Nickname);
      setAliveUsers(nicknames);  // Establecer el estado con los nicknames
      console.log('Usuarios activos:', nicknames);
    } else {
      setError('Error al obtener los usuarios activos');
      console.error('Error al obtener usuarios activos:', data.Message);
    }
  } catch (err) {
    setError('Error al obtener los usuarios activos: ' + err);
    console.error('Error al obtener usuarios activos:', err);
  }
};

// Verificar si el usuario está autenticado
useEffect(() => {
    // Obtener mensajes históricos al cargar la página
    const loadMessages = async (userObject: any) => {
    // Generar un UUID vacío para el primer request
    const emptyUUID: UUID = "00000000-0000-0000-0000-000000000000";
      
    try {
        // Hacer la primera llamada a la API para obtener los mensajes
        const messageList = await getMessages(userObject.token, userObject.idRoom, emptyUUID);
        
        if (!room) {
          // Mostrar una ventana emergente de error
          alert('El servicio de chat no está disponible. Disculpe las molestias. Por favor, intente ingresar nuevamente más tarde.');
  
          // Redirigir al login
          navigate('/Login');
          return;  // Detener la ejecución del código si room es null
        }
  
        // Convertir los mensajes a objetos
        room.updateMessages(messageList);
  
        // Actualizar el estado de los mensajes
        setMessages(room.messageList);
  
        console.log('Mensajes cargados:', room.messageList);
      } catch (err) {
        setError('El Servicio de Chat está temporalmente cerrado. Intente logarse mñas tarde');
        console.error('Error al cargar los mensajes:', err);
        setShowError(true);
      }

      // Verificar que el objeto del usuario esté disponible
      if (!userChat || !userChat.token || !userChat.roomId) {
        setErrorMessage('Datos del usuario no válidos');
        setShowError(true);
        navigate('/Login');
      } else {
        // Cargar los usuarios activos
        console.log ("Cargamos los usuarios activos:",userChat);
        loadAliveUsers(userChat); 
      }
  };
}, [messages, aliveUsers, navigate]);

// Función para enviar un mensaje
const handleSendMessage = async () => {
    console.log ("handleSendMessage: messageText: ",messageText )
    console.log ("handleSendMessage: userChat: ",userChat )

    try {
      if (!userChat || !userChat.token || !userChat.roomId) {
        setError('Datos del usuario no válidos');
        setErrorMessage('Datos del usuario no válidos');
        setShowError(true);
        navigate('/login');  // Redirigir al login si no hay datos del usuario
      } else {
        const newMessage = new Message(userChat?.nickname || '', messageText, new Date().toISOString(), 'newId', userChat?.roomId || '');
        console.log('Enviando mensaje:', newMessage);     
        await sendMessage(userChat.nickname, userChat.token, userChat.roomId, userChat.roomName, messageText);
        setMessages([...messages, newMessage]);
        setMessageText('');  // Limpiar el campo de mensaje después de enviarlo
        console.log('Mensaje enviado con éxito:', newMessage);
        if (messageText.trim() === '') return;
      }
    } catch (error) {
      setErrorMessage(`Error al enviar el mensaje`);
      setShowError(true);
    }
  };

  // Función para manejar el cambio en el campo de mensaje
  const handleMessageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setMessageText(e.target.value);
    console.log('Cambio en el mensaje:', e.target.value);
  };
  const logoutAndRedirect = (event: React.MouseEvent<HTMLButtonElement>) => {
    // Lógica para destruir los objetos Room, User, Messages
    localStorage.removeItem("Room");
    localStorage.removeItem("User");
    localStorage.removeItem("Messages");
  
    // Redirigir al login
    window.location.href = "/Login"; // Redirige a la página de login
  };

  return (
    <div className="chat-container">
      <div className="chat-left-column">
          <div className="chat-title">ChatSphere</div>
          <div className="chat-subtitle">GoChat ZeroMQ</div>

           
            <div className="chat-logo-container">
              <div className="chat-logo">
                <img src="../../public/images/logo.webp" alt="Logo" />
              </div>
            </div>
     
            <div className="chat-container">
            <div className="chat-banners">
              {/* Banner de Programación en la Nube */}
              <div className="chat-banner-programming">
                <BannerProgramming 
                  title="Desarrolla, Concurre y Escala en la Nube"
                  subtitle="Impulsa tu aplicación con Go, TypeScript, Node.js, React y más"
                  imageUrl="https://via.placeholder.com/600x250?text=Banner+de+Tecnolog%C3%ADas+de+Programaci%C3%B3n"
                />
              </div>

              {/* Banner de Cloud Computing */}
              <div className="chat-banner-cloud">
                <BannerCloud 
                  title="Escala sin Servidores"
                  subtitle="Soluciones eficientes con FaaS y Serverless"
                  imageUrl="https://via.placeholder.com/600x250?text=Banner+de+Cloud+Computing"
                />
              </div>
            </div>
          </div>

          <div className="footer">
            <p>&copy; 2024 José Javier Gutiérrez Gil</p>
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
              <button className="send-btn" onClick={handleSendMessage}>Enviar</button>
              <input
                type="text"
                placeholder="Escribe tu mensaje..."
                value={messageText}
                onChange={handleMessageChange}
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
        <div className="error-message">
          <WErrorMessage
            message              = {errorMessage || ''} 
            showError            = {showError}
            closeErrorMessage    = {closeErrorMessage}
            minimizeErrorMessage = {minimizeErrorMessage}
            restoreErrorMessage  = {restoreErrorMessage}
            minimized            = {minimized}
            iconType             = {showError ? "error" : "info"} 
          />
        </div>
        <button id="logoutButton" className="logout-btn" onClick={logoutAndRedirect}>Salir</button>
      </div>
      
     </div>
  );
};

export default Chat;
