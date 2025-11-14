#!/bin/bash
# Build and push Docker image for LeakyRepo
# Usage: ./scripts/docker-build.sh [version] [push]
# Example: ./scripts/docker-build.sh v1.1.0 push

set -e

VERSION=${1:-"latest"}
PUSH=${2:-""}
IMAGE_NAME="gittingsboyce/leakyrepo"

echo "🐳 Building Docker image: ${IMAGE_NAME}:${VERSION}"

# Check if buildx is available
if ! docker buildx version > /dev/null 2>&1; then
    echo "⚠️  buildx not available, using regular docker build (single platform)"
    docker build -t "${IMAGE_NAME}:${VERSION}" .
    if [ "$VERSION" != "latest" ]; then
        docker tag "${IMAGE_NAME}:${VERSION}" "${IMAGE_NAME}:latest"
    fi
else
    echo "Building multi-platform image (linux/amd64, linux/arm64)..."
    
    # Create builder if it doesn't exist
    docker buildx create --use --name leakyrepo-builder 2>/dev/null || docker buildx use leakyrepo-builder 2>/dev/null || true
    
    # Build for multiple platforms
    if [ "$PUSH" = "push" ]; then
        # Build and push multi-platform
        docker buildx build --platform linux/amd64,linux/arm64 \
            -t "${IMAGE_NAME}:${VERSION}" \
            -t "${IMAGE_NAME}:latest" \
            --push \
            .
    else
        # Build for local testing (amd64 only for speed)
        docker buildx build --platform linux/amd64 \
            -t "${IMAGE_NAME}:${VERSION}" \
            --load \
            .
        
        if [ "$VERSION" != "latest" ]; then
            docker tag "${IMAGE_NAME}:${VERSION}" "${IMAGE_NAME}:latest"
        fi
    fi
fi

echo "✅ Build complete!"
echo ""
echo "To test locally:"
echo "  docker run --rm ${IMAGE_NAME}:${VERSION} version"
echo ""
echo "To scan a directory:"
echo "  docker run --rm -v \$(pwd):/workspace ${IMAGE_NAME}:${VERSION} scan"
echo ""

# Push message (handled above if buildx is used)
if [ "$PUSH" = "push" ] && ! docker buildx version > /dev/null 2>&1; then
    echo "📤 Pushing to Docker Hub..."
    docker push "${IMAGE_NAME}:${VERSION}"
    if [ "$VERSION" != "latest" ]; then
        docker push "${IMAGE_NAME}:latest"
    fi
    echo "✅ Pushed to Docker Hub!"
    echo ""
    echo "Users can now use:"
    echo "  docker run --rm ${IMAGE_NAME}:${VERSION} scan"
fi

