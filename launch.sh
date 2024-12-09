#!/bin/bash

# Función para verificar si Docker está instalado
function check_docker() {
    if ! command -v docker &> /dev/null; then
        echo "Docker no está instalado. Instalando Docker..."
        sudo apt-get update
        sudo apt-get install -y docker.io
        sudo systemctl start docker
        sudo systemctl enable docker
    else
        echo "Docker ya está instalado."
    fi
}

# Función para verificar si la imagen MongoDB existe
function check_image() {
    if docker image inspect mongo &> /dev/null; then
        echo "La imagen MongoDB existe. Verificando integridad..."
        if ! docker run --rm mongo mongo --eval "db.stats()" &> /dev/null; then
            echo "La imagen MongoDB no está en buen estado. Eliminando y descargando nuevamente..."
            docker rmi -f mongo
            docker pull mongo
        else
            echo "La imagen MongoDB está en buen estado."
        fi
    else
        echo "La imagen MongoDB no existe. Descargando..."
        docker pull mongo
    fi
}

# Función para verificar si el contenedor MongoDB está levantado y en funcionamiento
function check_container() {
    # Verificamos si el contenedor está corriendo
    if [ "$(docker ps -q -f name=mongo-chat)" ]; then
        echo "El contenedor MongoDB ya está corriendo. Deteniéndolo y eliminándolo..."
        docker stop mongo-chat
        docker rm mongo-chat
    # Verificamos si el contenedor está detenido pero existe
    elif [ "$(docker ps -a -q -f name=mongo-chat)" ]; then
        echo "El contenedor MongoDB existe pero está detenido. Eliminando contenedor detenido..."
        docker rm mongo-chat
    else
        echo "El contenedor MongoDB no existe. Creando contenedor..."
    fi
    
    # Crear el contenedor MongoDB
    docker run --name mongo-chat -d -p 27017:27017 mongo
}
# Función para liberar el puerto 27017 si está ocupado
function check_port() {
    if lsof -i :27017 &> /dev/null; then
        echo "El puerto 27017 está ocupado por otro proceso. Deteniendo el servicio..."
        PID=$(lsof -t -i :27017)
        kill -9 $PID
        echo "Proceso con PID $PID detenido y puerto 27017 liberado."
    else
        echo "El puerto 27017 está libre."
    fi
}

# Función para verificar si el archivo GoChat está en el directorio correcto
function check_gochat() {
    if [ ! -f "backend/main.go" ]; then
        echo "El archivo GoChat no se encuentra en '/backend/main.go'. Asegúrate de que el archivo esté en el directorio correcto."
        exit 1
    else
        echo "El archivo GoChat está presente en '/backend/main.go'."
    fi
}

# Función para verificar si MongoDB está accesible en el puerto 27017
function check_mongo_access() {
    echo "Verificando si MongoDB está accesible en el puerto 27017..."
    if ! docker exec mongo-chat mongosh --eval "db.stats()" &> /dev/null; then
        echo "Error: MongoDB no es accesible en el puerto 27017. Abortando..."
        exit 1
    else
        echo "MongoDB está accesible en el puerto 27017."
    fi
}

# Función para verificar si el puerto del servidor GoChat está en uso
function check_chat_port() {
    read -p "Ingrese el puerto para el servidor GoChat (por defecto 8080): " CHAT_PORT
    CHAT_PORT=${CHAT_PORT:-8080}

    if lsof -i :$CHAT_PORT &> /dev/null; then
        echo "El puerto $CHAT_PORT está en uso por otro proceso."
        read -p "¿Deseas eliminar el proceso que está usando el puerto $CHAT_PORT? (s/n): " eliminar
        if [[ $eliminar =~ ^[Ss]$ ]]; then
            PID=$(lsof -t -i :$CHAT_PORT)
            kill -9 $PID
            echo "Proceso con PID $PID detenido y puerto $CHAT_PORT liberado."
        else
            echo "Por favor, elige otro puerto. El servidor no se puede iniciar en el puerto $CHAT_PORT."
            read -p "Introduce un nuevo puerto para el servidor GoChat: " NUEVO_PUERTO
            CHAT_PORT=$NUEVO_PUERTO
        fi
    else
        echo "El puerto $CHAT_PORT está libre."
    fi
}

# Función para iniciar el servidor GoChat
function start_gochat_server() {
    echo "Iniciando servidor de chat en Go en el puerto $CHAT_PORT..."
    go run backend/main.go -port $CHAT_PORT
}

# Iniciar proceso

# 1. Comprobar si Docker está instalado
check_docker

# 2. Comprobar imagen MongoDB
check_image

# 3. Comprobar si el contenedor MongoDB está levantado y en funcionamiento
check_container

# 4. Comprobar si el puerto 27017 está ocupado y liberarlo
check_port

# 5. Verificar si el archivo de GoChat está en el directorio correcto
check_gochat

# 6. Verificar si MongoDB está accesible en el puerto 27017
check_mongo_access

# 7. Comprobar si el puerto para GoChat está en uso
check_chat_port

# 8. Iniciar el servidor de chat en Go
start_gochat_server
