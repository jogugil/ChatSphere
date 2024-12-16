// src/utils/WebSocketManager.ts
export class WebSocketManager {
  private socket: WebSocket | null = null;
  private url: string;
  private messageQueue: string[] = [];
  private _isConnected: boolean = false;
  private retryCount: number = 0; // Contador de intentos de reconexión
  private maxRetries: number;
  private socketId: string;

  constructor(url: string, socketId: string) {
    this.url = url;
    this.socketId = socketId; // Se usa un ID para diferenciar las conexiones
    this.maxRetries = parseInt(import.meta.env.VITE_RECONNECT_WEBSOCKET, 10) || 3;
    this.connect();
  }

  private connect() {
    console.log(`Conectando a WebSocket para ${this.socketId}...`);
    this.socket = new WebSocket(this.url);

    this.socket.onopen = () => {
      console.log(`WebSocket conectado para ${this.socketId}`);
      this.isConnected = true;
      this.retryCount = 0;

      // Enviar mensajes en cola
      while (this.messageQueue.length > 0) {
        const message = this.messageQueue.shift();
        if (message) this.socket?.send(message);
      }
    };

    this.socket.onmessage = (event) => {
      console.log(`Mensaje recibido para ${this.socketId}:`, event.data);
    };

    this.socket.onclose = () => {
      console.log(`WebSocket desconectado para ${this.socketId}, intentando reconectar...`);
      this.isConnected = false;

      if (this.retryCount < this.maxRetries) {
        setTimeout(() => this.connect(), 2000);
        this.retryCount += 1;
      } else {
        throw new Error('No se pudo conectar al WebSocket después de varios intentos');
      }
    };

    this.socket.onerror = (error) => {
      console.error(`Error en WebSocket para ${this.socketId}:`, error);
    };
  }

  public sendMessage(message: string) {
    if (this.isConnected) {
      this.socket?.send(message);
    } else {
      this.messageQueue.push(message);
    }
  }

  public close() {
    if (this.socket) {
      this.socket.close();
    }
  }

  public get isConnected(): boolean {
    return this._isConnected;
  }

  public set isConnected(value: boolean) {
    this._isConnected = value;
  }

  public getSocket() {
    return this.socket;
  }
}

export default WebSocketManager;
