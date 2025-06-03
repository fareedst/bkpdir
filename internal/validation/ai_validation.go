// 🔺 DOC-011: AI validation integration - Core validation framework for AI assistants
package validation

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// AIValidationGateway provides the main entry point for AI assistant validation
type AIValidationGateway struct {
	validator         *DOC008Validator
	errorFormatter    *AIErrorFormatter
	complianceTracker *ComplianceMonitor
	bypassManager     *BypassController
}

// ValidationRequest represents a validation request from an AI assistant
type ValidationRequest struct {
	SourceFiles    []string          `json:"source_files"`
	ChangedTokens  []string          `json:"changed_tokens"`
	ValidationMode string            `json:"validation_mode"` // "standard", "strict", "legacy"
	RequestContext *AIRequestContext `json:"request_context"`
}

// ValidationResponse represents the response to an AI assistant validation request
type ValidationResponse struct {
	Status           string               `json:"status"` // "pass", "fail", "warning"
	Errors           []AIOptimizedError   `json:"errors"`
	Warnings         []AIOptimizedWarning `json:"warnings"`
	RemediationSteps []RemediationAction  `json:"remediation_steps"`
	ComplianceScore  float64              `json:"compliance_score"`
	ProcessingTime   time.Duration        `json:"processing_time"`
}

// AIRequestContext provides context about the AI assistant making the request
type AIRequestContext struct {
	AssistantID string    `json:"assistant_id"`
	SessionID   string    `json:"session_id"`
	Timestamp   time.Time `json:"timestamp"`
}

// AIOptimizedError represents an error formatted for AI assistant consumption
type AIOptimizedError struct {
	ErrorID       string            `json:"error_id"`
	Category      string            `json:"category"`
	Severity      string            `json:"severity"`
	Message       string            `json:"message"`
	FileReference *FileLocation     `json:"file_reference"`
	Remediation   *RemediationGuide `json:"remediation"`
	Context       map[string]string `json:"context"`
}

// AIOptimizedWarning represents a warning formatted for AI assistant consumption
type AIOptimizedWarning struct {
	WarningID     string            `json:"warning_id"`
	Category      string            `json:"category"`
	Message       string            `json:"message"`
	FileReference *FileLocation     `json:"file_reference"`
	Suggestion    *RemediationGuide `json:"suggestion"`
}

// FileLocation represents a specific location in a source file
type FileLocation struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}

// RemediationGuide provides step-by-step guidance for fixing validation issues
type RemediationGuide struct {
	Steps      []string `json:"steps"`
	References []string `json:"references"`
	Examples   []string `json:"examples"`
	Automation string   `json:"automation"` // Script or command for automated fix
}

// RemediationAction represents a specific action an AI assistant can take
type RemediationAction struct {
	Action      string `json:"action"`
	Description string `json:"description"`
	Command     string `json:"command"`
	Priority    int    `json:"priority"`
}

// BypassController manages validation bypass mechanisms
type BypassController struct {
	auditTrail []BypassEvent
}

// BypassEvent represents a validation bypass event
type BypassEvent struct {
	AssistantID string    `json:"assistant_id"`
	Reason      string    `json:"reason"`
	Timestamp   time.Time `json:"timestamp"`
}

// BypassRequest represents a request to bypass validation
type BypassRequest struct {
	AssistantID   string   `json:"assistant_id"`
	Reason        string   `json:"reason"`
	FilesAffected []string `json:"files_affected"`
	Justification string   `json:"justification"`
}

// ComplianceReport represents compliance monitoring data
type ComplianceReport struct {
	AssistantID        string    `json:"assistant_id"`
	TimeRange          string    `json:"time_range"`
	TotalValidations   int       `json:"total_validations"`
	SuccessfulPasses   int       `json:"successful_passes"`
	ValidationFailures int       `json:"validation_failures"`
	BypassUsage        int       `json:"bypass_usage"`
	ComplianceScore    float64   `json:"compliance_score"`
	GeneratedAt        time.Time `json:"generated_at"`
}

// NewBypassController creates a new bypass controller
func NewBypassController() *BypassController {
	return &BypassController{
		auditTrail: make([]BypassEvent, 0),
	}
}

// NewAIValidationGateway creates a new AI validation gateway
func NewAIValidationGateway() *AIValidationGateway {
	return &AIValidationGateway{
		validator:         NewDOC008Validator(),
		errorFormatter:    NewAIErrorFormatter(),
		complianceTracker: NewComplianceMonitor(),
		bypassManager:     NewBypassController(),
	}
}

// ProcessValidationRequest processes a validation request from an AI assistant
func (gw *AIValidationGateway) ProcessValidationRequest(ctx context.Context, request ValidationRequest) (*ValidationResponse, error) {
	startTime := time.Now()

	// 🔺 DOC-011: Pre-submission validation processing - 🔧 Request validation
	if err := gw.validateRequest(request); err != nil {
		return nil, fmt.Errorf("invalid validation request: %w", err)
	}

	// 🔺 DOC-011: DOC-008 framework integration - 🔍 Core validation execution
	validationResult, err := gw.validator.ValidateFiles(ctx, request.SourceFiles, request.ValidationMode)
	if err != nil {
		return nil, fmt.Errorf("validation execution failed: %w", err)
	}

	// 🔺 DOC-011: AI-optimized error processing - 📝 Error formatting for AI consumption
	response := &ValidationResponse{
		Status:           gw.determineStatus(validationResult),
		Errors:           gw.errorFormatter.FormatErrors(validationResult.Errors),
		Warnings:         gw.errorFormatter.FormatWarnings(validationResult.Warnings),
		RemediationSteps: gw.generateRemediationSteps(validationResult),
		ComplianceScore:  gw.calculateComplianceScore(validationResult),
		ProcessingTime:   time.Since(startTime),
	}

	// 🔺 DOC-011: Compliance monitoring - 📊 Track AI assistant behavior
	gw.complianceTracker.RecordValidationEvent(ValidationEvent{
		AssistantID:     request.RequestContext.AssistantID,
		SessionID:       request.RequestContext.SessionID,
		Status:          response.Status,
		ErrorCount:      len(response.Errors),
		WarningCount:    len(response.Warnings),
		ComplianceScore: response.ComplianceScore,
		Timestamp:       time.Now(),
	})

	return response, nil
}

// ValidatePreSubmission provides the main pre-submission validation interface
func (gw *AIValidationGateway) ValidatePreSubmission(ctx context.Context, sourceFiles []string) (*ValidationResponse, error) {
	request := ValidationRequest{
		SourceFiles:    sourceFiles,
		ValidationMode: "standard",
		RequestContext: &AIRequestContext{
			AssistantID: "pre-submission",
			SessionID:   fmt.Sprintf("pre-sub-%d", time.Now().Unix()),
			Timestamp:   time.Now(),
		},
	}

	return gw.ProcessValidationRequest(ctx, request)
}

// RequestValidationBypass requests a bypass for validation failures
func (gw *AIValidationGateway) RequestValidationBypass(ctx context.Context, request BypassRequest) error {
	// 🔺 DOC-011: Bypass mechanism - 🛡️ Safe override with documentation
	if request.Reason == "" {
		return fmt.Errorf("bypass reason is required")
	}

	if request.Justification == "" {
		return fmt.Errorf("bypass justification is required")
	}

	// Record bypass event for audit trail
	bypassEvent := BypassEvent{
		AssistantID: request.AssistantID,
		Reason:      fmt.Sprintf("%s: %s", request.Reason, request.Justification),
		Timestamp:   time.Now(),
	}

	gw.bypassManager.auditTrail = append(gw.bypassManager.auditTrail, bypassEvent)

	// 📊 Track bypass usage in compliance monitoring
	gw.complianceTracker.RecordValidationEvent(ValidationEvent{
		AssistantID:     request.AssistantID,
		SessionID:       fmt.Sprintf("bypass-%d", time.Now().Unix()),
		Status:          "bypass",
		ErrorCount:      0,
		WarningCount:    0,
		ComplianceScore: 0.5, // Reduced score for bypass usage
		Timestamp:       time.Now(),
	})

	return nil
}

// GetComplianceReport generates a compliance report for an AI assistant
func (gw *AIValidationGateway) GetComplianceReport(assistantID string, timeRange string) (*ComplianceReport, error) {
	// 🔺 DOC-011: Compliance monitoring - 📊 Generate compliance reports
	dashboard := gw.complianceTracker.GenerateDashboard()

	if metrics, exists := dashboard.Metrics[assistantID]; exists {
		return &ComplianceReport{
			AssistantID:        assistantID,
			TimeRange:          timeRange,
			TotalValidations:   metrics.TotalValidations,
			SuccessfulPasses:   metrics.SuccessfulPasses,
			ValidationFailures: metrics.FailureCount,
			BypassUsage:        len(gw.getBypassEventsForAssistant(assistantID)),
			ComplianceScore:    metrics.AverageScore,
			GeneratedAt:        time.Now(),
		}, nil
	}

	return &ComplianceReport{
		AssistantID:        assistantID,
		TimeRange:          timeRange,
		TotalValidations:   0,
		SuccessfulPasses:   0,
		ValidationFailures: 0,
		BypassUsage:        0,
		ComplianceScore:    1.0,
		GeneratedAt:        time.Now(),
	}, nil
}

// GetBypassAuditTrail returns the bypass audit trail
func (gw *AIValidationGateway) GetBypassAuditTrail() []BypassEvent {
	// 🔺 DOC-011: Bypass audit trail - 📝 Comprehensive audit tracking
	return gw.bypassManager.auditTrail
}

// ValidateWithStrictMode performs validation in strict mode
func (gw *AIValidationGateway) ValidateWithStrictMode(ctx context.Context, sourceFiles []string, assistantID string) (*ValidationResponse, error) {
	// 🔺 DOC-011: Strict validation mode - 🔍 Enhanced validation for critical changes
	request := ValidationRequest{
		SourceFiles:    sourceFiles,
		ValidationMode: "strict",
		RequestContext: &AIRequestContext{
			AssistantID: assistantID,
			SessionID:   fmt.Sprintf("strict-%d", time.Now().Unix()),
			Timestamp:   time.Now(),
		},
	}

	return gw.ProcessValidationRequest(ctx, request)
}

// validateRequest validates the structure and content of a validation request
func (gw *AIValidationGateway) validateRequest(request ValidationRequest) error {
	if len(request.SourceFiles) == 0 {
		return fmt.Errorf("no source files specified")
	}

	if request.ValidationMode == "" {
		request.ValidationMode = "standard"
	}

	validModes := []string{"standard", "strict", "legacy", "realtime"}
	isValid := false
	for _, mode := range validModes {
		if request.ValidationMode == mode {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("invalid validation mode: %s (valid: %v)", request.ValidationMode, validModes)
	}

	if request.RequestContext == nil {
		request.RequestContext = &AIRequestContext{
			AssistantID: "unknown",
			SessionID:   fmt.Sprintf("auto-%d", time.Now().Unix()),
			Timestamp:   time.Now(),
		}
	}

	return nil
}

// determineStatus determines the overall validation status
func (gw *AIValidationGateway) determineStatus(result *ValidationResult) string {
	if len(result.Errors) > 0 {
		return "fail"
	}
	if len(result.Warnings) > 5 {
		return "warning"
	}
	return "pass"
}

// calculateComplianceScore calculates a compliance score based on validation results
func (gw *AIValidationGateway) calculateComplianceScore(result *ValidationResult) float64 {
	totalIssues := len(result.Errors) + len(result.Warnings)
	if totalIssues == 0 {
		return 1.0
	}

	// Weight errors more heavily than warnings
	errorWeight := 1.0
	warningWeight := 0.3

	weightedScore := (float64(len(result.Errors)) * errorWeight) + (float64(len(result.Warnings)) * warningWeight)
	maxPossibleScore := float64(result.TotalTokens)

	if maxPossibleScore == 0 {
		return 1.0
	}

	score := 1.0 - (weightedScore / maxPossibleScore)
	if score < 0 {
		score = 0
	}

	return score
}

// generateRemediationSteps creates actionable remediation steps for AI assistants
func (gw *AIValidationGateway) generateRemediationSteps(result *ValidationResult) []RemediationAction {
	var actions []RemediationAction

	// 🔺 DOC-011: Intelligent remediation generation - 🔧 AI-friendly action steps
	for _, err := range result.Errors {
		switch err.Category {
		case "token_format":
			actions = append(actions, RemediationAction{
				Action:      "fix_token_format",
				Description: "Update implementation token to use standardized format with priority icons",
				Command:     "scripts/fix-token-format.sh",
				Priority:    1,
			})
		case "missing_icon":
			actions = append(actions, RemediationAction{
				Action:      "add_priority_icon",
				Description: "Add priority icon (⭐🔺🔶🔻) to implementation token",
				Command:     "scripts/add-priority-icons.sh",
				Priority:    2,
			})
		case "documentation_sync":
			actions = append(actions, RemediationAction{
				Action:      "sync_documentation",
				Description: "Update documentation cross-references to maintain consistency",
				Command:     "make validate-icon-enforcement",
				Priority:    3,
			})
		}
	}

	return actions
}

// getBypassEventsForAssistant returns bypass events for a specific assistant
func (gw *AIValidationGateway) getBypassEventsForAssistant(assistantID string) []BypassEvent {
	var events []BypassEvent
	for _, event := range gw.bypassManager.auditTrail {
		if event.AssistantID == assistantID {
			events = append(events, event)
		}
	}
	return events
}

// ExecuteMakefileValidation executes DOC-008 validation via Makefile integration
func (gw *AIValidationGateway) ExecuteMakefileValidation(mode string) (*ValidationResponse, error) {
	var cmd *exec.Cmd

	// 🔺 DOC-011: Makefile integration - 🔧 Seamless workflow integration
	switch mode {
	case "standard":
		cmd = exec.Command("make", "validate-icon-enforcement")
	case "strict":
		cmd = exec.Command("make", "validate-icons-strict")
	case "legacy":
		cmd = exec.Command("make", "validate-icons")
	default:
		return nil, fmt.Errorf("unsupported validation mode: %s", mode)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return &ValidationResponse{
			Status: "fail",
			Errors: []AIOptimizedError{
				{
					ErrorID:  "DOC011-MAKE-001",
					Category: "makefile_execution",
					Severity: "high",
					Message:  fmt.Sprintf("Makefile validation failed: %v", err),
					Context:  map[string]string{"output": string(output)},
				},
			},
			ComplianceScore: 0.0,
		}, nil
	}

	return &ValidationResponse{
		Status:          "pass",
		ComplianceScore: 1.0,
	}, nil
}
