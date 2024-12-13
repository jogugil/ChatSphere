package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Estructura de datos para el login
type LoginResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
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
	ListMessage []MessageResponse `json:"data,omitempty"` // Lista de mensajes si existen
}

// Estructura de datos para enviar en la solicitud
type RequestUserData struct {
	IdSala      string `json:"idSala"`
	TokenSesion string `json:"tokenSesion"`
	Nickname    string `json:"nickname"`
	Operacion   string `json:"operacion"`
}

// Estructura de la respuesta esperada
type ResponseUser struct {
	Status          string `json:"status"`
	Message         string `json:"message"`
	TokenSesion     string `json:"tokenSesion"`
	Nickname        string `json:"nickname"`
	IdSala          string `json:"idSala"`
	UsuariosActivos []struct {
		Nickname         string `json:"nickname"`
		HoraUltimaAccion string `json:"horaUltimaAccion"`
	} `json:"data,omitempty"`
}

func login(nickname string) (LoginResponse, error) {
	// Crear los datos para el login
	data := map[string]string{"nickname": nickname}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return LoginResponse{}, err
	}

	// Crear la solicitud HTTP
	req, err := http.NewRequest("POST", "http://localhost:8081/login", bytes.NewBuffer(jsonData))
	if err != nil {
		return LoginResponse{}, err
	}

	// Añadir la cabecera 'x-gochat' con el valor 'http://localhost:8081'
	req.Header.Set("x-gochat", "http://localhost:8081")
	req.Header.Set("Content-Type", "application/json")

	// Hacer la petición HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
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

func enviarMensaje(token, idsala, nickName, mensaje string) error {
	// Crear los datos para el mensaje
	data := map[string]string{
		"nickname":     nickName, // Cambiar por el nickname
		"idsala":       idsala,
		"namesala":     "SalaEjemplo", // Cambiar por el nombre de la sala
		"tokensession": token,
		"mensaje":      mensaje,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Crear la solicitud HTTP
	req, err := http.NewRequest("POST", "http://localhost:8081/newmessage", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Añadir la cabecera 'x-gochat' con el valor 'http://localhost:8081'
	req.Header.Set("x-gochat", "http://localhost:8081")
	req.Header.Set("Content-Type", "application/json")

	// Hacer la petición HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
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
		// Crear la solicitud HTTP personalizada para WebSocket
		req, err := http.NewRequest("GET", "ws://localhost:8081/ws", nil)
		if err != nil {
			log.Printf("Error al crear la solicitud HTTP: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Añadir la cabecera 'x-gochat' con el valor 'http://localhost:8081'
		req.Header.Set("x-gochat", "http://localhost:8081")

		// Establecer la conexión WebSocket con la solicitud personalizada
		conn, _, err := websocket.DefaultDialer.Dial(req.URL.String(), req.Header)
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
	return mensajeIndividual.ListMessage, nil
}

func obtenerUsuarios(conn *websocket.Conn, nickname, idsala, token string) (ResponseUser, error) {
	// Declarar la estructura RequestUserData fuera de la llamada de la función
	requestData := struct {
		IdSala      string `json:"idSala"`
		TokenSesion string `json:"tokenSesion"`
		Nickname    string `json:"nickname"`
		Operacion   string `json:"operacion"`
	}{
		IdSala:      idsala,
		TokenSesion: token,
		Nickname:    nickname,
		Operacion:   "listusers",
	}

	// Convertir a JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return ResponseUser{}, err // Devolver un ResponseUser vacío y el error
	}

	// Enviar solicitud de lista de usuarios
	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		return ResponseUser{}, err // Devolver un ResponseUser vacío y el error
	}

	// Leer la respuesta del WebSocket
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return ResponseUser{}, err // Devolver un ResponseUser vacío y el error
	}

	// Verificar si la respuesta es un objeto único (ResponseUser)
	var responUser ResponseUser
	err = json.Unmarshal(msg, &responUser)
	if err != nil {
		// Si no hubo error al deserializar como mensaje individual, devolver el error
		return ResponseUser{}, fmt.Errorf("error al deserializar la respuesta: %v", err)
	}

	// Devolver la respuesta deserializada
	return responUser, nil
}

// Ejecutar peticiones periódicas
func ejecutarPeticionesPeriodicas(conn *websocket.Conn, nickname, idsala, token string, stopCh chan bool) {
	ticker := time.NewTicker(10 * time.Second) // intervalo de 10 segundos
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Intentamos recuperar los usuarios activos
			usuarios, err_u := obtenerUsuarios(conn, nickname, idsala, token)
			if err_u != nil {
				log.Printf("Error al obtener usuarios: %v", err_u)
				// Si ocurre un error, intentamos reconectar o terminamos el hilo.
				// Aquí podrías implementar lógica de reconexión, o simplemente terminar.
				// Si deseas salir del bucle, podrías cerrar el canal stopCh.
				stopCh <- true
				return
			} else {
				// Imprimir los usuarios recibidos
				fmt.Println("Usuarios recibidos:")
				MostrarUsuarios(usuarios)
			}

			// Intentamos recuperar la lista de mensajes nuevos
			mensajes, err_m2 := obtenerMensajes(conn, nickname, idsala, token, "00000000-0000-0000-0000-000000000000")
			if err_m2 != nil {
				log.Printf("Error al obtener mensajes: %v", err_m2)
				stopCh <- true
				return
			} else {
				// Imprimir los mensajes recibidos
				fmt.Println("Mensajes recibidos:")
				MostrarMensajes(mensajes)
			}

		case <-stopCh:
			log.Println("Deteniendo el goroutine debido a la desconexión.")
			return
		}
	}
}
func main() {
	// Solicitar el nickname al usuario
	fmt.Print("Introduce tu nickname: ")
	var nickname string
	fmt.Scanln(&nickname)

	// Realizar login
	loginResp, err := login(nickname)
	if err != nil {
		log.Fatalf("Error en el login: %v", err)
	}
	fmt.Printf("LloginResp.status: %s \n", loginResp.Status)
	fmt.Printf("LloginResp.Message: %s \n", loginResp.Message)
	fmt.Printf("Login exitoso! Token: %s, Sala: %s\n", loginResp.Token, loginResp.Namesala)
	fmt.Printf("LloginResp.Idsala: %s \n", loginResp.Idsala)

	// Pedir el primer mensaje
	fmt.Print("Introduce el mensaje a enviar: ")
	var mensaje string
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	mensaje = scanner.Text()

	// Enviar el primer mensaje
	err = enviarMensaje(loginResp.Token, loginResp.Idsala, loginResp.Nickname, mensaje)
	if err != nil {
		log.Fatalf("Error al enviar el mensaje: %v", err)
	}
	fmt.Println("Mensaje enviado correctamente")

	// Conectar al WebSocket para obtener la lista de mensajes
	conn, err := conectarWebSocket()
	if err != nil {
		log.Fatalf("Error al conectar WebSocket: %v", err)
	}

	// Canal para controlar la detención del goroutine
	stopCh := make(chan bool)

	// Iniciar el goroutine para ejecutar las peticiones periódicas
	go ejecutarPeticionesPeriodicas(conn, loginResp.Nickname, loginResp.Idsala, loginResp.Token, stopCh)

	// Mantener el programa activo hasta recibir una señal de detención o desconexión
	// Función para esperar la detención del programa sin bloquear el hilo principal
	go func() {
		<-stopCh
		log.Println("Programa detenido debido a la desconexión.")
	}()

	//Intentamos recuperar los usuarios activos
	usuarios, err_u1 := obtenerUsuarios(conn, loginResp.Nickname, loginResp.Idsala, loginResp.Token)
	if err_u1 != nil {
		log.Fatalf("Error al obtener usuarios 1: %v", err_u1)
	} else {
		// Imprimir los mensajes recibidos
		fmt.Println("usuarios recibidos 1 :")
		MostrarUsuarios(usuarios)
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

	// Pedir el segundo mensaje
	fmt.Print("Introduce el segundo mensaje a enviar: ")
	scanner.Scan()
	mensaje2 := scanner.Text()

	// Enviar el segundo mensaje
	err = enviarMensaje(loginResp.Token, loginResp.Idsala, loginResp.Nickname, mensaje2)
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

	//Intentamos recuperar los usuarios activos
	usuarios, err_u := obtenerUsuarios(conn, loginResp.Nickname, loginResp.Idsala, loginResp.Token)
	if err_u != nil {
		log.Fatalf("Error al obtener usuarios 2: %v", err_u)
	} else {
		// Imprimir los mensajes recibidos
		fmt.Println("usuarios recibidos 2 :")
		MostrarUsuarios(usuarios)
	}

	// Ahora enviamos otro mensaje e intentamos recuperarlo
	time.Sleep(20 * time.Second) // Espera 2 segundos antes de cerrar

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

// Función MostrarUsuarios que filtra por idSala y muestra los usuarios activos
func MostrarUsuarios(usuarios ResponseUser) {
	// Mostrar el ID de la sala y los usuarios activos
	fmt.Printf("Sala: %s\n", usuarios.IdSala)
	// Iterar sobre cada ResponseUser en el slice de usuarios
	if len(usuarios.UsuariosActivos) == 0 {
		fmt.Println("No hay usuarios activos en esta sala.")
	} else {
		listuser := usuarios.UsuariosActivos
		for _, usuario := range listuser {
			// Mostrar el Nickname y la Hora de la última acción
			fmt.Printf("Usuario: %s, Última acción: %s\n", usuario.Nickname, usuario.HoraUltimaAccion)
		}
	}

	// Si no se encuentra ninguna sala con el idSala proporcionado
	fmt.Println("Sala no encontrada o no hay usuarios activos.")
}
