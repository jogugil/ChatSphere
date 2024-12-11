// Estructura de la respuesta
interface MessageResponse {
    status: string;
    messages: string;
    tokenSesion: string;
    nickname: string;
    idSala: string;
    ListMessage: MessageList[];
  }
  interface MessageList {
    idMensaje: string;
    nickname: string;
    texto: string;
  }