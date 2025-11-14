# Docker Usage Guide

LeakyRepo is available as a Docker image for easy use in CI/CD pipelines and containerized environments.

## Quick Start

### Basic Usage

```bash
# Scan current directory (staged files)
docker run --rm --platform linux/amd64 -v $(pwd):/workspace -w /workspace gittingsboyce/leakyrepo:latest scan

# Show version
docker run --rm gittingsboyce/leakyrepo:latest version

# Show help
docker run --rm gittingsboyce/leakyrepo:latest --help
```

### Scan Specific Files

```bash
# Scan specific files
docker run --rm --platform linux/amd64 -v $(pwd):/workspace -w /workspace gittingsboyce/leakyrepo:latest scan file1.env file2.yaml

# Interactive mode
docker run --rm -it --platform linux/amd64 -v $(pwd):/workspace -w /workspace gittingsboyce/leakyrepo:latest scan -i
```

## CI/CD Integration

### GitHub Actions

```yaml
name: Secret Detection

on: [push, pull_request]

jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Run LeakyRepo
        run: |
          docker run --rm --platform linux/amd64 \
            -v ${{ github.workspace }}:/workspace \
            -w /workspace \
            gittingsboyce/leakyrepo:latest scan
```

### GitLab CI

```yaml
secret_detection:
  image: gittingsboyce/leakyrepo:latest
  script:
    - leakyrepo scan
  only:
    - merge_requests
    - main
```

### Jenkins

```groovy
pipeline {
    agent any
    stages {
        stage('Scan for Secrets') {
            steps {
                sh '''
                    docker run --rm \
                      -v ${WORKSPACE}:/workspace \
                      gittingsboyce/leakyrepo:latest scan
                '''
            }
        }
    }
}
```

## Building the Image

### Build Locally

```bash
# Build image
docker build -t gittingsboyce/leakyrepo:latest .

# Test it
docker run --rm gittingsboyce/leakyrepo:latest version
```

### Build Script

```bash
# Build for specific version
./scripts/docker-build.sh v1.1.0

# Build and push to Docker Hub
./scripts/docker-build.sh v1.1.0 push
```

## Publishing to Docker Hub

### Prerequisites

1. Create account on Docker Hub: https://hub.docker.com
2. Create repository: `gittingsboyce/leakyrepo`
3. Login: `docker login`

### Publish

```bash
# Build and tag
docker build -t gittingsboyce/leakyrepo:v1.1.0 .
docker tag gittingsboyce/leakyrepo:v1.1.0 gittingsboyce/leakyrepo:latest

# Push
docker push gittingsboyce/leakyrepo:v1.1.0
docker push gittingsboyce/leakyrepo:latest
```

Or use the build script:

```bash
./scripts/docker-build.sh v1.1.0 push
```

## Usage Examples

### Scan Git Repository

```bash
# Mount the repo and scan staged files
docker run --rm \
  -v $(pwd):/workspace \
  -w /workspace \
  gittingsboyce/leakyrepo:latest scan
```

### Generate JSON Report

```bash
docker run --rm \
  -v $(pwd):/workspace \
  gittingsboyce/leakyrepo:latest scan --json results.json

# Results will be in your current directory
cat results.json
```

### Initialize Configuration

```bash
docker run --rm \
  -v $(pwd):/workspace \
  gittingsboyce/leakyrepo:latest init

# .leakyrepo.yml will be created in your current directory
```

## Image Details

- **Base Image**: `alpine:latest` (small, ~20MB)
- **Includes**: Git (for staged file detection)
- **Working Directory**: `/workspace` (mount your code here)
- **Entrypoint**: `leakyrepo` (can override with commands)

## Troubleshooting

### Permission Errors

If you get permission errors, ensure Docker has access:

```bash
# On Linux, you might need to add user to docker group
sudo usermod -aG docker $USER
```

### Volume Mounting

Make sure to mount your code to `/workspace`:

```bash
# Correct
docker run --rm -v $(pwd):/workspace gittingsboyce/leakyrepo:latest scan

# Wrong (no volume mount)
docker run --rm gittingsboyce/leakyrepo:latest scan
```

### Git Not Found

The image includes git, but if you have issues:

```bash
# Verify git is available
docker run --rm gittingsboyce/leakyrepo:latest sh -c "which git"
```

