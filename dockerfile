# setup base image
FROM golang:latest as builder
#FROM golang:1.17.0-alpine
RUN apt-get update && apt-get install -y ca-certificates openssl

ARG cert_location=/usr/local/share/ca-certificates

# Get certificate from "github.com"
RUN openssl s_client -showcerts -connect github.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/github.crt
# Get certificate from "proxy.golang.org"
RUN openssl s_client -showcerts -connect proxy.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM >  ${cert_location}/proxy.golang.crt
# Update certificates
RUN update-ca-certificates

WORKDIR /app

COPY ./ /app

COPY ./db/01-init.sh /docker-entrypoint-initdb.d/

#RUN apk --no-cache add ca-certificates
RUN go env -w GOPROXY=direct GOFLAGS="-insecure"
RUN go mod tidy

ENTRYPOINT [ "go", "run", "main.go" ]