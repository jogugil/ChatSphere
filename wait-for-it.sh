#!/bin/bash

# Variables con valores predeterminados
HOST=${HOST:-"mongodb"}
PORT=${PORT:-27017}
TIMEOUT=${TIMEOUT:-60}
CMD=${CMD:-"/gochat-server"}  # Usa CMD por defecto si no está especificado

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
    exit 1
fi
