# LeakyRepo - Developer Usage Guide

This guide explains how to use LeakyRepo as a developer to detect secrets before committing code.

## Quick Start (5 minutes)

### 1. Build the Tool

```bash
# Clone or download the repository
git clone https://github.com/lgboyce/leakyrepo.git
cd leakyrepo

# Build the binary
go build -o leakyrepo .

# (Optional) Add to your PATH
export PATH=$PATH:$(pwd)
# Or install globally:
sudo cp leakyrepo /usr/local/bin/
```

### 2. Initialize in Your Project

Navigate to your project directory:

```bash
cd /path/to/your/project

# Create default configuration
leakyrepo init
```

This creates `.leakyrepo.yml` with default rules for:
- AWS Access Keys
- Generic API Keys
- High-entropy strings

### 3. Test It Out

```bash
# Scan all staged files
leakyrepo scan

# Or scan specific files
leakyrepo scan config.env secrets.yaml
```

## Common Workflows

### Workflow 1: Manual Scanning Before Commit

**Before committing code, scan staged files:**

```bash
# Stage your files
git add .

# Scan staged files
leakyrepo scan

# If secrets are found, review and fix, then scan again
# Once clean, proceed with commit
git commit -m "Your commit message"
```

### Workflow 2: Automated Pre-commit Hook (Recommended)

**Install the pre-commit hook once, and it runs automatically:**

```bash
# Install the hook (do this once per repository)
leakyrepo install-hook

# Now every commit will be checked automatically
git add .
git commit -m "Your commit message"
# Hook runs automatically, blocks commit if secrets found
```

**What happens:**
- Every time you run `git commit`, the hook scans staged files
- If secrets are detected, the commit is blocked with an error message
- Fix the issues and try again

### Workflow 3: Scan Specific Files During Development

**Check files as you work on them:**

```bash
# Check a single file
leakyrepo scan config.env

# Check multiple files
leakyrepo scan .env.example .env.production secrets.yaml

# Check with explanations
leakyrepo scan config.env --explain
```

**Output example:**
```
⚠️  Found 1 potential secret(s):

🔒 [High] AWS Access Key found in config.env:42
   Match: AKIA****************
   Reason: Matched regex rule 'aws_access_key' (pattern: AKIA[0-9A-Z]{16})
```

### Workflow 4: Generate JSON Reports for CI/CD

**Generate machine-readable reports:**

```bash
# Scan and output to JSON
leakyrepo scan --json results.json

# View the report
cat results.json
```

**Example JSON output:**
```json
[
  {
    "file": "config.env",
    "line": 42,
    "rule_id": "aws_access_key",
    "severity": "high",
    "match": "AKIA****************"
  }
]
```

## Real-World Scenarios

### Scenario 1: New Developer Joining a Project

```bash
# 1. Clone the project
git clone https://github.com/company/project.git
cd project

# 2. Install LeakyRepo (if not already installed)
# (Follow installation steps above)

# 3. Initialize configuration (if not already present)
leakyrepo init

# 4. Install pre-commit hook
leakyrepo install-hook

# 5. Start working - the hook will protect you automatically
git checkout -b feature/new-feature
# ... make changes ...
git add .
git commit -m "Add new feature"
# Hook automatically checks for secrets!
```

### Scenario 2: Checking Existing Code

```bash
# Scan all files in a directory
find . -name "*.env" -o -name "*.yaml" | xargs leakyrepo scan

# Or scan everything (may be slow on large projects)
find . -type f -name "*.env" -exec leakyrepo scan {} \;
```

### Scenario 3: Updating Configuration

```bash
# Edit the configuration file
vi .leakyrepo.yml

# Add custom rules for your organization's secret patterns
# Example: Add Slack token detection
```

**Example custom rule:**
```yaml
rules:
  - id: slack_token
    description: "Slack Bot Token"
    severity: high
    pattern: 'xoxb-[0-9]{11}-[0-9]{11}-[a-zA-Z0-9]{24}'
    file_types:
      - .env
      - .yaml
```

### Scenario 4: Ignoring False Positives

**If you have files with false positives, add them to `.leakyrepoignore`:**

```bash
# Create or edit .leakyrepoignore
vi .leakyrepoignore
```

**Add patterns to ignore:**
```
# Ignore test fixtures
test/fixtures/*
*.test.env

# Ignore vendor directories
vendor/
node_modules/

# Ignore generated files
*.min.js
dist/
```

## Configuration Examples

### Minimal Configuration

```yaml
entropy_threshold: 4.5
rules: []
```

### Organization-Specific Configuration

```yaml
entropy_threshold: 4.5

rules:
  # AWS Access Keys
  - id: aws_access_key
    description: "AWS Access Key"
    severity: high
    pattern: 'AKIA[0-9A-Z]{16}'
    file_types: [.env, .yaml, .json, .py, .js]

  # Your company's API key format
  - id: company_api_key
    description: "Company API Key"
    severity: critical
    pattern: 'COMP_[a-zA-Z0-9]{32}'
    file_types: [.env, .yaml, .json]

  # Slack tokens
  - id: slack_token
    description: "Slack Bot Token"
    severity: high
    pattern: 'xoxb-[0-9]{11}-[0-9]{11}-[a-zA-Z0-9]{24}'
    file_types: [.env, .yaml, .json]

allowlist:
  files:
    - .leakyrepoignore
    - .git/**
    - vendor/**
  strings:
    - "EXAMPLE_KEY_NOT_REAL"
```

## Integration with CI/CD

### GitHub Actions Example

```yaml
name: Secret Detection

on: [push, pull_request]

jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.21
      
      - name: Build LeakyRepo
        run: |
          git clone https://github.com/lgboyce/leakyrepo.git
          cd leakyrepo
          go build -o ../leakyrepo .
          cd ..
      
      - name: Run LeakyRepo
        run: |
          ./leakyrepo init
          ./leakyrepo scan --json results.json || true
      
      - name: Upload results
        uses: actions/upload-artifact@v2
        if: always()
        with:
          name: leakyrepo-results
          path: results.json
```

### GitLab CI Example

```yaml
secret_detection:
  stage: test
  script:
    - go build -o leakyrepo /path/to/leakyrepo
    - ./leakyrepo init
    - ./leakyrepo scan --json results.json || true
  artifacts:
    when: always
    paths:
      - results.json
```

## Troubleshooting

### "Config file not found"

**Problem:** LeakyRepo can't find `.leakyrepo.yml`

**Solution:**
```bash
# Create default config
leakyrepo init
```

### "Not a git repository"

**Problem:** Running `leakyrepo scan` without arguments in a non-git directory

**Solution:**
```bash
# Specify files explicitly
leakyrepo scan file1 file2

# Or run from a git repository
cd /path/to/git/repo
leakyrepo scan
```

### Pre-commit hook not running

**Problem:** Hook installed but not executing

**Solution:**
```bash
# Check hook exists and is executable
ls -la .git/hooks/pre-commit

# Make executable if needed
chmod +x .git/hooks/pre-commit

# Test manually
.git/hooks/pre-commit
```

### False Positives

**Problem:** Legitimate strings being flagged

**Solution:**
1. Add to allowlist in `.leakyrepo.yml`:
   ```yaml
   allowlist:
     strings:
       - "your-legitimate-string"
   ```

2. Or add file pattern to `.leakyrepoignore`:
   ```
   test/fixtures/*
   ```

### Too many false positives from entropy detection

**Problem:** High-entropy string detection is too sensitive

**Solution:**
```yaml
# Increase the threshold in .leakyrepo.yml
entropy_threshold: 5.5  # Higher = less sensitive
```

## Best Practices

### 1. Install Hook Early

Install the pre-commit hook when you first clone a repository:
```bash
leakyrepo install-hook
```

### 2. Customize Rules for Your Stack

Edit `.leakyrepo.yml` to match your technology stack:
- Add rules for services you use (Slack, GitHub tokens, etc.)
- Adjust file types based on your languages

### 3. Use Allowlist for Test Data

Add test fixtures to `.leakyrepoignore`:
```
test/fixtures/
spec/fixtures/
__fixtures__/
```

### 4. Review Before Adding to Ignore

When you get a false positive:
1. First verify it's actually a false positive
2. Consider if the string should be in the code at all
3. Only then add to allowlist/ignore

### 5. Scan Before Committing

Make scanning part of your workflow:
```bash
# Good habit: scan before every commit
leakyrepo scan
git commit -m "Your message"
```

### 6. Share Configuration

Commit `.leakyrepo.yml` to your repository so all team members use the same rules:
```bash
git add .leakyrepo.yml
git commit -m "Add LeakyRepo configuration"
```

## Quick Reference

| Command | Description |
|---------|-------------|
| `leakyrepo init` | Create default `.leakyrepo.yml` |
| `leakyrepo scan` | Scan staged files |
| `leakyrepo scan file1 file2` | Scan specific files |
| `leakyrepo scan --json output.json` | Output JSON report |
| `leakyrepo scan --explain` | Show explanations |
| `leakyrepo install-hook` | Install pre-commit hook |

## Getting Help

- Check the README.md for detailed documentation
- Review `.leakyrepo.yml` for configuration options
- Check test files for examples of detection patterns

