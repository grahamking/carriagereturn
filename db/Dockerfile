# docker build -t=localhost:5000/crdb .
# docker run --rm --name=crdb -v /tmp/pdata:/var/lib/postgresql/data localhost:5000/crdb
FROM postgres:9.3
COPY cr.sql /var/lib/
COPY setup-cr-db.sh /docker-entrypoint-initdb.d/
