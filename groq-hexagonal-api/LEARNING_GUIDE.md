# 📚 Guía de Aprendizaje - Go y Arquitectura Hexagonal

Esta guía te llevará paso a paso a través del proyecto para que entiendas cada concepto de Go y la arquitectura hexagonal.

## 🎯 Objetivos de Aprendizaje

Al finalizar, comprenderás:
1. ✅ Conceptos fundamentales de Go
2. ✅ Arquitectura Hexagonal (Ports & Adapters)
3. ✅ Diseño de APIs RESTful
4. ✅ Dependency Injection manual
5. ✅ Manejo de errores en Go
6. ✅ Estructuración de proyectos Go

---

## 📖 Parte 1: Fundamentos de Go

### 1.1 Structs (Estructuras de Datos)

**Lee:** `internal/domain/chat.go` (líneas 1-50)

**Conceptos clave:**
```go
// Definir un struct
type ChatMessage struct {
    Role    string `json:"role"`    // Tags JSON
    Content string `json:"content"`
}

// Crear una instancia
msg := ChatMessage{
    Role:    "user",
    Content: "Hola",
}

// Acceder a campos
fmt.Println(msg.Role) // "user"
```

**Ejercicio:** Crea tu propio struct `User` con campos `Name`, `Email` y `Age`.

### 1.2 Interfaces (Contratos)

**Lee:** `internal/domain/ports.go`

**Conceptos clave:**
```go
// Definir una interfaz
type ChatService interface {
    SendMessage(ctx context.Context, message string) error
}

// Implementar implícitamente (no necesitas declararlo)
type MiServicio struct{}

func (s *MiServicio) SendMessage(ctx context.Context, message string) error {
    // Implementación
    return nil
}
// ¡MiServicio implementa ChatService automáticamente!
```

**Pregunta:** ¿Por qué usamos interfaces en lugar de structs concretos?
<details>
<summary>Respuesta</summary>
Para desacoplar el código. El dominio define QUÉ necesita (interfaz), la infraestructura define CÓMO se hace (implementación). Esto permite cambiar implementaciones sin tocar el dominio.
</details>

### 1.3 Punteros

**Lee:** `internal/application/chat_service.go` (líneas 40-60)

**Conceptos clave:**
```go
// Sin puntero (copia el valor)
func modificar(x int) {
    x = 10  // Solo modifica la copia
}
num := 5
modificar(num)
fmt.Println(num) // Sigue siendo 5

// Con puntero (modifica el original)
func modificarPtr(x *int) {
    *x = 10  // Modifica el original
}
num := 5
modificarPtr(&num)
fmt.Println(num) // Ahora es 10
```

**Regla de oro:** 
- Usa `*` cuando necesites modificar o cuando el struct sea grande
- Usa valor cuando sean datos pequeños e inmutables

### 1.4 Manejo de Errores

**Lee:** `internal/application/chat_service.go` (SendMessage)

**Conceptos clave:**
```go
// Go NO tiene try-catch, usa múltiples retornos
result, err := hacerAlgo()
if err != nil {
    // Manejar el error
    return nil, fmt.Errorf("falló: %w", err)
}
// Usar result

// Crear errores
var ErrNoEncontrado = errors.New("no encontrado")

// Wrapear errores (Go 1.13+)
return fmt.Errorf("operación falló: %w", err)
```

**Ejercicio:** Modifica `SendMessage` para agregar validación de longitud máxima del mensaje.

### 1.5 Context

**Lee:** Cualquier método que use `ctx context.Context`

**Conceptos clave:**
```go
// Crear contexto con timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel() // SIEMPRE cancelar

// Pasar a funciones
resultado, err := miServicio.Hacer(ctx)

// El contexto permite:
// 1. Cancelar operaciones
// 2. Establecer timeouts
// 3. Pasar valores (usar con precaución)
```

---

## 🏗️ Parte 2: Arquitectura Hexagonal

### 2.1 ¿Qué es la Arquitectura Hexagonal?

**Diagrama:**
```
┌─────────────────────────────────────────────┐
│         PUERTOS PRIMARIOS (Entrada)         │
│              ↓                              │
│         HTTP Handlers (adaptadores)          │
└──────────────────┬──────────────────────────┘
                   ↓
┌─────────────────────────────────────────────┐
│         APLICACIÓN (Casos de Uso)           │
│              ChatService                     │
└──────────────────┬──────────────────────────┘
                   ↓
┌─────────────────────────────────────────────┐
│         DOMINIO (Núcleo del negocio)        │
│    Entidades + Interfaces (Ports)           │
└──────────────────┬──────────────────────────┘
                   ↓
┌─────────────────────────────────────────────┐
│      PUERTOS SECUNDARIOS (Salida)           │
│              ↓                              │
│         GroqClient (adaptador)              │
└─────────────────────────────────────────────┘
```

**Principios clave:**
1. El **DOMINIO** es el centro (independiente de todo)
2. Define **INTERFACES** (ports) para comunicarse
3. La **INFRAESTRUCTURA** implementa esas interfaces
4. La **APLICACIÓN** orquesta el dominio

### 2.2 Capas del Proyecto

#### Capa de Dominio (`internal/domain/`)
**Lee:** `chat.go` y `ports.go`

- ✅ Entidades del negocio (ChatMessage, ChatRequest)
- ✅ Interfaces (ChatService, GroqRepository)
- ❌ NO debe conocer HTTP, DB, o cualquier framework
- ❌ NO debe tener imports externos (solo stdlib)

#### Capa de Aplicación (`internal/application/`)
**Lee:** `chat_service.go`

- ✅ Implementa casos de uso (SendMessage, GetModels)
- ✅ Orquesta el dominio
- ✅ Usa las interfaces del dominio
- ❌ NO debe conocer detalles de HTTP o DB

#### Capa de Infraestructura (`internal/infrastructure/`)
**Lee:** `groq/groq_client.go` y `http/handler.go`

- ✅ Implementa las interfaces del dominio
- ✅ Maneja detalles técnicos (HTTP, JSON, etc.)
- ✅ Adaptadores de entrada (HTTP handlers)
- ✅ Adaptadores de salida (GroqClient)

### 2.3 Flujo de una Petición

**Traza este flujo en el código:**

1. **Cliente HTTP** → `POST /api/v1/chat`
2. **Router** (`router.go`) → Detecta la ruta
3. **Handler** (`handler.go:HandleChat`) → Valida y parsea JSON
4. **Service** (`chat_service.go:SendMessage`) → Lógica de negocio
5. **Repository** (`groq_client.go:CreateChatCompletion`) → Llama API externa
6. **Respuesta** → Regresa por el mismo camino

**Ejercicio:** Agrega logs en cada paso para ver el flujo.

---

## 🔧 Parte 3: Patrones y Mejores Prácticas

### 3.1 Dependency Injection

**Lee:** `cmd/api/main.go` (líneas 50-80)

```go
// 1. Crear dependencias (de afuera hacia adentro)
groqClient := groq.NewGroqClient(apiKey, baseURL, timeout)

// 2. Inyectar en el siguiente nivel
chatService := application.NewChatService(groqClient, defaultModel)

// 3. Inyectar en los handlers
chatHandler := http.NewChatHandler(chatService)
```

**Beneficios:**
- ✅ Testeable (puedes inyectar mocks)
- ✅ Flexible (cambiar implementaciones)
- ✅ Explícito (ves todas las dependencias)

### 3.2 DTOs vs Entidades

**Lee:** `infrastructure/http/dto.go`

**DTO (Data Transfer Object):**
- Para comunicación HTTP
- Simplifica/transforma datos del dominio
- Puede tener validaciones específicas de HTTP

**Entidad de Dominio:**
- Representa conceptos del negocio
- Independiente de cómo se transportan
- Contiene lógica de negocio

```go
// DTO (HTTP)
type ChatRequest struct {
    Message string `json:"message"`
    Model   string `json:"model,omitempty"`
}

// Entidad (Dominio)
type ChatRequest struct {
    Messages    []ChatMessage
    Model       string
    Temperature *float64
    MaxTokens   int
}
```

### 3.3 Middleware Pattern

**Lee:** `infrastructure/http/router.go` (loggingMiddleware)

```go
func miMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Antes del handler
        log.Println("Antes")
        
        // Llamar al siguiente handler
        next.ServeHTTP(w, r)
        
        // Después del handler
        log.Println("Después")
    })
}

// Aplicar middleware
router.Use(miMiddleware)
```

**Ejercicio:** Crea un middleware de autenticación que verifique un header `X-API-Key`.

---

## 🚀 Parte 4: Ejecutar y Probar

### 4.1 Configurar el Proyecto

```bash
# 1. Clonar/copiar el proyecto
cd groq-hexagonal-api

# 2. Instalar dependencias
go mod download

# 3. Crear archivo .env
cp .env.example .env
# Editar .env y añadir tu GROQ_API_KEY

# 4. Ejecutar
go run cmd/api/main.go
```

### 4.2 Probar los Endpoints

```bash
# Health check
curl http://localhost:8080/health

# Listar modelos
curl http://localhost:8080/api/v1/models

# Chat
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hola!"}'
```

**O usa el script de ejemplos:**
```bash
chmod +x examples.sh
./examples.sh
```

---

## 🎓 Ejercicios Prácticos

### Ejercicio 1: Agregar Endpoint de History
Crea un endpoint que mantenga el historial de conversaciones.

**Pasos:**
1. Añadir struct `Conversation` en dominio
2. Crear método `GetHistory()` en la interfaz
3. Implementar en el servicio
4. Crear handler y ruta

### Ejercicio 2: Validación Avanzada
Agrega validación de temperatura (0-2) en el DTO.

### Ejercicio 3: Agregar Tests
Crea tests unitarios para `chat_service.go`.

**Hint:**
```go
func TestSendMessage(t *testing.T) {
    // Crear mock del repositorio
    mockRepo := &MockGroqRepository{}
    
    // Crear servicio
    service := NewChatService(mockRepo, "model")
    
    // Probar
    response, err := service.SendMessage(ctx, "test", "model")
    
    // Verificar
    if err != nil {
        t.Errorf("Error inesperado: %v", err)
    }
}
```

---

## 📚 Recursos Adicionales

### Documentación Oficial
- [Go Tour](https://tour.golang.org/) - Tutorial interactivo
- [Effective Go](https://go.dev/doc/effective_go) - Mejores prácticas
- [Go by Example](https://gobyexample.com/) - Ejemplos prácticos

### Arquitectura
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

### Videos Recomendados
- "Hexagonal Architecture in Go" - YouTube
- "Clean Architecture in Go" - YouTube

---

## 🎯 Checklist de Aprendizaje

Marca lo que ya dominas:

- [ ] Entiendo qué son los structs y cómo usarlos
- [ ] Sé crear y usar interfaces
- [ ] Comprendo cuándo usar punteros vs valores
- [ ] Manejo errores correctamente en Go
- [ ] Entiendo qué es context.Context y para qué sirve
- [ ] Comprendo la arquitectura hexagonal
- [ ] Puedo identificar las 3 capas principales
- [ ] Entiendo el flujo de una petición HTTP
- [ ] Sé hacer dependency injection manual
- [ ] Puedo crear nuevos endpoints
- [ ] Entiendo la diferencia entre DTOs y entidades

---

## 🚀 Próximos Pasos

Una vez domines este proyecto:

1. **Frontend:** Crear interfaz web (React/Vue)
2. **WebSockets:** Streaming de respuestas en tiempo real
3. **Database:** Persistir conversaciones (PostgreSQL)
4. **Authentication:** JWT o OAuth2
5. **Testing:** Tests unitarios e integración
6. **CI/CD:** GitHub Actions para deploy automático
7. **Monitoring:** Prometheus + Grafana

---

## 💡 Tips de Aprendizaje

1. **Lee el código en orden:** Dominio → Aplicación → Infraestructura
2. **Modifica y experimenta:** Cambia valores, agrega logs
3. **Rompe cosas:** Elimina código y ve qué falla (aprenderás las dependencias)
4. **Dibuja diagramas:** Visualiza el flujo de datos
5. **Pregunta "¿por qué?":** ¿Por qué interfaces? ¿Por qué punteros?
6. **Compara con otros lenguajes:** ¿Cómo harías esto en Python/Java?

---

## 🤔 Preguntas Frecuentes

**P: ¿Por qué Go no tiene clases?**
R: Go prefiere composición sobre herencia. Los structs + interfaces son más simples y flexibles.

**P: ¿Por qué todo ese código para manejo de errores?**
R: Hace el código más robusto y explícito. Los errores son valores, no excepciones.

**P: ¿Debo usar siempre arquitectura hexagonal?**
R: No. Para proyectos pequeños puede ser overkill. Para proyectos grandes, es invaluable.

**P: ¿Cómo aprendo más sobre Go?**
R: Practica, lee código de proyectos open source, contribuye a la comunidad.

---

¡Buena suerte en tu aprendizaje! 🚀
