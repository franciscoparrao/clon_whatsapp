#!/bin/bash

# Script de migración de base de datos para Clon de WhatsApp

echo "=== Script de Migración de Base de Datos ==="
echo ""

# Configuración predeterminada
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres123}"
DB_NAME="${DB_NAME:-whatsapp_clone}"

# Colores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "Configuración de base de datos:"
echo "  Host: $DB_HOST"
echo "  Puerto: $DB_PORT"
echo "  Usuario: $DB_USER"
echo "  Base de datos: $DB_NAME"
echo ""

# Función para verificar si estamos usando Docker
check_docker() {
    if docker ps --format "table {{.Names}}" | grep -q "whatsapp_postgres"; then
        echo -e "${GREEN}Detectado contenedor Docker 'whatsapp_postgres'${NC}"
        USE_DOCKER=true
    else
        USE_DOCKER=false
    fi
}

# Verificar si PostgreSQL está disponible
echo -n "Verificando conexión a PostgreSQL... "

# Primero intentar con Docker si está disponible
check_docker

if [ "$USE_DOCKER" = true ]; then
    if docker exec whatsapp_postgres pg_isready -U postgres &>/dev/null; then
        echo -e "${GREEN}OK (Docker)${NC}"
        # Redefinir comandos para usar Docker
        PSQL_CMD="docker exec -i whatsapp_postgres psql -U $DB_USER"
        CREATE_DB_CMD="docker exec whatsapp_postgres createdb -U $DB_USER"
    else
        echo -e "${YELLOW}Docker detectado pero PostgreSQL no responde${NC}"
        USE_DOCKER=false
    fi
fi

# Si no usamos Docker, intentar conexión local
if [ "$USE_DOCKER" = false ]; then
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c '\q' 2>/dev/null
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}OK (Local)${NC}"
        # Comandos para PostgreSQL local
        export PGPASSWORD=$DB_PASSWORD
        PSQL_CMD="psql -h $DB_HOST -p $DB_PORT -U $DB_USER"
        CREATE_DB_CMD="createdb -h $DB_HOST -p $DB_PORT -U $DB_USER"
    else
        echo -e "${RED}ERROR${NC}"
        echo ""
        echo "No se pudo conectar a PostgreSQL. Opciones:"
        echo ""
        echo "1. Iniciar PostgreSQL con Docker:"
        echo "   cd .. && ./start-db-docker.sh"
        echo ""
        echo "2. Iniciar PostgreSQL local:"
        echo "   sudo systemctl start postgresql"
        echo ""
        echo "3. Verificar credenciales en variables de entorno:"
        echo "   export DB_USER=tu_usuario"
        echo "   export DB_PASSWORD=tu_password"
        exit 1
    fi
fi

# Crear base de datos si no existe
echo -n "Creando base de datos '$DB_NAME' si no existe... "
if $PSQL_CMD -d postgres -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1; then
    echo -e "${GREEN}Ya existe${NC}"
else
    $CREATE_DB_CMD $DB_NAME
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}Creada${NC}"
    else
        echo -e "${RED}ERROR${NC}"
        echo "Intentando crear con comando alternativo..."
        $PSQL_CMD -d postgres -c "CREATE DATABASE $DB_NAME;"
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}Creada (método alternativo)${NC}"
        else
            echo -e "${RED}No se pudo crear la base de datos${NC}"
            exit 1
        fi
    fi
fi

# Ejecutar script de schema
echo -n "Ejecutando script de schema... "
if $PSQL_CMD -d $DB_NAME -f ./schema.sql > /tmp/migrate_output.log 2>&1; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${YELLOW}Advertencia${NC}"
    echo "Algunos comandos pueden haber fallado (esto es normal si las tablas ya existen)."
    echo ""
    echo "Verificando tablas creadas..."
    
    # Verificar si las tablas principales existen
    TABLES_COUNT=$($PSQL_CMD -d $DB_NAME -tc "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE';" | tr -d ' ')
    
    if [ "$TABLES_COUNT" -gt "0" ]; then
        echo -e "${GREEN}✓ Se encontraron $TABLES_COUNT tablas en la base de datos${NC}"
    else
        echo -e "${RED}ERROR: No se encontraron tablas${NC}"
        echo "Mostrando detalles del error:"
        cat /tmp/migrate_output.log
        exit 1
    fi
fi

# Limpiar archivo temporal
rm -f /tmp/migrate_output.log

echo ""
echo -e "${GREEN}¡Migración completada exitosamente!${NC}"
echo ""
echo "Tablas creadas:"
echo "  - users (usuarios del sistema)"
echo "  - privacy_settings (configuraciones de privacidad)"
echo "  - contacts (lista de contactos)"
echo "  - chats (conversaciones)"
echo "  - chat_participants (participantes de chats)"
echo "  - messages (mensajes)"
echo "  - message_status (estado de mensajes)"
echo "  - gig_workers (trabajadores)"
echo "  - worker_availability (disponibilidad de trabajadores)"
echo "  - tasks (tareas/turnos)"
echo "  - task_assignments (asignaciones de tareas)"
echo "  - attendance_history (historial de asistencia)"
echo "  - notifications (notificaciones)"
echo ""
echo "Para conectarte a la base de datos:"
if [ "$USE_DOCKER" = true ]; then
    echo "  docker exec -it whatsapp_postgres psql -U $DB_USER -d $DB_NAME"
else
    echo "  PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME"
fi

# Mostrar lista de tablas creadas
echo ""
echo "Verificando tablas..."
$PSQL_CMD -d $DB_NAME -c "\dt" | grep -E "^\s*public\s*\|" | awk '{print "  ✓ " $3}'

# Limpiar variable de entorno si no usamos Docker
if [ "$USE_DOCKER" = false ]; then
    unset PGPASSWORD
fi