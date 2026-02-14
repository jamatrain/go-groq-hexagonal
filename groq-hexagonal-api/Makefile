# Makefile para facilitar el desarrollo
# Uso: make <comando>

.PHONY: help run build test clean install

# Comando por defecto
.DEFAULT_GOAL := help

# Colores para output
GREEN  := \033[0;32m
YELLOW := \033[0;33m
NC     := \033[0m # No Color

## help: Muestra esta ayuda
help:
	@echo "$(GREEN)Comandos disponibles:$(NC)"
	@echo "  $(YELLOW)make run$(NC)      - Ejecutar la aplicación"
	@echo "  $(YELLOW)make build$(NC)    - Compilar la aplicación"
	@echo "  $(YELLOW)make test$(NC)     - Ejecutar tests"
	@echo "  $(YELLOW)make clean$(NC)    - Limpiar archivos compilados"
	@echo "  $(YELLOW)make install$(NC)  - Instalar dependencias"
	@echo "  $(YELLOW)make dev$(NC)      - Modo desarrollo (con hot reload)"

## install: Instala las dependencias del proyecto
install:
	@echo "$(GREEN)Instalando dependencias...$(NC)"
	go mod download
	go mod tidy
	@echo "$(GREEN)✓ Dependencias instaladas$(NC)"

## run: Ejecuta la aplicación
run:
	@echo "$(GREEN)Iniciando aplicación...$(NC)"
	go run cmd/api/main.go

## build: Compila la aplicación
build:
	@echo "$(GREEN)Compilando aplicación...$(NC)"
	go build -o bin/groq-api cmd/api/main.go
	@echo "$(GREEN)✓ Compilado en: bin/groq-api$(NC)"

## test: Ejecuta los tests
test:
	@echo "$(GREEN)Ejecutando tests...$(NC)"
	go test -v ./...

## test-coverage: Ejecuta tests con coverage
test-coverage:
	@echo "$(GREEN)Ejecutando tests con coverage...$(NC)"
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report: coverage.html$(NC)"

## clean: Limpia archivos compilados
clean:
	@echo "$(YELLOW)Limpiando archivos...$(NC)"
	rm -f bin/groq-api
	rm -f coverage.out coverage.html
	@echo "$(GREEN)✓ Limpieza completada$(NC)"

## fmt: Formatea el código
fmt:
	@echo "$(GREEN)Formateando código...$(NC)"
	go fmt ./...
	@echo "$(GREEN)✓ Código formateado$(NC)"

## lint: Ejecuta el linter (requiere golangci-lint instalado)
lint:
	@echo "$(GREEN)Ejecutando linter...$(NC)"
	golangci-lint run ./...

## dev: Modo desarrollo con hot reload (requiere air instalado)
dev:
	@echo "$(GREEN)Iniciando modo desarrollo...$(NC)"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(YELLOW)⚠️  'air' no está instalado. Instálalo con:$(NC)"; \
		echo "   go install github.com/air-verse/air@latest"; \
		echo "$(YELLOW)Ejecutando sin hot reload...$(NC)"; \
		go run cmd/api/main.go; \
	fi

## docker-build: Construye imagen Docker
docker-build:
	@echo "$(GREEN)Construyendo imagen Docker...$(NC)"
	docker build -t groq-hexagonal-api:latest .
	@echo "$(GREEN)✓ Imagen construida$(NC)"

## docker-run: Ejecuta container Docker
docker-run:
	@echo "$(GREEN)Ejecutando container Docker...$(NC)"
	docker run -p 8080:8080 --env-file .env groq-hexagonal-api:latest
