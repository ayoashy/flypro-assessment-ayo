# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Install goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o flypro-assessment ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Install postgresql-client for psql (optional, for manual DB access)
RUN apk add --no-cache postgresql-client

# Copy goose binary from builder
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Copy binary from builder
COPY --from=builder /app/flypro-assessment .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Copy Makefile
COPY --from=builder /app/Makefile .

EXPOSE 8080

CMD ["./flypro-assessment"]