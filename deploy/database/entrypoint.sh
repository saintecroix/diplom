#!/bin/bash
set -e

if [[ -z "$POSTGRES_HOST_AUTH_METHOD" ]]; then
  sed -i "s/#listen_addresses = 'localhost'/listen_addresses = '*'/" "$PG_CONF"
fi

# Правило с scram-sha-256
if ! grep -q "host    all             keril             0.0.0.0/0            scram-sha-256" "$PG_HBA_CONF"; then
  echo "host    all             keril             0.0.0.0/0            scram-sha-256" >> "$PG_HBA_CONF"
fi

psql_ctl reload