# 🌐 Configuración de ngrok para WhatsApp Clone

## Instalación de ngrok

### 1. Descargar ngrok
```bash
# Linux/Mac
wget https://bin.equinox.io/c/bNyj1mQVY4c/ngrok-v3-stable-linux-amd64.tgz
tar xvzf ngrok-v3-stable-linux-amd64.tgz

# O con snap
sudo snap install ngrok

# Windows: Descargar desde https://ngrok.com/download
```

### 2. Crear cuenta gratuita
- Ir a https://ngrok.com/signup
- Obtener tu authtoken

### 3. Configurar authtoken
```bash
ngrok config add-authtoken YOUR_AUTH_TOKEN
```

## Exponer la aplicación

### Opción 1: Exponer solo el frontend (Recomendado)
```bash
# Terminal 1: Backend
cd backend
go run main.go

# Terminal 2: Frontend
cd frontend
npm run dev

# Terminal 3: ngrok
ngrok http 3000
```

### Opción 2: Exponer backend y frontend por separado
```bash
# Terminal 1: Backend
cd backend
go run main.go

# Terminal 2: Frontend
cd frontend
npm run dev

# Terminal 3: ngrok para frontend
ngrok http 3000

# Terminal 4: ngrok para backend (requiere cuenta pagada para múltiples túneles)
ngrok http 5000
```

## Configuración después de iniciar ngrok

### 1. Actualizar configuración del frontend
Cuando ngrok te de una URL como `https://abc123.ngrok-free.app`, actualiza:

```javascript
// frontend/src/composables/useSocket.ts
const SOCKET_URL = 'https://abc123.ngrok-free.app'  // Para producción con ngrok
// O mantener localhost si usas el proxy de Vite
```

### 2. Compartir URL con compañeros
```
Tu app está disponible en: https://abc123.ngrok-free.app
```

## Solución de problemas

### Error: "Blocked request"
✅ Ya está solucionado en `vite.config.ts`

### Error de WebSocket
Si tienes problemas con WebSocket a través de ngrok:

1. Asegúrate de que el proxy esté configurado en `vite.config.ts`
2. O actualiza directamente la URL del socket en el frontend

### Límites de ngrok gratuito
- 1 túnel activo
- 40 conexiones por minuto
- URL cambia cada vez que reinicias

## Script de inicio rápido

Crea un archivo `start-with-ngrok.sh`:

```bash
#!/bin/bash

# Iniciar backend
cd backend && go run main.go &
BACKEND_PID=$!

# Iniciar frontend
cd ../frontend && npm run dev &
FRONTEND_PID=$!

# Esperar que los servicios inicien
sleep 5

# Iniciar ngrok
ngrok http 3000

# Limpiar al salir
trap "kill $BACKEND_PID $FRONTEND_PID" EXIT
```

Hazlo ejecutable:
```bash
chmod +x start-with-ngrok.sh
./start-with-ngrok.sh
```

## Tips para producción

1. **Para un demo más estable:**
   - Considera usar Cloudflare Tunnel (gratis)
   - O despliega en Railway/Render

2. **Seguridad:**
   - No compartas URLs de ngrok públicamente
   - Agrega autenticación básica si es necesario

3. **Rendimiento:**
   - ngrok puede ser lento para WebSockets
   - Considera alternativas como localtunnel o serveo

¡Ahora tu WhatsApp Clone es accesible desde internet! 🎉