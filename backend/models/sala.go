package models

import (
	"backend/entities"
	"log"
	"github.com/google/uuid"
)

// Alias del tipo Sala
type LocalSala entities.Sala

func (sala *LocalSala) AgregarUsuario(usuario entities.Usuario) {
	log.Printf("AgregarUsuario: Entrando con usuario %v a la sala %s", usuario, sala.Name)
	
	// Operación clave: Añadir el usuario a la lista
	sala.Usuarios = append(sala.Usuarios, usuario)
	
	log.Printf("AgregarUsuario: Usuario %v agregado correctamente a la sala %s", usuario, sala.Name)
}

func (sala *LocalSala) QuitarUsuario(usuario entities.Usuario) {
	log.Printf("QuitarUsuario: Entrando con usuario %v para quitar de la sala %s", usuario, sala.Name)
	
	// Recorre la lista de usuarios
	for i, u := range sala.Usuarios {
		// Si se encuentra al usuario, lo elimina de la lista
		if u == usuario {
			log.Printf("QuitarUsuario: Usuario %v encontrado en la posición %d", usuario, i)
			
			// Operación clave: Elimina el usuario manteniendo los demás
			sala.Usuarios = append(sala.Usuarios[:i], sala.Usuarios[i+1:]...)
			
			log.Printf("QuitarUsuario: Usuario %v eliminado correctamente de la sala %s", usuario, sala.Name)
			return
		}
	}
	
	log.Printf("QuitarUsuario: Usuario %v no encontrado en la sala %s", usuario, sala.Name)
}

func (sala *LocalSala) EnviarMensaje(usuario entities.Usuario, mensaje entities.Mensaje) {
	log.Printf("EnviarMensaje: Entrando con usuario %v y mensaje %v en la sala %s", usuario, mensaje, sala.Name)
	
	// Operación clave: Añadir el mensaje a la cola
	sala.HistoricoMensajes.Enqueue(mensaje)

	log.Printf("EnviarMensaje: Mensaje %v añadido correctamente a la cola de la sala %s", mensaje, sala.Name)
}

func (sala *LocalSala) ObtenerMensajesSala() []entities.Mensaje {
	log.Printf("ObtenerMensajesSala: Entrando para obtener mensajes de la sala %s", sala.Name)
	
	// Operación clave: Obtener todos los mensajes
	mensajes, err := sala.HistoricoMensajes.ObtenerTodos(sala.ID)
	if err != nil {
		log.Printf("ObtenerMensajesSala: Error al obtener mensajes de la sala %s: %v", sala.Name, err)
		return nil
	}
	
	log.Printf("ObtenerMensajesSala: Mensajes obtenidos correctamente de la sala %s", sala.Name)
	return mensajes
}

func (sala *LocalSala) ObtenerMensajesdesdeId(idMensaje uuid.UUID) []entities.Mensaje {
	log.Printf("ObtenerMensajesdesdeId: Entrando con idMensaje %v para la sala %s", idMensaje, sala.Name)
	
	// Operación clave: Obtener mensajes desde un ID específico
	mensajes, err := sala.HistoricoMensajes.ObtenerMensajesDesdeId(sala.ID, idMensaje)
	if err != nil {
		log.Printf("ObtenerMensajesdesdeId: Error al obtener mensajes desde ID %v en la sala %s: %v", idMensaje, sala.Name, err)
		return nil
	}
	
	log.Printf("ObtenerMensajesdesdeId: Mensajes obtenidos correctamente desde ID %v en la sala %s", idMensaje, sala.Name)
	return mensajes
}
