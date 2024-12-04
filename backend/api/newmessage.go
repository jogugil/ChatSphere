package api

import (
	"backend/services"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
)

func NewMessageHandler(c *gin.Context) {
	var requestData struct {
		Nickname    string `json:"nickname"`
		IdSala      string `json:"idsala"`
		NameSala    string `json:"namesala"`
		TokenSesion string `json:"tokensession"`
		Mensaje     string `json:"mensaje"`
	}

	// Decodificar los datos JSON de la solicitud
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// Log para mostrar error al procesar el JSON
		log.Printf("Error al procesar los datos JSON: %v", err)

		// Si hay un error en el body de la solicitud, devolver un error HTTP 400
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "nok",
			"message": err.Error(),
		})
		return
	}

	// Log para ver los datos recibidos
	log.Printf("Datos recibidos: %+v", requestData)

	// Obtener la instancia del singleton
	secMod, err := services.GetSecModServidorChat ()
	if err != nil {
		// Si el IdSala no es un UUID válido, devolver un error
		log.Printf("Error al obtener el servidor chat. : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "nok",
			"message": "Servicio  chat no disponible",
		})
		return
	}
	// Parsear el IdSala como UUID
	idSalaUUID, err := uuid.Parse(requestData.IdSala)
	if err != nil {
		// Si el IdSala no es un UUID válido, devolver un error
		log.Printf("Error al parsear IdSala: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "nok",
			"message": "IdSala no es un UUID válido",
		})
		return
	}

	 
	// Llamar al método para enviar el mensaje	
	err = secMod.EjecutarEnvioMensaje(requestData.Nickname, requestData.TokenSesion, requestData.Mensaje, idSalaUUID)
	if err != nil {
		// Log para mostrar error al intentar enviar el mensaje
		log.Printf("Error al enviar el mensaje: %v", err)

		// Si hay un error al enviar el mensaje, devolver un error con status 400
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "nok",
			"message": err.Error(),
		})
		return
	}

	// Log para confirmar que el mensaje fue enviado exitosamente
	log.Printf("Mensaje enviado exitosamente a la sala %s: %s", requestData.NameSala, requestData.Mensaje)

	// Responder con un JSON de éxito si el mensaje se envía correctamente
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}