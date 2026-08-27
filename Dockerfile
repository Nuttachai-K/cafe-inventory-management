# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /bootstrap ./cmd/bootstrap

FROM alpine:3.20

RUN apk add --no-cache ca-certificates \
    && addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=builder /server /app/server
COPY --from=builder /bootstrap /app/bootstrap

USER app
EXPOSE 8080

ENTRYPOINT ["/app/server"]