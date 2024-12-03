package services

import (
	"backend/entities"
	"backend/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type SecModServidorChat struct {
	GestionMensajes GestionMensajes
	GestionSalas    GestionSalas
	GestionUsuarios GestionUsuarios
}

// Método para crear una nueva instancia de SecModServidorChat
func NewSecModServidorChat() *SecModServidorChat {
	return &SecModServidorChat{
		GestionMensajes: GestionMensajes{}, // Inicializa tu gestión de mensajes aquí
		GestionSalas:    GestionSalas{},    // Inicializa tu gestión de salas aquí
		GestionUsuarios: GestionUsuarios{}, // Inicializa tu gestión de usuarios aquí
	}
}

// ValidarTokenAccion valida si un token es correcto y si el nickname está asociado a ese token
func (secMod *SecModServidorChat) ValidarTokenAccion(token, nickname, opt string) bool {
	// Aquí se debe validar el token
	usuario, err := secMod.GestionUsuarios.BuscarUsuarioPorToken(token)
	if err != nil || usuario.Nickname != nickname {
		return false
	}

	// Validación de la acción (opt) se puede agregar aquí según el tipo de acción
	switch opt {
	case "enviarMensaje":
		return true
	case "verSala":
		return true
	default:
		return false
	}
}

// CrearTokenSesion genera un nuevo token para un usuario al registrarse
func (secMod *SecModServidorChat) CrearTokenSesion(nickname string) string {
	// Crea un token de sesión único para el usuario
	return models.CrearTokenSesion(nickname)
}

// EjecutarLogin maneja el proceso de login
func (secMod *SecModServidorChat) EjecutarLogin(nickname string) (*entities.Usuario, error) {
	// Llamar a RegistrarUsuario para asegurarnos de que el usuario esté registrado
	if !secMod.GestionUsuarios.VerificarUsuarioExistente(nickname) {
		return nil, errors.New("El nickname ya está en uso.!!!")
	}
	token := secMod.CrearTokenSesion(nickname)
	if token == "" {
		return nil, errors.New("no se pudo crear el token")
	}

	newUser, err := secMod.GestionUsuarios.RegistrarUsuario(nickname, token, &secMod.GestionSalas.SalaPrincipal)
	if err != nil {
		return nil, fmt.Errorf("error al registrar el usuario %v", err)
	}
	secMod.GestionSalas.SalaPrincipal.Usuarios = append(secMod.GestionSalas.SalaPrincipal.Usuarios, *newUser)

	// Mostrar el mensaje de éxito
	fmt.Printf("Usuario %s logueado con el token %s\n", nickname, token)
	return newUser, nil
}

// EjecutarEnvioMensaje permite a un usuario enviar un mensaje
func (secMod *SecModServidorChat) EjecutarEnvioMensaje(nickname, token, mensaje string, idSala uuid.UUID) error {
	// Validar si el usuario está autorizado a enviar el mensaje2
	token_usuario, err := secMod.GestionUsuarios.ObtenerTokenDeUsuario(nickname)
	if !secMod.ValidarTokenAccion(token_usuario, nickname, "enviarMensaje") {
		return errors.New("token inválido o acción no permitida")
	}
	if token_usuario != token {
		return errors.New("token inválido o acción no permitida")
	}

	// Llamar a la lógica para enviar el mensaje
	err = secMod.GestionSalas.EnviarMensaje(idSala, nickname, mensaje)
	if err != nil {
		return err
	}

	fmt.Printf("Mensaje enviado por %s: %s\n", nickname, mensaje)
	return nil
}

// EjecutarVerSala permite a un usuario ver la lista de mensajes en una sala
func (secMod *SecModServidorChat) EjecutarVerSala(nickname string, idSala uuid.UUID) error {
	// Validar si el usuario está autorizado a ver la sala
	token, err := secMod.GestionUsuarios.ObtenerTokenDeUsuario(nickname)
	if err == nil {
		return errors.New("token inválido o acción no permitida")
	}
	if !secMod.ValidarTokenAccion(token, nickname, "verSala") {
		return errors.New("token inválido o acción no permitida")
	}

	// Llamar a la lógica para ver los mensajes de la sala
	mensajes, err := secMod.GestionSalas.ObtenerMensajesDesdeCola(idSala)
	if err != nil {
		return fmt.Errorf("token inválido o acción no permitida %v", err)
	}

	fmt.Println("Mensajes en la sala:")
	for _, msg := range mensajes {
		fmt.Printf("[%s]: %s\n", msg.Nickname, msg.Mensaje)
	}

	return nil
}
