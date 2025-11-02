# LeakyRepo v1.2.0

Release of Phase 2 features - Docker and GitHub Actions CI/CD integration.

## 🎉 New Features

### CI/CD Integration
- **Docker Support**: Run LeakyRepo in any Docker-compatible CI/CD environment
- **GitHub Actions**: Drop-in GitHub Action for easy CI/CD integration
- **Full Repository Scanning**: New `--all` flag to scan all tracked files (perfect for CI)

### Enhanced Scanning
- **`scan --all` flag**: Scan entire repository, not just staged files
- **Improved path handling**: Better support for absolute vs relative paths in ignore patterns
- **Better directory pattern matching**: Enhanced `.leakyrepoignore` directory pattern support

## 🐛 Bug Fixes

- Fixed directory ignore patterns (`scripts/**`, etc.) to work with absolute file paths
- Improved scanner to handle both relative and absolute path matching correctly

## 📚 Documentation

- Added comprehensive CI/CD integration guide
- Docker build instructions with Linux binary steps
- GitHub Actions example workflow
- Updated command reference with `--all` flag

## 🚀 Quick Start

### GitHub Actions

```yaml
name: LeakyRepo Secret Scan
on: [pull_request, push]
jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - name: Run LeakyRepo
        uses: ./github-action
        with:
          args: scan --all
```

### Docker

```bash
# Build Linux binary
GOOS=linux GOARCH=amd64 go build -o leakyrepo .

# Build Docker image
docker build -t leakyrepo .

# Run scan
docker run --rm -v $(pwd):/workspace -w /workspace leakyrepo scan --all
```

## 📦 Installation

### Homebrew
```bash
brew tap gittingsboyce/leakyrepo
brew upgrade leakyrepo
```

### Build from Source
```bash
git clone https://github.com/gittingsboyce/leakyrepo.git
cd leakyrepo
git checkout v1.2.0
go build -o leakyrepo .
```

## 🔄 Upgrading from v1.1.0

The main changes are new features, so upgrading should be seamless:
- Add `--all` flag when you need full repository scanning in CI
- Use new Docker and GitHub Actions options for CI/CD integration
- Ignore patterns will work better with directory patterns

## 📝 Full Changelog

### Added
- Dockerfile for containerized LeakyRepo
- GitHub Action metadata and example workflow
- `scan --all` flag for repository-wide scanning
- `git.GetAllTrackedFiles()` helper function
- Comprehensive CI/CD documentation

### Changed
- Improved directory ignore pattern matching
- Enhanced scanner path resolution logic
- Better Docker build documentation

### Fixed
- Directory patterns in `.leakyrepoignore` now work correctly
- Absolute file path handling in scanner

## 🏗️ Files Changed

**New Files:**
- `Dockerfile`
- `github-action/action.yml`
- `.github/workflows/leakyrepo-scan.yaml`

**Modified Files:**
- `cmd/scan.go` (added `--all` flag)
- `git/git.go` (added `GetAllTrackedFiles()`)
- `scanner/scanner.go` (improved path handling)
- `README.md` (CI/CD documentation)

## 🙏 Thank You

Thank you for using LeakyRepo! Please report any issues at https://github.com/gittingsboyce/leakyrepo/issues

