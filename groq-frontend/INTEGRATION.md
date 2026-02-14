# 🔗 Guía de Integración - Backend + Frontend

Esta guía te ayudará a conectar el backend de Go con el frontend de React.

## 🎯 Arquitectura Completa

```
┌─────────────────────────────────────────────────────────────┐
│                      NAVEGADOR                              │
│  ┌───────────────────────────────────────────────────────┐  │
│  │              React App (Puerto 3000)                  │  │
│  │  • Componentes UI                                     │  │
│  │  • Manejo de estado                                   │  │
│  │  • Fetch API                                          │  │
│  └─────────────────────┬─────────────────────────────────┘  │
└────────────────────────┼────────────────────────────────────┘
                         │
                         │ HTTP Requests
                         │ (fetch/axios)
                         ↓
┌─────────────────────────────────────────────────────────────┐
│              Backend Go API (Puerto 8080)                   │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  Router (Gorilla Mux)                                 │  │
│  │    ↓                                                   │  │
│  │  Handlers HTTP                                        │  │
│  │    ↓                                                   │  │
│  │  Service Layer (Lógica de negocio)                   │  │
│  │    ↓                                                   │  │
│  │  Groq Client                                          │  │
│  └─────────────────────┬─────────────────────────────────┘  │
└────────────────────────┼────────────────────────────────────┘
                         │
                         │ HTTPS
                         ↓
              ┌────────────────────┐
              │     Groq API       │
              │  (api.groq.com)    │
              └────────────────────┘
```

## 📋 Pasos de Integración

### 1️⃣ Configurar el Backend

```bash
# Terminal 1: Backend
cd groq-hexagonal-api

# Configurar .env
cp .env.example .env
# Editar .env y añadir tu GROQ_API_KEY

# Instalar dependencias
go mod download

# Ejecutar
go run cmd/api/main.go
```

**Verificar que funciona:**
```bash
# Test health check
curl http://localhost:8080/health

# Debería retornar:
# {"status":"healthy","timestamp":1234567890,"service":"groq-api"}
```

### 2️⃣ Configurar el Frontend

```bash
# Terminal 2: Frontend
cd groq-frontend

# Instalar dependencias
npm install

# Ejecutar
npm run dev
```

**Verificar que funciona:**
- Abre http://localhost:3000 en tu navegador
- Deberías ver la interfaz del chat

### 3️⃣ Probar la Integración

1. **Abrir el navegador** en http://localhost:3000
2. **Escribir un mensaje** en el input
3. **Presionar Enter** o click en "Enviar"
4. **Esperar la respuesta** del modelo

**Si funciona:**
- ✅ Verás el mensaje del usuario
- ✅ Aparece indicador de "escribiendo..."
- ✅ Llega respuesta del asistente

**Si no funciona:**
- ❌ Ver sección de Troubleshooting

---

## 🔧 Configuración de CORS

El backend ya tiene CORS configurado en `internal/infrastructure/http/router.go`:

```go
corsHandler := cors.New(cors.Options{
    AllowedOrigins: []string{"*"},  // Permite todos los orígenes
    AllowedMethods: []string{
        http.MethodGet,
        http.MethodPost,
        http.MethodPut,
        http.MethodDelete,
        http.MethodOptions,
    },
    AllowedHeaders: []string{
        "Content-Type",
        "Authorization",
        "X-Requested-With",
    },
    AllowCredentials: true,
    MaxAge: 300,
})
```

**Para producción**, cambia `AllowedOrigins`:
```go
AllowedOrigins: []string{"https://tudominio.com"},
```

---

## 📡 Endpoints y Formato de Datos

### POST /api/v1/chat

**Request:**
```json
{
  "message": "Tu mensaje aquí",
  "model": "llama-3.3-70b-versatile"
}
```

**Response (Éxito):**
```json
{
  "success": true,
  "message": "Respuesta del modelo aquí",
  "model": "llama-3.3-70b-versatile",
  "usage": {
    "prompt_tokens": 10,
    "completion_tokens": 50,
    "total_tokens": 60
  }
}
```

**Response (Error):**
```json
{
  "success": false,
  "error": "Descripción del error"
}
```

### GET /api/v1/models

**Response:**
```json
{
  "success": true,
  "models": [
    {
      "id": "llama-3.3-70b-versatile",
      "name": "llama-3.3-70b-versatile",
      "owned_by": "Meta"
    }
  ]
}
```

### GET /health

**Response:**
```json
{
  "status": "healthy",
  "timestamp": 1234567890,
  "service": "groq-api"
}
```

---

## 🐛 Troubleshooting

### Problema: "Failed to fetch"

**Causa:** El backend no está ejecutándose o hay problema de CORS

**Solución:**
1. Verificar que el backend está corriendo:
   ```bash
   curl http://localhost:8080/health
   ```

2. Verificar la consola del navegador (F12 → Console)

3. Verificar la pestaña Network (F12 → Network)

### Problema: "Network Error" o "CORS Error"

**Causa:** Configuración de CORS

**Solución:**
1. Verifica que el backend tenga CORS configurado
2. En desarrollo, verifica que el proxy de Vite esté configurado:
   ```javascript
   // vite.config.js
   server: {
     proxy: {
       '/api': 'http://localhost:8080'
     }
   }
   ```

### Problema: "API returned status 500"

**Causa:** Error en el backend

**Solución:**
1. Ver logs del backend (Terminal 1)
2. Verificar que la API key de Groq es válida
3. Verificar conexión a internet

### Problema: El mensaje no se envía

**Causa:** Validación en el frontend

**Solución:**
1. Verificar que el mensaje no esté vacío
2. Ver consola del navegador para errores
3. Verificar estado con React DevTools

### Problema: "Cannot read property of undefined"

**Causa:** Respuesta de la API en formato incorrecto

**Solución:**
1. Verificar logs del backend
2. Usar console.log para ver la respuesta exacta:
   ```jsx
   console.log('Response:', response)
   ```

---

## 🔍 Debugging

### Backend (Go)

```go
// En handler.go
log.Printf("Request recibido: %+v", req)
log.Printf("Response: %+v", response)
```

### Frontend (React)

```jsx
// En App.jsx
console.log('Enviando mensaje:', message)
console.log('Respuesta recibida:', response)
console.log('Estado actual:', messages)
```

### Network Inspector

1. Abrir DevTools (F12)
2. Ir a Network
3. Enviar mensaje
4. Click en la petición
5. Ver:
   - **Headers**: método, URL, headers
   - **Payload**: datos enviados
   - **Response**: datos recibidos

---

## 📊 Flujo Completo de una Petición

```
1. Usuario escribe "Hola" en React
   ↓
2. handleSubmit() en App.jsx
   ↓
3. fetch('http://localhost:8080/api/v1/chat', {...})
   ↓
4. [NETWORK] HTTP POST con JSON
   ↓
5. Backend Go recibe en puerto 8080
   ↓
6. Router detecta /api/v1/chat
   ↓
7. handler.HandleChat() procesa
   ↓
8. chatService.SendMessage()
   ↓
9. groqClient.CreateChatCompletion()
   ↓
10. [NETWORK] HTTPS POST a api.groq.com
    ↓
11. Groq procesa con el modelo LLM
    ↓
12. [NETWORK] Respuesta JSON de Groq
    ↓
13. groqClient parsea respuesta
    ↓
14. chatService retorna al handler
    ↓
15. handler mapea a DTO y serializa JSON
    ↓
16. [NETWORK] HTTP Response 200 OK
    ↓
17. fetch() en React recibe respuesta
    ↓
18. await response.json() parsea
    ↓
19. setMessages() actualiza estado
    ↓
20. React re-renderiza UI
    ↓
21. Usuario ve la respuesta
```

---

## 🚀 Optimizaciones

### 1. Caché de Modelos

```jsx
// App.jsx
const [models, setModels] = useState([])

useEffect(() => {
  fetch('http://localhost:8080/api/v1/models')
    .then(r => r.json())
    .then(data => setModels(data.models))
}, [])
```

### 2. Loading States Granulares

```jsx
const [isSending, setIsSending] = useState(false)
const [isLoadingModels, setIsLoadingModels] = useState(false)
```

### 3. Retry Logic

```jsx
const sendWithRetry = async (message, retries = 3) => {
  for (let i = 0; i < retries; i++) {
    try {
      return await sendMessageToAPI(message)
    } catch (error) {
      if (i === retries - 1) throw error
      await new Promise(r => setTimeout(r, 1000 * (i + 1)))
    }
  }
}
```

### 4. Timeout

```jsx
const fetchWithTimeout = (url, options, timeout = 30000) => {
  return Promise.race([
    fetch(url, options),
    new Promise((_, reject) =>
      setTimeout(() => reject(new Error('Timeout')), timeout)
    )
  ])
}
```

---

## 🏭 Preparación para Producción

### Backend

```bash
# Compilar
go build -o groq-api cmd/api/main.go

# Ejecutar
./groq-api
```

### Frontend

```bash
# Build de producción
npm run build

# Los archivos están en dist/
# Subir a un servidor estático o CDN
```

### Variables de Entorno

**Backend (.env):**
```bash
PORT=8080
GROQ_API_KEY=tu_api_key_real
GROQ_BASE_URL=https://api.groq.com/openai/v1
```

**Frontend:**
Crear `.env.production`:
```bash
VITE_API_URL=https://tu-backend.com
```

Actualizar en App.jsx:
```jsx
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
```

### Docker

**Backend:**
```bash
cd groq-hexagonal-api
docker build -t groq-backend .
docker run -p 8080:8080 --env-file .env groq-backend
```

**Frontend:**
```dockerfile
# Dockerfile para frontend
FROM node:18-alpine as build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

```bash
cd groq-frontend
docker build -t groq-frontend .
docker run -p 80:80 groq-frontend
```

---

## ✅ Checklist de Integración

### Desarrollo
- [ ] Backend ejecutándose en puerto 8080
- [ ] Frontend ejecutándose en puerto 3000
- [ ] CORS configurado correctamente
- [ ] API key de Groq configurada
- [ ] Health check funciona
- [ ] Puedo enviar mensajes
- [ ] Las respuestas llegan correctamente

### Producción
- [ ] Variables de entorno configuradas
- [ ] CORS restringido a dominio específico
- [ ] HTTPS configurado
- [ ] Build de frontend optimizado
- [ ] Backend compilado
- [ ] Logs configurados
- [ ] Monitoreo configurado

---

## 📚 Próximos Pasos

Una vez que la integración funcione:

1. **Streaming de Respuestas** (Server-Sent Events)
2. **Autenticación** (JWT tokens)
3. **Base de Datos** (PostgreSQL para persistir conversaciones)
4. **WebSockets** para chat en tiempo real
5. **Tests** (Jest para frontend, Go testing para backend)
6. **CI/CD** (GitHub Actions)
7. **Monitoring** (Prometheus, Grafana)

---

¡Ya tienes una aplicación fullstack completa funcionando! 🎉
