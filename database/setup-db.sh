#!/bin/bash

# Script para configurar la base de datos PostgreSQL

echo "=== Configuración de Base de Datos para WhatsApp Clone ==="
echo ""

# Colores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Función para mostrar el menú
show_menu() {
    echo "Selecciona una opción:"
    echo "1) Usar Docker (recomendado)"
    echo "2) Usar PostgreSQL local"
    echo "3) Salir"
}

# Opción 1: Configurar con Docker
setup_docker() {
    echo -e "${YELLOW}Configurando PostgreSQL con Docker...${NC}"
    
    # Verificar si Docker está instalado
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}Docker no está instalado. Por favor, instala Docker primero.${NC}"
        echo "Visita: https://docs.docker.com/get-docker/"
        exit 1
    fi
    
    # Ir al directorio raíz del proyecto
    cd ..
    
    # Detener contenedor existente si existe
    docker-compose -f docker-compose.db.yml down 2>/dev/null
    
    # Iniciar PostgreSQL con Docker Compose
    echo "Iniciando PostgreSQL con Docker..."
    docker-compose -f docker-compose.db.yml up -d
    
    # Esperar a que PostgreSQL esté listo
    echo -n "Esperando a que PostgreSQL esté listo..."
    for i in {1..30}; do
        if docker exec whatsapp_postgres pg_isready -U postgres &>/dev/null; then
            echo -e " ${GREEN}OK${NC}"
            break
        fi
        echo -n "."
        sleep 1
    done
    
    # Verificar si la base de datos se creó
    echo -n "Verificando base de datos... "
    if docker exec whatsapp_postgres psql -U postgres -d whatsapp_clone -c '\dt' &>/dev/null; then
        echo -e "${GREEN}Base de datos creada exitosamente${NC}"
    else
        echo -e "${YELLOW}Creando schema...${NC}"
        docker exec -i whatsapp_postgres psql -U postgres -d whatsapp_clone < database/schema.sql
    fi
    
    echo ""
    echo -e "${GREEN}¡PostgreSQL está listo!${NC}"
    echo ""
    echo "Información de conexión:"
    echo "  Host: localhost"
    echo "  Puerto: 5432"
    echo "  Usuario: postgres"
    echo "  Contraseña: postgres123"
    echo "  Base de datos: whatsapp_clone"
    echo ""
    echo "Para conectarte manualmente:"
    echo "  docker exec -it whatsapp_postgres psql -U postgres -d whatsapp_clone"
}

# Opción 2: Configurar con PostgreSQL local
setup_local() {
    echo -e "${YELLOW}Configurando PostgreSQL local...${NC}"
    
    # Solicitar credenciales
    read -p "Usuario de PostgreSQL (default: $USER): " pg_user
    pg_user=${pg_user:-$USER}
    
    read -sp "Contraseña de PostgreSQL: " pg_password
    echo ""
    
    # Configurar variables de entorno temporalmente
    export PGPASSWORD=$pg_password
    
    # Verificar conexión
    echo -n "Verificando conexión a PostgreSQL... "
    if psql -h localhost -U $pg_user -d postgres -c '\q' 2>/dev/null; then
        echo -e "${GREEN}OK${NC}"
    else
        echo -e "${RED}ERROR${NC}"
        echo "No se pudo conectar. Verifica las credenciales."
        exit 1
    fi
    
    # Crear base de datos si no existe
    echo -n "Creando base de datos 'whatsapp_clone'... "
    if psql -h localhost -U $pg_user -d postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'whatsapp_clone'" | grep -q 1; then
        echo -e "${GREEN}Ya existe${NC}"
    else
        createdb -h localhost -U $pg_user whatsapp_clone
        echo -e "${GREEN}Creada${NC}"
    fi
    
    # Ejecutar schema
    echo -n "Ejecutando schema... "
    if psql -h localhost -U $pg_user -d whatsapp_clone -f ./schema.sql > /dev/null 2>&1; then
        echo -e "${GREEN}OK${NC}"
    else
        echo -e "${RED}ERROR${NC}"
        echo "Ejecutando con output completo:"
        psql -h localhost -U $pg_user -d whatsapp_clone -f ./schema.sql
    fi
    
    # Limpiar variable de entorno
    unset PGPASSWORD
    
    echo ""
    echo -e "${GREEN}¡Base de datos configurada!${NC}"
    echo ""
    echo "Actualiza las variables de entorno en tu backend:"
    echo "  DB_USER=$pg_user"
    echo "  DB_PASSWORD=<tu_contraseña>"
    echo "  DB_NAME=whatsapp_clone"
}

# Menú principal
while true; do
    show_menu
    read -p "Opción: " choice
    
    case $choice in
        1)
            setup_docker
            break
            ;;
        2)
            setup_local
            break
            ;;
        3)
            echo "Saliendo..."
            exit 0
            ;;
        *)
            echo -e "${RED}Opción inválida${NC}"
            ;;
    esac
done