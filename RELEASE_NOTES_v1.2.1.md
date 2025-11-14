# LeakyRepo v1.2.1 Release Notes

## 🎯 Improved Entropy Detection - Reduced False Positives

This release focuses on making entropy-based secret detection more prescriptive and reducing false positives, especially in code files and binary assets.

### ✨ Key Improvements

#### **More Prescriptive Entropy Detection**
- **Increased entropy threshold**: Default threshold raised from 4.5 to 5.5 to reduce false positives
- **Longer minimum length**: Now requires at least 16 characters (up from 8) for entropy detection
- **Better code pattern filtering**: Automatically filters out common false positives:
  - Template strings (`${variable}`)
  - JSX/HTML tags (`</tag>`)
  - CSS classes with brackets (`w-[...]`)
  - Code with multiple parentheses
  - Common variable prefixes/suffixes

#### **Binary File Detection**
- Automatically detects and skips binary files (PNG, JPEG, GIF)
- Detects files with high non-printable character content
- Prevents false positives from image files and other binary assets

#### **Smarter Tokenization**
- Focuses on key-value patterns (`key=value`, `key:value`) common in config files
- Less aggressive splitting to avoid code constructs
- Better filtering during tokenization phase

#### **Default File Exclusions**
- Added common binary/image file extensions to default allowlist:
  - Images: `*.png`, `*.jpg`, `*.jpeg`, `*.gif`, `*.ico`, `*.svg`, `*.webp`
  - Archives: `*.pdf`, `*.zip`, `*.tar`, `*.gz`, `*.bz2`
  - Binaries: `*.exe`, `*.dll`, `*.so`, `*.dylib`
  - Fonts: `*.woff`, `*.woff2`, `*.ttf`, `*.eot`

### 🐛 Bug Fixes

- Fixed Git ownership error in Docker for `/workspace` mount path (in addition to `/github/workspace`)
- Improved entrypoint script to handle both standard GitHub Actions and custom workflows

### 📦 Installation

**Homebrew:**
```bash
brew upgrade leakyrepo
```

**Docker:**
```bash
docker pull gittingsboyce/leakyrepo:latest
```

**From Source:**
```bash
git clone https://github.com/gittingsboyce/leakyrepo.git
cd leakyrepo
git checkout v1.2.1
go build -o leakyrepo .
```

### 🔄 Migration Notes

If you have a custom `.leakyrepo.yml` with a lower `entropy_threshold`, you may want to increase it to 5.5 to match the new default and reduce false positives.

### 📝 Full Changelog

- Increased default entropy threshold from 4.5 to 5.5
- Increased minimum string length for entropy detection from 8 to 16 characters
- Added binary file detection (PNG, JPEG, GIF signatures)
- Added code pattern filtering for template strings, JSX, CSS classes
- Improved tokenization to focus on key-value patterns
- Added default file exclusions for common binary formats
- Fixed Docker entrypoint to handle `/workspace` mount path

---

**Thank you for using LeakyRepo!** 🚀

For issues or feature requests, please visit: https://github.com/gittingsboyce/leakyrepo/issues

