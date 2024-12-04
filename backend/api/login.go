package api

import (
	"backend/services"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	// Definir estructura para recibir datos de la solicitud
	var requestData struct {
		Nickname string `json:"nickname"`
	}

	// Log de entrada de la solicitud
	fmt.Println("Recibiendo datos de solicitud...")

	// Decodificar los datos JSON de la solicitud
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// Log de error en la decodificación de datos
		fmt.Println("Error al decodificar el JSON de la solicitud:", err)

		// Si hay un error en el body de la solicitud, devolver un error HTTP 400
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "nok",
			"message": err.Error(),
		})
		return
	}

	// Log de los datos recibidos
	fmt.Println("Datos recibidos de la solicitud:", requestData)

	// Obtener la instancia del singleton
	secMod,err := services.GetSecModServidorChat ()
	if err != nil {
		// Si el IdSala no es un UUID válido, devolver un error
		log.Printf("Error al obtener el servidor chat. : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "nok",
			"message": "Servicio  chat no disponible",
		})
		return
	}
	// Llamar a EjecutarLogin con el nickname recibido
	fmt.Println("Ejecutando login para el usuario:", requestData.Nickname)
	usuario, err := secMod.EjecutarLogin(requestData.Nickname)
	if err != nil {
		// Log del error en el proceso de login
		fmt.Println("Error al ejecutar login para el usuario:", err)

		// Si hay un error al hacer login, devolver el error con status 400
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "nok",
			"message": err.Error(),
		})
		return
	}

	// Log de los datos del usuario después del login
	fmt.Printf("Login exitoso. Datos del usuario: Token: %s, Nickname: %s, Sala ID: %v, Sala Name: %s\n",
		usuario.Token, usuario.Nickname, usuario.Sala.ID, usuario.Sala.Name)

	// Responder con un JSON de éxito si el login es exitoso
	responseData := gin.H{
		"status":   "ok",
		"token":    usuario.Token,
		"nickname": usuario.Nickname,
		"idsala":   usuario.Sala.ID,   // Sala por defecto
		"namesala": usuario.Sala.Name, // Nombre de la sala
	}

	// Log de la respuesta enviada
	fmt.Println("Enviando respuesta:", responseData)

	c.JSON(http.StatusOK, responseData)
}
