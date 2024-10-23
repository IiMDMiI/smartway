FROM golang:1.23.2-bookworm AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o goserver cmd/smartway/main.go


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/goserver .

EXPOSE 8080

CMD ["./goserver"]