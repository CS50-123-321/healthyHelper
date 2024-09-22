# Use a build argument for Go version
ARG GO_VERSION=1

# Stage 1: Build the Go binary
FROM golang:${GO_VERSION}-bookworm AS builder

# Set the working directory inside the container
WORKDIR /usr/src/app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy the source code
COPY . .

# Build the application
RUN go build -v -o run-app .

# Stage 2: Create the final lightweight image
FROM debian:bookworm

# Update CA certificates
RUN apt-get update && apt-get install -y ca-certificates

# Set the working directory
WORKDIR /usr/src/app

# Copy the Go binary from the builder stage to the working directory
COPY --from=builder /usr/src/app/run-app .

# Copy the static files from the builder stage to the working directory
COPY --from=builder /usr/src/app/static ./static

# Expose the application port
EXPOSE 8888

# Run the application
CMD ["./run-app"]