#!/bin/bash

# Nombre de la imagen de MongoDB
IMAGE="mongo"

# Nombre del contenedor
CONTAINER_NAME="mongo-container"

# Puerto que se expondrá
HOST_PORT=27017
CONTAINER_PORT=27017

# Volumen para persistir datos
VOLUME_PATH="./mongo-data"
mkdir -p $VOLUME_PATH

# Comando para lanzar MongoDB en Docker
docker run --name $CONTAINER_NAME -d \
  -p $HOST_PORT:$CONTAINER_PORT \
  -v $VOLUME_PATH:/data/db \
  $IMAGE

# Verificación de que el contenedor está corriendo
echo "MongoDB está corriendo en el contenedor $CONTAINER_NAME en el puerto $HOST_PORT"
docker ps -f name=$CONTAINER_NAME
