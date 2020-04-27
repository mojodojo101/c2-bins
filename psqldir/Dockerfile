FROM ubuntu:19.10

ARG C2DBUsername 
ARG C2DBPassword
ARG C2DBName

RUN apt-get update

RUN apt-get install -y postgresql postgresql-contrib

USER postgres


RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/11/main/pg_hba.conf

RUN echo "listen_addresses= '*'" >> /etc/postgresql/11/main/postgresql.conf

RUN    /etc/init.d/postgresql start &&\
    psql --command "CREATE USER $C2DBUsername WITH SUPERUSER PASSWORD '$C2DBPassword';" &&\
    psql --command "Create DATABASE $C2DBName;" &&\
    psql --command "grant all privileges on database $C2DBName to $C2DBUsername;"


EXPOSE 5432

VOLUME ["/etc/postgresql","/var/log/postgresql","/var/lib/postgresql"]

CMD ["/usr/lib/postgresql/11/bin/postgres","-D","/var/lib/postgresql/11/main","-c","config_file=/etc/postgresql/11/main/postgresql.conf"]
