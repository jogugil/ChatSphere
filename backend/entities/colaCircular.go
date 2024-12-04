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
	buffer       []Mensaje     // Buffer circular de mensajes
	head         int           // Índice del primer mensaje
	tail         int           // Índice del siguiente espacio disponible
	size         int           // Número actual de mensajes en la cola
	capacity     int           // Capacidad máxima de la cola (se lee desde el archivo .env)
	mu           sync.Mutex    // Mutex para la exclusión mutua
	persistencia *Persistencia // Para almacenar en la BD mongoDB
}

// NuevaColaCircular crea una nueva instancia de ColaCircular
func NuevaColaCircular(persistence *Persistencia) *ColaCircular {
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
		buffer:       make([]Mensaje, capacity),
		head:         0,
		tail:         0,
		size:         0,
		capacity:     capacity,
		persistencia: persistence,
	}
}

// Enqueue agrega un mensaje a la cola
func (q *ColaCircular) Enqueue(msg Mensaje) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Si la cola está llena, vacía los mensajes antes de agregar uno nuevo
	if q.size == q.capacity {
		log.Println("La cola está llena, procesando y vaciando en la base de datos.")
		q.flushToDatabase()
	}

	q.buffer[q.tail] = msg
	q.tail = (q.tail + 1) % q.capacity
	q.size++

	log.Printf("Mensaje añadido a la cola: %v", msg)
}

// flushToDatabase vacía la cola y guarda los mensajes en la base de datos
func (q *ColaCircular) flushToDatabase() {
	// Proceso de vaciado a base de datos
	var mensajes []Mensaje
	for i := 0; i < q.size; i++ {
		mensajes = append(mensajes, q.buffer[(q.head+i)%q.capacity])
	}
	// Guardar mensajes en la base de datos
	err := (*q.persistencia).GuardarMensajesEnBaseDeDatos(mensajes)
	if err != nil {
		log.Printf("Error al guardar los mensajes: %v", err)
		return
	}
	q.size = 0 // Ahora que se vació la cola, podemos resetear el tamaño
}

// ObtenerMensajes obtiene hasta 'cantidad' de mensajes de la cola circular sin eliminarlos.
func (q *ColaCircular) ObtenerMensajes(idSala uuid.UUID, cantidad int) ([]Mensaje, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Si la cola está vacía, retornamos un error
	if q.size == 0 {
		mensajes, err := (*q.persistencia).ObtenerMensajesDesdeSala(idSala)
		if err != nil {
			return nil, fmt.Errorf("la cola está vacía")
		}
		return mensajes, nil
	}

	// Si la cantidad solicitada es mayor que el tamaño de la cola, ajustamos la cantidad
	if cantidad > q.size {
		cantidad = q.size
	}

	// Creamos un slice para los mensajes a devolver
	var mensajes []Mensaje
	for i := 0; i < cantidad; i++ {
		// Calculamos el índice en el buffer circular
		mensaje := q.buffer[(q.head+i)%q.capacity]
		mensajes = append(mensajes, mensaje)
	}

	// Devolvemos los mensajes obtenidos
	return mensajes, nil
}

// ObtenerTodos obtiene todos los mensajes de la cola circular sin eliminarlos.
func (q *ColaCircular) ObtenerTodos(idSala uuid.UUID) ([]Mensaje, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Si la cola está vacía, retornamos un error
	if q.size == 0 {
		mensajes, err := (*q.persistencia).ObtenerMensajesDesdeSala(idSala)
		if err != nil {
			return nil, fmt.Errorf("la cola está vacía y hay un error en la BD")
		}
		return mensajes, nil
	}

	// Creamos un slice para los mensajes a devolver
	var mensajes []Mensaje

	// Recorremos todos los mensajes en la cola circular
	for i := 0; i < q.size; i++ {
		// Calculamos el índice en el buffer circular
		mensaje := q.buffer[(q.head+i)%q.capacity]
		mensajes = append(mensajes, mensaje)
	}

	// Devolvemos todos los mensajes obtenidos
	return mensajes, nil
}

// ObtenerMensajesDesdeId obtiene mensajes desde un id específico.
func (q *ColaCircular) ObtenerMensajesDesdeId(idSala uuid.UUID, idMensaje uuid.UUID) ([]Mensaje, error) {
	// Log para ver qué parámetros entran en la función
	fmt.Printf("ObtenerMensajesDesdeId - Parámetros de entrada: idSala=%s, idMensaje=%s\n", idSala, idMensaje)

	q.mu.Lock()
	defer q.mu.Unlock()

	// Si la cola está vacía, intentamos obtener mensajes desde la persistencia
	if q.size == 0 {
		return q.obtenerMensajesDesdePersistencia(idSala, idMensaje)
	}

	// Si idMensaje es uuid.Nil, devolver todos los mensajes de la cola
	if idMensaje == uuid.Nil {
		return q.obtenerTodosLosMensajes()
	}

	// Buscar el mensaje en la cola circular
	mensaje, index, err := q.buscarMensajeEnCola(idMensaje)
	if err != nil {
		// Si no lo encontramos, intentar obtenerlo desde la persistencia
		return q.obtenerMensajeDesdePersistencia(idSala, idMensaje)
	}

	// Recopilar los mensajes desde el índice encontrado hasta el final de la cola
	return q.recopilarMensajesDesdeIndice(index, mensaje)
}

// Función para obtener mensajes desde la persistencia
func (q *ColaCircular) obtenerMensajesDesdePersistencia(idSala uuid.UUID, idMensaje uuid.UUID) ([]Mensaje, error) {
	if idMensaje == uuid.Nil {
		fmt.Println("idMensaje es UUID de ceros, obteniendo todos los mensajes de la sala...")
		return (*q.persistencia).ObtenerMensajesDesdeSala(idSala)
	}
	fmt.Println("idMensaje no es UUID de ceros, obteniendo el mensaje con el id especificado...")
	return (*q.persistencia).ObtenerMensajesDesdeId(idSala, idMensaje)
}

// Función para obtener todos los mensajes de la cola
func (q *ColaCircular) obtenerTodosLosMensajes() ([]Mensaje, error) {
	fmt.Println("idMensaje es UUID de ceros, devolviendo todos los mensajes de la cola...")

	var mensajes []Mensaje
	for i := 0; i < q.size; i++ {
		mensaje := q.buffer[(q.head+i)%q.capacity]
		mensajes = append(mensajes, mensaje)
	}
	fmt.Printf("Todos los mensajes en la cola: %v\n", mensajes)
	return mensajes, nil
}

// Función para buscar un mensaje en la cola
func (q *ColaCircular) buscarMensajeEnCola(idMensaje uuid.UUID) (Mensaje, int, error) {
	fmt.Println("Buscando el mensaje en la cola circular...")

	for i := 0; i < q.size; i++ {
		mensaje := q.buffer[(q.head+i)%q.capacity]
		if mensaje.IDM == idMensaje {
			fmt.Printf("Mensaje encontrado en el índice: %d\n", i)
			return mensaje, i, nil
		}
	}

	// Si no se encuentra el mensaje
	fmt.Println("Mensaje no encontrado en la cola.")
	return Mensaje{}, -1, fmt.Errorf("mensaje con id %s no encontrado en la cola", idMensaje)
}

// Función para obtener el mensaje desde la persistencia
func (q *ColaCircular) obtenerMensajeDesdePersistencia(idSala uuid.UUID, idMensaje uuid.UUID) ([]Mensaje, error) {
	return q.obtenerMensajesDesdePersistencia(idSala, idMensaje)
}

// Función para recopilar mensajes desde un índice
func (q *ColaCircular) recopilarMensajesDesdeIndice(index int, mensaje Mensaje) ([]Mensaje, error) {
	var mensajes []Mensaje
	fmt.Println("Recopilando mensajes desde el índice encontrado hasta el final de la cola...")

	// Recopilamos el primer mensaje
	mensajes = append(mensajes, mensaje)

	// Recopilamos los mensajes desde el índice encontrado hasta el final de la cola
	for i := 1; i < q.size; i++ {
		mensaje := q.buffer[(index+i)%q.capacity]
		mensajes = append(mensajes, mensaje)
	}

	// Log para los mensajes que se devolverán
	fmt.Printf("Mensajes recopilados: %v\n", mensajes)

	return mensajes, nil
}
 

// ObtenerElementos devuelve un slice con los mensajes almacenados en la cola.
func (c *ColaCircular) ObtenerElementos() []Mensaje {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Si la cola está vacía, devolvemos un slice vacío.
	if c.size == 0 {
		return []Mensaje{}
	}

	// Crear un slice para almacenar los mensajes en el orden correcto.
	elementos := make([]Mensaje, 0, c.size)

	// Recorrer desde `head` hasta el final, manejando la circularidad.
	for i := 0; i < c.size; i++ {
		index := (c.head + i) % c.capacity
		elementos = append(elementos, c.buffer[index])
	}

	return elementos
}

// ComprobarYVisualizarMensajes verifica la cola y muestra los mensajes si existen.
func (cola *ColaCircular) ComprobarYVisualizarMensajes() (bool, []Mensaje, string) {
	// Verificar si la cola es nula
	if cola == nil {
		log.Println("Error: La cola es nula.")
		return false, nil, "La cola es nula."
	}

	cola.mu.Lock()
	defer cola.mu.Unlock()

	// Verificar si la cola tiene mensajes
	if cola.size == 0 {
		log.Println("Error: La cola está vacía.")
		return false, nil, "La cola está vacía."
	}

	// Obtener y visualizar los mensajes
	mensajes := cola.ObtenerElementos()
	log.Println("La cola contiene los siguientes mensajes:")
	for i, mensaje := range mensajes {
		log.Printf("Mensaje %d: %+v\n", i+1, mensaje)
	}

	return true, mensajes, "La cola contiene mensajes."
}
