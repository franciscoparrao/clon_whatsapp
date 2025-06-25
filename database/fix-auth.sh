#!/bin/bash

echo "=== Configurando autenticación PostgreSQL para conexiones TCP/IP ==="
echo ""

# Verificar la línea actual
echo "Configuración actual:"
sudo grep -E "^host.*all.*all.*127.0.0.1" /etc/postgresql/16/main/pg_hba.conf

echo ""
echo "Si ves 'peer' o 'ident' en lugar de 'md5' o 'scram-sha-256', ejecuta:"
echo ""
echo "sudo nano /etc/postgresql/16/main/pg_hba.conf"
echo ""
echo "Y asegúrate de que esta línea esté presente:"
echo "host    all             all             127.0.0.1/32            md5"
echo ""
echo "Luego reinicia PostgreSQL:"
echo "sudo systemctl restart postgresql"
echo ""
echo "O ejecuta este comando para hacerlo automáticamente:"
echo "sudo sed -i 's/host.*all.*all.*127.0.0.1.*/host    all             all             127.0.0.1\/32            md5/' /etc/postgresql/16/main/pg_hba.conf && sudo systemctl restart postgresql"