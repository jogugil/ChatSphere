package api

import (
	"backend/models"
	"backend/services"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
)

// Estructura para los usuarios activos
type UsuarioActivo struct {
	Nickname         string `json:"nickname"`
	HoraUltimaAccion string `json:"horaUltimaAccion"`
}

// Estructura para la respuesta general
type ResponseUser struct {
	Status          string          `json:"status"`
	Message         string          `json:"message"`
	TokenSesion     string          `json:"tokenSesion"`
	Nickname        string          `json:"nickname"`
	IdSala          string          `json:"idSala"`
	UsuariosActivos []UsuarioActivo `json:"data,omitempty"`
}

func PostUsersHandler(msg []byte) []byte {
	var requestData struct {
		IdSala      string `json:"idSala"`
		TokenSesion string `json:"tokenSesion"`
		Nickname    string `json:"nickname"`
	}

	// Crear un canal para pasar la respuesta
	respChan := make(chan Response)

	// Validar el token de sesión en una goroutine para no bloquear la respuesta
	go func() {
		// Decodificar la solicitud
		if err := json.Unmarshal(msg, &requestData); err != nil {
			log.Printf("Error al decodificar la solicitud: %v", err)
			respChan <- Response{
				Status:  "NOK",
				Message: "Solicitud incorrecta. Error al decodificar la solicitud cod:00",
			}
			return
		}
		log.Printf("PostListHandler: Datos de solicitud decodificados: %+v\n", requestData)
		// Obtener la instancia del singleton
		secMod,err  := services.GetSecModServidorChat()
		if err != nil {
			// Si el IdSala no es un UUID válido, devolver un error
			log.Printf("Error al obtener el servidor chat. : %v", err)
			respChan <- Response{
				Status:  "NOK",
				Message: "Servicio  chat no disponible",
			}
			return
		}
		log.Printf("Singleton de servidor de chat obtenido: %v", secMod)

		// Validar IdSala
		idSala, err := uuid.Parse(requestData.IdSala)
		if err != nil {
			log.Printf("Error al parsear IdSala: %v", err)
			respChan <- Response{
				Status:  "NOK",
				Message: "Error al parsear IdSala cod:00",
			}
			return
		}
		log.Printf("IdSala validado correctamente: %v", idSala)
		log.Println("PostUsersHandler: Validando token de sesión.")
		_, err = models.ValidarTokenSesion(requestData.TokenSesion)
		if err != nil {
			log.Printf("PostUsersHandler: Error al validar el token: %v", err)
			respChan <- Response{
				Status:  "NOK",
				Message: "Sesión de usuario inválida, por favor inicie sesión nuevamente.",
			}
			return
		}

		// Obtener usuarios activos
		usuarios, err := secMod.GestionUsuarios.ObtenerUsuariosActivos(idSala)
		if err != nil {
			log.Printf("Error al obtener usuarios activos: %v", err)

			respChan <- Response{
				Status:  "NOK",
				Message: fmt.Sprintf("Error al obtener usuarios activos: %s", err.Error()),
			}
			return
		}
		log.Printf("Usuarios activos obtenidos para la sala %v: %v", idSala, usuarios)

		// Construir la respuesta con la información de la sala y los usuarios activos
		var usuariosActivos []UsuarioActivo
		for _, usuario := range usuarios {
			usuariosActivos = append(usuariosActivos, UsuarioActivo{
				Nickname:         usuario.Nickname,
				HoraUltimaAccion: usuario.HoraUltimaAccion.Format("2006-01-02 15:04:05"), // Formato estándar de fecha y hora
			})
		}

		// Preparar la respuesta
		respChan := make(chan ResponseUser)

		respChan <- ResponseUser{
			Status:          "OK",
			Message:         "Usuarios activos obtenidos correctamente.",
			TokenSesion:     requestData.TokenSesion,
			Nickname:        requestData.Nickname,
			IdSala:          requestData.IdSala,
			UsuariosActivos: usuariosActivos,
		}
		log.Printf("Respuesta preparada: %+v", respChan)

	}()
	// Esperar la respuesta del canal
	response := <-respChan

	// Serializar la respuesta a JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("PostListHandler: Error al serializar la respuesta:%v", err)
		respChan <- Response{
			Status:  "NOK",
			Message: "Error interno del servidor. No hay mensajes nuevos.PostListHandler: Error al serializar la respuesta",
		}
		response := <-respChan
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Printf("PostListHandler: Error al serializar la respuesta:%v", err)
		}
		return responseJSON
	}

	// Enviar la respuesta por WebSocket
	return responseJSON
}
