// interfaces/persistencia.go
package entities

import (
	"github.com/google/uuid"
)

type Persistencia interface {
	GuardarUsuario(usuario *Usuario) error
	GuardarSala(sala Sala) error
	ObtenerSala(id uuid.UUID) (Sala, error)
	GuardarMensajesEnBaseDeDatos(mensajes []Mensaje) error
	ObtenerMensajesDesdeSala(idSala uuid.UUID) ([]Mensaje, error)
	ObtenerMensajesDesdeId(idSala uuid.UUID, idMensaje string) ([]Mensaje, error)
}
