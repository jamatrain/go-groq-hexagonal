# 📚 Guía de Aprendizaje - React y Frontend

Esta guía te ayudará a entender paso a paso cómo funciona la aplicación React.

## 🎯 Objetivos

Al finalizar, comprenderás:
1. ✅ Fundamentos de React (componentes, estado, efectos)
2. ✅ Comunicación con APIs REST
3. ✅ Manejo de eventos en React
4. ✅ Diseño responsive con CSS
5. ✅ Arquitectura de una aplicación frontend moderna

---

## 📖 Parte 1: Conceptos Fundamentales de React

### 1.1 ¿Qué es React?

React es una biblioteca de JavaScript para construir interfaces de usuario. Características principales:

- **Componentes**: Piezas reutilizables de UI
- **Estado**: Datos que pueden cambiar
- **Virtual DOM**: Actualizaciones eficientes de la UI
- **JSX**: Sintaxis similar a HTML en JavaScript
- **Unidireccional**: Flujo de datos predecible

### 1.2 JSX - JavaScript XML

JSX parece HTML pero es JavaScript:

```jsx
// JSX
const element = <h1 className="title">Hola</h1>

// Se convierte en:
const element = React.createElement('h1', {className: 'title'}, 'Hola')
```

**Reglas importantes:**
- `className` en vez de `class`
- `htmlFor` en vez de `for`
- Eventos en camelCase: `onClick`, `onChange`
- Cerrar todas las etiquetas: `<img />`, `<input />`
- Usar `{}` para expresiones JavaScript

**Ejemplos:**
```jsx
// Variables
const nombre = "Juan"
<h1>Hola {nombre}</h1>

// Expresiones
<p>{2 + 2}</p>  // Muestra: 4

// Condicionales
{isLoading && <Spinner />}
{isDark ? <Moon /> : <Sun />}

// Arrays
{items.map(item => <div key={item.id}>{item.name}</div>)}
```

### 1.3 Componentes

Un componente es una función que retorna JSX:

```jsx
// Componente simple
function Saludo() {
  return <h1>¡Hola!</h1>
}

// Componente con props
function Saludo({ nombre }) {
  return <h1>¡Hola {nombre}!</h1>
}

// Uso
<Saludo nombre="María" />
```

**Lee:** `src/App.jsx` - El componente principal

### 1.4 Estado (useState)

El estado son datos que pueden cambiar:

```jsx
import { useState } from 'react'

function Contador() {
  // [valor, función_para_cambiar] = useState(valor_inicial)
  const [count, setCount] = useState(0)
  
  return (
    <div>
      <p>Clicks: {count}</p>
      <button onClick={() => setCount(count + 1)}>
        Incrementar
      </button>
    </div>
  )
}
```

**Puntos clave:**
- `useState` retorna un array de 2 elementos
- Primer elemento: valor actual
- Segundo elemento: función para actualizar
- Cuando actualizas el estado, React re-renderiza el componente

**Ejemplo en nuestra app:**
```jsx
// En App.jsx
const [messages, setMessages] = useState([])
const [inputMessage, setInputMessage] = useState('')

// Añadir mensaje
setMessages([...messages, newMessage])

// Actualizar input
setInputMessage('Nuevo texto')
```

**Ejercicio:** Encuentra todos los `useState` en `App.jsx` y entiende qué controla cada uno.

### 1.5 Efectos (useEffect)

Los efectos ejecutan código cuando algo cambia:

```jsx
import { useEffect } from 'react'

function Ejemplo() {
  const [count, setCount] = useState(0)
  
  // Se ejecuta después de cada render
  useEffect(() => {
    console.log('Componente renderizado')
  })
  
  // Se ejecuta solo al montar
  useEffect(() => {
    console.log('Componente montado')
  }, [])
  
  // Se ejecuta cuando count cambia
  useEffect(() => {
    console.log('Count cambió:', count)
  }, [count])
  
  // Cleanup
  useEffect(() => {
    const timer = setInterval(() => {
      console.log('Tick')
    }, 1000)
    
    // Esta función se ejecuta al desmontar
    return () => {
      clearInterval(timer)
    }
  }, [])
}
```

**En nuestra app:**
```jsx
// Scroll automático cuando cambian los mensajes
useEffect(() => {
  scrollToBottom()
}, [messages])

// Aplicar tema
useEffect(() => {
  if (isDarkMode) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}, [isDarkMode])
```

### 1.6 Referencias (useRef)

Para acceder a elementos del DOM o guardar valores sin re-renderizar:

```jsx
import { useRef } from 'react'

function Ejemplo() {
  const inputRef = useRef(null)
  
  const hacerFocus = () => {
    // Acceder al elemento del DOM
    inputRef.current.focus()
  }
  
  return (
    <div>
      <input ref={inputRef} />
      <button onClick={hacerFocus}>Focus</button>
    </div>
  )
}
```

**Diferencias con useState:**
- `useRef` NO causa re-renderizado
- `useState` SÍ causa re-renderizado
- Usa `useRef` para valores que no afectan la UI

**En nuestra app:**
```jsx
const inputRef = useRef(null)
const messagesEndRef = useRef(null)

// Hacer focus en el input
inputRef.current?.focus()

// Scroll al final
messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
```

---

## 📖 Parte 2: Comunicación con APIs

### 2.1 Fetch API

Fetch es la forma estándar de hacer peticiones HTTP:

```jsx
// GET simple
fetch('https://api.ejemplo.com/datos')
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error(error))

// POST con async/await
async function enviarDatos() {
  try {
    const response = await fetch('https://api.ejemplo.com/datos', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ mensaje: 'Hola' })
    })
    
    if (!response.ok) {
      throw new Error('Error en la petición')
    }
    
    const data = await response.json()
    return data
  } catch (error) {
    console.error('Error:', error)
    throw error
  }
}
```

### 2.2 En Nuestra Aplicación

**Lee:** `src/App.jsx` función `sendMessageToAPI`

```jsx
const sendMessageToAPI = async (message, model) => {
  // 1. Hacer petición POST
  const response = await fetch(`${API_BASE_URL}/api/v1/chat`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      message: message,
      model: model
    })
  })
  
  // 2. Verificar respuesta
  if (!response.ok) {
    const errorData = await response.json()
    throw new Error(errorData.error || 'Error al comunicarse con la API')
  }
  
  // 3. Parsear y retornar
  return await response.json()
}
```

### 2.3 Manejo de Errores

Siempre usa try-catch con async/await:

```jsx
const handleSubmit = async (e) => {
  e.preventDefault()
  
  setIsLoading(true)
  
  try {
    // Intentar enviar mensaje
    const response = await sendMessageToAPI(message, model)
    
    // Si llega aquí, fue exitoso
    setMessages([...messages, response])
    
  } catch (error) {
    // Si hubo error
    console.error('Error:', error)
    setError(error.message)
    
  } finally {
    // Siempre se ejecuta (éxito o error)
    setIsLoading(false)
  }
}
```

---

## 📖 Parte 3: Eventos en React

### 3.1 Event Handlers

```jsx
// onClick
<button onClick={handleClick}>Click</button>

// onChange
<input onChange={handleChange} />

// onSubmit
<form onSubmit={handleSubmit}>

// Con parámetros
<button onClick={() => handleClick(id)}>

// Event object
const handleChange = (e) => {
  console.log(e.target.value)
}
```

### 3.2 Prevenir Comportamiento Default

```jsx
const handleSubmit = (e) => {
  e.preventDefault() // Evita recargar la página
  // Tu código aquí
}

const handleLink = (e) => {
  e.preventDefault() // Evita navegar
  // Tu código aquí
}
```

### 3.3 En Nuestra App

```jsx
// Enviar formulario
const handleSubmit = (e) => {
  e.preventDefault()
  // Enviar mensaje
}

// Cambiar input
const handleInputChange = (e) => {
  setInputMessage(e.target.value)
}

// Enter vs Shift+Enter
const handleKeyDown = (e) => {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSubmit(e)
  }
}
```

---

## 📖 Parte 4: Estilización con CSS

### 4.1 CSS Variables

Definidas en `:root`:

```css
:root {
  --primary: #6366f1;
  --bg-primary: #ffffff;
}

.button {
  background-color: var(--primary);
}
```

### 4.2 Tema Oscuro

```css
.dark {
  --bg-primary: #111827;
  --text-primary: #f9fafb;
}
```

En JavaScript:
```jsx
const toggleTheme = () => {
  if (isDarkMode) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}
```

### 4.3 Flexbox

```css
.container {
  display: flex;
  flex-direction: column; /* vertical */
  gap: 1rem; /* espacio entre elementos */
  align-items: center; /* centro horizontal */
  justify-content: space-between; /* espacio vertical */
}
```

### 4.4 Grid

```css
.grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr); /* 3 columnas iguales */
  gap: 1rem;
}
```

### 4.5 Responsive

```css
@media (max-width: 768px) {
  .container {
    padding: 0.5rem;
  }
}
```

---

## 📖 Parte 5: Flujo Completo

### Ejemplo: Usuario Envía Mensaje

```
1. Usuario escribe en textarea
   → handleInputChange actualiza inputMessage
   
2. Usuario presiona Enter
   → handleKeyDown detecta Enter
   → handleSubmit se ejecuta
   
3. handleSubmit
   → e.preventDefault() evita reload
   → Crea objeto userMessage
   → setMessages añade userMessage
   → setInputMessage('') limpia input
   → setIsLoading(true) activa spinner
   
4. sendMessageToAPI
   → fetch() hace POST al backend
   → await espera respuesta
   → Parsea JSON
   
5. Respuesta exitosa
   → Crea assistantMessage
   → setMessages añade assistantMessage
   → setIsLoading(false) desactiva spinner
   
6. React detecta cambio en messages
   → Re-renderiza componente
   → map() renderiza cada mensaje
   
7. useEffect detecta cambio en messages
   → scrollToBottom() se ejecuta
   → Scroll automático al final
```

---

## 🎓 Ejercicios Prácticos

### Ejercicio 1: Agregar Contador de Caracteres
Muestra cuántos caracteres tiene el mensaje actual.

**Hint:**
```jsx
<div>{inputMessage.length} / 1000 caracteres</div>
```

### Ejercicio 2: Botón para Copiar Mensaje
Agrega un botón para copiar respuestas del asistente.

**Hint:**
```jsx
const copiarTexto = (texto) => {
  navigator.clipboard.writeText(texto)
  alert('Copiado!')
}
```

### Ejercicio 3: Guardar en localStorage
Persiste las conversaciones en localStorage.

**Hint:**
```jsx
// Guardar
useEffect(() => {
  localStorage.setItem('messages', JSON.stringify(messages))
}, [messages])

// Cargar
useEffect(() => {
  const saved = localStorage.getItem('messages')
  if (saved) {
    setMessages(JSON.parse(saved))
  }
}, [])
```

### Ejercicio 4: Markdown en Respuestas
Renderiza Markdown en las respuestas del asistente.

**Hint:**
```bash
npm install react-markdown
```

```jsx
import ReactMarkdown from 'react-markdown'

<ReactMarkdown>{message.content}</ReactMarkdown>
```

---

## 🐛 Debugging Tips

### 1. Console.log
```jsx
console.log('Estado actual:', messages)
console.log('Props:', { message, model })
```

### 2. React DevTools
- Instala extensión del navegador
- Inspecciona componentes
- Ve el estado en tiempo real

### 3. Debugger
```jsx
const handleSubmit = (e) => {
  debugger // Pausa ejecución aquí
  // ...
}
```

### 4. Network Tab
- Abre DevTools → Network
- Ve peticiones HTTP
- Inspecciona request/response

---

## 📚 Recursos Adicionales

### Documentación
- [React Docs](https://react.dev/) - Documentación oficial
- [MDN](https://developer.mozilla.org/) - JavaScript reference

### Tutoriales
- [React Tutorial](https://react.dev/learn/tutorial-tic-tac-toe)
- [JavaScript.info](https://javascript.info/)

### Práctica
- [CodeSandbox](https://codesandbox.io/) - Editor online
- [StackBlitz](https://stackblitz.com/) - IDE online

---

## ✅ Checklist de Aprendizaje

- [ ] Entiendo qué es JSX
- [ ] Sé crear componentes
- [ ] Comprendo useState
- [ ] Comprendo useEffect
- [ ] Sé manejar eventos
- [ ] Puedo hacer peticiones HTTP
- [ ] Entiendo async/await
- [ ] Sé usar CSS moderno
- [ ] Puedo agregar nuevas features
- [ ] Puedo debuggear problemas

---

¡Sigue practicando y experimentando! 🚀
