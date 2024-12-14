package models

import (
	"backend/entities"
	"backend/types" // Importamos el paquete interfaces
	"log"
	"time"

	"github.com/google/uuid" // Creación de uuid's únicos
)

// Alias del tipo Usuario
type LocalUsuario entities.Usuario

// Implementación de los métodos de la interfaz UsuarioChat
func (u *LocalUsuario) IniciarSesion() bool {
	u.LastActionTime = time.Now()
	u.State = types.Activo
	return true
}

func (u *LocalUsuario) EliminarSesion() bool {
	u.State = types.Inactivo
	return true
}

func (u *LocalUsuario) ActualizarEstado() {
	u.State = types.Activo
}

func (u *LocalUsuario) UnirseASala(room *entities.Sala) {
	u.RoomId = room.ID
}

func (u *LocalUsuario) SalirDeSala() {
	roomId, err := uuid.Parse("00000000-0000-0000-0000-000000000000")
	if err != nil {
		log.Printf("error a poner nil una sla del usuario: %v", err)
	}
	u.RoomId = roomId // Ahora se elimina la referencia a la sala
}

func NewUsuarioGo(nickname string, room *entities.Sala) *entities.Usuario {
	return &entities.Usuario{
		UserId:         "usr-" + uuid.New().String(),
		Nickname:       nickname,
		Token:          CrearTokenSesion(nickname),
		Type:           "usuariochat",
		LastActionTime: time.Now(),
		State:          types.Activo,
		RoomId:         room.ID,
		RoomName:       room.Name,
	}
}
