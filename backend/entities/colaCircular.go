package entities

import (
	"backend/utils"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/google/uuid"
)

type ColaCircular struct {
	Buffer       []Mensaje    // Buffer circular de mensajes
	Head         int          // Índice del primer mensaje
	Tail         int          // Índice del siguiente espacio disponible
	Size         int          // Número actual de mensajes en la cola
	Capacity     int          // Capacidad máxima de la cola (se lee desde el archivo .env)
	Mu           sync.Mutex   // Mutex para la exclusión mutua
	Persistencia Persistencia // Para almacenar en la BD mongoDB
}

// NuevaColaCircular crea una nueva instancia de ColaCircular
func NuevaColaCircular(persistencia Persistencia) *ColaCircular {
	// Leer el tamaño máximo de la cola desde las variables de entorno
	val, err_env := utils.ObtenerVariableDeEntorno("SizeQueue")
	if err_env != nil {
		log.Fatalf("Error al leer EXP_SIZE_QMESSAGE: %v", err_env)
	}
	capacity, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("Error al strconv.Atoi (val): %v", err)
	}

	// Crear la cola con la capacidad configurada
	return &ColaCircular{
		Buffer:       make([]Mensaje, capacity),
		Head:         0,
		Tail:         0,
		Size:         0,
		Capacity:     capacity,
		Persistencia: persistencia,
	}
}

// Enqueue agrega un mensaje a la cola
func (q *ColaCircular) Enqueue(msg Mensaje) {
	q.Mu.Lock()
	defer q.Mu.Unlock()

	// Si la cola está llena, vacía los mensajes antes de agregar uno nuevo
	if q.Size == q.Capacity {
		log.Println("La cola está llena, procesando y vaciando en la base de datos.")
		q.flushToDatabase()
	}

	q.Buffer[q.Tail] = msg
	q.Tail = (q.Tail + 1) % q.Capacity
	q.Size++

	log.Printf("Mensaje añadido a la cola: %v", msg)
}

// flushToDatabase vacía la cola y guarda los mensajes en la base de datos
func (q *ColaCircular) flushToDatabase() {
	// Proceso de vaciado a base de datos
	var mensajes []Mensaje
	for i := 0; i < q.Size; i++ {
		mensajes = append(mensajes, q.Buffer[(q.Head+i)%q.Capacity])
	}
	// Guardar mensajes en la base de datos
	err := q.Persistencia.GuardarMensajesEnBaseDeDatos(mensajes)
	if err != nil {
		log.Printf("Error al guardar los mensajes: %v", err)
		return
	}
	q.Size = 0 // Ahora que se vació la cola, podemos resetear el tamaño
}

// ObtenerMensajes obtiene hasta 'cantidad' de mensajes de la cola circular sin eliminarlos.
func (q *ColaCircular) ObtenerMensajes(idSala uuid.UUID, cantidad int) ([]Mensaje, error) {
	q.Mu.Lock()
	defer q.Mu.Unlock()

	// Si la cola está vacía, retornamos un error
	if q.Size == 0 {
		mensajes, err := q.Persistencia.ObtenerMensajesDesdeSala(idSala)
		if err != nil {
			return nil, fmt.Errorf("la cola está vacía")
		}
		return mensajes, nil
	}

	// Si la cantidad solicitada es mayor que el tamaño de la cola, ajustamos la cantidad
	if cantidad > q.Size {
		cantidad = q.Size
	}

	// Creamos un slice para los mensajes a devolver
	var mensajes []Mensaje
	for i := 0; i < cantidad; i++ {
		// Calculamos el índice en el buffer circular
		mensaje := q.Buffer[(q.Head+i)%q.Capacity]
		mensajes = append(mensajes, mensaje)
	}

	// Devolvemos los mensajes obtenidos
	return mensajes, nil
}

// ObtenerTodos obtiene todos los mensajes de la cola circular sin eliminarlos.
func (q *ColaCircular) ObtenerTodos(idSala uuid.UUID) ([]Mensaje, error) {
	q.Mu.Lock()
	defer q.Mu.Unlock()

	// Si la cola está vacía, retornamos un error
	if q.Size == 0 {
		mensajes, err := q.Persistencia.ObtenerMensajesDesdeSala(idSala)
		if err != nil {
			return nil, fmt.Errorf("la cola está vacía y hay un error en la BD")
		}
		return mensajes, nil
	}

	// Creamos un slice para los mensajes a devolver
	var mensajes []Mensaje

	// Recorremos todos los mensajes en la cola circular
	for i := 0; i < q.Size; i++ {
		// Calculamos el índice en el buffer circular
		mensaje := q.Buffer[(q.Head+i)%q.Capacity]
		mensajes = append(mensajes, mensaje)
	}

	// Devolvemos todos los mensajes obtenidos
	return mensajes, nil
}

// ObtenerMensajesDesdeId obtiene mensajes desde un id específico.
func (q *ColaCircular) ObtenerMensajesDesdeId(idSala uuid.UUID, idMensaje uuid.UUID) ([]Mensaje, error) {
	q.Mu.Lock()
	defer q.Mu.Unlock()

	// Si la cola está vacía, retornamos un error
	if q.Size == 0 {
		return nil, fmt.Errorf("la cola está vacía")
	}

	// Buscamos el índice del mensaje con el id dado
	var index int = -1
	for i := 0; i < q.Size; i++ {
		mensaje := q.Buffer[(q.Head+i)%q.Capacity]
		if  mensaje.IDM == idMensaje { // Asumimos que `Mensaje` tiene un campo `ID`
			index = (q.Head + i) % q.Capacity
			break
		}
	}

	// Si no encontramos el mensaje, devolvemos un error
	// Ojo! el mensaje puede estar volcado en la BD mongoDB
	if index == -1 {
		mensajes, err := q.Persistencia.ObtenerMensajesDesdeId(idSala, idMensaje)
		if err != nil {
			return nil, fmt.Errorf("mensaje con id %s no encontrado", idMensaje)
		}
		return mensajes, nil
	}

	// Creamos un slice para los mensajes a devolver
	var mensajes []Mensaje

	// Recopilamos los mensajes desde el índice encontrado hasta el final de la cola
	for i := 0; i < q.Size; i++ {
		mensaje := q.Buffer[(index+i)%q.Capacity]
		mensajes = append(mensajes, mensaje)
	}

	return mensajes, nil
}

// ObtenerElementos devuelve un slice con los mensajes almacenados en la cola.
func (c *ColaCircular) ObtenerElementos() []Mensaje {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	// Si la cola está vacía, devolvemos un slice vacío.
	if c.Size == 0 {
		return []Mensaje{}
	}

	// Crear un slice para almacenar los mensajes en el orden correcto.
	elementos := make([]Mensaje, 0, c.Size)

	// Recorrer desde `head` hasta el final, manejando la circularidad.
	for i := 0; i < c.Size; i++ {
		index := (c.Head + i) % c.Capacity
		elementos = append(elementos, c.Buffer[index])
	}

	return elementos
}
