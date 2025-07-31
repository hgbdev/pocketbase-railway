# Build stage
FROM golang:1.23-alpine AS builder

ARG PB_VERSION=0.23.11
ARG PB_TOKEN_SECRET

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy Go module files
COPY go.mod ./

# Copy source code
COPY main.go .

# Generate go.sum and download dependencies
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o pocketbase-custom .

# Runtime stage
FROM alpine:latest

ARG PB_TOKEN_SECRET
ENV PB_TOKEN_SECRET=${PB_TOKEN_SECRET}

RUN apk add --no-cache ca-certificates

# Copy the built binary
COPY --from=builder /app/pocketbase-custom /usr/local/bin/pocketbase-custom

EXPOSE 8080

# start custom PocketBase
CMD ["/usr/local/bin/pocketbase-custom", "serve", "--http=0.0.0.0:8080"]
