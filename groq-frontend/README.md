# 🎨 Groq Chat Frontend

Interfaz web moderna para interactuar con la API de Groq construida con React.

## 🚀 Características

- ✅ Chat interactivo en tiempo real
- ✅ Selector de modelos de IA
- ✅ Diseño moderno y responsive
- ✅ Modo oscuro/claro
- ✅ Indicador de escritura
- ✅ Historial de conversación
- ✅ Sugerencias rápidas
- ✅ Información de tokens usados

## 📋 Requisitos Previos

- Node.js 18+ instalado
- Backend de Go ejecutándose en http://localhost:8080

## 🔧 Instalación

```bash
# 1. Instalar dependencias
npm install

# 2. Iniciar servidor de desarrollo
npm run dev

# La aplicación estará disponible en http://localhost:3000
```

## 🏗️ Estructura del Proyecto

```
groq-frontend/
├── index.html              # HTML principal
├── package.json            # Dependencias y scripts
├── vite.config.js          # Configuración de Vite
└── src/
    ├── main.jsx            # Punto de entrada React
    ├── App.jsx             # Componente principal
    └── App.css             # Estilos
```

## 🎯 Tecnologías Utilizadas

### Core
- **React 18** - Biblioteca de UI
- **Vite** - Build tool ultrarrápido
- **Lucide React** - Iconos modernos

### Características de React Usadas
- `useState` - Manejo de estado
- `useEffect` - Efectos secundarios
- `useRef` - Referencias al DOM
- Event Handlers
- Conditional Rendering
- Lists & Keys

## 📡 Comunicación con la API

El frontend se comunica con el backend Go mediante:

```javascript
// Enviar mensaje
POST http://localhost:8080/api/v1/chat
Body: {
  "message": "Tu mensaje aquí",
  "model": "llama-3.3-70b-versatile"
}

// Respuesta
{
  "success": true,
  "message": "Respuesta del modelo",
  "model": "llama-3.3-70b-versatile",
  "usage": {
    "prompt_tokens": 10,
    "completion_tokens": 20,
    "total_tokens": 30
  }
}
```

## 🎨 Personalización

### Cambiar Colores

Edita las variables CSS en `src/App.css`:

```css
:root {
  --primary: #6366f1;        /* Color principal */
  --primary-dark: #4f46e5;   /* Color principal oscuro */
  --primary-light: #818cf8;  /* Color principal claro */
  /* ... más variables ... */
}
```

### Agregar Nuevos Modelos

Edita el array `AVAILABLE_MODELS` en `src/App.jsx`:

```javascript
const AVAILABLE_MODELS = [
  { id: 'tu-modelo', name: 'Nombre del Modelo' },
  // ...
]
```

### Cambiar Puerto

Edita `vite.config.js`:

```javascript
server: {
  port: 3000, // Cambia este número
}
```

## 🔍 Conceptos de React Explicados

### 1. Estado (State)
Datos que pueden cambiar y causan re-renderizado:

```javascript
const [messages, setMessages] = useState([])
// messages: valor actual
// setMessages: función para actualizar
```

### 2. Efectos (Effects)
Código que se ejecuta en momentos específicos:

```javascript
useEffect(() => {
  // Código a ejecutar
}, [dependencias]) // Se ejecuta cuando las dependencias cambian
```

### 3. Referencias (Refs)
Acceso directo a elementos del DOM:

```javascript
const inputRef = useRef(null)
// Usar: inputRef.current.focus()
```

### 4. Event Handlers
Funciones que responden a eventos del usuario:

```javascript
const handleSubmit = (e) => {
  e.preventDefault() // Prevenir recarga
  // Tu lógica aquí
}
```

### 5. Renderizado Condicional
Mostrar/ocultar elementos según condiciones:

```javascript
{isLoading && <LoadingSpinner />}
{messages.length === 0 ? <EmptyState /> : <MessageList />}
```

## 📚 Estructura del Código

### App.jsx - Secciones Principales

```javascript
function App() {
  // 1. ESTADO - Variables que cambian
  const [messages, setMessages] = useState([])
  const [inputMessage, setInputMessage] = useState('')
  
  // 2. EFECTOS - Acciones automáticas
  useEffect(() => {
    scrollToBottom()
  }, [messages])
  
  // 3. FUNCIONES - Lógica de la app
  const sendMessageToAPI = async (message) => {
    // Llamada a la API
  }
  
  const handleSubmit = (e) => {
    // Enviar mensaje
  }
  
  // 4. RENDERIZADO - UI
  return (
    <div className="app">
      {/* JSX aquí */}
    </div>
  )
}
```

## 🎓 Flujo de una Interacción

```
1. Usuario escribe mensaje
   ↓
2. handleInputChange actualiza inputMessage
   ↓
3. Usuario presiona Enter o botón Enviar
   ↓
4. handleSubmit previene reload
   ↓
5. Crear mensaje del usuario
   ↓
6. Actualizar estado: setMessages([...prev, userMessage])
   ↓
7. Activar loading: setIsLoading(true)
   ↓
8. Llamar API: sendMessageToAPI(message)
   ↓
9. Esperar respuesta
   ↓
10. Crear mensaje del asistente con respuesta
    ↓
11. Actualizar estado: setMessages([...prev, assistantMessage])
    ↓
12. Desactivar loading: setIsLoading(false)
    ↓
13. React re-renderiza el componente con nuevos mensajes
    ↓
14. useEffect detecta cambio en messages
    ↓
15. scrollToBottom() se ejecuta
```

## 🐛 Debugging

### Ver Estado en React DevTools
1. Instalar React DevTools (extensión del navegador)
2. Abrir DevTools → Components
3. Ver el estado de cada componente

### Console.log Útiles
```javascript
// Ver mensajes
console.log('Mensajes:', messages)

// Ver respuesta de API
console.log('Respuesta:', response)

// Ver errores
console.error('Error:', error)
```

## 🚀 Scripts Disponibles

```bash
# Desarrollo (con hot reload)
npm run dev

# Compilar para producción
npm run build

# Preview de la build de producción
npm run preview

# Linting (verificar código)
npm run lint
```

## 📦 Dependencias

```json
{
  "react": "^18.2.0",           // Biblioteca de UI
  "react-dom": "^18.2.0",       // React para DOM
  "lucide-react": "^0.294.0"    // Iconos
}
```

## 🎨 Características de UI

### Animaciones
- Fade in para mensajes nuevos
- Indicador de escritura animado
- Transiciones suaves en hover
- Scroll automático

### Responsive Design
- Desktop: Layout amplio con sidebar potencial
- Tablet: Layout adaptado
- Mobile: Layout de una columna, botones optimizados

### Accesibilidad
- Contraste de colores adecuado
- Focus visible en elementos interactivos
- Aria labels donde sea necesario
- Navegación por teclado

## 🔒 Seguridad

### Validaciones
- Input no vacío antes de enviar
- Sanitización de HTML (React lo hace automáticamente)
- Manejo de errores de red

### Mejores Prácticas
- No guardar API keys en el frontend
- Validar respuestas de la API
- Manejar errores gracefully

## 🎯 Próximas Mejoras

Ideas para extender la aplicación:

- [ ] Persistencia local (localStorage)
- [ ] Exportar conversación (JSON/TXT)
- [ ] Streaming de respuestas (Server-Sent Events)
- [ ] Markdown rendering para respuestas
- [ ] Syntax highlighting para código
- [ ] Soporte para imágenes
- [ ] Historial de conversaciones múltiples
- [ ] Búsqueda en conversaciones
- [ ] Temas personalizables
- [ ] Atajos de teclado
- [ ] PWA (Progressive Web App)

## 📖 Recursos de Aprendizaje

### React
- [React Docs](https://react.dev/) - Documentación oficial
- [React Tutorial](https://react.dev/learn) - Tutorial interactivo

### JavaScript Moderno
- [MDN Web Docs](https://developer.mozilla.org/) - Referencia completa
- [JavaScript.info](https://javascript.info/) - Tutorial completo

### Vite
- [Vite Docs](https://vitejs.dev/) - Build tool

## 💡 Tips de Desarrollo

1. **Hot Reload**: Los cambios se reflejan automáticamente
2. **React DevTools**: Instala la extensión para debugging
3. **Console**: Usa console.log para ver valores
4. **Errores**: Lee los mensajes de error, son muy descriptivos
5. **ESLint**: Presta atención a los warnings

## 🤝 Integración con Backend

Asegúrate de que:
1. El backend esté ejecutándose en `http://localhost:8080`
2. CORS esté configurado correctamente en el backend
3. Los endpoints coincidan (`/api/v1/chat`)
4. Los formatos de request/response sean compatibles

## ✅ Checklist de Setup

- [ ] Node.js instalado
- [ ] Dependencias instaladas (`npm install`)
- [ ] Backend ejecutándose
- [ ] Puerto 3000 disponible
- [ ] Navegador moderno (Chrome, Firefox, Safari, Edge)

¡Disfruta construyendo tu aplicación! 🚀
