// [REQ-CUSTOMIZABLE_FORMAT_STRINGS] Simplified list formatting implementation
// This file contains a simplified, testable implementation for list command formatting
// CRITICAL: Only supports #{name} style placeholders, NOT printf-style (%s, %d)
package main

import (
	"bkpdir/pkg/formatter"
	"fmt"
	"os"
	"strings"
)

// formatListArchiveSimple is a simplified implementation that:
// 1. Always gathers file statistics (needed for template placeholders)
// 2. Uses FormatListArchive if it contains template placeholders, otherwise TemplateListArchive
// 3. Uses the simple ReplacePlaceholders function from pkg/formatter
// 4. Returns formatted string ready for output
func formatListArchiveSimple(cfg *Config, formatterInstance *FormatterAdapter, archivePath, creationTime string) string {
	// DEBUG: Log entry point
	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: formatListArchiveSimple called with path=%q, creationTime=%q\n", archivePath, creationTime)
	}

	// Step 1: Always gather file statistics (needed for template placeholders)
	statInfo, err := formatter.GatherFileStatInfo(archivePath)

	// Step 2: Build data map with all available information
	data := make(map[string]string)
	data["path"] = archivePath
	data["creation_time"] = creationTime

	// Extract filename from path
	parts := strings.Split(archivePath, "/")
	if len(parts) > 0 {
		data["name"] = parts[len(parts)-1]
	} else {
		data["name"] = archivePath
	}

	// Add file statistics if available
	if err == nil && statInfo != nil {
		data["size"] = fmt.Sprintf("%d", statInfo.Size)
		data["size_human"] = statInfo.SizeHuman
		data["mtime"] = statInfo.MTime.Format("2006-01-02 15:04:05")
		data["mtime_unix"] = fmt.Sprintf("%d", statInfo.MTimeUnix)
		data["mode"] = statInfo.Mode.String()
		data["type"] = statInfo.Type
		// Override name with stat name if available
		if statInfo.Name != "" {
			data["name"] = statInfo.Name
		}
	} else {
		// Provide defaults for missing stats
		data["size"] = "0"
		data["size_human"] = "unknown"
		data["mtime"] = creationTime
		data["mtime_unix"] = "0"
		data["mode"] = "unknown"
		data["type"] = "unknown"
	}

	// Step 3: Select format string
	// Priority: FormatListArchive (if contains #{) > TemplateListArchive > default
	var formatStr string

	if cfg.FormatListArchive != "" && strings.Contains(cfg.FormatListArchive, "#{") {
		// FormatListArchive contains template placeholders, use it
		formatStr = cfg.FormatListArchive
		if debug {
			fmt.Fprintf(os.Stderr, "DEBUG: formatListArchiveSimple using FormatListArchive=%q\n", formatStr)
		}
	} else if cfg.TemplateListArchive != "" {
		// Use TemplateListArchive if FormatListArchive is empty or doesn't have template placeholders
		formatStr = cfg.TemplateListArchive
		if debug {
			fmt.Fprintf(os.Stderr, "DEBUG: formatListArchiveSimple using TemplateListArchive=%q\n", formatStr)
		}
	} else {
		// Default template format when both are empty
		formatStr = "#{path} (size: #{size_human})\n"
		if debug {
			fmt.Fprintf(os.Stderr, "DEBUG: formatListArchiveSimple using default format=%q\n", formatStr)
		}
	}

	// Step 4: Use the simple ReplacePlaceholders function
	// This function has been thoroughly tested and is guaranteed to work
	result := formatter.ReplacePlaceholders(formatStr, data)
	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: formatListArchiveSimple after ReplacePlaceholders: result=%q\n", result)
		fmt.Fprintf(os.Stderr, "DEBUG: formatListArchiveSimple data map: %+v\n", data)
	}

	// CRITICAL: Final safety check - ALWAYS replace any remaining known placeholders
	// This prevents fmt.Sprintf errors when the result is printed
	// Extract filename for #{name}
	nameParts := strings.Split(archivePath, "/")
	nameValue := archivePath
	if len(nameParts) > 0 {
		nameValue = nameParts[len(nameParts)-1]
	}

	// ALWAYS perform final replacements, even if ReplacePlaceholders worked
	// This ensures no #{...} patterns remain that could cause fmt errors
	finalReplacements := map[string]string{
		"#{path}":          archivePath,
		"#{name}":          nameValue,
		"#{creation_time}": creationTime,
		"#{size_human}":    data["size_human"], // Use actual value from data map
		"#{size}":          data["size"],       // Use actual value from data map
		"#{mtime}":         data["mtime"],      // Use actual value from data map
		"#{mtime_unix}":    data["mtime_unix"], // Use actual value from data map
		"#{mode}":          data["mode"],       // Use actual value from data map
		"#{type}":          data["type"],       // Use actual value from data map
	}
	for placeholder, value := range finalReplacements {
		if strings.Contains(result, placeholder) {
			result = strings.ReplaceAll(result, placeholder, value)
		}
	}

	// CRITICAL: Final verification - if any #{...} patterns remain, replace with safe fallbacks
	if strings.Contains(result, "#{") {
		if debug {
			fmt.Fprintf(os.Stderr, "DEBUG: formatListArchiveSimple WARNING: result still contains #{} patterns before emergency replacement: %q\n", result)
		}
		// Emergency fallback: replace any remaining #{...} with safe values
		emergencyReplacements := map[string]string{
			"#{path}":          archivePath,
			"#{name}":          nameValue,
			"#{creation_time}": creationTime,
			"#{size_human}":    "unknown",
			"#{size}":          "0",
			"#{mtime}":         creationTime,
			"#{mtime_unix}":    "0",
			"#{mode}":          "unknown",
			"#{type}":          "unknown",
		}
		for placeholder, value := range emergencyReplacements {
			result = strings.ReplaceAll(result, placeholder, value)
		}
		if debug {
			fmt.Fprintf(os.Stderr, "DEBUG: formatListArchiveSimple after emergency replacement: result=%q\n", result)
		}
	}

	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: formatListArchiveSimple final result: %q\n", result)
		fmt.Fprintf(os.Stderr, "DEBUG: formatListArchiveSimple contains #{: %v\n", strings.Contains(result, "#{"))
	}

	return result
}
