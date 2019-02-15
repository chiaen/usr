FROM golang:1.11.5-alpine3.9

RUN apk add --no-cache ca-certificates curl git sed gcc libc-dev bash findutils openssh && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/github.com/chiaen/usr
COPY . .
ENV GO111MODULE=on
RUN go mod download
WORKDIR /root/
