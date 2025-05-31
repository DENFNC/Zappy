FROM golang:1.24.3-alpine3.21 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ARG POSTGRES_URL
ENV POSTGRES_URL=${POSTGRES_URL}

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -ldflags="-s -w" -o review_service ./cmd/app

FROM scratch

COPY --from=builder /app/review_service /usr/local/bin/review_service
EXPOSE 8082 50052

ENTRYPOINT ["/usr/local/bin/review_service"]