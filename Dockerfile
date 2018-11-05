FROM golang:latest

ADD . /go/src/github.com/cj123/basicauth-proxy
WORKDIR /go/src/github.com/cj123/basicauth-proxy

RUN go get .
RUN go build

EXPOSE 8766

ENTRYPOINT /go/bin/basicauth-proxy