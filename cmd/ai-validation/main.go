// [REQ:DOC_016] Semantic token coverage audit marker
// [ARCH:TOKEN_SYSTEM] Token system architecture
// [IMPL:TOKEN_COVERAGE_AUDIT] Applied audit annotation

// DOC-011: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
// Integrated with DOC-014: AI Assistant Decision Framework
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"bkpdir/internal/validation"

	"github.com/spf13/cobra"
)

var (
	validationMode string
	assistantID    string
	sessionID      string
	outputFormat   string
	timeout        time.Duration
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ai-validation",
	Short: "DOC-011: AI assistant validation tool [DECISION:validation]",
	Long: `AI Validation Tool for AI Assistants

This tool provides comprehensive validation services for AI assistants working with
code changes, including pre-submission validation, compliance monitoring, and
bypass mechanisms with comprehensive audit trails.

DOC-011: Token validation integration for AI assistants - Zero-friction validation workflow integration designed specifically for AI-first development environments. [DECISION:validation]`,
}

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate [files...]",
	Short: "Validate source files for AI assistant compliance",
	Long: `Validate source files for AI assistant compliance.

This command validates the specified files using the DOC-008 validation framework
with AI-optimized error reporting and remediation guidance.

DOC-011: Pre-submission validation - Zero-friction validation integration [DECISION:validation]`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// DOC-011: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
		gateway := validation.NewAIValidationGateway()

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		request := validation.ValidationRequest{
			SourceFiles:    args,
			ValidationMode: validationMode,
			RequestContext: &validation.AIRequestContext{
				AssistantID: assistantID,
				SessionID:   sessionID,
				Timestamp:   time.Now(),
			},
		}

		response, err := gateway.ProcessValidationRequest(ctx, request)
		if err != nil {
			log.Fatalf("❌ Validation failed: %v", err)
		}

		outputValidationResult(response)

		// Exit with appropriate code based on validation status
		switch response.Status {
		case "pass":
			os.Exit(0)
		case "warning":
			os.Exit(1)
		case "fail":
			os.Exit(2)
		default:
			os.Exit(3)
		}
	},
}

// preSubmitCmd represents the pre-submit validation command
var preSubmitCmd = &cobra.Command{
	Use:   "pre-submit [files...]",
	Short: "Pre-submission validation for AI assistants",
	Long: `Pre-submission validation for AI assistants.

This command provides pre-submission validation designed for AI assistant workflows,
ensuring all changes meet validation requirements before submission.

DOC-011: Pre-submission validation APIs - Zero-friction integration [DECISION:validation]`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// DOC-011: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
		gateway := validation.NewAIValidationGateway()

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		response, err := gateway.ValidatePreSubmission(ctx, args)
		if err != nil {
			log.Fatalf("❌ Pre-submission validation failed: %v", err)
		}

		outputValidationResult(response)

		if response.Status == "fail" {
			fmt.Fprintf(os.Stderr, "\n🚨 DOC-011: Pre-submission validation failed. Changes cannot be submitted.\n")
			os.Exit(2)
		}

		fmt.Printf("\n[SUCCESS] DOC-011: Pre-submission validation passed. Changes ready for submission.\n")
	},
}

// bypassCmd represents the bypass command
var bypassCmd = &cobra.Command{
	Use:   "bypass",
	Short: "Request validation bypass with audit trail",
	Long: `Request validation bypass with comprehensive audit trail.

This command allows AI assistants to request bypasses for validation failures
in exceptional cases, with mandatory documentation and audit trail recording.

DOC-011: Bypass mechanisms - Safe overrides with documentation [DECISION:validation]`,
	Run: func(cmd *cobra.Command, args []string) {
		// DOC-011: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
		reason, _ := cmd.Flags().GetString("reason")
		justification, _ := cmd.Flags().GetString("justification")
		files, _ := cmd.Flags().GetStringSlice("files")

		if reason == "" {
			log.Fatal("❌ Bypass reason is required (--reason)")
		}
		if justification == "" {
			log.Fatal("❌ Bypass justification is required (--justification)")
		}

		gateway := validation.NewAIValidationGateway()

		request := validation.BypassRequest{
			AssistantID:   assistantID,
			Reason:        reason,
			Justification: justification,
			FilesAffected: files,
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := gateway.RequestValidationBypass(ctx, request)
		if err != nil {
			log.Fatalf("❌ Bypass request failed: %v", err)
		}

		fmt.Printf("[SUCCESS] DOC-011: Validation bypass granted for assistant %s\n", assistantID)
		fmt.Printf("[NOTE] Reason: %s\n", reason)
		fmt.Printf("📋 Justification: %s\n", justification)
		fmt.Printf("DOC-011: Audit trail updated with bypass event [DECISION:validation]\n")
	},
}

// complianceCmd represents the compliance monitoring command
var complianceCmd = &cobra.Command{
	Use:   "compliance",
	Short: "Generate compliance report for AI assistant",
	Long: `Generate compliance report for AI assistant.

This command generates comprehensive compliance reports showing validation
behavior, success rates, and adherence to validation requirements.

DOC-011: Compliance monitoring - AI assistant behavior tracking [DECISION:validation]`,
	Run: func(cmd *cobra.Command, args []string) {
		// DOC-011: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
		timeRange, _ := cmd.Flags().GetString("time-range")

		gateway := validation.NewAIValidationGateway()

		report, err := gateway.GetComplianceReport(assistantID, timeRange)
		if err != nil {
			log.Fatalf("❌ Failed to generate compliance report: %v", err)
		}

		outputComplianceReport(report)
	},
}

// auditCmd represents the audit trail command
var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Display bypass audit trail",
	Long: `Display comprehensive bypass audit trail.

This command shows the complete audit trail of validation bypasses,
providing transparency and accountability for exceptional cases.

DOC-011: Bypass audit trails - Comprehensive tracking [DECISION:validation]`,
	Run: func(cmd *cobra.Command, args []string) {
		// DOC-011: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
		gateway := validation.NewAIValidationGateway()

		auditTrail := gateway.GetBypassAuditTrail()

		if len(auditTrail) == 0 {
			fmt.Println("📋 No bypass events in audit trail")
			return
		}

		fmt.Printf("[NOTE] DOC-011: Bypass Audit Trail (%d events)\n", len(auditTrail))
		fmt.Println("=" + strings.Repeat("=", 50))

		for i, event := range auditTrail {
			fmt.Printf("\n%d. Assistant: %s\n", i+1, event.AssistantID)
			fmt.Printf("   Timestamp: %s\n", event.Timestamp.Format("2006-01-02 15:04:05"))
			fmt.Printf("   Reason: %s\n", event.Reason)
		}
	},
}

// strictCmd represents the strict validation command
var strictCmd = &cobra.Command{
	Use:   "strict [files...]",
	Short: "Strict validation mode for critical changes",
	Long: `Strict validation mode for critical changes.

This command performs enhanced validation with stricter requirements,
designed for critical code changes that require higher validation standards.

DOC-011: Strict validation mode - Enhanced validation for critical changes [DECISION:validation]`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// DOC-011: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
		gateway := validation.NewAIValidationGateway()

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		response, err := gateway.ValidateWithStrictMode(ctx, args, assistantID)
		if err != nil {
			log.Fatalf("❌ Strict validation failed: %v", err)
		}

		fmt.Printf("[CHECK] DOC-011: Strict validation mode results:\n")
		outputValidationResult(response)

		// Exit with appropriate code based on validation status
		switch response.Status {
		case "pass":
			os.Exit(0)
		case "warning":
			os.Exit(1)
		case "fail":
			os.Exit(2)
		default:
			os.Exit(3)
		}
	},
}

// outputValidationResult outputs the validation result in the specified format
func outputValidationResult(response *validation.ValidationResponse) {
	switch outputFormat {
	case "json":
		// DOC-011: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
		data, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			log.Fatalf("❌ Failed to marshal JSON: %v", err)
		}
		fmt.Println(string(data))
	case "summary":
		// DOC-011: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
		fmt.Printf("[INFO] Validation Status: %s\n", response.Status)
		fmt.Printf("[INFO] Compliance Score: %.2f\n", response.ComplianceScore)
		fmt.Printf("[INFO] Processing Time: %v\n", response.ProcessingTime)
		fmt.Printf("[INFO] Errors: %d, Warnings: %d\n", len(response.Errors), len(response.Warnings))

		if len(response.RemediationSteps) > 0 {
			fmt.Printf("\n[ACTION] Remediation Steps:\n")
			for i, step := range response.RemediationSteps {
				fmt.Printf("  %d. %s: %s\n", i+1, step.Action, step.Description)
				if step.Command != "" {
					fmt.Printf("     Command: %s\n", step.Command)
				}
			}
		}
	default:
		// DOC-011: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
		fmt.Printf("DOC-011: AI Validation Results [DECISION:validation]\n")
		fmt.Println("=" + strings.Repeat("=", 40))
		fmt.Printf("Status: %s\n", response.Status)
		fmt.Printf("Compliance Score: %.2f\n", response.ComplianceScore)
		fmt.Printf("Processing Time: %v\n", response.ProcessingTime)

		if len(response.Errors) > 0 {
			fmt.Printf("\n❌ Errors (%d):\n", len(response.Errors))
			for i, err := range response.Errors {
				fmt.Printf("  %d. [%s] %s: %s\n", i+1, err.Severity, err.Category, err.Message)
				if err.FileReference != nil {
					fmt.Printf("     File: %s:%d:%d\n", err.FileReference.File, err.FileReference.Line, err.FileReference.Column)
				}
			}
		}

		if len(response.Warnings) > 0 {
			fmt.Printf("\n⚠️ Warnings (%d):\n", len(response.Warnings))
			for i, warn := range response.Warnings {
				fmt.Printf("  %d. %s: %s\n", i+1, warn.Category, warn.Message)
				if warn.FileReference != nil {
					fmt.Printf("     File: %s:%d:%d\n", warn.FileReference.File, warn.FileReference.Line, warn.FileReference.Column)
				}
			}
		}

		if len(response.RemediationSteps) > 0 {
			fmt.Printf("\n[ACTION] Remediation Steps (%d):\n", len(response.RemediationSteps))
			for i, step := range response.RemediationSteps {
				fmt.Printf("  %d. %s (Priority %d)\n", i+1, step.Description, step.Priority)
				if step.Command != "" {
					fmt.Printf("     Command: %s\n", step.Command)
				}
			}
		}
	}
}

// outputComplianceReport outputs the compliance report
func outputComplianceReport(report *validation.ComplianceReport) {
	fmt.Printf("[INFO] DOC-011: Compliance Report for %s\n", report.AssistantID)
	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Printf("Time Range: %s\n", report.TimeRange)
	fmt.Printf("Generated: %s\n", report.GeneratedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("\n📈 Statistics:\n")
	fmt.Printf("  Total Validations: %d\n", report.TotalValidations)
	fmt.Printf("  Successful Passes: %d\n", report.SuccessfulPasses)
	fmt.Printf("  Validation Failures: %d\n", report.ValidationFailures)
	fmt.Printf("  Bypass Usage: %d\n", report.BypassUsage)
	fmt.Printf("  Compliance Score: %.2f\n", report.ComplianceScore)

	if report.TotalValidations > 0 {
		successRate := float64(report.SuccessfulPasses) / float64(report.TotalValidations) * 100
		fmt.Printf("  Success Rate: %.1f%%\n", successRate)
	}
}

func init() {
	// DOC-011: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(preSubmitCmd)
	rootCmd.AddCommand(bypassCmd)
	rootCmd.AddCommand(complianceCmd)
	rootCmd.AddCommand(auditCmd)
	rootCmd.AddCommand(strictCmd)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&assistantID, "assistant-id", "ai-assistant", "AI assistant identifier")
	rootCmd.PersistentFlags().StringVar(&sessionID, "session-id", "", "Session identifier (auto-generated if not provided)")
	rootCmd.PersistentFlags().StringVar(&outputFormat, "format", "detailed", "Output format: detailed, summary, json")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 30*time.Second, "Validation timeout")

	// Validate command flags
	validateCmd.Flags().StringVar(&validationMode, "mode", "standard", "Validation mode: standard, strict, legacy")

	// Bypass command flags
	bypassCmd.Flags().String("reason", "", "Reason for bypass (required)")
	bypassCmd.Flags().String("justification", "", "Detailed justification for bypass (required)")
	bypassCmd.Flags().StringSlice("files", []string{}, "Files affected by bypass")

	// Compliance command flags
	complianceCmd.Flags().String("time-range", "24h", "Time range for compliance report")
}

func main() {
	// DOC-011: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	if sessionID == "" {
		sessionID = fmt.Sprintf("cli-%d", time.Now().Unix())
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}
}
