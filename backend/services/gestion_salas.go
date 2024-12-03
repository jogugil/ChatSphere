package services

import (
	"backend/entities"
 
 
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
)

// GestionSalas estructura para gestionar las salas
type GestionSalas struct {
	SalasFijas    map[uuid.UUID]entities.Sala // Mapa de salas fijas en memoria
	SalaPrincipal entities.Sala               // Sala principal única
	mu            sync.Mutex                  // Mutex para proteger las estructuras de datos
	persistencia  entities.Persistencia     // Persistencia para almacenar los datos
	once          sync.Once                   // Para garantizar inicialización única por instancia
}

// Instancia única de GestionSalas
var instancia *GestionSalas

func NuevaGestionSalas(persistencia entities.Persistencia, configFile string) *GestionSalas {
	if instancia == nil {
		// Creamos una nueva instancia si no existe
		instancia = &GestionSalas{
			SalasFijas:   make(map[uuid.UUID]entities.Sala),
			persistencia: persistencia,
		}
	}
	// Garantizamos que la inicialización ocurra una sola vez
	instancia.once.Do(func() {
		// Solo se carga el archivo de configuración en la primera inicialización
		if configFile != "" {
			err := instancia.CargarSalasFijasDesdeArchivo(configFile)
			if err != nil {
				log.Fatalf("Error al cargar configuración: %v", err)
			}
		}
	})
	return instancia
}

// Función para cargar las salas fijas desde un archivo de configuración
func (sm *GestionSalas) CargarSalasFijasDesdeArchivo(configFile string) error {
	// Configurar la sala principal (siempre debe existir)
 	sm.SalaPrincipal = entities.Sala{
		ID:               uuid.New(),
		Name:              "Sala Principal",
		SalaTipo:          "Principal",
		HistoricoMensajes: entities.NuevaColaCircular(sm.persistencia),
	}
	// Leer y cargar el archivo de configuración para las salas fijas
	file, err := os.Open(configFile)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el archivo de configuración: %v", err)
	}
	defer file.Close()

	// Leemos el archivo
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("no se pudo leer el archivo de configuración: %v", err)
	}

	// Parseamos el JSON
	var salas []map[string]interface{}
	err = json.Unmarshal(byteValue, &salas)
	if err != nil {
		return fmt.Errorf("error al parsear el archivo de configuración: %v", err)
	}

	// Cargamos las salas fijas desde el archivo
	for _, salaData := range salas {
		// Parseamos el ID de la sala
		salaID, err := uuid.Parse(salaData["id"].(string))
		if err != nil {
			return fmt.Errorf("error al parsear el ID de la sala: %v", err)
		}

		// Crear la sala fija
		sala := entities.Sala{
			ID:                salaID,
			Name:              salaData["nombre"].(string),
			SalaTipo:          "Fija",
			HistoricoMensajes: entities.NuevaColaCircular(sm.persistencia),
		}

		// Agregar la sala a las salas fijas
		sm.SalasFijas[salaID] = sala

	}

	return nil
}

// CrearSalaTemporal crea una sala temporal y su cola circular asociada
func (sm *GestionSalas) CrearSalaTemporal(nombre string) entities.Sala {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Asegurarse de que la persistencia esté configurada
	if sm.persistencia == nil {
		log.Fatalf("Persistencia no inicializada. Asegúrate de configurar la persistencia antes de crear salas.")
	}

	// Crear una nueva sala temporal
	sala := entities.Sala{
		ID:                uuid.New(),
		Name:              nombre,
		SalaTipo:          "Temporal", // Definimos el tipo como temporal
		HistoricoMensajes: entities.NuevaColaCircular(sm.persistencia),
	}

	// Retornar la sala creada
	return sala
}

// ObtenerSalaPorID busca una sala por su ID en SalaPrincipal o SalasFijas
func (sm *GestionSalas) ObtenerSalaPorID(idSala uuid.UUID) (*entities.Sala, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Verificar si el ID corresponde a SalaPrincipal
	if sm.SalaPrincipal.ID == idSala {
		return &sm.SalaPrincipal, nil
	}

	// Buscar en SalasFijas
	if sala, existe := sm.SalasFijas[idSala]; existe {
		return &sala, nil
	}

	return nil, fmt.Errorf("la sala con ID %s no existe", idSala)
}

// Función para enviar un mensaje a una sala
func (sm *GestionSalas) EnviarMensaje(idSala uuid.UUID, nickname, mensaje string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sala, err := sm.ObtenerSalaPorID(idSala)
	if err != nil {
		return fmt.Errorf("la sala con ID %s no existe", idSala)
	}
	cola := sala.HistoricoMensajes
	// Crear el mensaje
	nuevoMensaje := entities.Mensaje{
		Id:       uuid.New().String(),
		Nickname: nickname,
		Mensaje:  mensaje,
	}

	// Añadir el mensaje a la cola
	cola.Enqueue(nuevoMensaje)

	return nil
}

// Función para obtener mensajes de una sala
func (sm *GestionSalas) ObtenerMensajesDesdeId(idSala uuid.UUID, idMensaje string) ([]entities.Mensaje, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sala, err := sm.ObtenerSalaPorID(idSala)
	if err != nil {
		return nil, fmt.Errorf("la sala con ID %s no existe", idSala)
	}

	cola := sala.HistoricoMensajes
	// Obtener los mensajes desde la cola
	mensajes, err := cola.ObtenerMensajesDesdeId(idSala, idMensaje)
	if err != nil {
		return nil, fmt.Errorf(" Error en  ObtenerPorUltimoIdMensaje: %v", err)
	}

	return mensajes, nil
}

// Función para obtener mensajes de una sala
func (sm *GestionSalas) ObtenerMensajes(idSala uuid.UUID) ([]entities.Mensaje, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sala, err := sm.ObtenerSalaPorID(idSala)
	if err != nil {
		return nil, fmt.Errorf("la sala con ID %s no existe", idSala)
	}
	cola := sala.HistoricoMensajes
	// Obtener los mensajes desde la cola
	mensajes, err := cola.ObtenerTodos(idSala)
	if err != nil {
		return nil, fmt.Errorf(" Error en  ObtenerMensajes: %v", err)
	}
	return mensajes, nil
}

// Función para obtener mensajes de una sala
func (sm *GestionSalas) ObtenerMensajesCantidad(idSala uuid.UUID, cantidad int) ([]entities.Mensaje, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sala, err := sm.ObtenerSalaPorID(idSala)
	if err != nil {
		return nil, fmt.Errorf("la sala con ID %s no existe", idSala)
	}
	cola := sala.HistoricoMensajes
	// Obtener los mensajes desde la cola
	mensajes, err := cola.ObtenerMensajes(idSala, cantidad)
	if err != nil {
		return nil, fmt.Errorf(" Error en  ObtenerMensajes: %v", err)
	}
	return mensajes, nil
}

// Función para obtener la sala en la que un usuario debería unirse
func (sm *GestionSalas) UnirseASala(usuarioId string, salaFijaID string) (entities.Sala, error) {
	if salaFijaID != "" {
		// Si hay una sala fija, el usuario se une a esa sala
		salaUUID, err := uuid.Parse(salaFijaID)
		if err != nil {
			return entities.Sala{}, fmt.Errorf("el ID de la sala fija %s no es un UUID válido: %v", salaFijaID, err)
		}
		sala, existe := sm.SalasFijas[salaUUID]
		if !existe {
			return entities.Sala{}, fmt.Errorf("la sala fija con ID %s no existe", salaFijaID)
		}
		return sala, nil
	}

	// Si no hay sala fija, el usuario se unirá a la sala principal
	return sm.SalaPrincipal, nil
}

// ObtenerTodosLosMensajes combina los mensajes de la base de datos y la cola, sin duplicados
func (gestion *GestionSalas) ObtenerTodosLosMensajes(idSala uuid.UUID) ([]entities.Mensaje, error) {
	// Obtener la persistencia de MongoDB
	if gestion.persistencia == nil {
		return nil, fmt.Errorf("persistencia no inicializada")
	}

	// Obtener los mensajes desde la base de datos
	mensajesBD, err := gestion.persistencia.ObtenerMensajesDesdeSala(idSala)
	if err != nil {
		return nil, fmt.Errorf("error al obtener mensajes de la base de datos: %v", err)
	}
	sala, err := gestion.ObtenerSalaPorID(idSala)
	if err != nil {
		return nil, fmt.Errorf("la sala con ID %s no existe", idSala)
	}
	// Simulación: Obtener mensajes desde la cola
	mensajesCola := sala.HistoricoMensajes.ObtenerElementos()

	// Combinar mensajes, eliminando duplicados
	uniqueMessages := make(map[string]entities.Mensaje) // Usamos id_mensaje como clave
	// Primero, agregar los mensajes de la cola
	for _, mensaje := range mensajesCola {
		uniqueMessages[mensaje.Id] = mensaje
	}

	// Luego, agregar los mensajes de la base de datos
	for _, mensaje := range mensajesBD {
		uniqueMessages[mensaje.Id] = mensaje
	}
	// Convertir el mapa a un slice
	var mensajesUnicos []entities.Mensaje
	for _, mensaje := range uniqueMessages {
		mensajesUnicos = append(mensajesUnicos, mensaje)
	}

	return mensajesUnicos, nil
}

// Simulación: Obtener mensajes desde la cola (reemplazar con tu lógica de cola)
func (gestion *GestionSalas) ObtenerMensajesDesdeCola(idSala uuid.UUID) ([]entities.Mensaje, error) {

	// Obtener la cola correspondiente a la sala
	gestionSalas := &GestionSalas{} // Referencia a GestionSalas (ajustar según tu implementación)
	gestionSalas.mu.Lock()
	sala, err := gestion.ObtenerSalaPorID(idSala)
	if err != nil {
		return nil, fmt.Errorf("la sala con ID %s no existe", idSala)
	}
	mensajesCola := sala.HistoricoMensajes.ObtenerElementos()
	return mensajesCola, nil
}
