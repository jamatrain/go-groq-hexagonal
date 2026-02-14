// Package config maneja la configuración de la aplicación
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// ============================================================================
// CONFIGURATION STRUCT
// ============================================================================

// Config contiene toda la configuración de la aplicación
// Centraliza todos los valores configurables en un solo lugar
type Config struct {
	// Server configuración
	Port string
	
	// Groq API configuración
	GroqAPIKey   string
	GroqBaseURL  string
	DefaultModel string
	HTTPTimeout  time.Duration
}

// ============================================================================
// CONSTRUCTOR
// ============================================================================

// Load carga la configuración desde variables de entorno
// Intenta cargar .env primero, pero no falla si no existe
//
// Retorna:
//   - *Config: configuración cargada
//   - error: error si falta alguna variable requerida
func Load() (*Config, error) {
	// ========================================================================
	// 1. CARGAR .env (si existe)
	// ========================================================================
	
	// godotenv.Load() carga variables desde .env
	// Si el archivo no existe, no es un error crítico
	// Las variables ya podrían estar en el entorno del sistema
	if err := godotenv.Load(); err != nil {
		// No es fatal, solo advertir
		fmt.Println("⚠️  Advertencia: archivo .env no encontrado, usando variables de entorno del sistema")
	}
	
	// ========================================================================
	// 2. LEER VARIABLES DE ENTORNO
	// ========================================================================
	
	config := &Config{
		Port:         getEnv("PORT", "8080"),              // Default: 8080
		GroqAPIKey:   getEnv("GROQ_API_KEY", ""),          // Sin default (requerido)
		GroqBaseURL:  getEnv("GROQ_BASE_URL", "https://api.groq.com/openai/v1"),
		DefaultModel: getEnv("DEFAULT_MODEL", "llama-3.3-70b-versatile"),
		HTTPTimeout:  getEnvAsDuration("HTTP_TIMEOUT", 30*time.Second),
	}
	
	// ========================================================================
	// 3. VALIDAR CONFIGURACIÓN
	// ========================================================================
	
	if err := config.Validate(); err != nil {
		return nil, err
	}
	
	return config, nil
}

// ============================================================================
// MÉTODOS DE VALIDACIÓN
// ============================================================================

// Validate verifica que la configuración sea válida
func (c *Config) Validate() error {
	// Verificar que el API key no esté vacío
	if c.GroqAPIKey == "" {
		return fmt.Errorf("GROQ_API_KEY es requerido")
	}
	
	// Verificar que el base URL no esté vacío
	if c.GroqBaseURL == "" {
		return fmt.Errorf("GROQ_BASE_URL es requerido")
	}
	
	// Verificar que el puerto sea válido
	if c.Port == "" {
		return fmt.Errorf("PORT es requerido")
	}
	
	// Verificar que el timeout sea positivo
	if c.HTTPTimeout <= 0 {
		return fmt.Errorf("HTTP_TIMEOUT debe ser mayor a 0")
	}
	
	return nil
}

// ============================================================================
// MÉTODOS DE UTILIDAD
// ============================================================================

// GetServerAddress retorna la dirección completa del servidor
// Útil para iniciar el servidor HTTP
func (c *Config) GetServerAddress() string {
	// Format: :8080 (escucha en todas las interfaces)
	return ":" + c.Port
}

// Print imprime la configuración (sin información sensible)
// Útil para debugging y logs de inicio
func (c *Config) Print() {
	fmt.Println("📋 Configuración cargada:")
	fmt.Printf("   • Puerto: %s\n", c.Port)
	fmt.Printf("   • Groq Base URL: %s\n", c.GroqBaseURL)
	fmt.Printf("   • Modelo por defecto: %s\n", c.DefaultModel)
	fmt.Printf("   • HTTP Timeout: %v\n", c.HTTPTimeout)
	// NO imprimir el API key por seguridad
	fmt.Printf("   • API Key: %s\n", maskAPIKey(c.GroqAPIKey))
}

// ============================================================================
// FUNCIONES AUXILIARES (helpers)
// ============================================================================

// getEnv obtiene una variable de entorno o retorna un valor por defecto
//
// Parámetros:
//   - key: nombre de la variable
//   - defaultValue: valor por defecto si no existe
//
// Retorna:
//   - string: valor de la variable o default
func getEnv(key, defaultValue string) string {
	// os.Getenv() obtiene una variable de entorno
	// Retorna "" si no existe
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt obtiene una variable de entorno como int
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	
	// strconv.Atoi() convierte string a int
	// Retorna error si no es un número válido
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	
	return value
}

// getEnvAsDuration obtiene una variable de entorno como time.Duration
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	
	// Intentar parsear como número de segundos
	seconds, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	
	// Convertir segundos a Duration
	return time.Duration(seconds) * time.Second
}

// maskAPIKey oculta parcialmente el API key para logs
// Muestra solo los primeros y últimos caracteres
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		// Si es muy corta, ocultar todo
		return "***"
	}
	
	// Mostrar primeros 4 y últimos 4 caracteres
	return key[:4] + "..." + key[len(key)-4:]
}

// ============================================================================
// CONCEPTOS CLAVE DE GO EXPLICADOS:
// ============================================================================
//
// 1. VARIABLES DE ENTORNO:
//    - os.Getenv("KEY"): obtiene una variable de entorno
//    - Forma estándar de configurar aplicaciones
//    - Sigue los principios de 12-factor app
//
// 2. .ENV FILES:
//    - godotenv.Load() carga variables desde .env
//    - Útil para desarrollo local
//    - En producción, usar variables del sistema
//
// 3. VALORES POR DEFECTO:
//    - Siempre proveer defaults razonables
//    - Excepto para valores críticos (API keys)
//    - Hace la aplicación más fácil de usar
//
// 4. VALIDACIÓN:
//    - Validar temprano (al inicio de la aplicación)
//    - Fallar rápido si falta configuración
//    - Mejor que fallar en runtime
//
// 5. STRCONV:
//    - strconv.Atoi(): string → int
//    - strconv.Itoa(): int → string
//    - strconv.ParseFloat(), ParseBool(), etc.
//
// 6. TIME.DURATION:
//    - Tipo para representar duraciones
//    - time.Second, time.Minute, etc.
//    - Operaciones: 30 * time.Second
//
// 7. FMT.PRINTF:
//    - fmt.Printf() imprime con formato
//    - %s: string
//    - %d: int
//    - %v: cualquier valor
//    - %+v: struct con nombres de campos
//
// 8. SEGURIDAD:
//    - NUNCA loguear API keys completas
//    - Usar maskAPIKey() o similar
//    - Importante para seguridad
//
// ============================================================================

// ============================================================================
// EJEMPLO DE .env FILE:
// ============================================================================
//
// # .env
// PORT=8080
// GROQ_API_KEY=gsk_abc123...
// GROQ_BASE_URL=https://api.groq.com/openai/v1
// DEFAULT_MODEL=llama-3.3-70b-versatile
// HTTP_TIMEOUT=30
//
// ============================================================================

// ============================================================================
// EJEMPLO DE USO:
// ============================================================================
//
// func main() {
//     // Cargar configuración
//     cfg, err := config.Load()
//     if err != nil {
//         log.Fatalf("Error al cargar configuración: %v", err)
//     }
//     
//     // Imprimir configuración
//     cfg.Print()
//     
//     // Usar la configuración
//     groqClient := groq.NewGroqClient(
//         cfg.GroqAPIKey,
//         cfg.GroqBaseURL,
//         cfg.HTTPTimeout,
//     )
//     
//     // Iniciar servidor
//     log.Printf("Servidor escuchando en %s", cfg.GetServerAddress())
//     http.ListenAndServe(cfg.GetServerAddress(), router)
// }
//
// ============================================================================

// ============================================================================
// MEJORES PRÁCTICAS:
// ============================================================================
//
// 1. CENTRALIZAR CONFIGURACIÓN:
//    - Un solo lugar para toda la config
//    - Fácil de encontrar y modificar
//
// 2. VALIDACIÓN TEMPRANA:
//    - Validar al inicio de la app
//    - Fallar rápido y claro
//
// 3. VALORES POR DEFECTO SENSATOS:
//    - Defaults que funcionen "out of the box"
//    - Excepto valores sensibles
//
// 4. SEGURIDAD:
//    - NO commitear .env al repositorio
//    - Usar .env.example con valores de ejemplo
//    - Ocultar API keys en logs
//
// 5. DOCUMENTACIÓN:
//    - Documentar cada variable
//    - Explicar valores válidos
//    - Proveer ejemplos
//
// 6. 12-FACTOR APP:
//    - Config en variables de entorno
//    - Separar config de código
//    - Una config por entorno (dev, staging, prod)
//
// ============================================================================
