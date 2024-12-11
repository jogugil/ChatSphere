// user.ts

import jwt, { JwtPayload } from 'jsonwebtoken'; // Utilizanmos jsonwebtoken para manejar JWT
import { UserChat } from '../types/types';      
import { validate as validateUUID } from 'uuid';

// Definir la interfaz del payload del token
const payload = {
    userId: "0", // ID del usuario en cero
    username: "default", // Nombre de usuario en cero
    exp: Math.floor(Date.now() / 1000) + 60 * 60 *24, // Expiración en 1 día (el tiempo actual + 24 horas)
  };
export type UUID = string;
export class User {
  nickname: string;
  id: string;
  status: string;
  private _token: string = jwt.sign(payload, "nulo");;
  private _idRoom: UUID = "00000000-0000-0000-0000-000000000000"; 

  constructor(nickname: string, id: string, status: string, idRoom: string, token: string) {
    this.nickname = nickname;
    this.id = id;
    this.status = status;
    this.idRoom = idRoom; // Esto llamará al setter que valida el UUID
    if (token) {        
        this.token = jwt.sign(payload, "nulo" );
    } else  {
        this.token = token;
    }
  }
  // Getter para _token
  get token(): string {
    return this._token;
  }

  // Setter para _token con validación
  set token(value: string) {
    if (this.isValidToken(value)) {
      this._token = value; // Solo asigna si es un token válido
    } else {
      throw new Error('El token no es válido');
    }
  }
  // Método para validar si el token es válido
  private isValidToken(token: string): boolean {
    try {
      const decoded = jwt.decode(token) as JwtPayload;
      return decoded != null && this.tokenIsValid(decoded.exp);
    } catch (error) {
      console.error('Error al validar el token:', error);
      return false;
    }
  }

    // Verifica si el token ha expirado
    private tokenIsValid(expiration: number | undefined): boolean {
        if (expiration === undefined) {
        return false; // Si no hay fecha de expiración, no es válido
        }

        const now = Math.floor(Date.now() / 1000); // Obtenemos la fecha y hora actual en segundos
        const expirationDatePlusOneDay = expiration + 60 * 60 * 24; // Añadir 24 horas (un día)

        return now < expirationDatePlusOneDay; // Verificar si la fecha actual es antes de la fecha de expiración + 1 día
    }

  // Método para obtener el ID del usuario decodificado desde el token
  getUserIdFromToken(): string | null {
    try {
      const decoded = jwt.decode(this._token) as JwtPayload;
      return decoded ? decoded.userId : null;
    } catch (error) {
      console.error('Error al obtener el ID del token:', error);
      return null;
    }
  } 

    // Getter y Setter para idSala
    get idRoom(): UUID {
        return this._idRoom;
    }
    set idRoom(value: string) {
        if (!validateUUID(value)) {
            throw new Error('El ID de la sala no es un UUID válido');
        }
        this._idRoom = value; // Si la validación pasa, asigna el valor
    }
}