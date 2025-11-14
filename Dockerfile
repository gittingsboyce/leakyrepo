# Multi-stage build for LeakyRepo
# Stage 1: Build
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG BUILDPLATFORM

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary for the target platform
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-s -w" \
    -o leakyrepo .

# Stage 2: Runtime
FROM --platform=$TARGETPLATFORM alpine:latest

ARG TARGETPLATFORM=linux/amd64

# Install git (needed for pre-commit hooks and staged file detection)
# Git will be installed for the correct platform automatically
RUN apk --no-cache add git ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/leakyrepo /usr/local/bin/leakyrepo

# Copy entrypoint script (for GitHub Actions argument handling)
COPY entrypoint.sh /usr/local/bin/entrypoint.sh

# Make binaries executable
RUN chmod +x /usr/local/bin/leakyrepo /usr/local/bin/entrypoint.sh

# Set working directory to /workspace (where users will mount their code)
WORKDIR /workspace

# Use entrypoint script to handle GitHub Actions argument parsing
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
