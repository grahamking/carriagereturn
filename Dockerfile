# docker build -t=grahamking/carriagereturn .
# docker run --name=carriagereturn -d --link crdb:crdb -p 8082:8082 grahamking/carriagereturn
FROM debian:stable
MAINTAINER Graham King <graham@gkgk.org>
RUN mkdir -p /opt/cr && chown www-data:www-data /opt/cr
COPY index.atom /opt/cr/
COPY index.html /opt/cr/
COPY cr /opt/cr/
RUN chown www-data:www-data /opt/cr/*
EXPOSE 8082
USER www-data
CMD ["/opt/cr/cr", "-p", "8082", "-r", "/opt/cr/", "-h", "crdb"]
