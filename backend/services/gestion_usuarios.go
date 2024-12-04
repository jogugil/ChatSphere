package services

import (
	"backend/entities"
	"backend/models"
	"backend/persistence"
	"backend/types"
	"errors"
	"fmt"
	"log"
	"sync"
)

type GestionUsuarios struct {
	Usuarios []*entities.Usuario // Lista de usuarios
	onceUser sync.Once           // Para garantizar la inicialización única de la instancia
}

var instanciaUser *GestionUsuarios // Instancia única de GestionUsuarios

// NuevaGestionUsuarios crea y devuelve una instancia Singleton de GestionUsuarios
func NuevaGestionUsuarios() *GestionUsuarios {
	// Si la instancia no ha sido creada, se crea una nueva
	if instanciaUser == nil {
		instanciaUser = &GestionUsuarios{
			Usuarios: make([]*entities.Usuario, 0), // Inicializamos la lista de usuarios vacía
		}
	}
	// Usamos once.Do para asegurarnos que la configuración solo se realice una vez
	instanciaUser.onceUser.Do(func() {
		log.Println("NuevaGestionUsuarios: Inicializando la instancia de Gestión de Usuarios.")
		// Aquí podrías cargar los usuarios desde una base de datos o archivo si es necesario
		// Ejemplo de carga de usuarios (esto depende de tu implementación y necesidades)
		// instancia.CargarUsuariosDesdeArchivo("usuarios.txt")
		log.Println("NuevaGestionUsuarios: Inicialización completada.")
	})
	return instanciaUser
}

func (gestion *GestionUsuarios) BuscarUsuarioPorToken(token string) (entities.Usuario, error) {
	fmt.Println("Iniciando búsqueda de usuario por token:", token)
	for _, usuario := range gestion.Usuarios {
		if usuario.Token == token {
			fmt.Println("Usuario encontrado:", usuario.Nickname)
			return *usuario, nil
		}
	}
	fmt.Println("Usuario no encontrado para el token:", token)
	return entities.Usuario{}, errors.New("usuario no encontrado")
}

func (gestion *GestionUsuarios) ObtenerTokenDeUsuario(nickname string) (string, error) {
	fmt.Println("Obteniendo token para el usuario:", nickname)
	for _, usuario := range gestion.Usuarios {
		if usuario.Nickname == nickname {
			fmt.Println("Token encontrado para el usuario:", usuario.Token)
			return usuario.Token, nil
		}
	}
	fmt.Println("Usuario no encontrado para el nickname:", nickname)
	return "", errors.New("usuario no encontrado")
}

// VerificarUsuarioExistente verifica si un usuario con el nickname dado ya está registrado
func (gestion *GestionUsuarios) VerificarUsuarioExistente(nickname string) bool {
	fmt.Println("Verificando si el usuario ya existe:", nickname)
	for _, usuario := range gestion.Usuarios {
		if usuario.Nickname == nickname {
			fmt.Println("Usuario ya registrado:", nickname)
			return false // Usuario ya registrado
		}
	}
	fmt.Println("Usuario no registrado:", nickname)
	return true // Usuario no registrado
}

// Función que registra un usuario y lo guarda en la base de datos
func (gestion *GestionUsuarios) RegistrarUsuario(nickname string, token string, sala *entities.Sala) (*entities.Usuario, error) {
	fmt.Println("Registrando nuevo usuario:", nickname)

	// Crear un nuevo usuario con la sala proporcionada
	nuevoUsuario := models.NewUsuarioGo(nickname, sala)
	fmt.Println("Nuevo usuario creado:", nuevoUsuario.Nickname)

	// Crear una nueva instancia de MongoPersistencia
	mongoPersistencia, err_per := persistence.ObtenerInstanciaDB()
	if err_per != nil {
		fmt.Println("Error al crear la instancia de MongoPersistencia:", err_per)
		return nil, err_per
	}

	// Guardar el usuario en la base de datos
	fmt.Println("Guardando usuario en la base de datos...")
	err_per = (*mongoPersistencia).GuardarUsuario(nuevoUsuario)
	if err_per != nil {
		fmt.Println("Error al guardar el usuario:", err_per)
		return nil, err_per
	} else {
		fmt.Println("Usuario guardado correctamente en la base de datos")
	}

	// Añadir el usuario a la lista de usuarios en memoria
	gestion.Usuarios = append(gestion.Usuarios, nuevoUsuario)
	fmt.Println("Usuario añadido a la lista en memoria:", nuevoUsuario.Nickname)
	return nuevoUsuario, nil
}

// Pra otras versiones pasar el uuid de la sala y devolver lso usuarios presentes en dicha sala
// El objeto slaa tiene el lisado de usuarios activos en la sala
// Método ObtenerUsuariosActivos que devuelve una lista de usuarios activos
func (g *GestionUsuarios) ObtenerUsuariosActivos() []*entities.Usuario {
	var usuariosActivos []*entities.Usuario
	for _, usuario := range g.Usuarios {
		if usuario.Estado == types.Activo {
			usuariosActivos = append(usuariosActivos, usuario)
		}
	}
	return usuariosActivos
}
