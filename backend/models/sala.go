package models

import (
	"backend/entities"
 
	"log"
)

// Alias del tipo Sala
type LocalSala entities.Sala

func (sala *LocalSala) AgregarUsuario(usuario entities.Usuario) {
	sala.Usuarios = append(sala.Usuarios, usuario)
}

func (sala *LocalSala) QuitarUsuario(usuario entities.Usuario) {
	// Lógica para eliminar un usuario de la sala

	// Recorre la lista de usuarios
	for i, u := range sala.Usuarios {
		// Si se encuentra al usuario, lo elimina de la lista
		if u == usuario {
			// Elimina el usuario de la lista manteniendo los demás
			sala.Usuarios = append(sala.Usuarios[:i], sala.Usuarios[i+1:]...)
			return
		}
	}
}

func (sala *LocalSala) EnviarMensaje(usuario entities.Usuario, mensaje entities.Mensaje) {
	// Añadir el mensaje a la cola
	sala.HistoricoMensajes.Enqueue(mensaje)
}

func (sala *LocalSala) ObtenerMensajesSala() []entities.Mensaje {
	mensajes, err := sala.HistoricoMensajes.ObtenerTodos(sala.ID)
	if err != nil {
		log.Printf("error al obtener lso mensjaes de la cola circular de la sala %s", sala.Name)
		return nil
	}
	return mensajes
}
func (sala *LocalSala) ObtenerMensajesdesdeId(idMensaje string) []entities.Mensaje {
	mensjes, err := sala.HistoricoMensajes.ObtenerMensajesDesdeId(sala.ID, idMensaje)
	if err != nil {
		log.Printf("error al obtener lso mensjaes de la cola circular de la sala %s", sala.Name)
		return nil
	}
	return mensjes
}
