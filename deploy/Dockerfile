# Building
FROM golang:1.27.1-trixie AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY /cmd      ./cmd
COPY /internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Launching
FROM alpine:3.24.2 AS runtime

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]