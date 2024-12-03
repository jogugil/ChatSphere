package comm

import (
	"backend/api"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Permitir conexiones de cualquier origen
	},
}

// Manejador para WebSocket
func WebSocketHandler(c *gin.Context) {
	// Establecer la conexión WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error al establecer WebSocket:", err)
		return
	}
	defer conn.Close()

	// Leer el mensaje del cliente
	// Leer el mensaje del cliente
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error al leer mensaje WebSocket:", err)
			break
		}

		// Procesar el mensaje dependiendo de la petición
		go func() {
			var response http.ResponseWriter
			// En este punto, identificas el tipo de mensaje (por ejemplo, /listmenssage o /listusers)
			// y llamas a la lógica correspondiente
			if string(msg) == "/listmenssage" {
				// Llamamos al handler del API para obtener los mensajes
				callPostListHandler(c.Request, response) // Función que invoca la lógica de PostListHandler
			} else if string(msg) == "/listusers" {
				// Llamamos al handler del API para obtener los usuarios
				callPostUsersHandler(c.Request, response) // Función que invoca la lógica de PostUsersHandler
			} else if string(msg) == "/heat" {
				// Llamamos al handler del API para obtener los usuarios
				callheatHandler(c.Request, response) // Función que invoca la lógica de PostUsersHandler
			}

			// Enviar la respuesta al cliente WebSocket
			err := conn.WriteJSON(response)
			if err != nil {
				log.Println("Error al enviar mensaje WebSocket:", err)
			}
		}()
	}
}
func callPostListHandler(requestData *http.Request, w http.ResponseWriter) {
	api.PostListHandler(w, requestData)
}
func callPostUsersHandler(requestData *http.Request, w http.ResponseWriter) {
	api.PostUsersHandler(w, requestData)
}

// Tu función que maneja el "heat" y devuelve el código 202
func callheatHandler(requestData *http.Request, w http.ResponseWriter) error {
	// Si todo va bien, se establece el código de estado 202 y se envía la respuesta
	w.WriteHeader(http.StatusAccepted)                 // Código 202
	w.Header().Set("Content-Type", "application/json") // Si la respuesta es JSON, ajusta el Content-Type

	// O puedes enviar un mensaje de éxito vacío o algún contenido si es necesario
	w.Write([]byte(`{"status": "OK", "message": "Request accepted for processing"}`))

	return nil
}
