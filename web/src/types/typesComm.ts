import { UUID } from "../models/User";

export const TOKEN_NULO = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIwIiwidXNlcm5hbWUiOiJkZWZhdWx0IiwiZXhwIjoxNjk4Mzg2Mjc5fQ.MGkJd5EYHQffQ9jrUzX7Djgmd4mOuH3aPRvcOP61TnM";
// Estructura de la respuesta
export interface MessageResponse {
    status: string;
    messages: string;
    tokenSesion: string;
    nickname: string;
    idSala: string;
    ListMessage: MessageList[];
}

export interface MessageList {
    idMessage: string;
    nickName: string;
    message: string;
}

// Estructuras que defines para usuarios y respuestas
export interface  UsersAlive {
    nickname: string;
    acttionLastDate: string;
  };
  
export interface ResponseUser  {
    status: string;
    message: string;
    tokenSesion: string;
    nickname: string;
    roomId: string;
    usersAlive: UsersAlive[];
  };


export interface LoginResponse {
    status: string;
    message: string;
    token: string;
    nickname: string;
    idsala: UUID;
    namesala: string;
  }

export interface JwtPayload {
    userId: string;
    username: string;
    exp?: number;  // La expiración es opcional
  }