# LeakyRepo v1.0.0

Initial release of LeakyRepo - a secrets detection tool that catches API keys, tokens, and credentials before they're committed to version control.

## Features

### CLI Scanner
- Scan files staged for git commit or any files specified in command-line arguments
- Human-readable output with color highlighting
- JSON output format for CI/CD integration
- `--explain` flag to show why secrets were detected

### Detection Methods
- **Regex-based Detection**: Customizable patterns for known secret formats (AWS keys, API keys, etc.)
- **Entropy-based Detection**: Uses Shannon entropy to detect high-entropy strings that may be secrets

### Pre-commit Hook
- Automatically install Git pre-commit hook with `leakyrepo install-hook`
- Blocks commits if secrets are detected in staged files

### Configuration
- `.leakyrepo.yml` configuration file with customizable rules
- Support for file type-specific rules
- Configurable entropy threshold
- Allowlist support for ignoring false positives

### File Ignoring
- `.leakyrepoignore` file support for excluding files/patterns

## Installation

### Homebrew
```bash
brew tap gittingsboyce/leakyrepo
brew install leakyrepo
```

### Build from Source
```bash
git clone https://github.com/gittingsboyce/leakyrepo.git
cd leakyrepo
go build -o leakyrepo .
```

## Quick Start

```bash
# Initialize configuration
leakyrepo init

# Scan staged files
leakyrepo scan

# Install pre-commit hook
leakyrepo install-hook
```

## Documentation

- [README.md](README.md) - Overview and installation
- [USAGE.md](USAGE.md) - Comprehensive usage guide
- [QUICKSTART.md](QUICKSTART.md) - 2-minute setup guide

## Default Detection Rules

- AWS Access Key (AKIA...)
- Generic API Key patterns
- High-entropy string detection (threshold: 4.5)

## Supported Platforms

- macOS (amd64, arm64)
- Linux (amd64, arm64)
- Windows (amd64, arm64)

## Commands

- `leakyrepo scan [files...]` - Scan files for secrets
- `leakyrepo scan --json <file>` - Output JSON report
- `leakyrepo scan --explain` - Show explanations
- `leakyrepo init` - Create default configuration
- `leakyrepo install-hook` - Install Git pre-commit hook
- `leakyrepo version` - Show version information

