# ---------- Build Stage ----------
# Usa una imagen oficial de Golang para compilar la aplicación
FROM golang:1.22 AS builder

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos de dependencias (go.mod y go.sum)
COPY go.mod go.sum ./
# Descarga las dependencias antes de copiar el resto del código 
# (mejora el uso del caché)
RUN go mod download

# Copia el resto del código fuente al contenedor
COPY . .

# Compilación estática para evitar dependencia de glibc
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd

# ---------- Final Image ----------
# Usa una imagen ligera de Debian para producción
FROM debian:bullseye-slim

# Instala certificados raíz del sistema (necesarios para conexiones 
# HTTPS seguras)
RUN apt-get update && apt-get install -y \
ca-certificates && rm -rf /var/lib/apt/lists/*

# Establece el directorio de trabajo en la imagen final
WORKDIR /app

# Copiar el binario
COPY --from=builder /app/app .

# Copiar las plantillas views
COPY --from=builder /app/views ./views

# Copiar los archivos CSS
COPY --from=builder /app/css ./css

# Expone el puerto 3030 (para documentación, útil en Docker Compose)
EXPOSE 3030

# Comando por defecto al iniciar el contenedor: ejecuta la app
CMD ["./app"]