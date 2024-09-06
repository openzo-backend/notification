# Stage 1: Build the Go application
FROM golang:1.20-alpine AS builder

# Install git (needed for go mod download)
RUN apk add --no-cache git

# Set the working directory to /app
WORKDIR /app

# Copy go.mod and go.sum first, and download dependencies (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main .

# Stage 2: Create the final lightweight image
FROM alpine:latest

# Set the working directory in the minimal container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the application's port
EXPOSE 50053

# Command to run the executable
CMD ["./main"]
