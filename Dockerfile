# Build stage
FROM golang:1.23-alpine AS builder

ARG PB_VERSION=0.29.0

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy Go module files
COPY go.mod ./

# Copy source code
COPY main.go .

# Generate go.sum and download dependencies
RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o pocketbase-custom .

# Runtime stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

# Copy the built binary
COPY --from=builder /app/pocketbase-custom /usr/local/bin/pocketbase-custom

# Create data directory
RUN mkdir -p /pb/pb_data

# Set volume for data persistence
VOLUME ["/pb/pb_data"]

EXPOSE 8080

# start custom PocketBase with data directory
CMD ["/usr/local/bin/pocketbase-custom", "serve", "--http=0.0.0.0:8080", "--dir=/pb/pb_data"]
