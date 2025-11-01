# Step 2: Create Homebrew Tap - Quick Guide

After you've published your GitHub release, follow these steps:

## 1. Create GitHub Repository

1. Go to: https://github.com/new
2. Repository name: `homebrew-leakyrepo` (must start with `homebrew-`)
3. Make it **public**
4. Description: "Homebrew tap for LeakyRepo"
5. **Don't** initialize with README, .gitignore, or license
6. Click "Create repository"

## 2. Clone and Set Up

```bash
# Clone the empty repository
git clone https://github.com/gittingsboyce/homebrew-leakyrepo.git
cd homebrew-leakyrepo

# Create Formula directory
mkdir -p Formula
```

## 3. Copy and Update Formula

```bash
# Copy formula from your main repo
cp /Users/lgboyce/Desktop/projects/leakyrepo/Formula/leakyrepo.rb Formula/leakyrepo.rb
```

The formula is already updated with:
- ✅ Correct URL: `https://github.com/gittingsboyce/leakyrepo/archive/v1.0.0.tar.gz`
- ✅ SHA256: `4698027617a50a4abcdf096ad9b6d3477ac7b910fc04939ba1eeca7a423e03b9`

## 4. Commit and Push

```bash
cd homebrew-leakyrepo
git add Formula/leakyrepo.rb
git commit -m "Add LeakyRepo formula v1.0.0"
git push origin main
```

## 5. Test Installation

```bash
# Tap your repository
brew tap gittingsboyce/leakyrepo

# Install
brew install leakyrepo

# Verify it works
leakyrepo version
```

## Done! 🎉

Users can now install with:
```bash
brew tap gittingsboyce/leakyrepo
brew install leakyrepo
```

