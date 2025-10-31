# LeakyRepo - Quick Start Guide

Get up and running with LeakyRepo in under 2 minutes!

## Step 1: Build the Tool

```bash
git clone https://github.com/gittingsboyce/leakyrepo.git
cd leakyrepo
go build -o leakyrepo .
```

## Step 2: Use in Your Project

```bash
# Navigate to your project
cd /path/to/your/project

# Initialize configuration
leakyrepo init

# Install pre-commit hook (recommended)
leakyrepo install-hook

# That's it! Now every commit is automatically checked
```

## Step 3: Test It

```bash
# Create a test file with a secret
echo "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE" > test.env

# Try to commit it
git add test.env
git commit -m "Test commit"
# Commit will be blocked if secret detected!
```

## Daily Usage

Once set up, you don't need to do anything - the pre-commit hook runs automatically!

For manual scanning:
```bash
# Scan staged files
leakyrepo scan

# Scan specific files
leakyrepo scan config.env

# Get JSON output
leakyrepo scan --json results.json
```

That's it! For more details, see [USAGE.md](USAGE.md) or [README.md](README.md).

