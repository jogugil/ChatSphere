#!/bin/bash
# wait-for-it.sh

# Este script espera hasta que el servicio MongoDB (host:port) esté disponible
# Sin parámetros, se utilizarán valores de variables de entorno o predeterminados.

TIMEOUT=30  # Tiempo de espera máximo (en segundos)

# Variables de entorno que puedes definir en tu archivo docker-compose.yml o .env
HOST=${WAIT_FOR_HOST:-"mongodb"}  # Predeterminado "mongodb" si no se especifica (nombre del servicio en Docker Compose)
PORT=${WAIT_FOR_PORT:-"27017"}   # Puerto por defecto de MongoDB (27017)
CMD=${WAIT_FOR_CMD}               # El comando a ejecutar después de esperar

# Función para verificar si el puerto está disponible
wait_for() {
    local elapsed=0
    until echo > /dev/tcp/$HOST/$PORT; do
        if [ $elapsed -ge $TIMEOUT ]; then
            echo "Error: Tiempo de espera agotado. $HOST:$PORT no está disponible."
            exit 1
        fi
        # Reducir los logs de espera solo a una vez si se desea
        if [ $elapsed -eq 0 ] || [ $(($elapsed % 5)) -eq 0 ]; then
            echo "Esperando por $HOST:$PORT..."
        fi
        sleep 1
        ((elapsed++))
    done
    echo "$HOST:$PORT ahora está disponible."
}

# Esperamos hasta el timeout
wait_for

# Ejecutar el comando una vez que el servicio esté disponible
if [ -n "$CMD" ]; then
    echo "Ejecutando el comando: $CMD"
    exec $CMD
else
    echo "No se ha especificado un comando para ejecutar."
fi
