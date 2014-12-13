#!/bin/bash
echo "Creating database..."
gosu postgres postgres --single <<- EOSQL
	CREATE DATABASE carriagereturn WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8';
	ALTER DATABASE carriagereturn OWNER TO postgres;
EOSQL
echo "Created"
echo "Running cr.sql..."
gosu postgres postgres --single -j carriagereturn < /var/lib/cr.sql
echo "Done"
