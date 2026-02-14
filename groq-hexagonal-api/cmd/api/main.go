// Package main es el punto de entrada de la aplicación
// Aquí se ensamblan todas las piezas de la arquitectura hexagonal
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"groq-hexagonal-api/internal/application"
	"groq-hexagonal-api/internal/config"
	"groq-hexagonal-api/internal/infrastructure/groq"
	httpInfra "groq-hexagonal-api/internal/infrastructure/http"
)

// ============================================================================
// MAIN FUNCTION
// ============================================================================

// main es la función de entrada de cualquier programa Go
// Se ejecuta automáticamente cuando inicias la aplicación
func main() {
	// ========================================================================
	// 1. BANNER DE INICIO
	// ========================================================================
	printBanner()
	
	// ========================================================================
	// 2. CARGAR CONFIGURACIÓN
	// ========================================================================
	
	fmt.Println("🔧 Cargando configuración...")
	cfg, err := config.Load()
	if err != nil {
		// log.Fatalf() imprime el error y termina el programa con exit code 1
		log.Fatalf("❌ Error al cargar configuración: %v", err)
	}
	
	// Imprimir configuración (sin info sensible)
	cfg.Print()
	
	// ========================================================================
	// 3. INICIALIZAR DEPENDENCIAS (Dependency Injection)
	// ========================================================================
	//
	// Aquí ensamblamos la arquitectura hexagonal:
	// 1. Infraestructura (adaptadores externos)
	// 2. Aplicación (casos de uso)
	// 3. HTTP (adaptadores de entrada)
	//
	// El orden es importante: primero lo más externo (infra),
	// luego lo que depende de ello (aplicación), y finalmente
	// lo que expone la funcionalidad (HTTP)
	// ========================================================================
	
	fmt.Println("🔌 Inicializando dependencias...")
	
	// CAPA DE INFRAESTRUCTURA - Adaptador Groq (puerto secundario)
	// Este es el adaptador que se comunica con la API externa de Groq
	groqClient := groq.NewGroqClient(
		cfg.GroqAPIKey,
		cfg.GroqBaseURL,
		cfg.HTTPTimeout,
	)
	fmt.Println("   ✓ Cliente Groq inicializado")
	
	// CAPA DE APLICACIÓN - Servicio de Chat (lógica de negocio)
	// Inyectamos el groqClient al servicio
	// El servicio solo conoce la interfaz, no la implementación
	chatService := application.NewChatService(groqClient, cfg.DefaultModel)
	fmt.Println("   ✓ Servicio de chat inicializado")
	
	// CAPA DE INFRAESTRUCTURA - Handler HTTP (puerto primario)
	// Inyectamos el chatService al handler
	chatHandler := httpInfra.NewChatHandler(chatService)
	fmt.Println("   ✓ Handlers HTTP inicializados")
	
	// CAPA DE INFRAESTRUCTURA - Router HTTP
	// Configuramos todas las rutas
	router := httpInfra.SetupRouter(chatHandler)
	fmt.Println("   ✓ Router configurado")
	
	// ========================================================================
	// 4. CONFIGURAR SERVIDOR HTTP
	// ========================================================================
	
	// http.Server permite configurar timeouts y otras opciones
	// Esto es mejor que usar http.ListenAndServe() directamente
	server := &http.Server{
		Addr:    cfg.GetServerAddress(), // ej: ":8080"
		Handler: router,                 // El router configurado
		
		// Timeouts importantes para seguridad y performance
		ReadTimeout:  15 * time.Second, // Tiempo máx para leer el request
		WriteTimeout: 15 * time.Second, // Tiempo máx para escribir la response
		IdleTimeout:  60 * time.Second, // Tiempo máx que una conexión keep-alive puede estar idle
	}
	
	// ========================================================================
	// 5. INICIAR SERVIDOR EN GOROUTINE
	// ========================================================================
	//
	// Usamos una goroutine para que el servidor no bloquee
	// Esto nos permite manejar señales de shutdown más adelante
	//
	go func() {
		fmt.Println()
		fmt.Printf("🚀 Servidor escuchando en http://localhost%s\n", cfg.GetServerAddress())
		fmt.Println("📡 Endpoints disponibles:")
		fmt.Printf("   • POST http://localhost%s/api/v1/chat\n", cfg.GetServerAddress())
		fmt.Printf("   • GET  http://localhost%s/api/v1/models\n", cfg.GetServerAddress())
		fmt.Printf("   • GET  http://localhost%s/health\n", cfg.GetServerAddress())
		fmt.Println()
		fmt.Println("👉 Presiona Ctrl+C para detener el servidor")
		fmt.Println()
		
		// ListenAndServe() bloquea hasta que el servidor se detenga
		// Retorna error si falla al iniciar (ej: puerto ocupado)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Error al iniciar servidor: %v", err)
		}
	}()
	
	// ========================================================================
	// 6. GRACEFUL SHUTDOWN
	// ========================================================================
	//
	// Manejar señales del sistema para shutdown gracioso
	// Esto permite que las peticiones en curso terminen antes de cerrar
	//
	waitForShutdown(server)
}

// ============================================================================
// FUNCIONES AUXILIARES
// ============================================================================

// waitForShutdown espera una señal de interrupción y hace shutdown gracioso
func waitForShutdown(server *http.Server) {
	// Crear un canal para recibir señales del sistema
	// make(chan os.Signal, 1) crea un canal con buffer de 1
	quit := make(chan os.Signal, 1)
	
	// signal.Notify() envía señales al canal
	// SIGINT es Ctrl+C
	// SIGTERM es kill
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	// Bloquear hasta recibir una señal
	// <-quit lee del canal (bloquea hasta que llegue algo)
	sig := <-quit
	fmt.Printf("\n🛑 Señal recibida: %v\n", sig)
	fmt.Println("🔄 Apagando servidor graciosamente...")
	
	// Crear un contexto con timeout para el shutdown
	// 30 segundos para que las peticiones en curso terminen
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel() asegura que se liberen recursos
	defer cancel()
	
	// server.Shutdown() intenta cerrar el servidor graciosamente
	// Espera a que las conexiones activas terminen (hasta el timeout)
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("❌ Error durante shutdown: %v", err)
	}
	
	fmt.Println("✅ Servidor detenido correctamente")
	fmt.Println("👋 ¡Hasta luego!")
}

// printBanner imprime el banner de inicio de la aplicación
func printBanner() {
	banner := `
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║   ██████╗ ██████╗  ██████╗  ██████╗     █████╗ ██████╗ ██╗║
║  ██╔════╝ ██╔══██╗██╔═══██╗██╔═══██╗   ██╔══██╗██╔══██╗██║║
║  ██║  ███╗██████╔╝██║   ██║██║   ██║   ███████║██████╔╝██║║
║  ██║   ██║██╔══██╗██║   ██║██║▄▄ ██║   ██╔══██║██╔═══╝ ██║║
║  ╚██████╔╝██║  ██║╚██████╔╝╚██████╔╝   ██║  ██║██║     ██║║
║   ╚═════╝ ╚═╝  ╚═╝ ╚═════╝  ╚══════╝   ╚═╝  ╚═╝╚═╝     ╚═╝║
║                                                           ║
║       API REST con Arquitectura Hexagonal en Go          ║
║                    Powered by Groq                        ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}

// ============================================================================
// CONCEPTOS CLAVE DE GO EXPLICADOS:
// ============================================================================
//
// 1. PACKAGE MAIN:
//    - package main es especial: define un ejecutable
//    - func main() es el punto de entrada
//    - Solo puede haber un package main por programa
//
// 2. IMPORTS:
//    - import "fmt": librería estándar
//    - import "groq-hexagonal-api/internal/config": import interno
//    - Alias: httpInfra "groq.../http" evita conflictos con "net/http"
//
// 3. GOROUTINES:
//    - go func() { ... }() ejecuta función concurrentemente
//    - Es como un thread ligero
//    - No bloquea el código siguiente
//
// 4. CHANNELS:
//    - make(chan tipo, buffer) crea un canal
//    - ch <- value: enviar al canal
//    - value := <-ch: recibir del canal
//    - Usados para comunicación entre goroutines
//
// 5. SIGNALS:
//    - os.Signal: señales del sistema operativo
//    - SIGINT: Ctrl+C
//    - SIGTERM: kill
//    - signal.Notify() escucha señales
//
// 6. CONTEXT:
//    - context.WithTimeout(): crea contexto con timeout
//    - Usado para cancelaciones y timeouts
//    - defer cancel(): siempre cancelar para liberar recursos
//
// 7. GRACEFUL SHUTDOWN:
//    - server.Shutdown() cierra graciosamente
//    - Espera a que las conexiones terminen
//    - Importante para no perder requests
//
// 8. LOG vs FMT:
//    - log.Printf(): incluye timestamp
//    - fmt.Printf(): solo el mensaje
//    - log.Fatalf(): imprime y termina programa (exit 1)
//
// 9. DEPENDENCY INJECTION:
//    - Manual en Go (sin frameworks)
//    - Inyectar dependencias en constructores
//    - Principio: depender de interfaces, no implementaciones
//
// 10. ERROR HANDLING:
//     - Siempre verificar errores
//     - log.Fatalf() para errores fatales
//     - log.Printf() para errores no fatales
//
// ============================================================================

// ============================================================================
// ARQUITECTURA HEXAGONAL - ENSAMBLAJE:
// ============================================================================
//
//              ┌─────────────────────────┐
//              │    ENTRADA (HTTP)       │
//              │  • Router               │
//              │  • Handlers             │
//              └───────────┬─────────────┘
//                          │
//                          ↓
//              ┌─────────────────────────┐
//              │   APLICACIÓN (Casos de  │
//              │        uso)             │
//              │  • ChatService          │
//              └───────────┬─────────────┘
//                          │
//                          ↓
//              ┌─────────────────────────┐
//              │    DOMINIO (Core)       │
//              │  • Entidades            │
//              │  • Interfaces (Ports)   │
//              └───────────┬─────────────┘
//                          │
//                          ↓
//              ┌─────────────────────────┐
//              │  INFRAESTRUCTURA        │
//              │   (Adaptadores)         │
//              │  • GroqClient           │
//              └─────────────────────────┘
//
// Flujo de una petición:
// 1. HTTP Request → Router
// 2. Router → Handler
// 3. Handler → ChatService (aplicación)
// 4. ChatService → GroqRepository (interfaz del dominio)
// 5. GroqClient → API de Groq (implementación de infraestructura)
// 6. Respuesta en sentido inverso
//
// ============================================================================

// ============================================================================
// EJEMPLO DE FLUJO COMPLETO:
// ============================================================================
//
// 1. Usuario: POST /api/v1/chat {"message": "Hola"}
// 2. Router: detecta ruta, llama a handler.HandleChat()
// 3. Handler: valida JSON, llama a chatService.SendMessage()
// 4. Service: crea ChatRequest, llama a groqRepo.CreateChatCompletion()
// 5. GroqClient: hace HTTP POST a api.groq.com
// 6. Groq API: procesa y retorna respuesta
// 7. GroqClient: parsea JSON, retorna ChatResponse
// 8. Service: valida respuesta, retorna al handler
// 9. Handler: mapea a DTO, serializa a JSON
// 10. Router: envía respuesta HTTP al usuario
//
// ============================================================================

// ============================================================================
// MEJORES PRÁCTICAS:
// ============================================================================
//
// 1. GRACEFUL SHUTDOWN:
//    - Siempre implementar shutdown gracioso
//    - Evita pérdida de requests
//    - Importante en producción
//
// 2. TIMEOUTS:
//    - Configurar ReadTimeout, WriteTimeout, IdleTimeout
//    - Previene ataques slowloris
//    - Libera recursos
//
// 3. LOGGING:
//    - Log al inicio de la aplicación
//    - Log configuración (sin info sensible)
//    - Log errores con contexto
//
// 4. ERROR HANDLING:
//    - Validar configuración al inicio
//    - Fallar rápido si algo está mal
//    - Mensajes de error claros
//
// 5. DEPENDENCY INJECTION:
//    - Inyectar todas las dependencias
//    - No crear dependencias dentro de funciones
//    - Facilita testing
//
// 6. SEPARACIÓN DE RESPONSABILIDADES:
//    - main.go solo ensambla
//    - Lógica en otros packages
//    - Mantiene main.go simple y claro
//
// ============================================================================
