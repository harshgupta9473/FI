# syntax=docker/dockerfile:1
FROM golang:1.21-alpine

# Set working directory
WORKDIR /app

# Install Go dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the full source code
COPY . .

# Build the Go app
RUN go build -o main ./cmd/server

# Expose the port your app runs on
EXPOSE 8080

# Run the app
CMD ["./main"]
