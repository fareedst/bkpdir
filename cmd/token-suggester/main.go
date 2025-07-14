// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
var (
	verbose    bool
	outputJSON bool
	configFile string
	dryRun     bool
)

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
var rootCmd = &cobra.Command{
	Use:   "token-suggester",
	Short: "Automated token format suggestion engine for AI assistants",
	Long: `Token Suggester analyzes Go source code and suggests appropriate 
implementation token formats following DOC-007/DOC-008 standardization.

This tool helps AI assistants create consistently formatted implementation 
tokens with correct priority levels ([CRITICAL][HIGH][MEDIUM][LOW]) and action types ([DECISION:discovery][DECISION:format-processing][DECISION:core-functionality][DECISION:validation]).`,
	Example: `  token-suggester analyze ./pkg/config/
  token-suggester suggest-function main.go:45
  token-suggester validate-tokens . --dry-run
  token-suggester batch-suggest . --output-json`,
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
var analyzeCmd = &cobra.Command{
	Use:   "analyze [directory|file]",
	Short: "Analyze code for token suggestion opportunities",
	Long: `Analyze Go source code to identify functions that need implementation 
tokens or have incorrectly formatted tokens. Provides suggestions based on 
function signatures, behavior patterns, and feature tracking mappings.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		analyzer := NewTokenAnalyzer()

		// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
		if verbose {
			fmt.Printf("[CHECK] Analyzing %s for token suggestions...\n", args[0])
		}

		results, err := analyzer.AnalyzeTarget(args[0])
		if err != nil {
			return fmt.Errorf("analysis failed: %w", err)
		}

		// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
		if outputJSON {
			return outputResultsJSON(results)
		}
		return outputResultsText(results)
	},
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
var suggestFunctionCmd = &cobra.Command{
	Use:   "suggest-function [file:line]",
	Short: "Suggest token format for specific function",
	Long: `Analyze a specific function at the given file and line number to 
provide detailed token format suggestions including priority and action icons.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		parts := strings.Split(args[0], ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid format: expected file:line, got %s", args[0])
		}

		analyzer := NewTokenAnalyzer()

		// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
		suggestion, err := analyzer.SuggestForFunction(parts[0], parts[1])
		if err != nil {
			return fmt.Errorf("suggestion failed: %w", err)
		}

		// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
		if outputJSON {
			return outputSuggestionJSON(suggestion)
		}
		return outputSuggestionText(suggestion)
	},
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
var validateCmd = &cobra.Command{
	Use:   "validate-tokens [directory]",
	Short: "Validate existing token formats and suggest improvements",
	Long: `Scan existing implementation tokens in the codebase and validate 
them against DOC-007/DOC-008 standards. Provide suggestions for improvements.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		validator := NewTokenValidator()

		// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
		if verbose {
			fmt.Printf("[ACTION:validation] Validating tokens in %s...\n", args[0])
		}

		violations, err := validator.ValidateTokens(args[0])
		if err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
		if outputJSON {
			return outputViolationsJSON(violations)
		}
		return outputViolationsText(violations, dryRun)
	},
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
var batchSuggestCmd = &cobra.Command{
	Use:   "batch-suggest [directory]",
	Short: "Generate token suggestions for entire codebase",
	Long: `Perform comprehensive analysis of entire codebase to generate 
token format suggestions for all functions. Useful for large-scale 
standardization efforts.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		processor := NewBatchProcessor()

		// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
		if verbose {
			fmt.Printf("🚀 Processing batch suggestions for %s...\n", args[0])
		}

		batchResults, err := processor.ProcessDirectory(args[0])
		if err != nil {
			return fmt.Errorf("batch processing failed: %w", err)
		}

		// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
		if outputJSON {
			return outputBatchResultsJSON(batchResults)
		}
		return outputBatchResultsText(batchResults)
	},
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func outputResultsJSON(results *AnalysisResults) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

func outputResultsText(results *AnalysisResults) error {
	fmt.Printf("[INFO] Token Analysis Results\n")
	fmt.Printf("========================\n\n")

	fmt.Printf("[DIR] Analyzed: %s\n", results.Target)
	fmt.Printf("[CHECK] Functions analyzed: %d\n", results.FunctionsAnalyzed)
	fmt.Printf("🆕 Missing tokens: %d\n", results.MissingTokens)
	fmt.Printf("⚠️  Format violations: %d\n", results.FormatViolations)
	fmt.Printf("💡 Suggestions generated: %d\n\n", len(results.Suggestions))

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	for _, suggestion := range results.Suggestions {
		fmt.Printf("📍 %s:%d\n", suggestion.FilePath, suggestion.LineNumber)
		fmt.Printf("   Function: %s\n", suggestion.FunctionName)
		fmt.Printf("   Priority: %s (%s)\n", suggestion.PriorityIcon, suggestion.PriorityReason)
		fmt.Printf("   Action: %s (%s)\n", suggestion.ActionIcon, suggestion.ActionReason)
		fmt.Printf("   Suggested: %s\n", suggestion.SuggestedToken)
		fmt.Printf("   Confidence: %.1f%%\n\n", suggestion.Confidence*100)
	}

	return nil
}

func outputSuggestionJSON(suggestion *TokenSuggestion) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(suggestion)
}

func outputSuggestionText(suggestion *TokenSuggestion) error {
	fmt.Printf("🎯 Token Suggestion for %s\n", suggestion.FunctionName)
	fmt.Printf("================================\n\n")

	fmt.Printf("[DIR] File: %s:%d\n", suggestion.FilePath, suggestion.LineNumber)
	fmt.Printf("[ACTION] Function: %s\n", suggestion.FunctionName)
	fmt.Printf("🎯 Feature ID: %s\n\n", suggestion.FeatureID)

	fmt.Printf("💡 Suggested Token:\n")
	fmt.Printf("   %s\n\n", suggestion.SuggestedToken)

	fmt.Printf("[INFO] Analysis Details:\n")
	fmt.Printf("   Priority: %s (%s)\n", suggestion.PriorityIcon, suggestion.PriorityReason)
	fmt.Printf("   Action: %s (%s)\n", suggestion.ActionIcon, suggestion.ActionReason)
	fmt.Printf("   Confidence: %.1f%%\n\n", suggestion.Confidence*100)

	fmt.Printf("[CHECK] Function Analysis:\n")
	fmt.Printf("   Return Type: %s\n", suggestion.FunctionSignature.ReturnType)
	fmt.Printf("   Parameters: %d\n", len(suggestion.FunctionSignature.Parameters))
	fmt.Printf("   Complexity: %s\n", suggestion.ComplexityLevel)

	return nil
}

func outputViolationsJSON(violations []TokenViolation) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(violations)
}

func outputViolationsText(violations []TokenViolation, dryRun bool) error {
	fmt.Printf("[ACTION:validation] Token Validation Results\n")
	fmt.Printf("============================\n\n")

	if len(violations) == 0 {
		fmt.Printf("[SUCCESS] No violations found - all tokens comply with standards!\n")
		return nil
	}

	fmt.Printf("⚠️  Found %d token violations:\n\n", len(violations))

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	for _, violation := range violations {
		fmt.Printf("📍 %s:%d\n", violation.FilePath, violation.LineNumber)
		fmt.Printf("   Issue: %s\n", violation.ViolationType)
		fmt.Printf("   Current: %s\n", violation.CurrentToken)
		fmt.Printf("   Suggested: %s\n", violation.SuggestedFix)
		fmt.Printf("   Severity: %s\n", violation.Severity)

		if dryRun {
			fmt.Printf("   [NOTE] DRY RUN - would apply fix\n")
		}
		fmt.Printf("\n")
	}

	return nil
}

func outputBatchResultsJSON(results *BatchResults) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

func outputBatchResultsText(results *BatchResults) error {
	fmt.Printf("🚀 Batch Processing Results\n")
	fmt.Printf("===========================\n\n")

	fmt.Printf("[DIR] Directory: %s\n", results.Directory)
	fmt.Printf("[INFO] Files processed: %d\n", results.FilesProcessed)
	fmt.Printf("[CHECK] Functions analyzed: %d\n", results.TotalFunctions)
	fmt.Printf("💡 Suggestions generated: %d\n", results.TotalSuggestions)
	fmt.Printf("⚠️  Violations found: %d\n\n", results.TotalViolations)

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	fmt.Printf("🎯 Priority Breakdown:\n")
	fmt.Printf("   [CRITICAL]: %d suggestions\n", results.PriorityBreakdown.Critical)
	fmt.Printf("   [HIGH]: %d suggestions\n", results.PriorityBreakdown.High)
	fmt.Printf("   [MEDIUM]: %d suggestions\n", results.PriorityBreakdown.Medium)
	fmt.Printf("   [LOW]: %d suggestions\n\n", results.PriorityBreakdown.Low)

	fmt.Printf("[ACTION] Action Breakdown:\n")
	fmt.Printf("   [CHECK] Analysis: %d suggestions\n", results.ActionBreakdown.Analysis)
	fmt.Printf("   [NOTE] Documentation: %d suggestions\n", results.ActionBreakdown.Documentation)
	fmt.Printf("   [ACTION] Configuration: %d suggestions\n", results.ActionBreakdown.Configuration)
	fmt.Printf("   [ACTION:validation] Protection: %d suggestions\n\n", results.ActionBreakdown.Protection)

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	if len(results.TopSuggestions) > 0 {
		fmt.Printf("💡 Top Suggestions (by confidence):\n")
		for i, suggestion := range results.TopSuggestions {
			if i >= 5 { // Limit to top 5
				break
			}
			fmt.Printf("   %d. %s:%d - %s (%.1f%%)\n",
				i+1, filepath.Base(suggestion.FilePath), suggestion.LineNumber,
				suggestion.FunctionName, suggestion.Confidence*100)
		}
	}

	return nil
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func init() {
	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "Output results in JSON format")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Configuration file path")

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	validateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be changed without making changes")

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(suggestFunctionCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(batchSuggestCmd)
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func main() {
	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}
}
