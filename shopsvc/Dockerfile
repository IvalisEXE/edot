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

# Install bash (required for running shell scripts)
RUN apk update && apk add --no-cache bash postgresql-client redis

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/api .
# Copy the .env file from the builder stage (if needed at runtime)
COPY --from=builder /app/.env .

# Copy the JSON file into the container
COPY pkg/data/shop.json /root/pkg/data/shop.json
RUN chmod +x /root/pkg/data/shop.json

# Copy the entrypoint.sh script and give it execution permission
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose port 3004
EXPOSE 3004

# Command to run the binary
CMD ["/bin/sh", "/entrypoint.sh"]
