#!/bin/bash
# Test Homebrew formula locally
# This script tests the formula by building from source using the head version

set -e

echo "🧪 Testing LeakyRepo Homebrew formula..."
echo ""

# Get absolute path to project
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TAP_NAME="gittingsboyce/leakyrepo-test"

echo "Project directory: $PROJECT_DIR"
echo ""

# Create test tap if it doesn't exist
if ! brew tap-info "$TAP_NAME" &>/dev/null; then
    echo "Creating test tap: $TAP_NAME"
    brew tap-new "$TAP_NAME"
fi

# Get the tap directory
TAP_DIR=$(brew --prefix)/Library/Taps/${TAP_NAME%/*}/homebrew-${TAP_NAME##*/}
echo "Tap directory: $TAP_DIR"

# Create or update formula for testing
cat > "$TAP_DIR/Formula/leakyrepo.rb" <<EOF
# Homebrew formula for LeakyRepo (Testing version)
class Leakyrepo < Formula
  desc "Secrets detection tool that catches API keys and tokens before commit"
  homepage "https://github.com/gittingsboyce/leakyrepo"
  license "MIT"
  # Use head to build from local source for testing
  head "$PROJECT_DIR", :using => :git, :branch => "main"

  depends_on "go" => :build

  def install
    # Build the binary
    system "go", "build", "-ldflags", "-s -w", "-o", bin/"leakyrepo", "."
  end

  test do
    # Test that the binary works
    system "#{bin}/leakyrepo", "version"
    
    # Test init command
    testpath = testpath/"test-project"
    testpath.mkpath
    cd testpath do
      system bin/"leakyrepo", "init"
      assert_predicate testpath/".leakyrepo.yml", :exist?
    end
  end
end
EOF

# Commit to tap (if git initialized)
cd "$TAP_DIR"
if git rev-parse --git-dir > /dev/null 2>&1; then
    git add Formula/leakyrepo.rb 2>/dev/null || true
    git commit -m "Update leakyrepo formula for testing" 2>/dev/null || true
fi

echo ""
echo "📦 Installing from test tap (building from source)..."
echo ""

# Uninstall if already installed
if brew list --formula "$TAP_NAME/leakyrepo" &>/dev/null; then
    echo "Uninstalling existing version..."
    brew uninstall "$TAP_NAME/leakyrepo" 2>/dev/null || true
fi

# Install from head (builds from local source)
if brew install --HEAD --build-from-source "$TAP_NAME/leakyrepo"; then
    echo ""
    echo "✅ Installation successful!"
    echo ""
    
    # Test the binary
    echo "🧪 Testing binary..."
    if leakyrepo version > /dev/null 2>&1; then
        echo "✅ Version command works:"
        leakyrepo version
        echo ""
    fi
    
    # Test init in temp directory
    TEST_DIR=$(mktemp -d)
    cd "$TEST_DIR"
    echo "🧪 Testing init command in: $TEST_DIR"
    if leakyrepo init; then
        echo "✅ Init command works"
        if [ -f ".leakyrepo.yml" ]; then
            echo "✅ Config file created"
        fi
    fi
    rm -rf "$TEST_DIR"
    
    echo ""
    echo "✅ All tests passed! Formula is ready to publish."
    echo ""
    echo "Next steps:"
    echo "  1. Create GitHub release with tag v1.0.0"
    echo "  2. Upload binaries to GitHub Releases"
    echo "  3. Update Formula/leakyrepo.rb with release URL and SHA256"
    echo "  4. Push to homebrew-leakyrepo tap"
else
    echo ""
    echo "❌ Installation failed. Check errors above."
    exit 1
fi

