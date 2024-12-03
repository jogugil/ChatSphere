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
	defer conn.Close()

	// Leer el mensaje del cliente
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error al leer mensaje WebSocket:", err)
			break
		}

		// Procesar el mensaje dependiendo de la petición
		go func(msg []byte, requestData *http.Request) {
			var response []byte

			// Aquí procesas el mensaje según la lógica de la aplicación
			if string(msg) == "/listmenssage" {
				// Llamamos al handler del API para obtener los mensajes
				response = callPostListHandler(requestData)
			} else if string(msg) == "/listusers" {
				// Llamamos al handler del API para obtener los usuarios
				response = callPostUsersHandler(requestData)
			} else if string(msg) == "/heat" {
				// Llamamos al handler del API para obtener el "heat"
				response = callheatHandler(requestData)
			} else {
				response,err = json.Marshal(`{"status": "error", "message": "Comando no reconocido"}`)
			}

			// Enviar la respuesta al cliente WebSocket
			err := conn.WriteMessage(websocket.TextMessage, []byte(response))
			if err != nil {
				log.Println("Error al enviar mensaje WebSocket:", err)
			}
		}(msg, c.Request)
	}
}

// Aquí solo retornas una cadena de texto (JSON) que representa la respuesta
func callPostListHandler(requestData *http.Request) []byte {
	// Llamamos a la función correspondiente del API
	// Aquí debes definir tu lógica de la función, por ejemplo:
	result := api.PostListHandler(requestData) // Llamar a la función de API pasando la request
	return result // Retornamos la respuesta como un string (puede ser un JSON)
}

func callPostUsersHandler(requestData *http.Request) []byte {
	// Llamamos a la función correspondiente del API
	// Aquí debes definir tu lógica de la función
	result := api.PostUsersHandler (requestData) // Llamar a la función de API pasando la request
	return result // Retornamos la respuesta como un string (puede ser un JSON)
}

// Tu función que maneja el "heat" y devuelve el código 202
func callheatHandler(requestData *http.Request) []byte {
	log.Printf("callheatHandler procesando solicitud: %v", requestData)
	// Si todo va bien, se establece el código de estado 202 y se envía la respuesta
	errorJSON, err := json.Marshal(`{"status": "OK", "message": "Request accepted for processing"}`)
	if err != nil {
		log.Println("sendErrorResponse: Error al serializar el mensaje de error:", err)
		return errorJSON
	}
	return errorJSON
}
 