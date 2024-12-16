package persistence

import (
	"backend/entities"
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Configuración de MongoDB (puede ser configurada en la propia estructura de MongoPersistencia)

var mongoInstance *MongoPersistencia
var onceMongodb sync.Once

type MongoPersistencia struct {
	client *mongo.Client
	db     *mongo.Database
}

// NuevaMongoPersistencia crea o devuelve una instancia de MongoPersistencia con pool de conexiones
func NuevaMongoPersistencia(uri, dbName string) (*entities.Persistencia, error) {
	fmt.Printf("Creando nueva instancia de MongoPersistencia con URI: %s y DB: %s\n", uri, dbName)
	onceMongodb.Do(func() {
		// Configura el cliente MongoDB con un pool de conexiones (máximo 10 conexiones).
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetMaxPoolSize(20))
		if err != nil {
			fmt.Printf("Error de conexión con MongoDB: %v\n", err)
			return
		}
		// Verifica la conexión
		err = client.Ping(context.Background(), nil)
		if err != nil {
			fmt.Printf("Error al hacer ping a MongoDB: %v\n", err)
			return
		}
		db := client.Database(dbName)
		mongoInstance = &MongoPersistencia{client: client, db: db}
	})

	if mongoInstance == nil {
		return nil, fmt.Errorf("no se pudo crear la instancia de MongoPersistencia")
	}

	var persistenciaPersist entities.Persistencia = mongoInstance
	fmt.Println("Instancia de MongoPersistencia creada correctamente")
	return &persistenciaPersist, nil
}

// ObtenerInstanciaDB obtiene la instancia de la base de datos
func ObtenerInstanciaDB() (*entities.Persistencia, error) {
	fmt.Println("ObtenerInstanciaDB: Obteniendo instancia de base de datos MongoDB...")
	if mongoInstance == nil {
		// Regresamos un error si la instancia no ha sido inicializada
		return nil, errors.New("ObtenerInstanciaDB: la instancia de MongoPersistencia no ha sido inicializada")
	}

	// Devolvemos mongoInstance como un tipo Persistencia, que implementa la interfaz Persistencia
	var persistenciaPersist entities.Persistencia = mongoInstance
	return &persistenciaPersist, nil
}

// Implementa GuardarSala
func (m *MongoPersistencia) GuardarSala(sala entities.Room) error {
	fmt.Printf("Guardando sala: %+v\n", sala)
	collection := m.db.Collection("salas")
	_, err := collection.InsertOne(context.TODO(), sala)
	if err != nil {
		fmt.Printf("Error al guardar la sala: %v\n", err)
	}
	return err
}

// Implementa ObtenerSala
func (m *MongoPersistencia) ObtenerSala(id uuid.UUID) (entities.Room, error) {
	fmt.Printf("Obteniendo sala con ID: %s\n", id)
	var sala entities.Room
	collection := m.db.Collection("salas")
	err := collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&sala)
	if err != nil {
		fmt.Printf("Error al obtener la sala: %v\n", err)
		return entities.Room{}, err
	}
	fmt.Printf("Sala obtenida: %+v\n", sala)
	return sala, nil
}

// GuardarMensaje guarda un solo mensaje en la base de datos MongoDB
func (mp *MongoPersistencia) GuardarMensaje(mensaje *entities.Message) error {
	fmt.Printf("Guardando mensaje: %+v\n", mensaje)
	collection := mp.db.Collection("mensajes")
	_, err := collection.InsertOne(context.TODO(), mensaje)
	if err != nil {
		fmt.Printf("Error al guardar el mensaje en MongoDB: %v\n", err)
		return fmt.Errorf("error al guardar el mensaje en MongoDB: %v", err)
	}
	fmt.Println("Mensaje guardado correctamente")
	return nil
}

// GuardarMensajesEnBaseDeDatos guarda una lista de mensajes en MongoDB
func (mp *MongoPersistencia) GuardarMensajesEnBaseDeDatos(mensajes []entities.Message) error {
	fmt.Printf("Guardando lista de %d mensajes\n", len(mensajes))
	collection := mp.db.Collection("mensajes")

	// Convertimos los mensajes a una lista de interfaces{}
	var documentos []interface{}
	for _, mensaje := range mensajes {
		documentos = append(documentos, bson.D{
			{Key: "nick_usuario", Value: mensaje.Nickname},
			{Key: "fecha_envio", Value: mensaje.SendDate},
			{Key: "id_sala", Value: mensaje.RoomID},
			{Key: "nombre_sala", Value: mensaje.RoomName},
		})
	}

	// Insertamos todos los documentos en la colección
	_, err := collection.InsertMany(context.TODO(), documentos)
	if err != nil {
		fmt.Printf("Error al guardar los mensajes en MongoDB: %v\n", err)
		return fmt.Errorf("error al guardar los mensajes en MongoDB: %v", err)
	}

	// Todo salió bien
	fmt.Println("Mensajes guardados correctamente")
	return nil
}

func (mp *MongoPersistencia) ObtenerMensajesDesdeSala(idSala uuid.UUID) ([]entities.Message, error) {
	fmt.Printf("Obteniendo mensajes desde la sala con ID: %s\n", idSala)
	collection := mp.db.Collection("mensajes")

	// Buscar el mensaje base
	var mensajeBase entities.Message
	filterBase := bson.D{{Key: "id_sala", Value: idSala}}
	err := collection.FindOne(context.TODO(), filterBase).Decode(&mensajeBase)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("No se encontraron mensajes en la sala: %v\n", idSala)
			return nil, fmt.Errorf("no se  encontraron mensajes en la sala: %v", idSala)
		}
		fmt.Printf("Error al obtener el mensaje base: %v\n", err)
		return nil, fmt.Errorf("error al obtener el mensaje base: %v", err)
	}

	cursor, err := collection.Find(context.TODO(), filterBase)
	if err != nil {
		fmt.Printf("Error al ejecutar la consulta: %v\n", err)
		return nil, fmt.Errorf("error al ejecutar la consulta: %v", err)
	}
	defer cursor.Close(context.TODO())

	var mensajes []entities.Message
	for cursor.Next(context.TODO()) {
		var mensaje entities.Message
		if err := cursor.Decode(&mensaje); err != nil {
			fmt.Printf("Error al decodificar el mensaje: %v\n", err)
			return nil, fmt.Errorf("error al decodificar el mensaje: %v", err)
		}
		mensajes = append(mensajes, mensaje)
	}

	if err := cursor.Err(); err != nil {
		fmt.Printf("Error al iterar sobre el cursor: %v\n", err)
		return nil, fmt.Errorf("error al iterar sobre el cursor: %v", err)
	}

	fmt.Printf("Mensajes obtenidos: %d mensajes\n", len(mensajes))
	return mensajes, nil
}

func (mp *MongoPersistencia) ObtenerMensajesDesdeId(idSala uuid.UUID, idMensaje uuid.UUID) ([]entities.Message, error) {
	fmt.Printf("Obteniendo mensajes desde la sala con ID: %s y mensaje con ID: %s\n", idSala, idMensaje)
	collection := mp.db.Collection("mensajes")

	// Buscar el mensaje base
	var mensajeBase entities.Message
	filterBase := bson.D{{Key: "id_mensaje", Value: idMensaje}, {Key: "id_sala", Value: idSala}}
	err := collection.FindOne(context.TODO(), filterBase).Decode(&mensajeBase)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("No se encontró el mensaje con id_mensaje: %v en la sala: %v\n", idMensaje, idSala)
			return nil, fmt.Errorf("no se encontró el mensaje con id_mensaje: %v en la sala: %v", idMensaje, idSala)
		}
		fmt.Printf("Error al obtener el mensaje base: %v\n", err)
		return nil, fmt.Errorf("error al obtener el mensaje base: %v", err)
	}

	// Filtro para mensajes posteriores
	filter := bson.D{
		{Key: "id_sala", Value: idSala},
		{Key: "fecha_envio", Value: bson.D{{Key: "$gte", Value: mensajeBase.SendDate}}}}
	findOptions := options.Find().SetSort(bson.D{{Key: "fecha_envio", Value: 1}})

	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		fmt.Printf("Error al ejecutar la consulta: %v\n", err)
		return nil, fmt.Errorf("error al ejecutar la consulta: %v", err)
	}
	defer cursor.Close(context.TODO())

	var mensajes []entities.Message
	for cursor.Next(context.TODO()) {
		var mensaje entities.Message
		if err := cursor.Decode(&mensaje); err != nil {
			fmt.Printf("Error al decodificar el mensaje: %v\n", err)
			return nil, fmt.Errorf("error al decodificar el mensaje: %v", err)
		}
		mensajes = append(mensajes, mensaje)
	}

	if err := cursor.Err(); err != nil {
		fmt.Printf("Error al iterar sobre el cursor: %v\n", err)
		return nil, fmt.Errorf("error al iterar sobre el cursor: %v", err)
	}

	fmt.Printf("Mensajes obtenidos después del ID de mensaje: %d\n", len(mensajes))
	return mensajes, nil
}
func (mp *MongoPersistencia) GuardarUsuario(usuario *entities.User) error {
	fmt.Printf("Iniciando el guardado del usuario con ID: %s\n", usuario.UserId)

	// Crear un documento BSON a partir del usuario
	documento := bson.D{
		{Key: "id_usuario", Value: usuario.UserId},
		{Key: "nickname", Value: usuario.Nickname},
		{Key: "token", Value: usuario.Token},
		{Key: "hora_ultima_accion", Value: usuario.LastActionTime},
		{Key: "estado", Value: usuario.State},
		{Key: "tipo", Value: usuario.Type},
		{Key: "idsala", Value: usuario.RoomId},     // La sala puede ser nil
		{Key: "namesala", Value: usuario.RoomName}, // La sala puede ser nil
	}
	fmt.Printf("Documento BSON para guardar usuario: %+v\n", documento)

	// Insertar el documento en la colección
	collection := mp.db.Collection("usuarios")
	fmt.Printf("Colección recogida de usuarios: %+v\n", documento)
	_, err := collection.InsertOne(context.Background(), documento)
	if err != nil {
		fmt.Printf("Error al guardar el usuario en MongoDB: %v\n", err)
		return fmt.Errorf("error al guardar el usuario en MongoDB: %v", err)
	}

	// Confirmación de éxito
	fmt.Printf("Usuario con ID: %s guardado correctamente en MongoDB\n", usuario.UserId)
	return nil
}
