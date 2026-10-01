FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o social-graph-service cmd/social-graph-service/main.go

FROM alpine:latest

COPY --from=builder /app/social-graph-service /usr/local/bin/social-graph-service

EXPOSE 8080

CMD ["social-graph-service"]







