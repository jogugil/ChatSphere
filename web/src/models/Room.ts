import { validate as validateUUID } from 'uuid';
import { RoomChat } from '../types/types'; // Ajusta la ruta según la estructura de tu proyecto
import { MessageResponse, MessageList } from '../types/typesComm'; // Ajusta la ruta según la estructura de tu proyecto
import { ResponseUser, UsersAlive } from '../types/typesComm'; // Ajusta la ruta según la estructura de tu proyecto
import { Message } from '../models/Message'; // Ajusta la ruta según la estructura de tu proyecto

export type UUID = string;

// Clase Room que implementa la interfaz RoomChat y valida el idRoom
export class Room implements RoomChat {
  _roomid: UUID = "00000000-0000-0000-0000-000000000000"; 
  nombre: string;
  usuarios: string[];
  messages: Message[] = [];
  private lastIdMessage: UUID | null = null;

  constructor(roomId: string, nombre: string) {
    this.roomId = roomId; // Valida el ID de la sala al inicializar
    this.nombre = nombre;
    this.usuarios = []; // Inicializa la lista de usuarios vacía
  }

  // Getter y Setter para idRoom
  get roomId(): UUID {
    return this._roomid;
  }

  set roomId(value: string) {
    // Utilizando validateUUID del paquete uuid
    if (!validateUUID(value)) {
      throw new Error('El ID de la sala no es un UUID válido');
    }
    this._roomid = value; // Si la validación pasa, asigna el valor
  }

  // Getter y Setter para nombre
  get roomName(): string {
    return this.nombre;
  }

  set roomName(value: string) {
    this.nombre = value;
  }

  // Getter y Setter para usuarios
  get userList(): string[] {
    return this.usuarios;
  }

  set userList(usuarios: string[]) {
    this.usuarios = usuarios;
  }

  // Getter y Setter para messages
  get messageList(): Message[] {
    return this.messages;
  }

  set messageList(messages: Message[]) {
    this.messages = messages;
  }

  // Getter y Setter para lastIdMessage
  get lastIDMessageId(): UUID | null {
    return this.lastIdMessage;
  }

  set lastIDMessageId(id: UUID | null) {
    this.lastIdMessage = id;
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
  // Nueva función para procesar el JSON y añadir usuarios
  addUsersFromJson(json: string): void {
    try {
      const response: ResponseUser = JSON.parse(json);

      // Validamos que la sala del JSON coincida con la sala de la instancia
      if (response.roomid !== this.roomId) {
        console.log("La sala en el JSON no coincide con la sala de la instancia.");
        return;
      }

      // Añadir los usuarios activos a la lista de usuarios de la sala
      response.usersalive.forEach(usuario => {
        this.addUser(usuario.nickname); // Añade cada usuario a la sala
      });

      console.log(`${response.usersalive.length} usuarios añadidos a la sala.`);
    } catch (error) {
      console.error("Error al procesar el JSON:", error);
    }
  }
  // Método para actualizar los mensajes de la sala con los datos de la respuesta
  updateMessages(responseMessage: string) {
      // Parseamos el string JSON a un objeto MessageResponse
      let response: MessageResponse;

      try {
          response = JSON.parse(responseMessage);
      } catch (error) {
          console.log("Error al parsear el JSON:", error);
          return;
      }

      // Verificamos si el roomId coincide
      if (response.roomId !== this.roomid) {
          console.log("El roomId no coincide.");
          return;
      }

      // Si hay mensajes, actualizamos la lista de mensajes y el último id de mensaje
      if (response.ListMessage.length > 0) {
          // Añadimos todos los mensajes a la lista de mensajes
          response.ListMessage.forEach((msg: MessageList) => {
              const message = new Message(
                  msg.nickName,
                  msg.message,
                  new Date().toISOString(), // Asignamos la fecha actual como timestamp
                  msg.idMessage,
                  this.roomId
              );
              this.messages.push(message); // Guardamos todos los mensajes
          });

          // Extraemos el último mensaje de la lista y actualizamos el id
          const lastMessage = response.ListMessage[response.ListMessage.length - 1];
          this.lastIDMessageId = lastMessage.idMessage; // Actualizamos solo el último id de mensaje
      }
  }

  // Método para obtener el último mensaje
  getLastMessage() {
    if (this.lastIdMessage) {
      return this.messages.find(msg => msg.idMessage === this.lastIdMessage);
    }
    return null;
  }

  // Método para obtener todos los usuarios de la sala
  getUsers(): string[] {
    return this.usuarios;
  }
}
