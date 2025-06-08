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

### Using Pre-compiled Binaries (Recommended)

BkpDir provides pre-compiled binaries for Ubuntu and macOS. This is the fastest and easiest installation method.

#### Ubuntu (20.04, 22.04, 24.04)

**Quick install for Ubuntu 24.04:**
```bash
# Download and install in one command
curl -L https://github.com/fareedst/bkpdir/raw/main/bin/bkpdir-ubuntu24.04 -o bkpdir && \
chmod +x bkpdir && \
sudo mv bkpdir /usr/local/bin/
```

**For other Ubuntu versions:**
```bash
# Ubuntu 22.04
curl -L https://github.com/fareedst/bkpdir/raw/main/bin/bkpdir-ubuntu22.04 -o bkpdir && \
chmod +x bkpdir && sudo mv bkpdir /usr/local/bin/

# Ubuntu 20.04  
curl -L https://github.com/fareedst/bkpdir/raw/main/bin/bkpdir-ubuntu20.04 -o bkpdir && \
chmod +x bkpdir && sudo mv bkpdir /usr/local/bin/
```

**Alternative with wget:**
```bash
# Replace 'ubuntu24.04' with your Ubuntu version (ubuntu20.04, ubuntu22.04, ubuntu24.04)
wget https://github.com/fareedst/bkpdir/raw/main/bin/bkpdir-ubuntu24.04 -O bkpdir && \
chmod +x bkpdir && \
sudo mv bkpdir /usr/local/bin/
```

#### macOS (Intel and Apple Silicon)

**For Apple Silicon (M1/M2) Macs:**
```bash
curl -L https://github.com/fareedst/bkpdir/raw/main/bin/bkpdir-macos-arm64 -o bkpdir && \
chmod +x bkpdir && \
sudo mv bkpdir /usr/local/bin/
```

**For Intel Macs:**
```bash
curl -L https://github.com/fareedst/bkpdir/raw/main/bin/bkpdir-macos-amd64 -o bkpdir && \
chmod +x bkpdir && \
sudo mv bkpdir /usr/local/bin/
```

**Auto-detect macOS architecture:**
```bash
# This script automatically detects your Mac architecture
arch=$(uname -m)
if [[ "$arch" == "arm64" ]]; then
    binary="bkpdir-macos-arm64"
else
    binary="bkpdir-macos-amd64"
fi
curl -L "https://github.com/fareedst/bkpdir/raw/main/bin/$binary" -o bkpdir && \
chmod +x bkpdir && \
sudo mv bkpdir /usr/local/bin/
```

#### User-specific Installation (No sudo required)

If you prefer to install without sudo access, you can install to your personal bin directory:

```bash
# Create personal bin directory if it doesn't exist
mkdir -p ~/bin

# For Ubuntu (replace with your version)
curl -L https://github.com/fareedst/bkpdir/raw/main/bin/bkpdir-ubuntu24.04 -o ~/bin/bkpdir && chmod +x ~/bin/bkpdir

# For macOS (replace with your architecture)  
curl -L https://github.com/fareedst/bkpdir/raw/main/bin/bkpdir-macos-arm64 -o ~/bin/bkpdir && chmod +x ~/bin/bkpdir

# Add ~/bin to your PATH (add this to your ~/.bashrc or ~/.zshrc)
export PATH="$HOME/bin:$PATH"
```

#### Verify Installation

After installation, verify bkpdir is working:
```bash
bkpdir --help
bkpdir --version  # If version flag is supported
```

### Alternative Installation Methods

#### Using Go
```bash
go install github.com/fareedst/bkpdir@latest
```

#### Manual Build from Source
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

#### Using Homebrew (macOS)
```bash
brew untap fareedst/bkpdir
brew tap fareedst/bkpdir  
brew install bkpdir
```

### Troubleshooting

**Permission denied when running bkpdir:**
```bash
chmod +x /usr/local/bin/bkpdir
# or if installed in ~/bin
chmod +x ~/bin/bkpdir
```

**Command not found:**
- For system-wide installation: Ensure `/usr/local/bin` is in your PATH
- For user installation: Ensure `~/bin` is in your PATH with `export PATH="$HOME/bin:$PATH"`

**Download fails:**
- Check your internet connection
- Try using `wget` instead of `curl`
- Manually download from: https://github.com/fareedst/bkpdir/tree/main/bin

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