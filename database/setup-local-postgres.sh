#!/bin/bash

# Script para configurar PostgreSQL local

echo "=== Configuración de PostgreSQL Local ==="
echo ""

# Colores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# Verificar si PostgreSQL está instalado
if ! command -v psql &> /dev/null; then
    echo -e "${RED}PostgreSQL no está instalado${NC}"
    echo "Instálalo con:"
    echo "  sudo apt update && sudo apt install postgresql postgresql-client"
    exit 1
fi

echo -e "${GREEN}✓ PostgreSQL está instalado${NC}"

# Verificar si PostgreSQL está ejecutándose
echo -n "Verificando servicio PostgreSQL... "
if systemctl is-active --quiet postgresql; then
    echo -e "${GREEN}Activo${NC}"
else
    echo -e "${YELLOW}Inactivo${NC}"
    echo "Iniciando PostgreSQL..."
    sudo systemctl start postgresql
    sleep 2
fi

echo ""
echo "Configurando PostgreSQL para el proyecto..."
echo ""

# Opción 1: Crear usuario y base de datos con sudo
echo "Método 1: Usando usuario postgres del sistema"
echo "========================================"

# Crear base de datos y usuario
sudo -u postgres psql <<EOF
-- Crear usuario si no existe
DO \$\$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_user WHERE usename = 'postgres') THEN
        CREATE USER postgres WITH PASSWORD 'postgres123';
    END IF;
END\$\$;

-- Dar permisos de superusuario
ALTER USER postgres WITH SUPERUSER;

-- Crear base de datos si no existe
SELECT 'CREATE DATABASE whatsapp_clone'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'whatsapp_clone')\gexec

-- Dar permisos sobre la base de datos
GRANT ALL PRIVILEGES ON DATABASE whatsapp_clone TO postgres;

-- Mostrar configuración
\echo ''
\echo 'Configuración creada:'
\echo '===================='
\echo 'Usuario: postgres'
\echo 'Contraseña: postgres123'
\echo 'Base de datos: whatsapp_clone'
\echo ''
EOF

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Base de datos y usuario configurados${NC}"
else
    echo -e "${YELLOW}Advertencia: Puede que ya existan${NC}"
fi

# Verificar autenticación en pg_hba.conf
echo ""
echo "Verificando configuración de autenticación..."

PG_VERSION=$(sudo -u postgres psql -t -c "SELECT version();" | grep -oP 'PostgreSQL \K[0-9]+')
HBA_FILE="/etc/postgresql/$PG_VERSION/main/pg_hba.conf"

if [ -f "$HBA_FILE" ]; then
    echo "Archivo de configuración: $HBA_FILE"
    
    # Verificar si ya está configurado para md5/password
    if grep -q "local.*all.*all.*md5\|local.*all.*all.*password" "$HBA_FILE"; then
        echo -e "${GREEN}✓ Autenticación por contraseña ya configurada${NC}"
    else
        echo -e "${YELLOW}Configurando autenticación por contraseña...${NC}"
        
        # Hacer backup
        sudo cp "$HBA_FILE" "$HBA_FILE.backup.$(date +%Y%m%d_%H%M%S)"
        
        # Agregar línea para autenticación por contraseña
        sudo sed -i '/^local.*all.*postgres.*peer/a local   all             all                                     md5' "$HBA_FILE"
        
        echo "Reiniciando PostgreSQL..."
        sudo systemctl restart postgresql
        sleep 2
    fi
fi

# Probar conexión
echo ""
echo "Probando conexión..."
PGPASSWORD=postgres123 psql -h localhost -U postgres -d whatsapp_clone -c '\dt' &>/dev/null

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Conexión exitosa${NC}"
    echo ""
    echo "📊 Base de datos lista para usar:"
    echo "   Host: localhost"
    echo "   Puerto: 5432"
    echo "   Usuario: postgres"
    echo "   Contraseña: postgres123"
    echo "   Base de datos: whatsapp_clone"
    echo ""
    echo "🚀 Siguiente paso:"
    echo "   cd /home/franciscoparrao/proyectos/clon_whatsapp/database"
    echo "   ./migrate.sh"
else
    echo -e "${YELLOW}⚠️  La conexión requiere configuración adicional${NC}"
    echo ""
    echo "Opción alternativa: Crear usuario con tu nombre de usuario del sistema"
    echo ""
    
    # Crear usuario con el nombre del usuario actual
    sudo -u postgres createuser -s $USER 2>/dev/null
    
    # Crear base de datos con el usuario actual
    createdb whatsapp_clone 2>/dev/null
    
    echo "Intenta ejecutar migrate.sh con:"
    echo "   export DB_USER=$USER"
    echo "   export DB_PASSWORD=''"
    echo "   ./migrate.sh"
fi

echo ""
echo -e "${BLUE}Nota: Si prefieres usar Docker:${NC}"
echo "  cd .. && ./start-db-docker.sh (selecciona opción 3 para usar puerto 5433)"