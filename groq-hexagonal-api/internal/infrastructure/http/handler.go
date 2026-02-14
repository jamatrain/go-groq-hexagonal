// Package http - Handlers HTTP
package http

import (
	"encoding/json"
	"groq-hexagonal-api/internal/domain"
	"log"
	"net/http"
	"time"
)

// ============================================================================
// HANDLER STRUCT
// ============================================================================

// ChatHandler maneja las peticiones HTTP relacionadas con chat
// Actúa como adaptador entre HTTP y el servicio de aplicación
type ChatHandler struct {
	// chatService es la dependencia del servicio de aplicación
	// Usamos la interfaz, no la implementación concreta
	chatService domain.ChatService
}

// ============================================================================
// CONSTRUCTOR
// ============================================================================

// NewChatHandler crea un nuevo handler con el servicio inyectado
func NewChatHandler(service domain.ChatService) *ChatHandler {
	if service == nil {
		panic("chatService no puede ser nil")
	}
	
	return &ChatHandler{
		chatService: service,
	}
}

// ============================================================================
// HTTP HANDLERS
// ============================================================================

// HandleChat maneja POST /api/v1/chat
// Este es el endpoint principal para enviar mensajes
//
// En Go, los handlers HTTP tienen esta firma:
// func(w http.ResponseWriter, r *http.Request)
//   - w: para escribir la respuesta
//   - r: contiene la petición del cliente
func (h *ChatHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	// ========================================================================
	// 1. LOGGING (opcional pero recomendado)
	// ========================================================================
	log.Printf("[%s] %s - HandleChat", r.Method, r.URL.Path)
	
	// ========================================================================
	// 2. VALIDAR MÉTODO HTTP
	// ========================================================================
	
	// Verificar que sea POST
	if r.Method != http.MethodPost {
		// Escribir error con status 405 Method Not Allowed
		h.writeErrorResponse(w, "método no permitido", http.StatusMethodNotAllowed)
		return
	}
	
	// ========================================================================
	// 3. DECODIFICAR EL BODY JSON
	// ========================================================================
	
	// Crear una variable para el DTO
	var req ChatRequest
	
	// json.NewDecoder() lee del body de la petición
	// .Decode(&req) parsea el JSON a la struct
	// &req es un puntero porque Decode necesita modificar el struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// Cerrar el body (buena práctica)
	// defer lo ejecuta al final de la función
	defer r.Body.Close()
	
	// ========================================================================
	// 4. VALIDAR EL REQUEST
	// ========================================================================
	
	if err := req.Validate(); err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// ========================================================================
	// 5. LLAMAR AL SERVICIO DE APLICACIÓN
	// ========================================================================
	
	// r.Context() obtiene el contexto de la petición HTTP
	// Este contexto se cancela automáticamente si el cliente cierra la conexión
	ctx := r.Context()
	
	// Llamar al servicio con el mensaje y modelo
	response, err := h.chatService.SendMessage(ctx, req.Message, req.Model)
	if err != nil {
		// Error del servicio -> 500 Internal Server Error
		log.Printf("Error en servicio: %v", err)
		h.writeErrorResponse(w, "error al procesar el mensaje", http.StatusInternalServerError)
		return
	}
	
	// ========================================================================
	// 6. MAPEAR DOMINIO → DTO
	// ========================================================================
	
	// Convertir la respuesta del dominio a DTO HTTP
	chatResponse := NewChatResponse(
		response.GetResponseContent(),
		response.Model,
		&UsageInfo{
			PromptTokens:     response.Usage.PromptTokens,
			CompletionTokens: response.Usage.CompletionTokens,
			TotalTokens:      response.Usage.TotalTokens,
		},
	)
	
	// ========================================================================
	// 7. ESCRIBIR LA RESPUESTA JSON
	// ========================================================================
	
	h.writeJSONResponse(w, chatResponse, http.StatusOK)
}

// HandleGetModels maneja GET /api/v1/models
// Retorna la lista de modelos disponibles
func (h *ChatHandler) HandleGetModels(w http.ResponseWriter, r *http.Request) {
	// ========================================================================
	// 1. LOGGING
	// ========================================================================
	log.Printf("[%s] %s - HandleGetModels", r.Method, r.URL.Path)
	
	// ========================================================================
	// 2. VALIDAR MÉTODO
	// ========================================================================
	
	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, "método no permitido", http.StatusMethodNotAllowed)
		return
	}
	
	// ========================================================================
	// 3. LLAMAR AL SERVICIO
	// ========================================================================
	
	ctx := r.Context()
	response, err := h.chatService.GetAvailableModels(ctx)
	if err != nil {
		log.Printf("Error al obtener modelos: %v", err)
		h.writeErrorResponse(w, "error al obtener modelos", http.StatusInternalServerError)
		return
	}
	
	// ========================================================================
	// 4. MAPEAR A DTO
	// ========================================================================
	
	// Convertir []domain.Model a []ModelInfo
	modelInfos := make([]ModelInfo, len(response.Data))
	for i, model := range response.Data {
		modelInfos[i] = ModelInfo{
			ID:      model.ID,
			Name:    model.ID, // Usamos el ID como nombre
			OwnedBy: model.OwnedBy,
		}
	}
	
	modelsResponse := NewModelsResponse(modelInfos)
	
	// ========================================================================
	// 5. ESCRIBIR RESPUESTA
	// ========================================================================
	
	h.writeJSONResponse(w, modelsResponse, http.StatusOK)
}

// HandleHealth maneja GET /health
// Endpoint simple para verificar que la API está funcionando
func (h *ChatHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	// Crear respuesta de health check
	health := NewHealthResponse("healthy", "groq-api", time.Now().Unix())
	
	// Escribir respuesta
	h.writeJSONResponse(w, health, http.StatusOK)
}

// ============================================================================
// MÉTODOS AUXILIARES (helpers)
// ============================================================================

// writeJSONResponse escribe una respuesta JSON
// Es un método privado (empieza con minúscula)
func (h *ChatHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	// Establecer Content-Type
	w.Header().Set("Content-Type", "application/json")
	
	// Establecer status code
	w.WriteHeader(statusCode)
	
	// Serializar y escribir JSON
	// json.NewEncoder() crea un encoder que escribe directamente a w
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Si falla la serialización, registrar el error
		log.Printf("Error al escribir JSON: %v", err)
	}
}

// writeErrorResponse escribe una respuesta de error
func (h *ChatHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	// Crear el DTO de error
	errorResponse := NewErrorResponse(message, statusCode)
	
	// Escribir la respuesta
	h.writeJSONResponse(w, errorResponse, statusCode)
}

// ============================================================================
// CONCEPTOS CLAVE DE GO EXPLICADOS:
// ============================================================================
//
// 1. HTTP HANDLERS EN GO:
//    - Firma: func(w http.ResponseWriter, r *http.Request)
//    - w: para escribir la respuesta (WriteHeader, Write, etc.)
//    - r: contiene la petición (Method, URL, Body, Headers, etc.)
//
// 2. http.ResponseWriter:
//    - Interfaz para escribir respuestas HTTP
//    - Header(): obtiene/modifica headers
//    - WriteHeader(code): establece el status code
//    - Write([]byte): escribe el body
//
// 3. *http.Request:
//    - Method: string con el método HTTP (GET, POST, etc.)
//    - URL: *url.URL con la URL de la petición
//    - Body: io.ReadCloser para leer el body
//    - Header: http.Header con los headers
//    - Context(): obtiene el contexto de la petición
//
// 4. JSON EN GO:
//    - json.NewDecoder(r).Decode(&v): parsea JSON del reader
//    - json.NewEncoder(w).Encode(v): escribe JSON al writer
//    - Alternativas: json.Marshal() y json.Unmarshal()
//
// 5. DEFER:
//    - defer r.Body.Close() asegura que el body se cierre
//    - Se ejecuta al final de la función
//    - Importante para evitar memory leaks
//
// 6. CONTEXT:
//    - r.Context() obtiene el contexto de la petición
//    - Se cancela si el cliente cierra la conexión
//    - Debe pasarse a todas las llamadas que puedan tardar
//
// 7. LOGGING:
//    - log.Printf() escribe a stdout con timestamp
//    - Útil para debugging y monitoring
//    - En producción, usa librerías más robustas (zap, logrus)
//
// 8. ERROR HANDLING:
//    - Siempre verifica errores: if err != nil { ... }
//    - Retorna respuestas HTTP apropiadas
//    - Log errores internos, no los expongas al cliente
//
// 9. STATUS CODES HTTP:
//    - 200 OK: éxito
//    - 400 Bad Request: error del cliente (validación)
//    - 405 Method Not Allowed: método HTTP incorrecto
//    - 500 Internal Server Error: error del servidor
//
// 10. DEPENDENCY INJECTION:
//     - El handler recibe el servicio en el constructor
//     - No crea dependencias, las recibe
//     - Facilita testing y desacoplamiento
//
// ============================================================================

// ============================================================================
// FLUJO DE UNA PETICIÓN:
// ============================================================================
//
// 1. Cliente envía: POST /api/v1/chat
//    Body: {"message": "Hola", "model": "llama-3.3-70b-versatile"}
//
// 2. Router llama a: handler.HandleChat(w, r)
//
// 3. Handler:
//    a. Valida método HTTP (POST)
//    b. Decodifica JSON a ChatRequest
//    c. Valida el ChatRequest
//    d. Llama a chatService.SendMessage()
//    e. Mapea domain.ChatResponse a http.ChatResponse
//    f. Escribe JSON al cliente
//
// 4. Cliente recibe:
//    {
//      "success": true,
//      "message": "Hola! ¿Cómo puedo ayudarte?",
//      "model": "llama-3.3-70b-versatile",
//      "usage": {
//        "prompt_tokens": 10,
//        "completion_tokens": 20,
//        "total_tokens": 30
//      }
//    }
//
// ============================================================================

// ============================================================================
// MEJORES PRÁCTICAS:
// ============================================================================
//
// 1. VALIDACIÓN TEMPRANA:
//    - Validar método HTTP primero
//    - Validar JSON antes de procesar
//    - Retornar errores claros al cliente
//
// 2. MANEJO DE ERRORES:
//    - Log errores internos
//    - No expongas detalles internos al cliente
//    - Usa códigos HTTP apropiados
//
// 3. CONTEXTO:
//    - Siempre pasa r.Context() al servicio
//    - Respeta cancelaciones del cliente
//    - Establece timeouts si es necesario
//
// 4. LOGGING:
//    - Log todas las peticiones
//    - Log errores con contexto
//    - En producción, usa structured logging
//
// 5. SEPARACIÓN DE RESPONSABILIDADES:
//    - Handler solo maneja HTTP
//    - Lógica de negocio en el servicio
//    - Mapeo entre DTOs y dominio en el handler
//
// ============================================================================
