# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Install required build tools
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./

# Izinkan penggunaan versi Go yang lebih baru
ENV GOTOOLCHAIN=auto

RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Final stage
FROM alpine:3.18
WORKDIR /app

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]