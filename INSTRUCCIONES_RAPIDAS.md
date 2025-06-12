# 🚀 Instrucciones Rápidas para Compañeros

## Opción 1: Usar la App en la Nube (Más Fácil)

Si uno de ustedes despliega la app en Railway o Render:

1. **Para conectarse:**
   - Abrir el navegador
   - Ir a: `https://tu-app.railway.app` (o la URL que les compartan)
   - Registrarse con email y contraseña
   - ¡Listo para chatear!

## Opción 2: Ejecutar Localmente con Docker

### Requisitos
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- Git

### Pasos
```bash
# 1. Clonar el proyecto
git clone [URL_DEL_REPOSITORIO]
cd clon_whatsapp

# 2. Configurar ambiente
cp .env.example .env

# 3. Ejecutar con Docker
docker-compose up -d

# 4. Abrir en navegador
# http://localhost:3000
```

## Opción 3: Ejecutar sin Docker (Desarrollo)

### Requisitos
- Go 1.22+
- Node.js 18+
- PostgreSQL

### Backend
```bash
cd backend
go mod download
go run main.go
```

### Frontend (nueva terminal)
```bash
cd frontend
npm install
npm run dev
```

## 🌐 Para Acceso Remoto (Red Local)

Si quieren que compañeros en la misma red WiFi se conecten:

1. **Obtener tu IP local:**
   ```bash
   # Windows
   ipconfig
   
   # Mac/Linux
   ifconfig
   ```

2. **Compartir esta URL con compañeros:**
   ```
   http://TU_IP_LOCAL:3000
   Ejemplo: http://192.168.1.100:3000
   ```

3. **Configurar CORS en el backend:**
   - Editar `backend/main.go`
   - Agregar la IP a `AllowOrigins`

## 📱 Funcionalidades

- ✅ Registro/Login
- ✅ Lista de chats
- ✅ Mensajes en tiempo real
- ✅ Indicadores de estado (enviado/leído)
- ✅ Crear chats grupales
- ✅ Enviar mensajes

## 🔧 Solución de Problemas

### Puerto ocupado
```bash
# Cambiar puerto en .env
FRONTEND_PORT=3001
BACKEND_PORT=5001
```

### Error de conexión WebSocket
- Verificar que el backend esté corriendo
- Revisar la URL del WebSocket en `frontend/src/composables/useSocket.ts`

### Base de datos
```bash
# Crear base de datos
docker-compose exec postgres psql -U postgres
CREATE DATABASE whatsapp_clone;
```

## 💡 Tips para el Proyecto

1. **Para Plantillas de Mensajes:**
   - Crear endpoint `/api/templates`
   - Usar variables como `{{nombre}}`, `{{fecha}}`
   - Guardar plantillas en la BD

2. **Para Integración con WhatsApp Business:**
   - Usar webhooks para recibir mensajes
   - Implementar cola de mensajes
   - Agregar límites de rate

3. **Mejoras sugeridas:**
   - Agregar autenticación OAuth
   - Implementar cifrado E2E
   - Añadir notificaciones push
   - Soporte para archivos multimedia

## 📞 Soporte

Si tienen problemas:
1. Revisar logs: `docker-compose logs -f`
2. Verificar que todos los servicios estén up: `docker-compose ps`
3. Reiniciar servicios: `docker-compose restart`

¡Éxito con el proyecto! 🎉