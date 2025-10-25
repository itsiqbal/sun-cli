# Minimal Dockerfile using scratch (smallest possible image)
# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build arguments
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${BUILD_DATE} -extldflags '-static'" \
    -tags="netgo,osusergo" \
    -installsuffix netgo \
    -o sun-cli .

# Final stage - scratch (absolute minimal)
FROM scratch

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data (optional, if your app needs it)
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binary
COPY --from=builder /app/sun-cli /sun-cli

# Create /tmp directory for temporary files
COPY --from=builder --chown=65534:65534 /tmp /tmp

# Use nobody user (non-root)
USER 65534:65534

# Set entrypoint
ENTRYPOINT ["/sun-cli"]

# Default command
CMD ["--help"]