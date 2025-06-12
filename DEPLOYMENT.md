# WhatsApp Clone - Guía de Implementación y Configuración

Esta guía te ayudará a desplegar el clon de WhatsApp para que puedas usarlo con tus compañeros de clase.

## 📋 Requisitos Previos

- Docker y Docker Compose instalados
- PostgreSQL (incluido en Docker Compose)
- Un dominio o IP pública (para acceso remoto)
- Certificado SSL (recomendado para producción)

## 🚀 Despliegue Local (Para Pruebas)

### 1. Clonar el Repositorio

```bash
git clone <tu-repositorio>
cd clon_whatsapp
```

### 2. Configurar Variables de Entorno

```bash
cp .env.example .env
```

Edita el archivo `.env` con tus configuraciones:

```bash
nano .env
```

### 3. Iniciar con Docker Compose

```bash
docker-compose up -d
```

La aplicación estará disponible en:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

### 4. Verificar el Estado

```bash
docker-compose ps
docker-compose logs -f
```

## 🌐 Despliegue en Servidor (VPS/Cloud)

### Opción A: Despliegue en VPS (DigitalOcean, Linode, etc.)

#### 1. Preparar el Servidor

```bash
# Actualizar el sistema
sudo apt update && sudo apt upgrade -y

# Instalar Docker
curl -fsSL https://get.docker.com | sh

# Instalar Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Agregar tu usuario al grupo docker
sudo usermod -aG docker $USER
```

#### 2. Configurar Firewall

```bash
# Permitir SSH
sudo ufw allow 22

# Permitir HTTP y HTTPS
sudo ufw allow 80
sudo ufw allow 443

# Permitir WebSocket (si usas puerto diferente)
sudo ufw allow 8080

# Habilitar firewall
sudo ufw enable
```

#### 3. Clonar y Configurar

```bash
# Clonar repositorio
git clone <tu-repositorio>
cd clon_whatsapp

# Configurar variables de entorno
cp .env.example .env
nano .env
```

#### 4. Configurar Nginx como Proxy Reverso

```bash
# Instalar Nginx
sudo apt install nginx -y

# Crear configuración
sudo nano /etc/nginx/sites-available/whatsapp-clone
```

Contenido del archivo:

```nginx
server {
    listen 80;
    server_name tu-dominio.com;

    # Frontend
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # Backend API
    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # WebSocket
    location /ws {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

```bash
# Habilitar sitio
sudo ln -s /etc/nginx/sites-available/whatsapp-clone /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

#### 5. Configurar SSL con Let's Encrypt

```bash
# Instalar Certbot
sudo apt install certbot python3-certbot-nginx -y

# Obtener certificado SSL
sudo certbot --nginx -d tu-dominio.com
```

### Opción B: Despliegue en Railway/Render

#### Railway

1. Crea una cuenta en [Railway](https://railway.app)
2. Conecta tu repositorio de GitHub
3. Añade las variables de entorno en el dashboard
4. Railway detectará automáticamente los Dockerfiles

#### Render

1. Crea una cuenta en [Render](https://render.com)
2. Crea un nuevo Web Service
3. Conecta tu repositorio
4. Configura las variables de entorno
5. Render construirá y desplegará automáticamente

## 👥 Configuración para Múltiples Usuarios

### 1. Base de Datos

La aplicación usa PostgreSQL. Para producción:

```sql
-- Crear base de datos
CREATE DATABASE whatsapp_clone;

-- Crear usuario
CREATE USER whatsapp_user WITH ENCRYPTED PASSWORD 'tu_password_segura';

-- Dar permisos
GRANT ALL PRIVILEGES ON DATABASE whatsapp_clone TO whatsapp_user;
```

### 2. Escalabilidad

Para soportar múltiples usuarios simultáneos:

```yaml
# docker-compose.yml - Aumentar réplicas
services:
  backend:
    deploy:
      replicas: 3
```

### 3. Redis para Sesiones (Opcional)

Si necesitas manejar muchos usuarios, añade Redis:

```yaml
# docker-compose.yml
redis:
  image: redis:alpine
  ports:
    - "6379:6379"
```

## 🔒 Seguridad

### 1. Variables de Entorno Seguras

```env
# .env
JWT_SECRET=genera_una_clave_segura_aqui
DB_PASSWORD=usa_una_contraseña_fuerte
CORS_ORIGINS=https://tu-dominio.com
```

### 2. Limitar Acceso

```nginx
# Limitar requests
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

location /api {
    limit_req zone=api_limit burst=20 nodelay;
    # ... resto de configuración
}
```

### 3. Monitoreo

```bash
# Instalar herramientas de monitoreo
docker-compose logs -f backend
docker stats
```

## 📱 Instrucciones para Compañeros

### Para Conectarse:

1. **Acceder a la Aplicación**
   ```
   https://tu-dominio.com
   ```

2. **Crear una Cuenta**
   - Click en "Registrarse"
   - Ingresar nombre de usuario y contraseña
   - La contraseña debe tener al menos 6 caracteres

3. **Iniciar Sesión**
   - Usar las credenciales creadas
   - El sistema recordará la sesión

4. **Usar el Chat**
   - Ver lista de usuarios conectados
   - Click en un usuario para iniciar chat
   - Los mensajes se envían en tiempo real

### Requisitos del Cliente:
- Navegador moderno (Chrome, Firefox, Safari, Edge)
- Conexión a internet estable
- JavaScript habilitado

## 🛠️ Solución de Problemas

### Error de Conexión WebSocket

```bash
# Verificar que el backend esté corriendo
docker-compose ps
docker-compose logs backend
```

### Base de Datos no Conecta

```bash
# Verificar PostgreSQL
docker-compose exec db psql -U postgres
# Verificar tablas
\dt
```

### Puerto en Uso

```bash
# Cambiar puertos en docker-compose.yml
ports:
  - "3001:3000"  # Frontend
  - "8081:8080"  # Backend
```

## 📊 Monitoreo y Logs

### Ver Logs en Tiempo Real

```bash
# Todos los servicios
docker-compose logs -f

# Solo backend
docker-compose logs -f backend

# Solo frontend
docker-compose logs -f frontend
```

### Estadísticas de Uso

```bash
# Ver uso de recursos
docker stats

# Ver procesos
docker-compose top
```

## 🔄 Actualizaciones

Para actualizar la aplicación:

```bash
# Detener servicios
docker-compose down

# Actualizar código
git pull origin main

# Reconstruir imágenes
docker-compose build

# Iniciar servicios
docker-compose up -d
```

## 🆘 Soporte

Si tienes problemas:

1. Revisa los logs: `docker-compose logs`
2. Verifica las variables de entorno
3. Asegúrate de que los puertos estén abiertos
4. Contacta al administrador con el error específico

## 📝 Notas Adicionales

- La aplicación guarda todos los mensajes en la base de datos
- Los usuarios pueden ver el historial de conversaciones
- El sistema soporta múltiples conversaciones simultáneas
- Los mensajes se entregan en tiempo real vía WebSocket