FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o api-gateway cmd/api-gateway/main.go


FROM alpine:latest

COPY --from=builder /app/api-gateway /usr/local/bin/api-gateway

EXPOSE 8080

CMD ["api-gateway"]

