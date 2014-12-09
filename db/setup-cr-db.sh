#!/bin/bash
echo "Creating database..."
gosu postgres postgres --single <<- EOSQL
   CREATE DATABASE carriagereturn;
   CREATE USER postgres;
   GRANT ALL PRIVILEGES ON DATABASE carriagereturn to postgres;
EOSQL
echo "Created"
echo "Importing db.sql..."
gosu postgres postgres --single -j carriagereturn < /var/lib/cr.sql
echo "Imported"
