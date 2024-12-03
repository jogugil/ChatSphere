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
	IdMensaje uuid.UUID `json:"idMensaje"`
	Nickname  string    `json:"nickname"`
	Texto     string    `json:"texto"`
}

// Convertir los mensajes a la estructura de respuesta
func ConvertirMensajes(mensajes []entities.Mensaje) []MensajeResponse {
	log.Println("ConvertirMensajes: Iniciando conversión de mensajes.")
	var respuesta []MensajeResponse
	for _, mensaje := range mensajes {
		respuesta = append(respuesta, MensajeResponse{
			IdMensaje: mensaje.IDM,
			Nickname:  mensaje.Nickname,
			Texto:     mensaje.Mensaje,
		})
	}
	log.Printf("ConvertirMensajes: Se han convertido %d mensajes.\n", len(respuesta))
	return respuesta
}

// Manejador para la ruta POST /messagelist
func PostListHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("PostListHandler: Iniciando el manejo de la solicitud POST para la lista de mensajes.")

	// Decodificar el cuerpo de la solicitud
	var requestData RequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		log.Println("PostListHandler: Error al decodificar los datos:", err)
		http.Error(w, "Solicitud incorrecta", http.StatusBadRequest)
		return
	}
	log.Printf("PostListHandler: Datos de solicitud decodificados: %+v\n", requestData)

	// Obtener la instancia del singleton
	secMod := services.GetSecModServidorChat()
	log.Println("PostListHandler: Obteniendo usuario por token de sesión.")
	user, err := secMod.GestionUsuarios.BuscarUsuarioPorToken(requestData.TokenSesion)
	if err != nil {
		log.Println("PostListHandler: Error al buscar el usuario por token:", err)
		response := Response{
			Status:  "NOK",
			Message: "Sesión de usuario inválida, por favor inicie sesión nuevamente. cod:01",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Println("PostListHandler: Validando token de sesión.")
	_, err = models.ValidarTokenSesion(requestData.TokenSesion)
	if err != nil {
		log.Println("PostListHandler: Error al validar el token:", err)
		response := Response{
			Status:  "NOK",
			Message: "Sesión de usuario inválida, por favor inicie sesión nuevamente. cod:02",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	tokenUser := user.Token
	if tokenUser != requestData.TokenSesion {
		log.Println("PostListHandler: Error de validación de token: tokenUser != requestData.TokenSesion")
		response := Response{
			Status:  "NOK",
			Message: "Sesión de usuario inválida, por favor inicie sesión nuevamente. cod:03",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Println("PostListHandler: Validando ID de sala.")
	idSalaUUID, err := uuid.Parse(requestData.IdSala)
	if err != nil {
		log.Println("PostListHandler: Error al validar el IdSala:", err)
		http.Error(w, "IdSala no es un UUID válido", http.StatusBadRequest)
		response := Response{
			Status:  "NOK",
			Message: "Sala de chat inválida, por favor inicie sesión nuevamente. cod:04",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Println("PostListHandler: Llamando al servicio para obtener mensajes desde el IdUltimoMensaje.")
	idmensaje, err:= uuid.Parse(requestData.IdUltimoMensaje)
	if err != nil {
		log.Println("PostListHandler: Error al validar el idmensaje:", err)
		http.Error(w, "idmensaje no es un UUID válido", http.StatusBadRequest)
		response := Response{
			Status:  "NOK",
			Message: "Sala de chat inválida, por favor inicie sesión nuevamente. cod:04",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	mensajes, err := secMod.GestionSalas.ObtenerMensajesDesdeId(idSalaUUID, idmensaje)
	if err != nil {
		log.Println("PostListHandler: Error al obtener los mensajes:", err)
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
		log.Println("PostListHandler: No hay mensajes nuevos.")
		response := Response{
			Status:  "OK",
			Message: "No hay mensajes nuevos.",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convertir los mensajes a la respuesta adecuada
	log.Println("PostListHandler: Convertiendo mensajes a la estructura de respuesta.")
	mensajesResponse := ConvertirMensajes(mensajes)

	// Si hay mensajes, devolver la lista con un OK
	log.Println("PostListHandler: Devolviendo mensajes obtenidos correctamente.")
	response := Response{
		Status:  "OK",
		Message: "Mensajes obtenidos correctamente.",
		Data:    mensajesResponse,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
