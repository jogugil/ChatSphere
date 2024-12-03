package models

import (
 
	"time"

	"backend/entities"

	"github.com/google/uuid"
)

// Crear un nuevo mensaje con los datos proporcionados
func CrearMensaje (mensajeText string, usuario entities.Usuario, idSala uuid.UUID) *entities.Mensaje {
	return &entities.Mensaje{
		IDM:         generarIDMensaje(), // Función para generar IDs únicos. Generado en el Servidor
		Tipo:        entities.Ordinario, // Tipo predeterminado: Ordinario.En esta versión siempre es ordinario
		FechaEnvio:  time.Now(),         // Fecha del cliente. Se genera en el cleinte antes de enviar el mensaje
		FechaServer: time.Now(),         // Fecha del servidor. Se genera en el srvidor antes de giardalo en el biffer
		Nickname:    usuario.Nickname,   // Viene del cliente
		Token:       usuario.Token,      // Viene del cliente
		Mensaje:     mensajeText,        // Viene del cliente
		IdSala:      idSala,             // Viene del cliente
	}
}

// Crear un nuevo mensaje con datos enviados por el cliente
func CrearMensajeConFecha(mensajeText string, usuario entities.Usuario, idSala uuid.UUID, fechaEnvio time.Time) *entities.Mensaje {
	return &entities.Mensaje{
		IDM:          generarIDMensaje(),
		Tipo:        entities.Ordinario, // Tipo predeterminado: Ordinario
		FechaEnvio:  fechaEnvio,         // Fecha enviada por el cliente
		FechaServer: time.Now(),         // Fecha generada por el servidor
		Nickname:    usuario.Nickname,
		Token:       usuario.Token,
		Mensaje:     mensajeText,
		IdSala:      idSala,
	}
}


// Generador de ID de mensajes
func generarIDMensaje() uuid.UUID  {
	return  uuid.New() 
}
