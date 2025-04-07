# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy source code first
COPY . .

# Initialize go module if it doesn't exist
RUN if [ ! -f go.mod ]; then \
    go mod init fe-kafka-board && \
    go mod tidy; \
    fi

# Change to scr directory and build
WORKDIR /app/scr
RUN CGO_ENABLED=0 GOOS=linux go build -o ../main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port 9080
EXPOSE 9080

# Run the application
CMD ["./main"] 