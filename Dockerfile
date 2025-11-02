# LeakyRepo Docker Image
# A minimal, secure container for running LeakyRepo in CI/CD pipelines

# Build stage - compile the binary
FROM golang:1.21-alpine AS builder

# Install git for go modules (needed for some dependencies)
RUN apk --no-cache add git

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# Use -ldflags to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o leakyrepo \
    .

# Final stage - minimal runtime image
FROM alpine:latest

# Metadata for maintainability
LABEL maintainer="LeakyRepo Team"
LABEL description="LeakyRepo secrets detection tool in a container"

# Install required packages:
# - ca-certificates for HTTPS operations
# - git for repository scanning functionality
# - BusyBox provides basic Unix utilities needed by LeakyRepo
RUN apk --no-cache add ca-certificates git

# Copy the compiled binary from builder stage
COPY --from=builder /build/leakyrepo /usr/local/bin/leakyrepo

# Ensure the binary is executable
RUN chmod +x /usr/local/bin/leakyrepo

# Set working directory to allow mounting project files
WORKDIR /workspace

# Note: We run as root for Docker container actions compatibility
# GitHub Actions may mount volumes that require root permissions

# Set the entrypoint to leakyrepo
# All arguments passed to docker run will be forwarded to the binary
ENTRYPOINT ["/usr/local/bin/leakyrepo"]

# Default command if none provided
# Users can override this with docker run arguments
CMD ["scan", "--help"]

