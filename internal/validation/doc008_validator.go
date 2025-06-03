// 🔺 DOC-011: DOC-008 validation framework integration - 🔧 Bridge to existing validation system
package validation

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// DOC008Validator provides integration with the existing DOC-008 validation framework
type DOC008Validator struct {
	scriptPath string
}

// ValidationResult represents the result of a validation operation
type ValidationResult struct {
	Errors      []ValidationError   `json:"errors"`
	Warnings    []ValidationWarning `json:"warnings"`
	TotalTokens int                 `json:"total_tokens"`
	Status      string              `json:"status"`
}

// ValidationError represents a validation error
type ValidationError struct {
	ErrorID  string            `json:"error_id"`
	Category string            `json:"category"`
	Severity string            `json:"severity"`
	Message  string            `json:"message"`
	File     string            `json:"file"`
	Line     int               `json:"line"`
	Context  map[string]string `json:"context"`
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	WarningID string            `json:"warning_id"`
	Category  string            `json:"category"`
	Message   string            `json:"message"`
	File      string            `json:"file"`
	Line      int               `json:"line"`
	Context   map[string]string `json:"context"`
}

// NewDOC008Validator creates a new DOC-008 validator
func NewDOC008Validator() *DOC008Validator {
	return &DOC008Validator{
		scriptPath: "scripts/validate-icon-enforcement.sh",
	}
}

// ValidateFiles validates the specified files using DOC-008 validation framework
func (v *DOC008Validator) ValidateFiles(ctx context.Context, files []string, mode string) (*ValidationResult, error) {
	// 🔺 DOC-011: DOC-008 framework bridge - 🔍 Execute existing validation
	cmd := v.buildValidationCommand(mode)

	// Execute validation with context support
	output, err := v.executeValidationCommand(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("validation command failed: %w", err)
	}

	// 🔺 DOC-011: Validation result parsing - 📝 Parse DOC-008 output for AI consumption
	result, err := v.parseValidationOutput(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse validation output: %w", err)
	}

	return result, nil
}

// buildValidationCommand constructs the appropriate validation command
func (v *DOC008Validator) buildValidationCommand(mode string) *exec.Cmd {
	switch mode {
	case "strict":
		return exec.Command("make", "validate-icons-strict")
	case "legacy":
		return exec.Command("make", "validate-icons")
	default: // "standard"
		return exec.Command("make", "validate-icon-enforcement")
	}
}

// executeValidationCommand executes the validation command with context support
func (v *DOC008Validator) executeValidationCommand(ctx context.Context, cmd *exec.Cmd) ([]byte, error) {
	// 🔺 DOC-011: Context-aware validation execution - 🛡️ Cancellation support
	cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)

	output, err := cmd.CombinedOutput()
	return output, err
}

// parseValidationOutput parses the validation output into structured format
func (v *DOC008Validator) parseValidationOutput(output []byte) (*ValidationResult, error) {
	result := &ValidationResult{
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
		Status:   "pass",
	}

	lines := strings.Split(string(output), "\n")

	// 🔺 DOC-011: Validation output parsing - 📝 Convert DOC-008 output to AI format
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse different types of validation output
		if strings.Contains(line, "ERROR:") {
			error := v.parseErrorLine(line, i+1)
			if error != nil {
				result.Errors = append(result.Errors, *error)
			}
		} else if strings.Contains(line, "WARNING:") {
			warning := v.parseWarningLine(line, i+1)
			if warning != nil {
				result.Warnings = append(result.Warnings, *warning)
			}
		} else if strings.Contains(line, "Total tokens:") {
			if tokens, err := v.parseTokenCount(line); err == nil {
				result.TotalTokens = tokens
			}
		}
	}

	// Determine overall status
	if len(result.Errors) > 0 {
		result.Status = "fail"
	} else if len(result.Warnings) > 5 {
		result.Status = "warning"
	}

	return result, nil
}

// parseErrorLine parses an error line from validation output
func (v *DOC008Validator) parseErrorLine(line string, lineNum int) *ValidationError {
	// 🔺 DOC-011: Error parsing - 🔍 Extract structured error information
	parts := strings.Split(line, ":")
	if len(parts) < 3 {
		return &ValidationError{
			ErrorID:  fmt.Sprintf("DOC011-PARSE-%d", lineNum),
			Category: "parse_error",
			Severity: "medium",
			Message:  line,
			Context:  map[string]string{"raw_line": line},
		}
	}

	return &ValidationError{
		ErrorID:  fmt.Sprintf("DOC008-ERR-%d", lineNum),
		Category: v.determineErrorCategory(line),
		Severity: "high",
		Message:  strings.Join(parts[1:], ":"),
		File:     v.extractFileName(line),
		Line:     v.extractLineNumber(line),
		Context:  map[string]string{"validation_line": strconv.Itoa(lineNum)},
	}
}

// parseWarningLine parses a warning line from validation output
func (v *DOC008Validator) parseWarningLine(line string, lineNum int) *ValidationWarning {
	parts := strings.Split(line, ":")
	if len(parts) < 3 {
		return &ValidationWarning{
			WarningID: fmt.Sprintf("DOC011-PARSE-WARN-%d", lineNum),
			Category:  "parse_warning",
			Message:   line,
			Context:   map[string]string{"raw_line": line},
		}
	}

	return &ValidationWarning{
		WarningID: fmt.Sprintf("DOC008-WARN-%d", lineNum),
		Category:  v.determineWarningCategory(line),
		Message:   strings.Join(parts[1:], ":"),
		File:      v.extractFileName(line),
		Line:      v.extractLineNumber(line),
		Context:   map[string]string{"validation_line": strconv.Itoa(lineNum)},
	}
}

// determineErrorCategory categorizes errors for AI assistant understanding
func (v *DOC008Validator) determineErrorCategory(line string) string {
	line = strings.ToLower(line)

	if strings.Contains(line, "token") && strings.Contains(line, "format") {
		return "token_format"
	}
	if strings.Contains(line, "icon") && strings.Contains(line, "missing") {
		return "missing_icon"
	}
	if strings.Contains(line, "cross-reference") || strings.Contains(line, "documentation") {
		return "documentation_sync"
	}
	if strings.Contains(line, "priority") {
		return "priority_icon"
	}

	return "general_validation"
}

// determineWarningCategory categorizes warnings for AI assistant understanding
func (v *DOC008Validator) determineWarningCategory(line string) string {
	line = strings.ToLower(line)

	if strings.Contains(line, "deprecated") {
		return "deprecated_usage"
	}
	if strings.Contains(line, "consistency") {
		return "consistency_warning"
	}
	if strings.Contains(line, "suggestion") {
		return "improvement_suggestion"
	}

	return "general_warning"
}

// extractFileName extracts filename from validation output
func (v *DOC008Validator) extractFileName(line string) string {
	// Look for file patterns in validation output
	parts := strings.Fields(line)
	for _, part := range parts {
		if strings.Contains(part, ".go") || strings.Contains(part, ".md") {
			return part
		}
	}
	return ""
}

// extractLineNumber extracts line number from validation output
func (v *DOC008Validator) extractLineNumber(line string) int {
	// Look for line number patterns like ":123:" or "(line 123)"
	parts := strings.Split(line, ":")
	for _, part := range parts {
		if num, err := strconv.Atoi(strings.TrimSpace(part)); err == nil && num > 0 {
			return num
		}
	}
	return 0
}

// parseTokenCount extracts token count from validation output
func (v *DOC008Validator) parseTokenCount(line string) (int, error) {
	parts := strings.Fields(line)
	for i, part := range parts {
		if part == "tokens:" && i+1 < len(parts) {
			return strconv.Atoi(parts[i+1])
		}
	}
	return 0, fmt.Errorf("token count not found in line: %s", line)
}
