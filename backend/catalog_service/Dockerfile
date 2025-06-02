FROM golang:1.24.3-alpine3.21 AS builder
ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN apk add --no-cache git \
    && go mod download \
    && go install -tags 'postgres,pgx' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .
RUN go build -ldflags="-s -w" -o catalog_service ./cmd/app

FROM alpine:3.21
RUN apk add --no-cache ca-certificates postgresql-client

WORKDIR /app

COPY --from=builder /app/migrations /migrations
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

COPY --from=builder /app/catalog_service /usr/local/bin/catalog_service

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
EXPOSE 8080 50050
