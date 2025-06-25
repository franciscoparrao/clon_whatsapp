# 📱 Clon de WhatsApp con Gestión de Gig Workers

Este proyecto es un clon de WhatsApp desarrollado con Go (backend) y Vue 3 (frontend), que incluye funcionalidades adicionales para la gestión de trabajadores gig (gig workers) según los requerimientos del informe académico.

## 🏗️ Arquitectura del Proyecto

```
clon_whatsapp/
├── backend/          # API REST en Go con Gin
├── frontend/         # SPA en Vue 3 + TypeScript + Vite
├── database/         # Scripts y esquemas de PostgreSQL
├── docker-compose.yml        # Configuración completa con Docker
├── docker-compose.db.yml     # Solo PostgreSQL con Docker
└── start-db-docker.sh        # Script para iniciar PostgreSQL
```

## 🚀 Guía de Instalación Rápida

### Opción A: Instalación con Docker (Recomendado)

```bash
# 1. Clonar el repositorio
git clone <url-del-repositorio>
cd clon_whatsapp

# 2. Iniciar PostgreSQL con Docker
./start-db-docker.sh

# 3. Cargar el esquema de base de datos
cd database
./load-schema-docker.sh

# 4. Iniciar el backend
cd ../backend
go run main.go

# 5. En otra terminal, iniciar el frontend
cd ../frontend
npm install
npm run dev
```

### Opción B: Instalación con PostgreSQL Local

```bash
# 1. Clonar el repositorio
git clone <url-del-repositorio>
cd clon_whatsapp

# 2. Configurar PostgreSQL local
cd database
./quick-setup.sh

# 3. Seguir los pasos 4 y 5 de la Opción A
```

## 📋 Requisitos Previos

### Software Necesario

- **Node.js** 18+ y npm
- **Go** 1.22+
- **PostgreSQL** 15+ o **Docker** con Docker Compose
- **Git**

### Instalación de Requisitos

#### En Ubuntu/Debian:

```bash
# Actualizar sistema
sudo apt update && sudo apt upgrade -y

# Instalar Node.js 18+
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs

# Instalar Go
sudo snap install go --classic
# O descarga desde https://go.dev/dl/

# Opción 1: Instalar PostgreSQL
sudo apt install -y postgresql postgresql-client

# Opción 2: Instalar Docker
curl -fsSL https://get.docker.com | sudo bash
sudo usermod -aG docker $USER
# Cerrar sesión y volver a entrar

# Instalar Git
sudo apt install -y git
```

#### En macOS:

```bash
# Instalar Homebrew si no lo tienes
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Instalar requisitos
brew install node
brew install go
brew install postgresql  # O docker
brew install git
```

#### En Windows:

1. Instalar [Node.js](https://nodejs.org/)
2. Instalar [Go](https://go.dev/dl/)
3. Instalar [Docker Desktop](https://www.docker.com/products/docker-desktop/) o [PostgreSQL](https://www.postgresql.org/download/windows/)
4. Instalar [Git](https://git-scm.com/download/win)

## 🔧 Configuración Detallada

### 1. Base de Datos PostgreSQL

#### Opción 1: Usar Docker (Recomendado)

```bash
# Iniciar PostgreSQL con Docker
./start-db-docker.sh

# Si el puerto 5432 está ocupado, el script te dará opciones:
# 1) Detener PostgreSQL local y usar Docker
# 2) Usar PostgreSQL local existente
# 3) Usar Docker en otro puerto (5433)
```

#### Opción 2: Configurar PostgreSQL Local

```bash
cd database

# Método rápido (intenta con tu usuario del sistema)
./quick-setup.sh

# Método completo (configura usuario postgres con contraseña)
sudo ./setup-local-postgres.sh
```

#### Cargar el Esquema de Base de Datos

```bash
cd database

# Si usas Docker
./load-schema-docker.sh

# Si usas PostgreSQL local
./migrate.sh
```

#### Verificar la Base de Datos

```bash
# Verificar estado y tablas
./check-db.sh

# Conectarse manualmente (Docker)
sudo docker exec -it whatsapp_postgres psql -U postgres -d whatsapp_clone

# Conectarse manualmente (Local)
psql -U postgres -d whatsapp_clone
```

### 2. Configuración del Backend

#### Variables de Entorno

Crea un archivo `.env` en el directorio `backend/` (ya existe uno de ejemplo):

```env
# Configuración de Base de Datos
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres123
DB_NAME=whatsapp_clone
DB_SSLMODE=disable

# Configuración JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRY=24h

# Puerto del servidor
PORT=8080

# Otros
UPLOAD_DIR=./uploads
MAX_MESSAGE_SIZE=10485760
```

#### Instalar Dependencias e Iniciar

```bash
cd backend

# Instalar dependencias
go mod download

# Iniciar el servidor
go run main.go

# El servidor estará disponible en http://localhost:8080
```

### 3. Configuración del Frontend

#### Variables de Entorno

Crea un archivo `.env` en el directorio `frontend/`:

```env
VITE_API_URL=http://localhost:8080/api
VITE_WS_URL=ws://localhost:8080
```

#### Instalar Dependencias e Iniciar

```bash
cd frontend

# Instalar dependencias
npm install

# Iniciar en modo desarrollo
npm run dev

# La aplicación estará disponible en http://localhost:5173
```

## 📱 Funcionalidades Principales

### WhatsApp Clone
- ✅ Registro e inicio de sesión
- ✅ Chat en tiempo real con WebSocket
- ✅ Lista de conversaciones
- ✅ Envío de mensajes
- ✅ Estado en línea/última vez
- ✅ Perfil de usuario
- ✅ 3 Chatbots integrados (Soporte, Ventas, Plantillas)

### Gestión de Gig Workers
- ✅ Registro como trabajador
- ✅ Gestión de disponibilidad
- ✅ Asignación de tareas
- ✅ Seguimiento de asignaciones
- ✅ Historial de asistencia
- ✅ Sistema de notificaciones
- ✅ Estadísticas y métricas

## 🗂️ Estructura de la Base de Datos

### Tablas Principales

1. **users** - Usuarios del sistema
2. **chats** - Conversaciones
3. **messages** - Mensajes
4. **gig_workers** - Trabajadores
5. **tasks** - Tareas/turnos
6. **task_assignments** - Asignaciones
7. **worker_availability** - Disponibilidad
8. **attendance_history** - Historial de asistencia
9. **notifications** - Notificaciones

## 🛠️ Solución de Problemas

### Error: "No se pudo conectar a PostgreSQL"

```bash
# Verificar si PostgreSQL está ejecutándose
sudo systemctl status postgresql  # Linux
brew services list                 # macOS

# Iniciar PostgreSQL
sudo systemctl start postgresql    # Linux
brew services start postgresql     # macOS

# O usar Docker
./start-db-docker.sh
```

### Error: "password authentication failed for user postgres"

```bash
# Configurar contraseña del usuario postgres
sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'postgres123';"

# O usar el script de configuración
cd database
sudo ./fix-postgres-auth.sh
```

### Error: "address already in use :5432"

```bash
# Opción 1: Detener PostgreSQL local
sudo systemctl stop postgresql

# Opción 2: Usar Docker en otro puerto
# El script start-db-docker.sh te dará esta opción
```

### Error en el Frontend: "Cannot find module axios"

```bash
cd frontend
npm install axios
```

## 🐳 Uso con Docker Compose Completo

Para ejecutar todo el proyecto con Docker:

```bash
# Construir e iniciar todos los servicios
docker-compose up -d

# Ver logs
docker-compose logs -f

# Detener todos los servicios
docker-compose down
```

## 📚 Documentación de API

Ver [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) para la documentación completa de los endpoints.

### Endpoints Principales

- `POST /api/auth/register` - Registro de usuario
- `POST /api/auth/login` - Inicio de sesión
- `GET /api/chats` - Obtener conversaciones
- `POST /api/chats/:id/messages` - Enviar mensaje
- `POST /api/workers/register` - Registrarse como trabajador
- `GET /api/tasks` - Ver tareas disponibles
- `POST /api/workers/availability` - Actualizar disponibilidad

## 🧪 Pruebas

### Probar la API con curl

```bash
# Registrar usuario
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123","email":"test@example.com"}'

# Iniciar sesión
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'
```

## 🚀 Despliegue en Producción

### Consideraciones de Seguridad

1. **Cambiar todas las contraseñas y secrets por defecto**
2. **Configurar HTTPS con certificados SSL**
3. **Habilitar CORS solo para dominios específicos**
4. **Configurar un firewall**
5. **Usar variables de entorno seguras**

### Ejemplo de Configuración para Producción

```env
# .env.production
DB_HOST=tu-servidor-db.com
DB_PASSWORD=contraseña-segura-generada
JWT_SECRET=secret-muy-largo-y-aleatorio
GIN_MODE=release
```

## 📄 Licencia

Este proyecto es parte de un trabajo académico para la Universidad de Santiago de Chile.

## 👥 Autores

- Estudiantes del curso de Innovación y Emprendimiento
- Profesor: Francisco Parra

## 🤝 Contribuciones

Para contribuir al proyecto:

1. Fork el repositorio
2. Crea una rama (`git checkout -b feature/nueva-funcionalidad`)
3. Commit tus cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Crea un Pull Request

## 📞 Soporte

Si encuentras problemas:

1. Revisa la sección de [Solución de Problemas](#-solución-de-problemas)
2. Busca en los [Issues](https://github.com/tu-usuario/clon_whatsapp/issues) existentes
3. Crea un nuevo Issue con detalles del problema

---

**Nota**: Este proyecto requiere PostgreSQL 15+ y fue desarrollado con Go 1.22 y Node.js 18+. Asegúrate de tener las versiones correctas instaladas.