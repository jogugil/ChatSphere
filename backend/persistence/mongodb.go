package persistence

import (
	"backend/entities"
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoPersistencia es el contenedor de la conexión y operaciones con MongoDB
type MongoPersistencia struct {
	client *mongo.Client
	db     *mongo.Database
}

// Configuración de MongoDB (puede ser configurada en la propia estructura de MongoPersistencia)
var uri = "mongodb://localhost:27017" // URI del servidor MongoDB
var dbName = "nombreBaseDeDatos"      // Nombre de la base de datos

var mongoInstance *MongoPersistencia
var once sync.Once

// NewMongoPersistencia crea una nueva instancia de MongoPersistencia, asegurando que solo haya una conexión
func NewMongoPersistencia() (*MongoPersistencia, error) {
	once.Do(func() {
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err != nil {
			fmt.Println("Error de conexión con MongoDB:", err)
			return
		}
		db := client.Database(dbName)
		mongoInstance = &MongoPersistencia{client: client, db: db}
	})

	if mongoInstance == nil {
		return nil, fmt.Errorf("no se pudo crear la instancia de MongoPersistencia")
	}

	return mongoInstance, nil
}

// GuardarMensaje guarda un solo mensaje en la base de datos MongoDB
func (mp *MongoPersistencia) GuardarMensaje(mensaje *entities.Mensaje) error {
	collection := mp.db.Collection("mensajes")
	_, err := collection.InsertOne(context.TODO(), mensaje)
	if err != nil {
		return fmt.Errorf("error al guardar el mensaje en MongoDB: %v", err)
	}
	return nil
}

// GuardarMensajesEnBaseDeDatos guarda una lista de mensajes en MongoDB
func (mp *MongoPersistencia) GuardarMensajesEnBaseDeDatos(mensajes []*entities.Mensaje) error {
	collection := mp.db.Collection("mensajes")

	// Convertimos los mensajes a una lista de interfaces{}
	var documentos []interface{}
	for _, mensaje := range mensajes {
		documentos = append(documentos, bson.D{
			{Key: "nick_usuario", Value: mensaje.Nickname},
			{Key: "fecha_envio", Value: mensaje.FechaEnvio},
			{Key: "id_sala", Value: mensaje.IdSala},
			{Key: "nombre_sala", Value: mensaje.NombreSala},
			{Key: "token_sesion", Value: mensaje.Token},
			{Key: "texto_mensaje", Value: mensaje.Mensaje},
		})
	}

	// Usamos InsertMany para insertar varios mensajes a la vez
	_, err := collection.InsertMany(context.TODO(), documentos)
	if err != nil {
		return fmt.Errorf("error al guardar los mensajes en MongoDB: %v", err)
	}

	// Todo salió bien
	fmt.Println("Mensajes guardados correctamente")
	return nil
}
func (mp *MongoPersistencia) ObtenerMensajesDesdeSala(idSala uuid.UUID) ([]entities.Mensaje, error) {
	collection := mp.db.Collection("mensajes")

	// Buscar el mensaje base
	var mensajeBase entities.Mensaje
	filterBase := bson.D{{Key: "id_sala", Value: idSala}}
	err := collection.FindOne(context.TODO(), filterBase).Decode(&mensajeBase)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no se  encontrsron mensajes en la sala: %v", idSala)
		}
		return nil, fmt.Errorf("error al obtener el mensaje base: %v", err)
	}

	cursor, err := collection.Find(context.TODO(), filterBase)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta: %v", err)
	}
	defer cursor.Close(context.TODO())

	var mensajes []entities.Mensaje
	for cursor.Next(context.TODO()) {
		var mensaje entities.Mensaje
		if err := cursor.Decode(&mensaje); err != nil {
			return nil, fmt.Errorf("error al decodificar el mensaje: %v", err)
		}
		mensajes = append(mensajes, mensaje)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre el cursor: %v", err)
	}

	return mensajes, nil
}
func (mp *MongoPersistencia) ObtenerMensajesDesdeId(idSala uuid.UUID, idMensaje string) ([]entities.Mensaje, error) {
	collection := mp.db.Collection("mensajes")

	// Buscar el mensaje base
	var mensajeBase entities.Mensaje
	filterBase := bson.D{{Key: "id_mensaje", Value: idMensaje}, {Key: "id_sala", Value: idSala}}
	err := collection.FindOne(context.TODO(), filterBase).Decode(&mensajeBase)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no se encontró el mensaje con id_mensaje: %v en la sala: %v", idMensaje, idSala)
		}
		return nil, fmt.Errorf("error al obtener el mensaje base: %v", err)
	}

	// Filtro para mensajes posteriores
	filter := bson.D{
		{Key: "id_sala", Value: idSala},
		{Key: "fecha_envio", Value: bson.D{{Key: "$gte", Value: mensajeBase.FechaEnvio}}},
	}
	findOptions := options.Find().SetSort(bson.D{{Key: "fecha_envio", Value: 1}})

	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta: %v", err)
	}
	defer cursor.Close(context.TODO())

	var mensajes []entities.Mensaje
	for cursor.Next(context.TODO()) {
		var mensaje entities.Mensaje
		if err := cursor.Decode(&mensaje); err != nil {
			return nil, fmt.Errorf("error al decodificar el mensaje: %v", err)
		}
		mensajes = append(mensajes, mensaje)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre el cursor: %v", err)
	}

	return mensajes, nil
}
func (mp *MongoPersistencia) GuardarUsuario(usuario *entities.Usuario) error {
	collection := mp.db.Collection("usuarios")

	// Crear un documento BSON a partir del usuario
	documento := bson.D{
		{Key: "id_usuario", Value: usuario.IdUsuario},
		{Key: "nickname", Value: usuario.Nickname},
		{Key: "token", Value: usuario.Token},
		{Key: "hora_ultima_accion", Value: usuario.HoraUltimaAccion},
		{Key: "estado", Value: usuario.Estado},
		{Key: "tipo", Value: usuario.Tipo},
		{Key: "sala", Value: usuario.Sala}, // La sala puede ser nil
	}

	// Insertar el documento en la colección
	_, err := collection.InsertOne(context.TODO(), documento)
	if err != nil {
		return fmt.Errorf("error al guardar el usuario en MongoDB: %v", err)
	}

	return nil
}
