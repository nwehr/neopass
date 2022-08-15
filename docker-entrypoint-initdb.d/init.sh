#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB"  POSTGRES_PASSWORD<<-EOSQL
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    create table if not exists entries (
        client_uuid uuid not null,
        name varchar(255) not null,
        password text
    );
EOSQL