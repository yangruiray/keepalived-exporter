# Download pkg env

FROM golang:1.12
RUN go version
RUN mkdir -p /go/src/github.com/keepalived-exporter
RUN mkdir -p /etc/keepalived
COPY . /go/src/github.com/keepalived-exporter
WORKDIR /go/src/github.com/keepalived-exporter

ARG GOOS=linux
ARG GOARCH=amd64

# Maintainer
MAINTAINER yangrui@kpaas.io

# Set env
#ENV PORT 9999
#EXPOSE $PORT

# build 
RUN GOOS=linux GOARCH=amd64 go build -o /go/src/github.com/keepalived-exporter/cmd/keepalived-exporter/keepalived-exporter /go/src/github.com/keepalived-exporter/cmd/keepalived-exporter/main.go && cp /go/src/github.com/keepalived-exporter/cmd/keepalived-exporter/keepalived-exporter /usr/local/bin/ && chmod +x /usr/local/bin/keepalived-exporter

# Run container command 
ENTRYPOINT ["/usr/local/bin/keepalived-exporter"]

