# Build stage
FROM golang:1.22-alpine3.19 AS builder

# Install build dependencies
RUN apk update && apk add --no-cache curl ca-certificates

# Set the working directory
WORKDIR /app

# Copy all the Go files and modules into the container
COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

# Build the Go binary with CGO disabled for portability
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

# Final stage (Runtime)
FROM alpine:3.19

# Install required dependencies for the runtime environment
RUN apk add --no-cache libc6-compat ca-certificates && \
    update-ca-certificates

# Copy the Go binary from the build stage
COPY --from=builder /app/myapp /usr/local/bin/myapp

# Set the entry point to the compiled Go binary
CMD ["/usr/local/bin/myapp"]
