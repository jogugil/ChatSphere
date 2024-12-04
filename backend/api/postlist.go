package api

import (
	"backend/entities"
	"backend/models"
	"backend/services"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
)

// Estructura para recibir los datos de la solicitud
type RequestData struct {
	Operacion       string `json:"operacion"`
	IdUltimoMensaje string `json:"idUltimoMensaje"`
	TokenSesion     string `json:"tokenSesion"`
	Nickname        string `json:"nickName"`
	IdSala          string `json:"idSala"`
}

type MessageResponse struct {
	IdMensaje uuid.UUID `json:"idMensaje"`
	Nickname  string    `json:"nickname"`
	Texto     string    `json:"texto"`
}

// Estructura para la respuesta
type Response struct {
	Status      string            `json:"status"`
	Message     string            `json:"message"`
	TokenSesion string            `json:"tokenSesion"`
	Nickname    string            `json:"nickname"`
	IdSala      string            `json:"idSala"`
	Data        []MessageResponse `json:"data,omitempty"` // Lista de mensajes si existen
}

// Convertir los mensajes a la estructura de respuesta
func ConvertirMensajes(mensajes []entities.Mensaje) []MessageResponse {
	log.Println("ConvertirMensajes: Iniciando conversión de mensajes.")
	var respuesta []MessageResponse
	for _, mensaje := range mensajes {
		respuesta = append(respuesta, MessageResponse{
			IdMensaje: mensaje.IDM,
			Nickname:  mensaje.Nickname,
			Texto:     mensaje.Mensaje,
		})
	}
	log.Printf("ConvertirMensajes: Se han convertido %d mensajes.\n", len(respuesta))
	return respuesta
}

// Manejador para la ruta POST /messagelist
func PostListHandler(msg []byte) []byte {
	log.Println("PostListHandler: Iniciando el manejo de la solicitud POST para la lista de mensajes.")

	// Decodificar el cuerpo de la solicitud
	var requestData RequestData

	// Crear un canal para pasar la respuesta
	respChan := make(chan Response)

	// Ejecutar el manejo de la solicitud en una goroutine
	go func() {
		err := json.Unmarshal(msg, &requestData)
		log.Printf("PostListHandler: Datos de solicitud decodificados: %+v\n", requestData)
		if err != nil {
			log.Printf ("PostListHandler: Error al decodificar los datos: %v", err)
			// Enviar error por WebSocket
			respChan <- Response{
				Status:  "NOK",
				Message: "Solicitud incorrecta. cod:00",
			}
			return
		}
		// Obtener la instancia del singleton
		secMod,err := services.GetSecModServidorChat()
		if err != nil {
			// Si el IdSala no es un UUID válido, devolver un error
			log.Printf("Error al obtener el servidor chat. : %v", err)
			respChan <- Response{
				Status:  "NOK",
				Message: "Servicio  chat no disponible",
			}
			return
		}
		log.Println("PostListHandler: Obteniendo usuario por token de sesión.")
		user, err := secMod.GestionUsuarios.BuscarUsuarioPorToken(requestData.TokenSesion)
		if err != nil {
			log.Printf ("PostListHandler: Error al buscar el usuario por token: %v", err)
			respChan <- Response{
				Status:  "NOK",
				Message: "Sesión de usuario inválida, por favor inicie sesión nuevamente. cod:01",
			}
			return
		}

		log.Println("PostListHandler: Validando token de sesión.")
		_, err = models.ValidarTokenSesion(requestData.TokenSesion)
		if err != nil {
			log.Printf("PostListHandler: Error al validar el token: %v", err)
			respChan <- Response{
				Status:  "NOK",
				Message: "Sesión de usuario inválida, por favor inicie sesión nuevamente. cod:02",
			}
			return
		}

		tokenUser := user.Token
		if tokenUser != requestData.TokenSesion {
			log.Println("PostListHandler: Error de validación de token: tokenUser != requestData.TokenSesion")
			respChan <- Response{
				Status:  "NOK",
				Message: "Sesión de usuario inválida, por favor inicie sesión nuevamente. cod:03",
			}
			return
		}

		log.Println("PostListHandler: Validando ID de sala.")
		idSalaUUID, err := uuid.Parse(requestData.IdSala)
		if err != nil {
			log.Printf("PostListHandler: Error al validar el IdSala: %v", err)
			respChan <- Response{
				Status:  "NOK",
				Message: "Sala de chat inválida, por favor inicie sesión nuevamente. cod:04",
			}
			return
		}

		log.Println("PostListHandler: Llamando al servicio para obtener mensajes desde el IdUltimoMensaje.")
		idmensaje, err := uuid.Parse(requestData.IdUltimoMensaje)
		if err != nil {
			log.Printf("PostListHandler: Error al validar el idmensaje: %v", err)
			respChan <- Response{
				Status:  "NOK",
				Message: "idmensaje no es un UUID válido",
			}
			return
		}

		mensajes, err := secMod.GestionSalas.ObtenerMensajesDesdeId(idSalaUUID, idmensaje)
		if err != nil {
			log.Printf("PostListHandler: Error al obtener los mensajes: %v", err)
			respChan <- Response{
				Status:  "NOK",
				Message: fmt.Sprintf("Error al obtener los mensajes: %v", err),
			}
			return
		}

		// Si no hay mensajes nuevos, devolver solo OK
		if len(mensajes) == 0 {
			log.Println("PostListHandler: No hay mensajes nuevos.")
			respChan <- Response{
				Status:  "OK",
				Message: "No hay mensajes nuevos.",
			}
			return
		}

		// Convertir los mensajes a la respuesta adecuada
		log.Println("PostListHandler: Convertiendo mensajes a la estructura de respuesta.")
		mensajesResponse := ConvertirMensajes(mensajes)

		// Si hay mensajes, devolver la lista con un OK
		log.Println("PostListHandler: Devolviendo mensajes obtenidos correctamente.")
		respChan <- Response{
			Status:      "OK",
			Message:     "Mensajes obtenidos correctamente.",
			TokenSesion: requestData.TokenSesion,
			Nickname:    requestData.Nickname,
			IdSala:      requestData.IdSala,
			Data:        mensajesResponse,
		}
	}()

	// Esperar la respuesta del canal
	response := <-respChan

	// Serializar la respuesta a JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf ("PostListHandler: Error al serializar la respuesta: %v", err)
		respChan <- Response{
			Status:  "NOK",
			Message: "Error interno del servidor. No hay mensajes nuevos.PostListHandler: Error al serializar la respuesta",
		}
		response := <-respChan
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Printf ("PostListHandler: Error al serializar la respuesta: %v", err)
		}
		return responseJSON
	}

	// Enviar la respuesta por WebSocket
	return responseJSON
}
