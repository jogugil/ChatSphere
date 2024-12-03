package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Estructura de datos para el login
type LoginResponse struct {
	Status   string `json:"status"`
	Token    string `json:"token"`
	Nickname string `json:"nickname"`
	Idsala   string `json:"idsala"`
	Namesala string `json:"namesala"`
}

// Estructura de datos para el mensaje
type MessageResponse struct {
	Nickname   string `json:"nickname"`
	Idsala     string `json:"idsala"`
	Namesala   string `json:"namesala"`
	Idmensaje  string `json:"idmensaje"`
	Mensaje    string `json:"mensaje"`
	De         string `json:"de"`
}

func login(nickname string) (LoginResponse, error) {
	// Crear los datos para el login
	data := map[string]string{"nickname": nickname}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return LoginResponse{}, err
	}

	// Hacer la petición HTTP para login
	resp, err := http.Post("http://localhost:8081/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return LoginResponse{}, err
	}
	defer resp.Body.Close()

	// Leer la respuesta
	var loginResponse LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		return LoginResponse{}, err
	}

	// Retornar el resultado
	return loginResponse, nil
}

func enviarMensaje(token, idsala, mensaje string) error {
	// Crear los datos para el mensaje
	data := map[string]string{
		"nickname":     "usuarioEjemplo", // Cambiar por el nickname
		"idsala":       idsala,
		"namesala":     "SalaEjemplo", // Cambiar por el nombre de la sala
		"tokensession": token,
		"mensaje":      mensaje,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Hacer la petición HTTP para enviar el mensaje
	resp, err := http.Post("http://localhost:8081/newmessage", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Verificar que la respuesta sea correcta
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error al enviar el mensaje, status: %s", resp.Status)
	}
	return nil
}

func conectarWebSocket() (*websocket.Conn, error) {
	// Conectar al servidor WebSocket
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8081/ws", nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func obtenerMensajes(conn *websocket.Conn, nickname, idsala, token, ultimoIdMensaje string) ([]MessageResponse, error) {
	// Crear el mensaje de solicitud
	requestData := map[string]string{
		"nickName":        nickname,
		"idSala":          idsala,
		"tokenSesion":     token,
		"idUltimoMensaje": ultimoIdMensaje,
	}

	// Convertir a JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	// Enviar solicitud de lista de mensajes
	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		return nil, err
	}

	// Leer la respuesta del WebSocket
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	// Decodificar la respuesta en la lista de mensajes
	var mensajes []MessageResponse
	if err := json.Unmarshal(msg, &mensajes); err != nil {
		return nil, err
	}

	return mensajes, nil
}

func main() {
	// Realizar login
	loginResp, err := login("usuarioEjemplo")
	if err != nil {
		log.Fatalf("Error en el login: %v", err)
	}
	fmt.Printf("Login exitoso! Token: %s, Sala: %s\n", loginResp.Token, loginResp.Namesala)
	fmt.Printf("LloginResp.Idsala: %s \n", loginResp.Idsala)

	// Enviar un mensaje
	err = enviarMensaje(loginResp.Token, loginResp.Idsala, "Hola, este es un mensaje de prueba.")
	if err != nil {
		log.Fatalf("Error al enviar el mensaje: %v", err)
	}
	fmt.Println("Mensaje enviado correctamente")

	// Conectar al WebSocket para obtener la lista de mensajes
	conn, err := conectarWebSocket()
	if err != nil {
		log.Fatalf("Error al conectar WebSocket: %v", err)
	}
	defer conn.Close()

	// Obtener los mensajes
	mensajes, err := obtenerMensajes(conn, loginResp.Nickname, loginResp.Idsala, loginResp.Token, "00000000")
	if err != nil {
		log.Fatalf("Error al obtener mensajes: %v", err)
	}

	// Imprimir los mensajes recibidos
	fmt.Println("Mensajes recibidos:")
	for _, mensaje := range mensajes {
		fmt.Printf("ID Mensaje: %s, Mensaje: %s, Enviado por: %s\n", mensaje.Idmensaje, mensaje.Mensaje, mensaje.De)
	}

	// Ahora enviamso oro mensaje e intentamos recuperarlo
	time.Sleep(2 * time.Second) // Espera por ejemplo 2 segundos antes de cerrar

	// Enviar un mensaje
	err = enviarMensaje(loginResp.Token, loginResp.Idsala, "Hola, este es un segundo mensaje de prueba.")
	if err != nil {
		log.Fatalf("Error al enviar el mensaje: %v", err)
	}
	fmt.Println("Mensaje enviado correctamente")

	//intentamso recuperar la lsita d emensjaes nuevos que es el segundo enviado.
	// Obtener los mensajes
	mensajes1, err := obtenerMensajes(conn, loginResp.Nickname, loginResp.Idsala, loginResp.Token, "00000000")
	if err != nil {
		log.Fatalf("Error al obtener mensajes: %v", err)
	}
	// Imprimir los mensajes recibidos
	fmt.Println("Mensajes recibidos:")
	for _, mensaje := range mensajes1 {
		fmt.Printf("ID Mensaje: %s, Mensaje: %s, Enviado por: %s\n", mensaje.Idmensaje, mensaje.Mensaje, mensaje.De)
	}
}