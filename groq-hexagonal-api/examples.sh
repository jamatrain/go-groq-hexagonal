#!/bin/bash

# ============================================================================
# EJEMPLOS DE USO DE LA API - GROQ HEXAGONAL
# ============================================================================
# 
# Este script contiene ejemplos de cómo usar la API con curl
# Puedes ejecutar cada comando individualmente copiando y pegando
#
# ============================================================================

# Colores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# URL base de la API
BASE_URL="http://localhost:8080"

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     EJEMPLOS DE USO - GROQ HEXAGONAL API                   ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# ============================================================================
# 1. HEALTH CHECK
# ============================================================================

echo -e "${GREEN}📊 1. Health Check${NC}"
echo -e "${YELLOW}GET /health${NC}"
echo ""
curl -X GET "${BASE_URL}/health" | jq '.'
echo -e "\n"

# ============================================================================
# 2. INFORMACIÓN DE LA API
# ============================================================================

echo -e "${GREEN}ℹ️  2. Información de la API${NC}"
echo -e "${YELLOW}GET /${NC}"
echo ""
curl -X GET "${BASE_URL}/" | jq '.'
echo -e "\n"

# ============================================================================
# 3. LISTAR MODELOS DISPONIBLES
# ============================================================================

echo -e "${GREEN}🤖 3. Listar modelos disponibles${NC}"
echo -e "${YELLOW}GET /api/v1/models${NC}"
echo ""
curl -X GET "${BASE_URL}/api/v1/models" | jq '.'
echo -e "\n"

# ============================================================================
# 4. CHAT SIMPLE
# ============================================================================

echo -e "${GREEN}💬 4. Chat simple (modelo por defecto)${NC}"
echo -e "${YELLOW}POST /api/v1/chat${NC}"
echo ""
curl -X POST "${BASE_URL}/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Explica qué es la arquitectura hexagonal en 3 líneas"
  }' | jq '.'
echo -e "\n"

# ============================================================================
# 5. CHAT CON MODELO ESPECÍFICO
# ============================================================================

echo -e "${GREEN}💬 5. Chat con modelo específico${NC}"
echo -e "${YELLOW}POST /api/v1/chat${NC}"
echo ""
curl -X POST "${BASE_URL}/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "¿Qué es Go y por qué es popular?",
    "model": "llama-3.3-70b-versatile"
  }' | jq '.'
echo -e "\n"

# ============================================================================
# 6. CHAT CON PARÁMETROS AVANZADOS
# ============================================================================

echo -e "${GREEN}💬 6. Chat con temperatura personalizada${NC}"
echo -e "${YELLOW}POST /api/v1/chat${NC}"
echo ""
curl -X POST "${BASE_URL}/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Escribe un poema corto sobre Go",
    "model": "llama-3.3-70b-versatile",
    "temperature": 0.9,
    "max_tokens": 200
  }' | jq '.'
echo -e "\n"

# ============================================================================
# 7. CHAT - PREGUNTA TÉCNICA
# ============================================================================

echo -e "${GREEN}💬 7. Pregunta técnica sobre Go${NC}"
echo -e "${YELLOW}POST /api/v1/chat${NC}"
echo ""
curl -X POST "${BASE_URL}/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "¿Cuál es la diferencia entre un puntero y un valor en Go?",
    "model": "llama-3.3-70b-versatile"
  }' | jq '.'
echo -e "\n"

# ============================================================================
# 8. MANEJO DE ERRORES - Mensaje vacío
# ============================================================================

echo -e "${GREEN}❌ 8. Manejo de errores - Mensaje vacío${NC}"
echo -e "${YELLOW}POST /api/v1/chat${NC}"
echo ""
curl -X POST "${BASE_URL}/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "message": ""
  }' | jq '.'
echo -e "\n"

# ============================================================================
# 9. MANEJO DE ERRORES - JSON inválido
# ============================================================================

echo -e "${GREEN}❌ 9. Manejo de errores - JSON inválido${NC}"
echo -e "${YELLOW}POST /api/v1/chat${NC}"
echo ""
curl -X POST "${BASE_URL}/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d 'mensaje invalido' | jq '.'
echo -e "\n"

# ============================================================================
# 10. MANEJO DE ERRORES - Método incorrecto
# ============================================================================

echo -e "${GREEN}❌ 10. Manejo de errores - Método HTTP incorrecto${NC}"
echo -e "${YELLOW}GET /api/v1/chat (debería ser POST)${NC}"
echo ""
curl -X GET "${BASE_URL}/api/v1/chat" | jq '.'
echo -e "\n"

# ============================================================================
# FINALIZADO
# ============================================================================

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     EJEMPLOS COMPLETADOS                                   ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}💡 Tip: Usa 'jq' para formatear el JSON:${NC}"
echo -e "   curl ... | jq '.'"
echo ""
echo -e "${GREEN}💡 Tip: Guarda la respuesta en un archivo:${NC}"
echo -e "   curl ... > response.json"
echo ""

# ============================================================================
# CONCEPTOS EXPLICADOS:
# ============================================================================
#
# CURL FLAGS USADOS:
# -X: especifica el método HTTP (GET, POST, etc.)
# -H: añade un header a la petición
# -d: envía datos en el body (para POST)
# | jq '.': formatea el JSON de respuesta
#
# JQ:
# - Herramienta de línea de comandos para procesar JSON
# - Instalar: apt install jq (Linux) o brew install jq (Mac)
# - jq '.' formatea y colorea el JSON
# - jq '.message' extrae solo el campo "message"
#
# CONTENT-TYPE:
# - "application/json" indica que enviamos JSON
# - El servidor debe saber qué tipo de datos recibe
#
# TEMPERATURA:
# - Controla la aleatoriedad del modelo
# - 0.0 = muy determinista, siempre la misma respuesta
# - 1.0 = balance entre creatividad y coherencia
# - 2.0 = muy creativo/aleatorio
#
# MAX_TOKENS:
# - Límite de tokens en la respuesta
# - 1 token ≈ 0.75 palabras en inglés
# - Controla la longitud máxima de la respuesta
#
# ============================================================================
