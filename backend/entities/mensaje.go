package entities

import (
	"time"

	"github.com/google/uuid"
)

// Tipo de mensaje
type TipoMensaje int

const (
	Ordinario TipoMensaje = iota + 1
	Administracion
	Info
)

// Clase Mensaje
type Mensaje struct {
	IDM         uuid.UUID   `json:"id"`
	Tipo        TipoMensaje `json:"tipo"`
	FechaEnvio  time.Time   `json:"fechaEnvio"`
	FechaServer time.Time   `json:"fechaServer"`
	Nickname    string      `json:"nickname"`
	Token       string      `json:"tokenjwt"`
	Mensaje     string      `json:"mensaje"`
	IdSala      uuid.UUID   `json:"idSala"`
	NombreSala  string      `json:"nombresala"`
}

// Actualizar el tipo de mensaje
func (m *Mensaje) ActualizarTipo(tipo TipoMensaje) {
	m.Tipo = tipo
}
