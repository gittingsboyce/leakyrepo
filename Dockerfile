# LeakyRepo Docker Image
# A minimal, secure container for running LeakyRepo in CI/CD pipelines

# Use Alpine Linux for minimal size (~5MB base)
FROM alpine:latest

# Metadata for maintainability
LABEL maintainer="LeakyRepo Team"
LABEL description="LeakyRepo secrets detection tool in a container"

# Install required packages:
# - ca-certificates for HTTPS operations
# - git for repository scanning functionality
# - BusyBox provides basic Unix utilities needed by LeakyRepo
RUN apk --no-cache add ca-certificates git

# Copy the compiled binary into the container
# IMPORTANT: Build the Linux binary first using:
#   GOOS=linux GOARCH=amd64 go build -o leakyrepo .
# Or use the build-release.sh script to create all platform binaries
COPY leakyrepo /usr/local/bin/leakyrepo

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

