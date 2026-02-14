@echo off
setlocal enabledelayedexpansion

title Groq App Auto-Launcher

:: Ensure we are running from the script's directory
cd /d "%~dp0"

:: Check if we are in the parent directory by mistake (user placed script outside)
if exist "go-groq-hexagonal" (
    echo [INFO] Detected execution from parent directory. Entering project folder...
    cd "go-groq-hexagonal"
)

echo ===================================================
echo   Groq Hexagonal App - One Click Installer & Run
echo ===================================================
echo.

:: -------------------------
:: 0. ENVIRONMENT CHECK
:: -------------------------
echo [INFO] Checking environment...

:: Check Node (try to fix path if missing)
node -v >nul 2>&1
if %errorlevel% neq 0 (
    echo [WARN] Node.js not found in PATH. Attempting to add standard path...
    set "PATH=C:\Program Files\nodejs;%PATH%"
)

:: Re-check Node
node -v >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Node.js is still not found.
    echo         Please install Node.js from https://nodejs.org/
    echo         or restart your computer if you just installed it.
    pause
    exit /b 1
) else (
    echo [OK] Node.js found.
)

:: Check Go
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Go is not found. Please install Go from https://go.dev/dl/
    pause
    exit /b 1
) else (
    echo [OK] Go found.
)

:: -------------------------
:: 1. BUILD BACKEND
:: -------------------------
echo.
echo [1/4] Building Backend (Go)...
if not exist "groq-hexagonal-api" (
    echo [ERROR] Directory 'groq-hexagonal-api' not found in %CD%
    echo         Please ensure the script is in the root of the project.
    dir
    pause
    exit /b 1
)
cd "groq-hexagonal-api"

if not exist .env (
    echo [WARN] .env file not found! Copying .env.example...
    copy .env.example .env
    echo [IMPORTANT] Please update .env with your GROQ_API_KEY later.
)

if not exist ..\bin mkdir ..\bin
go build -o ..\bin\groq-api.exe cmd/api/main.go
if %errorlevel% neq 0 (
    echo [ERROR] Failed to compile backend.
    echo         Current directory: %CD%
    echo         Trying to build: cmd/api/main.go
    pause
    exit /b 1
)
cd ..
echo [OK] Backend built successfully.

:: -------------------------
:: 2. SETUP FRONTEND
:: -------------------------
echo.
echo [2/4] Installing Frontend Dependencies...
cd groq-frontend
call npm install
if %errorlevel% neq 0 (
    echo [ERROR] Failed to install npm dependencies.
    pause
    exit /b 1
)

echo.
echo [3/4] Building Frontend...
call npm run build
if %errorlevel% neq 0 (
    echo [ERROR] Failed to build frontend.
    pause
    exit /b 1
)
cd ..
echo [OK] Frontend built successfully.

:: -------------------------
:: 3. LAUNCH APPS
:: -------------------------
echo.
echo [4/4] Launching Applications...

:: Launch Backend
cd groq-hexagonal-api
start "Groq Backend API" ..\bin\groq-api.exe

:: Launch Frontend (Preview)
cd ..\groq-frontend
:: 'preview' serves the 'dist' folder created by 'build'
:: Use /k to keep window open if it fails
start "Groq Frontend" cmd /k "npm run preview"

:: Open Browser (Wait a few seconds for servers to start)
echo.
echo Waiting for servers to start...
timeout /t 5 /nobreak >nul
echo Opening browser...
start http://localhost:4173

echo.
echo ===================================================
echo [SUCCESS] All services started!
echo ===================================================
echo.
echo Backend URL:   http://localhost:8080
echo Frontend URL:  http://localhost:4173
echo.
echo Close this window to keep servers running, or close the spawned windows to stop them.
echo.
pause
