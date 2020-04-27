FROM golang:latest

RUN useradd -d /home/c2beacon/ -m -p beacon -s /bin/bash beacon

RUN echo "beacon:beacon" | chpasswd
ADD . /go/src/github.com/mojodojo101/beacon

RUN go install /go/src/github.com/mojodojo101/beacon

COPY mojogreetings.txt /home/c2beacon/mojogreetings.txt

CMD /go/bin/beacon
