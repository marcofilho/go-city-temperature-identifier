FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN go build -o main ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/main .

ENTRYPOINT ["./main"]