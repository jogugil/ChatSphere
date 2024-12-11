# Usar una imagen base de Go para la construcción
FROM golang:1.23-bullseye AS builder

# Establecer el directorio de trabajo en la imagen de construcción
RUN mkdir -p /app
WORKDIR /app

# Copiar los archivos necesarios
COPY ./backend/go.mod ./backend/go.sum ./

# Descargar las dependencias
RUN go mod tidy

# Copiar el resto de los archivos del backend
COPY ./backend ./

# Compilar el binario del servidor GoChat
RUN go build -o gochat-server main.go

# Asegurarse de que el binario sea ejecutable
RUN chmod +x gochat-server

# Instalar herramientas adicionales
RUN apt-get update && apt-get install -y netcat-openbsd

# Imagen de producción, solo copiamos el binario compilado
FROM debian:bullseye-slim

# Instalar dependencias mínimas (si es necesario)
RUN apt-get update && apt-get install -y \
    libc6 \
    netcat-openbsd \
    && rm -rf /var/lib/apt/lists/*  # Limpiar el caché de apt para reducir el tamaño

# Copiar el binario desde la imagen de construcción
COPY --from=builder /app/gochat-server /gochat-server

# Verificar el binario en la etapa final (depuración)
RUN ls -l /gochat-server

# Exponer el puerto que usa GoChat
EXPOSE 8081

# Comando para ejecutar la aplicación Go
CMD ["/gochat-server"]
