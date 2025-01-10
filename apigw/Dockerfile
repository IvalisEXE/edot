# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

# Set environment variables
ENV GO111MODULE=on

# Create app directory
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Copy the .env file (if needed during build)
COPY .env .env

# Build the Go application
RUN go build -o api .

# Stage 2: Create a minimal image
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/api .
# Copy the .env file from the builder stage (if needed at runtime)
COPY --from=builder /app/.env .

# Expose port 3000
EXPOSE 3000

# Command to run the binary
CMD ["./api"]
