# 🌌 Go Groq Hexagonal

Este proyecto es una aplicación **Full Stack** diseñada para enseñar y demostrar la implementación de **Arquitectura Hexagonal (Ports & Adapters)** en Go, consumiendo la poderosa API de inteligencia artificial de **Groq**, con un frontend moderno en **React**.

## 🏗️ Estructura del Proyecto

El repositorio está organizado como un monorepo con dos componentes principales:

*   **`groq-hexagonal-api/`** (Backend): API RESTful escrita en **Go** siguiendo estrictamente la Arquitectura Hexagonal. Maneja la lógica de negocio, la comunicación con Groq y expone endpoints para el frontend.
*   **`groq-frontend/`** (Frontend): Interfaz de usuario moderna construida con **React** y **Vite**. Permite interactuar con la IA de manera visual y amigable.

## 🚀 Requisitos Previos

Antes de comenzar, asegúrate de tener instalado:

*   [Go](https://go.dev/dl/) (versión 1.22 o superior)
*   [Node.js](https://nodejs.org/) (versión 18 o superior)
*   Una **API Key de Groq** (puedes obtenerla gratis en [console.groq.com](https://console.groq.com))

---

## ⚡ Guía de Inicio Rápido

Sigue estos pasos para levantar todo el entorno de desarrollo.

### 1. Configurar y Ejecutar el Backend (API)

```bash
cd groq-hexagonal-api

# 1. Configurar variables de entorno
cp .env.example .env
# ⚠️ IMPORTANTE: Abre el archivo .env y pega tu GROQ_API_KEY

# 2. Instalar dependencias
make install  # O usa: go mod download

# 3. Ejecutar el servidor
make run      # O usa: go run cmd/api/main.go
```
*El backend estará corriendo en `http://localhost:8080`*

### 2. Configurar y Ejecutar el Frontend

En una **nueva terminal**:

```bash
cd groq-frontend

# 1. Instalar dependencias
npm install API

# 2. Iniciar el servidor de desarrollo
npm run dev
```
*El frontend estará disponible en `http://localhost:3000` (o el puerto que indique Vite)*

---

## 🏛️ Arquitectura

### Backend (Go)
El backend está estructurado para desacoplar el dominio de la infraestructura:
*   **Domain (`internal/domain`)**: Entidades (`Chat`) e interfaces (`Ports`). No tiene dependencias externas.
*   **Application (`internal/application`)**: Casos de uso (`ChatService`). Orquesta la lógica sin saber de HTTP o bases de datos excesivas.
*   **Infrastructure (`internal/infrastructure`)**: Implementaciones concretas (Cliente HTTP de Groq, Handlers HTTP, Router).

### Frontend (React)
Una SPA (Single Page Application) ligera que consume la API del backend. Utiliza `axios` para las peticiones HTTP y mantiene el estado de la conversación localmente.

---

## 📚 Documentación Detallada

Para más detalles sobre cada parte del proyecto, consulta los READMEs específicos:

*   👉 [Documentación del Backend (API)](./groq-hexagonal-api/README.md) - Detalles sobre endpoints, estructura de carpetas y comandos de Makefile.
*   👉 [Documentación del Frontend](./groq-frontend/README.md) - Componentes, estructura de React y configuración de Vite.

## 🤝 Contribuir

¡Las contribuciones son bienvenidas! Si encuentras un bug o tienes una idea para mejorar la arquitectura, no dudes en abrir un issue o un pull request.
