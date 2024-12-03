package services

import (
	"backend/entities"
	"backend/models"
	"backend/persistence"
	"backend/types"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type GestionUsuarios struct {
	Usuarios []*entities.Usuario
}

func (gestion *GestionUsuarios) BuscarUsuarioPorToken(token string) (entities.Usuario, error) {
	for _, usuario := range gestion.Usuarios {
		if usuario.Token == token {
			return *usuario, nil
		}
	}
	return entities.Usuario{}, errors.New("usuario no encontrado")
}

func (gestion *GestionUsuarios) ObtenerTokenDeUsuario(nickname string) (string, error) {
	for _, usuario := range gestion.Usuarios {
		if usuario.Nickname == nickname {
			return usuario.Token, nil
		}
	}
	return "", errors.New("usuario no encontrado")
}

// VerificarUsuarioExistente verifica si un usuario con el nickname dado ya está registrado
func (gestion *GestionUsuarios) VerificarUsuarioExistente(nickname string) bool {
	// Verificar si el usuario ya existe
	for _, usuario := range gestion.Usuarios {
		if usuario.Nickname == nickname {
			return false // Usuario ya registrado
		}
	}
	return true // Usuario no registrado
}
// Función que registra un usuario y lo guarda en la base de datos
func (gestion *GestionUsuarios) RegistrarUsuario(nickname string, token string, sala *entities.Sala) (*entities.Usuario, error) {

	// Crear un nuevo usuario con la sala proporcionada
	nuevoUsuario := models.NewUsuarioGo(nickname, sala)

	// Crear una nueva instancia de MongoPersistenciae
	mongoPersistencia, err_per := persistence.NewMongoPersistencia()
	if err_per != nil {
		fmt.Println("Error al crear la instancia de MongoPersistencia:", err_per)
		return nil, err_per
	}

	// Guardar el usuario en la base de datos
	err_per = mongoPersistencia.GuardarUsuario(nuevoUsuario)
	if err_per != nil {
		fmt.Println("Error al guardar el usuario:", err_per)
		return nil, err_per
	} else {
		fmt.Println("Usuario guardado correctamente")
	}

	// Añadir el usuario a la lista de usuarios en memoria
	gestion.Usuarios = append(gestion.Usuarios, nuevoUsuario)
	return nuevoUsuario, nil
}
func (gestion *GestionUsuarios) ObtenerUsuariosActivos(idSala uuid.UUID) ([]*entities.Usuario, error) {
	// Lista para almacenar los usuarios activos
	var usuariosActivos []*entities.Usuario

	// Iterar sobre la lista de usuarios en memoria
	for _, usuario := range gestion.Usuarios {
		// Verificar si el usuario está activo y pertenece a la sala especificada
		if usuario.Estado == types.Activo && usuario.Sala != nil &&  usuario.Sala.ID  == idSala {
			usuariosActivos = append(usuariosActivos, usuario)
		}
	}

	// Retornar error si no se encuentran usuarios activos
	if len(usuariosActivos) == 0 {
		return nil, fmt.Errorf("no se encontraron usuarios activos en la sala con ID %s", idSala)
	}

	// Retornar la lista de usuarios activos
	return usuariosActivos, nil
}