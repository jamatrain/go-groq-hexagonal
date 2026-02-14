// ============================================================================
// APP.JSX - Componente Principal de la Aplicación
// ============================================================================
//
// Este es el componente raíz de nuestra aplicación React
// Maneja toda la lógica del chat y la comunicación con la API
//
// ============================================================================

import { useState, useEffect, useRef } from 'react'
import { 
  Send, 
  Bot, 
  User, 
  Loader2, 
  Moon, 
  Sun, 
  Sparkles,
  MessageSquare,
  Trash2,
  RefreshCw
} from 'lucide-react'
import './App.css'

// ============================================================================
// CONFIGURACIÓN
// ============================================================================

const API_BASE_URL = 'http://localhost:8080'

// Modelos disponibles (esto también podría venir de la API)
const AVAILABLE_MODELS = [
  { id: 'llama-3.3-70b-versatile', name: 'Llama 3.3 70B (Recomendado)' },
  { id: 'llama-3.1-8b-instant', name: 'Llama 3.1 8B (Rápido)' },
  { id: 'mixtral-8x7b-32768', name: 'Mixtral 8x7B' },
]

// ============================================================================
// COMPONENTE PRINCIPAL
// ============================================================================

function App() {
  // ==========================================================================
  // ESTADO (State)
  // ==========================================================================
  // En React, el estado es información que puede cambiar y que dispara
  // re-renderizado cuando cambia. Usamos useState para crear estado.
  
  // Lista de mensajes del chat
  const [messages, setMessages] = useState([])
  
  // Texto que el usuario está escribiendo
  const [inputMessage, setInputMessage] = useState('')
  
  // Modelo de IA seleccionado
  const [selectedModel, setSelectedModel] = useState(AVAILABLE_MODELS[0].id)
  
  // Estado de carga (cuando esperamos respuesta de la API)
  const [isLoading, setIsLoading] = useState(false)
  
  // Tema (claro/oscuro)
  const [isDarkMode, setIsDarkMode] = useState(false)
  
  // Estado de error
  const [error, setError] = useState(null)
  
  // ==========================================================================
  // REFS
  // ==========================================================================
  // useRef nos permite acceder a elementos del DOM o mantener valores
  // entre renders sin disparar re-renderizado
  
  // Referencia al contenedor de mensajes para hacer scroll automático
  const messagesEndRef = useRef(null)
  
  // Referencia al textarea para hacer focus
  const inputRef = useRef(null)

  // ==========================================================================
  // EFECTOS (Effects)
  // ==========================================================================
  // useEffect ejecuta código cuando el componente se monta o cuando
  // cambian ciertas dependencias
  
  // Scroll automático al último mensaje
  useEffect(() => {
    scrollToBottom()
  }, [messages])
  
  // Aplicar tema oscuro/claro
  useEffect(() => {
    if (isDarkMode) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }, [isDarkMode])
  
  // Focus en el input al cargar
  useEffect(() => {
    inputRef.current?.focus()
  }, [])

  // ==========================================================================
  // FUNCIONES AUXILIARES
  // ==========================================================================
  
  // Hacer scroll al final de los mensajes
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }
  
  // Generar ID único para mensajes
  const generateId = () => {
    return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
  }
  
  // Formatear fecha/hora
  const formatTime = (date) => {
    return new Intl.DateTimeFormat('es', {
      hour: '2-digit',
      minute: '2-digit'
    }).format(date)
  }

  // ==========================================================================
  // COMUNICACIÓN CON LA API
  // ==========================================================================
  
  // Enviar mensaje a la API
  const sendMessageToAPI = async (message, model) => {
    // Realizamos una petición POST a nuestra API de Go
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
    
    // Si la respuesta no es OK, lanzar error
    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Error al comunicarse con la API')
    }
    
    // Parsear y retornar la respuesta
    return await response.json()
  }

  // ==========================================================================
  // MANEJADORES DE EVENTOS (Event Handlers)
  // ==========================================================================
  
  // Manejar envío del formulario
  const handleSubmit = async (e) => {
    // Prevenir el comportamiento por defecto del form (recargar página)
    e.preventDefault()
    
    // Validar que hay mensaje
    if (!inputMessage.trim()) return
    
    // Limpiar error previo
    setError(null)
    
    // Crear mensaje del usuario
    const userMessage = {
      id: generateId(),
      role: 'user',
      content: inputMessage,
      timestamp: new Date()
    }
    
    // Añadir mensaje del usuario a la lista
    setMessages(prev => [...prev, userMessage])
    
    // Guardar el mensaje antes de limpiar el input
    const messageToSend = inputMessage
    
    // Limpiar el input
    setInputMessage('')
    
    // Activar estado de carga
    setIsLoading(true)
    
    try {
      // Enviar mensaje a la API
      const response = await sendMessageToAPI(messageToSend, selectedModel)
      
      // Crear mensaje del asistente con la respuesta
      const assistantMessage = {
        id: generateId(),
        role: 'assistant',
        content: response.message,
        model: response.model,
        timestamp: new Date(),
        usage: response.usage
      }
      
      // Añadir respuesta del asistente
      setMessages(prev => [...prev, assistantMessage])
      
    } catch (err) {
      // Manejar error
      console.error('Error al enviar mensaje:', err)
      setError(err.message)
      
      // Crear mensaje de error
      const errorMessage = {
        id: generateId(),
        role: 'assistant',
        content: `❌ Error: ${err.message}`,
        timestamp: new Date(),
        isError: true
      }
      
      setMessages(prev => [...prev, errorMessage])
      
    } finally {
      // Desactivar estado de carga
      setIsLoading(false)
      
      // Volver a hacer focus en el input
      inputRef.current?.focus()
    }
  }
  
  // Manejar cambio en el textarea
  const handleInputChange = (e) => {
    setInputMessage(e.target.value)
  }
  
  // Manejar tecla Enter (enviar) vs Shift+Enter (nueva línea)
  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSubmit(e)
    }
  }
  
  // Limpiar chat
  const handleClearChat = () => {
    if (window.confirm('¿Estás seguro de que quieres borrar toda la conversación?')) {
      setMessages([])
      setError(null)
    }
  }
  
  // Alternar tema
  const toggleTheme = () => {
    setIsDarkMode(prev => !prev)
  }
  
  // Manejar sugerencia rápida
  const handleSuggestion = (text) => {
    setInputMessage(text)
    inputRef.current?.focus()
  }

  // ==========================================================================
  // RENDERIZADO (Render)
  // ==========================================================================
  
  return (
    <div className="app">
      {/* ====================================================================
          HEADER - Barra superior
          ==================================================================== */}
      <header className="header">
        <div className="header-title">
          <Sparkles size={28} />
          <div>
            <div>Groq Chat</div>
            <div className="header-subtitle">Powered by Groq API</div>
          </div>
        </div>
        
        <div style={{ display: 'flex', gap: '0.5rem' }}>
          {/* Botón para limpiar chat */}
          {messages.length > 0 && (
            <button
              className="icon-button"
              onClick={handleClearChat}
              title="Limpiar conversación"
            >
              <Trash2 size={20} />
            </button>
          )}
          
          {/* Botón para cambiar tema */}
          <button
            className="icon-button"
            onClick={toggleTheme}
            title={isDarkMode ? 'Modo claro' : 'Modo oscuro'}
          >
            {isDarkMode ? <Sun size={20} /> : <Moon size={20} />}
          </button>
        </div>
      </header>

      {/* ====================================================================
          ÁREA DE CHAT - Mensajes
          ==================================================================== */}
      <div className="chat-container">
        <div className="messages-container">
          <div className="messages-wrapper">
            {/* Estado vacío - Cuando no hay mensajes */}
            {messages.length === 0 ? (
              <div className="empty-state">
                <Bot size={64} className="empty-state-icon" />
                <h2 className="empty-state-title">¡Hola! Soy tu asistente de IA</h2>
                <p className="empty-state-description">
                  Pregúntame cualquier cosa. Puedo ayudarte con programación,
                  explicaciones, creatividad y mucho más.
                </p>
                
                {/* Sugerencias rápidas */}
                <div className="suggestions">
                  <button
                    className="suggestion-card"
                    onClick={() => handleSuggestion('Explica qué es la arquitectura hexagonal')}
                  >
                    <div className="suggestion-title">🏗️ Arquitectura</div>
                    <div className="suggestion-text">Arquitectura hexagonal</div>
                  </button>
                  
                  <button
                    className="suggestion-card"
                    onClick={() => handleSuggestion('¿Cuáles son las ventajas de usar Go?')}
                  >
                    <div className="suggestion-title">💻 Programación</div>
                    <div className="suggestion-text">Ventajas de Go</div>
                  </button>
                  
                  <button
                    className="suggestion-card"
                    onClick={() => handleSuggestion('Escribe un poema corto sobre la tecnología')}
                  >
                    <div className="suggestion-title">✨ Creatividad</div>
                    <div className="suggestion-text">Poema sobre tecnología</div>
                  </button>
                </div>
              </div>
            ) : (
              /* Lista de mensajes */
              <>
                {messages.map((message) => (
                  <div
                    key={message.id}
                    className={`message ${message.role}`}
                  >
                    {/* Avatar */}
                    <div className="message-avatar">
                      {message.role === 'user' ? (
                        <User size={20} />
                      ) : (
                        <Bot size={20} />
                      )}
                    </div>
                    
                    {/* Contenido del mensaje */}
                    <div className="message-content">
                      <div className="message-text">{message.content}</div>
                      
                      {/* Información adicional */}
                      <div className="message-info">
                        <span>{formatTime(message.timestamp)}</span>
                        {message.model && (
                          <>
                            <span>•</span>
                            <span>{message.model}</span>
                          </>
                        )}
                        {message.usage && (
                          <>
                            <span>•</span>
                            <span>{message.usage.total_tokens} tokens</span>
                          </>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
                
                {/* Indicador de "escribiendo..." */}
                {isLoading && (
                  <div className="message assistant">
                    <div className="message-avatar">
                      <Bot size={20} />
                    </div>
                    <div className="typing-indicator">
                      <div className="typing-dot"></div>
                      <div className="typing-dot"></div>
                      <div className="typing-dot"></div>
                    </div>
                  </div>
                )}
                
                {/* Elemento invisible para scroll automático */}
                <div ref={messagesEndRef} />
              </>
            )}
          </div>
        </div>

        {/* ================================================================
            ÁREA DE INPUT - Formulario de envío
            ================================================================ */}
        <div className="input-area">
          <div className="input-wrapper">
            {/* Selector de modelo */}
            <div className="input-controls">
              <MessageSquare size={20} style={{ color: 'var(--text-secondary)' }} />
              <select
                className="model-select"
                value={selectedModel}
                onChange={(e) => setSelectedModel(e.target.value)}
                disabled={isLoading}
              >
                {AVAILABLE_MODELS.map(model => (
                  <option key={model.id} value={model.id}>
                    {model.name}
                  </option>
                ))}
              </select>
            </div>
            
            {/* Formulario de envío */}
            <form onSubmit={handleSubmit} className="input-form">
              <textarea
                ref={inputRef}
                className="input-field"
                value={inputMessage}
                onChange={handleInputChange}
                onKeyDown={handleKeyDown}
                placeholder="Escribe tu mensaje aquí... (Enter para enviar, Shift+Enter para nueva línea)"
                disabled={isLoading}
                rows={1}
              />
              
              <button
                type="submit"
                className="send-button"
                disabled={isLoading || !inputMessage.trim()}
              >
                {isLoading ? (
                  <>
                    <Loader2 size={20} className="spin" />
                    <span>Enviando...</span>
                  </>
                ) : (
                  <>
                    <Send size={20} />
                    <span>Enviar</span>
                  </>
                )}
              </button>
            </form>
            
            {/* Mensaje de error */}
            {error && (
              <div style={{
                padding: '0.75rem',
                backgroundColor: '#fee2e2',
                color: '#991b1b',
                borderRadius: 'var(--radius-md)',
                fontSize: '0.875rem'
              }}>
                {error}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

// ============================================================================
// CONCEPTOS DE REACT EXPLICADOS:
// ============================================================================
//
// 1. COMPONENTES:
//    - Funciones que retornan JSX (parecido a HTML)
//    - Pueden tener estado y lógica
//    - Se pueden reutilizar
//
// 2. JSX:
//    - Sintaxis similar a HTML pero en JavaScript
//    - Permite usar expresiones JS con {}
//    - className en vez de class
//    - Eventos con camelCase: onClick, onChange
//
// 3. ESTADO (useState):
//    - Datos que pueden cambiar
//    - Cuando cambian, el componente se re-renderiza
//    - const [valor, setValor] = useState(inicial)
//
// 4. EFECTOS (useEffect):
//    - Ejecutar código cuando algo cambia
//    - useEffect(() => { código }, [dependencias])
//    - [] = solo al montar, [x] = cuando x cambia
//
// 5. REFS (useRef):
//    - Acceder a elementos del DOM
//    - Mantener valores entre renders sin re-renderizar
//    - inputRef.current es el elemento
//
// 6. EVENT HANDLERS:
//    - Funciones que responden a eventos del usuario
//    - onClick, onChange, onSubmit, etc.
//    - e.preventDefault() previene comportamiento default
//
// 7. RENDERIZADO CONDICIONAL:
//    - {condicion && <Componente />} - render si true
//    - {condicion ? <A /> : <B />} - render A o B
//
// 8. LISTAS:
//    - array.map() para renderizar múltiples elementos
//    - Cada elemento necesita una key única
//
// 9. ASYNC/AWAIT:
//    - Para operaciones asíncronas (API calls)
//    - try/catch para manejar errores
//    - fetch() para hacer HTTP requests
//
// ============================================================================

export default App
