// Obtener los elementos del DOM
const messageInput = document.getElementById('messageInput') as HTMLInputElement;
const sendMessageButton = document.getElementById('sendMessageButton') as HTMLButtonElement;
const messagesZone = document.getElementById('messagesZone') as HTMLElement;
const usersZone = document.getElementById('usersZone') as HTMLElement;

// Recuperar los datos de la sesión desde el localStorage
const token = localStorage.getItem('token');
const nickname = localStorage.getItem('nickname');
const roomId = localStorage.getItem('roomId');
const roomName = localStorage.getItem('roomName');

// Verificar si la sesión es válida
if (!token || !nickname || !roomId || !roomName) {
  window.location.href = 'pantallaInicio.html';  // Redirigir si no hay sesión activa
}

// Crear una instancia de ClienteChat
const clienteChat = new ClienteChat(nickname!, token!, roomId!, roomName!);

// Función para actualizar los mensajes y usuarios activos
const actualizarChat = async () => {
  const mensajes = await clienteChat.obtenerMensajes();
  const usuarios = await clienteChat.obtenerUsuarios();

  // Mostrar mensajes en la zona de mensajes
  messagesZone.innerHTML = '';
  mensajes.forEach((msg: { nickname: string, message: string }) => {
    const messageElement = document.createElement('div');
    messageElement.classList.add('message');
    messageElement.innerHTML = `<strong>${msg.nickname}</strong>: ${msg.message}`;
    messagesZone.appendChild(messageElement);
  });

  // Mostrar usuarios en la zona de usuarios
  usersZone.innerHTML = '<h3>Usuarios Activos</h3>';
  usuarios.forEach((user: string) => {
    const userElement = document.createElement('div');
    userElement.classList.add('user');
    userElement.innerText = user;
    usersZone.appendChild(userElement);
  });
};

// Iniciar polling para actualizar mensajes y usuarios
setInterval(actualizarChat, 5000);  // Cada 5 segundos

// Función para enviar un mensaje
const sendMessage = async () => {
  const message = messageInput.value.trim();
  if (message) {
    await clienteChat.enviarMensaje(message);
    messageInput.value = '';  // Limpiar el campo de entrada después de enviar
  }
};

// Event listener para el botón de enviar mensaje
sendMessageButton.addEventListener('click', sendMessage);

// Cargar datos al cargar la página
window.onload = actualizarChat;
