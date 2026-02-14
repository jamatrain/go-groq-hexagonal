// Package domain contiene las entidades y reglas de negocio
// Esta es la CAPA MÁS IMPORTANTE - no depende de nada externo
package domain

import "time"

// ============================================================================
// ENTIDADES DEL DOMINIO
// ============================================================================

// ChatMessage representa un mensaje en una conversación
// En Go, los structs son como clases pero sin herencia
type ChatMessage struct {
	// Role puede ser: "system", "user", o "assistant"
	// La etiqueta `json:"role"` indica cómo se serializa a JSON
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest representa una solicitud de chat completa
type ChatRequest struct {
	// Slice (array dinámico) de mensajes
	// Los slices son una de las estructuras más usadas en Go
	Messages []ChatMessage `json:"messages"`
	
	// Modelo de IA a usar (ej: "llama-3.3-70b-versatile")
	Model string `json:"model"`
	
	// Temperatura controla la creatividad (0.0 - 2.0)
	// El * indica que es un puntero (puede ser nil/null)
	// Se usa puntero para campos opcionales
	Temperature *float64 `json:"temperature,omitempty"`
	
	// Máximo de tokens a generar
	// omitempty significa que si es 0, no se incluye en el JSON
	MaxTokens int `json:"max_tokens,omitempty"`
}

// ChatResponse representa la respuesta de la API de Groq
type ChatResponse struct {
	// ID único de la respuesta
	ID string `json:"id"`
	
	// Tipo de objeto (siempre "chat.completion")
	Object string `json:"object"`
	
	// Timestamp de creación (Unix timestamp)
	Created int64 `json:"created"`
	
	// Modelo usado
	Model string `json:"model"`
	
	// Array de opciones de respuesta (normalmente solo una)
	Choices []Choice `json:"choices"`
	
	// Información de uso de tokens
	Usage Usage `json:"usage"`
}

// Choice representa una opción de respuesta del modelo
type Choice struct {
	// Índice de la opción
	Index int `json:"index"`
	
	// Mensaje de respuesta del asistente
	Message ChatMessage `json:"message"`
	
	// Razón por la que terminó (ej: "stop", "length")
	FinishReason string `json:"finish_reason"`
}

// Usage contiene información sobre tokens usados
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`      // Tokens del input
	CompletionTokens int `json:"completion_tokens"`  // Tokens del output
	TotalTokens      int `json:"total_tokens"`       // Total
}

// Model representa un modelo de IA disponible
type Model struct {
	ID      string    `json:"id"`       // ID del modelo
	Object  string    `json:"object"`   // Tipo de objeto
	Created time.Time `json:"created"`  // Fecha de creación
	OwnedBy string    `json:"owned_by"` // Propietario del modelo
}

// ModelsResponse contiene la lista de modelos disponibles
type ModelsResponse struct {
	Object string  `json:"object"` // Tipo de objeto (ej: "list")
	Data   []Model `json:"data"`   // Array de modelos
}

// ============================================================================
// FUNCIONES DE AYUDA (Helper Functions)
// ============================================================================

// NewChatMessage es una función constructora para crear mensajes
// En Go, las funciones que empiezan con "New" son constructores por convención
func NewChatMessage(role, content string) ChatMessage {
	// Retorna un nuevo ChatMessage con los valores dados
	return ChatMessage{
		Role:    role,
		Content: content,
	}
}

// NewChatRequest crea una nueva solicitud de chat
func NewChatRequest(model string, messages []ChatMessage) ChatRequest {
	return ChatRequest{
		Model:    model,
		Messages: messages,
	}
}

// AddMessage añade un mensaje a la solicitud
// El receiver (c *ChatRequest) es como "this" en otros lenguajes
// El * indica que es un pointer receiver, por lo que puede modificar el struct
func (c *ChatRequest) AddMessage(role, content string) {
	// append() añade elementos a un slice
	c.Messages = append(c.Messages, NewChatMessage(role, content))
}

// SetTemperature configura la temperatura del modelo
func (c *ChatRequest) SetTemperature(temp float64) {
	// &temp crea un puntero a la variable temp
	c.Temperature = &temp
}

// SetMaxTokens configura el máximo de tokens
func (c *ChatRequest) SetMaxTokens(max int) {
	c.MaxTokens = max
}

// GetResponseContent extrae el contenido de la primera respuesta
func (c *ChatResponse) GetResponseContent() string {
	// Verificar que hay al menos una opción
	if len(c.Choices) > 0 {
		return c.Choices[0].Message.Content
	}
	return ""
}

// IsComplete verifica si la respuesta está completa
func (c *ChatResponse) IsComplete() bool {
	// Retorna true si hay opciones y la primera terminó con "stop"
	return len(c.Choices) > 0 && c.Choices[0].FinishReason == "stop"
}

// ============================================================================
// CONCEPTOS CLAVE DE GO EXPLICADOS:
// ============================================================================
//
// 1. STRUCTS: Son tipos de datos compuestos (como clases sin métodos heredados)
//    type Person struct { Name string; Age int }
//
// 2. TAGS DE STRUCT: Las etiquetas `json:"..."` definen cómo se serializa
//    - `json:"name"` -> campo se llama "name" en JSON
//    - `json:"name,omitempty"` -> se omite si está vacío
//
// 3. PUNTEROS (*): Variables que almacenan direcciones de memoria
//    - *float64 es un puntero a float64
//    - &variable obtiene la dirección
//    - *puntero obtiene el valor
//
// 4. SLICES ([]): Arrays dinámicos que pueden crecer
//    - []string es un slice de strings
//    - append(slice, elemento) añade elementos
//    - len(slice) obtiene la longitud
//
// 5. RECEIVERS: Permiten añadir métodos a structs
//    - (c ChatRequest) -> value receiver (no modifica)
//    - (c *ChatRequest) -> pointer receiver (puede modificar)
//
// 6. FUNCIONES CONSTRUCTORAS: Por convención empiezan con "New"
//    - NewChatMessage() crea y retorna un ChatMessage
//
// 7. EXPORTS: En Go, si algo empieza con mayúscula es público (exportado)
//    - ChatMessage es público
//    - chatMessage sería privado (solo dentro del package)
//
// ============================================================================
