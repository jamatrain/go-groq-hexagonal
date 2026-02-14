# 🗂️ Estructura Completa del Proyecto

```
groq-hexagonal-api/
│
├── 📄 README.md                          # Documentación principal del proyecto
├── 📄 LEARNING_GUIDE.md                  # Guía de aprendizaje paso a paso
├── 📄 .env.example                       # Ejemplo de variables de entorno
├── 📄 .gitignore                         # Archivos a ignorar en git
├── 📄 go.mod                             # Definición del módulo y dependencias
├── 📄 Makefile                           # Comandos útiles (make run, make build, etc.)
├── 📄 Dockerfile                         # Imagen Docker para containerización
├── 📄 examples.sh                        # Scripts de ejemplo para probar la API
│
├── 📁 cmd/                               # PUNTO DE ENTRADA
│   └── 📁 api/
│       └── 📄 main.go                    # Función main - ensambla toda la app
│
└── 📁 internal/                          # CÓDIGO PRIVADO (no importable)
    │
    ├── 📁 domain/                        # 🎯 CAPA DE DOMINIO (núcleo)
    │   ├── 📄 chat.go                    # Entidades (ChatMessage, ChatRequest, etc.)
    │   └── 📄 ports.go                   # Interfaces (ChatService, GroqRepository)
    │
    ├── 📁 application/                   # 💼 CAPA DE APLICACIÓN (casos de uso)
    │   └── 📄 chat_service.go            # Lógica de negocio (SendMessage, GetModels)
    │
    ├── 📁 infrastructure/                # 🔌 CAPA DE INFRAESTRUCTURA (adaptadores)
    │   ├── 📁 groq/                      # Adaptador de salida (API externa)
    │   │   └── 📄 groq_client.go         # Cliente HTTP para Groq API
    │   │
    │   └── 📁 http/                      # Adaptador de entrada (HTTP)
    │       ├── 📄 dto.go                 # Data Transfer Objects
    │       ├── 📄 handler.go             # Manejadores HTTP
    │       └── 📄 router.go              # Configuración de rutas
    │
    └── 📁 config/                        # ⚙️ CONFIGURACIÓN
        └── 📄 config.go                  # Carga de variables de entorno
```

---

## 📊 Arquitectura Visual

```
┌─────────────────────────────────────────────────────────────────────┐
│                         CLIENTE HTTP                                │
│                    (Browser, Postman, curl)                         │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ↓
┌─────────────────────────────────────────────────────────────────────┐
│                    INFRASTRUCTURE LAYER (HTTP)                      │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │  router.go          →  Enrutamiento de peticiones           │  │
│  │  handler.go         →  Validación y parseo HTTP             │  │
│  │  dto.go             →  Transformación de datos              │  │
│  └──────────────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ↓
┌─────────────────────────────────────────────────────────────────────┐
│                     APPLICATION LAYER                               │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │  chat_service.go    →  Lógica de negocio                    │  │
│  │                        Casos de uso                          │  │
│  │                        Orquestación                          │  │
│  └──────────────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ↓
┌─────────────────────────────────────────────────────────────────────┐
│                        DOMAIN LAYER                                 │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │  chat.go            →  Entidades del negocio                 │  │
│  │  ports.go           →  Interfaces (contratos)                │  │
│  │                                                               │  │
│  │  ⚡ Núcleo independiente - Sin dependencias externas          │  │
│  └──────────────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ↓
┌─────────────────────────────────────────────────────────────────────┐
│                  INFRASTRUCTURE LAYER (Groq)                        │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │  groq_client.go     →  Cliente HTTP                          │  │
│  │                        Comunicación con API externa          │  │
│  │                        Implementa GroqRepository             │  │
│  └──────────────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ↓
┌─────────────────────────────────────────────────────────────────────┐
│                         GROQ API                                    │
│                  (https://api.groq.com)                             │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 🔄 Flujo de Datos Detallado

### Request (Cliente → Servidor → API Externa)

```
1. Cliente envía:
   POST /api/v1/chat
   Body: {"message": "Hola", "model": "llama-3.3-70b-versatile"}
   
   ↓

2. router.go (línea 40)
   - Detecta ruta /api/v1/chat
   - Aplica middlewares (logging, recovery)
   - Llama a handler.HandleChat()
   
   ↓

3. handler.go (línea 30)
   - Valida método HTTP (POST)
   - Decodifica JSON a dto.ChatRequest
   - Valida el DTO
   - Llama a chatService.SendMessage()
   
   ↓

4. chat_service.go (línea 60)
   - Valida mensaje no vacío
   - Crea domain.ChatRequest
   - Llama a groqRepo.CreateChatCompletion()
   
   ↓

5. groq_client.go (línea 70)
   - Serializa a JSON
   - Crea petición HTTP POST
   - Añade headers (Authorization, Content-Type)
   - Ejecuta request a api.groq.com
   - Parsea respuesta JSON
   
   ↓

6. API de Groq
   - Procesa el mensaje con el modelo LLM
   - Retorna respuesta
```

### Response (API Externa → Servidor → Cliente)

```
6. API de Groq retorna:
   {"id": "...", "choices": [{"message": {"content": "¡Hola! ..."}}], ...}
   
   ↑

5. groq_client.go
   - Recibe JSON
   - Parsea a domain.ChatResponse
   - Retorna al servicio
   
   ↑

4. chat_service.go
   - Valida que hay respuesta
   - Retorna domain.ChatResponse al handler
   
   ↑

3. handler.go
   - Mapea domain.ChatResponse → dto.ChatResponse
   - Serializa a JSON
   - Escribe respuesta HTTP
   
   ↑

2. router.go
   - Aplica middlewares de salida (logging)
   - Añade headers CORS
   
   ↑

1. Cliente recibe:
   200 OK
   Body: {
     "success": true,
     "message": "¡Hola! ¿Cómo puedo ayudarte?",
     "model": "llama-3.3-70b-versatile",
     "usage": {...}
   }
```

---

## 🎯 Responsabilidades de Cada Archivo

### `cmd/api/main.go` (164 líneas)
- ✅ Punto de entrada de la aplicación
- ✅ Ensambla todas las dependencias (DI)
- ✅ Configura el servidor HTTP
- ✅ Maneja graceful shutdown
- ❌ NO contiene lógica de negocio

### `internal/domain/chat.go` (150 líneas)
- ✅ Define entidades del negocio
- ✅ Métodos auxiliares de las entidades
- ✅ Constructores de entidades
- ❌ NO tiene dependencias externas
- ❌ NO conoce HTTP, JSON, o DB

### `internal/domain/ports.go` (120 líneas)
- ✅ Define interfaces (contratos)
- ✅ Documenta qué necesita la aplicación
- ❌ NO implementa nada
- ❌ NO tiene dependencias

### `internal/application/chat_service.go` (180 líneas)
- ✅ Implementa ChatService (puerto primario)
- ✅ Contiene lógica de negocio
- ✅ Orquesta llamadas
- ✅ Valida reglas de negocio
- ❌ NO conoce HTTP
- ❌ NO conoce implementación de GroqRepository

### `internal/infrastructure/groq/groq_client.go` (200 líneas)
- ✅ Implementa GroqRepository (puerto secundario)
- ✅ Maneja comunicación HTTP
- ✅ Serialización/deserialización JSON
- ✅ Manejo de errores HTTP
- ❌ NO contiene lógica de negocio

### `internal/infrastructure/http/dto.go` (150 líneas)
- ✅ Define DTOs para HTTP
- ✅ Validación de entrada
- ✅ Factory functions
- ❌ NO es lo mismo que domain entities

### `internal/infrastructure/http/handler.go` (180 líneas)
- ✅ Maneja peticiones HTTP
- ✅ Valida y parsea JSON
- ✅ Mapea DTOs ↔ Entidades
- ✅ Maneja errores HTTP
- ❌ NO contiene lógica de negocio

### `internal/infrastructure/http/router.go` (160 líneas)
- ✅ Configura rutas
- ✅ Aplica middlewares
- ✅ Configura CORS
- ❌ NO maneja lógica de handlers

### `internal/config/config.go` (120 líneas)
- ✅ Carga variables de entorno
- ✅ Valida configuración
- ✅ Provee defaults
- ❌ NO contiene lógica de negocio

---

## 📦 Dependencias del Proyecto

```go
// go.mod
require (
    github.com/gorilla/mux v1.8.1    // Router HTTP
    github.com/joho/godotenv v1.5.1  // Cargar .env
    github.com/rs/cors v1.10.1       // CORS middleware
)
```

**Nota:** Solo 3 dependencias externas. Go incluye lo demás en su stdlib:
- `net/http` - Servidor y cliente HTTP
- `encoding/json` - JSON serialization
- `context` - Timeouts y cancelaciones
- `errors` - Manejo de errores

---

## 🚀 Comandos Rápidos

```bash
# Instalar dependencias
go mod download

# Ejecutar aplicación
go run cmd/api/main.go

# Compilar binario
go build -o bin/groq-api cmd/api/main.go

# Ejecutar tests
go test ./...

# Formatear código
go fmt ./...

# Ver documentación
go doc internal/domain

# Usar Makefile
make run      # Ejecutar
make build    # Compilar
make test     # Tests
```

---

## 📈 Métricas del Proyecto

- **Total de archivos Go:** 9
- **Líneas de código:** ~1,500
- **Packages:** 6 (main, domain, application, groq, http, config)
- **Interfaces:** 2 (ChatService, GroqRepository)
- **Structs:** 15+
- **Endpoints:** 4 (/, /health, /api/v1/chat, /api/v1/models)

---

## 🎓 Complejidad por Archivo (para aprendizaje)

**🟢 Fácil (empieza aquí):**
1. `internal/domain/chat.go` - Solo structs
2. `internal/config/config.go` - Variables de entorno
3. `internal/domain/ports.go` - Solo interfaces

**🟡 Intermedio:**
4. `internal/infrastructure/http/dto.go` - DTOs y validación
5. `internal/application/chat_service.go` - Lógica simple
6. `internal/infrastructure/http/handler.go` - HTTP handlers

**🔴 Avanzado:**
7. `internal/infrastructure/groq/groq_client.go` - HTTP client
8. `internal/infrastructure/http/router.go` - Routing y middleware
9. `cmd/api/main.go` - Dependency injection

---

Este proyecto está diseñado para aprender Go y arquitectura hexagonal de forma progresiva. ¡Disfruta el viaje! 🚀
