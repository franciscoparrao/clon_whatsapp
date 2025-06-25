#!/bin/bash

# Script rápido para probar la conexión a PostgreSQL

echo "🔍 Probando conexión a PostgreSQL..."
echo ""

# Configuración
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres123}"

# Colores
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Método 1: Docker
echo "1. Intentando con Docker..."
if docker exec whatsapp_postgres pg_isready -U postgres &>/dev/null; then
    echo -e "   ${GREEN}✓ Conexión exitosa vía Docker${NC}"
    echo "   Comando: docker exec -it whatsapp_postgres psql -U postgres"
    exit 0
else
    echo "   ❌ No disponible vía Docker"
fi

# Método 2: Local con credenciales por defecto
echo ""
echo "2. Intentando conexión local (credenciales por defecto)..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c '\q' 2>/dev/null
if [ $? -eq 0 ]; then
    echo -e "   ${GREEN}✓ Conexión exitosa${NC}"
    echo "   Host: $DB_HOST"
    echo "   Puerto: $DB_PORT"
    echo "   Usuario: $DB_USER"
    exit 0
else
    echo "   ❌ Fallo con credenciales por defecto"
fi

# Método 3: Sin contraseña (peer authentication)
echo ""
echo "3. Intentando sin contraseña..."
psql -U $USER -d postgres -c '\q' 2>/dev/null
if [ $? -eq 0 ]; then
    echo -e "   ${GREEN}✓ Conexión exitosa con usuario: $USER${NC}"
    echo "   Nota: Deberás actualizar DB_USER=$USER en tu configuración"
    exit 0
else
    echo "   ❌ Fallo sin contraseña"
fi

# Si todo falla
echo ""
echo -e "${RED}❌ No se pudo conectar a PostgreSQL${NC}"
echo ""
echo "Opciones:"
echo "1. Iniciar con Docker:"
echo "   cd .. && ./start-db-docker.sh"
echo ""
echo "2. Iniciar PostgreSQL local:"
echo "   sudo systemctl start postgresql"
echo ""
echo "3. Verificar que PostgreSQL esté instalado:"
echo "   sudo apt install postgresql postgresql-client  # Ubuntu/Debian"
echo "   brew install postgresql                        # macOS"

exit 1