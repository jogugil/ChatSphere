# Usar una imagen base de Go para la construcción
FROM golang:1.22 as builder

# Establecer el directorio de trabajo en la imagen de construcción
RUN mkdir -p /app
WORKDIR /app

# Copy the Go module files (go.mod and go.sum)
COPY ./backend/go.mod ./backend/go.sum ./

# Run `go mod tidy` to download dependencies
RUN go mod tidy

# Copy the rest of the backend files
COPY ./backend ./
# Compilar el binario del servidor GoChat
RUN go build -o gochat-server main.go  
RUN apt-get update && apt-get install -y netcat-openbsd
# Imagen de producción, solo copiamos el binario compilado
FROM debian:bullseye-slim

# Instalar dependencias mínimas (si es necesario)
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copiar el binario desde la imagen de construcción
COPY --from=builder /app/gochat-server /gochat-server

# Exponer el puerto que usa GoChat
EXPOSE 8081

# Comando para ejecutar la aplicación Go
CMD ["./gochat-server"]  
