# Multi-stage build for Go backend
FROM golang:1.20-alpine AS build-stage

WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o booksrendezvous-backend .

# Production stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary and data files
COPY --from=build-stage /app/booksrendezvous-backend .
COPY --from=build-stage /app/data ./data

# Change ownership
RUN chown -R appuser:appgroup /root/

USER appuser

EXPOSE 6050

CMD ["./booksrendezvous-backend"]

