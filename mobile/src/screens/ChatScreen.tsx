import React, { useState, useEffect } from 'react';
import { View, TextInput, Button, Text } from 'react-native';
import { sendMessage, getMessages, getUsers } from '../../shared/communication/chatService';

const ChatScreen: React.FC = () => {
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
    <View>
      <TextInput
        value={newMessage}
        onChangeText={setNewMessage}
        placeholder="Type your message..."
      />
      <Button title="Send" onPress={handleSendMessage} />
      <Text>Messages</Text>
      {messages.map((msg, idx) => (
        <Text key={idx}>
          <strong>{msg.nickname}</strong>: {msg.message}
        </Text>
      ))}
      <Text>Active Users</Text>
      {users.map((user, idx) => (
        <Text key={idx}>{user.nickname}</Text>
      ))}
    </View>
  );
};

export default ChatScreen;
