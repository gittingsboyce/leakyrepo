# Publishing LeakyRepo to Homebrew

Step-by-step guide to publish LeakyRepo as a Homebrew package.

## Prerequisites

- [ ] GitHub repository exists and is public
- [ ] Go 1.21+ installed
- [ ] Homebrew installed (for testing)
- [ ] Git configured

## Step 1: Create GitHub Release

### 1.1 Build Release Binaries

```bash
# Build binaries for all platforms
./scripts/build-release.sh v1.0.0
```

This creates binaries in `build/release/`:
- `leakyrepo-v1.0.0-darwin-amd64.tar.gz`
- `leakyrepo-v1.0.0-darwin-arm64.tar.gz`
- `leakyrepo-v1.0.0-linux-amd64.tar.gz`
- `leakyrepo-v1.0.0-linux-arm64.tar.gz`
- `leakyrepo-v1.0.0-windows-amd64.zip`
- `leakyrepo-v1.0.0-windows-arm64.zip`
- `checksums.txt`

### 1.2 Create Git Tag

```bash
# Tag the release
git tag v1.0.0
git push origin v1.0.0
```

### 1.3 Create GitHub Release

1. Go to your GitHub repository
2. Click "Releases" → "Create a new release"
3. Choose tag: `v1.0.0`
4. Release title: `v1.0.0`
5. Description: Release notes

### 1.4 Upload Binaries

Upload these files to the GitHub release:
- `build/release/leakyrepo-v1.0.0-darwin-amd64.tar.gz`
- `build/release/leakyrepo-v1.0.0-darwin-arm64.tar.gz`
- (Optional: Linux and Windows binaries for other users)

**Important:** Get the SHA256 checksums from `build/release/checksums.txt`:
```bash
cat build/release/checksums.txt
```

Save these - you'll need them for the formula.

## Step 2: Create Homebrew Tap Repository

### 2.1 Create GitHub Repository

1. Go to GitHub and create a new repository
2. Repository name: `homebrew-leakyrepo` (must start with `homebrew-`)
3. Make it public
4. Description: "Homebrew tap for LeakyRepo"
5. Don't initialize with README

### 2.2 Clone and Set Up Tap

```bash
# Clone the tap repository
git clone https://github.com/gittingsboyce/homebrew-leakyrepo.git
cd homebrew-leakyrepo

# Create Formula directory
mkdir -p Formula
```

### 2.3 Update Formula with Release Info

Copy the formula from your main repo:

```bash
# Copy formula from main repo
cp /path/to/leakyrepo/Formula/leakyrepo.rb Formula/leakyrepo.rb
```

Then edit `Formula/leakyrepo.rb` to update:

1. **URL**: Point to your GitHub release
   ```ruby
   url "https://github.com/gittingsboyce/leakyrepo/archive/v1.0.0.tar.gz"
   ```

2. **SHA256**: Update with checksums from Step 1.4
   ```ruby
   sha256 "abc123..."  # Replace with actual SHA256 from checksums.txt
   ```

3. **Version**: Make sure version matches
   ```ruby
   # Version is in the URL, so it's already set
   ```

**Example updated formula:**

```ruby
class Leakyrepo < Formula
  desc "Secrets detection tool that catches API keys and tokens before commit"
   homepage "https://github.com/gittingsboyce/leakyrepo"
   url "https://github.com/gittingsboyce/leakyrepo/archive/v1.0.0.tar.gz"
  sha256 "abc123def456..."  # From checksums.txt
  license "MIT"
   head "https://github.com/gittingsboyce/leakyrepo.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags", "-s -w", "-o", bin/"leakyrepo", "."
  end

  test do
    system bin/"leakyrepo", "version"
    
    testpath = testpath/"test-project"
    testpath.mkpath
    cd testpath do
      system bin/"leakyrepo", "init"
      assert_predicate testpath/".leakyrepo.yml", :exist?
    end
  end
end
```

### 2.4 Commit and Push

```bash
cd homebrew-leakyrepo
git add Formula/leakyrepo.rb
git commit -m "Add LeakyRepo formula v1.0.0"
git push origin main
```

## Step 3: Test Installation

Before announcing, test that users can install:

```bash
# Tap your repository
brew tap gittingsboyce/leakyrepo

# Install
brew install leakyrepo

# Verify it works
leakyrepo version
leakyrepo init
```

## Step 4: Update Documentation

Update your main repository README:

```markdown
## Installation

### Homebrew (Recommended)

```bash
brew tap gittingsboyce/leakyrepo
brew install leakyrepo
```
```

## Step 5: Announce!

Share the installation command:

```bash
brew tap gittingsboyce/leakyrepo && brew install leakyrepo
```

## Updating for New Versions

When releasing v1.1.0:

1. **Build new release:**
   ```bash
   ./scripts/build-release.sh v1.1.0
   ```

2. **Create GitHub release** with new binaries

3. **Update formula in tap:**
   ```bash
   cd homebrew-leakyrepo
   # Edit Formula/leakyrepo.rb
   # Update url, sha256, version
   git commit -m "Update to v1.1.0"
   git push
   ```

4. **Users update with:**
   ```bash
   brew upgrade leakyrepo
   ```

## Troubleshooting

### Formula won't install

- Check URL is accessible: `curl -I <url>`
- Verify SHA256 matches: `shasum -a 256 <file>`
- Check Go version requirement

### Build fails

- Ensure dependencies are available
- Check Go version: `go version`

### Test fails

- Verify binary works locally: `./leakyrepo version`
- Check test paths in formula

## Checklist

Before publishing:

- [ ] Build release binaries with `./scripts/build-release.sh`
- [ ] Test binaries locally
- [ ] Create GitHub release tag
- [ ] Upload binaries to GitHub Releases
- [ ] Get SHA256 checksums from `checksums.txt`
- [ ] Create `homebrew-leakyrepo` repository
- [ ] Update formula with correct URL and SHA256
- [ ] Commit and push formula
- [ ] Test installation: `brew tap gittingsboyce/leakyrepo && brew install leakyrepo`
- [ ] Verify binary works: `leakyrepo version`
- [ ] Update main README with Homebrew install instructions

