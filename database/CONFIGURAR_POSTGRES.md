# Configurar PostgreSQL para el proyecto

Para que funcione la conexión con usuario `postgres` y contraseña `postgres`, ejecuta estos comandos:

## 1. Configurar la contraseña del usuario postgres:

```bash
sudo -u postgres psql -f set-postgres-password.sql
```

## 2. Verificar/actualizar la autenticación en PostgreSQL:

Busca tu versión de PostgreSQL:
```bash
ls /etc/postgresql/
```

Edita el archivo pg_hba.conf (reemplaza XX con tu versión):
```bash
sudo nano /etc/postgresql/XX/main/pg_hba.conf
```

Busca estas líneas:
```
local   all             postgres                                peer
local   all             all                                     peer
```

Y cámbialas a:
```
local   all             postgres                                md5
local   all             all                                     md5
```

También asegúrate de que esta línea esté presente para conexiones por TCP/IP:
```
host    all             all             127.0.0.1/32            md5
```

## 3. Reiniciar PostgreSQL:

```bash
sudo systemctl restart postgresql
```

## 4. Probar la conexión:

```bash
PGPASSWORD=postgres psql -h localhost -U postgres -d whatsapp_clone
```

Si funciona, deberías ver el prompt de PostgreSQL. Escribe `\q` para salir.

## 5. Ejecutar el backend:

```bash
cd ../backend
go run main.go
```

---

## Alternativa: Usar Docker

Si prefieres evitar configurar PostgreSQL local:

```bash
# Detener PostgreSQL local
sudo systemctl stop postgresql

# Iniciar con Docker
cd ..
./start-db-docker.sh
```