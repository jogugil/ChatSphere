package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type LoginResponse struct {
	Status   string `json:"status"`
	Token    string `json:"token"`
	Nickname string `json:"nickname"`
	IdSala   string `json:"idsala"`
	Namesala string `json:"namesala"`
}

func login(nickname string) (*LoginResponse, error) {
	url := "http://localhost:8081/login" // URL del servidor de chat
	// Crear el cuerpo de la solicitud con el nickname
	data := map[string]string{
		"nickname": nickname,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Realizar la solicitud POST
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decodificar la respuesta JSON
	var loginResponse LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		return nil, err
	}

	// Verificar si el login fue exitoso
	if loginResponse.Status != "ok" {
		return nil, fmt.Errorf("error en el login: %s", loginResponse.Status)
	}

	return &loginResponse, nil
}

func main() {
	nickname := "usuarioEjemplo"
	loginResponse, err := login(nickname)
	if err != nil {
		log.Fatalf("Error al iniciar sesión: %v", err)
	}

	fmt.Printf("Login exitoso! Token: %s, Sala: %s\n", loginResponse.Token, loginResponse.Namesala)
	// Ahora tienes el token, el id de la sala y el nickname
}