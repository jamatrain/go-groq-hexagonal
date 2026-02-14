// Package groq implementa el adaptador para comunicarse con la API de Groq
// Esta es la CAPA DE INFRAESTRUCTURA - contiene detalles de implementación
package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"groq-hexagonal-api/internal/domain"
	"io"
	"net/http"
	"time"
)

// ============================================================================
// CONSTANTES
// ============================================================================

const (
	// Endpoints de la API de Groq
	ChatCompletionsEndpoint = "/chat/completions"
	ModelsEndpoint          = "/models"
	
	// Headers HTTP
	ContentTypeJSON   = "application/json"
	AuthorizationHeader = "Authorization"
)

// ============================================================================
// CLIENT STRUCT
// ============================================================================

// GroqClient es el adaptador HTTP que implementa domain.GroqRepository
// Implementa la interfaz implícitamente (no necesita declararlo)
type GroqClient struct {
	// httpClient es el cliente HTTP estándar de Go
	// Lo reutilizamos para todas las peticiones (connection pooling)
	httpClient *http.Client
	
	// baseURL es la URL base de la API (ej: https://api.groq.com/openai/v1)
	baseURL string
	
	// apiKey es la clave de autenticación
	apiKey string
}

// ============================================================================
// CONSTRUCTOR
// ============================================================================

// NewGroqClient crea un nuevo cliente para la API de Groq
//
// Parámetros:
//   - apiKey: tu API key de Groq
//   - baseURL: URL base de la API
//   - timeout: tiempo máximo de espera para requests
//
// Retorna:
//   - domain.GroqRepository: retornamos la interfaz (buena práctica)
func NewGroqClient(apiKey, baseURL string, timeout time.Duration) domain.GroqRepository {
	// Validación básica
	if apiKey == "" {
		panic("apiKey no puede estar vacía")
	}
	if baseURL == "" {
		panic("baseURL no puede estar vacía")
	}
	
	// Crear el cliente HTTP con timeout
	// &http.Client{...} crea un puntero a http.Client
	httpClient := &http.Client{
		Timeout: timeout, // Timeout total para cada request
		
		// Transport controla cómo se hacen las conexiones HTTP
		Transport: &http.Transport{
			// Configuración de connection pooling
			MaxIdleConns:        100,              // Máx. conexiones idle totales
			MaxIdleConnsPerHost: 10,               // Máx. conexiones idle por host
			IdleConnTimeout:     90 * time.Second, // Tiempo antes de cerrar conexión idle
		},
	}
	
	return &GroqClient{
		httpClient: httpClient,
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

// ============================================================================
// IMPLEMENTACIÓN DE domain.GroqRepository
// ============================================================================

// CreateChatCompletion implementa la interfaz GroqRepository
// Envía una petición POST a /chat/completions
func (c *GroqClient) CreateChatCompletion(
	ctx context.Context,
	request domain.ChatRequest,
) (*domain.ChatResponse, error) {
	// Construir la URL completa
	// c.baseURL + ChatCompletionsEndpoint
	url := c.baseURL + ChatCompletionsEndpoint
	
	// Serializar el request a JSON
	// json.Marshal() convierte un struct Go a JSON bytes
	jsonData, err := json.Marshal(request)
	if err != nil {
		// fmt.Errorf() crea un error con formato
		// %w preserva el error original para wrapping
		return nil, fmt.Errorf("error al serializar request: %w", err)
	}
	
	// Hacer la petición HTTP POST
	response, err := c.doRequest(ctx, http.MethodPost, url, jsonData)
	if err != nil {
		return nil, fmt.Errorf("error en la petición HTTP: %w", err)
	}
	
	// Parsear la respuesta
	var chatResponse domain.ChatResponse
	if err := json.Unmarshal(response, &chatResponse); err != nil {
		return nil, fmt.Errorf("error al parsear respuesta: %w", err)
	}
	
	// Retornar la respuesta parseada
	// &chatResponse crea un puntero al chatResponse
	return &chatResponse, nil
}

// ListModels implementa la interfaz GroqRepository
// Envía una petición GET a /models
func (c *GroqClient) ListModels(ctx context.Context) (*domain.ModelsResponse, error) {
	// Construir la URL completa
	url := c.baseURL + ModelsEndpoint
	
	// Hacer la petición HTTP GET
	// nil porque GET no lleva body
	response, err := c.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error al obtener modelos: %w", err)
	}
	
	// Parsear la respuesta
	var modelsResponse domain.ModelsResponse
	if err := json.Unmarshal(response, &modelsResponse); err != nil {
		return nil, fmt.Errorf("error al parsear modelos: %w", err)
	}
	
	return &modelsResponse, nil
}

// ============================================================================
// MÉTODOS PRIVADOS (helpers)
// ============================================================================

// doRequest es un método privado que realiza la petición HTTP
// Los métodos privados empiezan con minúscula en Go
//
// Parámetros:
//   - ctx: contexto para cancelaciones
//   - method: método HTTP (GET, POST, etc.)
//   - url: URL completa
//   - body: datos a enviar (nil para GET)
//
// Retorna:
//   - []byte: respuesta del servidor en bytes
//   - error: error si algo falla
func (c *GroqClient) doRequest(
	ctx context.Context,
	method string,
	url string,
	body []byte,
) ([]byte, error) {
	// ========================================================================
	// 1. CREAR LA PETICIÓN HTTP
	// ========================================================================
	
	// bytes.NewBuffer() crea un io.Reader desde []byte
	// io.Reader es una interfaz que http.NewRequest espera
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewBuffer(body)
	}
	
	// Crear la petición HTTP
	// http.NewRequestWithContext incluye el contexto para cancelaciones
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %w", err)
	}
	
	// ========================================================================
	// 2. CONFIGURAR HEADERS
	// ========================================================================
	
	// Establecer Content-Type
	req.Header.Set("Content-Type", ContentTypeJSON)
	
	// Establecer Authorization
	// La API de Groq usa Bearer token
	req.Header.Set(AuthorizationHeader, "Bearer "+c.apiKey)
	
	// ========================================================================
	// 3. EJECUTAR LA PETICIÓN
	// ========================================================================
	
	// c.httpClient.Do() ejecuta la petición HTTP
	// Usa el contexto para timeouts y cancelaciones
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar request: %w", err)
	}
	
	// defer asegura que el body se cierre al final de la función
	// Esto es CRÍTICO para no tener memory leaks
	// defer se ejecuta cuando la función retorna (como finally)
	defer resp.Body.Close()
	
	// ========================================================================
	// 4. LEER LA RESPUESTA
	// ========================================================================
	
	// io.ReadAll() lee todo el body de la respuesta
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %w", err)
	}
	
	// ========================================================================
	// 5. VERIFICAR STATUS CODE
	// ========================================================================
	
	// Verificar si la respuesta es exitosa (2xx)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Si no es 2xx, retornar error con el status y el body
		return nil, fmt.Errorf(
			"API retornó status %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}
	
	// ========================================================================
	// 6. RETORNAR RESPUESTA
	// ========================================================================
	
	return responseBody, nil
}

// ============================================================================
// CONCEPTOS CLAVE DE GO EXPLICADOS:
// ============================================================================
//
// 1. HTTP CLIENT:
//    - http.Client: cliente HTTP reutilizable con connection pooling
//    - http.NewRequestWithContext(): crea requests con contexto
//    - resp.Body.Close(): SIEMPRE cerrar el body (usa defer)
//
// 2. JSON EN GO:
//    - json.Marshal(v): convierte struct → JSON bytes
//    - json.Unmarshal(data, &v): convierte JSON bytes → struct
//    - Las tags `json:"..."` controlan la serialización
//
// 3. BYTES y IO:
//    - []byte: slice de bytes (array de bytes)
//    - io.Reader: interfaz para leer datos
//    - io.ReadAll(): lee todo de un Reader
//    - bytes.NewBuffer(): crea Reader desde []byte
//
// 4. DEFER:
//    - defer posterga la ejecución hasta que la función retorne
//    - Se usa para cleanup (cerrar archivos, conexiones, etc.)
//    - Múltiples defers se ejecutan en orden LIFO (último primero)
//    
//    Ejemplo:
//    func readFile() {
//        f, _ := os.Open("file.txt")
//        defer f.Close() // Se ejecuta al retornar
//        // ... usar f ...
//    }
//
// 5. ERROR WRAPPING:
//    - fmt.Errorf() con %w preserva el error original
//    - Permite usar errors.Is() y errors.As()
//    - Crea una cadena de errores para debugging
//
// 6. CONTEXTO (context.Context):
//    - Propaga cancelaciones y timeouts
//    - ctx.Done(): canal que se cierra cuando se cancela
//    - ctx.Err(): retorna el error de cancelación
//
// 7. INTERFACES IMPLÍCITAS:
//    - GroqClient implementa domain.GroqRepository sin declararlo
//    - Solo necesita tener los métodos correctos
//    - Esto permite desacoplamiento total
//
// 8. CONNECTION POOLING:
//    - http.Client reutiliza conexiones TCP
//    - Mejora el rendimiento significativamente
//    - Configurar MaxIdleConns para optimizar
//
// 9. CONSTANTS:
//    - const define valores inmutables
//    - Se evalúan en tiempo de compilación
//    - Mejor performance que variables
//
// 10. TIME.DURATION:
//     - Tipo para representar duraciones
//     - time.Second, time.Minute, etc.
//     - Se pueden multiplicar: 30 * time.Second
//
// ============================================================================

// ============================================================================
// EJEMPLO DE USO:
// ============================================================================
//
// // Crear el cliente
// client := NewGroqClient(
//     "tu-api-key",
//     "https://api.groq.com/openai/v1",
//     30 * time.Second,
// )
//
// // Crear una petición
// request := domain.NewChatRequest("llama-3.3-70b-versatile", []domain.ChatMessage{
//     domain.NewChatMessage("user", "Hola!"),
// })
//
// // Hacer la petición
// ctx := context.Background()
// response, err := client.CreateChatCompletion(ctx, request)
// if err != nil {
//     log.Fatal(err)
// }
//
// fmt.Println(response.GetResponseContent())
//
// ============================================================================
