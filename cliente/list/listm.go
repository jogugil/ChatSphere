package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func connectWebSocket() (*websocket.Conn, error) {
	url := "ws://localhost:8081/ws" // Dirección del WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func sendCommand(conn *websocket.Conn, command string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte(command))
	if err != nil {
		log.Println("Error al enviar comando:", err)
		return
	}
	fmt.Println("Comando enviado:", command)
}

func readMessages(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error al leer mensaje:", err)
			return
		}
		fmt.Println("Mensaje recibido:", string(msg))
	}
}

func main() {
	// Establecer la conexión WebSocket
	conn, err := connectWebSocket()
	if err != nil {
		log.Fatalf("Error al conectar WebSocket: %v", err)
	}
	defer conn.Close()

	// Iniciar el proceso de lectura de mensajes
	go readMessages(conn)

	// Enviar comandos para obtener la lista de mensajes y usuarios
	time.Sleep(1 * time.Second) // Esperar un momento antes de enviar los comandos
	sendCommand(conn, "/listmenssage")
	sendCommand(conn, "/listusers")

	// Mantener el proceso en ejecución para seguir recibiendo mensajes
	select {}
}
