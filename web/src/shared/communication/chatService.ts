// chatService.ts
const API_BASE_URL = 'http://localhost:3000';  // Asegúrate de que sea la URL correcta de tu API

// Función para obtener los mensajes
export const getMessages = async (): Promise<any[]> => {
  try {
    const response = await fetch(`${API_BASE_URL}/messages`);
    if (!response.ok) {
      throw new Error('Failed to fetch messages');
    }
    return await response.json();  // Asumimos que la respuesta es un array de mensajes
  } catch (error) {
    console.error('Error fetching messages:', error);
    return [];
  }
};

// Función para obtener los usuarios
export const getUsers = async (): Promise<any[]> => {
  try {
    const response = await fetch(`${API_BASE_URL}/users`);
    if (!response.ok) {
      throw new Error('Failed to fetch users');
    }
    return await response.json();  // Asumimos que la respuesta es un array de usuarios
  } catch (error) {
    console.error('Error fetching users:', error);
    return [];
  }
};

// Función para enviar un mensaje
export const sendMessage = async (nickname: string, message: string, token: string, roomId: string): Promise<void> => {
  try {
    const response = await fetch(`${API_BASE_URL}/sendMessage`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify({
        nickname,
        message,
        roomId,
      }),
    });
    if (!response.ok) {
      throw new Error('Failed to send message');
    }
    console.log('Message sent successfully');
  } catch (error) {
    console.error('Error sending message:', error);
  }
};
