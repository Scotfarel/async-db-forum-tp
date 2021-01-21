FROM ubuntu:18.04

LABEL MAINTAINER="Ivan Zakharov <scotfarel@gmail.com>"

RUN ln -snf /usr/share/zoneinfo/Europe/Moscow /etc/localtime && echo Europe/Moscow > /etc/timezone

RUN apt-get -y update
RUN apt install -y git wget gcc gnupg


#
# Install postgresql, get and set a key.
#
ENV PGVER 12
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ bionic-pgdg main" > /etc/apt/sources.list.d/pgdg.list
RUN wget https://www.postgresql.org/media/keys/ACCC4CF8.asc && apt-key add ACCC4CF8.asc
RUN apt-get update && apt-get install -y postgresql-$PGVER
#ENV PGVER 10
#RUN apt -y update && apt install -y postgresql-$PGVER

#
# Install Go
#
RUN wget https://dl.google.com/go/go1.14.linux-amd64.tar.gz
RUN tar -xvf go1.14.linux-amd64.tar.gz
RUN mv go /usr/local

ENV GOROOT /usr/local/go
ENV GOPATH $HOME/go
ENV PATH $GOPATH/bin:$GOROOT/bin:$PATH

WORKDIR /server
COPY . .

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER iszakharov WITH SUPERUSER PASSWORD 'iszakharov';" &&\
    createdb -O iszakharov forums &&\
    psql forums -f /server/configs/init.sql &&\
    /etc/init.d/postgresql stop

RUN echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "listen_addresses='*'\nsynchronous_commit = off\nfsync = off\nshared_buffers = 512MB\neffective_cache_size = 1024MB\n" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "wal_buffers = 1MB\nwal_writer_delay = 2000ms\nrandom_page_cost = 1.0\nautovacuum = off\nmax_connections = 100\nwork_mem = 8MB\nmaintenance_work_mem = 128MB\ncpu_tuple_cost = 0.0030\ncpu_index_tuple_cost = 0.0010\ncpu_operator_cost = 0.0005" >> /etc/postgresql/$PGVER/main/postgresql.conf

RUN echo "log_statement = none" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_duration = off " >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_lock_waits = on" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_min_duration_statement = 50" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_filename = 'query.log'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_directory = '/var/log/postgresql'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_destination = 'csvlog'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "logging_collector = on" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_temp_files = '-1'" >> /etc/postgresql/$PGVER/main/postgresql.conf

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root
RUN go mod vendor
RUN go build -mod=vendor /server/cmd/main.go
CMD service postgresql start && ./main

EXPOSE 5000
