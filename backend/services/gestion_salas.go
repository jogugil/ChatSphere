package services

import (
	"backend/models"
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
	mu            sync.RWMutex                 // Mutex lectura/excritura
	persistencia  entities.Persistencia       // Persistencia para almacenar los datos
	once          sync.Once                   // Para garantizar inicialización única por instancia
}

// Instancia única de GestionSalas
var instancia *GestionSalas

func NuevaGestionSalas(persistencia entities.Persistencia, configFile string) *GestionSalas {
	log.Println("NuevaGestionSalas: Iniciando la creación de una nueva instancia.")
	if instancia == nil {
		instancia = &GestionSalas{
			SalasFijas:   make(map[uuid.UUID]entities.Sala),
			persistencia: persistencia,
		}
	}
	instancia.once.Do(func() {
		log.Println("NuevaGestionSalas: Cargando configuración desde archivo.")
		if configFile != "" {
			err := instancia.CargarSalasFijasDesdeArchivo(configFile)
			if err != nil {
				log.Fatalf("NuevaGestionSalas: Error al cargar configuración: %v", err)
			}
		}
	})
	log.Println("NuevaGestionSalas: Creación completada.")
	return instancia
}

func (sm *GestionSalas) CargarSalasFijasDesdeArchivo(configFile string) error {
	log.Printf("CargarSalasFijasDesdeArchivo: Cargando salas desde archivo: %s", configFile)
	if sm.SalasFijas == nil {
		sm.SalasFijas = make(map[uuid.UUID]entities.Sala)
	}

	sm.SalaPrincipal = entities.Sala{
		ID:                uuid.New(),
		Name:              "Sala Principal",
		SalaTipo:          "Principal",
		HistoricoMensajes: entities.NuevaColaCircular(sm.persistencia),
	}

	file, err := os.Open(configFile)
	if err != nil {
		log.Printf("CargarSalasFijasDesdeArchivo: Error al abrir archivo: %v", err)
		return err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("CargarSalasFijasDesdeArchivo: Error al leer archivo: %v", err)
		return err
	}

	var salas []map[string]interface{}
	err = json.Unmarshal(byteValue, &salas)
	if err != nil {
		log.Printf("CargarSalasFijasDesdeArchivo: Error al parsear JSON: %v", err)
		return err
	}

	// Usar RLock para leer de SalasFijas de forma concurrente.
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for _, salaData := range salas {
		salaID, err := uuid.Parse(salaData["id"].(string))
		if err != nil {
			log.Printf("CargarSalasFijasDesdeArchivo: Error al parsear ID: %v", err)
			return err
		}

		sala := entities.Sala{
			ID:                salaID,
			Name:              salaData["nombre"].(string),
			SalaTipo:          "Fija",
			HistoricoMensajes: entities.NuevaColaCircular(sm.persistencia),
		}
		sm.SalasFijas[salaID] = sala
	}
	log.Println("CargarSalasFijasDesdeArchivo: Carga completada.")
	return nil
}

func (sm *GestionSalas) CrearSalaTemporal(nombre string) entities.Sala {
	log.Printf("CrearSalaTemporal: Creando sala temporal con nombre: %s", nombre)
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.persistencia == nil {
		log.Fatal("CrearSalaTemporal: Persistencia no inicializada.")
	}

	sala := entities.Sala{
		ID:                uuid.New(),
		Name:              nombre,
		SalaTipo:          "Temporal",
		HistoricoMensajes: entities.NuevaColaCircular(sm.persistencia),
	}
	log.Printf("CrearSalaTemporal: Sala temporal creada con ID: %s", sala.ID)
	return sala
}

func (sm *GestionSalas) ObtenerSalaPorID(idSala uuid.UUID) (*entities.Sala, error) {
	log.Printf("ObtenerSalaPorID: Buscando sala con ID: %s", idSala)
	 
	sm.mu.RLock()
	log.Printf("Lock adquirido para lectura de la sala con ID %s", idSala)
	defer func() {
		sm.mu.RUnlock()
		log.Printf("Lock liberado para lectura de la sala con ID %s", idSala)
	}()

	if sm.SalaPrincipal.ID == idSala {
		log.Println("ObtenerSalaPorID: Sala principal encontrada.")
		return &sm.SalaPrincipal, nil
	}

	if sala, existe := sm.SalasFijas[idSala]; existe {
		log.Printf("ObtenerSalaPorID: Sala fija encontrada con ID: %s", idSala)
		return &sala, nil
	}

	log.Printf("ObtenerSalaPorID: Sala no encontrada con ID: %s", idSala)
	return nil, fmt.Errorf("la sala con ID %s no existe", idSala)
}

func (sm *GestionSalas) EnviarMensaje(idSala uuid.UUID, nickname, mensaje string, usuario entities.Usuario) error {
    log.Printf("EnviarMensaje: Enviando mensaje a la sala con ID: %s", idSala)

    // Primero obtenemos la sala con RLock, ya que estamos haciendo una lectura
    sala, err := sm.ObtenerSalaPorID(idSala)
    if err != nil {
        log.Printf("EnviarMensaje: Error al obtener sala: %v", err)
        return err
    }
 
    nuevoMensaje := models.CrearMensajeConFecha (mensaje,usuario,idSala,usuario.HoraUltimaAccion)
  
    // Modificamos el historial de mensajes
    sala.HistoricoMensajes.Enqueue(*nuevoMensaje)
    log.Printf("EnviarMensaje: Mensaje enviado a la sala con ID: %s", idSala)
    return nil
}

// Función para obtener mensajes de una sala
func (sm *GestionSalas) ObtenerMensajesDesdeId(idSala uuid.UUID, idMensaje uuid.UUID) ([]entities.Mensaje, error) {
	log.Printf("ObtenerMensajesDesdeId: Obteniendo mensajes desde ID %s en la sala %s", idMensaje, idSala)
	sm.mu.RLock() // Lectura, se puede hacer concurrentemente
	defer sm.mu.RUnlock()

	sala, err := sm.ObtenerSalaPorID(idSala)
	if err != nil {
		log.Printf("ObtenerMensajesDesdeId: Error al obtener sala: %v", err)
		return nil, fmt.Errorf("la sala con ID %s no existe", idSala)
	}

	cola := sala.HistoricoMensajes
	mensajes, err := cola.ObtenerMensajesDesdeId(idSala, idMensaje)
	if err != nil {
		log.Printf("ObtenerMensajesDesdeId: Error al obtener mensajes: %v", err)
		return nil, fmt.Errorf("Error en ObtenerPorUltimoIdMensaje: %v", err)
	}

	log.Printf("ObtenerMensajesDesdeId: Mensajes obtenidos correctamente (%d mensajes)", len(mensajes))
	return mensajes, nil
}

// Función para obtener todos los mensajes de una sala
func (sm *GestionSalas) ObtenerMensajes(idSala uuid.UUID) ([]entities.Mensaje, error) {
	log.Printf("ObtenerMensajes: Obteniendo todos los mensajes de la sala %s", idSala)
	sm.mu.RLock() // Lectura, se puede hacer concurrentemente
	defer sm.mu.RUnlock()

	sala, err := sm.ObtenerSalaPorID(idSala)
	if err != nil {
		log.Printf("ObtenerMensajes: Error al obtener sala: %v", err)
		return nil, fmt.Errorf("la sala con ID %s no existe", idSala)
	}
	cola := sala.HistoricoMensajes
	mensajes, err := cola.ObtenerTodos(idSala)
	if err != nil {
		log.Printf("ObtenerMensajes: Error al obtener mensajes: %v", err)
		return nil, fmt.Errorf("Error en ObtenerMensajes: %v", err)
	}

	log.Printf("ObtenerMensajes: Mensajes obtenidos correctamente (%d mensajes)", len(mensajes))
	return mensajes, nil
}


func (sm *GestionSalas) ObtenerMensajesCantidad(idSala uuid.UUID, cantidad int) ([]entities.Mensaje, error) {
	log.Printf("ObtenerMensajesCantidad: Obteniendo %d mensajes de la sala %s", cantidad, idSala)
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Validación previa
	sala, err := sm.ObtenerSalaPorID(idSala)
	if err != nil {
		log.Printf("ObtenerMensajesCantidad: Error al obtener sala: %v", err)
		return nil, fmt.Errorf("la sala con ID %s no existe", idSala)
	}

	// Mejorar rendimiento si la cantidad es mayor que los mensajes disponibles
	cola := sala.HistoricoMensajes
	mensajes, err := cola.ObtenerMensajes(idSala, cantidad)
	if err != nil {
		log.Printf("ObtenerMensajesCantidad: Error al obtener mensajes: %v", err)
		return nil, fmt.Errorf("Error en ObtenerMensajes: %v", err)
	}

	log.Printf("ObtenerMensajesCantidad: Mensajes obtenidos correctamente (%d mensajes)", len(mensajes))
	return mensajes, nil
}

func (sm *GestionSalas) UnirseASala(usuarioId string, salaFijaID string) (entities.Sala, error) {
	log.Printf("UnirseASala: Usuario %s intentando unirse a la sala fija %s", usuarioId, salaFijaID)

	if salaFijaID != "" {
		salaUUID, err := uuid.Parse(salaFijaID)
		if err != nil {
			log.Printf("UnirseASala: Error al parsear ID de sala fija: %v", err)
			return entities.Sala{}, fmt.Errorf("el ID de la sala fija %s no es un UUID válido: %v", salaFijaID, err)
		}

		// Mejorar la búsqueda con un "if" temprano
		sala, existe := sm.SalasFijas[salaUUID]
		if !existe {
			log.Printf("UnirseASala: Sala fija no encontrada con ID: %s", salaFijaID)
			return entities.Sala{}, fmt.Errorf("la sala fija con ID %s no existe", salaFijaID)
		}

		log.Printf("UnirseASala: Usuario %s unido a la sala fija %s", usuarioId, salaFijaID)
		return sala, nil
	}

	log.Printf("UnirseASala: Usuario %s unido a la sala principal", usuarioId)
	return sm.SalaPrincipal, nil
}

func (gestion *GestionSalas) ObtenerTodosLosMensajes(idSala uuid.UUID) ([]entities.Mensaje, error) {
	log.Printf("ObtenerTodosLosMensajes: Obteniendo todos los mensajes únicos de la sala %s", idSala)

	if gestion.persistencia == nil {
		log.Println("ObtenerTodosLosMensajes: Persistencia no inicializada")
		return nil, fmt.Errorf("persistencia no inicializada")
	}

	mensajesBD, err := gestion.persistencia.ObtenerMensajesDesdeSala(idSala)
	if err != nil {
		log.Printf("ObtenerTodosLosMensajes: Error al obtener mensajes de la base de datos: %v", err)
		return nil, fmt.Errorf("error al obtener mensajes de la base de datos: %v", err)
	}

	sala, err := gestion.ObtenerSalaPorID(idSala)
	if err != nil {
		log.Printf("ObtenerTodosLosMensajes: Error al obtener sala: %v", err)
		return nil, fmt.Errorf("la sala con ID %s no existe", idSala)
	}

	mensajesCola := sala.HistoricoMensajes.ObtenerElementos()
	uniqueMessages := make(map[uuid.UUID]entities.Mensaje)

	// Usando un map para eliminar duplicados
	for _, mensaje := range mensajesCola {
		uniqueMessages[mensaje.IDM] = mensaje
	}
	for _, mensaje := range mensajesBD {
		uniqueMessages[mensaje.IDM] = mensaje
	}

	var mensajesUnicos []entities.Mensaje
	for _, mensaje := range uniqueMessages {
		mensajesUnicos = append(mensajesUnicos, mensaje)
	}

	log.Printf("ObtenerTodosLosMensajes: Mensajes únicos obtenidos (%d mensajes)", len(mensajesUnicos))
	return mensajesUnicos, nil
}

func (gestion *GestionSalas) ObtenerMensajesDesdeCola(idSala uuid.UUID) ([]entities.Mensaje, error) {
	log.Printf("ObtenerMensajesDesdeCola: Obteniendo mensajes desde la cola de la sala %s", idSala)

	sala, err := gestion.ObtenerSalaPorID(idSala)
	if err != nil {
		log.Printf("ObtenerMensajesDesdeCola: Error al obtener sala: %v", err)
		return nil, fmt.Errorf("la sala con ID %s no existe", idSala)
	}

	mensajesCola := sala.HistoricoMensajes.ObtenerElementos()
	log.Printf("ObtenerMensajesDesdeCola: Mensajes obtenidos desde la cola (%d mensajes)", len(mensajesCola))
	return mensajesCola, nil
}