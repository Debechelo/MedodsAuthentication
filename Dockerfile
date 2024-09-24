FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth-service ./cmd/main.go


FROM alpine:latest

RUN apk add --no-cache bash

ENV DB_HOST=postgres
ENV DB_PORT=5432
ENV DB_USER=user
ENV DB_PASSWORD=password
ENV DB_NAME=authdb

COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh
COPY --from=builder /app/auth-service /usr/local/bin/auth-service

CMD ["wait-for-it.sh", "postgres:5432", "--", "/usr/local/bin/auth-service"]
#CMD ["postgres:5432", "--", "/usr/local/bin/auth-service"]