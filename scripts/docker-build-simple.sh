#!/bin/bash
# Simple Docker build script (single platform, faster for testing)
# Usage: ./scripts/docker-build-simple.sh [version] [platform]
# Example: ./scripts/docker-build-simple.sh v1.1.0 linux/amd64

set -e

VERSION=${1:-"latest"}
PLATFORM=${2:-"linux/amd64"}
IMAGE_NAME="gittingsboyce/leakyrepo"

echo "🐳 Building Docker image: ${IMAGE_NAME}:${VERSION} for ${PLATFORM}"

# Build for specific platform
docker buildx build --platform ${PLATFORM} \
    -t "${IMAGE_NAME}:${VERSION}" \
    --load \
    .

echo "✅ Build complete!"
echo ""
echo "To test:"
echo "  docker run --rm --platform ${PLATFORM} ${IMAGE_NAME}:${VERSION} version"
echo ""
echo "To push:"
echo "  docker push ${IMAGE_NAME}:${VERSION}"

