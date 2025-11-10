FROM golang:1.25-alpine

# Устанавливаем build-base, который включает gcc
RUN apk add --no-cache build-base

# Устанавливаем strace
RUN apk add strace

# Устанавливаем goimports
RUN go install golang.org/x/tools/cmd/goimports@latest

RUN touch /var/log/tmp.log

WORKDIR /app

# RUN go mod download && go mod verify
CMD tail -f /var/log/tmp.log
