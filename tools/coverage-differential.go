// 🔺 COV-002: Differential coverage reporting tool - 🔧
// Compares current coverage against baseline and reports changes for modified code only
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// 🔺 COV-002: Configuration structure for differential reporting - 🔧
type DifferentialConfig struct {
	General struct {
		Version          string  `toml:"version"`
		BaselineDate     string  `toml:"baseline_date"`
		BaselineCoverage float64 `toml:"baseline_coverage"`
		DifferentialMode bool    `toml:"differential_mode"`
		TrendTracking    bool    `toml:"trend_tracking"`
	} `toml:"general"`

	QualityGates struct {
		NewCodeThreshold       float64 `toml:"new_code_threshold"`
		CriticalPathThreshold  float64 `toml:"critical_path_threshold"`
		LegacyPreservation     bool    `toml:"legacy_preservation"`
		OverallRegressionLimit float64 `toml:"overall_regression_limit"`
	} `toml:"quality_gates"`

	Differential struct {
		Enabled              bool   `toml:"enabled"`
		CompareMode          string `toml:"compare_mode"`
		IncludeLineChanges   bool   `toml:"include_line_changes"`
		ExcludePureAdditions bool   `toml:"exclude_pure_additions"`
		FocusOnModified      bool   `toml:"focus_on_modified"`
	} `toml:"differential"`

	Exclusions struct {
		ExcludeFiles     []string `toml:"exclude_files"`
		ExcludeFunctions []string `toml:"exclude_functions"`
	} `toml:"exclusions"`

	CriticalFunctions struct {
		CriticalPaths []string `toml:"critical_paths"`
	} `toml:"critical_functions"`

	Reporting struct {
		GenerateHTML               bool   `toml:"generate_html"`
		GenerateJSON               bool   `toml:"generate_json"`
		GenerateDiffHTML           bool   `toml:"generate_diff_html"`
		GenerateDiffJSON           bool   `toml:"generate_diff_json"`
		GenerateBaselineComparison bool   `toml:"generate_baseline_comparison"`
		OutputDir                  string `toml:"output_dir"`
		DifferentialReport         string `toml:"differential_report"`
		TrendReport                string `toml:"trend_report"`
	} `toml:"reporting"`
}

// 🔺 COV-002: Coverage data structures for differential analysis - 🔧
type FunctionCoverage struct {
	File       string  `json:"file"`
	Function   string  `json:"function"`
	Line       int     `json:"line"`
	Coverage   float64 `json:"coverage"`
	Statements int     `json:"statements"`
	Covered    int     `json:"covered"`
}

type FileCoverage struct {
	File         string             `json:"file"`
	Coverage     float64            `json:"coverage"`
	Functions    []FunctionCoverage `json:"functions"`
	TotalStmts   int                `json:"total_statements"`
	CoveredStmts int                `json:"covered_statements"`
}

type BaselineCoverage struct {
	Timestamp       time.Time               `json:"timestamp"`
	OverallCoverage float64                 `json:"overall_coverage"`
	Files           map[string]FileCoverage `json:"files"`
	GitCommit       string                  `json:"git_commit,omitempty"`
	GitBranch       string                  `json:"git_branch,omitempty"`
}

type DifferentialReport struct {
	Timestamp          time.Time             `json:"timestamp"`
	BaselineCommit     string                `json:"baseline_commit"`
	CurrentCommit      string                `json:"current_commit"`
	ModifiedFiles      []string              `json:"modified_files"`
	OverallChange      float64               `json:"overall_change"`
	NewCodeCoverage    float64               `json:"new_code_coverage"`
	QualityGatesPassed bool                  `json:"quality_gates_passed"`
	FileChanges        map[string]FileChange `json:"file_changes"`
	CriticalPathStatus map[string]bool       `json:"critical_path_status"`
	Recommendations    []string              `json:"recommendations"`
}

type FileChange struct {
	File             string           `json:"file"`
	BaselineCoverage float64          `json:"baseline_coverage"`
	CurrentCoverage  float64          `json:"current_coverage"`
	Change           float64          `json:"change"`
	IsNewFile        bool             `json:"is_new_file"`
	IsModified       bool             `json:"is_modified"`
	MeetsThreshold   bool             `json:"meets_threshold"`
	FunctionChanges  []FunctionChange `json:"function_changes,omitempty"`
}

type FunctionChange struct {
	Function         string  `json:"function"`
	BaselineCoverage float64 `json:"baseline_coverage"`
	CurrentCoverage  float64 `json:"current_coverage"`
	Change           float64 `json:"change"`
	IsCriticalPath   bool    `json:"is_critical_path"`
	MeetsThreshold   bool    `json:"meets_threshold"`
}

// 🔺 COV-002: Git integration for detecting modified files - 🔍
func getModifiedFiles() ([]string, error) {
	// Get list of modified files since last commit
	cmd := exec.Command("git", "diff", "--name-only", "HEAD~1")
	output, err := cmd.Output()
	if err != nil {
		// Fall back to comparing with main branch
		cmd = exec.Command("git", "diff", "--name-only", "origin/main")
		output, err = cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to get modified files: %v", err)
		}
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	var goFiles []string
	for _, file := range files {
		if strings.HasSuffix(file, ".go") && !strings.HasSuffix(file, "_test.go") {
			goFiles = append(goFiles, file)
		}
	}

	return goFiles, nil
}

// 🔺 COV-002: Parse current coverage profile - 🔍
func parseCurrentCoverage(profilePath string) (map[string]FileCoverage, error) {
	file, err := os.Open(profilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open coverage profile: %v", err)
	}
	defer file.Close()

	files := make(map[string]FileCoverage)
	funcRegex := regexp.MustCompile(`^(.+):(\d+):\s*(.+)\s+([\d.]+)%$`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "mode:") {
			continue
		}

		matches := funcRegex.FindStringSubmatch(line)
		if len(matches) == 5 {
			fileName := matches[1]
			lineNum, _ := strconv.Atoi(matches[2])
			functionName := matches[3]
			coverage, _ := strconv.ParseFloat(matches[4], 64)

			if _, exists := files[fileName]; !exists {
				files[fileName] = FileCoverage{
					File:      fileName,
					Functions: []FunctionCoverage{},
				}
			}

			fileCov := files[fileName]
			fileCov.Functions = append(fileCov.Functions, FunctionCoverage{
				File:     fileName,
				Function: functionName,
				Line:     lineNum,
				Coverage: coverage,
			})
			files[fileName] = fileCov
		}
	}

	// Calculate file-level coverage averages
	for fileName, fileCov := range files {
		if len(fileCov.Functions) > 0 {
			total := 0.0
			for _, fn := range fileCov.Functions {
				total += fn.Coverage
			}
			fileCov.Coverage = total / float64(len(fileCov.Functions))
			files[fileName] = fileCov
		}
	}

	return files, scanner.Err()
}

// 🔺 COV-002: Load baseline coverage from baseline file - 🔧
func loadBaselineCoverage() (BaselineCoverage, error) {
	var baseline BaselineCoverage

	// Try to load from JSON history file first
	historyPath := "docs/coverage-history.json"
	if data, err := os.ReadFile(historyPath); err == nil {
		var histories []BaselineCoverage
		if json.Unmarshal(data, &histories) == nil && len(histories) > 0 {
			// Use the most recent baseline
			baseline = histories[len(histories)-1]
			return baseline, nil
		}
	}

	// Fall back to parsing baseline documentation
	baselinePath := "docs/coverage-baseline.md"
	data, err := os.ReadFile(baselinePath)
	if err != nil {
		return baseline, fmt.Errorf("failed to read baseline file: %v", err)
	}

	// Parse baseline coverage from markdown
	content := string(data)
	overallRegex := regexp.MustCompile(`\*\*Overall Coverage:\*\*\s+` + "`" + `([0-9.]+)%` + "`")
	if matches := overallRegex.FindStringSubmatch(content); len(matches) > 1 {
		if coverage, err := strconv.ParseFloat(matches[1], 64); err == nil {
			baseline.OverallCoverage = coverage
		}
	}

	baseline.Timestamp = time.Now()
	baseline.Files = make(map[string]FileCoverage)

	return baseline, nil
}

// 🔺 COV-002: Generate differential report - 🔧
func generateDifferentialReport(config DifferentialConfig, currentCoverage map[string]FileCoverage, baseline BaselineCoverage, modifiedFiles []string) DifferentialReport {
	report := DifferentialReport{
		Timestamp:          time.Now(),
		ModifiedFiles:      modifiedFiles,
		FileChanges:        make(map[string]FileChange),
		CriticalPathStatus: make(map[string]bool),
		Recommendations:    []string{},
	}

	// Get current overall coverage
	currentOverall := calculateOverallCoverage(currentCoverage)
	report.OverallChange = currentOverall - baseline.OverallCoverage

	// Calculate new code coverage (modified files only)
	if len(modifiedFiles) > 0 {
		newCodeTotal := 0.0
		count := 0
		for _, file := range modifiedFiles {
			if fileCov, exists := currentCoverage[file]; exists {
				newCodeTotal += fileCov.Coverage
				count++
			}
		}
		if count > 0 {
			report.NewCodeCoverage = newCodeTotal / float64(count)
		}
	} else {
		report.NewCodeCoverage = currentOverall
	}

	// Analyze file changes
	for _, file := range modifiedFiles {
		change := FileChange{
			File:       file,
			IsModified: true,
		}

		if currentCov, exists := currentCoverage[file]; exists {
			change.CurrentCoverage = currentCov.Coverage
			change.MeetsThreshold = currentCov.Coverage >= config.QualityGates.NewCodeThreshold

			if baselineCov, baselineExists := baseline.Files[file]; baselineExists {
				change.BaselineCoverage = baselineCov.Coverage
				change.Change = currentCov.Coverage - baselineCov.Coverage
			} else {
				change.IsNewFile = true
				change.BaselineCoverage = 0
				change.Change = currentCov.Coverage
			}

			// Analyze function-level changes
			change.FunctionChanges = analyzeFunctionChanges(currentCov, baseline.Files[file], config.CriticalFunctions.CriticalPaths)
		}

		report.FileChanges[file] = change
	}

	// Check critical path coverage
	for _, criticalPath := range config.CriticalFunctions.CriticalPaths {
		parts := strings.Split(criticalPath, ":")
		if len(parts) == 2 {
			file := parts[0]
			function := parts[1]

			if fileCov, exists := currentCoverage[file]; exists {
				for _, fn := range fileCov.Functions {
					if fn.Function == function {
						report.CriticalPathStatus[criticalPath] = fn.Coverage >= config.QualityGates.CriticalPathThreshold
						break
					}
				}
			}
		}
	}

	// Evaluate quality gates
	report.QualityGatesPassed = evaluateQualityGates(config, report)

	// Generate recommendations
	report.Recommendations = generateRecommendations(config, report)

	return report
}

// 🔺 COV-002: Calculate overall coverage from file coverage data - 🔧
func calculateOverallCoverage(files map[string]FileCoverage) float64 {
	if len(files) == 0 {
		return 0.0
	}

	total := 0.0
	for _, file := range files {
		total += file.Coverage
	}
	return total / float64(len(files))
}

// 🔺 COV-002: Analyze function-level changes - 🔍
func analyzeFunctionChanges(currentFile FileCoverage, baselineFile FileCoverage, criticalPaths []string) []FunctionChange {
	var changes []FunctionChange

	// Create baseline function map
	baselineFuncs := make(map[string]float64)
	for _, fn := range baselineFile.Functions {
		baselineFuncs[fn.Function] = fn.Coverage
	}

	// Check critical paths
	criticalSet := make(map[string]bool)
	for _, path := range criticalPaths {
		parts := strings.Split(path, ":")
		if len(parts) == 2 && parts[0] == currentFile.File {
			criticalSet[parts[1]] = true
		}
	}

	// Analyze current functions
	for _, fn := range currentFile.Functions {
		change := FunctionChange{
			Function:        fn.Function,
			CurrentCoverage: fn.Coverage,
			IsCriticalPath:  criticalSet[fn.Function],
		}

		if baselineCov, exists := baselineFuncs[fn.Function]; exists {
			change.BaselineCoverage = baselineCov
			change.Change = fn.Coverage - baselineCov
		} else {
			change.BaselineCoverage = 0
			change.Change = fn.Coverage
		}

		threshold := 70.0 // Default threshold
		if change.IsCriticalPath {
			threshold = 80.0 // Critical path threshold
		}
		change.MeetsThreshold = fn.Coverage >= threshold

		changes = append(changes, change)
	}

	return changes
}

// 🔺 COV-002: Evaluate quality gates - 🔍
func evaluateQualityGates(config DifferentialConfig, report DifferentialReport) bool {
	// Check overall regression limit
	if report.OverallChange < config.QualityGates.OverallRegressionLimit {
		return false
	}

	// Check new code threshold
	if report.NewCodeCoverage < config.QualityGates.NewCodeThreshold {
		return false
	}

	// Check critical path thresholds
	for _, passed := range report.CriticalPathStatus {
		if !passed {
			return false
		}
	}

	// Check individual file thresholds for modified files
	for _, change := range report.FileChanges {
		if change.IsModified && !change.MeetsThreshold {
			return false
		}
	}

	return true
}

// 🔺 COV-002: Generate recommendations based on analysis - 📝
func generateRecommendations(config DifferentialConfig, report DifferentialReport) []string {
	var recommendations []string

	if report.OverallChange < 0 {
		recommendations = append(recommendations, fmt.Sprintf("Overall coverage decreased by %.1f%%. Consider adding tests for modified code.", -report.OverallChange))
	}

	if report.NewCodeCoverage < config.QualityGates.NewCodeThreshold {
		recommendations = append(recommendations, fmt.Sprintf("New code coverage (%.1f%%) is below threshold (%.1f%%). Add tests for modified functions.", report.NewCodeCoverage, config.QualityGates.NewCodeThreshold))
	}

	for path, passed := range report.CriticalPathStatus {
		if !passed {
			recommendations = append(recommendations, fmt.Sprintf("Critical path %s does not meet required threshold. Increase test coverage for this function.", path))
		}
	}

	lowCoverageFiles := 0
	for _, change := range report.FileChanges {
		if change.IsModified && !change.MeetsThreshold {
			lowCoverageFiles++
		}
	}

	if lowCoverageFiles > 0 {
		recommendations = append(recommendations, fmt.Sprintf("%d modified files have coverage below threshold. Focus testing efforts on these files.", lowCoverageFiles))
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "All quality gates passed! Coverage changes look good.")
	}

	return recommendations
}

// 🔺 COV-002: Generate HTML differential report - 📝
func generateHTMLReport(report DifferentialReport, outputPath string) error {
	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Differential Coverage Report (COV-002)</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background: #f0f0f0; padding: 20px; border-radius: 5px; }
        .summary { display: flex; gap: 20px; margin: 20px 0; }
        .metric { background: #e8f4fd; padding: 15px; border-radius: 5px; flex: 1; }
        .pass { background: #d4edda; }
        .fail { background: #f8d7da; }
        table { width: 100%%; border-collapse: collapse; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background: #f2f2f2; }
        .recommendations { background: #fff3cd; padding: 15px; border-radius: 5px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Differential Coverage Report (COV-002)</h1>
        <p>Generated: %s</p>
        <p>Quality Gates: <span class="%s">%s</span></p>
    </div>

    <div class="summary">
        <div class="metric">
            <h3>Overall Change</h3>
            <p><strong>%.2f%%</strong></p>
        </div>
        <div class="metric">
            <h3>New Code Coverage</h3>
            <p><strong>%.2f%%</strong></p>
        </div>
        <div class="metric">
            <h3>Modified Files</h3>
            <p><strong>%d</strong></p>
        </div>
    </div>

    <h2>File Changes</h2>
    <table>
        <tr>
            <th>File</th>
            <th>Baseline Coverage</th>
            <th>Current Coverage</th>
            <th>Change</th>
            <th>Status</th>
        </tr>`,
		report.Timestamp.Format("2006-01-02 15:04:05"),
		getStatusClass(report.QualityGatesPassed),
		getStatusText(report.QualityGatesPassed),
		report.OverallChange,
		report.NewCodeCoverage,
		len(report.ModifiedFiles))

	// Add file changes
	for file, change := range report.FileChanges {
		statusClass := getStatusClass(change.MeetsThreshold)
		statusText := getStatusText(change.MeetsThreshold)

		htmlContent += fmt.Sprintf(`
        <tr>
            <td>%s</td>
            <td>%.1f%%%%</td>
            <td>%.1f%%%%</td>
            <td>%+.1f%%%%</td>
            <td class="%s">%s</td>
        </tr>`,
			file, change.BaselineCoverage, change.CurrentCoverage, change.Change, statusClass, statusText)
	}

	htmlContent += `
    </table>

    <h2>Recommendations</h2>
    <div class="recommendations">
        <ul>`

	for _, rec := range report.Recommendations {
		htmlContent += fmt.Sprintf(`<li>%s</li>`, rec)
	}

	htmlContent += `
        </ul>
    </div>

</body>
</html>`

	return os.WriteFile(outputPath, []byte(htmlContent), 0644)
}

func getStatusClass(passed bool) string {
	if passed {
		return "pass"
	}
	return "fail"
}

func getStatusText(passed bool) string {
	if passed {
		return "PASS"
	}
	return "FAIL"
}

// 🔺 COV-002: Save coverage history for trend tracking - 🔧
func saveCoverageHistory(baseline BaselineCoverage) error {
	historyPath := "docs/coverage-history.json"

	var histories []BaselineCoverage

	// Load existing history
	if data, err := os.ReadFile(historyPath); err == nil {
		json.Unmarshal(data, &histories)
	}

	// Add current baseline
	histories = append(histories, baseline)

	// Keep only last 50 entries
	if len(histories) > 50 {
		histories = histories[len(histories)-50:]
	}

	// Save back to file
	data, err := json.MarshalIndent(histories, "", "  ")
	if err != nil {
		return err
	}

	// Ensure directory exists
	os.MkdirAll(filepath.Dir(historyPath), 0755)

	return os.WriteFile(historyPath, data, 0644)
}

// 🔺 COV-002: Main function - 📝
func main() {
	fmt.Println("# Differential Coverage Report (COV-002: Baseline and Selective Reporting)")
	fmt.Println()

	// Load configuration
	var config DifferentialConfig
	if _, err := toml.DecodeFile("coverage.toml", &config); err != nil {
		log.Printf("Warning: Could not load coverage.toml: %v", err)
		// Use default configuration
		config.General.DifferentialMode = true
		config.QualityGates.NewCodeThreshold = 70.0
		config.QualityGates.CriticalPathThreshold = 80.0
		config.Differential.Enabled = true
		config.Reporting.GenerateJSON = true
		config.Reporting.GenerateDiffHTML = true
		config.Reporting.OutputDir = "coverage_reports"
	}

	if !config.General.DifferentialMode || !config.Differential.Enabled {
		fmt.Println("Differential mode is disabled. Run standard coverage instead.")
		return
	}

	// Get modified files
	modifiedFiles, err := getModifiedFiles()
	if err != nil {
		log.Printf("Warning: Could not detect modified files: %v", err)
		modifiedFiles = []string{} // Continue with empty list
	}

	fmt.Printf("Modified files detected: %d\n", len(modifiedFiles))
	for _, file := range modifiedFiles {
		fmt.Printf("  - %s\n", file)
	}
	fmt.Println()

	// Parse current coverage
	currentCoverage, err := parseCurrentCoverage("coverage.out")
	if err != nil {
		log.Fatalf("Failed to parse current coverage: %v", err)
	}

	// Load baseline
	baseline, err := loadBaselineCoverage()
	if err != nil {
		log.Printf("Warning: Could not load baseline coverage: %v", err)
		// Create minimal baseline
		baseline.OverallCoverage = 73.3 // From config
		baseline.Files = make(map[string]FileCoverage)
		baseline.Timestamp = time.Now()
	}

	// Generate differential report
	report := generateDifferentialReport(config, currentCoverage, baseline, modifiedFiles)

	// Print summary
	fmt.Printf("Quality Gates: %s\n", getStatusText(report.QualityGatesPassed))
	fmt.Printf("Overall Change: %+.2f%%\n", report.OverallChange)
	fmt.Printf("New Code Coverage: %.2f%%\n", report.NewCodeCoverage)
	fmt.Println()

	// Print recommendations
	fmt.Println("Recommendations:")
	for _, rec := range report.Recommendations {
		fmt.Printf("  • %s\n", rec)
	}
	fmt.Println()

	// Generate reports
	os.MkdirAll(config.Reporting.OutputDir, 0755)

	if config.Reporting.GenerateDiffJSON {
		jsonPath := filepath.Join(config.Reporting.OutputDir, "differential.json")
		if data, err := json.MarshalIndent(report, "", "  "); err == nil {
			os.WriteFile(jsonPath, data, 0644)
			fmt.Printf("Differential JSON report: %s\n", jsonPath)
		}
	}

	if config.Reporting.GenerateDiffHTML {
		htmlPath := filepath.Join(config.Reporting.OutputDir, "differential.html")
		if err := generateHTMLReport(report, htmlPath); err == nil {
			fmt.Printf("Differential HTML report: %s\n", htmlPath)
		}
	}

	// Update coverage history
	if config.General.TrendTracking {
		currentBaseline := BaselineCoverage{
			Timestamp:       time.Now(),
			OverallCoverage: calculateOverallCoverage(currentCoverage),
			Files:           make(map[string]FileCoverage),
		}

		for file, cov := range currentCoverage {
			currentBaseline.Files[file] = cov
		}

		if err := saveCoverageHistory(currentBaseline); err == nil {
			fmt.Printf("Coverage history updated: docs/coverage-history.json\n")
		}
	}

	fmt.Println("✓ Differential coverage analysis complete")
}
