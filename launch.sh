#!/bin/bash
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


# ==============================================
# Función: manage_container
# Descripción:
#   Gestiona un contenedor Docker basado en su nombre.
#   - Si el contenedor está en funcionamiento, lo detiene,
#     espera a que se detenga completamente y luego lo elimina.
#   - Si el contenedor está detenido, lo elimina directamente.
#   - Si el contenedor no existe, informa al usuario.
#
# Parámetros:
#   $1: Nombre del contenedor a gestionar.
#
# Uso:
#   manage_container <nombre_del_contenedor>
#
# Notas:
#   - Requiere permisos para ejecutar comandos Docker.
#   - Utiliza un bucle para esperar a que el contenedor se detenga
#     antes de intentar eliminarlo.
#
# Autor: José Javier Gutiérrez Giul
# Emails: jogugil@gmail.com, jogugi@posgrado.upv.es
# ==============================================

manage_container() {
    local container_name=$1
    if [ -z "$container_name" ]; then
        echo "Por favor, proporciona el nombre del contenedor."
        return 1
    fi

    echo "Verificando estado del contenedor '${container_name}'..."

    # Verificar si el contenedor está en ejecución
    if docker ps --filter "name=^${container_name}$" --format '{{.Names}}' | grep -wq "${container_name}"; then
        echo "El contenedor '${container_name}' está en funcionamiento. Deteniéndolo..."
        docker stop "${container_name}"

        # Esperar a que el contenedor se detenga
        echo "Esperando a que el contenedor '${container_name}' se detenga..."
        while docker ps --filter "name=^${container_name}$" --format '{{.Names}}' | grep -wq "${container_name}"; do
            sleep 1
        done
        echo "El contenedor '${container_name}' se ha detenido."

        # Verificar los puertos antes de eliminar el contenedor
        exposed_ports=$(docker inspect --format '{{range .NetworkSettings.Ports}}{{.}}{{end}}' "${container_name}" | cut -d ':' -f 2 | cut -d '/' -f 1)
        liberar_puertos "$exposed_ports"

        # Eliminar contenedor
        docker rm -f "${container_name}"
        echo "Contenedor '${container_name}' eliminado."
    else
        echo "El contenedor '${container_name}' no está en ejecución. Verificando si existe..."
        if docker ps -a --filter "name=^${container_name}$" --format '{{.Names}}' | grep -wq "${container_name}"; then
            echo "Contenedor '${container_name}' detenido. Eliminándolo..."
            docker rm -f "${container_name}"
            echo "Contenedor '${container_name}' eliminado."
        else
            echo "No se encontró el contenedor '${container_name}'."
        fi
    fi
}
# ==============================================
# Función: cleanup_ports
# Descripción:
#   Gestiona contenedores Docker basados en su nombre.
#   - Si el contenedor está en funcionamiento, obtiene los puertos expuestos,
#     los libera, detiene y elimina el contenedor.
#   - Si el contenedor no está en funcionamiento, verifica si existe una imagen
#     con ese nombre, obtiene los puertos expuestos y los libera.
#   - Si no existe contenedor ni imagen, permite al usuario ingresar una lista
#     de puertos para liberarlos.
#
# Parámetros:
#   $1: Nombre del contenedor o imagen a gestionar.
#
# Uso:
#   cleanup_ports <nombre_del_contenedor_o_imagen>
#
# Notas:
#   - Requiere permisos para ejecutar comandos Docker.
#   - Utiliza lsof para liberar los puertos ocupados.
#
# Autor: José Javier Gutiérrez Giul
# Emails: jogugil@gmail.com, jogugi@posgrado.upv.es
# ==============================================
# Función para ver puertos y gestionar contenedores
liberar_puertos() {
    local puertos=$1
    for puerto in $puertos; do
        echo "Liberando puerto $puerto..."
        if lsof -i :$puerto &> /dev/null; then
            PID=$(lsof -t -i :$puerto)
            kill -9 $PID
            echo "Puerto $puerto liberado. Proceso $PID detenido."
        else
            echo "El puerto $puerto ya está libre."
        fi
    done
}
cleanup_ports() {
    local container_name=$1
    echo "Verificando si existe un contenedor con el nombre '${container_name}'..."

    if docker ps --filter "name=^${container_name}$" --format '{{.Names}}' | grep -wq "${container_name}"; then
        echo "El contenedor '${container_name}' está en funcionamiento."

        # Obtener puertos expuestos
        exposed_ports=$(docker inspect --format '{{range .NetworkSettings.Ports}}{{.}}{{end}}' "${container_name}" | cut -d ':' -f 2 | cut -d '/' -f 1)
        liberar_puertos "$exposed_ports"

        # Detener y eliminar contenedor
        docker stop "${container_name}"
        docker rm "${container_name}"
        echo "Contenedor '${container_name}' detenido y eliminado."
    else
        echo "Contenedor '${container_name}' no encontrado. Verificando imagen..."

        if docker image ls --filter "reference=${container_name}" --format '{{.Repository}}' | grep -wq "${container_name}"; then
            echo "Imagen '${container_name}' encontrada."
            exposed_ports=$(docker image inspect "${container_name}" --format '{{json .Config.ExposedPorts}}' | grep -oP '\d+/tcp' | cut -d '/' -f 1)
            liberar_puertos "$exposed_ports"
        else
            echo "No se encontró ni contenedor ni imagen. Solicitará puertos manualmente."
            echo "Proporcione una lista de puertos a liberar (separados por comas):"
            read -p "Puertos: " port_list
            IFS=',' read -r -a ports <<< "$port_list"
            liberar_puertos "${ports[@]}"
        fi
    fi
}
# ==============================================
# Función: manage_image
# Descripción:
#   Gestiona la imagen Docker proporcionada.
#   - Si la imagen existe y es válida, se intenta ejecutar un contenedor
#     de prueba para comprobar su funcionalidad.
#   - Si la imagen no se puede ejecutar, se elimina.
#   - Si la imagen está parada o no está en ejecución, se elimina.
#   - Si hay una nueva versión de la imagen disponible, se pregunta al usuario
#     si desea actualizarla.
#
# Parámetros:
#   $1: Nombre del contenedor (para verificar la imagen asociada).
#   $2: Nombre de la imagen a gestionar.
#
# Uso:
#   manage_image <nombre_del_contenedor> <nombre_de_imagen>
#
# Autor: José Javier Gutiérrez Giul
# Emails: jogugil@gmail.com, jogugi@posgrado.upv.es
# ==============================================

# Función para gestionar la imagen Docker
manage_image() {
    container_name=$1
    image_name=$2

    # Verificar si la imagen existe
    if docker image ls --filter "reference=${image_name}" --format '{{.Repository}}' | grep -wq "${image_name}"; then
        echo "La imagen '${image_name}' está disponible."

        # Verificar integridad de la imagen (comprobar que se puede ejecutar un contenedor)
        if docker inspect "${image_name}" &> /dev/null; then
            echo "La imagen '${image_name}' es válida y puede usarse para crear un contenedor."

            # Intentar ejecutar un contenedor de prueba con esta imagen
            if docker run --rm "${image_name}" echo "Contenedor de prueba ejecutado correctamente"; then
                echo "La imagen '${image_name}' se puede ejecutar correctamente."
            else
                echo "No se pudo ejecutar un contenedor con la imagen '${image_name}'. Eliminando la imagen..."
                docker rmi "${image_name}"
            fi
        else
            echo "La imagen '${image_name}' no es válida. Eliminando la imagen..."
            docker rmi "${image_name}"
        fi

        # Verificar si hay una actualización de la imagen
        echo "Verificando si hay actualizaciones de la imagen '${image_name}'..."
        if docker pull "${image_name}" &> /dev/null; then
            echo "Una nueva versión de la imagen '${image_name}' está disponible."
            read -p "¿Deseas actualizar a la nueva versión? (s/n): " update_response
            if [[ "$update_response" == "s" || "$update_response" == "S" ]]; then
                echo "Actualizando la imagen '${image_name}'..."
                docker pull "${image_name}"
            else
                echo "No se actualizará la imagen."
            fi
        else
            echo "No hay nuevas versiones de la imagen '${image_name}'."
        fi
    else
        echo "No se encontró la imagen '${image_name}' en el repositorio local."
    fi
}
# ==============================================
# Función: setup_docker_image
# Descripción:
#   - Verifica si una imagen Docker está disponible localmente.
#   - Si está disponible, se pregunta si desea actualizarla.
#   - Si no está disponible, la descarga.
#   - Arranca la imagen para comprobar si se puede ejecutar un contenedor con ella.
#   - Verifica el estado de la imagen y puertos.
#   - Si todo es correcto, informa al usuario que la imagen está lista para usarse.
#   - Si falla alguna comprobación, muestra un error y detiene el proceso.
#
# Parámetros:
#   $1: Nombre de la imagen a gestionar.
#
# Uso:
#   setup_docker_image <nombre_imagen>
#
# Autor: José Javier Gutiérrez Giul
# Emails: jogugil@gmail.com, jogugi@posgrado.upv.es
# ==============================================
setup_docker_image() {
    local image_name=$1

    if [ -z "$image_name" ]; then
        echo "Por favor, proporciona el nombre de la imagen."
        return 1
    fi

    # Verificar si la imagen ya está descargada localmente
    if docker image ls --filter "reference=${image_name}" --format '{{.Repository}}' | grep -wq "${image_name}"; then
        echo "La imagen '${image_name}' ya está disponible localmente."

        # Verificar si hay actualizaciones para la imagen
        echo "Verificando si hay una actualización para la imagen '${image_name}'..."
        if docker pull "${image_name}" &> /dev/null; then
            echo "Una nueva versión de la imagen '${image_name}' está disponible."
            read -p "¿Deseas actualizar a la nueva versión? (s/n): " update_response
            if [[ "$update_response" == "s" || "$update_response" == "S" ]]; then
                echo "Actualizando la imagen '${image_name}'..."
                docker pull "${image_name}"
            else
                echo "No se actualizará la imagen."
            fi
        else
            echo "No hay nuevas versiones de la imagen '${image_name}'."
        fi
    else
        # Si la imagen no está disponible localmente, descargarla
        echo "La imagen '${image_name}' no está disponible localmente. Descargándola..."
        if docker pull "${image_name}"; then
            echo "Imagen '${image_name}' descargada con éxito."
        else
            echo "Error al descargar la imagen '${image_name}'."
            return 1
        fi
    fi

    # Intentar arrancar un contenedor con la imagen para verificar su estado
    echo "Intentando arrancar un contenedor de prueba con la imagen '${image_name}'..."
    container_id=$(docker run -d --rm "${image_name}" sleep 30)

    if [ -z "$container_id" ]; then
        echo "Error: No se pudo iniciar un contenedor con la imagen '${image_name}'."
        return 1
    fi

    # Verificar el estado del contenedor
    echo "Verificando el estado del contenedor..."
    container_status=$(docker inspect -f '{{.State.Status}}' "$container_id")

    if [ "$container_status" != "running" ]; then
        echo "El contenedor no se está ejecutando correctamente. Estado: $container_status"
        docker rm -f "$container_id"
        return 1
    fi

    # Verificar que los puertos están correctamente asignados
    echo "Verificando puertos del contenedor..."
    container_ports=$(docker inspect -f '{{range .NetworkSettings.Ports}}{{.}}{{end}}' "$container_id")

    if [ -z "$container_ports" ]; then
        echo "No se asignaron puertos al contenedor. Verifique la configuración de puertos."
        docker rm -f "$container_id"
        return 1
    fi

    echo "La imagen '${image_name}' está lista para usarse. El contenedor está corriendo y los puertos están configurados correctamente."

    # Detener el contenedor de prueba
    docker rm -f "$container_id"
}

# ==============================================
# Función: setup_docker_container
# Descripción:
#   Gestiona un contenedor Docker basado en su nombre.
#   - Verifica si el contenedor está en funcionamiento y los puertos expuestos.
#   - Si el contenedor no está en ejecución, intenta arrancarlo.
#   - Verifica el estado del contenedor y los puertos, si todo está bien, indica que todo está ok.
#   - Si el contenedor no arranca correctamente o el servicio no es accesible, informa del error.
#
# Parámetros:
#   $1: Nombre del contenedor a gestionar.
#   $2: Nombre de la imagen asociada al contenedor.
#   $3: URL del servicio a verificar (si aplica).
#
# Uso:
#   setup_docker_container <nombre_del_contenedor> <nombre_de_la_imagen> [<url_del_servicio>]
#
# Autor: José Javier Gutiérrez Giul
# Emails: jogugil@gmail.com, jogugi@posgrado.upv.es
# ==============================================
setup_docker_image() {
    local image_name=$1

    if [ -z "$image_name" ]; then
        echo "Por favor, proporciona el nombre de la imagen."
        return 1
    fi

    # Verificar si la imagen ya está descargada localmente
    if docker image ls --filter "reference=${image_name}" --format '{{.Repository}}' | grep -wq "${image_name}"; then
        echo "La imagen '${image_name}' ya está disponible localmente."

        # Verificar si hay actualizaciones para la imagen
        echo "Verificando si hay una actualización para la imagen '${image_name}'..."
        if docker pull "${image_name}" &> /dev/null; then
            echo "Una nueva versión de la imagen '${image_name}' está disponible."
            read -p "¿Deseas actualizar a la nueva versión? (s/n): " update_response
            if [[ "$update_response" == "s" || "$update_response" == "S" ]]; then
                echo "Actualizando la imagen '${image_name}'..."
                docker pull "${image_name}"
            else
                echo "No se actualizará la imagen."
            fi
        else
            echo "No hay nuevas versiones de la imagen '${image_name}'."
        fi
    else
        # Si la imagen no está disponible localmente, descargarla
        echo "La imagen '${image_name}' no está disponible localmente. Descargándola..."
        if docker pull "${image_name}"; then
            echo "Imagen '${image_name}' descargada correctamente."
        else
            echo "Error al descargar la imagen '${image_name}'."
        fi
    fi
}
# ==============================================
# Función: priobipal
# Descripción:
#   Lee un archivo de configuración, procesa cada línea y gestiona contenedores, imágenes y puertos.
#
# Parámetros:
#   $1: Ruta al archivo de configuración
#
# Uso:
#   priobipal <archivo_de_configuración>
#
# Autor: José Javier Gutiérrez Giul
# Emails: jogugil@gmail.com, jogugi@posgrado.upv.es
# ==============================================
priobipal() {
    local config_file=$1

    if [ ! -f "$config_file" ]; then
        echo "El archivo de configuración no existe."
        return 1
    fi

    # Leer cada línea del archivo de configuración
    while IFS= read -r linea; do
        # Extraer la imagen, el contenedor, los puertos y el servicio
        # Se asume que la estructura de la línea es: "imagen, contenedor, [puertos], servicio, URL"
        image_name=$(echo "$linea" | cut -d ';' -f 1 | xargs)  # Extraer imagen
        container_name=$(echo "$linea" | cut -d ';' -f 2 | xargs)  # Extraer contenedor
        expose_ports=$(echo "$linea" | cut -d ';' -f 3 | xargs | tr -d '[]')  # Extraer puertos (eliminando corchetes)
        check_service=$(echo "$linea" | cut -d ';' -f 4 | xargs)  # Extraer servicio
        service_url=$(echo "$linea" | cut -d ';' -f 5 | xargs)  # Extraer URL

        # Convertir los puertos en un array
        IFS=', ' read -r -a ports_array <<< "$expose_ports"

        # Mostrar los resultados
        echo "Procesando imagen: $image_name, contenedor: $container_name"
        echo "Puertos: ${ports_array[@]}"
        echo "Comprobar servicio: $check_service, URL del servicio: $service_url"

        # Llamar a la función que manejará cada puerto
        port_args=""
        for puerto in "${ports_array[@]}"; do
            echo "Procesando puerto: $puerto"
            port_args+="-p $puerto "
        done

        # Mostrar la información del contenedor
        echo "Procesando imagen: $image_name, contenedor: $container_name, puertos: $expose_ports, Comprobar servicio: $check_service, URL del servicio: $service_url"

        # Limpiar puertos, contenedor e imagen si es necesario
        echo "Limpiando puertos y contenedor..."
        cleanup_ports "$container_name"

        echo "Eliminando el contenedor..."
        manage_container "$container_name"

        echo "Eliminando la imagen..."
        manage_image "$container_name" "$image_name"

        # Configurar la imagen Docker
        echo "Configurando la imagen Docker..."
        setup_docker_image "$image_name"

        # Configurar el contenedor Docker
        echo "Configurando el contenedor Docker..."
        docker run -d --name "$container_name" $port_args "$image_name"
        echo "Contenedor '$container_name' en ejecución con la imagen '$image_name' en los puertos $expose_ports."

        # Verificar el servicio si es necesario
        if [ "$check_service" == "sí" ] && [ -n "$service_url" ]; then
            echo "Verificando el servicio en $service_url..."
            # Código para verificar si el servicio está funcionando correctamente
            curl --silent --head "$service_url" | grep "200 OK" > /dev/null
            if [ $? -eq 0 ]; then
                echo "El servicio en $service_url está funcionando correctamente."
            else
                echo "El servicio en $service_url no está disponible."
            fi
        fi
    done < "$config_file"
}



# Iniciar proceso

# 1. Comprobar los Docker's usados en el entorno
priobipal dockers_gochat_conf.txt

# Verificar si el paso anterior fue exitoso
if [ $? -ne 0 ]; then
    echo "Hubo un error al procesar los contenedores Docker. Abortando el proceso."
    exit 1
fi

# 2. Comprobar y levantar el servidor GoChat
./start_gochat_server.sh
