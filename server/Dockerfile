FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/server
RUN CGO_ENABLED=0 GOOS=linux go build -o kv-server .

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/server/kv-server .

EXPOSE 8080
CMD ["./kv-server"]
