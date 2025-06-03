// 🔶 DOC-010: Token analysis engine - 🔧 CLI application for automated token format suggestions
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// 🔶 DOC-010: CLI application structure - 🔧 Main application framework
var (
	verbose    bool
	outputJSON bool
	configFile string
	dryRun     bool
)

// 🔶 DOC-010: Main command definition - 🔧 Primary CLI interface
var rootCmd = &cobra.Command{
	Use:   "token-suggester",
	Short: "Automated token format suggestion engine for AI assistants",
	Long: `Token Suggester analyzes Go source code and suggests appropriate 
implementation token formats following DOC-007/DOC-008 standardization.

This tool helps AI assistants create consistently formatted implementation 
tokens with correct priority icons (⭐🔺🔶🔻) and action icons (🔍📝🔧🛡️).`,
	Example: `  token-suggester analyze ./pkg/config/
  token-suggester suggest-function main.go:45
  token-suggester validate-tokens . --dry-run
  token-suggester batch-suggest . --output-json`,
}

// 🔶 DOC-010: Analysis command - 🔍 Code analysis and pattern recognition
var analyzeCmd = &cobra.Command{
	Use:   "analyze [directory|file]",
	Short: "Analyze code for token suggestion opportunities",
	Long: `Analyze Go source code to identify functions that need implementation 
tokens or have incorrectly formatted tokens. Provides suggestions based on 
function signatures, behavior patterns, and feature tracking mappings.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		analyzer := NewTokenAnalyzer()

		// 🔶 DOC-010: Pattern recognition analysis - 🔍 Code structure analysis
		if verbose {
			fmt.Printf("🔍 Analyzing %s for token suggestions...\n", args[0])
		}

		results, err := analyzer.AnalyzeTarget(args[0])
		if err != nil {
			return fmt.Errorf("analysis failed: %w", err)
		}

		// 🔶 DOC-010: Results presentation - 📝 Output formatting
		if outputJSON {
			return outputResultsJSON(results)
		}
		return outputResultsText(results)
	},
}

// 🔶 DOC-010: Function suggestion command - 🔧 Individual function analysis
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

		// 🔶 DOC-010: Function-specific analysis - 🔍 Detailed function examination
		suggestion, err := analyzer.SuggestForFunction(parts[0], parts[1])
		if err != nil {
			return fmt.Errorf("suggestion failed: %w", err)
		}

		// 🔶 DOC-010: Function suggestion output - 📝 Detailed suggestion formatting
		if outputJSON {
			return outputSuggestionJSON(suggestion)
		}
		return outputSuggestionText(suggestion)
	},
}

// 🔶 DOC-010: Token validation command - 🛡️ Token format validation
var validateCmd = &cobra.Command{
	Use:   "validate-tokens [directory]",
	Short: "Validate existing token formats and suggest improvements",
	Long: `Scan existing implementation tokens in the codebase and validate 
them against DOC-007/DOC-008 standards. Provide suggestions for improvements.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		validator := NewTokenValidator()

		// 🔶 DOC-010: Token validation process - 🛡️ Standards compliance checking
		if verbose {
			fmt.Printf("🛡️ Validating tokens in %s...\n", args[0])
		}

		violations, err := validator.ValidateTokens(args[0])
		if err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		// 🔶 DOC-010: Validation results output - 📝 Violation reporting
		if outputJSON {
			return outputViolationsJSON(violations)
		}
		return outputViolationsText(violations, dryRun)
	},
}

// 🔶 DOC-010: Batch suggestion command - 🚀 Mass suggestion processing
var batchSuggestCmd = &cobra.Command{
	Use:   "batch-suggest [directory]",
	Short: "Generate token suggestions for entire codebase",
	Long: `Perform comprehensive analysis of entire codebase to generate 
token format suggestions for all functions. Useful for large-scale 
standardization efforts.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		processor := NewBatchProcessor()

		// 🔶 DOC-010: Batch processing execution - 🚀 Comprehensive analysis
		if verbose {
			fmt.Printf("🚀 Processing batch suggestions for %s...\n", args[0])
		}

		batchResults, err := processor.ProcessDirectory(args[0])
		if err != nil {
			return fmt.Errorf("batch processing failed: %w", err)
		}

		// 🔶 DOC-010: Batch results output - 📊 Comprehensive reporting
		if outputJSON {
			return outputBatchResultsJSON(batchResults)
		}
		return outputBatchResultsText(batchResults)
	},
}

// 🔶 DOC-010: Output formatting functions - 📝 Results presentation
func outputResultsJSON(results *AnalysisResults) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

func outputResultsText(results *AnalysisResults) error {
	fmt.Printf("📊 Token Analysis Results\n")
	fmt.Printf("========================\n\n")

	fmt.Printf("📁 Analyzed: %s\n", results.Target)
	fmt.Printf("🔍 Functions analyzed: %d\n", results.FunctionsAnalyzed)
	fmt.Printf("🆕 Missing tokens: %d\n", results.MissingTokens)
	fmt.Printf("⚠️  Format violations: %d\n", results.FormatViolations)
	fmt.Printf("💡 Suggestions generated: %d\n\n", len(results.Suggestions))

	// 🔶 DOC-010: Suggestion display - 📝 Human-readable output
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

	fmt.Printf("📁 File: %s:%d\n", suggestion.FilePath, suggestion.LineNumber)
	fmt.Printf("🔧 Function: %s\n", suggestion.FunctionName)
	fmt.Printf("🎯 Feature ID: %s\n\n", suggestion.FeatureID)

	fmt.Printf("💡 Suggested Token:\n")
	fmt.Printf("   %s\n\n", suggestion.SuggestedToken)

	fmt.Printf("📊 Analysis Details:\n")
	fmt.Printf("   Priority: %s (%s)\n", suggestion.PriorityIcon, suggestion.PriorityReason)
	fmt.Printf("   Action: %s (%s)\n", suggestion.ActionIcon, suggestion.ActionReason)
	fmt.Printf("   Confidence: %.1f%%\n\n", suggestion.Confidence*100)

	fmt.Printf("🔍 Function Analysis:\n")
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
	fmt.Printf("🛡️ Token Validation Results\n")
	fmt.Printf("============================\n\n")

	if len(violations) == 0 {
		fmt.Printf("✅ No violations found - all tokens comply with standards!\n")
		return nil
	}

	fmt.Printf("⚠️  Found %d token violations:\n\n", len(violations))

	// 🔶 DOC-010: Violation categorization - 🛡️ Issue classification
	for _, violation := range violations {
		fmt.Printf("📍 %s:%d\n", violation.FilePath, violation.LineNumber)
		fmt.Printf("   Issue: %s\n", violation.ViolationType)
		fmt.Printf("   Current: %s\n", violation.CurrentToken)
		fmt.Printf("   Suggested: %s\n", violation.SuggestedFix)
		fmt.Printf("   Severity: %s\n", violation.Severity)

		if dryRun {
			fmt.Printf("   📝 DRY RUN - would apply fix\n")
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

	fmt.Printf("📁 Directory: %s\n", results.Directory)
	fmt.Printf("📊 Files processed: %d\n", results.FilesProcessed)
	fmt.Printf("🔍 Functions analyzed: %d\n", results.TotalFunctions)
	fmt.Printf("💡 Suggestions generated: %d\n", results.TotalSuggestions)
	fmt.Printf("⚠️  Violations found: %d\n\n", results.TotalViolations)

	// 🔶 DOC-010: Priority breakdown - 📊 Suggestion categorization
	fmt.Printf("🎯 Priority Breakdown:\n")
	fmt.Printf("   ⭐ Critical: %d suggestions\n", results.PriorityBreakdown.Critical)
	fmt.Printf("   🔺 High: %d suggestions\n", results.PriorityBreakdown.High)
	fmt.Printf("   🔶 Medium: %d suggestions\n", results.PriorityBreakdown.Medium)
	fmt.Printf("   🔻 Low: %d suggestions\n\n", results.PriorityBreakdown.Low)

	fmt.Printf("🔧 Action Breakdown:\n")
	fmt.Printf("   🔍 Analysis: %d suggestions\n", results.ActionBreakdown.Analysis)
	fmt.Printf("   📝 Documentation: %d suggestions\n", results.ActionBreakdown.Documentation)
	fmt.Printf("   🔧 Configuration: %d suggestions\n", results.ActionBreakdown.Configuration)
	fmt.Printf("   🛡️ Protection: %d suggestions\n\n", results.ActionBreakdown.Protection)

	// 🔶 DOC-010: Top suggestions display - 💡 Most important suggestions
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

// 🔶 DOC-010: Application initialization - 🚀 Setup and configuration
func init() {
	// 🔶 DOC-010: Global flags - ⚙️ CLI configuration
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "Output results in JSON format")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Configuration file path")

	// 🔶 DOC-010: Command-specific flags - ⚙️ Command configuration
	validateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be changed without making changes")

	// 🔶 DOC-010: Command registration - 🔧 CLI command structure
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(suggestFunctionCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(batchSuggestCmd)
}

// 🔶 DOC-010: Main application entry point - 🚀 CLI execution
func main() {
	// 🔶 DOC-010: Error handling and execution - 🛡️ Application reliability
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}
}
