import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {  getMessages, sendMessage, getActiveUsers } from '../api/api';
import { getToken } from '../utils/jwtUtils';
import usePolling from '../hooks/usePolling';
import  {Room} from '../models/Room'; // Ajusta la ruta según la estructura de tu proyecto
import  {Message, UUID} from "../models/Message"
import  {User} from "../models/User"
import { MessageResponse } from '../types/typesComm';
import { useAuth } from './AuthContext';

const Chat = () => {
    const [message, setMessage] = useState('');
    const [messages, setMessages] = useState([]);
    const [activeUsers, setActiveUsers] = useState<User[]>([]);
    const [room, setRoom] = useState<Room | null>(null);
    const [currentTime, setCurrentTime] = useState<string>(new Date().toLocaleString());
    const navigate = useNavigate();
  
     
    const { token, nickname, roomId, roomName } = useAuth();
    
    // Esta función asegura que room esté inicializado antes de ser utilizado
    useEffect(
      () => { 
        if (roomId && roomName && token && nickname) { 
          const roomObj = new Room(roomId, roomName); 
          const userObj = new User(nickname, 'someUserId', 'active', roomId, token); 
          setRoom(roomObj); 
          setUser(userObj); 
          // Actualización cada minuto de la hora 
          const interval = setInterval(() => { setCurrentTime(new Date().toLocaleString()); }, 60000); 
          // Limpiar el intervalo cuando el componente se desmonte 
          return () => clearInterval(interval); } 
      }, [roomId, roomName, token, nickname]);

    if (!token) {
        navigate('/'); // Redirige al login si no hay token
        return null;
    }
    

    // Función para actualizar la lista de mensajes con nuevos mensajes
    const updateMessagesFromPolling = (newListMessages: string) => {
        if (room && room.updateMessages) {
            room.updateMessages(newListMessages); // Llamar al método updateMessages de la clase Room
        }
    };
    const updateUsersFromPolling = (newListMessages: string) => {
      if (room && room.addUsersFromJson) {
          room.addUsersFromJson(newListMessages); // Llamar al método updateMessages de la clase Room
      }
  };
    
    // Obtiene los mensajes y usuarios activos
    // Hook para realizar el polling cada 5 segundos
    // Hook para obtener los mensajes, solo si room no es null
    usePolling(
      () => {
        if (room && room.lastIDMessageId) {
          // Llama a la función para obtener mensajes de la sala
          return getMessages(token, room.idRoom, room.lastIDMessageId);
        }
        return null; // Devuelve null si room es null
      },
      updateMessagesFromPolling,  // Esta función es el callback para actualizar los mensajes
      5000 // Intervalo en milisegundos (5 segundos)
    );   
    
    // Hook para obtener los usuarios activos, solo si roomId está disponible
    usePolling(
      () => {
        if (room) {
          return getActiveUsers(token, room._idRoom); // Usa el id de la sala para obtener los usuarios activos
        }
        return null; // Devuelve null si room es null
      },
      updateUsersFromPolling,  // Esta función es el callback para actualizar los mensajes
      5000 // Intervalo en milisegundos (5 segundos)
    )  
    
  
    const handleSendMessage = async () => {
      if (message.trim() === '') return;
      try {
        await sendMessage(token, roomId, message);
        setMessage(''); // Limpiar el campo de mensaje
      } catch (error) {
        console.error('Error al enviar el mensaje', error);
      }
    };
  
    const handleLogout = () => {
      localStorage.removeItem('token');
      navigate('/'); // Redirige al login
    };
  
    const formatMessage = (msg: { timestamp: string; nickname: string; message: string }) => {
      const date = new Date(msg.timestamp);
      return `${date.toLocaleDateString()} ${date.toLocaleTimeString()} - ${msg.nickname}: ${msg.message}`;
    };
      


 
    return (
      <div className="chat-container">
        <div className="left-sidebar">
          <h2>GoChat ZeroMq</h2>
          <p>{currentTime}</p>
          <div>
            <strong>{room?.nombre}</strong>
          </div>
          <div className="stats">
            <p>{activeUsers.length} usuarios activos</p>
            <p>{messages.length} mensajes en la última hora</p>
          </div>
          <button onClick={handleLogout}>Salir</button>
          <div className="banners">
            <img src="banner1.png" alt="Banner 1" />
            <img src="banner2.png" alt="Banner 2" />
          </div>
          <footer>© SmartIAServices, {new Date().getFullYear()}</footer>
        </div>
  
        <div className="chat-room">
          <div className="messages-display">
            <h3>Mensajes</h3>
            <ul>
              {messages.map((msg, index) => (
                <li key={index}>
                  {formatMessage(msg)}
                </li>
              ))}
            </ul>
          </div>
          <div className="input-section">
            <input
              type="text"
              placeholder="Escribe tu mensaje..."
              value={message}
              onChange={(e) => setMessage(e.target.value)}
            />
            <button onClick={handleSendMessage}>Enviar</button>
          </div>
        </div>
  
        <div className="right-sidebar">
          <h3>Usuarios en la sala</h3>
          <ul>
            {activeUsers.map((user, index) => (
              <li key={index}>{user.nickname}</li>
            ))}
          </ul>
        </div>
      </div>
    );
  };
  
  export default Chat;

function setUser(userObj: User) {
  throw new Error('Function not implemented.');
}


function setLoading(arg0: boolean) {
  throw new Error('Function not implemented.');
}
