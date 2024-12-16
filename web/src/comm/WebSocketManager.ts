import { getSystemErrorName } from "util";

export class WebSocketManager {
  private socket: WebSocket | null = null;
  private url: string;
  private messageQueue: string[] = [];
  private _isConnected: boolean = false;
  private retryCount: number = 0;
  private maxRetries: number;
  private socketId: string;
  private onErrorCallback: ((error: string) => void) | null = null;
 
   // Agregar una propiedad para el callback
   onOpenCallback?: () => void;

  constructor(url: string, socketId: string, onErrorCallback?: (error: string) => void) {
    this.url = url;
    this.socketId = socketId;
    this.maxRetries = parseInt(import.meta.env.VITE_RECONNECT_WEBSOCKET, 10) || 3;
    if (onErrorCallback) {
      this.onErrorCallback = onErrorCallback;  // Guardamos el callback de error
    }
    this.connect();
  }

  private connect() {
    console.log(`Conectando a WebSocket para ${this.socketId}...`);
    this.socket = new WebSocket(this.url);

    this.socket.onopen = () => {
      console.log(`WebSocket conectado para ${this.socketId}`);
      this.isConnected = true;
      this.retryCount = 0;
      console.log('Cola de mensajes en onopen:', this.messageQueue);
      while (this.messageQueue.length > 0) {
        const message = this.messageQueue.shift();
        if (message) this.socket?.send(message);
      }

      // Ejecutar el callback si está definido
      if (this.onOpenCallback) {
        this.onOpenCallback();
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
        if (this.onErrorCallback) {
          this.onErrorCallback(`No se pudo conectar al WebSocket después de ${this.maxRetries} intentos.`);
        }
      }
    };

    this.socket.onerror = (error) => {
      console.error(`Error en WebSocket para ${this.socketId}:`, error);
      this.onErrorCallback && this.onErrorCallback("Error en WebSocket: " + (error instanceof Error ? error.message : "Error desconocido"));
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

  // Método para asignar el callback onOpen desde fuera de la clase
  setOnOpenCallback(callback: () => void) {
    this.onOpenCallback = callback;
  }
}