package entities

import (
	"backend/types" // Importamos el paquete interfaces
	"time"
	// Creación de uuid's únicos
)

type Usuario struct {
	IdUsuario        string
	Nickname         string
	Token            string
	HoraUltimaAccion time.Time
	Estado           types.EstadoUsuario
	Tipo             string
	Sala             *Sala // Cambié esto para que Sala sea una referencia a un objeto Sala
}
