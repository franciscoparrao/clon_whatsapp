# Configuración de Base de Datos PostgreSQL

Este directorio contiene los scripts necesarios para configurar la base de datos PostgreSQL del proyecto Clon de WhatsApp.

## Opciones de Configuración

### Opción 1: Usar Docker (Recomendado) 🐳

Esta es la opción más sencilla y no requiere instalar PostgreSQL localmente.

```bash
# Desde el directorio database/
./setup-db.sh

# Selecciona opción 1
```

Esto iniciará PostgreSQL en un contenedor Docker con las siguientes credenciales:
- **Host:** localhost
- **Puerto:** 5432
- **Usuario:** postgres
- **Contraseña:** postgres123
- **Base de datos:** whatsapp_clone

### Opción 2: Usar PostgreSQL Local

Si ya tienes PostgreSQL instalado localmente:

```bash
# Desde el directorio database/
./setup-db.sh

# Selecciona opción 2
# Ingresa tus credenciales de PostgreSQL
```

### Opción 3: Configuración Manual

Si prefieres configurar manualmente:

1. Conéctate a PostgreSQL:
   ```bash
   psql -U tu_usuario
   ```

2. Crea la base de datos:
   ```sql
   CREATE DATABASE whatsapp_clone;
   \c whatsapp_clone
   ```

3. Ejecuta el schema:
   ```bash
   psql -U tu_usuario -d whatsapp_clone -f schema.sql
   ```

## Verificar la Instalación

### Con Docker:
```bash
docker exec -it whatsapp_postgres psql -U postgres -d whatsapp_clone -c '\dt'
```

### Con PostgreSQL local:
```bash
psql -U postgres -d whatsapp_clone -c '\dt'
```

Deberías ver todas las tablas creadas:
- users
- privacy_settings
- contacts
- chats
- chat_participants
- messages
- message_status
- gig_workers
- worker_availability
- tasks
- task_assignments
- attendance_history
- notifications

## Configuración del Backend

Crea un archivo `.env` en el directorio `backend/` con las credenciales correctas:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres123
DB_NAME=whatsapp_clone
DB_SSLMODE=disable
```

## Solución de Problemas

### Error: "No se pudo conectar a PostgreSQL"

1. **Si usas Docker:**
   - Verifica que Docker esté ejecutándose: `docker ps`
   - Inicia el contenedor: `docker-compose -f ../docker-compose.db.yml up -d`

2. **Si usas PostgreSQL local:**
   - Verifica que PostgreSQL esté ejecutándose:
     ```bash
     # En Ubuntu/Debian
     sudo systemctl status postgresql
     sudo systemctl start postgresql
     
     # En macOS con Homebrew
     brew services list
     brew services start postgresql
     ```

### Error: "FATAL: password authentication failed"

- Verifica que las credenciales en tu archivo `.env` coincidan con las de tu instalación de PostgreSQL
- Si usas Docker, las credenciales por defecto son: usuario `postgres`, contraseña `postgres123`

### Error: "database whatsapp_clone does not exist"

Ejecuta:
```bash
createdb -U postgres whatsapp_clone
```

## Scripts Disponibles

- `schema.sql`: Contiene toda la estructura de tablas e índices
- `migrate.sh`: Script básico de migración (requiere PostgreSQL local)
- `setup-db.sh`: Script interactivo para configurar la base de datos
- `../docker-compose.db.yml`: Configuración de Docker para PostgreSQL