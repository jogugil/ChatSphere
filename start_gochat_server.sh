#!/bin/bash

# ==============================================
# Proyecto: Servidor GoChat
# Descripción:
#   Este script gestiona la instalación de las dependencias necesarias para ejecutar
#   el servidor GoChat. Verifica la instalación de Go y Docker, comprueba si el puerto
#   está disponible y configura el archivo .env antes de ejecutar el servidor GoChat
#   en segundo plano. Después de la ejecución, restaura los valores por defecto del
#   archivo de configuración.
#
# Licencia Open Source:
#   Este código está bajo la Licencia Apache 2.0.
#
#   © 2024 SMARTIASERVICE's / FaaioFlex
#   Autor: José Javier Gutiérrez Gil
#   Emails: jogugi@posgrado.upv.es, jogugil@gmail.com
#
# Licencia Open Source:
#     Licencia Apache 2.0 - Apache Software Foundation
#
# ────────────────────────────────────────────────────────────────────────────────────────────────────────
#      ███████╗████████╗ ██████╗  █████╗ ██████╗██╗████████╗ █████╗
#      ██╔════╝╚══██╔══╝██╔══██╗██╔══██╗██╔══██╗██║╚══██╔══╝██╔══╝
#      █████╗     ██║   ██║  ██║███████║██████╔╝██║   ██║   █████╗
#      ██╔══╝     ██║   ██║  ██║██╔══██║██╔══██╗██║   ██║   ██╔══╝
#      ███████╗   ██║   ██████╔╝██║  ██║██████╔╝██║   ██║   ███████╗
#      ╚══════╝   ╚═╝   ╚═════╝ ╚═╝  ╚═╝╚══════╝ ╚═╝   ╚═╝   ╚══════╝
#
#  ____                  __   ____ ___    ____                 _           _
#  / __/__ _  ___ _ ____ / /_ /  _// _ |  / __/___  ____ _  __ (_)____ ___ ( )___
# _\ \ /  ' \/ _ `// __// __/_/ / / __ | _\ \ / -_)/ __/| |/ // // __// -_)|/(_-<
#/___//_/_/_/\_,_//_/   \__//___//_/ |_|/___/ \__//_/   |___//_/ \__/ \__/  /___/
#
#     ,d8888b                        d8,           ,d8888bd8b
#     88P'                          `8P            88P'   88P
#  d888888P                                     d888888P d88
#    ?88'd888b8b   d888b8b   .d888b, 88b d8888b   ?88'   888   d8888b?88,  88P
#    88Pd8P' ?88  d8P' ?88   ?8b,    88Pd8P' ?88  88P    ?88  d8b_,dP `?8bd8P'
#   d88 88b  ,88b 88b  ,88b    `?8b d88 88b  d88 d88      88b 88b     d8P?8b,
#  d88' `?88P'`88b`?88P'`88b`?888P'd88' `?8888P'd88'       88b`?888P'd8P' `?8b
#
#
# ────────────────────────────────────────────────────────────────────────────────────────────────────────
#     SMARTIASERVICE's - FaaioFlex (© 2024)
#     Este código está bajo la Licencia Open Source Apache 2.0, otorgad  por SMARTIASERVICE's.
#
#	Ver detalles completos en: https://www.apache.org/licenses/LICENSE-2.0
#
# =======================================================================================================





# ==============================================
# Función: check_go
# Descripción:
#   Verifica si Go está instalado en el sistema. Si no lo está, solicita al usuario 
#   la instalación de Go. Si el usuario acepta, instala Go utilizando los comandos 
#   apropiados para Ubuntu.
#
# Parámetros: Ninguno
#
# Uso:
#   check_go
#
# Autor: José Javier Gutiérrez Gil
# ==============================================
function check_go() {
    if ! command -v go &> /dev/null; then
        echo "Go no está instalado en tu sistema."
        read -p "¿Deseas instalar Go? (Y/n): " instalar_go
        if [[ $instalar_go =~ ^[Yy]$ ]]; then
            # Comando para instalar Go (depende de la distribución de Linux)
            echo "Instalando Go..."
            sudo apt update
            sudo apt install -y golang-1.22-go
            if ! command -v go &> /dev/null; then
                echo "Hubo un error al instalar Go. Asegúrate de tener permisos de root."
                exit 1
            fi
            echo "Go ha sido instalado correctamente."
        else
            echo "No se puede continuar sin Go. El script se detendrá."
            exit 1
        fi
    else
        echo "Go ya está instalado en tu sistema."
    fi
}

# ==============================================
# Función: check_docker
# Descripción:
#   Verifica si Docker está instalado en el sistema. Si no lo está, solicita al usuario 
#   la instalación de Docker. Si el usuario acepta, instala Docker utilizando los comandos 
#   apropiados para Ubuntu.
#
# Parámetros: Ninguno
#
# Uso:
#   check_docker
#
# Autor: José Javier Gutiérrez Gil
# ==============================================
function check_docker() {
    if ! command -v docker &> /dev/null; then
        echo "Docker no está instalado en tu sistema."
        read -p "¿Deseas instalar Docker? (Y/n): " instalar_docker
        if [[ $instalar_docker =~ ^[Yy]$ ]]; then
            # Comando para instalar Docker en Ubuntu
            echo "Instalando Docker..."
            sudo apt update
            sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
            curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
            echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
            sudo apt update
            sudo apt install -y docker-ce
            if ! command -v docker &> /dev/null; then
                echo "Hubo un error al instalar Docker. Asegúrate de tener permisos de root."
                exit 1
            fi
            echo "Docker ha sido instalado correctamente."
        else
            echo "No se puede continuar sin Docker. El script se detendrá."
            exit 1
        fi
    else
        echo "Docker ya está instalado en tu sistema."
    fi
}

# ==============================================
# Función: check_dependencies
# Descripción:
#   Verifica si las dependencias necesarias para ejecutar GoChat están instaladas.
#   Si el archivo go.mod no está presente, ofrece la opción de crear uno. Si las dependencias
#   no están instaladas, permite al usuario instalarlas con 'go mod tidy'.
#
# Parámetros: Ninguno
#
# Uso:
#   check_dependencies
#
# Autor: José Javier Gutiérrez Gil
# ==============================================
function check_dependencies() {
    echo "Verificando dependencias del servidor GoChat..."

    # Guardar el directorio actual para regresar en caso de error
    ROOT_DIR=$(pwd)

    # Verificar si el directorio 'backend' existe
    if [ ! -d "backend" ]; then
        echo "El directorio 'backend' no existe. El script se detendrá."
        exit 1
    fi

    # Cambiar al directorio backend
    cd backend

    # Verificar si go.mod existe
    if [ ! -f "go.mod" ]; then
        echo "No se encuentra el archivo go.mod. ¿Quieres crear uno? (Y/n):"
        read crear_go_mod
        if [[ $crear_go_mod =~ ^[Yy]$ ]]; then
            # Crear el archivo go.mod con el nombre del módulo 'backend'
            go mod init backend
            echo "Archivo go.mod creado."

            # Verificar e instalar dependencias necesarias
            echo "Instalando dependencias necesarias..."
            go mod tidy
            if [ $? -ne 0 ]; then
                echo "Hubo un error al instalar las dependencias. Asegúrate de tener una conexión a Internet."
                cd "$ROOT_DIR"  # Regresar al directorio raíz
                exit 1
            fi
            echo "Las dependencias se han instalado correctamente."
        else
            echo "El servidor GoChat no se puede iniciar sin un archivo go.mod."
            cd "$ROOT_DIR"  # Regresar al directorio raíz
            exit 1
        fi
    else
        # Eliminar go.mod y go.sum si existen y volver a crearlos
        rm -f go.mod go.sum
        go mod init backend
        go mod tidy
        if [ $? -ne 0 ]; then
            echo "Hubo un error al instalar las dependencias. Asegúrate de tener una conexión a Internet."
            cd "$ROOT_DIR"  # Regresar al directorio raíz
            exit 1
        fi
        echo "Las dependencias se han instalado correctamente."
    fi

    # Volver al directorio raíz
    cd "$ROOT_DIR"
}

# ==============================================
# Función: check_chat_port
# Descripción:
#   Verifica si el puerto indicado por el usuario está libre. Si el puerto está en uso,
#   ofrece la opción de eliminar el proceso que está utilizando el puerto.
#
# Parámetros: Ninguno
#
# Uso:
#   check_chat_port
#
# Autor: José Javier Gutiérrez Gil
# ==============================================
# Función para actualizar el puerto en el archivo .env
function update_env_port() {
    echo "Actualizando archivo .env con el puerto configurado..."
    cp backend/.env backend/.env.bak  # Hacer una copia de seguridad del archivo .env
    sed -i "s/^CHAT_PORT=.*$/CHAT_PORT=$CHAT_PORT/" backend/.env
    echo "Archivo .env actualizado con el puerto: $CHAT_PORT"
}

# Función para verificar el puerto del servidor GoChat
function check_chat_port() {
    read -p "Ingrese el puerto para el servidor GoChat (por defecto 8081): " CHAT_PORT
    CHAT_PORT=${CHAT_PORT:-8081}  # Asigna el puerto por defecto si no se ingresa uno

    echo "Puerto seleccionado: $CHAT_PORT"

    # Comprobar si el puerto está en uso
    if lsof -i :$CHAT_PORT &> /dev/null; then
        echo "El puerto $CHAT_PORT está en uso por otro proceso."
        read -p "¿Deseas eliminar el proceso que está usando el puerto $CHAT_PORT? (Y/n): " eliminar
        if [[ $eliminar =~ ^[Yy]$ ]]; then
            PID=$(lsof -t -i :$CHAT_PORT)
            if [ -n "$PID" ]; then  # Verificar si PID no está vacío
                sudo kill -9 $PID
                echo "Proceso con PID $PID detenido y puerto $CHAT_PORT liberado."
            else
                echo "No se pudo obtener el PID del proceso que usa el puerto $CHAT_PORT."
            fi
        else
            echo "El script se detendrá debido a que el puerto está en uso."
            exit 1
        fi
    fi

    # Llamar a la función para actualizar el archivo .env si el puerto no es 8081
    if [ "$CHAT_PORT" -ne 8081 ]; then
        update_env_port
    fi
}
# ==============================================
# Función: start_backend
# Descripción:
#   Inicia el servidor GoChat en segundo plano usando Docker Compose y verifica su estado.
#
# Parámetros:
#   Ninguno
#
# Uso:
#   start_backend
#
# Autor: José Javier Gutiérrez Gil
# ==============================================
function start_backend() {
    echo "Iniciando el servidor GoChat en segundo plano..."

    # Verificar si el subdirectorio 'backend' existe
    if [ ! -d "./backend" ]; then
        echo "El directorio 'backend' no existe. Asegúrate de que el directorio esté presente."
        exit 1
    fi

    # Ir al directorio raíz para asegurarnos de que estamos en el lugar correcto
    CURRENT_DIR=$(pwd)
    cd  ./backend

    pwd
 
    # Verificar si el archivo main.go existe en el subdirectorio backend
    if [ ! -f "main.go" ]; then
        echo "El archivo 'backend/main.go' no se encuentra. Asegúrate de que el archivo exista."
        cd "$CURRENT_DIR"  # Regresar al directorio original antes de salir
        exit 1
    fi

    # Crear un archivo de log para capturar el output
    LOG_FILE="gochat_server.log"

    # Ejecutar el servidor GoChat y redirigir el output a un archivo de log
    go run main.go > "$LOG_FILE" 2>&1 &

    # Obtener el PID del servidor
    SERVER_PID=$!

    # Volver al directorio original
    cd "$CURRENT_DIR"

    # Esperar un poco para asegurarnos de que el servidor tenga tiempo de iniciar
    sleep 2

    # Verificar si el servidor se está ejecutando
    if ps -p $SERVER_PID > /dev/null; then
        echo "Servidor GoChat iniciado correctamente en el puerto $CHAT_PORT."
    else
        echo "Hubo un error al iniciar el servidor GoChat. Verifica los logs en $LOG_FILE para más detalles."
        cat "$LOG_FILE"  # Mostrar el contenido del log para depurar
        cd "$CURRENT_DIR"  # Regresar al directorio original antes de salir
        exit 1
    fi
    cd "$CURRENT_DIR"  # Regresar al directorio original antes de salir
    pwd
}

# ==============================================
# Función: restore_env
# Descripción:
#   Restaura el archivo .env desde su copia de seguridad.
#
# Parámetros:
#   Ninguno
#
# Uso:
#   restore_env
#
# Autor: José Javier Gutiérrez Gil
# ==============================================
function restore_env() {
    echo "Restaurando el archivo .env desde la copia de seguridad..."
    cp backend/.env.bak backend/.env
    echo "Archivo .env restaurado correctamente."
}
# ==============================================
# Función: check_docker_mongodb
# Descripción:
#   Verifica si MongoDB está disponible en Docker.
#
# Parámetros: Ninguno
#
# Uso:
#   check_docker_mongodb
#
# Autor: José Javier Gutiérrez Gil
# ==============================================
function check_docker_mongodb() {
    # Comprobar si MongoDB ya está disponible en Docker
    docker ps -a | grep -i mongo &> /dev/null
    if [[ $? -eq 0 ]]; then
        echo "MongoDB ya está disponible en un contenedor Docker."
        return 0  # MongoDB está disponible en Docker
    else
        echo "MongoDB no está disponible en Docker."
        return 1  # MongoDB no está en Docker
    fi
}
# ==============================================
# Función: check_and_install_mongodb
# Descripción:
#   Comprueba si MongoDB está instalado y, si no lo está, lo instala.
#   Luego, verifica si el servicio de MongoDB puede levantarse y realiza
#   una comprobación básica de su funcionamiento.
#
# Parámetros:
#   Ninguno
#
# Uso:
#   check_and_install_mongodb
#
# Autor: José Javier Gutiérrez Gil
# ==============================================
check_and_install_mongodb() {
    # Verificar si MongoDB está instalado
    if ! command -v mongod &>/dev/null; then
        echo "MongoDB no está instalado. Instalándolo ahora..."
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
            # Instalar MongoDB en Linux
            sudo apt-get update
            sudo apt-get install -y mongodb
        elif [[ "$OSTYPE" == "darwin"* ]]; then
            # Instalar MongoDB en macOS
            brew tap mongodb/brew
            brew install mongodb-community@6.0
        else
            echo "Sistema operativo no soportado para instalación automática."
            exit 1
        fi
    else
        echo "MongoDB ya está instalado."
    fi

    # Verificar si el servicio de MongoDB puede levantarse
    if ! pgrep mongod &>/dev/null; then
        echo "Iniciando servicio MongoDB..."
        sudo systemctl start mongod
        sudo systemctl enable mongod
    fi

    # Comprobación básica de funcionamiento
    if mongo --eval "db.runCommand({ ping: 1 })" &>/dev/null; then
        echo "MongoDB está funcionando correctamente."
    else
        echo "Error: MongoDB no se pudo iniciar correctamente."
        exit 1
    fi
}

# ==============================================
# Función: create_database
# Descripción:
#   Solicita al usuario un nombre para una nueva base de datos en MongoDB.
#   Si no se proporciona un nombre, utiliza "MongodbGoChat" como predeterminado.
#   Crea la base de datos y una colección inicial llamada "default_collection".
#
# Parámetros:
#   Ninguno
#
# Uso:
#   create_database
#
# Autor: José Javier Gutiérrez Gil
# ==============================================
create_database() {
    # Solicitar al usuario un nombre para la base de datos o usar el predeterminado
    read -p "Ingrese el nombre de la base de datos (por defecto: MongodbGoChat): " db_name
    db_name=${db_name:-MongodbGoChat}

    # Crear la base de datos
    echo "Creando la base de datos '$db_name'..."
    if mongo --eval "use $db_name; db.createCollection('default_collection');" &>/dev/null; then
        echo "Base de datos '$db_name' creada exitosamente."
    else
        echo "Error al crear la base de datos '$db_name'."
        exit 1
    fi
}
# ==============================================
# Función: install_mongodb_system
# Descripción:
#   Instala MongoDB en el sistema.
#
# Parámetros: Ninguno
#
# Uso:
#   install_mongodb_system
#
# Autor: José Javier Gutiérrez Gil
# ==============================================
function install_mongodb_system() {
    echo "Instalando MongoDB en el sistema..."

    # Instalar MongoDB en el sistema
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Instalar MongoDB en Linux
        sudo apt-get update
        sudo apt-get install -y mongodb
        echo "MongoDB ha sido instalado en el sistema."
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        # Instalar MongoDB en macOS
        brew tap mongodb/brew
        brew install mongodb-community@6.0
        echo "MongoDB ha sido instalado en macOS."
    else
        echo "Sistema operativo no soportado para instalación automática."
        exit 1
    fi

    # Iniciar el servicio de MongoDB
    sudo systemctl start mongodb
    echo "MongoDB ha comenzado a ejecutarse en el sistema."
}
# ==============================================
# Función: install_mongodb_docker
# Descripción:
#   Instala MongoDB en Docker.
#
# Parámetros: Ninguno
#
# Uso:
#   install_mongodb_docker
#
# Autor: José Javier Gutiérrez Gil
# ==============================================
function install_mongodb_docker() {
    echo "Instalando MongoDB en Docker..."

    # Descargar la imagen de MongoDB
    docker pull mongo:latest

    # Crear y ejecutar el contenedor de MongoDB
    docker run -d --name mongodb -p 27017:27017 mongo:latest
    echo "MongoDB está ejecutándose en Docker en el puerto 27017."
}
# ==============================================
# Función: install_mongodb_option
# Descripción:
#   Pregunta al usuario si desea instalar MongoDB en su sistema o en Docker.
#
# Parámetros: Ninguno
#
# Uso:
#   install_mongodb_option
#
# Autor: José Javier Gutiérrez Gil
# ==============================================
function install_mongodb_option() {
    # Verificar si MongoDB está disponible en Docker
    check_docker_mongodb
    if [[ $? -eq 0 ]]; then
        echo "¿Quieres usar MongoDB desde Docker o instalarlo en tu sistema?"
    else
        echo "MongoDB no está disponible en Docker. ¿Quieres instalarlo en tu sistema o en Docker?"
    fi

    # Preguntar al usuario
    read -p "Selecciona una opción: (1) Instalar MongoDB en el sistema (2) Usar MongoDB en Docker: " opcion
    if [[ "$opcion" == "1" ]]; then
        install_mongodb_system
    elif [[ "$opcion" == "2" ]]; then
        install_mongodb_docker
    else
        echo "Opción no válida. El script se detendrá."
        exit 1
    fi
}
# ==============================================
# Función: main
# Descripción:
#   Ejecuta las funciones principales del script.
#
# Parámetros: Ninguno
#
# Uso:
#   main
#
# Autor: José Javier Gutiérrez Gil
# ==============================================
# Función para iniciar el servidor GoChat
function main() {
    # Preguntar si se debe instalar MongoDB en el sistema o en Docker
    install_mongodb_option

    # Ahora, procedemos con el resto del flujo del backend de GoChat
    check_go
    check_dependencies
    check_chat_port
     #check_and_install_mongodb
     #create_database


    echo "Iniciando el servidor GoChat..."
    # El resto de la implementación del backend sigue aquí...
    start_backend
    # restore_env
}


# Ejecutar la función principal
main
