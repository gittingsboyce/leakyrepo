# LeakyRepo

A secrets detection tool that catches API keys, tokens, and credentials before they're committed to version control.

## Features

- **CLI Scanner**: Scan staged files or specific files
- **Interactive Mode**: Easily ignore false positives interactively (`leakyrepo scan -i`)
- **Regex & Entropy Detection**: Custom regex patterns + Shannon entropy for high-entropy strings
- **Pre-commit Hook**: Automatically block commits with secrets
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

## Documentation

- **[USAGE.md](USAGE.md)**: Comprehensive usage guide with examples
- **[QUICKSTART.md](QUICKSTART.md)**: 2-minute setup guide

## Testing

```bash
go test ./...
```

## License

See [LICENSE](LICENSE) file for details.
