// Package http contiene los adaptadores HTTP (handlers, routers, DTOs)
// Esta es parte de la CAPA DE INFRAESTRUCTURA
package http

// ============================================================================
// DATA TRANSFER OBJECTS (DTOs)
// ============================================================================
//
// Los DTOs son objetos que definen la estructura de datos que se envía/recibe
// a través de HTTP. Son diferentes de las entidades del dominio porque:
//
// 1. Están diseñados para la API REST (no para lógica de negocio)
// 2. Pueden tener validaciones específicas de HTTP
// 3. Pueden simplificar o enriquecer los datos del dominio
// 4. Permiten versionado de la API sin afectar el dominio
//
// ============================================================================

// ============================================================================
// REQUEST DTOs (lo que el cliente envía)
// ============================================================================

// ChatRequest es el DTO para solicitudes de chat desde el cliente HTTP
type ChatRequest struct {
	// Message es el mensaje del usuario
	// validate:"required" podría usarse con librerías de validación
	Message string `json:"message" example:"Explica qué es Go"`
	
	// Model es el modelo de IA a usar (opcional, hay default)
	// omitempty: si está vacío, no se incluye en el JSON
	Model string `json:"model,omitempty" example:"llama-3.3-70b-versatile"`
	
	// Parámetros opcionales avanzados
	Temperature *float64 `json:"temperature,omitempty" example:"0.7"`
	MaxTokens   int      `json:"max_tokens,omitempty" example:"1000"`
}

// ============================================================================
// RESPONSE DTOs (lo que el servidor retorna)
// ============================================================================

// ChatResponse es el DTO para respuestas de chat al cliente HTTP
type ChatResponse struct {
	// Success indica si la operación fue exitosa
	Success bool `json:"success"`
	
	// Message contiene el mensaje de respuesta del modelo
	Message string `json:"message"`
	
	// Model indica qué modelo se usó
	Model string `json:"model"`
	
	// Usage contiene información sobre tokens usados
	Usage *UsageInfo `json:"usage,omitempty"`
	
	// Error contiene el mensaje de error si success=false
	// omitempty: solo se incluye si hay error
	Error string `json:"error,omitempty"`
}

// UsageInfo contiene información sobre el uso de tokens
type UsageInfo struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ModelsResponse es el DTO para la lista de modelos
type ModelsResponse struct {
	Success bool          `json:"success"`
	Models  []ModelInfo   `json:"models,omitempty"`
	Error   string        `json:"error,omitempty"`
}

// ModelInfo contiene información sobre un modelo
type ModelInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OwnedBy string `json:"owned_by"`
}

// ============================================================================
// GENERIC RESPONSE DTOs
// ============================================================================

// ErrorResponse es una respuesta genérica de error
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
}

// SuccessResponse es una respuesta genérica de éxito
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// HealthResponse es la respuesta del endpoint de health check
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
	Service   string `json:"service"`
}

// ============================================================================
// MÉTODOS DE VALIDACIÓN
// ============================================================================

// Validate valida el ChatRequest
// Retorna error si algo está mal
func (r *ChatRequest) Validate() error {
	// Verificar que el mensaje no esté vacío
	if r.Message == "" {
		return ErrEmptyMessage
	}
	
	// Validar temperatura si está presente
	if r.Temperature != nil {
		temp := *r.Temperature
		if temp < 0 || temp > 2 {
			return ErrInvalidTemperature
		}
	}
	
	// Validar max_tokens si está presente
	if r.MaxTokens < 0 {
		return ErrInvalidMaxTokens
	}
	
	return nil
}

// ============================================================================
// ERRORES DE VALIDACIÓN
// ============================================================================

// Definimos errores personalizados para validación
// Estos son específicos de la capa HTTP
var (
	ErrEmptyMessage        = NewValidationError("el mensaje no puede estar vacío")
	ErrInvalidTemperature  = NewValidationError("la temperatura debe estar entre 0 y 2")
	ErrInvalidMaxTokens    = NewValidationError("max_tokens debe ser mayor o igual a 0")
)

// ValidationError es un tipo de error personalizado para validaciones
type ValidationError struct {
	Message string
}

// Error implementa la interfaz error
// Este método es necesario para que ValidationError sea un error
func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError crea un nuevo error de validación
func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

// ============================================================================
// FACTORY FUNCTIONS (funciones para crear DTOs)
// ============================================================================

// NewChatResponse crea una respuesta de chat exitosa
func NewChatResponse(message, model string, usage *UsageInfo) *ChatResponse {
	return &ChatResponse{
		Success: true,
		Message: message,
		Model:   model,
		Usage:   usage,
	}
}

// NewChatErrorResponse crea una respuesta de error de chat
func NewChatErrorResponse(errorMsg string) *ChatResponse {
	return &ChatResponse{
		Success: false,
		Error:   errorMsg,
	}
}

// NewModelsResponse crea una respuesta de modelos exitosa
func NewModelsResponse(models []ModelInfo) *ModelsResponse {
	return &ModelsResponse{
		Success: true,
		Models:  models,
	}
}

// NewModelsErrorResponse crea una respuesta de error de modelos
func NewModelsErrorResponse(errorMsg string) *ModelsResponse {
	return &ModelsResponse{
		Success: false,
		Error:   errorMsg,
	}
}

// NewErrorResponse crea una respuesta genérica de error
func NewErrorResponse(message string, code int) *ErrorResponse {
	return &ErrorResponse{
		Success: false,
		Error:   message,
		Code:    code,
	}
}

// NewHealthResponse crea una respuesta de health check
func NewHealthResponse(status, service string, timestamp int64) *HealthResponse {
	return &HealthResponse{
		Status:    status,
		Service:   service,
		Timestamp: timestamp,
	}
}

// ============================================================================
// CONCEPTOS CLAVE DE GO EXPLICADOS:
// ============================================================================
//
// 1. DTOs (Data Transfer Objects):
//    - Separan la API REST del dominio
//    - Permiten evolucionar la API sin cambiar el dominio
//    - Simplifican las respuestas para el cliente
//
// 2. STRUCT TAGS:
//    - `json:"field_name"`: nombre del campo en JSON
//    - `json:",omitempty"`: omite si está vacío
//    - `example:"value"`: ejemplo para documentación (Swagger)
//    
//    Ejemplo:
//    type User struct {
//        ID   int    `json:"id"`
//        Name string `json:"name,omitempty"`
//    }
//
// 3. PUNTEROS EN CAMPOS OPCIONALES:
//    - *float64 permite nil (campo no presente)
//    - float64 siempre tiene un valor (0.0 por defecto)
//    - Usa punteros para distinguir "no enviado" vs "enviado con 0"
//
// 4. INTERFACE{} (any en Go 1.18+):
//    - interface{} acepta cualquier tipo
//    - Útil para campos genéricos
//    - Pierde type safety, usar con precaución
//
// 5. ERROR INTERFACE:
//    - Cualquier tipo con método Error() string es un error
//    - Permite crear tipos de error personalizados
//    - ValidationError implementa error
//
// 6. FACTORY FUNCTIONS:
//    - Funciones que crean y retornan structs
//    - Mejor que constructores (Go no tiene constructores)
//    - Permiten validación y valores default
//
// 7. VALIDACIÓN:
//    - Validate() es un método que verifica datos
//    - Se llama antes de usar el DTO
//    - Retorna error si algo está mal
//
// 8. COMPOSICIÓN DE RESPUESTAS:
//    - Success + Data + Error
//    - Patrón común en APIs REST
//    - Facilita manejo de errores en el frontend
//
// ============================================================================

// ============================================================================
// MEJORES PRÁCTICAS:
// ============================================================================
//
// 1. SEPARAR DTOs del DOMINIO:
//    - DTOs para HTTP (esta capa)
//    - Entidades para dominio (internal/domain)
//    - Mapear entre ellos en los handlers
//
// 2. VALIDACIÓN EN DTOs:
//    - Validar temprano (en el DTO)
//    - No dejar que datos inválidos lleguen al dominio
//    - Retornar errores claros al cliente
//
// 3. RESPUESTAS CONSISTENTES:
//    - Siempre incluir "success"
//    - Incluir "error" si success=false
//    - Usar estructuras similares para todas las respuestas
//
// 4. DOCUMENTACIÓN:
//    - Usar tags `example` para Swagger
//    - Comentar los campos públicos
//    - Explicar valores posibles
//
// 5. VERSIONADO:
//    - Los DTOs permiten versionar la API
//    - Puedes tener ChatRequestV1 y ChatRequestV2
//    - El dominio permanece estable
//
// ============================================================================
