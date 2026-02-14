// ============================================================================
// MAIN.JSX - Punto de Entrada de la Aplicación React
// ============================================================================
//
// Este archivo es el punto de entrada de nuestra aplicación React
// Monta el componente App en el DOM
//
// ============================================================================

import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.jsx'

// ============================================================================
// RENDERIZADO DE LA APLICACIÓN
// ============================================================================

// ReactDOM.createRoot() crea un "root" de React
// document.getElementById('root') busca el elemento con id="root" en index.html
ReactDOM.createRoot(document.getElementById('root')).render(
  // React.StrictMode es un wrapper que ayuda a detectar problemas
  // Solo funciona en desarrollo, no en producción
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)

// ============================================================================
// CONCEPTOS EXPLICADOS:
// ============================================================================
//
// 1. REACT DOM:
//    - Biblioteca que conecta React con el DOM del navegador
//    - createRoot() crea un punto de montaje para React
//    - render() renderiza el componente en el DOM
//
// 2. STRICT MODE:
//    - Herramienta de desarrollo para detectar problemas
//    - Activa advertencias adicionales
//    - Hace render doble en desarrollo (para detectar efectos)
//    - Se quita automáticamente en producción
//
// 3. IMPORTS:
//    - import React from 'react' - Biblioteca principal
//    - import ReactDOM from 'react-dom/client' - Para renderizar
//    - import App from './App.jsx' - Nuestro componente
//
// ============================================================================
