// 🔶 DOC-012: Real-time validation CLI - ⚡ Live validation service command-line interface
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"bkpdir/internal/validation"

	"github.com/spf13/cobra"
)

var (
	serverPort    int
	outputFormat  string
	watchMode     bool
	subscribeMode bool
	metricsMode   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "realtime-validator",
	Short: "🔶 DOC-012: Real-time icon validation service",
	Long: `Real-time Icon Validation Service

This service provides live validation feedback with sub-second response times
for icon compliance and intelligent correction suggestions. Part of the DOC-012
Real-time Icon Validation Feedback system.

🔶 DOC-012: Real-time validation - Enhanced developer experience with immediate feedback`,
}

// serverCmd starts the real-time validation server
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "🔧 Start real-time validation HTTP server",
	Long: `Start the real-time validation HTTP server.

This command starts an HTTP server that provides real-time validation APIs
for editors and development tools to integrate with.

🔶 DOC-012: Live validation service - HTTP API server for real-time validation`,
	Run: func(cmd *cobra.Command, args []string) {
		// 🔶 DOC-012: Server startup - 🔧 HTTP service initialization
		fmt.Printf("🔶 DOC-012: Starting real-time validation server on port %d\n", serverPort)

		validator := validation.NewRealTimeValidator()

		// Set up graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			if err := validator.StartValidationServer(serverPort); err != nil {
				log.Fatalf("❌ Failed to start validation server: %v", err)
			}
		}()

		fmt.Printf("✅ Real-time validation server running at http://localhost:%d\n", serverPort)
		fmt.Println("📊 Available endpoints:")
		fmt.Printf("  POST /validate - Real-time file validation\n")
		fmt.Printf("  GET  /subscribe - WebSocket subscription\n")
		fmt.Printf("  GET  /status - Validation status indicators\n")
		fmt.Printf("  GET  /suggestions - Intelligent suggestions\n")
		fmt.Printf("  GET  /metrics - Performance metrics\n")
		fmt.Println("\n🔶 DOC-012: Server ready for real-time validation requests")

		// Wait for shutdown signal
		<-sigChan
		fmt.Println("\n🔶 DOC-012: Shutting down real-time validation server...")
	},
}

// validateCmd validates files in real-time
var validateCmd = &cobra.Command{
	Use:   "validate [files...]",
	Short: "🔍 Validate files with real-time feedback",
	Long: `Validate files with real-time feedback and intelligent suggestions.

This command provides immediate validation feedback for the specified files
with intelligent correction suggestions and visual status indicators.

🔶 DOC-012: Real-time validation - Immediate feedback for development workflow`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 🔶 DOC-012: File validation - 🔍 Real-time validation processing
		validator := validation.NewRealTimeValidator()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		for _, filePath := range args {
			content, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("❌ Error reading file %s: %v", filePath, err)
				continue
			}

			update, err := validator.ValidateRealtimeFile(ctx, filePath, string(content))
			if err != nil {
				log.Printf("❌ Validation failed for %s: %v", filePath, err)
				continue
			}

			outputValidationUpdate(filePath, update)
		}
	},
}

// statusCmd shows validation status indicators
var statusCmd = &cobra.Command{
	Use:   "status [files...]",
	Short: "📊 Show validation status indicators",
	Long: `Show visual validation status indicators for files.

This command displays compliance levels, error counts, and visual feedback
elements for the specified files.

🔶 DOC-012: Status indicators - Visual compliance feedback`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 🔶 DOC-012: Status display - 📊 Visual status indicators
		validator := validation.NewRealTimeValidator()

		indicator := validator.GetValidationStatusIndicator(args)
		outputStatusIndicator(indicator)
	},
}

// watchCmd watches files for changes and provides real-time validation
var watchCmd = &cobra.Command{
	Use:   "watch [files...]",
	Short: "👁️ Watch files for real-time validation",
	Long: `Watch files for changes and provide real-time validation feedback.

This command monitors the specified files for changes and automatically
validates them, providing immediate feedback during development.

🔶 DOC-012: Editor integration - File watching for live feedback`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 🔶 DOC-012: File watching - 👁️ Live monitoring and validation
		fmt.Printf("🔶 DOC-012: Starting file watch mode for %d files\n", len(args))
		fmt.Println("📝 Watching for changes... (Press Ctrl+C to stop)")

		validator := validation.NewRealTimeValidator()

		// Subscribe to validation updates
		updates := validator.SubscribeToValidation("watch-cli", args)

		// Set up graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			for update := range updates {
				fmt.Printf("\n🔄 File changed: %s\n", update.File)
				outputValidationUpdate(update.File, update)
			}
		}()

		// Simulate file watching (in real implementation, use file system watcher)
		fmt.Println("⚠️  Note: File system watching not implemented in this demo")
		fmt.Println("🔶 DOC-012: Watch mode ready - subscription active")

		// Wait for shutdown
		<-sigChan
		fmt.Println("\n🔶 DOC-012: Stopping watch mode...")
		validator.UnsubscribeFromValidation("watch-cli")
	},
}

// metricsCmd shows performance metrics
var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "📈 Show real-time validation performance metrics",
	Long: `Show performance metrics for the real-time validation system.

This command displays validation performance statistics including response
times, cache hit rates, and throughput metrics.

🔶 DOC-012: Performance optimization - Metrics monitoring`,
	Run: func(cmd *cobra.Command, args []string) {
		// 🔶 DOC-012: Metrics display - 📈 Performance monitoring
		fmt.Println("🔶 DOC-012: Real-time Validation Performance Metrics")
		fmt.Println(strings.Repeat("=", 55))

		validator := validation.NewRealTimeValidator()

		// Note: In real implementation, connect to running server for metrics
		fmt.Println("📊 Current Session Metrics:")
		fmt.Println("  Total Validations: 0")
		fmt.Println("  Average Response Time: N/A")
		fmt.Println("  Cache Hit Rate: N/A")
		fmt.Println("  Active Subscribers: 0")
		fmt.Println("  Validations/Second: N/A")
		fmt.Println()
		fmt.Println("⚡ Target Performance Goals:")
		fmt.Println("  Response Time: < 1 second")
		fmt.Println("  Cache Hit Rate: > 80%")
		fmt.Println("  Availability: 99.9%")
		fmt.Println()
		fmt.Printf("🔶 DOC-012: Metrics available at server endpoint /metrics\n")

		_ = validator // Prevent unused variable warning
	},
}

// Helper functions for output formatting
func outputValidationUpdate(filePath string, update *validation.RealTimeValidationUpdate) {
	// 🔶 DOC-012: Output formatting - 📝 Structured validation results
	fmt.Printf("\n📄 File: %s\n", filePath)
	fmt.Printf("⏱️  Processing Time: %v\n", update.ProcessingTime)

	// Status with visual indicators
	if update.StatusIndicator != nil {
		fmt.Printf("📊 Status: %s %s (%s compliance)\n",
			update.StatusIndicator.VisualElements.StatusIcon,
			strings.ToUpper(update.StatusIndicator.OverallStatus),
			update.StatusIndicator.ComplianceLevel)

		if update.StatusIndicator.ErrorCount > 0 {
			fmt.Printf("❌ Errors: %d\n", update.StatusIndicator.ErrorCount)
		}
		if update.StatusIndicator.WarningCount > 0 {
			fmt.Printf("⚠️  Warnings: %d\n", update.StatusIndicator.WarningCount)
		}
	}

	// Output format selection
	switch outputFormat {
	case "json":
		data, _ := json.MarshalIndent(update, "", "  ")
		fmt.Println(string(data))
	case "summary":
		outputSummaryFormat(update)
	default:
		outputDetailedFormat(update)
	}
}

func outputSummaryFormat(update *validation.RealTimeValidationUpdate) {
	// 🔶 DOC-012: Output formatting - 📋 Summary format display
	fmt.Printf("🔍 Validation Summary:\n")
	fmt.Printf("  Status: %s\n", update.Status)
	fmt.Printf("  Errors: %d, Warnings: %d\n", len(update.Errors), len(update.Warnings))
	fmt.Printf("  Suggestions: %d\n", len(update.Suggestions))

	if len(update.Suggestions) > 0 {
		fmt.Println("💡 Top Suggestions:")
		for i, suggestion := range update.Suggestions[:min(3, len(update.Suggestions))] {
			fmt.Printf("  %d. %s (%.0f%% confidence)\n", i+1, suggestion.Suggested, suggestion.Confidence*100)
		}
	}
}

func outputDetailedFormat(update *validation.RealTimeValidationUpdate) {
	// 🔶 DOC-012: Output formatting - 📋 Detailed format display
	// Errors
	if len(update.Errors) > 0 {
		fmt.Println("\n❌ Validation Errors:")
		for i, err := range update.Errors {
			fmt.Printf("  %d. [%s] %s\n", i+1, err.Category, err.Message)
			if err.FileReference != nil {
				fmt.Printf("     📍 Line %d, Column %d\n", err.FileReference.Line, err.FileReference.Column)
			}
		}
	}

	// Warnings
	if len(update.Warnings) > 0 {
		fmt.Println("\n⚠️  Validation Warnings:")
		for i, warning := range update.Warnings {
			fmt.Printf("  %d. [%s] %s\n", i+1, warning.Category, warning.Message)
		}
	}

	// Intelligent suggestions
	if len(update.Suggestions) > 0 {
		fmt.Println("\n💡 Intelligent Suggestions:")
		for i, suggestion := range update.Suggestions {
			fmt.Printf("  %d. %s\n", i+1, suggestion.Suggested)
			fmt.Printf("     🎯 Confidence: %.0f%% | Type: %s\n", suggestion.Confidence*100, suggestion.Type)
			if suggestion.AutoApply {
				fmt.Printf("     ⚡ Auto-apply available\n")
			}
		}
	}
}

func outputStatusIndicator(indicator *validation.ValidationStatusIndicator) {
	// 🔶 DOC-012: Status display - 📊 Visual status presentation
	fmt.Println("🔶 DOC-012: Validation Status Indicators")
	fmt.Println(strings.Repeat("=", 45))

	if indicator.VisualElements != nil {
		fmt.Printf("🎨 Visual Status: %s %s\n",
			indicator.VisualElements.StatusIcon,
			strings.ToUpper(indicator.OverallStatus))
		fmt.Printf("📊 Compliance Level: %s (%d%%)\n",
			indicator.ComplianceLevel,
			indicator.VisualElements.ProgressBar)
		fmt.Printf("🏷️  Badge: %s\n", indicator.VisualElements.BadgeText)
		fmt.Printf("💡 Tooltip: %s\n", indicator.VisualElements.Tooltip)
	}

	fmt.Printf("\n📈 Summary:\n")
	fmt.Printf("  Errors: %d\n", indicator.ErrorCount)
	fmt.Printf("  Warnings: %d\n", indicator.WarningCount)
	fmt.Printf("  Overall: %s\n", indicator.OverallStatus)

	if len(indicator.FileStatus) > 0 {
		fmt.Printf("\n📄 File Status:\n")
		for file, status := range indicator.FileStatus {
			fmt.Printf("  %s: %s\n", file, status)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func init() {
	// 🔶 DOC-012: CLI initialization - 🔧 Command setup
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(watchCmd)
	rootCmd.AddCommand(metricsCmd)

	// Server command flags
	serverCmd.Flags().IntVarP(&serverPort, "port", "p", 8080, "Server port")

	// Validate command flags
	validateCmd.Flags().StringVarP(&outputFormat, "format", "f", "detailed", "Output format (detailed, summary, json)")

	// Watch command flags
	watchCmd.Flags().BoolVarP(&watchMode, "continuous", "c", false, "Continuous monitoring mode")

	// Status command flags
	statusCmd.Flags().StringVarP(&outputFormat, "format", "f", "detailed", "Output format (detailed, summary, json)")

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&metricsMode, "metrics", "m", false, "Show metrics with output")
}

func main() {
	// 🔶 DOC-012: CLI main - ⚡ Real-time validator entry point
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}
}
