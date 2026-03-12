// [REQ-DOC_015] Semantic Token Validator
// Validates semantic token usage and consistency across the codebase
package validation

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// SemanticTokenValidator validates semantic token usage and consistency
type SemanticTokenValidator struct {
	projectRoot string
}

// SemanticTokenValidationResult represents the result of semantic token validation
type SemanticTokenValidationResult struct {
	TotalFilesScanned   int
	FilesWithREQTokens  int
	FilesWithARCHTokens int
	FilesWithIMPLTokens int
	TotalREQTokens      int
	TotalARCHTokens     int
	TotalIMPLTokens     int
	MissingTokens       []TokenError
	MalformedTokens     []TokenError
	InconsistentTokens  []TokenError
	Errors              []ValidationError
	Warnings            []ValidationWarning
	Status              string
}

// TokenError represents a token validation error
type TokenError struct {
	File     string
	Line     int
	Token    string
	Message  string
	Category string
}

// NewSemanticTokenValidator creates a new semantic token validator
func NewSemanticTokenValidator(projectRoot string) *SemanticTokenValidator {
	return &SemanticTokenValidator{
		projectRoot: projectRoot,
	}
}

// ValidateSemanticTokens validates semantic token usage across the codebase
// [REQ-DOC_015] Validates REQ/ARCH/IMPL token usage and consistency
func (v *SemanticTokenValidator) ValidateSemanticTokens(ctx context.Context) (*SemanticTokenValidationResult, error) {
	result := &SemanticTokenValidationResult{
		MissingTokens:      []TokenError{},
		MalformedTokens:    []TokenError{},
		InconsistentTokens: []TokenError{},
		Errors:             []ValidationError{},
		Warnings:           []ValidationWarning{},
		Status:             "pass",
	}

	// [REQ-DOC_015] Semantic token patterns (supports both colon and hyphen delimiters)
	reqTokenPattern := regexp.MustCompile(`\[REQ[:\-]([A-Z0-9_]+)\]`)
	archTokenPattern := regexp.MustCompile(`\[ARCH[:\-]([A-Z0-9_]+)\]`)
	implTokenPattern := regexp.MustCompile(`\[IMPL[:\-]([A-Z0-9_]+)\]`)

	// Walk through all relevant files
	err := filepath.Walk(v.projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Only process Go and Markdown files
		if info.IsDir() {
			// Skip certain directories
			if strings.Contains(path, ".git") || strings.Contains(path, "node_modules") ||
				strings.Contains(path, "vendor") || strings.Contains(path, "bin") {
				return filepath.SkipDir
			}
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".go" && ext != ".md" {
			return nil
		}

		result.TotalFilesScanned++

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		lines := strings.Split(string(content), "\n")

		// Validate tokens in file
		for lineNum, line := range lines {
			// Check for REQ tokens
			reqMatches := reqTokenPattern.FindAllStringSubmatch(line, -1)
			for _, match := range reqMatches {
				if len(match) > 1 {
					result.TotalREQTokens++
					if !v.isValidTokenIdentifier(match[1]) {
						result.MalformedTokens = append(result.MalformedTokens, TokenError{
							File:     path,
							Line:     lineNum + 1,
							Token:    match[0],
							Message:  fmt.Sprintf("Invalid REQ token identifier: %s", match[1]),
							Category: "malformed_token",
						})
					}
				}
			}

			// Check for ARCH tokens
			archMatches := archTokenPattern.FindAllStringSubmatch(line, -1)
			for _, match := range archMatches {
				if len(match) > 1 {
					result.TotalARCHTokens++
					if !v.isValidTokenIdentifier(match[1]) {
						result.MalformedTokens = append(result.MalformedTokens, TokenError{
							File:     path,
							Line:     lineNum + 1,
							Token:    match[0],
							Message:  fmt.Sprintf("Invalid ARCH token identifier: %s", match[1]),
							Category: "malformed_token",
						})
					}
				}
			}

			// Check for IMPL tokens
			implMatches := implTokenPattern.FindAllStringSubmatch(line, -1)
			for _, match := range implMatches {
				if len(match) > 1 {
					result.TotalIMPLTokens++
					if !v.isValidTokenIdentifier(match[1]) {
						result.MalformedTokens = append(result.MalformedTokens, TokenError{
							File:     path,
							Line:     lineNum + 1,
							Token:    match[0],
							Message:  fmt.Sprintf("Invalid IMPL token identifier: %s", match[1]),
							Category: "malformed_token",
						})
					}
				}
			}
		}

		// Count files with tokens (supports both colon and hyphen delimiters)
		if strings.Contains(string(content), "[REQ-") || strings.Contains(string(content), "[REQ:") {
			result.FilesWithREQTokens++
		}
		if strings.Contains(string(content), "[ARCH-") || strings.Contains(string(content), "[ARCH:") {
			result.FilesWithARCHTokens++
		}
		if strings.Contains(string(content), "[IMPL-") || strings.Contains(string(content), "[IMPL:") {
			result.FilesWithIMPLTokens++
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk project directory: %w", err)
	}

	// Convert token errors to validation errors
	for _, tokenErr := range result.MalformedTokens {
		result.Errors = append(result.Errors, ValidationError{
			ErrorID:  fmt.Sprintf("SEMTOKEN-ERR-%d", len(result.Errors)+1),
			Category: tokenErr.Category,
			Severity: "high",
			Message:  tokenErr.Message,
			File:     tokenErr.File,
			Line:     tokenErr.Line,
			Context:  map[string]string{"token": tokenErr.Token},
		})
	}

	// Determine overall status
	if len(result.Errors) > 0 {
		result.Status = "fail"
	} else if len(result.Warnings) > 5 {
		result.Status = "warning"
	}

	return result, nil
}

// isValidTokenIdentifier validates that a token identifier follows the naming convention
// [REQ-DOC_015] Token identifiers must use UPPER_SNAKE_CASE
func (v *SemanticTokenValidator) isValidTokenIdentifier(identifier string) bool {
	// Must be non-empty
	if len(identifier) == 0 {
		return false
	}

	// Must match UPPER_SNAKE_CASE pattern: [A-Z][A-Z0-9_]*[A-Z0-9]
	pattern := regexp.MustCompile(`^[A-Z][A-Z0-9_]*[A-Z0-9]$`)
	return pattern.MatchString(identifier)
}

// ValidateTokenConsistency validates cross-layer token consistency
// [REQ-DOC_015] Ensures tokens are used consistently across requirements, architecture, and implementation
func (v *SemanticTokenValidator) ValidateTokenConsistency(ctx context.Context) ([]TokenError, error) {
	// This would check that:
	// - REQ tokens have corresponding ARCH tokens
	// - ARCH tokens have corresponding IMPL tokens
	// - IMPL tokens reference valid ARCH and REQ tokens
	// Implementation would require reading semantic-tokens.md registry
	// For now, return empty slice as placeholder
	return []TokenError{}, nil
}
