# 🏗️ API Groq con Arquitectura Hexagonal en Go

Este proyecto implementa una API RESTful en Go que consume la API de Groq usando **Arquitectura Hexagonal** (también conocida como Ports & Adapters).

## 📚 ¿Qué es la Arquitectura Hexagonal?

La arquitectura hexagonal separa la lógica de negocio (dominio) de los detalles de implementación (infraestructura). Beneficios:

- **Independencia**: El dominio no depende de frameworks o tecnologías externas
- **Testeable**: Fácil de testear cada capa por separado
- **Mantenible**: Los cambios en una capa no afectan a las demás
- **Flexible**: Puedes cambiar la base de datos, el framework HTTP, etc., sin tocar el dominio

## 🎯 Estructura del Proyecto

```
groq-hexagonal-api/
├── cmd/
│   └── api/
│       └── main.go                 # Punto de entrada de la aplicación
├── internal/
│   ├── domain/                     # CAPA DE DOMINIO (núcleo del negocio)
│   │   ├── chat.go                 # Entidad Chat
│   │   └── ports.go                # Interfaces (contratos)
│   ├── application/                # CAPA DE APLICACIÓN (casos de uso)
│   │   └── chat_service.go         # Lógica de negocio
│   ├── infrastructure/             # CAPA DE INFRAESTRUCTURA (detalles técnicos)
│   │   ├── groq/
│   │   │   └── groq_client.go      # Cliente HTTP para Groq API
│   │   └── http/
│   │       ├── handler.go          # Manejadores HTTP
│   │       ├── router.go           # Configuración de rutas
│   │       └── dto.go              # Data Transfer Objects
│   └── config/
│       └── config.go               # Configuración de la app
├── .env.example                    # Ejemplo de variables de entorno
├── go.mod                          # Dependencias del proyecto
└── README.md                       # Este archivo
```

## 🔄 Flujo de una Petición

```
1. Cliente HTTP → 2. Handler (HTTP) → 3. Service (Aplicación) → 4. GroqClient (Infraestructura) → 5. API Groq
                ↓                       ↓                          ↓
                HTTP                    Dominio                    HTTP Client
```

## 🚀 Instalación

```bash
# Clonar el repositorio
git clone <tu-repo>
cd groq-hexagonal-api

# Instalar dependencias
go mod download

# Configurar variables de entorno
cp .env.example .env
# Edita .env y añade tu GROQ_API_KEY

# Ejecutar la aplicación
go run cmd/api/main.go
```

## 📡 Endpoints Disponibles

### 1. Chat Completion
```bash
POST /api/v1/chat
Content-Type: application/json

{
  "message": "Explica qué es Go en 3 líneas",
  "model": "llama-3.3-70b-versatile"
}
```

### 2. Listar Modelos
```bash
GET /api/v1/models
```

### 3. Health Check
```bash
GET /health
```

## 🧪 Ejemplos de Uso

```bash
# Chat con el modelo
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "¿Qué es la arquitectura hexagonal?",
    "model": "llama-3.3-70b-versatile"
  }'

# Listar modelos disponibles
curl http://localhost:8080/api/v1/models
```

## 📖 Conceptos de Go Explicados en el Código

- **Structs**: Estructuras de datos personalizadas
- **Interfaces**: Contratos que definen comportamiento
- **Pointer Receivers**: Métodos que pueden modificar el struct
- **Error Handling**: Manejo explícito de errores
- **Context**: Propagación de timeouts y cancelaciones
- **Goroutines**: Concurrencia (si se implementa)
- **Channels**: Comunicación entre goroutines
- **Defer**: Ejecución diferida de funciones

## 🎓 Orden de Lectura Recomendado

Para aprender el proyecto, lee los archivos en este orden:

1. `internal/domain/chat.go` - Entidades del dominio
2. `internal/domain/ports.go` - Interfaces (contratos)
3. `internal/application/chat_service.go` - Lógica de negocio
4. `internal/infrastructure/groq/groq_client.go` - Cliente HTTP
5. `internal/infrastructure/http/dto.go` - DTOs
6. `internal/infrastructure/http/handler.go` - Manejadores
7. `internal/infrastructure/http/router.go` - Rutas
8. `internal/config/config.go` - Configuración
9. `cmd/api/main.go` - Punto de entrada

## 🌟 Próximos Pasos

Después de dominar esta API, avanzaremos a:
- Frontend con interfaz gráfica (React/Vue/Svelte)
- Websockets para streaming de respuestas
- Base de datos para guardar conversaciones
- Sistema de autenticación
- Caché con Redis
- Tests unitarios e integración

## 📚 Recursos

- [Documentación de Go](https://go.dev/doc/)
- [Groq API Docs](https://console.groq.com/docs)
- [Arquitectura Hexagonal](https://alistair.cockburn.us/hexagonal-architecture/)
