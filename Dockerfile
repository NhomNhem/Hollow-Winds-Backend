# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Invalidate cache if needed
ARG BUILD_DATE=unknown
RUN echo "Building at ${BUILD_DATE}"

# Build the application
RUN rm -f server && CGO_ENABLED=0 GOOS=linux go build -v -o server ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/server .

# Expose port
EXPOSE 8080

# Run
CMD ["./server"]
