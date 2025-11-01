# Homebrew Installation Guide

This guide explains how to set up LeakyRepo for distribution via Homebrew.

## Option 1: Personal Tap (Easiest - Recommended)

A personal tap allows you to host the formula in your own GitHub repository.

### Step 1: Create a Tap Repository

Create a new GitHub repository called `homebrew-leakyrepo`:

```bash
# Create the tap repository on GitHub (via web UI)
# Repository name: homebrew-leakyrepo
# Description: Homebrew tap for LeakyRepo
# Make it public
```

### Step 2: Clone and Set Up the Tap

```bash
# Clone the tap repository
git clone https://github.com/gittingsboyce/homebrew-leakyrepo.git
cd homebrew-leakyrepo
```

### Step 3: Create Release and Update Formula

1. **Create a GitHub release:**
   ```bash
   # Tag the release
   git tag v1.0.0
   git push origin v1.0.0
   
   # Build release binaries
   ./scripts/build-release.sh v1.0.0
   ```

2. **Upload to GitHub Releases:**
   - Go to your repository's Releases page
   - Create a new release with tag `v1.0.0`
   - Upload the `leakyrepo-v1.0.0-darwin-amd64.tar.gz` and `leakyrepo-v1.0.0-darwin-arm64.tar.gz` files

3. **Get SHA256 checksums:**
   ```bash
   # From the build output
   cat build/release/checksums.txt
   # Or download and check:
   shasum -a 256 leakyrepo-v1.0.0-darwin-amd64.tar.gz
   shasum -a 256 leakyrepo-v1.0.0-darwin-arm64.tar.gz
   ```

4. **Update the formula:**
   ```bash
   # Copy the formula to your tap
   cp Formula/leakyrepo.rb /path/to/homebrew-leakyrepo/Formula/leakyrepo.rb
   ```

5. **Edit the formula:**
   - Update the `url` to point to your GitHub release
   - Update the `sha256` values for both architectures
   - Update the version

6. **Commit and push:**
   ```bash
   git add Formula/leakyrepo.rb
   git commit -m "Add LeakyRepo formula v1.0.0"
   git push origin main
   ```

### Step 4: Users Install

Users can now install with:

```bash
brew tap gittingsboyce/leakyrepo
brew install leakyrepo
```

Or in one command:
```bash
brew install gittingsboyce/leakyrepo/leakyrepo
```

## Option 2: Homebrew Core (Official)

To get into the official Homebrew repository:

1. **Meet Requirements:**
   - At least 50 stars on GitHub
   - Not a tool that's specific to one project
   - Useful to a general audience

2. **Submit PR:**
   - Fork homebrew-core
   - Add your formula
   - Submit pull request

3. **Review Process:**
   - Homebrew maintainers review
   - May take several weeks

For Phase 1, Option 1 (personal tap) is recommended.

## Updating the Formula for New Releases

1. Build new release:
   ```bash
   ./scripts/build-release.sh v1.1.0
   ```

2. Create GitHub release and upload binaries

3. Update formula:
   - Update `version` and `url`
   - Update `sha256` checksums

4. Commit and push to tap:
   ```bash
   git add Formula/leakyrepo.rb
   git commit -m "Update to v1.1.0"
   git push
   ```

5. Users update with:
   ```bash
   brew upgrade leakyrepo
   ```

## Testing the Formula Locally

Before publishing, test locally using the provided test script:

```bash
# Run the test script (builds from local source)
./scripts/test-homebrew.sh
```

This script will:
- Create a test tap
- Build the formula from your local source (HEAD)
- Install it
- Run all tests to verify it works

**Manual testing:**

```bash
# Test from test tap (builds from local source)
brew install --HEAD --build-from-source gittingsboyce/leakyrepo-test/leakyrepo

# Test that it works
leakyrepo version
leakyrepo init
```

**Or test from a local formula file (for official releases):**

```bash
# Test installation from local formula
brew install --build-from-source Formula/leakyrepo.rb

# Or test from tap (if already published)
brew tap gittingsboyce/leakyrepo
brew install --build-from-source gittingsboyce/leakyrepo/leakyrepo
```

## Formula Template

The formula file is located at `Formula/leakyrepo.rb`. Key sections:

- **Description**: What the tool does
- **Homepage**: GitHub repository URL
- **URL**: Link to release tarball
- **SHA256**: Checksum for security
- **License**: MIT
- **Dependencies**: Go (for building from source)
- **Install**: Build command
- **Test**: Verification tests

## Troubleshooting

### Formula won't install

- Check that the URL is accessible
- Verify SHA256 checksums match
- Ensure Go version requirements are met

### Build fails

- Check Go version: `go version`
- Ensure dependencies are available: `go mod download`

### Test fails

- Verify the binary works: `./leakyrepo version`
- Check that test paths are correct

