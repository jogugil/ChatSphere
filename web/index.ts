import express, { Request, Response } from 'express';
import bodyParser from 'body-parser';
import { chatRouter } from './routes/routes';

const app = express();
const port = 8080;

// Middleware para manejar JSON
app.use(bodyParser.json());

// Rutas
app.use('/api', chatRouter);

// Iniciar servidor
app.listen(port, () => {
  console.log(`Servidor corriendo en http://localhost:${port}`);
});
