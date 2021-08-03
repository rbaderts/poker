#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE ROLE poker WITH LOGIN PASSWORD 'poker' ;
    CREATE DATABASE poker;
    GRANT ALL PRIVILEGES ON DATABASE poker TO poker;

    CREATE ROLE keycloak WITH LOGIN PASSWORD 'keycloak';
    CREATE DATABASE keycloak;
    GRANT ALL PRIVILEGES ON DATABASE keycloak TO keycloak;
EOSQL

