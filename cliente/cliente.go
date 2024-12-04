package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
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

type MessageResponse struct {
	IdMensaje uuid.UUID `json:"idMensaje"`
	Nickname  string    `json:"nickname"`
	Texto     string    `json:"texto"`
}

// Estructura para la respuesta
type Response struct {
	Status      string            `json:"status"`
	Message     string            `json:"message"`
	TokenSesion string            `json:"tokenSesion"`
	Nickname    string            `json:"nickname"`
	IdSala      string            `json:"idSala"`
	Data        []MessageResponse `json:"data,omitempty"` // Lista de mensajes si existen
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
	// Intentar reconectar en caso de desconexión
	for {
		conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8081/ws", nil)
		if err != nil {
			log.Printf("Error al conectar WebSocket, reintentando: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		return conn, nil
	}
}

func obtenerMensajes(conn *websocket.Conn, nickname, idsala, token, ultimoIdMensaje string) ([]MessageResponse, error) {
	// Crear el mensaje de solicitud
	requestData := map[string]string{
		"operacion":       "listmenssage",
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

	// Verificar si la respuesta es un objeto único o un array de mensajes
	var mensajeIndividual Response
	err = json.Unmarshal(msg, &mensajeIndividual)
	if err != nil {
		// Si no hubo error al deserializar como mensaje individual, lo empaquetamos en un slice
		return nil, fmt.Errorf("error al deserializar la respuesta lista de mensajes nuevos: %v", err)
	}
	return mensajeIndividual.Data, nil
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

	// Obtener los mensajes
	mensajes, err := obtenerMensajes(conn, loginResp.Nickname, loginResp.Idsala, loginResp.Token, "00000000-0000-0000-0000-000000000000")
	if err != nil {
		log.Fatalf("Error al obtener mensajes: %v", err)
	} else {
		// Imprimir los mensajes recibidos
		fmt.Println("Mensajes recibidos 1:")
		MostrarMensajes(mensajes)
	}

	// Ahora enviamos otro mensaje e intentamos recuperarlo
	time.Sleep(2 * time.Second) // Espera 2 segundos antes de cerrar

	// Enviar un segundo mensaje
	err = enviarMensaje(loginResp.Token, loginResp.Idsala, "Hola, este es un segundo mensaje de prueba.")
	if err != nil {
		log.Fatalf("Error al enviar el mensaje: %v", err)
	}
	fmt.Println("Segundo mensaje enviado correctamente")

	// Intentamos recuperar la lista de mensajes nuevos (el segundo mensaje enviado)
	mensajes, err_m2 := obtenerMensajes(conn, loginResp.Nickname, loginResp.Idsala, loginResp.Token, "00000000-0000-0000-0000-000000000000")
	if err_m2 != nil {
		log.Fatalf("Error al obtener mensajes: %v", err_m2)
	} else {
		// Imprimir los mensajes recibidos
		fmt.Println("Mensajes recibidos 2:")
		MostrarMensajes(mensajes)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error al cerrar la conexión WebSocket Cliente: %v", err)
		} else {
			log.Println("Conexión WebSocket Cliente cerrada correctamente")
		}
	}()
}

// Función para mostrar los elementos
func MostrarMensajes(mensajes []MessageResponse) {
	for i, mensaje := range mensajes {
		fmt.Printf("Mensaje %d:\n", i+1)
		fmt.Printf("  ID: %s\n", mensaje.IdMensaje.String())
		fmt.Printf("  Nickname: %s\n", mensaje.Nickname)
		fmt.Printf("  Texto: %s\n", mensaje.Texto)
	}
}
