# BkpDir

BkpDir is a command-line tool for archiving directories on macOS and Linux. It supports full and incremental backups, customizable exclusion patterns, and Git-aware archive naming.

## Features
- Full and incremental directory archiving
- Exclusion patterns (like .gitignore)
- Git branch/hash in archive names
- Dry-run mode
- Archive listing

## Usage

```
bkpdir full [--note NOTE] [--dry-run]
bkpdir inc [--note NOTE] [--dry-run]
bkpdir list
```

## Configuration
Place a `.bkpdir.yaml` file in the root of your directory. See the documentation for options. 