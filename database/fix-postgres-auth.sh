#!/bin/bash

echo "=== Configurando autenticación de PostgreSQL ==="
echo ""

# Colores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Configurar contraseña para usuario postgres
echo "Configurando contraseña para usuario postgres..."
sudo -u postgres psql <<EOF
ALTER USER postgres PASSWORD 'postgres';
\q
EOF

echo -e "${GREEN}✓ Contraseña configurada${NC}"

# Encontrar y actualizar pg_hba.conf
echo ""
echo "Actualizando configuración de autenticación..."

# Buscar el archivo pg_hba.conf
PG_VERSION=$(ls /etc/postgresql/ | head -1)
HBA_FILE="/etc/postgresql/$PG_VERSION/main/pg_hba.conf"

if [ -f "$HBA_FILE" ]; then
    echo "Archivo encontrado: $HBA_FILE"
    
    # Backup
    sudo cp "$HBA_FILE" "$HBA_FILE.backup.$(date +%Y%m%d_%H%M%S)"
    
    # Cambiar método de autenticación local de peer a md5
    sudo sed -i 's/local   all             postgres                                peer/local   all             postgres                                md5/' "$HBA_FILE"
    sudo sed -i 's/local   all             all                                     peer/local   all             all                                     md5/' "$HBA_FILE"
    
    echo -e "${GREEN}✓ Configuración actualizada${NC}"
    
    # Reiniciar PostgreSQL
    echo "Reiniciando PostgreSQL..."
    sudo systemctl restart postgresql
    sleep 3
    
    echo -e "${GREEN}✓ PostgreSQL reiniciado${NC}"
else
    echo -e "${YELLOW}No se encontró pg_hba.conf${NC}"
fi

# Probar conexión
echo ""
echo "Probando conexión..."
PGPASSWORD=postgres psql -h localhost -U postgres -d whatsapp_clone -c '\dt' > /dev/null 2>&1

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Conexión exitosa${NC}"
    echo ""
    echo "Base de datos lista para usar con:"
    echo "  Usuario: postgres"
    echo "  Contraseña: postgres"
else
    echo -e "${YELLOW}La conexión aún requiere configuración${NC}"
    echo ""
    echo "Alternativa: Usa tu usuario del sistema"
    echo "  Actualiza el archivo backend/.env:"
    echo "  DB_USER=$USER"
    echo "  DB_PASSWORD="
fi