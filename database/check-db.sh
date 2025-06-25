#!/bin/bash

# Script para verificar el estado de la base de datos

echo "=== Verificación de Base de Datos ==="
echo ""

# Configuración predeterminada
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres123}"
DB_NAME="${DB_NAME:-whatsapp_clone}"

# Colores
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Verificar Docker
echo -n "🐳 Verificando Docker... "
if command -v docker &> /dev/null; then
    echo -e "${GREEN}Instalado${NC}"
    
    echo -n "   Contenedor PostgreSQL... "
    if docker ps --format "table {{.Names}}" | grep -q "whatsapp_postgres"; then
        echo -e "${GREEN}Ejecutándose${NC}"
        USE_DOCKER=true
        PSQL_CMD="docker exec -i whatsapp_postgres psql -U $DB_USER"
    else
        echo -e "${YELLOW}No encontrado${NC}"
        USE_DOCKER=false
    fi
else
    echo -e "${YELLOW}No instalado${NC}"
    USE_DOCKER=false
fi

# Si no hay Docker, verificar PostgreSQL local
if [ "$USE_DOCKER" = false ]; then
    echo -n "🐘 Verificando PostgreSQL local... "
    if command -v psql &> /dev/null; then
        echo -e "${GREEN}Instalado${NC}"
        
        echo -n "   Conexión... "
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c '\q' 2>/dev/null
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}OK${NC}"
            export PGPASSWORD=$DB_PASSWORD
            PSQL_CMD="psql -h $DB_HOST -p $DB_PORT -U $DB_USER"
        else
            echo -e "${RED}Fallo${NC}"
            echo ""
            echo "💡 Sugerencias:"
            echo "   1. Inicia PostgreSQL con Docker: cd .. && ./start-db-docker.sh"
            echo "   2. O inicia PostgreSQL local: sudo systemctl start postgresql"
            exit 1
        fi
    else
        echo -e "${RED}No instalado${NC}"
        exit 1
    fi
fi

# Verificar base de datos
echo ""
echo -n "📊 Verificando base de datos '$DB_NAME'... "
if $PSQL_CMD -d postgres -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1; then
    echo -e "${GREEN}Existe${NC}"
else
    echo -e "${RED}No existe${NC}"
    echo "   Ejecuta: ./migrate.sh para crearla"
    exit 1
fi

# Contar tablas
echo ""
echo "📋 Verificando tablas..."
TABLES_COUNT=$($PSQL_CMD -d $DB_NAME -tc "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE';" | tr -d ' ')

if [ "$TABLES_COUNT" -eq "0" ]; then
    echo -e "   ${RED}No se encontraron tablas${NC}"
    echo "   Ejecuta: ./migrate.sh para crear las tablas"
    exit 1
else
    echo -e "   ${GREEN}$TABLES_COUNT tablas encontradas${NC}"
fi

# Listar tablas principales
echo ""
echo "✅ Tablas principales:"
$PSQL_CMD -d $DB_NAME -tc "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE' AND table_name IN ('users', 'chats', 'messages', 'gig_workers', 'tasks', 'notifications') ORDER BY table_name;" | grep -v "^$" | sed 's/^/   ✓ /'

# Verificar tablas de gig workers
echo ""
echo "👷 Tablas de Gig Workers:"
$PSQL_CMD -d $DB_NAME -tc "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE' AND table_name LIKE '%worker%' OR table_name LIKE '%task%' OR table_name = 'attendance_history' ORDER BY table_name;" | grep -v "^$" | sed 's/^/   ✓ /'

# Mostrar estadísticas
echo ""
echo "📈 Estadísticas:"

# Contar registros en tablas principales
for table in users gig_workers tasks notifications; do
    COUNT=$($PSQL_CMD -d $DB_NAME -tc "SELECT COUNT(*) FROM $table;" 2>/dev/null | tr -d ' ')
    if [ $? -eq 0 ]; then
        printf "   %-20s %s registros\n" "$table:" "$COUNT"
    fi
done

# Información de conexión
echo ""
echo -e "${BLUE}ℹ️  Información de conexión:${NC}"
echo "   Host: $DB_HOST"
echo "   Puerto: $DB_PORT"
echo "   Usuario: $DB_USER"
echo "   Base de datos: $DB_NAME"

if [ "$USE_DOCKER" = true ]; then
    echo ""
    echo "🔗 Conectar manualmente:"
    echo "   docker exec -it whatsapp_postgres psql -U $DB_USER -d $DB_NAME"
else
    echo ""
    echo "🔗 Conectar manualmente:"
    echo "   PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME"
fi

# Limpiar
if [ "$USE_DOCKER" = false ]; then
    unset PGPASSWORD
fi

echo ""
echo -e "${GREEN}✅ Verificación completa${NC}"