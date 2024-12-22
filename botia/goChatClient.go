package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Reutilizamos las estructuras existentes para LoginResponse y otros datos

// Función para interactuar con la API de IA
func obtenerRespuestaIA(input string) (string, error) {
	apiURL := "https://api.openai.com/v1/chat/completions"
	apiKey := "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" // Añade tu API key aquí

	// Crear el payload
	payload := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "system", "content": "Eres un bot amigable en un chat grupal."},
			{"role": "user", "content": input},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Crear la solicitud HTTP
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Hacer la petición
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Leer la respuesta
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Extraer la respuesta generada por la IA
	choices := result["choices"].([]interface{})
	message := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)

	return message, nil
}

// Función principal
func main() {
	log.Print("Introduce tu nickname: ")
	var nickname string
	fmt.Scanln(&nickname)

	// Realizar login
	loginResp, err := login(nickname)
	if err != nil {
		log.Fatalf("Error en el login: %v", err)
	}

	log.Printf("Login exitoso! Token: %s, Sala: %s\n", loginResp.Token, loginResp.RoomName)

	// Conectar al WebSocket
	conn, err := conectarWebSocket()
	if err != nil {
		log.Fatalf("Error al conectar al WebSocket: %v", err)
	}
	defer conn.Close()

	// Escuchar mensajes y responder
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error al leer mensaje: %v", err)
			break
		}

		// Procesar mensaje
		var incomingMessage MessageResponse
		if err := json.Unmarshal(msg, &incomingMessage); err != nil {
			log.Printf("Error al procesar mensaje: %v", err)
			continue
		}

		log.Printf("Mensaje recibido de %s: %s", incomingMessage.Nickname, incomingMessage.MessageText)

		// Evitar responder a mensajes propios
		if incomingMessage.Nickname == nickname {
			continue
		}

		// Obtener respuesta de IA
		respuesta, err := obtenerRespuestaIA(incomingMessage.MessageText)
		if err != nil {
			log.Printf("Error al generar respuesta de IA: %v", err)
			continue
		}

		// Enviar respuesta al chat
		err = enviarMensaje(loginResp.Token, loginResp.RoomId, nickname, respuesta)
		if err != nil {
			log.Printf("Error al enviar mensaje: %v", err)
		}
	}
}
