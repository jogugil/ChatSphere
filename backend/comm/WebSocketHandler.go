package comm

import (
	"backend/api"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Manejador para WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Permitir conexiones de cualquier origen
	},
}

func WebSocketHandler(c *gin.Context) {
	// Establecer la conexión WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error al establecer WebSocket:", err)
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("Error al cerrar la conexión WebSocket:", err)
		} else {
			log.Println("Conexión WebSocket cerrada correctamente")
		}
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error al leer mensaje WebSocket:", err)
			break
		}

		var data map[string]string
		err = json.Unmarshal(msg, &data)
		if err != nil {
			log.Println("Error al deserializar el mensaje:", err)
			msgr, _ := json.Marshal(`{"status": "error", "message": "Error al deserializar el mensaje"}`)
			conn.WriteMessage(websocket.TextMessage, msgr)
			continue
		}

		operacion, exists := data["operacion"]
		if !exists {
			log.Println("No se encontró la clave 'operacion' en el mensaje")
			msgr, _ := json.Marshal(`{"status": "error", "message": "No se encontró la clave 'operacion' en el mensaje"}`)
			conn.WriteMessage(websocket.TextMessage, msgr)
			continue
		}

		var response []byte
		switch operacion {
		case "listmenssage":
			response = callPostListHandler(msg)
		case "listusers":
			response = callPostUsersHandler(msg)
		case "heat":
			response = callheatHandler(msg)
		default:
			response, _ = json.Marshal(`{"status": "error", "message": "Comando no reconocido"}`)
		}

		err = conn.WriteMessage(websocket.TextMessage, response)
		if err != nil {
			log.Println("Error al enviar mensaje WebSocket:", err)
		}
	}
}

// Aquí solo retornas una cadena de texto (JSON) que representa la respuesta
func callPostListHandler(msg []byte) []byte {
	// Llamamos a la función correspondiente del API
	// Aquí debes definir tu lógica de la función, por ejemplo:
	result := api.PostListHandler(msg) // Llamar a la función de API pasando la request
	return result                      // Retornamos la respuesta como un string (puede ser un JSON)
}

func callPostUsersHandler(msg []byte) []byte {
	// Llamamos a la función correspondiente del API
	// Aquí debes definir tu lógica de la función
	result := api.PostUsersHandler(msg) // Llamar a la función de API pasando la request
	return result                       // Retornamos la respuesta como un string (puede ser un JSON)
}

// Tu función que maneja el "heat" y devuelve el código 202
func callheatHandler(msg []byte) []byte {
	log.Printf("callheatHandler procesando solicitud: %v", msg)
	// Si todo va bien, se establece el código de estado 202 y se envía la respuesta
	errorJSON, err := json.Marshal(`{"status": "OK", "message": "Request accepted for processing"}`)
	if err != nil {
		log.Println("sendErrorResponse: Error al serializar el mensaje de error:", err)
		return errorJSON
	}
	return errorJSON
}
