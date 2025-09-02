# Use official Golang image as the build environment
FROM golang:1.20-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Install git for go modules and other dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire backend source code
COPY . .

# Build the Go app
RUN go build -o api.exe ./cmd/api/main.go

# Use a minimal image for the final container
FROM alpine:latest

# Install ca-certificates for HTTPS support
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/api.exe .

# Expose the port the app runs on (default 8080 or as per config)
EXPOSE 8080

# Command to run the executable
CMD ["./api.exe"]
