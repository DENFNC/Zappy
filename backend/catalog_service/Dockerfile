FROM golang:1.24.3-alpine3.21 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -ldflags="-s -w" -o catalog_service ./cmd/app

FROM scratch

COPY --from=builder /app/catalog_service /usr/local/bin/catalog_service
EXPOSE 8082 50052

ENTRYPOINT ["/usr/local/bin/catalog_service"]