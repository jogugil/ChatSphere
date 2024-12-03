package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func sendMessage(token string, idsala string, nickname string, mensaje string) error {
	url := "http://localhost:8081/newmessage" // URL del servidor de chat
	// Crear el cuerpo de la solicitud con los datos del mensaje
	data := map[string]string{
		"nickname":     nickname,
		"idsala":       idsala,
		"namesala":     "Sala Ejemplo", // Puedes poner el nombre de la sala si lo tienes
		"tokensession": token,
		"mensaje":      mensaje,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Realizar la solicitud POST
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Verificar el estado de la respuesta
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error al enviar mensaje: %s", resp.Status)
	}

	fmt.Println("Mensaje enviado exitosamente")
	return nil
}

func main() {
	// Suponiendo que ya has iniciado sesión y obtenido el token y sala
	token := "tokenDelUsuario" // Token obtenido del login
	idsala := "idDeLaSala"     // ID de la sala obtenida del login
	nickname := "usuarioEjemplo"
	mensaje := "Hola, este es un mensaje de prueba"

	err := sendMessage(token, idsala, nickname, mensaje)
	if err != nil {
		log.Fatalf("Error al enviar mensaje: %v", err)
	}
}
