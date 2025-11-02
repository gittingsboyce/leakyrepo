# LeakyRepo

A secrets detection tool that catches API keys, tokens, and credentials before they're committed to version control.

You can learn more about our project and its future at https://leakyrepo.com/

## Features

- **CLI Scanner**: Scan staged files, all tracked files, or specific files
- **Interactive Mode**: Easily ignore false positives interactively (`leakyrepo scan -i`)
- **Regex & Entropy Detection**: Custom regex patterns + Shannon entropy for high-entropy strings
- **Pre-commit Hook**: Automatically block commits with secrets
- **CI/CD Integration**: Docker and GitHub Actions support
- **JSON Output**: Machine-readable reports for CI/CD
- **Configurable**: Customize rules, thresholds, and ignore patterns

## Installation

### Homebrew (Recommended)

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

# Scan staged files (interactive mode recommended)
leakyrepo scan -i

# Install pre-commit hook (recommended)
leakyrepo install-hook
```

> **📖 See [USAGE.md](USAGE.md) for detailed instructions and examples**

## Commands

| Command | Description |
|---------|-------------|
| `leakyrepo scan [files...]` | Scan files (or staged files if none specified) |
| `leakyrepo scan --all` | Scan all tracked files in the repository (useful for CI/CD) |
| `leakyrepo scan -i` | **Interactive mode** - prompt to ignore false positives |
| `leakyrepo scan --json <file>` | Output JSON report |
| `leakyrepo scan --explain` | Show explanation for each detection |
| `leakyrepo ignore <file>` | Quick command to ignore a file or pattern |
| `leakyrepo init` | Create default `.leakyrepo.yml` |
| `leakyrepo install-hook` | Install Git pre-commit hook |

## Configuration

### `.leakyrepo.yml`

```yaml
entropy_threshold: 4.5

rules:
  - id: aws_access_key
    description: "AWS Access Key"
    severity: high
    pattern: 'AKIA[0-9A-Z]{16}'
    file_types: [.env, .yaml, .json, .py, .js]

allowlist:
  files: [.leakyrepoignore, .git/**]
  strings: []
```

### `.leakyrepoignore`

```
node_modules/
vendor/
*.lock
dist/
```

## Example Output

```
⚠️  Found 1 potential secret(s):

🔒 [High] AWS Access Key found in config.env:42
   Match: AKIA****************
```

### Interactive Mode (Recommended)

When false positives are detected, use interactive mode to quickly ignore them:

```bash
leakyrepo scan -i
```

Interactive mode:
- Shows all findings grouped by file
- Prompts to ignore false positives
- Automatically updates `.leakyrepoignore`
- Re-scans to verify all issues resolved
- Handles multiple secrets efficiently

## How It Works

- **Regex Detection**: Matches known secret patterns (AWS keys, API keys, etc.)
- **Entropy Detection**: Detects high-entropy strings using Shannon entropy

Default entropy threshold: 4.5 (configurable)

## CI/CD Integration

LeakyRepo can be easily integrated into your CI/CD pipelines using Docker or GitHub Actions.

### GitHub Actions

Use the LeakyRepo GitHub Action to automatically scan your repository:

```yaml
name: LeakyRepo Secret Scan

on:
  pull_request:
    types: [opened, synchronize, reopened]
  push:
    branches: [main]

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

**What gets scanned?**
- By default, `scan --all` scans all tracked files in your repository
- Configuration files (`.leakyrepo.yml` and `.leakyrepoignore`) are automatically picked up
- The Action exits with code 1 if secrets are found, failing the CI job

**Customizing the scan:**
```yaml
- Scan with JSON output
  args: scan --all --json secrets-report.json

- Scan with explanations
  args: scan --all --explain

- Scan only staged files
  args: scan
```

See [`.github/workflows/leakyrepo-scan.yaml`](.github/workflows/leakyrepo-scan.yaml) for a complete example.

### Docker

Run LeakyRepo in any Docker-compatible CI/CD environment:

```bash
# First, build the Linux binary
GOOS=linux GOARCH=amd64 go build -o leakyrepo .

# Build the Docker image
docker build -t leakyrepo .

# Run scan (config files from host are mounted automatically)
docker run --rm -v $(pwd):/workspace -w /workspace leakyrepo scan --all

# Run with JSON output
docker run --rm -v $(pwd):/workspace -w /workspace leakyrepo scan --all --json report.json
```

**Exit Codes:**
- `0`: No secrets found
- `1`: Secrets detected (CI job should fail)

## Documentation

- **[USAGE.md](USAGE.md)**: Comprehensive usage guide with examples
- **[QUICKSTART.md](QUICKSTART.md)**: 2-minute setup guide

## Testing

```bash
go test ./...
```

## License

See [LICENSE](LICENSE) file for details.
