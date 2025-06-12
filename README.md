# WhatsApp Clone

Un clon de WhatsApp desarrollado con Go (backend) y Vue.js (frontend) que incluye chatbots integrados.

## 🚀 Características

- ✅ Mensajería en tiempo real con WebSocket
- ✅ Autenticación JWT
- ✅ 3 Chatbots integrados (Soporte, Ventas, Plantillas)
- ✅ Interfaz similar a WhatsApp
- ✅ Indicadores de estado de mensaje
- ✅ Sistema de plantillas de mensajes
- ✅ Compatible con ngrok para acceso remoto

## 📁 Estructura del Proyecto

```
clon_whatsapp/
├── backend/          # API REST y WebSocket server en Go
│   ├── config/       # Configuración
│   ├── handlers/     # Controladores HTTP y WebSocket
│   ├── middleware/   # Autenticación y logging
│   ├── models/       # Modelos de datos
│   ├── services/     # Lógica de negocio y chatbots
│   └── main.go       # Punto de entrada
└── frontend/         # Aplicación Vue.js
    ├── src/
    │   ├── components/   # Componentes Vue
    │   ├── composables/  # Hooks reutilizables
    │   ├── router/       # Configuración de rutas
    │   ├── services/     # Servicios API
    │   ├── types/        # Tipos TypeScript
    │   └── views/        # Vistas principales
    └── package.json
```

## 🛠️ Instalación

### Prerrequisitos
- Go 1.22+
- Node.js 18+
- Git

### Backend

```bash
cd backend
go mod download
go run main.go
```

El servidor se ejecutará en `http://localhost:5000`

### Frontend

```bash
cd frontend
npm install
npm run dev
```

La aplicación se ejecutará en `http://localhost:3000`

## 🤖 Chatbots Integrados

El sistema incluye 3 chatbots que se crean automáticamente al registrarse:

### 1. Soporte Técnico
Palabras clave: `hola`, `ayuda`, `precio`, `horario`, `contacto`, `problema`, `pedido`, `envio`, `devolucion`

### 2. Ventas Bot
Palabras clave: `hola`, `oferta`, `descuento`, `catalogo`, `pago`, `cuotas`

### 3. Template Bot
Palabras clave: `hola`, `ayuda`, `nueva`, `listar`, `variables`, `ejemplo`

## 🌐 Uso con ngrok

Para compartir la aplicación con otros:

```bash
# Terminal 1: Backend
cd backend && go run main.go

# Terminal 2: Frontend
cd frontend && npm run dev

# Terminal 3: ngrok
ngrok http 3000
```

## 📝 Variables de Entorno

Crea un archivo `.env` en el backend:

```env
PORT=5000
JWT_SECRET=your-secret-key-change-this
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=whatsapp_clone
```

## 🔧 API Endpoints

### Autenticación
- `POST /api/auth/register` - Registro de usuario
- `POST /api/auth/login` - Login

### Chats (requieren autenticación)
- `GET /api/chats` - Obtener chats del usuario
- `POST /api/chats` - Crear nuevo chat
- `GET /api/chats/:id/messages` - Obtener mensajes
- `POST /api/chats/:id/messages` - Enviar mensaje

### WebSocket
- `GET /ws` - Conexión WebSocket para mensajes en tiempo real

## 👥 Autores

- Francisco Parra - [franciscoparrao](https://github.com/franciscoparrao)

## 📄 Licencia

Este proyecto está bajo la Licencia MIT.