# ---------- Build stage ----------
FROM golang:1.23-alpine AS builder

# Install git (needed if some modules are from private repos)
RUN apk add --no-cache git

# Set working directory inside the container
# Usually better to keep it under /app for clarity
WORKDIR /app

# Copy go mod files first (for better layer caching)
COPY go.mod go.sum ./

# Download dependencies (cached if go.mod hasn't changed)
RUN go mod download

# Copy the rest of your source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ---------- Final stage ----------
FROM alpine:latest

# Install certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port (change if needed)
EXPOSE 8080

# Run the app
CMD ["./main"]
