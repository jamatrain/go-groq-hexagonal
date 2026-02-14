// go.mod define el módulo de Go y sus dependencias
// El nombre del módulo es usado para imports internos
module groq-hexagonal-api

// Versión de Go requerida
go 1.22

// Dependencias externas necesarias
require (
	github.com/gorilla/mux v1.8.1 // Router HTTP potente y flexible
	github.com/joho/godotenv v1.5.1 // Cargar variables de entorno desde .env
	github.com/rs/cors v1.10.1 // Manejo de CORS para APIs
)
