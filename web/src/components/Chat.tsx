import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { getMessages, sendMessage, getActiveUsers } from '../api/api';
import { getToken } from '../utils/jwtUtils';
import usePolling from '../hooks/usePolling';
import  {Room} from '../models/Room'; // Ajusta la ruta según la estructura de tu proyecto
import  {Message} from "../models/Message"
import  {User} from "../models/User"

const Chat = () => {
    const [message, setMessage] = useState('');
    const [messages, setMessages] = useState([]);
    const [activeUsers, setActiveUsers] = useState<User[]>([]);
    const [room, setRoom] = useState<Room | null>(null);
    const [currentTime, setCurrentTime] = useState<string>(new Date().toLocaleString());
    const navigate = useNavigate();
  
    const token = getToken();
    const roomId = '12345'; // Este ID lo puedes obtener de la URL o de un contexto
  
    if (!token) {
      navigate('/'); // Redirige al login si no hay token
      return null;
    }
  
    // Obtiene los mensajes y usuarios activos
    // Hook para realizar el polling cada 5 segundos
    usePolling(() => getMessages(token, roomId, lastSeenMessageId), updateMessages, 5000);
    usePolling(() => getActiveUsers(token, roomId), (users: React.SetStateAction<User[]>) => setActiveUsers(users), 5000);
  
    // Obteniendo los detalles de la sala (nombre, número de usuarios, etc.)
    useEffect(() => {
      // Simulamos obtener los detalles de la sala (esto sería de una API)
      setRoom({
        idRoom: roomId,
        nombre: 'El Refugio Digital',
      });
      // Actualización cada minuto de la hora
      const interval = setInterval(() => setCurrentTime(new Date().toLocaleString()), 60000);
      return () => clearInterval(interval); // Limpiar el intervalo al desmontarse el componente
    }, []);
  
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
      
    // Función para actualizar la lista de mensajes con nuevos mensajes
    const updateMessages = (newMessages) => {
        // Puedes hacer un merge de los mensajes antiguos con los nuevos
        setMessages((prevMessages) => [...prevMessages, ...newMessages]);
    };
    // Función que realiza la petición para obtener los mensajes
    const getMessages = async (token, roomId, lastSeenMessageId) => {
        try {
        setLoading(true);
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

        const data = await response.json();
        return data.messages; // Asumiendo que los mensajes vienen en la propiedad 'messages'
        } catch (error) {
        console.error('Error fetching messages:', error);
        } finally {
        setLoading(false);
        }
    };

    // Función para actualizar la lista de mensajes con nuevos mensajes
    const updateMessages = (newMessages) => {
        // Puedes hacer un merge de los mensajes antiguos con los nuevos
        setMessages((prevMessages) => [...prevMessages, ...newMessages]);
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