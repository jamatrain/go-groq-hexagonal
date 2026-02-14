// Package domain - Continuación con las interfaces (Ports)
package domain

import "context"

// ============================================================================
// PORTS (INTERFACES)
// ============================================================================
//
// En arquitectura hexagonal, los "ports" son interfaces que definen
// cómo el dominio se comunica con el exterior (sin conocer los detalles)
//
// Beneficios:
// - El dominio no depende de implementaciones concretas
// - Podemos cambiar la implementación sin tocar el dominio
// - Facilita el testing (podemos crear mocks fácilmente)
// ============================================================================

// ChatService define los casos de uso de negocio
// Esta es una interfaz de PUERTO PRIMARIO (driving port)
// Los puertos primarios son invocados por actores externos (ej: HTTP handlers)
type ChatService interface {
	// SendMessage envía un mensaje y obtiene respuesta del modelo
	// context.Context permite cancelaciones, timeouts y propagación de valores
	// error es el tipo estándar de Go para manejar errores
	SendMessage(ctx context.Context, message string, model string) (*ChatResponse, error)
	
	// GetAvailableModels obtiene la lista de modelos disponibles
	GetAvailableModels(ctx context.Context) (*ModelsResponse, error)
}

// GroqRepository define cómo accedemos a la API de Groq
// Esta es una interfaz de PUERTO SECUNDARIO (driven port)
// Los puertos secundarios son implementados por adaptadores externos
type GroqRepository interface {
	// CreateChatCompletion realiza una petición de chat completion
	CreateChatCompletion(ctx context.Context, request ChatRequest) (*ChatResponse, error)
	
	// ListModels obtiene todos los modelos disponibles
	ListModels(ctx context.Context) (*ModelsResponse, error)
}

// ============================================================================
// CONCEPTOS CLAVE DE GO - INTERFACES
// ============================================================================
//
// 1. INTERFACES EN GO: Son tipos que definen un conjunto de métodos
//    - Son "satisfechas implícitamente" (no necesitas declarar que implementas una interfaz)
//    - Si un tipo tiene todos los métodos de una interfaz, la implementa automáticamente
//
//    Ejemplo:
//    type Writer interface {
//        Write([]byte) (int, error)
//    }
//    
//    // Cualquier tipo con este método implementa Writer automáticamente:
//    func (f *File) Write(data []byte) (int, error) { ... }
//
// 2. CONTEXT: Es un tipo especial en Go para:
//    - Cancelar operaciones de larga duración
//    - Establecer timeouts
//    - Pasar valores entre funciones (con precaución)
//
//    Uso común:
//    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//    defer cancel() // Siempre cancelar para liberar recursos
//
// 3. ERROR HANDLING: Go no tiene excepciones, usa múltiples valores de retorno
//    - Por convención, el error es el último valor de retorno
//    - nil significa "sin error"
//    - Siempre debes verificar errores explícitamente
//
//    Ejemplo:
//    result, err := SomeFunction()
//    if err != nil {
//        // Manejar el error
//        return nil, err
//    }
//    // Usar result
//
// 4. PUNTEROS EN RETORNOS: Se usan para:
//    - Evitar copiar estructuras grandes
//    - Permitir valores nil (para indicar "no hay resultado")
//    
//    *ChatResponse puede ser nil o apuntar a un ChatResponse
//
// 5. CONVENCIONES DE NOMBRES:
//    - Interfaces con un solo método terminan en "-er": Reader, Writer, Stringer
//    - Interfaces con varios métodos usan nombres descriptivos: ChatService
//
// 6. ARQUITECTURA HEXAGONAL - PUERTOS:
//    - Puertos Primarios (Driving): Definen qué puede hacer la aplicación
//      → Ejemplo: ChatService (casos de uso que los handlers invocan)
//    
//    - Puertos Secundarios (Driven): Definen qué necesita la aplicación
//      → Ejemplo: GroqRepository (cómo acceder a recursos externos)
//
//    El DOMINIO define las interfaces
//    La INFRAESTRUCTURA las implementa
//    La APLICACIÓN las usa
//
// 7. DEPENDENCY INJECTION: En Go se hace manualmente, sin frameworks
//    - Pasas las dependencias como parámetros al constructor
//    - Las estructuras guardan las dependencias como campos
//
//    Ejemplo:
//    type Service struct {
//        repo GroqRepository // Dependencia (interfaz)
//    }
//    
//    func NewService(repo GroqRepository) *Service {
//        return &Service{repo: repo}
//    }
//
// ============================================================================

// ============================================================================
// EJEMPLO DE IMPLEMENTACIÓN (para entender cómo se usa)
// ============================================================================
//
// Así se vería una implementación de estas interfaces:
//
// // Implementación del servicio (capa de aplicación)
// type chatServiceImpl struct {
//     groqRepo GroqRepository  // Dependencia inyectada
// }
//
// func (s *chatServiceImpl) SendMessage(ctx context.Context, message string, model string) (*ChatResponse, error) {
//     request := NewChatRequest(model, []ChatMessage{
//         NewChatMessage("user", message),
//     })
//     return s.groqRepo.CreateChatCompletion(ctx, request)
// }
//
// // Implementación del repositorio (capa de infraestructura)
// type groqHTTPClient struct {
//     httpClient *http.Client
//     apiKey     string
// }
//
// func (c *groqHTTPClient) CreateChatCompletion(ctx context.Context, request ChatRequest) (*ChatResponse, error) {
//     // Aquí iría el código HTTP para llamar a la API de Groq
//     // ...
// }
//
// ============================================================================
