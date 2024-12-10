
# ==============================================
# Proyecto: Servidor GoChat
# Descripción: setup_docker_gochat.sh
#
# Script de instalación y configuración de Docker y Docker Compose para el proyecto GoChat.
#
# Este script comprueba si Docker y Docker Compose están instalados en el sistema. Si no están presentes, 
# ofrece la opción de instalarlos. Además, realiza las siguientes acciones:
#
# 1. Verifica la instalación de Docker y Docker Compose.
# 2. Si Docker Compose no está instalado o tiene problemas, lo reinstala.
# 3. Ejecuta el comando `docker-compose up --build` para construir y ejecutar los contenedores
#    definidos en el archivo `docker-compose.yml` de este proyecto.
#
# Requiere permisos de sudo para instalar software y ejecutar ciertos comandos.
#
# Uso:
# - Ejecutar el script como superusuario para asegurar que Docker y Docker Compose puedan instalarse
#   correctamente.
# - El script también puede reiniciar Docker Compose en caso de que no funcione correctamente.
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


# Función para verificar si Docker está instalado
check_docker() {
    if command -v docker &> /dev/null
    then
        echo "Docker está instalado."
        return 0  # Docker está instalado correctamente
    else
        echo "Docker no está instalado."
        return 1  # Docker no está instalado
    fi
}

# Función para verificar si Docker Compose está instalado y funcionando correctamente
# Función para verificar si Docker Compose está instalado y funcionando correctamente
check_docker_compose() {
    if command -v docker-compose &> /dev/null  # Verifica si docker-compose está instalado
    then
        echo "Docker Compose está instalado."
        docker-compose --version &> /dev/null  # Verifica si Docker Compose puede ejecutarse
        if [ $? -eq 0 ]; then  # Si Docker Compose funciona correctamente
            echo "Docker Compose está funcionando correctamente."
            return 0  # Docker Compose está funcionando correctamente
        else  # Si Docker Compose no funciona correctamente
            echo "Docker Compose no está funcionando correctamente. Se procederá a reinstalar."
            return 2  # Docker Compose no está funcionando correctamente
        fi
    else  # Si Docker Compose no está instalado
        echo "Docker Compose no está instalado. Se procederá a instalar."
        return 1  # Docker Compose no está instalado
    fi
}
# Función para desinstalar Docker
uninstall_docker() {
    echo "Eliminando Docker..."
    sudo apt-get purge -y docker.io  # Elimina Docker
    sudo apt-get autoremove -y  # Elimina dependencias no necesarias
    echo "Docker eliminado correctamente."
}
# Función para desinstalar Docker Compose
uninstall_docker_compose() {
    echo "Eliminando Docker Compose..."
    sudo rm -f /usr/local/bin/docker-compose  # Elimina Docker Compose
    echo "Docker Compose eliminado correctamente."
}

#Funcion para instalar Docker
install_docker() {
    echo "Instalando Docker..."
    sudo apt-get update
    sudo apt-get install -y docker.io
    sudo systemctl enable --now docker
    echo "Docker instalado correctamente."
}

# Función para instalar jq (si no está presente)
install_jq() {
    echo "Instalando jq..."
    sudo apt-get install -y jq
    echo "jq instalado correctamente."
}

# Función para instalar Docker Compose
install_docker_compose() {
    echo "Instalando Docker Compose..."

    # Instalar jq si no está presente
    if ! command -v jq &> /dev/null
    then
        install_jq
    fi

    # Obtener la última versión de Docker Compose
    LATEST_VERSION=$(curl -s https://api.github.com/repos/docker/compose/releases/latest | jq -r .tag_name)
    sudo curl -L "https://github.com/docker/compose/releases/download/$LATEST_VERSION/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

    # Darle permisos de ejecución
    sudo chmod +x /usr/local/bin/docker-compose
    echo "Docker Compose instalado correctamente."
}
# Función principal que coordina la verificación e instalación
setup_environment() {
    # Verificar Docker
    check_docker
    if [ $? -ne 0 ]; then
        read -p "¿Quieres instalar Docker? (Y/n): " INSTALL_DOCKER
        if [[ "$INSTALL_DOCKER" == "Y" || "$INSTALL_DOCKER" == "y" || "$INSTALL_DOCKER" == "" ]]; then
            install_docker
        else
            echo "Docker no se instalará. Saliendo..."
            exit 1
        fi
    fi

    # Verificar Docker Compose
    check_docker_compose
    result=$?  # Guardar el código de salida de check_docker_compose

    if [ $result -eq 2 ]; then
        # Si Docker Compose no funciona correctamente, eliminar e instalar de nuevo
        echo "Docker Compose no está funcionando correctamente. Se procederá a eliminar e instalar de nuevo."
        remove_docker_compose
        install_docker_compose
    elif [ $result -eq 1 ]; then
        # Si Docker Compose no está instalado, preguntar si desea instalarlo
        read -p "¿Quieres instalar Docker Compose? (Y/n): " INSTALL_COMPOSE
        if [[ "$INSTALL_COMPOSE" == "Y" || "$INSTALL_COMPOSE" == "y" || "$INSTALL_COMPOSE" == "" ]]; then
            install_docker_compose
        else
            echo "Docker Compose no se instalará. Saliendo..."
            exit 1
        fi
    else
        # Si Docker Compose está instalado y funcionando correctamente, no hacer nada
        echo "Docker Compose ya está funcionando correctamente. No se requiere instalación."
    fi

    # Verificar que todo está correctamente instalado y funcionando
    docker --version
    docker-compose --version
}
# Función para ejecutar docker-compose up --build
run_docker_compose() {
    echo "Ejecutando docker-compose up --build..."
    docker-compose up --build
}
# Función para ejecutar docker-compose up --build
run_docker_compose() {
    echo "Ejecutando docker-compose up --build..."
    docker-compose up --build
}

# Paso 1: Configuración del entorno
setup_environment

# Paso 2: Ejecutar docker-compose
run_docker_compose
