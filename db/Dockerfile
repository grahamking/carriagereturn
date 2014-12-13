# docker build -t=grahamking/crdb .
# docker run --rm --name=crdb -v /tmp/pdata:/var/lib/postgresql/data grahamking/crdb
FROM postgres:9.3
COPY cr.sql /var/lib/
COPY setup-cr-db.sh /docker-entrypoint-initdb.d/
