import { Request, Response } from 'express';
import jwt from 'jsonwebtoken';

const users = new Set<string>();  // Simulamos una lista de usuarios registrados
const secretKey = 'your-secret-key';  // Clave secreta para firmar el token

// Función para generar un token JWT
const generateToken = (nickname: string) => {
  const payload = { nickname };
  const token = jwt.sign(payload, secretKey, { expiresIn: '1h' });  // Token válido por 1 hora
  return token;
};

export const login = (req: Request, res: Response) => {
  const { nickname } = req.body;

  if (users.has(nickname)) {
    // Si el nickname ya está en uso, devolver error 409
    return res.status(409).json({ message: 'El nickname ya está en uso. Intenta con otro.' });
  }

  // Si el nickname está disponible, agregarlo a la lista de usuarios y generar el token
  users.add(nickname);

  const token = generateToken(nickname);
  const roomName = 'Sala Principal';  // Nombre de la sala
  const roomId = 'room_123';  // ID de la sala (esto puede ser dinámico)

  return res.status(200).json({
    message: 'Nickname disponible. Bienvenido al chat.',
    token,
    roomName,
    roomId
  });
};
