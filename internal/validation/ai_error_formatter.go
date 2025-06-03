// 🔺 DOC-011: AI-optimized error processing - 📝 Error formatting for AI assistant consumption
package validation

import (
	"fmt"
)

// AIErrorFormatter formats validation errors for optimal AI assistant consumption
type AIErrorFormatter struct {
	remediationGuides map[string]*RemediationGuide
}

// NewAIErrorFormatter creates a new AI error formatter
func NewAIErrorFormatter() *AIErrorFormatter {
	formatter := &AIErrorFormatter{
		remediationGuides: make(map[string]*RemediationGuide),
	}

	// 🔺 DOC-011: Pre-defined remediation guides - 🔧 AI-friendly fix instructions
	formatter.initializeRemediationGuides()

	return formatter
}

// FormatErrors converts validation errors to AI-optimized format
func (f *AIErrorFormatter) FormatErrors(errors []ValidationError) []AIOptimizedError {
	var aiErrors []AIOptimizedError

	for _, err := range errors {
		// 🔺 DOC-011: Error transformation - 📝 Convert to AI-readable format
		aiError := AIOptimizedError{
			ErrorID:  err.ErrorID,
			Category: err.Category,
			Severity: err.Severity,
			Message:  f.enhanceErrorMessage(err.Message, err.Category),
			FileReference: &FileLocation{
				File:   err.File,
				Line:   err.Line,
				Column: 0, // Will be enhanced in future versions
			},
			Remediation: f.getRemediationGuide(err.Category),
			Context:     f.enhanceContext(err.Context, err.Category),
		}

		aiErrors = append(aiErrors, aiError)
	}

	return aiErrors
}

// FormatWarnings converts validation warnings to AI-optimized format
func (f *AIErrorFormatter) FormatWarnings(warnings []ValidationWarning) []AIOptimizedWarning {
	var aiWarnings []AIOptimizedWarning

	for _, warning := range warnings {
		// 🔺 DOC-011: Warning transformation - 📝 Convert to AI-readable format
		aiWarning := AIOptimizedWarning{
			WarningID: warning.WarningID,
			Category:  warning.Category,
			Message:   f.enhanceWarningMessage(warning.Message, warning.Category),
			FileReference: &FileLocation{
				File: warning.File,
				Line: warning.Line,
			},
			Suggestion: f.getSuggestionGuide(warning.Category),
		}

		aiWarnings = append(aiWarnings, aiWarning)
	}

	return aiWarnings
}

// enhanceErrorMessage improves error messages for AI assistant understanding
func (f *AIErrorFormatter) enhanceErrorMessage(message, category string) string {
	// 🔺 DOC-011: Message enhancement - 🔧 AI-friendly error descriptions
	switch category {
	case "token_format":
		return fmt.Sprintf("Implementation token format validation failed: %s. "+
			"Expected format: '// [PRIORITY_ICON] FEATURE-ID: Description - [ACTION_ICON] Context'", message)
	case "missing_icon":
		return fmt.Sprintf("Priority icon missing from implementation token: %s. "+
			"Required icons: ⭐ (CRITICAL), 🔺 (HIGH), 🔶 (MEDIUM), 🔻 (LOW)", message)
	case "documentation_sync":
		return fmt.Sprintf("Documentation cross-reference inconsistency detected: %s. "+
			"Check feature-tracking.md and related documentation files for consistency", message)
	case "priority_icon":
		return fmt.Sprintf("Priority icon validation failed: %s. "+
			"Ensure priority icon matches feature priority in feature-tracking.md", message)
	default:
		return fmt.Sprintf("Validation error: %s", message)
	}
}

// enhanceWarningMessage improves warning messages for AI assistant understanding
func (f *AIErrorFormatter) enhanceWarningMessage(message, category string) string {
	switch category {
	case "deprecated_usage":
		return fmt.Sprintf("Deprecated pattern detected: %s. "+
			"Consider updating to current standardized format", message)
	case "consistency_warning":
		return fmt.Sprintf("Consistency warning: %s. "+
			"Review related documentation for alignment", message)
	case "improvement_suggestion":
		return fmt.Sprintf("Improvement opportunity: %s. "+
			"Optional enhancement for better AI comprehension", message)
	default:
		return fmt.Sprintf("Validation warning: %s", message)
	}
}

// enhanceContext adds AI-relevant context to error information
func (f *AIErrorFormatter) enhanceContext(originalContext map[string]string, category string) map[string]string {
	context := make(map[string]string)

	// Copy original context
	for k, v := range originalContext {
		context[k] = v
	}

	// 🔺 DOC-011: Context enhancement - 🔍 Add AI-relevant information
	context["error_category"] = category
	context["ai_assistance_level"] = f.determineAssistanceLevel(category)
	context["automation_available"] = f.hasAutomation(category)

	switch category {
	case "token_format":
		context["doc_reference"] = "docs/context/source-code-icon-guidelines.md"
		context["validation_script"] = "scripts/validate-icon-enforcement.sh"
	case "missing_icon":
		context["icon_reference"] = "docs/context/README.md#priority-icons"
		context["auto_fix_script"] = "scripts/add-priority-icons.sh"
	case "documentation_sync":
		context["sync_target"] = "docs/context/feature-tracking.md"
		context["validation_command"] = "make validate-icon-enforcement"
	}

	return context
}

// initializeRemediationGuides sets up pre-defined remediation guides
func (f *AIErrorFormatter) initializeRemediationGuides() {
	// 🔺 DOC-011: Remediation guide initialization - 🛡️ AI assistant fix instructions
	f.remediationGuides["token_format"] = &RemediationGuide{
		Steps: []string{
			"1. Identify the implementation token in the source file",
			"2. Add priority icon based on feature priority (⭐🔺🔶🔻)",
			"3. Add action icon based on function behavior (🔍📝🔧🛡️)",
			"4. Ensure format: '// [PRIORITY] FEATURE-ID: Description - [ACTION] Context'",
			"5. Run validation: make validate-icon-enforcement",
		},
		References: []string{
			"docs/context/source-code-icon-guidelines.md",
			"docs/context/feature-tracking.md",
			"docs/context/ai-assistant-compliance.md",
		},
		Examples: []string{
			"// ⭐ ARCH-001: Archive naming convention - 🔧 Core functionality",
			"// 🔺 CFG-003: Template formatting logic - 📝 Configuration processing",
			"// 🔶 GIT-004: Git submodule support - 🔍 Discovery and validation",
		},
		Automation: "scripts/fix-token-format.sh",
	}

	f.remediationGuides["missing_icon"] = &RemediationGuide{
		Steps: []string{
			"1. Locate implementation token missing priority icon",
			"2. Check feature priority in docs/context/feature-tracking.md",
			"3. Add appropriate priority icon: ⭐(CRITICAL), 🔺(HIGH), 🔶(MEDIUM), 🔻(LOW)",
			"4. Verify icon matches feature priority level",
			"5. Run validation to confirm fix",
		},
		References: []string{
			"docs/context/README.md#priority-icons",
			"docs/context/feature-tracking.md",
		},
		Examples: []string{
			"Before: // ARCH-001: Archive naming",
			"After:  // ⭐ ARCH-001: Archive naming",
		},
		Automation: "scripts/add-priority-icons.sh",
	}

	f.remediationGuides["documentation_sync"] = &RemediationGuide{
		Steps: []string{
			"1. Identify the inconsistent cross-reference",
			"2. Check feature-tracking.md for correct feature information",
			"3. Update related documentation files for consistency",
			"4. Verify all cross-references are aligned",
			"5. Run comprehensive validation",
		},
		References: []string{
			"docs/context/feature-tracking.md",
			"docs/context/ai-assistant-protocol.md",
		},
		Examples: []string{
			"Update specification.md to match feature-tracking.md entries",
			"Align architecture.md with requirements.md cross-references",
		},
		Automation: "make validate-icon-enforcement",
	}

	f.remediationGuides["priority_icon"] = &RemediationGuide{
		Steps: []string{
			"1. Compare implementation token priority with feature-tracking.md",
			"2. Identify priority mismatch between code and documentation",
			"3. Update implementation token to match documented priority",
			"4. Ensure consistency across all related tokens",
			"5. Validate priority alignment",
		},
		References: []string{
			"docs/context/feature-tracking.md",
			"docs/context/source-code-icon-guidelines.md",
		},
		Automation: "scripts/priority-icon-inference.sh",
	}
}

// getRemediationGuide retrieves the appropriate remediation guide for an error category
func (f *AIErrorFormatter) getRemediationGuide(category string) *RemediationGuide {
	if guide, exists := f.remediationGuides[category]; exists {
		return guide
	}

	// 🔺 DOC-011: Fallback remediation - 🔧 Generic AI guidance
	return &RemediationGuide{
		Steps: []string{
			"1. Review the validation error message carefully",
			"2. Check relevant documentation for context",
			"3. Apply the suggested fix or correction",
			"4. Run validation to verify the fix",
		},
		References: []string{
			"docs/context/ai-assistant-compliance.md",
			"docs/context/feature-tracking.md",
		},
		Automation: "make validate-icon-enforcement",
	}
}

// getSuggestionGuide retrieves improvement suggestions for warnings
func (f *AIErrorFormatter) getSuggestionGuide(category string) *RemediationGuide {
	switch category {
	case "deprecated_usage":
		return &RemediationGuide{
			Steps: []string{
				"1. Identify the deprecated pattern or usage",
				"2. Check current standards documentation",
				"3. Update to use current standardized approach",
				"4. Verify compliance with current guidelines",
			},
			References: []string{
				"docs/context/source-code-icon-guidelines.md",
			},
		}
	case "consistency_warning":
		return &RemediationGuide{
			Steps: []string{
				"1. Review related documentation for alignment",
				"2. Identify inconsistencies in formatting or structure",
				"3. Apply consistent formatting across related content",
				"4. Validate consistency improvements",
			},
		}
	default:
		return &RemediationGuide{
			Steps: []string{
				"1. Review the warning context",
				"2. Consider the suggested improvement",
				"3. Apply changes if beneficial for AI comprehension",
			},
		}
	}
}

// determineAssistanceLevel determines the level of AI assistance needed
func (f *AIErrorFormatter) determineAssistanceLevel(category string) string {
	switch category {
	case "token_format", "missing_icon":
		return "automated" // Can be fixed automatically
	case "documentation_sync":
		return "guided" // Requires AI guidance but is systematic
	case "priority_icon":
		return "manual" // Requires understanding of feature context
	default:
		return "manual"
	}
}

// hasAutomation determines if automated fixing is available
func (f *AIErrorFormatter) hasAutomation(category string) string {
	automationMap := map[string]string{
		"token_format":       "scripts/fix-token-format.sh",
		"missing_icon":       "scripts/add-priority-icons.sh",
		"priority_icon":      "scripts/priority-icon-inference.sh",
		"documentation_sync": "make validate-icon-enforcement",
	}

	if script, exists := automationMap[category]; exists {
		return script
	}

	return "none"
}
