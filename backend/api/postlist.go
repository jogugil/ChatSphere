package api

import (
	"backend/models"  
	"backend/entities"   // Aquí se manejaría la validación del token de sesión
	"backend/services" // Asegúrate de importar el paquete de servicios
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// Estructura para recibir los datos de la solicitud
type RequestData struct {
	IdUltimoMensaje string `json:"idUltimoMensaje"`
	TokenSesion     string `json:"tokenSesion"`
	IdSala          string `json:"idSala"`
	SalaNombre      string `json:"salaNombre"`
}

// Estructura para la respuesta
type Response struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    []MensajeResponse   `json:"data,omitempty"` // Lista de mensajes si existen
}

type MensajeResponse struct {
	IdMensaje string `json:"idMensaje"`
	Nickname  string `json:"nickname"`
	Texto     string `json:"texto"`
}

func ConvertirMensajes(mensajes []entities.Mensaje) []MensajeResponse {
	var respuesta []MensajeResponse
	for _, mensaje := range mensajes {
		respuesta = append(respuesta, MensajeResponse{
			IdMensaje: mensaje.Id,
			Nickname:  mensaje.Nickname,
			Texto:     mensaje.Mensaje,
		})
	}
	return respuesta
}

// Manejador para la ruta POST /messagelist
func PostListHandler(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la solicitud
	var requestData RequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		log.Println("Error al decodificar los datos:", err)
		http.Error(w, "Solicitud incorrecta", http.StatusBadRequest)
		return
	}

	// Obtener la instancia del singleton
	secMod := services.NewSecModServidorChat()
	user, err := secMod.GestionUsuarios.BuscarUsuarioPorToken (requestData.TokenSesion)
	if err != nil {
		log.Println("Error BuscarUsuarioPorToken:", err)
		// Si el token no es válido, devolver error
		response := Response{
			Status:  "NOK",
			Message: "Sesión de usuario inválida, por favor inicie sesión nuevamente. cod:01",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
 
	}
	_, err = models.ValidarTokenSesion(requestData.TokenSesion)
	if err != nil {
		fmt.Println("Error al validar el token:", err)
		// Si el token no es válido, devolver error
		response := Response{
			Status:  "NOK",
			Message: "Sesión de usuario inválida, por favor inicie sesión nuevamente.cod:02",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	tokenUser := user.Token
	if tokenUser != requestData.TokenSesion {
		fmt.Println("Error al validar el token: tokenUser != requestData.TokenSesion")
		// Si el token no es válido, devolver error
		response := Response{
			Status:  "NOK",
			Message: "Sesión de usuario inválida, por favor inicie sesión nuevamente.cod:03",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}
	idSalaUUID, err := uuid.Parse(requestData.IdSala)
	if err != nil {
		// Si el IdSala no es un UUID válido, devolver un error
		http.Error(w, "IdSala no es un UUID válido", http.StatusBadRequest)
				// Si el token no es válido, devolver error
				response := Response{
					Status:  "NOK",
					Message: "Sala de chat inválida, por favor inicie sesión nuevamente.cod:04",
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response)
				return
	}
	// Llamar al servicio para obtener los mensajes desde el idUltimoMensaje hasta el actual
	mensajes, err := secMod.GestionSalas.ObtenerMensajesDesdeId(idSalaUUID, requestData.IdUltimoMensaje)
	if err != nil {
		// Si hubo un error al obtener los mensajes
		response := Response{
			Status:  "NOK",
			Message: fmt.Sprintf("Error al obtener los mensajes: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Si no hay mensajes nuevos, devolver solo OK
	if len(mensajes) == 0 {
		response := Response{
			Status:  "OK",
			Message: "No hay mensajes nuevos.",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
	// Convertir a MensajeResponse
	mensajesResponse := ConvertirMensajes(mensajes)
	// Si hay mensajes, devolver la lista con un OK
	response := Response{
		Status:  "OK",
		Message: "Mensajes obtenidos correctamente.",
		Data:    mensajesResponse,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}