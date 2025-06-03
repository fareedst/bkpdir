// 🔶 DOC-010: Token analysis engine - 🔍 Core analysis and suggestion logic
package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 🔶 DOC-010: Core token analyzer - 🔍 Main analysis engine
type TokenAnalyzer struct {
	config     *AnalysisConfig
	fileSet    *token.FileSet
	featureMap map[string]FeatureMapping
}

// 🔶 DOC-010: Token validator - 🛡️ Standards compliance validator
type TokenValidator struct {
	config *AnalysisConfig
}

// 🔶 DOC-010: Batch processor - 🚀 Mass processing engine
type BatchProcessor struct {
	analyzer  *TokenAnalyzer
	validator *TokenValidator
}

// 🔶 DOC-010: Analyzer constructor - 🔧 Initialize analysis engine
func NewTokenAnalyzer() *TokenAnalyzer {
	config := DefaultAnalysisConfig()
	analyzer := &TokenAnalyzer{
		config:  config,
		fileSet: token.NewFileSet(),
	}

	// 🔶 DOC-010: Feature tracking integration - 🎯 Load feature mappings
	if err := analyzer.loadFeatureMappings(); err != nil {
		// 🔶 DOC-010: Fallback to empty mapping - 🛡️ Graceful degradation
		analyzer.featureMap = make(map[string]FeatureMapping)
	}

	return analyzer
}

// 🔶 DOC-010: Validator constructor - 🛡️ Initialize validation engine
func NewTokenValidator() *TokenValidator {
	return &TokenValidator{
		config: DefaultAnalysisConfig(),
	}
}

// 🔶 DOC-010: Batch processor constructor - 🚀 Initialize batch processing
func NewBatchProcessor() *BatchProcessor {
	return &BatchProcessor{
		analyzer:  NewTokenAnalyzer(),
		validator: NewTokenValidator(),
	}
}

// 🔶 DOC-010: Feature mapping loader - 🎯 Feature tracking integration
func (ta *TokenAnalyzer) loadFeatureMappings() error {
	ta.featureMap = make(map[string]FeatureMapping)

	// 🔶 DOC-010: Feature tracking file processing - 📝 Parse feature tracking
	file, err := os.Open(ta.config.FeatureTrackingFile)
	if err != nil {
		return fmt.Errorf("failed to open feature tracking file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	featureRegex := regexp.MustCompile(`\|\s*([A-Z]+-[0-9]+)\s*\|.*\|\s*(⭐|🔺|🔶|🔻)\s*(CRITICAL|HIGH|MEDIUM|LOW)`)

	// 🔶 DOC-010: Feature extraction from table - 🔍 Parse feature entries
	for scanner.Scan() {
		line := scanner.Text()
		matches := featureRegex.FindStringSubmatch(line)
		if len(matches) >= 4 {
			featureID := matches[1]
			priorityIcon := matches[2]
			priorityText := matches[3]

			ta.featureMap[featureID] = FeatureMapping{
				FeatureID:   featureID,
				Priority:    priorityText,
				Status:      "UNKNOWN", // Could be enhanced to parse status
				Description: "",        // Could be enhanced to parse description
				Category:    strings.Split(featureID, "-")[0],
			}

			// 🔶 DOC-010: Priority icon validation - 🛡️ Consistency checking
			if priorityIcon != "" {
				// Store the icon for consistency checking
				if ta.featureMap[featureID].Priority == "" {
					mapping := ta.featureMap[featureID]
					mapping.Priority = priorityText
					ta.featureMap[featureID] = mapping
				}
			}
		}
	}

	return scanner.Err()
}

// 🔶 DOC-010: Target analysis - 🔍 Analyze directory or file
func (ta *TokenAnalyzer) AnalyzeTarget(target string) (*AnalysisResults, error) {
	startTime := time.Now()

	// 🔶 DOC-010: Target validation - 🛡️ Input validation
	info, err := os.Stat(target)
	if err != nil {
		return nil, fmt.Errorf("target not found: %w", err)
	}

	results := &AnalysisResults{
		Target:            target,
		FunctionsAnalyzed: 0,
		MissingTokens:     0,
		FormatViolations:  0,
		Suggestions:       make([]TokenSuggestion, 0),
		AnalysisTimestamp: time.Now(),
	}

	// 🔶 DOC-010: Directory vs file processing - 🔧 Path handling
	if info.IsDir() {
		err = ta.analyzeDirectory(target, results)
	} else {
		err = ta.analyzeFile(target, results)
	}

	if err != nil {
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	// 🔶 DOC-010: Results post-processing - 📊 Statistics calculation
	ta.calculateStatistics(results)
	results.ProcessingTime = time.Since(startTime)

	return results, nil
}

// 🔶 DOC-010: Directory analysis - 🔍 Recursive directory processing
func (ta *TokenAnalyzer) analyzeDirectory(dir string, results *AnalysisResults) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 🔶 DOC-010: Go file filtering - 📁 File type validation
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			// Skip test files unless explicitly requested
			if strings.HasSuffix(path, "_test.go") {
				return nil
			}

			return ta.analyzeFile(path, results)
		}

		return nil
	})
}

// 🔶 DOC-010: File analysis - 🔍 Individual file processing
func (ta *TokenAnalyzer) analyzeFile(filePath string, results *AnalysisResults) error {
	// 🔶 DOC-010: AST parsing - 🔍 Code structure analysis
	src, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	file, err := parser.ParseFile(ta.fileSet, filePath, src, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	// 🔶 DOC-010: Function discovery - 🔍 Function extraction
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:
			if node.Name != nil && node.Name.IsExported() {
				results.FunctionsAnalyzed++

				// 🔶 DOC-010: Function analysis - 🔍 Individual function processing
				suggestion, err := ta.analyzeFunctionDecl(filePath, node, src)
				if err == nil && suggestion != nil {
					results.Suggestions = append(results.Suggestions, *suggestion)
				}
			}
		}
		return true
	})

	return nil
}

// 🔶 DOC-010: Function declaration analysis - 🔍 Function-specific analysis
func (ta *TokenAnalyzer) analyzeFunctionDecl(filePath string, funcDecl *ast.FuncDecl, src []byte) (*TokenSuggestion, error) {
	// 🔶 DOC-010: Function signature extraction - 📝 Signature analysis
	signature := ta.extractFunctionSignature(funcDecl)

	// 🔶 DOC-010: Line number calculation - 📍 Position tracking
	position := ta.fileSet.Position(funcDecl.Pos())

	// 🔶 DOC-010: Context analysis - 🔍 Surrounding code analysis
	context, err := ta.analyzeContext(filePath, position.Line, src)
	if err != nil {
		return nil, fmt.Errorf("context analysis failed: %w", err)
	}

	// 🔶 DOC-010: Token suggestion generation - 💡 Suggestion creation
	suggestion := &TokenSuggestion{
		FilePath:          filePath,
		LineNumber:        position.Line,
		FunctionName:      signature.Name,
		FunctionSignature: signature,
		Context:           context,
		Timestamp:         time.Now(),
	}

	// 🔶 DOC-010: Priority assignment - ⭐ Priority determination
	suggestion.PriorityIcon, suggestion.PriorityReason = ta.determinePriority(signature, context)

	// 🔶 DOC-010: Action assignment - 🔧 Action determination
	suggestion.ActionIcon, suggestion.ActionReason = ta.determineAction(signature, context)

	// 🔶 DOC-010: Feature mapping - 🎯 Feature ID assignment
	suggestion.FeatureID = ta.determineFeatureID(signature, context)

	// 🔶 DOC-010: Confidence calculation - 📊 Quality scoring
	suggestion.Confidence = ta.calculateConfidence(signature, context, suggestion)

	// 🔶 DOC-010: Complexity analysis - 📊 Complexity determination
	suggestion.ComplexityLevel = ta.analyzeComplexity(funcDecl)

	// 🔶 DOC-010: Token format generation - 📝 Final token creation
	suggestion.SuggestedToken = fmt.Sprintf("// %s %s: %s - %s %s",
		suggestion.PriorityIcon,
		suggestion.FeatureID,
		ta.generateDescription(signature, context),
		suggestion.ActionIcon,
		suggestion.ActionReason)

	return suggestion, nil
}

// 🔶 DOC-010: Function signature extraction - 📝 Signature parsing
func (ta *TokenAnalyzer) extractFunctionSignature(funcDecl *ast.FuncDecl) FunctionSignature {
	signature := FunctionSignature{
		Name:       funcDecl.Name.Name,
		IsExported: funcDecl.Name.IsExported(),
		Parameters: make([]Parameter, 0),
	}

	// 🔶 DOC-010: Receiver analysis - 🔍 Method receiver extraction
	if funcDecl.Recv != nil {
		for _, field := range funcDecl.Recv.List {
			if len(field.Names) > 0 {
				signature.Receiver = field.Names[0].Name
			}
		}
	}

	// 🔶 DOC-010: Parameter extraction - 📝 Parameter analysis
	if funcDecl.Type.Params != nil {
		for _, field := range funcDecl.Type.Params.List {
			paramType := "unknown"
			if field.Type != nil {
				paramType = fmt.Sprintf("%T", field.Type)
			}

			if len(field.Names) > 0 {
				for _, name := range field.Names {
					signature.Parameters = append(signature.Parameters, Parameter{
						Name: name.Name,
						Type: paramType,
					})
				}
			} else {
				signature.Parameters = append(signature.Parameters, Parameter{
					Name: "",
					Type: paramType,
				})
			}
		}
	}

	// 🔶 DOC-010: Return type extraction - 📝 Return type analysis
	if funcDecl.Type.Results != nil {
		returnTypes := make([]string, 0)
		for _, field := range funcDecl.Type.Results.List {
			if field.Type != nil {
				returnTypes = append(returnTypes, fmt.Sprintf("%T", field.Type))
			}
		}
		signature.ReturnType = strings.Join(returnTypes, ", ")
	} else {
		signature.ReturnType = "void"
	}

	return signature
}

// 🔶 DOC-010: Context analysis - 🔍 Surrounding code analysis
func (ta *TokenAnalyzer) analyzeContext(filePath string, lineNumber int, src []byte) (map[string]string, error) {
	context := make(map[string]string)

	// 🔶 DOC-010: Source code line extraction - 📝 Line-based analysis
	lines := strings.Split(string(src), "\n")
	if lineNumber > len(lines) {
		return context, fmt.Errorf("line number %d exceeds file length", lineNumber)
	}

	// 🔶 DOC-010: Context window extraction - 🔍 Surrounding code collection
	start := max(0, lineNumber-5)
	end := min(len(lines), lineNumber+5)

	contextLines := lines[start:end]
	context["surrounding_code"] = strings.Join(contextLines, "\n")

	// 🔶 DOC-010: Pattern matching - 🔍 Code pattern recognition
	fullCode := strings.Join(lines, "\n")

	// 🔶 DOC-010: Error handling detection - 🛡️ Error pattern analysis
	if strings.Contains(fullCode, "error") || strings.Contains(fullCode, "err") {
		context["error_handling"] = "true"
	}

	// 🔶 DOC-010: Resource management detection - 🛡️ Resource pattern analysis
	if strings.Contains(fullCode, "defer") || strings.Contains(fullCode, "Close") {
		context["resource_management"] = "true"
	}

	// 🔶 DOC-010: Configuration detection - 🔧 Config pattern analysis
	if strings.Contains(fullCode, "config") || strings.Contains(fullCode, "Config") {
		context["configuration_related"] = "true"
	}

	// 🔶 DOC-010: File operation detection - 🔍 File pattern analysis
	if strings.Contains(fullCode, "file") || strings.Contains(fullCode, "File") {
		context["file_operations"] = "true"
	}

	return context, nil
}

// 🔶 DOC-010: Priority determination - ⭐ Priority assignment logic
func (ta *TokenAnalyzer) determinePriority(signature FunctionSignature, context map[string]string) (string, string) {
	funcName := strings.ToLower(signature.Name)

	// 🔶 DOC-010: Critical priority patterns - ⭐ Critical function detection
	for _, pattern := range ta.config.PriorityRules.CriticalPatterns {
		if strings.Contains(funcName, strings.ToLower(pattern)) {
			return "⭐", fmt.Sprintf("Critical operation: %s", pattern)
		}
	}

	// 🔶 DOC-010: High priority patterns - 🔺 High priority detection
	for _, pattern := range ta.config.PriorityRules.HighPatterns {
		if strings.Contains(funcName, strings.ToLower(pattern)) {
			return "🔺", fmt.Sprintf("High priority: %s", pattern)
		}
	}

	// 🔶 DOC-010: Medium priority patterns - 🔶 Medium priority detection
	for _, pattern := range ta.config.PriorityRules.MediumPatterns {
		if strings.Contains(funcName, strings.ToLower(pattern)) {
			return "🔶", fmt.Sprintf("Medium priority: %s", pattern)
		}
	}

	// 🔶 DOC-010: Context-based priority adjustment - 🔍 Context priority analysis
	if context["error_handling"] == "true" {
		return "🔺", "High priority: error handling function"
	}

	if context["resource_management"] == "true" {
		return "🔺", "High priority: resource management function"
	}

	// 🔶 DOC-010: Default low priority - 🔻 Default assignment
	return "🔻", "Low priority: utility function"
}

// 🔶 DOC-010: Action determination - 🔧 Action assignment logic
func (ta *TokenAnalyzer) determineAction(signature FunctionSignature, context map[string]string) (string, string) {
	funcName := strings.ToLower(signature.Name)

	// 🔶 DOC-010: Analysis action patterns - 🔍 Analysis function detection
	for _, pattern := range ta.config.ActionRules.AnalysisPatterns {
		if strings.Contains(funcName, strings.ToLower(pattern)) {
			return "🔍", fmt.Sprintf("Analysis operation: %s", pattern)
		}
	}

	// 🔶 DOC-010: Documentation action patterns - 📝 Documentation function detection
	for _, pattern := range ta.config.ActionRules.DocumentationPatterns {
		if strings.Contains(funcName, strings.ToLower(pattern)) {
			return "📝", fmt.Sprintf("Documentation operation: %s", pattern)
		}
	}

	// 🔶 DOC-010: Configuration action patterns - 🔧 Configuration function detection
	for _, pattern := range ta.config.ActionRules.ConfigurationPatterns {
		if strings.Contains(funcName, strings.ToLower(pattern)) {
			return "🔧", fmt.Sprintf("Configuration operation: %s", pattern)
		}
	}

	// 🔶 DOC-010: Protection action patterns - 🛡️ Protection function detection
	for _, pattern := range ta.config.ActionRules.ProtectionPatterns {
		if strings.Contains(funcName, strings.ToLower(pattern)) {
			return "🛡️", fmt.Sprintf("Protection operation: %s", pattern)
		}
	}

	// 🔶 DOC-010: Context-based action assignment - 🔍 Context action analysis
	if context["error_handling"] == "true" {
		return "🛡️", "Protection: error handling"
	}

	if context["configuration_related"] == "true" {
		return "🔧", "Configuration: config management"
	}

	// 🔶 DOC-010: Default configuration action - 🔧 Default assignment
	return "🔧", "Configuration: general operation"
}

// 🔶 DOC-010: Feature ID determination - 🎯 Feature mapping logic
func (ta *TokenAnalyzer) determineFeatureID(signature FunctionSignature, context map[string]string) string {
	funcName := strings.ToLower(signature.Name)

	// 🔶 DOC-010: Direct feature mapping - 🎯 Feature tracking lookup
	for featureID, mapping := range ta.featureMap {
		category := strings.ToLower(mapping.Category)
		if strings.Contains(funcName, category) {
			return featureID
		}
	}

	// 🔶 DOC-010: Category-based assignment - 🎯 Category inference
	if strings.Contains(funcName, "config") {
		return "CFG-NEW"
	}
	if strings.Contains(funcName, "archive") {
		return "ARCH-NEW"
	}
	if strings.Contains(funcName, "backup") {
		return "FILE-NEW"
	}
	if strings.Contains(funcName, "git") {
		return "GIT-NEW"
	}
	if strings.Contains(funcName, "test") {
		return "TEST-NEW"
	}

	// 🔶 DOC-010: Default generic assignment - 🎯 Fallback assignment
	return "UTIL-NEW"
}

// 🔶 DOC-010: Confidence calculation - 📊 Quality scoring algorithm
func (ta *TokenAnalyzer) calculateConfidence(signature FunctionSignature, context map[string]string, suggestion *TokenSuggestion) float64 {
	confidence := 0.0
	weights := ta.config.ConfidenceWeights

	// 🔶 DOC-010: Signature matching confidence - 📝 Signature quality assessment
	if signature.IsExported && len(signature.Parameters) > 0 {
		confidence += weights.SignatureMatch
	}

	// 🔶 DOC-010: Pattern matching confidence - 🔍 Pattern recognition quality
	if suggestion.PriorityReason != "Low priority: utility function" {
		confidence += weights.PatternMatch
	}

	// 🔶 DOC-010: Context matching confidence - 🔍 Context analysis quality
	if len(context) > 2 {
		confidence += weights.ContextMatch
	}

	// 🔶 DOC-010: Feature mapping confidence - 🎯 Feature tracking quality
	if !strings.HasSuffix(suggestion.FeatureID, "-NEW") {
		confidence += weights.FeatureMapping
	}

	// 🔶 DOC-010: Complexity factor - 📊 Complexity-based adjustment
	if len(signature.Parameters) > 3 || signature.ReturnType != "void" {
		confidence += weights.ComplexityFactor
	}

	// 🔶 DOC-010: Confidence normalization - 📊 Score normalization
	return minFloat64(1.0, confidence)
}

// 🔶 DOC-010: Complexity analysis - 📊 Function complexity assessment
func (ta *TokenAnalyzer) analyzeComplexity(funcDecl *ast.FuncDecl) string {
	paramCount := 0
	if funcDecl.Type.Params != nil {
		paramCount = len(funcDecl.Type.Params.List)
	}

	returnCount := 0
	if funcDecl.Type.Results != nil {
		returnCount = len(funcDecl.Type.Results.List)
	}

	// 🔶 DOC-010: Complexity scoring - 📊 Simple complexity heuristics
	complexityScore := paramCount + returnCount

	if complexityScore >= 8 {
		return "CRITICAL"
	} else if complexityScore >= 5 {
		return "HIGH"
	} else if complexityScore >= 3 {
		return "MEDIUM"
	}

	return "LOW"
}

// 🔶 DOC-010: Description generation - 📝 Human-readable description
func (ta *TokenAnalyzer) generateDescription(signature FunctionSignature, context map[string]string) string {
	funcName := signature.Name

	// 🔶 DOC-010: Context-based description - 📝 Contextual description generation
	if context["error_handling"] == "true" {
		return fmt.Sprintf("%s with error handling", funcName)
	}

	if context["resource_management"] == "true" {
		return fmt.Sprintf("%s with resource management", funcName)
	}

	if context["configuration_related"] == "true" {
		return fmt.Sprintf("%s configuration operation", funcName)
	}

	// 🔶 DOC-010: Default description - 📝 Generic description
	return fmt.Sprintf("%s implementation", funcName)
}

// 🔶 DOC-010: Statistics calculation - 📊 Results post-processing
func (ta *TokenAnalyzer) calculateStatistics(results *AnalysisResults) {
	if len(results.Suggestions) == 0 {
		return
	}

	// 🔶 DOC-010: Confidence statistics - 📈 Confidence analysis
	var totalConfidence float64
	minConf := 1.0
	maxConf := 0.0
	highCount := 0
	mediumCount := 0
	lowCount := 0

	for _, suggestion := range results.Suggestions {
		totalConfidence += suggestion.Confidence
		minConf = minFloat64(minConf, suggestion.Confidence)
		maxConf = maxFloat64(maxConf, suggestion.Confidence)

		if suggestion.Confidence > 0.8 {
			highCount++
		} else if suggestion.Confidence > 0.5 {
			mediumCount++
		} else {
			lowCount++
		}
	}

	results.ConfidenceStats = ConfidenceStats{
		AverageConfidence: totalConfidence / float64(len(results.Suggestions)),
		MinConfidence:     minConf,
		MaxConfidence:     maxConf,
		HighConfidence:    highCount,
		MediumConfidence:  mediumCount,
		LowConfidence:     lowCount,
	}

	// 🔶 DOC-010: Suggestion sorting - 📊 Quality-based ordering
	sort.Slice(results.Suggestions, func(i, j int) bool {
		return results.Suggestions[i].Confidence > results.Suggestions[j].Confidence
	})
}

// 🔶 DOC-010: Function-specific suggestion - 🔍 Individual function analysis
func (ta *TokenAnalyzer) SuggestForFunction(filePath, lineStr string) (*TokenSuggestion, error) {
	lineNumber, err := strconv.Atoi(lineStr)
	if err != nil {
		return nil, fmt.Errorf("invalid line number: %w", err)
	}

	// 🔶 DOC-010: File parsing for specific function - 🔍 Targeted analysis
	src, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	file, err := parser.ParseFile(ta.fileSet, filePath, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	// 🔶 DOC-010: Function location matching - 📍 Position-based lookup
	var targetFunc *ast.FuncDecl
	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			pos := ta.fileSet.Position(funcDecl.Pos())
			if pos.Line <= lineNumber && lineNumber <= ta.fileSet.Position(funcDecl.End()).Line {
				targetFunc = funcDecl
				return false
			}
		}
		return true
	})

	if targetFunc == nil {
		return nil, fmt.Errorf("no function found at line %d", lineNumber)
	}

	// 🔶 DOC-010: Targeted function analysis - 🔍 Single function processing
	return ta.analyzeFunctionDecl(filePath, targetFunc, src)
}

// 🔶 DOC-010: Token validation - 🛡️ Standards compliance checking
func (tv *TokenValidator) ValidateTokens(directory string) ([]TokenViolation, error) {
	violations := make([]TokenViolation, 0)

	// 🔶 DOC-010: Directory traversal for validation - 🛡️ Comprehensive validation
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			fileViolations, err := tv.validateFile(path)
			if err != nil {
				return err
			}
			violations = append(violations, fileViolations...)
		}

		return nil
	})

	return violations, err
}

// 🔶 DOC-010: File validation - 🛡️ Individual file compliance checking
func (tv *TokenValidator) validateFile(filePath string) ([]TokenViolation, error) {
	violations := make([]TokenViolation, 0)

	// 🔶 DOC-010: File content analysis - 📝 Token extraction and validation
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	tokenRegex := regexp.MustCompile(`//\s*([⭐🔺🔶🔻])?\s*([A-Z]+-[0-9]+):?\s*(.*)`)

	// 🔶 DOC-010: Line-by-line token validation - 🛡️ Format compliance checking
	for i, line := range lines {
		if strings.Contains(line, "//") && (strings.Contains(line, "-") || strings.Contains(line, ":")) {
			matches := tokenRegex.FindStringSubmatch(line)
			if len(matches) < 4 {
				violations = append(violations, TokenViolation{
					FilePath:      filePath,
					LineNumber:    i + 1,
					ViolationType: "INVALID_FORMAT",
					CurrentToken:  strings.TrimSpace(line),
					SuggestedFix:  "// 🔶 FEATURE-ID: Description - 🔧 Action description",
					Severity:      "WARNING",
					RuleID:        "FORMAT-001",
					Description:   "Token does not match required format",
				})
			}
		}
	}

	return violations, nil
}

// 🔶 DOC-010: Batch processing - 🚀 Directory-wide analysis
func (bp *BatchProcessor) ProcessDirectory(directory string) (*BatchResults, error) {
	startTime := time.Now()

	results := &BatchResults{
		Directory:         directory,
		FilesProcessed:    0,
		TotalFunctions:    0,
		FileResults:       make(map[string]FileResult),
		PriorityBreakdown: PriorityBreakdown{},
		ActionBreakdown:   ActionBreakdown{},
		TopSuggestions:    make([]TokenSuggestion, 0),
	}

	// 🔶 DOC-010: Directory traversal for batch processing - 🚀 Comprehensive processing
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			fileResult, err := bp.processFile(path)
			if err != nil {
				return err
			}

			results.FilesProcessed++
			results.TotalFunctions += fileResult.FunctionsFound
			results.TotalSuggestions += fileResult.SuggestionsCount
			results.TotalViolations += fileResult.ViolationsCount
			results.FileResults[path] = *fileResult

			// 🔶 DOC-010: Priority and action breakdown - 📊 Category statistics
			for _, suggestion := range fileResult.Suggestions {
				bp.updateBreakdowns(suggestion, results)
				results.TopSuggestions = append(results.TopSuggestions, suggestion)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 🔶 DOC-010: Results finalization - 📊 Final processing
	bp.finalizeResults(results)
	results.ProcessingTime = time.Since(startTime)

	return results, nil
}

// 🔶 DOC-010: Individual file processing for batch - 🔍 File-specific batch processing
func (bp *BatchProcessor) processFile(filePath string) (*FileResult, error) {
	startTime := time.Now()

	// 🔶 DOC-010: Combined analysis and validation - 🔍 Comprehensive file processing
	analysisResults, err := bp.analyzer.AnalyzeTarget(filePath)
	if err != nil {
		return nil, err
	}

	violations, err := bp.validator.validateFile(filePath)
	if err != nil {
		return nil, err
	}

	return &FileResult{
		FilePath:         filePath,
		FunctionsFound:   analysisResults.FunctionsAnalyzed,
		SuggestionsCount: len(analysisResults.Suggestions),
		ViolationsCount:  len(violations),
		Suggestions:      analysisResults.Suggestions,
		Violations:       violations,
		ProcessingTime:   time.Since(startTime),
	}, nil
}

// 🔶 DOC-010: Breakdown updates - 📊 Category statistics maintenance
func (bp *BatchProcessor) updateBreakdowns(suggestion TokenSuggestion, results *BatchResults) {
	// 🔶 DOC-010: Priority breakdown updates - ⭐ Priority statistics
	switch suggestion.PriorityIcon {
	case "⭐":
		results.PriorityBreakdown.Critical++
	case "🔺":
		results.PriorityBreakdown.High++
	case "🔶":
		results.PriorityBreakdown.Medium++
	case "🔻":
		results.PriorityBreakdown.Low++
	}

	// 🔶 DOC-010: Action breakdown updates - 🔧 Action statistics
	switch suggestion.ActionIcon {
	case "🔍":
		results.ActionBreakdown.Analysis++
	case "📝":
		results.ActionBreakdown.Documentation++
	case "🔧":
		results.ActionBreakdown.Configuration++
	case "🛡️":
		results.ActionBreakdown.Protection++
	}
}

// 🔶 DOC-010: Results finalization - 📊 Final results processing
func (bp *BatchProcessor) finalizeResults(results *BatchResults) {
	// 🔶 DOC-010: Top suggestions sorting - 📊 Quality-based ranking
	sort.Slice(results.TopSuggestions, func(i, j int) bool {
		return results.TopSuggestions[i].Confidence > results.TopSuggestions[j].Confidence
	})

	// 🔶 DOC-010: Top suggestions limiting - 📊 Results optimization
	if len(results.TopSuggestions) > 10 {
		results.TopSuggestions = results.TopSuggestions[:10]
	}
}

// 🔶 DOC-010: Utility functions - 🔧 Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
