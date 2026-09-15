FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o auth-service cmd/auth-service/main.go

FROM alpine:latest

COPY --from=builder /app/auth-service /usr/local/bin/auth-service

EXPOSE 8080

CMD ["auth-service"]

