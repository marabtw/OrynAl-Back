#!/bin/bash

PG_USER="postgres"
PG_DB="orynal_db"
PG_CONTAINER="orynal_pg"

SQL_QUERY="INSERT INTO users (name, surname, email, phone, role, password) VALUES ('Admin', 'Admin', 'admin@example.com', '1234567890', 'admin', '\$2a\$10\$HcCB.d98tfCe6gX6HexEHOfaFsK20bdB09YkI.1fHlx/3WxDhFxpO');"

docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -c "$SQL_QUERY"
