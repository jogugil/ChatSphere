package api

import (
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		// Si hay un error en el body de la solicitud, devolver un error HTTP 400
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "nok",
			"message": err.Error(),
		})
		return
	}

	// Obtener la instancia del singleton
	secMod := services.NewSecModServidorChat()

	// Parsear el IdSala como UUID
	idSalaUUID, err := uuid.Parse(requestData.IdSala)
	if err != nil {
		// Si el IdSala no es un UUID válido, devolver un error
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "nok",
			"message": "IdSala no es un UUID válido",
		})
		return
	}

	// Llamar al método para enviar el mensaje
	err = secMod.GestionSalas.EnviarMensaje(idSalaUUID, requestData.Nickname, requestData.Mensaje)
	if err != nil {
		// Si hay un error al enviar el mensaje, devolver un error con status 400
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "nok",
			"message": err.Error(),
		})
		return
	}

	// Responder con un JSON de éxito si el mensaje se envía correctamente
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
