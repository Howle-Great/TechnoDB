FROM ubuntu:18.04
MAINTAINER Howle
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y gnupg
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y git

# Клонируем проект
USER root
RUN git clone https://github.com/Howle/TechnoDB.git
WORKDIR TechnoDB

# Устанавливаем PostgreSQL
RUN apt-get -y update
RUN apt-get -y install apt-transport-https git wget
RUN echo 'deb http://apt.postgresql.org/pub/repos/apt/ bionic-pgdg main' >> /etc/apt/sources.list.d/pgdg.list
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt-get -y update
ENV PGVERSION 11
RUN apt-get -y install postgresql-$PGVERSION postgresql-contrib

# Подключаемся к PostgreSQL и создаем БД
USER postgres
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    psql -d docker -c "CREATE EXTENSION IF NOT EXISTS citext;" &&\
    psql docker -a -f  db/init.sql &&\
    /etc/init.d/postgresql stop

RUN cat db/init.sql

USER root
# Настраиваем сеть и БД
RUN echo "local all all md5" > /etc/postgresql/$PGVERSION/main/pg_hba.conf &&\
    echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/$PGVERSION/main/pg_hba.conf
RUN cat database/postgresql.conf >> /etc/postgresql/$PGVERSION/main/postgresql.conf
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
EXPOSE 5432

# Устанавливаем Golang 
ENV GOVERSION 1.11.4
USER root
RUN wget https://storage.googleapis.com/golang/go$GOVERSION.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go$GOVERSION.linux-amd64.tar.gz && \
    mkdir go && mkdir go/src && mkdir go/bin && mkdir go/pkg
ENV GOROOT /usr/local/go
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/bin" "$GOPATH/src"
RUN apt-get -y install gcc musl-dev && GO11MODULE=on
ENV GOBIN $GOPATH/bin
RUN go get
RUN go build main.go
EXPOSE 5000

# Запускаем PostgreSQL и api сервер
CMD service postgresql start && ./main