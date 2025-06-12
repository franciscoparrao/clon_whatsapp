# 🧪 Guía de Prueba - WhatsApp Clone

## Verificar que todo funcione

### 1. Verificar Backend
```bash
# Probar health check
curl http://localhost:5000/api/health

# Debería responder:
# {"service":"whatsapp-clone-backend","status":"healthy","timestamp":"..."}
```

### 2. Probar desde ngrok

Una vez que tengas tu URL de ngrok (ej: `https://abc123.ngrok-free.app`):

#### a) Registro de usuario
1. Abre la URL de ngrok en el navegador
2. Te redirigirá a `/login`
3. Click en "¿No tienes cuenta? Regístrate"
4. Completa el formulario:
   - Nombre: Juan
   - Email: juan@test.com
   - Contraseña: 123456
5. Click en "Registrarse"

#### b) Login
1. Usa las credenciales que creaste
2. Click en "Iniciar Sesión"

#### c) Crear un chat
1. Click en el botón "+" (arriba a la derecha en la lista de chats)
2. Selecciona un usuario de la lista
3. Click en "Crear Chat"

#### d) Enviar mensaje
1. Selecciona un chat de la lista
2. Escribe un mensaje en el campo de texto
3. Presiona Enter o click en el botón enviar

### 3. Probar con múltiples usuarios

Para probar la funcionalidad en tiempo real:

1. **Usuario 1**: Abre en Chrome normal
2. **Usuario 2**: Abre en Chrome incógnito
3. Registra dos usuarios diferentes
4. Crea un chat entre ellos
5. Envía mensajes y verifica que aparezcan en tiempo real

### 4. Verificar en la consola del navegador

Abre las DevTools (F12) y verifica:

```javascript
// En la consola deberías ver:
// "Connected to WebSocket server"
// "Authenticated successfully"
```

### 5. Solución de problemas comunes

#### Error: "Network Error" al hacer login
- Verifica que el backend esté corriendo: `ps aux | grep "go run"`
- Revisa los logs del backend

#### WebSocket no conecta
- Verifica en Network tab que `/ws` esté intentando conectar
- Debe mostrar status 101 (Switching Protocols)

#### No se ven los mensajes en tiempo real
- Verifica que ambos usuarios estén en el mismo chat
- Revisa la consola por errores de WebSocket

### 6. URLs de prueba directa

Con tu URL de ngrok:
- Login: `https://abc123.ngrok-free.app/login`
- Home: `https://abc123.ngrok-free.app/` (requiere login)
- API Health: `https://abc123.ngrok-free.app/api/health`

### 7. Datos de prueba

Para pruebas rápidas, puedes usar:

**Usuario 1:**
- Email: test1@example.com
- Password: password123

**Usuario 2:**
- Email: test2@example.com  
- Password: password123

## 🎯 Checklist de funcionalidades

- [ ] Registro de usuario
- [ ] Login
- [ ] Logout (botón arriba a la derecha)
- [ ] Ver lista de chats
- [ ] Crear nuevo chat
- [ ] Seleccionar chat
- [ ] Enviar mensaje
- [ ] Recibir mensaje en tiempo real
- [ ] Ver indicadores de estado (✓ ✓✓)

¡Si todo funciona, tu WhatsApp Clone está listo! 🎉