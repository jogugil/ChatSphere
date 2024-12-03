import React, { useState, useEffect } from 'react';
import { sendMessage, getMessages, getUsers } from '../../shared/communication/chatService';

const Chat: React.FC = () => {
  const [nickname, setNickname] = useState('');
  const [messages, setMessages] = useState<any[]>([]);
  const [users, setUsers] = useState<any[]>([]);
  const [newMessage, setNewMessage] = useState('');

  useEffect(() => {
    // Cargar mensajes y usuarios al iniciar
    const fetchData = async () => {
      const msgs = await getMessages();
      const usr = await getUsers();
      setMessages(msgs);
      setUsers(usr);
    };

    fetchData();
  }, []);

  const handleSendMessage = async () => {
    // Suponiendo que ya tienes un token JWT y roomId
    const token = 'generated_token';
    const roomId = 'room_1';
    await sendMessage(nickname, newMessage, token, roomId);
    setNewMessage('');
  };

  return (
    <div>
      <h1>Chat Room</h1>
      <div>
        <input
          type="text"
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          placeholder="Type your message..."
        />
        <button onClick={handleSendMessage}>Send</button>
      </div>
      <div>
        <h2>Messages</h2>
        {messages.map((msg, idx) => (
          <div key={idx}>
            <strong>{msg.nickname}</strong>: {msg.message}
          </div>
        ))}
      </div>
      <div>
        <h2>Active Users</h2>
        {users.map((user, idx) => (
          <div key={idx}>{user.nickname}</div>
        ))}
      </div>
    </div>
  );
};

export default Chat;
import React, { useState, useEffect } from 'react';
import { sendMessage, getMessages, getUsers } from '../../shared/communication/chatService';

const Chat: React.FC = () => {
  const [nickname, setNickname] = useState('');
  const [messages, setMessages] = useState<any[]>([]);
  const [users, setUsers] = useState<any[]>([]);
  const [newMessage, setNewMessage] = useState('');

  useEffect(() => {
    // Cargar mensajes y usuarios al iniciar
    const fetchData = async () => {
      const msgs = await getMessages();
      const usr = await getUsers();
      setMessages(msgs);
      setUsers(usr);
    };

    fetchData();
  }, []);

  const handleSendMessage = async () => {
    // Suponiendo que ya tienes un token JWT y roomId
    const token = 'generated_token';
    const roomId = 'room_1';
    await sendMessage(nickname, newMessage, token, roomId);
    setNewMessage('');
  };

  return (
    <div>
      <h1>Chat Room</h1>
      <div>
        <input
          type="text"
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          placeholder="Type your message..."
        />
        <button onClick={handleSendMessage}>Send</button>
      </div>
      <div>
        <h2>Messages</h2>
        {messages.map((msg, idx) => (
          <div key={idx}>
            <strong>{msg.nickname}</strong>: {msg.message}
          </div>
        ))}
      </div>
      <div>
        <h2>Active Users</h2>
        {users.map((user, idx) => (
          <div key={idx}>{user.nickname}</div>
        ))}
      </div>
    </div>
  );
};

export default Chat;
