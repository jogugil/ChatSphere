package main

import (
	"backend/api"
	"backend/comm"
	"backend/services"
	"backend/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
 
)

// Para manejar las conexiones WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Asegúrate de que esta configuración se adapte a tus necesidades de seguridad
		return true
	},
}
 

func main() {
	log.SetFlags(log.Lshortfile)
	utils.CargarVariablesDeEntorno() 
	// Creamos un nuevo gestor de salas dentro del singleton que ocntrola lalogica del servidor
	// Obtener la instancia del singleton
	secMod := services.GetSecModServidorChat()
	salasManager := &secMod.GestionSalas

	// Cargamos las salas fijas desde el archivo de configuración
	err := salasManager.CargarSalasFijasDesdeArchivo("salas_config.json")
	if err != nil {
		log.Fatalf("Error al cargar las salas: %v", err)
	}

	// Imprimimos la sala principal
	fmt.Printf("Sala Principal: %s, ID: %s\n", salasManager.SalaPrincipal.Name, salasManager.SalaPrincipal.ID)

	// Imprimimos las salas fijas
	for id, sala := range salasManager.SalasFijas {
		fmt.Printf("Sala Fija: %s, ID: %s\n", sala.Name, id)
	}

	r := gin.Default()

	// Configurar las rutas para la API REST
	r.POST("/login", api.LoginHandler)
	r.POST("/newmessage", api.NewMessageHandler)

	// Configurar WebSocket para manejar las conexiones
	r.GET("/ws", comm.WebSocketHandler)
	// Configura las rutas de la aplicación
	server, err := utils.ObtenerVariableDeEntorno("NameServer")
	if err != nil {
		log.Fatalf("Error al iniciar el servidor. Nombre del servidor no configurado: %v", err)
		server = "localhost" // Valor por defecto si no se configura
	}

	// Inicia el servidor en el puerto 8081
	port, err := utils.ObtenerVariableDeEntorno("PortServer")
	if err != nil {
		port = "8081"
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

	// Pasar las variables de entorno a la función que inicia el servidor
	address := fmt.Sprintf("%s:%s", server, port)
	log.Printf("Servidor escuchando en %s", address)

	// Iniciar el servidor de forma concurrente, y manejar tanto HTTP como WebSocket
	go func() {
		err := r.Run(address)
		if err != nil {
			log.Fatalf("Error al iniciar el servidor HTTP: %v", err)
		}
	}()

	// Puedes poner más lógica aquí si necesitas un servidor HTTP y WebSocket independiente
	select {} // Mantener el servidor activo
}
