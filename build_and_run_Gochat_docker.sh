# ==============================================
# Proyecto: Servidor GoChat
# Script: launch_gochat.sh
# Descripción:
#   Este script gestiona la instalación y dependencias de contenedores dockers usados en el chat Gochat.
#   Comprueba si los contendores e imagenes asociadas existen y su integridad es correcta. Sino se reiunstalan
#   y vuelve a cargarse el contenedor, comprobando lso servicios que aportan. Tambiñén se comprueba los puertos
#   expuestos para el Gochat y si el usuario lo desea los libera para volverse a utilizar.
#
#  Despues de tener en marcha los contenedores. Si no hay problemas, se ejecuta el scripot "start_gochat_Server.go"
#  Este script hace un setup de las dependencias gochat y levanta el servidopr en el puerto que indique el usuario
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
#       Ver detalles completos en: https://www.apache.org/licenses/LICENSE-2.0
#
# =======================================================================================================

#!/bin/bash

# Función para verificar si el archivo docker-compose.yml existe
check_docker_compose_file() {
    if [ ! -f "docker-compose.yml" ]; then
        echo "❌ Archivo docker-compose.yml no encontrado. Abortando."
        exit 1
    fi
}

# Función para liberar puertos ocupados
free_ports() {
    local port=$1
    echo "🔍 Verificando si el puerto $port está en uso..."
    if lsof -iTCP:"$port" -sTCP:LISTEN >/dev/null 2>&1; then
        echo "⚠️  Puerto $port en uso. Intentando liberar..."
        # Identificar y detener el contenedor que utiliza este puerto
        container_id=$(docker ps --filter "publish=$port" --format "{{.ID}}")
        if [ -n "$container_id" ]; then
            echo "🛑 Deteniendo contenedor con ID: $container_id"
            docker stop "$container_id"
            docker rm "$container_id"
        else
            echo "⚠️  No se encontró un contenedor usando el puerto $port. Puede ser otro proceso."
            echo "🔍 Identificando proceso externo..."
            pid=$(lsof -iTCP:"$port" -sTCP:LISTEN -t)
            if [ -n "$pid" ]; then
                echo "🛑 Terminando proceso con PID: $pid"
                kill -9 "$pid"
            fi
        fi
    else
        echo "✅ Puerto $port está libre."
    fi
}
# Función para eliminar contenedores conflictivos
remove_container() {
    local container_name=$1
    echo "🔍 Verificando si el contenedor $container_name está en ejecución..."
    container_id=$(docker ps -q --filter "name=$container_name")
    if [ -n "$container_id" ]; then
        echo "⚠️  El contenedor $container_name está en ejecución. Eliminándolo..."
        docker stop "$container_id" && docker rm "$container_id"
    else
        echo "✅ El contenedor $container_name no está en ejecución."
    fi
}

# Paso 1: Comprobar existencia del archivo docker-compose.yml
echo "🔍 Comprobando si el archivo docker-compose.yml existe..."
check_docker_compose_file

# Paso 2: Extraer puertos expuestos del docker-compose.yml
echo "🔍 Extrayendo puertos expuestos del archivo docker-compose.yml..."
ports=$(grep -E '^\s+- "[0-9]+:[0-9]+"' docker-compose.yml | awk -F '"' '{print $2}' | cut -d: -f1 | sort -u)

if [ -z "$ports" ]; then
    echo "⚠️  No se encontraron puertos expuestos en el archivo docker-compose.yml."
else
    echo "📋 Puertos expuestos encontrados: $ports"
    for port in $ports; do
        free_ports "$port"
    done
fi

# Paso 3: Construir y levantar servicios
echo "🚀 Levantando servicios con docker-compose up..."
docker-compose up --build
if [ $? -ne 0 ]; then
    echo "❌ Error al iniciar servicios. Abortando."
    exit 1
fi

# Paso 4: Construir y levantar servicios
echo "🚀 Construyendo imágenes con --no-cache..."
if ! docker-compose build --no-cache; then
    echo "❌ Error en la construcción. Abortando."
    exit 1
fi

echo "🚀 Levantando servicios con docker-compose up..."
if docker-compose up -d; then
    echo "✔️ Servicios iniciados correctamente."
else
    echo "❌ Error al iniciar servicios. Abortando."
    exit 1
fi

# Paso 5: Comprobación final de estado
echo "🔍 Comprobando estado de servicios..."
if docker ps | grep -q "mongo"; then
    echo "✔️ MongoDB está en ejecución."
else
    echo "❌ MongoDB no se ha iniciado correctamente."
fi

if docker ps | grep -q "gochat"; then
    echo "✔️ GoChat está en ejecución."
else
    echo "❌ GoChat no se ha iniciado correctamente."
fi
