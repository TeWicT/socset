FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o profile-service cmd/profile-service/main.go

FROM alpine:latest

COPY --from=builder /app/profile-service /usr/local/bin/profile-service

EXPOSE 8080

CMD ["profile-service"]

