// src/server.ts
import express, { Request, Response } from 'express';
const app = express();
const port = 8080;

// Middleware para parsear JSON
app.use(express.json());

// Ruta de ejemplo
app.get('/', (req: Request, res: Response) => {
  res.send('¡Bienvenido a GoChat!');
});

// Ruta para obtener los mensajes (por ejemplo, en una API REST)
app.get('/api/messages', (req: Request, res: Response) => {
  res.json([
    { user: 'Usuario1', message: 'Hola, ¿cómo estás?' },
    { user: 'Usuario2', message: '¡Hola! Estoy bien, gracias.' },
  ]);
});

// Iniciar el servidor
app.listen(port, () => {
  console.log(`Servidor escuchando en http://localhost:${port}`);
});
