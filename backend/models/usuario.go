package models

import (
	"backend/types" // Importamos el paquete interfaces
	"backend/entities"
	"time"

	"github.com/google/uuid" // Creación de uuid's únicos
)
// Alias del tipo Usuario
type LocalUsuario entities.Usuario

// Implementación de los métodos de la interfaz UsuarioChat
func (u *LocalUsuario) IniciarSesion() bool {
	u.HoraUltimaAccion = time.Now()
	u.Estado = types.Activo
	return true
}

func (u *LocalUsuario) EliminarSesion() bool {
	u.Estado = types.Inactivo
	return true
}

func (u *LocalUsuario) ActualizarEstado() {
	u.Estado = types.Activo
}

func (u *LocalUsuario) UnirseASala(sala *entities.Sala) {
	u.Sala = sala
}

func (u *LocalUsuario) SalirDeSala() {
	u.Sala = nil // Ahora se elimina la referencia a la sala
}

func NewUsuarioGo(nickname string, sala *entities.Sala) *entities.Usuario {
	return &entities.Usuario{
		IdUsuario:        "usr-" + uuid.New().String(),
		Nickname:         nickname,
		Token:            CrearTokenSesion(nickname),
		Tipo:             "usuariochat",
		HoraUltimaAccion: time.Now(),
		Estado:           types.Activo,
		Sala:             sala, // Se pasa la referencia de la sala
	}
}
