# BkpDir

BkpDir is a command-line tool for archiving directories on macOS and Linux. It supports full and incremental backups, customizable exclusion patterns, Git-aware archive naming, and archive verification.

## Features
- Full and incremental directory archiving
- Exclusion patterns (like .gitignore)
- Git branch/hash in archive names
- Dry-run mode
- Archive listing
- Archive verification with checksums
- Automatic verification after creation

## Installation

### Using Go (Recommended)
```bash
go install github.com/fareedst/bkpdir@latest
```

### Manual Installation
1. Clone the repository:
```bash
git clone https://github.com/fareedst/bkpdir.git
cd bkpdir
```

2. Build the binary:
```bash
go build -o bkpdir
```

3. Move the binary to your PATH:
```bash
# For macOS/Linux
sudo mv bkpdir /usr/local/bin/
```

### Using Homebrew (macOS)
```bash
brew untap fareedst/bkpdir
brew tap fareedst/bkpdir
brew install bkpdir
```

## Usage

```
bkpdir full [--note NOTE] [--dry-run] [--verify]
bkpdir inc [--note NOTE] [--dry-run] [--verify]
bkpdir list
bkpdir verify ARCHIVE_NAME [--checksum]
```

## Configuration
Place a `.bkpdir.yml` file in the root of your directory. See the documentation for options.

### Verification Configuration
```yaml
verification:
  verify_on_create: false  # Automatically verify archives after creation
  checksum_algorithm: "sha256"  # Algorithm used for checksums
```

## Verification
BkpDir provides several ways to verify the integrity of your archives:

1. **Automatic Verification**: Enable `verify_on_create` in your configuration to automatically verify archives after creation.

2. **Manual Verification**: Use the `--verify` flag with `full` or `inc` commands to verify archives after creation.

3. **Standalone Verification**: Use the `verify` command to verify any archive:
   ```
   bkpdir verify archive-name.zip
   ```

4. **Checksum Verification**: Add the `--checksum` flag to verify file contents:
   ```
   bkpdir verify archive-name.zip --checksum
   ```

5. **Verification Status**: The `list` command shows verification status for each archive:
   ```
   archive-name.zip [VERIFIED]
   archive-name.zip [UNVERIFIED]
   archive-name.zip [FAILED]
   ``` 