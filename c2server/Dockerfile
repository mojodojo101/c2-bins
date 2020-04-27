FROM golang:latest

RUN useradd -d /home/c2server/ -m -p c2server -s /bin/bash c2server

RUN echo "c2server:c2server" | chpasswd

COPY . src/github.com/mojodojo101/c2server

WORKDIR src/github.com/mojodojo101/c2server/cmd/c2api

RUN go build github.com/mojodojo101/c2server/cmd/c2api

USER c2server

RUN mkdir /home/c2server/api

RUN mkdir /home/c2server/api/beacons 

RUN mkdir /home/c2server/api/targets

EXPOSE 8888

RUN chmod 744 /home/c2server/api/*

RUN cp server.key /home/c2server/api/server.key

RUN cp server.crt /home/c2server/api/server.crt

RUN cp c2api /home/c2server/api/c2api

WORKDIR /home/c2server/api

CMD ./c2api

