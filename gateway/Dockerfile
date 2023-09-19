# Build stage
FROM golang:1.20.7-alpine3.18 AS builder
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

# Run stage
FROM alpine:3.18
WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8081

ENTRYPOINT [ "./main" ]
