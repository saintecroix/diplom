#!/bin/bash
set -e

# Явные пути к конфигам
PG_CONF="/var/lib/postgresql/data/postgresql.conf"
PG_HBA_CONF="/var/lib/postgresql/data/pg_hba.conf"

# Настройка listen_addresses
if [[ -z "$POSTGRES_HOST_AUTH_METHOD" ]]; then
  sed -i "s/#listen_addresses = 'localhost'/listen_addresses = '*'/" "$PG_CONF"
fi

# Перезагрузка конфигурации
pg_ctl reload