# 🔍 Debug de ngrok - Error 403

## Verificaciones rápidas

### 1. Probar el health check directamente
```bash
# Desde tu máquina local
curl http://localhost:5000/api/health

# Desde ngrok (reemplaza con tu URL)
curl https://tu-url.ngrok-free.app/api/health
```

### 2. Verificar en el navegador
1. Abre las DevTools (F12)
2. Ve a la pestaña Network
3. Intenta registrarte
4. Busca la petición que falla (status 403)
5. Click en ella y revisa:
   - Headers enviados
   - Response headers
   - Response body

### 3. Posibles causas del error 403

#### A. ngrok requiere header especial
Algunas versiones de ngrok requieren un header especial. Prueba agregar esto al frontend:

```javascript
// En frontend/src/services/api.ts
headers['ngrok-skip-browser-warning'] = 'true'
```

#### B. El backend está rechazando la petición
Revisa los logs del backend cuando hagas la petición.

### 4. Solución alternativa - Sin autenticación para pruebas

Temporalmente, puedes deshabilitar la autenticación para las rutas de registro/login:

```go
// En backend/main.go
api.POST("/auth/register", handlers.Register) // Sin middleware de auth
api.POST("/auth/login", handlers.Login)       // Sin middleware de auth
```

### 5. Usar ngrok con configuración personalizada

Crea un archivo `ngrok.yml`:

```yaml
version: "2"
tunnels:
  whatsapp:
    proto: http
    addr: 3004
    host_header: "localhost:3004"
```

Ejecuta:
```bash
ngrok start --config=ngrok.yml whatsapp
```

### 6. Verificar que el frontend use las rutas correctas

En la consola del navegador, verifica:
```javascript
// Debería mostrar las peticiones a /api/...
// NO a http://localhost:5000/api/...
```

### 7. Logs detallados

En el backend, los logs deberían mostrar algo como:
```
2025/06/11 18:57:43 Request: POST /api/auth/register from https://tu-url.ngrok-free.app
```

Si ves el origen correcto pero sigue fallando, el problema puede ser el middleware de autenticación.

## Solución rápida

Si nada funciona, prueba esto:

1. **En el frontend**, agrega el header de ngrok:
```javascript
// frontend/src/services/api.ts
headers['ngrok-skip-browser-warning'] = 'true'
```

2. **Reinicia el frontend**

3. **Limpia caché del navegador** (Ctrl+Shift+R)

4. **Intenta de nuevo**