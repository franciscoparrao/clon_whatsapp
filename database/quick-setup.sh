#!/bin/bash

# Script rápido de configuración para PostgreSQL local

echo "=== Configuración Rápida de PostgreSQL ==="
echo ""

# Colores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Intentar crear base de datos con el usuario actual
echo "Creando base de datos con tu usuario: $USER"

# Método 1: Con tu usuario del sistema
createdb whatsapp_clone 2>/dev/null

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Base de datos creada exitosamente${NC}"
    
    # Ejecutar schema
    echo "Ejecutando schema..."
    psql -d whatsapp_clone -f ./schema.sql
    
    echo ""
    echo -e "${GREEN}¡Listo! Base de datos configurada${NC}"
    echo ""
    echo "Configuración para el backend:"
    echo "  export DB_USER=$USER"
    echo "  export DB_PASSWORD=''"
    echo "  export DB_HOST=localhost"
    echo "  export DB_PORT=5432"
    echo "  export DB_NAME=whatsapp_clone"
    echo ""
    echo "O crea un archivo .env en el directorio backend con:"
    echo "  DB_USER=$USER"
    echo "  DB_PASSWORD="
    echo "  DB_HOST=localhost"
    echo "  DB_PORT=5432"
    echo "  DB_NAME=whatsapp_clone"
    echo "  DB_SSLMODE=disable"
else
    echo -e "${YELLOW}No se pudo crear la base de datos${NC}"
    echo ""
    echo "Intentando método alternativo con sudo..."
    
    # Método 2: Con sudo
    sudo -u postgres createdb whatsapp_clone 2>/dev/null
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Base de datos creada con usuario postgres${NC}"
        
        # Dar permisos a tu usuario
        sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE whatsapp_clone TO $USER;" 2>/dev/null
        
        # Ejecutar schema
        sudo -u postgres psql -d whatsapp_clone -f ./schema.sql
        
        echo ""
        echo "Usa estas credenciales:"
        echo "  DB_USER=postgres"
        echo "  DB_PASSWORD=postgres"
    else
        echo ""
        echo "Opciones:"
        echo "1. Ejecuta: sudo ./setup-local-postgres.sh"
        echo "2. O usa Docker: cd .. && ./start-db-docker.sh"
    fi
fi