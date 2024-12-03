package api

import (
	"backend/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginHandler(c *gin.Context) {
	var requestData struct {
		Nickname string `json:"nickname"`
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

	// Llamar a EjecutarLogin con un nickname
	usuario, err := secMod.EjecutarLogin("usuarioEjemplo")
	if err != nil {
		// Si hay un error al hacer login, devolver el error con status 400
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "nok",
			"message": err.Error(),
		})
		return
	}

	// Responder con un JSON de éxito si el login es exitoso
	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"token":    usuario.Token,
		"nickname": usuario.Nickname,
		"idsala":   usuario.Sala.ID,    // Sala por defecto
		"namesala": usuario.Sala.Name,  // Nombre de la sala
	})
}
