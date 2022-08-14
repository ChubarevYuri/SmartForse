# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

ENV old ''
ENV new ''

WORKDIR /

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go build -o /docker-gs-ping

CMD ["/docker-gs-ping"]