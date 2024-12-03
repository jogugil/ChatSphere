package entities

import (
	"time"
	"github.com/google/uuid"
)

type Sala struct {
	ID                  uuid.UUID           `json:"id"`
	Name                string              `json:"name"`
	SalaTipo            string              `json:"salatipo"`
	Usuarios            []Usuario           `json:"usuariosChat"`
	TiempoActividad     int                 `json:"tiempoActividad"`
	Estado              bool                `json:"estado"`
	UltimaActualizacion time.Time           `json:"ultimaActualizacion"`
	FechaCreacion       time.Time           `json:"fechaCreacion"`
	FechaDesactivacion  time.Time           `json:"fechaDesactivacion"`
	HistoricoMensajes   *ColaCircular       `json:"historicoMensajes"` // Cola Circular
}