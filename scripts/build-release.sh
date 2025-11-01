#!/bin/bash
# Build release binaries for LeakyRepo
# Usage: ./scripts/build-release.sh [version]
# Example: ./scripts/build-release.sh v1.0.0

set -e

VERSION=${1:-"v1.0.0"}
BINARY_NAME="leakyrepo"
BUILD_DIR="build"
RELEASE_DIR="${BUILD_DIR}/release"

# Clean previous builds
rm -rf ${BUILD_DIR}
mkdir -p ${RELEASE_DIR}

# Get the git commit hash
COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build info
LDFLAGS="-s -w -X 'main.version=${VERSION}' -X 'main.commit=${COMMIT_HASH}' -X 'main.date=${BUILD_DATE}'"

echo "Building LeakyRepo ${VERSION}..."
echo "Commit: ${COMMIT_HASH}"
echo "Date: ${BUILD_DATE}"
echo ""

# Build for multiple platforms
PLATFORMS=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "windows/arm64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
    PLATFORM_SPLIT=(${PLATFORM//\// })
    GOOS=${PLATFORM_SPLIT[0]}
    GOARCH=${PLATFORM_SPLIT[1]}
    OUTPUT_NAME="${BINARY_NAME}"
    
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="${BINARY_NAME}.exe"
    fi
    
    OUTPUT_PATH="${RELEASE_DIR}/${BINARY_NAME}-${VERSION}-${GOOS}-${GOARCH}/${OUTPUT_NAME}"
    ARCHIVE_NAME="${BINARY_NAME}-${VERSION}-${GOOS}-${GOARCH}"
    
    echo "Building for ${GOOS}/${GOARCH}..."
    
    env GOOS=${GOOS} GOARCH=${GOARCH} go build \
        -ldflags "${LDFLAGS}" \
        -o "${OUTPUT_PATH}" \
        .
    
    # Create archive
    cd ${RELEASE_DIR}
    if [ "$GOOS" = "windows" ]; then
        zip -r "${ARCHIVE_NAME}.zip" "${ARCHIVE_NAME}" > /dev/null
    else
        tar -czf "${ARCHIVE_NAME}.tar.gz" "${ARCHIVE_NAME}" > /dev/null
    fi
    cd - > /dev/null
    
    echo "  ✓ Created ${RELEASE_DIR}/${ARCHIVE_NAME}.tar.gz or .zip"
done

# Create checksums
echo ""
echo "Creating checksums..."
cd ${RELEASE_DIR}
sha256sum *.tar.gz *.zip > checksums.txt 2>/dev/null || shasum -a 256 *.tar.gz *.zip > checksums.txt
cd - > /dev/null

echo ""
echo "✅ Build complete!"
echo ""
echo "Release files are in: ${RELEASE_DIR}/"
echo ""
echo "To create a GitHub release:"
echo "  1. Create a release tag: git tag ${VERSION}"
echo "  2. Push the tag: git push origin ${VERSION}"
echo "  3. Upload files from ${RELEASE_DIR}/ to GitHub Releases"
echo ""
echo "Then update Formula/leakyrepo.rb with the new SHA256 from checksums.txt"

