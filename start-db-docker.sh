#!/bin/bash

# Script rápido para iniciar PostgreSQL con Docker

echo "🚀 Iniciando PostgreSQL con Docker..."

# Colores
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m'

# Verificar si el puerto 5432 está en uso
echo -n "🔍 Verificando puerto 5432... "
if lsof -Pi :5432 -sTCP:LISTEN -t >/dev/null 2>&1 || netstat -tuln 2>/dev/null | grep -q ':5432 '; then
    echo -e "${YELLOW}En uso${NC}"
    echo ""
    echo "⚠️  El puerto 5432 ya está en uso. Opciones:"
    echo ""
    echo "1) Detener PostgreSQL local y usar Docker"
    echo "2) Usar PostgreSQL local existente"
    echo "3) Cambiar puerto de Docker (ej: 5433)"
    echo ""
    read -p "Selecciona una opción (1-3): " choice
    
    case $choice in
        1)
            echo "Deteniendo PostgreSQL local..."
            sudo systemctl stop postgresql 2>/dev/null || sudo service postgresql stop 2>/dev/null
            sleep 2
            ;;
        2)
            echo -e "${GREEN}Usando PostgreSQL local existente${NC}"
            echo ""
            echo "📊 Configuración:"
            echo "   Host: localhost"
            echo "   Puerto: 5432"
            echo "   Usuario: postgres (o tu usuario local)"
            echo ""
            echo "Para crear la base de datos:"
            echo "   cd database && ./migrate.sh"
            exit 0
            ;;
        3)
            echo "Configurando Docker para usar puerto 5433..."
            # Crear archivo temporal con puerto modificado
            sed 's/5432:5432/5433:5432/g' docker-compose.db.yml > docker-compose.db.tmp.yml
            COMPOSE_FILE="docker-compose.db.tmp.yml"
            DOCKER_PORT="5433"
            ;;
        *)
            echo -e "${RED}Opción inválida${NC}"
            exit 1
            ;;
    esac
else
    echo -e "${GREEN}Libre${NC}"
fi

# Usar archivo de compose correcto
COMPOSE_FILE=${COMPOSE_FILE:-"docker-compose.db.yml"}
DOCKER_PORT=${DOCKER_PORT:-"5432"}

# Detener contenedor si existe
echo "🛑 Deteniendo contenedor existente..."
sudo docker compose -f $COMPOSE_FILE down 2>/dev/null

# Iniciar PostgreSQL
echo "🚀 Iniciando PostgreSQL con Docker..."
sudo docker compose -f $COMPOSE_FILE up -d

# Esperar a que esté listo
echo -n "⏳ Esperando a que PostgreSQL esté listo..."
for i in {1..30}; do
    if docker exec whatsapp_postgres pg_isready -U postgres &>/dev/null; then
        echo " ✅"
        break
    fi
    echo -n "."
    sleep 1
done

echo ""
echo "📊 Base de datos PostgreSQL iniciada:"
echo "   Host: localhost"
echo "   Puerto: $DOCKER_PORT"
echo "   Usuario: postgres"
echo "   Contraseña: postgres123"
echo "   Base de datos: whatsapp_clone"
echo ""
echo "🔗 Para conectarte:"
echo "   docker exec -it whatsapp_postgres psql -U postgres -d whatsapp_clone"
echo ""
echo "🛑 Para detener:"
echo "   sudo docker compose -f $COMPOSE_FILE down"

# Si usamos puerto alternativo, recordar actualizar la configuración
if [ "$DOCKER_PORT" != "5432" ]; then
    echo ""
    echo -e "${YELLOW}⚠️  IMPORTANTE: Usando puerto $DOCKER_PORT${NC}"
    echo "   Actualiza tu archivo .env o variables de entorno:"
    echo "   DB_PORT=$DOCKER_PORT"
fi

# Limpiar archivo temporal si existe
[ -f "docker-compose.db.tmp.yml" ] && rm -f docker-compose.db.tmp.yml