// src/client.ts
const messageInput = document.getElementById('message-input') as HTMLInputElement;
const sendButton = document.getElementById('send-button') as HTMLButtonElement;
const messagesContainer = document.getElementById('messages') as HTMLElement;

sendButton.addEventListener('click', () => {
  const message = messageInput.value.trim();
  if (message) {
    // Aquí deberías enviar el mensaje al backend (usando Fetch o WebSocket)
    displayMessage('Tú', message);
    messageInput.value = ''; // Limpiar el input
  }
});

async function loadMessages() {
  const response = await fetch('/api/messages');
  const messages = await response.json();

  messages.forEach((msg: { user: string, message: string }) => {
    displayMessage(msg.user, msg.message);
  });
}

function displayMessage(user: string, message: string) {
  const messageElement = document.createElement('div');
  messageElement.innerHTML = `<strong>${user}:</strong> ${message}`;
  messagesContainer.appendChild(messageElement);
}

// Cargar los mensajes cuando la página cargue
loadMessages();




