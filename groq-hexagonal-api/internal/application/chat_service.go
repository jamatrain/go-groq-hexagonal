// Package application contiene la lógica de negocio (casos de uso)
// Esta capa orquesta el dominio y usa los puertos secundarios
package application

import (
	"context"
	"errors"
	"fmt"
	"groq-hexagonal-api/internal/domain"
)

// ============================================================================
// ERRORES PERSONALIZADOS
// ============================================================================
//
// En Go, es buena práctica definir errores específicos como variables
// Esto permite comparar errores específicos en lugar de strings
//
var (
	ErrEmptyMessage = errors.New("el mensaje no puede estar vacío")
	ErrEmptyModel   = errors.New("el modelo no puede estar vacío")
	ErrAPIFailure   = errors.New("fallo al comunicarse con la API de Groq")
)

// ============================================================================
// IMPLEMENTACIÓN DEL SERVICIO
// ============================================================================

// ChatServiceImpl es la implementación concreta de domain.ChatService
// Nota: No necesitamos declarar explícitamente que implementa la interfaz
// Go lo detecta automáticamente si tiene los métodos correctos
type ChatServiceImpl struct {
	// groqRepo es la dependencia inyectada (puerto secundario)
	// Es una interfaz, no una implementación concreta
	// Esto permite flexibilidad y testing
	groqRepo domain.GroqRepository
	
	// defaultModel es el modelo a usar si no se especifica uno
	defaultModel string
}

// ============================================================================
// CONSTRUCTOR
// ============================================================================

// NewChatService crea una nueva instancia del servicio
// Este es el patrón de Dependency Injection en Go
//
// Parámetros:
//   - repo: implementación del repositorio (inyección de dependencia)
//   - defaultModel: modelo por defecto a usar
//
// Retorna:
//   - domain.ChatService: retornamos la interfaz, no la implementación
//     Esto es una buena práctica: "programa contra interfaces, no implementaciones"
func NewChatService(repo domain.GroqRepository, defaultModel string) domain.ChatService {
	// Validación básica
	if repo == nil {
		// panic() es como throw en otros lenguajes, pero solo para errores irrecuperables
		panic("groqRepo no puede ser nil")
	}
	
	// Retornamos un puntero a la struct
	// El & crea un puntero, similar a "new" en otros lenguajes
	return &ChatServiceImpl{
		groqRepo:     repo,
		defaultModel: defaultModel,
	}
}

// ============================================================================
// IMPLEMENTACIÓN DE LOS MÉTODOS DE LA INTERFAZ
// ============================================================================

// SendMessage implementa el caso de uso de enviar un mensaje
//
// Parámetros:
//   - ctx: contexto para cancelaciones y timeouts
//   - message: mensaje del usuario
//   - model: modelo de IA a usar (vacío = usar default)
//
// Retorna:
//   - *domain.ChatResponse: respuesta del modelo
//   - error: error si algo falla (nil si todo OK)
//
// Nota sobre el receiver (s *ChatServiceImpl):
// - s: nombre de la variable (como "this" o "self")
// - *ChatServiceImpl: tipo del receiver (puntero)
// - Usamos puntero porque el struct tiene campos que queremos acceder
func (s *ChatServiceImpl) SendMessage(
	ctx context.Context,
	message string,
	model string,
) (*domain.ChatResponse, error) {
	// ========================================================================
	// 1. VALIDACIÓN DE ENTRADA
	// ========================================================================
	
	// Validar que el mensaje no esté vacío
	// strings.TrimSpace() elimina espacios al inicio y final
	if len(message) == 0 {
		// Retornamos nil y un error
		// En Go, siempre retornas (nil, error) o (valor, nil)
		return nil, ErrEmptyMessage
	}
	
	// Si no se especificó modelo, usar el default
	if model == "" {
		model = s.defaultModel
	}
	
	// Validar que tengamos un modelo
	if model == "" {
		return nil, ErrEmptyModel
	}
	
	// ========================================================================
	// 2. CONSTRUCCIÓN DE LA PETICIÓN
	// ========================================================================
	
	// Crear el mensaje del usuario
	userMessage := domain.NewChatMessage("user", message)
	
	// Crear la petición de chat con un slice de mensajes
	// []domain.ChatMessage{...} crea un slice con un elemento
	request := domain.NewChatRequest(model, []domain.ChatMessage{userMessage})
	
	// Opcionalmente, podemos configurar parámetros adicionales
	// Descomentar estas líneas si quieres personalizar:
	// request.SetTemperature(0.7)
	// request.SetMaxTokens(1000)
	
	// ========================================================================
	// 3. LLAMADA AL REPOSITORIO (puerto secundario)
	// ========================================================================
	
	// Llamamos al repositorio pasando el contexto y la petición
	// El repositorio se encarga de los detalles de comunicación HTTP
	response, err := s.groqRepo.CreateChatCompletion(ctx, request)
	
	// ========================================================================
	// 4. MANEJO DE ERRORES
	// ========================================================================
	
	// Verificar si hubo error
	if err != nil {
		// fmt.Errorf() crea un nuevo error wrapeando el original
		// %w es el verbo especial para wrap errors (Go 1.13+)
		// Esto permite usar errors.Is() y errors.As() después
		return nil, fmt.Errorf("error al obtener respuesta de Groq: %w", err)
	}
	
	// ========================================================================
	// 5. VALIDACIÓN DE RESPUESTA
	// ========================================================================
	
	// Verificar que la respuesta tenga contenido
	// len() obtiene la longitud de un slice
	if len(response.Choices) == 0 {
		return nil, errors.New("la respuesta no contiene opciones")
	}
	
	// ========================================================================
	// 6. RETORNO EXITOSO
	// ========================================================================
	
	// Todo OK, retornar la respuesta
	return response, nil
}

// GetAvailableModels implementa el caso de uso de listar modelos
//
// Este método es más simple porque solo delega al repositorio
func (s *ChatServiceImpl) GetAvailableModels(ctx context.Context) (*domain.ModelsResponse, error) {
	// Llamar directamente al repositorio
	models, err := s.groqRepo.ListModels(ctx)
	
	// Propagar el error si existe
	if err != nil {
		// fmt.Errorf con %w preserva el error original
		return nil, fmt.Errorf("error al obtener modelos: %w", err)
	}
	
	// Retornar los modelos
	return models, nil
}

// ============================================================================
// CONCEPTOS CLAVE DE GO EXPLICADOS:
// ============================================================================
//
// 1. ERROR HANDLING: Go usa múltiples retornos en lugar de excepciones
//    - Siempre verifica: if err != nil { ... }
//    - Retorna errores explícitamente
//    - Usa fmt.Errorf() con %w para wrapear errores
//
// 2. CONTEXTO (context.Context):
//    - Se pasa como primer parámetro por convención
//    - Permite cancelar operaciones: ctx.Done()
//    - Puede llevar valores: ctx.Value("key")
//    - Propaga timeouts y cancelaciones
//
// 3. PUNTEROS (*):
//    - *ChatServiceImpl es un puntero a ChatServiceImpl
//    - Usamos punteros en receivers para evitar copias
//    - nil es el valor "vacío" de un puntero
//
// 4. DEPENDENCY INJECTION:
//    - Pasamos dependencias en el constructor
//    - Usamos interfaces, no implementaciones concretas
//    - Facilita testing y desacoplamiento
//
// 5. MÉTODOS vs FUNCIONES:
//    - Método: tiene un receiver (s *ChatServiceImpl)
//    - Función: no tiene receiver
//    - Los métodos se llaman: objeto.Metodo()
//
// 6. NIL en Go:
//    - nil es el valor cero para punteros, interfaces, slices, maps, channels
//    - Siempre verifica nil antes de usar: if x != nil { ... }
//
// 7. NAMING CONVENTIONS:
//    - Exportado (público): empieza con mayúscula (ChatService)
//    - No exportado (privado): empieza con minúscula (groqRepo)
//
// 8. ARQUITECTURA HEXAGONAL - CAPA DE APLICACIÓN:
//    - Implementa los casos de uso del negocio
//    - Usa las interfaces del dominio (ports)
//    - No conoce detalles de infraestructura (HTTP, DB, etc.)
//    - Orquesta llamadas entre dominio e infraestructura
//
// ============================================================================

// ============================================================================
// EJEMPLO DE USO:
// ============================================================================
//
// // Crear el servicio
// groqClient := groq.NewGroqClient("tu-api-key", "https://api.groq.com/openai/v1")
// chatService := NewChatService(groqClient, "llama-3.3-70b-versatile")
//
// // Usar el servicio
// ctx := context.Background()
// response, err := chatService.SendMessage(ctx, "Hola, ¿cómo estás?", "")
// if err != nil {
//     log.Fatal(err)
// }
// fmt.Println(response.GetResponseContent())
//
// ============================================================================
