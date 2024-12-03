class ClienteChat {
    // Propiedades
    private token: string;
    private nickname: string;
    private roomId: string;
    private roomName: string;
  
    constructor(nickname: string, token: string, roomId: string, roomName: string) {
      this.nickname = nickname;
      this.token = token;
      this.roomId = roomId;
      this.roomName = roomName;
    }
  
    // Método para crear la sesión del usuario
    crearSesion(): boolean {
      if (this.token && this.nickname && this.roomId) {
        // Guardar la sesión en el localStorage o en memoria
        localStorage.setItem('token', this.token);
        localStorage.setItem('nickname', this.nickname);
        localStorage.setItem('roomId', this.roomId);
        localStorage.setItem('roomName', this.roomName);
        return true;
      }
      return false;
    }
  
    // Método para conectar al chat
    conectarChat(): boolean {
      // Verificar que la sesión esté activa
      const token = localStorage.getItem('token');
      const nickname = localStorage.getItem('nickname');
      const roomId = localStorage.getItem('roomId');
      const roomName = localStorage.getItem('roomName');
  
      if (token && nickname && roomId && roomName) {
        // Si la sesión está activa, inicializamos la conexión
        this.token = token;
        this.nickname = nickname;
        this.roomId = roomId;
        this.roomName = roomName;
        return true;
      }
      return false;
    }
  
    // Método para obtener el historial de mensajes
    async obtenerMensajes() {
      try {
        const response = await fetch(`http://localhost:8080/api/getMessages?roomId=${this.roomId}`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${this.token}`,
          }
        });
  
        const data = await response.json();
        if (response.ok) {
          return data.messages;
        } else {
          console.error("Error obteniendo mensajes:", data.message);
          return [];
        }
      } catch (error) {
        console.error("Error al hacer la solicitud de mensajes:", error);
        return [];
      }
    }
  
    // Método para obtener la lista de usuarios activos
    async obtenerUsuarios() {
      try {
        const response = await fetch(`http://localhost:8080/api/getUsers?roomId=${this.roomId}`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${this.token}`,
          }
        });
  
        const data = await response.json();
        if (response.ok) {
          return data.users;
        } else {
          console.error("Error obteniendo usuarios:", data.message);
          return [];
        }
      } catch (error) {
        console.error("Error al hacer la solicitud de usuarios:", error);
        return [];
      }
    }
  
    // Método para enviar un nuevo mensaje
    async enviarMensaje(message: string) {
      try {
        const response = await fetch('http://localhost:8080/api/sendMessage', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${this.token}`,
          },
          body: JSON.stringify({
            nickname: this.nickname,
            message: message,
            roomId: this.roomId,
          })
        });
  
        const data = await response.json();
        if (response.ok) {
          console.log("Mensaje enviado con éxito:", data);
        } else {
          console.error("Error enviando mensaje:", data.message);
        }
      } catch (error) {
        console.error("Error al enviar mensaje:", error);
      }
    }
  }
  