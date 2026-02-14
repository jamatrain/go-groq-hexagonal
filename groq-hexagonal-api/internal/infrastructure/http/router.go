// Package http - Router y configuración de rutas
package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// ============================================================================
// ROUTER SETUP
// ============================================================================

// SetupRouter configura y retorna el router HTTP con todas las rutas
//
// Parámetros:
//   - handler: el ChatHandler con todos los handlers
//
// Retorna:
//   - http.Handler: router configurado y listo para usar
func SetupRouter(handler *ChatHandler) http.Handler {
	// ========================================================================
	// 1. CREAR EL ROUTER
	// ========================================================================

	// gorilla/mux es un router HTTP potente y flexible
	// Es más poderoso que el http.ServeMux estándar
	router := mux.NewRouter()

	// ========================================================================
	// 2. CONFIGURAR MIDDLEWARES GLOBALES
	// ========================================================================

	// Middleware de logging para todas las rutas
	router.Use(loggingMiddleware)

	// Middleware de recovery para capturar panics
	router.Use(recoveryMiddleware)

	// ========================================================================
	// 3. DEFINIR RUTAS
	// ========================================================================

	// API v1 subrouter
	// Esto crea un "sub-router" que maneja todas las rutas bajo /api/v1
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// POST /api/v1/chat - Enviar mensaje al modelo
	apiV1.HandleFunc("/chat", handler.HandleChat).Methods(http.MethodPost)

	// GET /api/v1/models - Obtener modelos disponibles
	apiV1.HandleFunc("/models", handler.HandleGetModels).Methods(http.MethodGet)

	// Health check endpoint (fuera de /api/v1)
	// GET /health - Verificar estado del servicio
	router.HandleFunc("/health", handler.HandleHealth).Methods(http.MethodGet)

	// Ruta raíz (opcional)
	router.HandleFunc("/", handleRoot).Methods(http.MethodGet)

	// ========================================================================
	// 4. CONFIGURAR CORS
	// ========================================================================

	// CORS permite que frontends en otros dominios accedan a la API
	corsHandler := cors.New(cors.Options{
		// AllowedOrigins: dominios permitidos
		// ["*"] permite todos (OK para desarrollo, restringir en producción)
		AllowedOrigins: []string{"*"},

		// AllowedMethods: métodos HTTP permitidos
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},

		// AllowedHeaders: headers permitidos en las peticiones
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
		},

		// ExposedHeaders: headers que el cliente puede leer
		ExposedHeaders: []string{
			"Content-Length",
		},

		// AllowCredentials: permitir cookies
		AllowCredentials: true,

		// MaxAge: tiempo que el browser cachea la respuesta preflight
		MaxAge: 300, // 5 minutos
	})

	// Envolver el router con el handler de CORS
	return corsHandler.Handler(router)
}

// ============================================================================
// HANDLERS AUXILIARES
// ============================================================================

// handleRoot maneja GET /
// Retorna información básica sobre la API
func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Escribir un JSON simple
	// En producción, podrías usar un struct y json.NewEncoder()
	response := `{
		"name": "Groq Hexagonal API",
		"version": "1.0.0",
		"description": "API REST para interactuar con Groq usando Arquitectura Hexagonal",
		"endpoints": {
			"chat": "POST /api/v1/chat",
			"models": "GET /api/v1/models",
			"health": "GET /health"
		},
		"documentation": "https://github.com/tu-usuario/groq-hexagonal-api"
	}`

	w.Write([]byte(response))
}

// ============================================================================
// MIDDLEWARES
// ============================================================================
//
// Los middlewares son funciones que se ejecutan antes/después de los handlers
// Permiten añadir funcionalidad cross-cutting (logging, auth, etc.)
//
// Firma del middleware:
// func(http.Handler) http.Handler
//
// El middleware recibe un handler y retorna otro handler que:
// 1. Hace algo antes (ej: logging)
// 2. Llama al handler original
// 3. Hace algo después (ej: medir tiempo)
// ============================================================================

// loggingMiddleware registra todas las peticiones HTTP
func loggingMiddleware(next http.Handler) http.Handler {
	// http.HandlerFunc convierte una función en un Handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Registrar el inicio de la petición
		start := time.Now()

		log.Printf(
			"[%s] %s %s - Iniciando",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
		)

		// Llamar al siguiente handler en la cadena
		next.ServeHTTP(w, r)

		// Registrar el final de la petición con duración
		duration := time.Since(start)
		log.Printf(
			"[%s] %s %s - Completado en %v",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			duration,
		)
	})
}

// recoveryMiddleware captura panics y previene que crashee el servidor
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// defer con recover() captura panics
		defer func() {
			// recover() retorna nil si no hay panic, o el valor del panic
			if err := recover(); err != nil {
				// Registrar el panic
				log.Printf("PANIC: %v", err)

				// Retornar error 500 al cliente
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"success": false, "error": "internal server error"}`))
			}
		}()

		// Llamar al siguiente handler
		next.ServeHTTP(w, r)
	})
}

// ============================================================================
// CONCEPTOS CLAVE DE GO EXPLICADOS:
// ============================================================================
//
// 1. GORILLA MUX:
//    - Router HTTP más potente que el estándar
//    - Soporta variables en rutas: /users/{id}
//    - Permite restringir por método HTTP
//    - Soporta subrouters para organizar rutas
//
//    Ejemplo:
//    router := mux.NewRouter()
//    router.HandleFunc("/users/{id}", getUser).Methods("GET")
//
// 2. SUB-ROUTERS:
//    - Agrupan rutas bajo un prefijo común
//    - Útil para versionado: /api/v1, /api/v2
//
//    apiV1 := router.PathPrefix("/api/v1").Subrouter()
//    apiV1.HandleFunc("/users", getUsers)  // -> /api/v1/users
//
// 3. MIDDLEWARES:
//    - Funciones que envuelven handlers
//    - Permiten ejecutar código antes/después del handler
//    - Se "encadenan" formando un pipeline
//
//    Orden: Middleware1 -> Middleware2 -> Handler
//
// 4. CORS (Cross-Origin Resource Sharing):
//    - Política de seguridad del navegador
//    - Controla qué dominios pueden acceder a la API
//    - Necesario para APIs consumidas desde el browser
//
// 5. PANIC Y RECOVER:
//    - panic() es como throw en otros lenguajes
//    - recover() captura el panic (como catch)
//    - Solo funciona dentro de defer
//
//    defer func() {
//        if err := recover(); err != nil {
//            // Manejar el panic
//        }
//    }()
//
// 6. http.Handler vs http.HandlerFunc:
//    - http.Handler: interfaz con método ServeHTTP()
//    - http.HandlerFunc: tipo función que implementa Handler
//    - HandlerFunc() convierte funciones en Handlers
//
// 7. MÉTODOS HTTP:
//    - GET: obtener recursos
//    - POST: crear recursos
//    - PUT: actualizar recursos (completo)
//    - PATCH: actualizar recursos (parcial)
//    - DELETE: eliminar recursos
//    - OPTIONS: petición preflight de CORS
//
// 8. VERSIONADO DE APIs:
//    - Usar prefijos: /api/v1, /api/v2
//    - Permite evolucionar la API sin romper clientes
//    - Mantener v1 mientras migras a v2
//
// ============================================================================

// ============================================================================
// EJEMPLO DE MIDDLEWARE PERSONALIZADO:
// ============================================================================
//
// // authMiddleware verifica autenticación
// func authMiddleware(next http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         // Obtener token del header
//         token := r.Header.Get("Authorization")
//
//         // Validar token
//         if !isValidToken(token) {
//             http.Error(w, "no autorizado", http.StatusUnauthorized)
//             return // No llamar a next
//         }
//
//         // Token válido, continuar
//         next.ServeHTTP(w, r)
//     })
// }
//
// // Uso:
// router.Use(authMiddleware)  // Aplica a todas las rutas
// // O solo a algunas:
// apiV1.HandleFunc("/admin", handler).Methods("GET")
// apiV1.Use(authMiddleware)  // Solo /api/v1/*
//
// ============================================================================

// ============================================================================
// ORDEN DE EJECUCIÓN:
// ============================================================================
//
// Para una petición POST /api/v1/chat:
//
// 1. CORS Handler (preflight check)
// 2. loggingMiddleware (log inicio)
// 3. recoveryMiddleware (preparar recover)
// 4. handler.HandleChat (procesar petición)
// 5. recoveryMiddleware (verificar panic)
// 6. loggingMiddleware (log fin + duración)
// 7. CORS Handler (añadir headers CORS)
//
// ============================================================================

// ============================================================================
// MEJORES PRÁCTICAS:
// ============================================================================
//
// 1. USAR SUB-ROUTERS:
//    - Organiza rutas por recurso o versión
//    - Más fácil de mantener
//
// 2. APLICAR MIDDLEWARES ESTRATÉGICAMENTE:
//    - Globales: logging, recovery, CORS
//    - Específicos: autenticación solo donde se necesita
//
// 3. RESTRINGIR MÉTODOS HTTP:
//    - .Methods("GET") evita otros métodos
//    - Retorna 405 automáticamente
//
// 4. VERSIONADO:
//    - Siempre versiona desde el inicio
//    - Facilita evolución sin breaking changes
//
// 5. HEALTH CHECKS:
//    - Siempre incluye /health
//    - Útil para load balancers y monitoring
//
// 6. DOCUMENTACIÓN EN ROOT:
//    - GET / retorna info de la API
//    - Lista de endpoints disponibles
//
// ============================================================================
