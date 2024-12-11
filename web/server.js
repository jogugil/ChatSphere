 
// Reemplaza las importaciones por require
// Reemplaza las importaciones por require
const express = require('express');
const path = require('path');  // No olvides importar 'path'
const app = express();
const cors = require('cors');
// Define el puerto donde correrá el servidor
const PORT = process.env.PORT || 3000;
// Datos simulados
let messages = [
  { nickname: 'Alice', message: 'Hello world' },
  { nickname: 'Bob', message: 'Hi there!' }
];

let users = [
  { nickname: 'Alice' },
  { nickname: 'Bob' }
];


// Middleware para parsear JSON
app.use(express.json());
app.use(cors());

// Ruta para la página principal (home)
app.get('/tete', (req, res) => {
  // Envía un archivo HTML al navegador
  res.sendFile(path.join(__dirname, 'index.html'));
});
// Ruta para la página principal (home)
app.get('/', (req, res) => {
  // Envía un archivo HTML al navegador
  res.sendFile(path.join(__dirname, 'index.html'));
});
// Ruta adicional de ejemplo, por ejemplo para una página de chat
app.get('/chat', (req, res) => {
  res.send('Página de chat');
});
// Configura la ruta para servir archivos estáticos (CSS, JS, imágenes) desde la carpeta "public"
app.use(express.static(path.join(__dirname, 'public')));
 
// Inicia el servidor y escucha en el puerto
app.listen(PORT, () => {
  console.log(`Servidor corriendo en el puerto ${PORT}`);
});
// Endpoint para obtener mensajes
app.get('/messages', (req, res) => {
  res.json(messages);
});
// Endpoint para obtener usuarios
app.get('/users', (req, res) => {
  res.json(users);
});
// Endpoint para enviar mensaje
app.post('/sendMessage', (req, res) => {
  const { nickname, message, roomId } = req.body;
  if (!nickname || !message || !roomId) {
    return res.status(400).send('Missing parameters');
  }
  messages.push({ nickname, message });
  res.status(200).send('Message sent');
});
 

