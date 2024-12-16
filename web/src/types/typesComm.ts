import { UUID } from "../models/User";

export const TOKEN_NULO = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIwIiwidXNlcm5hbWUiOiJkZWZhdWx0IiwiZXhwIjoxNjk4Mzg2Mjc5fQ.MGkJd5EYHQffQ9jrUzX7Djgmd4mOuH3aPRvcOP61TnM";
// Estructura de la respuesta
export interface MessageResponse {
    status: string;
    messages: string;
    tokenSesion: string;
    nickname: string;
    roomId: string;
    ListMessage: MessageList[];
    x_gochat: string;
}

export interface MessageList {
    idMessage: string;
    nickName: string;
    message: string;
}

// Estructuras que defines para usuarios y respuestas
export interface  UsersAlive {
    nickname: string;
    acttionlastdate: string;
  };
  
export interface ResponseUser  {
    status: string;
    message: string;
    tokenSesion: string;
    nickname: string;
    roomid: string;
    usersalive: UsersAlive[];
  };


export interface LoginResponse {
    status: string;
    message: string;
    token: string;
    nickname: string;
    roomid: UUID;
    roomname: string;
  }

export interface JwtPayload {
    userid: string;
    username: string;
    exp?: number;  // La expiración es opcional
  }