package entities

import (
	"backend/types" // Importamos el paquete interfaces
	"time"

	"github.com/google/uuid"
	// Creación de uuid's únicos
)

type Usuario struct {
	IdUsuario        string
	Nickname         string
	Token            string
	HoraUltimaAccion time.Time
	Estado           types.EstadoUsuario
	Tipo             string
	IdSala           uuid.UUID
	NameSala         string
}
