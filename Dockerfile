# Build stage
FROM golang:1.19-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o fofa-bot cmd/bot/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /build/fofa-bot .

# Create directories for cache and data
RUN mkdir -p /app/fofa_cache

# Volume for configuration and cache
VOLUME ["/app/config", "/app/fofa_cache"]

# Run as non-root user
RUN addgroup -g 1000 botuser && \
    adduser -D -u 1000 -G botuser botuser && \
    chown -R botuser:botuser /app

USER botuser

# Set timezone
ENV TZ=Asia/Shanghai

CMD ["./fofa-bot"]
