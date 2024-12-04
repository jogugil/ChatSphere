package comm

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // Mapa de clientes conectados

// Maneja las nuevas conexiones WebSocket
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// Establecer la conexión WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Agregar el cliente a la lista de clientes
	clients[conn] = true

	// Escuchar mensajes del cliente
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error al leer mensaje:%v", err)
			break
		}

		// Aquí procesas el mensaje recibido
		// Ejemplo: reenviar el mensaje a todos los clientes conectados
		for client := range clients {
			if err := client.WriteMessage(messageType, p); err != nil {
				fmt.Printf("Error enviando mensaje:%v", err)
				client.Close()
				delete(clients, client) // Eliminar cliente si no se puede enviar el mensaje
			}
		}
	}
}

// Función para manejar el cierre de la conexión de un cliente
func RemoveClient(client *websocket.Conn) {
	client.Close()
	delete(clients, client)
}
