FROM golang:1.23.1-alpine3.20 AS builder

WORKDIR /app

COPY shared /app/shared

COPY proto_buffer /app/proto_buffer

COPY order /app/order

WORKDIR /app/order

RUN go mod tidy

RUN go mod download

WORKDIR /app/order

RUN CGO_ENABLED=0 go build -o /app/bin .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin .

EXPOSE 8080

CMD ["./bin"]