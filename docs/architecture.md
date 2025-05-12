## Data Objects

1. **Config**
   - `ArchiveDirPath`: string - Path where archives are stored
   - `UseCurrentDirName`: bool - Whether to use current directory name in archive path
   - `ExcludePatterns`: []string - Patterns for files to exclude

2. **Archive**
   - `Name`: string - Name of the archive file
   - `Path`: string - Full path to the archive file
   - `CreationTime`: time.Time - When the archive was created
   - `IsIncremental`: bool - Whether this is an incremental archive
   - `GitBranch`: string - Current git branch (if in git repo)
   - `GitHash`: string - Current git commit hash (if in git repo)
   - `Note`: string - Optional note for the archive
   - `BaseArchive`: string - Name of the base archive for incremental backups

## Core Functions

1. **Configuration Management**
   - `DefaultConfig() *Config`: Creates default configuration
   - `LoadConfig(root string) (*Config, error)`: Loads config from YAML or uses defaults

2. **Git Integration**
   - `IsGitRepository(dir string) bool`: Checks if directory is a git repo
   - `GetGitBranch(dir string) string`: Gets current branch name
   - `GetGitShortHash(dir string) string`: Gets current commit hash

3. **File System Operations**
   - `ShouldExcludeFile(filePath string, patterns []string) bool`: Checks if file matches exclude patterns
   - File listing is handled directly in `CreateFullArchive` and `CreateIncrementalArchive` using `filepath.Walk`

4. **Archive Management**
   - `GenerateArchiveName(prefix, timestamp, gitBranch, gitHash, note string, isGit, isIncremental bool, baseName string) string`: Generates archive filename
   - `ListArchives(archiveDir string) ([]Archive, error)`: Gets all archives in directory
   - `CreateFullArchive(cfg *Config, note string, dryRun bool) error`: Creates full archive
   - `CreateIncrementalArchive(cfg *Config, note string, dryRun bool) error`: Creates incremental archive

5. **Command Handlers**
   - `fullCmd() *cobra.Command`: Handles full archive creation
   - `incCmd() *cobra.Command`: Handles incremental archive creation
   - `listCmd() *cobra.Command`: Lists all archives

## Main Application Structure

1. **CLI Interface**
   - Uses `cobra` for command-line interface
   - Global flags:
     - `--dry-run`: Show what would be done without creating archives
   - Commands:
     - `full [NOTE]`: Create full archive with optional note
     - `inc [NOTE]`: Create incremental archive with optional note
     - `list`: List all archives

2. **Workflow Implementation**
   - For full archive:
     1. Load config
     2. Walk directory collecting files (excluding patterns)
     3. Get git info if applicable
     4. Generate archive name
     5. Create zip archive (or simulate in dry-run)

   - For incremental archive:
     1. Load config
     2. Find most recent full archive
     3. Get files changed since that archive
     4. Get git info if applicable
     5. Generate incremental archive name
     6. Create zip archive (or simulate in dry-run)

   - For listing archives:
     1. Load config
     2. Get archive directory path
     3. List all .zip files
     4. Display archive names

3. **Utility Functions**
   - File exclusion uses `doublestar` for glob pattern matching
   - Archive naming follows format:
     - Full: `[prefix-]timestamp[=branch=hash][=note].zip`
     - Incremental: `baseName_update=timestamp[=branch=hash][=note].zip`

The main differences from the original architecture are:
1. File listing is integrated into archive creation functions rather than being separate
2. Configuration is simpler with just three fields
3. Archive management functions are more focused on the specific needs of the application
4. The CLI interface is implemented using `cobra` with a simpler structure
5. File exclusion uses `doublestar` for more powerful glob pattern matching
