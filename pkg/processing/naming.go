// ⭐ EXTRACT-007: Naming conventions - Timestamp-based naming patterns extracted from archive.go and backup.go - 🔧
package processing

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// NamingProviderInterface defines the interface for timestamp-based naming
type NamingProviderInterface interface {
	GenerateName(template *NamingTemplate) (string, error)
	ParseName(name string, pattern string) (*NameComponents, error)
	ValidateName(name string, pattern string) error
	GetSupportedFormats() []string
}

// NamingTemplate represents a template for generating names with timestamps and metadata
type NamingTemplate struct {
	// Basic components
	Prefix    string    `json:"prefix"`
	Timestamp time.Time `json:"timestamp"`
	Note      string    `json:"note,omitempty"`

	// Git information (extracted from archive.go patterns)
	GitBranch          string `json:"git_branch,omitempty"`
	GitHash            string `json:"git_hash,omitempty"`
	GitIsClean         bool   `json:"git_is_clean"`
	ShowGitDirtyStatus bool   `json:"show_git_dirty_status"`

	// Additional metadata
	Metadata map[string]string `json:"metadata,omitempty"`

	// Format configuration
	TimestampFormat string `json:"timestamp_format"`
	IsIncremental   bool   `json:"is_incremental"`
	BaseName        string `json:"base_name,omitempty"`
}

// NameComponents represents parsed components from a filename
type NameComponents struct {
	Prefix    string            `json:"prefix"`
	Timestamp time.Time         `json:"timestamp"`
	Note      string            `json:"note,omitempty"`
	GitBranch string            `json:"git_branch,omitempty"`
	GitHash   string            `json:"git_hash,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Extension string            `json:"extension,omitempty"`
}

// NamingProvider implements timestamp-based naming conventions
type NamingProvider struct {
	// Timestamp formats (extracted from archive.go and backup.go)
	ArchiveTimestampFormat string // ISO 8601: YYYY-MM-DDTHHmmss
	BackupTimestampFormat  string // Backup format: YYYY-MM-DD-HH-MM

	// Validation patterns
	patterns map[string]*regexp.Regexp
}

// NewNamingProvider creates a new naming provider with default formats
func NewNamingProvider() *NamingProvider {
	provider := &NamingProvider{
		ArchiveTimestampFormat: "2006-01-02T150405", // ISO 8601 format from archive.go
		BackupTimestampFormat:  "2006-01-02-15-04",  // Backup format from backup.go
		patterns:               make(map[string]*regexp.Regexp),
	}

	// Initialize validation patterns
	provider.initializePatterns()
	return provider
}

// initializePatterns sets up regex patterns for name validation and parsing
func (np *NamingProvider) initializePatterns() {
	// Archive pattern: {prefix}-{timestamp}-{git_branch}-{git_hash}-{note}.zip
	// Example: test-2024-01-01T120000-main-abc123-backup.zip
	np.patterns["archive"] = regexp.MustCompile(
		`^(?P<prefix>[^-]+)-(?P<timestamp>\d{4}-\d{2}-\d{2}T\d{6})-(?P<git_branch>[^-]+)-(?P<git_hash>[a-f0-9]+)-(?P<note>[^-.]+)\.(?P<extension>zip)$`,
	)

	// Backup pattern: {filename}-{timestamp}[={note}]
	// Example: document.txt-2024-03-20-14-30=before-changes
	np.patterns["backup"] = regexp.MustCompile(
		`^(?P<filename>.+)-(?P<timestamp>\d{4}-\d{2}-\d{2}-\d{2}-\d{2})(?:=(?P<note>.+))?$`,
	)

	// Incremental archive pattern: {prefix}-{timestamp}-{git_info}-inc-{base}-{note}.zip
	np.patterns["incremental"] = regexp.MustCompile(
		`^(?P<prefix>[^-]+)-(?P<timestamp>\d{4}-\d{2}-\d{2}T\d{6})(?:-(?P<git_branch>[^-]+)-(?P<git_hash>[a-f0-9]+)(?P<git_dirty>-dirty)?)?-inc-(?P<base>[^-]+)(?:-(?P<note>[^-.]+))?\\.(?P<extension>zip)$`,
	)
}

// GenerateName creates a name using the provided template
func (np *NamingProvider) GenerateName(template *NamingTemplate) (string, error) {
	if template == nil {
		return "", NewProcessingError("INVALID_TEMPLATE", "GenerateName", "template cannot be nil")
	}

	// Use template timestamp format or default
	timestampFormat := template.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = np.ArchiveTimestampFormat
	}

	// Format timestamp
	timestamp := template.Timestamp.Format(timestampFormat)

	// Build name components
	var parts []string

	// Add prefix
	if template.Prefix != "" {
		parts = append(parts, template.Prefix)
	}

	// Add timestamp
	parts = append(parts, timestamp)

	// Add Git information if available
	if template.GitBranch != "" && template.GitHash != "" {
		gitInfo := template.GitBranch + "-" + template.GitHash
		if template.ShowGitDirtyStatus && !template.GitIsClean {
			gitInfo += "-dirty"
		}
		parts = append(parts, gitInfo)
	}

	// Add incremental marker and base name
	if template.IsIncremental && template.BaseName != "" {
		parts = append(parts, "inc", template.BaseName)
	}

	// Add note if provided
	if template.Note != "" {
		parts = append(parts, template.Note)
	}

	// Join parts with separator
	name := strings.Join(parts, "-")

	return name, nil
}

// GenerateArchiveName creates an archive name using archive.go patterns
func (np *NamingProvider) GenerateArchiveName(prefix, timestamp, gitBranch, gitHash, note string, gitIsClean, showGitDirtyStatus, isIncremental bool, baseName string) string {
	template := &NamingTemplate{
		Prefix:             prefix,
		Timestamp:          parseTimestamp(timestamp, np.ArchiveTimestampFormat),
		GitBranch:          gitBranch,
		GitHash:            gitHash,
		GitIsClean:         gitIsClean,
		ShowGitDirtyStatus: showGitDirtyStatus,
		Note:               note,
		IsIncremental:      isIncremental,
		BaseName:           baseName,
		TimestampFormat:    np.ArchiveTimestampFormat,
	}

	name, _ := np.GenerateName(template)
	return name + ".zip"
}

// GenerateBackupName creates a backup name using backup.go patterns
func (np *NamingProvider) GenerateBackupName(sourcePath, timestamp, note string) string {
	template := &NamingTemplate{
		Prefix:          sourcePath,
		Timestamp:       parseTimestamp(timestamp, np.BackupTimestampFormat),
		Note:            note,
		TimestampFormat: np.BackupTimestampFormat,
	}

	name, _ := np.GenerateName(template)

	// Add note with equals sign for backup format
	if note != "" {
		name += "=" + note
	}

	return name
}

// ParseName extracts components from a filename using the specified pattern
func (np *NamingProvider) ParseName(name string, pattern string) (*NameComponents, error) {
	regex, exists := np.patterns[pattern]
	if !exists {
		return nil, NewProcessingError("INVALID_PATTERN", "ParseName", fmt.Sprintf("unsupported pattern: %s", pattern))
	}

	matches := regex.FindStringSubmatch(name)
	if matches == nil {
		return nil, NewProcessingError("PARSE_FAILED", "ParseName", fmt.Sprintf("name does not match pattern %s: %s", pattern, name))
	}

	// Extract named groups
	result := &NameComponents{
		Metadata: make(map[string]string),
	}

	for i, groupName := range regex.SubexpNames() {
		if i == 0 || groupName == "" || i >= len(matches) {
			continue
		}

		value := matches[i]
		switch groupName {
		case "prefix", "filename":
			result.Prefix = value
		case "timestamp":
			timestamp, err := time.Parse(getTimestampFormat(pattern), value)
			if err != nil {
				return nil, NewProcessingError("TIMESTAMP_PARSE", "ParseName", fmt.Sprintf("failed to parse timestamp %s: %v", value, err))
			}
			result.Timestamp = timestamp
		case "note":
			result.Note = value
		case "note_no_git":
			result.Note = value
		case "git_branch":
			result.GitBranch = value
		case "git_hash":
			result.GitHash = value
		case "extension":
			result.Extension = value
		default:
			result.Metadata[groupName] = value
		}
	}

	return result, nil
}

// ValidateName checks if a name matches the specified pattern
func (np *NamingProvider) ValidateName(name string, pattern string) error {
	_, err := np.ParseName(name, pattern)
	return err
}

// GetSupportedFormats returns the list of supported naming formats
func (np *NamingProvider) GetSupportedFormats() []string {
	formats := make([]string, 0, len(np.patterns))
	for format := range np.patterns {
		formats = append(formats, format)
	}
	return formats
}

// Helper functions

// parseTimestamp parses a timestamp string using the provided format
func parseTimestamp(timestamp, format string) time.Time {
	if timestamp == "" {
		return time.Now()
	}

	parsed, err := time.Parse(format, timestamp)
	if err != nil {
		return time.Now()
	}

	return parsed
}

// getTimestampFormat returns the appropriate timestamp format for a pattern
func getTimestampFormat(pattern string) string {
	switch pattern {
	case "archive", "incremental":
		return "2006-01-02T150405"
	case "backup":
		return "2006-01-02-15-04"
	default:
		return "2006-01-02T150405"
	}
}

// FormatTimestamp formats a time using the specified format type
func (np *NamingProvider) FormatTimestamp(t time.Time, formatType string) string {
	switch formatType {
	case "archive", "incremental":
		return t.Format(np.ArchiveTimestampFormat)
	case "backup":
		return t.Format(np.BackupTimestampFormat)
	default:
		return t.Format(np.ArchiveTimestampFormat)
	}
}

// CreateTimestamp creates a current timestamp in the specified format
func (np *NamingProvider) CreateTimestamp(formatType string) string {
	return np.FormatTimestamp(time.Now(), formatType)
}
