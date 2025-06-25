#!/bin/bash

echo "📋 Cargando schema en PostgreSQL Docker..."
echo ""

# Verificar que el contenedor esté ejecutándose
if ! sudo docker ps | grep -q whatsapp_postgres; then
    echo "❌ El contenedor whatsapp_postgres no está ejecutándose"
    echo "   Ejecuta: cd .. && ./start-db-docker.sh"
    exit 1
fi

# Cargar schema
echo "Ejecutando schema.sql..."
sudo docker exec -i whatsapp_postgres psql -U postgres -d whatsapp_clone < schema.sql

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Schema cargado exitosamente"
    
    # Mostrar tablas creadas
    echo ""
    echo "Tablas creadas:"
    sudo docker exec whatsapp_postgres psql -U postgres -d whatsapp_clone -c '\dt' | grep "public\." | awk '{print "  ✓ " $3}'
else
    echo "❌ Error al cargar el schema"
fi

echo ""
echo "Para conectarte a la base de datos:"
echo "  sudo docker exec -it whatsapp_postgres psql -U postgres -d whatsapp_clone"