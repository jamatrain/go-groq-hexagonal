# ============================================================================
# DOCKERFILE MULTI-STAGE para Go
# ============================================================================
# 
# Multi-stage build permite:
# 1. Compilar en una imagen con todas las herramientas
# 2. Ejecutar en una imagen mínima solo con el binario
# Resultado: imagen final muy pequeña y segura
#
# ============================================================================

# ============================================================================
# STAGE 1: BUILD
# ============================================================================
# Usar imagen oficial de Go para compilar
FROM golang:1.22-alpine AS builder

# Instalar dependencias del sistema necesarias para compilar
RUN apk add --no-cache git

# Establecer directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias primero (mejor caching)
COPY go.mod go.sum ./

# Descargar dependencias
# Esto se cachea si go.mod y go.sum no cambian
RUN go mod download

# Copiar todo el código fuente
COPY . .

# Compilar la aplicación
# -o: nombre del binario de salida
# -ldflags="-s -w": reduce el tamaño del binario
# CGO_ENABLED=0: compilar sin dependencias C (binario estático)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o groq-api cmd/api/main.go

# ============================================================================
# STAGE 2: RUNTIME
# ============================================================================
# Usar imagen mínima para ejecutar
FROM alpine:latest

# Instalar CA certificates (necesario para HTTPS)
RUN apk --no-cache add ca-certificates

# Crear usuario no-root por seguridad
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Establecer directorio de trabajo
WORKDIR /home/appuser

# Copiar el binario compilado desde la stage anterior
COPY --from=builder /app/groq-api .

# Cambiar ownership del binario
RUN chown appuser:appgroup groq-api

# Cambiar a usuario no-root
USER appuser

# Exponer el puerto
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./groq-api"]

# ============================================================================
# CONCEPTOS EXPLICADOS:
# ============================================================================
#
# 1. MULTI-STAGE BUILD:
#    - Primera etapa: compila (imagen grande con herramientas)
#    - Segunda etapa: ejecuta (imagen pequeña solo con binario)
#    - COPY --from=builder copia entre stages
#
# 2. ALPINE LINUX:
#    - Distribución mínima de Linux (~5MB)
#    - Ideal para containers
#    - Usa apk como package manager
#
# 3. CGO_ENABLED=0:
#    - Desactiva CGO (interoperabilidad con C)
#    - Produce binario estático (sin dependencias dinámicas)
#    - Puede ejecutarse en cualquier Linux
#
# 4. LDFLAGS "-s -w":
#    - -s: quita symbol table
#    - -w: quita debug info
#    - Reduce significativamente el tamaño del binario
#
# 5. CA CERTIFICATES:
#    - Necesarios para hacer peticiones HTTPS
#    - Alpine no los incluye por defecto
#
# 6. USUARIO NO-ROOT:
#    - Buena práctica de seguridad
#    - No ejecutar como root en containers
#
# 7. CACHING DE LAYERS:
#    - COPY go.mod go.sum antes que el resto
#    - go mod download se cachea si no cambian
#    - Acelera builds posteriores
#
# ============================================================================

# ============================================================================
# COMANDOS ÚTILES:
# ============================================================================
#
# Construir imagen:
# docker build -t groq-hexagonal-api .
#
# Ejecutar container:
# docker run -p 8080:8080 --env-file .env groq-hexagonal-api
#
# Ver logs:
# docker logs <container-id>
#
# Entrar al container:
# docker exec -it <container-id> sh
#
# Detener container:
# docker stop <container-id>
#
# ============================================================================
