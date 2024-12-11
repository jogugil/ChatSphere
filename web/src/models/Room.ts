import { validate as validateUUID } from 'uuid';
import { RoomChat} from '../types/types'; // Ajusta la ruta según la estructura de tu proyecto

export type UUID = string;
// Clase RoomChat que implementa la interfaz Room y valida el idRoom
export class Room implements RoomChat {
    _idRoom: UUID = "00000000-0000-0000-0000-000000000000"; 
    nombre: string;
    usuarios: string[];
  
    constructor(idRoom: string, nombre: string) {
      this.idRoom   = idRoom ; // Valida el ID de la sala al inicializar
      this.nombre   = nombre;
      this.usuarios = []; // Inicializa la lista de usuarios vacía
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
  
    // Método para añadir un usuario a la sala
    addUser(nickname: string): void {
      if (!this.usuarios.includes(nickname)) {
        this.usuarios.push(nickname);
      } else {
        console.log('El usuario ya está en la sala');
      }
    }
  
    // Método para eliminar un usuario de la sala
    removeUser(nickname: string): void {
      const index = this.usuarios.indexOf(nickname);
      if (index > -1) {
        this.usuarios.splice(index, 1);
      } else {
        console.log('El usuario no está en la sala');
      }
    }
  
    // Método para obtener todos los usuarios de la sala
    getUsers(): string[] {
      return this.usuarios;
    }
  }