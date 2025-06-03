// 🔶 DOC-010: Core data types - 🔧 Token suggestion data structures
package main

import (
	"time"
)

// 🔶 DOC-010: Function signature analysis - 🔍 Code structure representation
type FunctionSignature struct {
	Name       string      `json:"name"`
	ReturnType string      `json:"return_type"`
	Parameters []Parameter `json:"parameters"`
	Receiver   string      `json:"receiver,omitempty"`
	IsExported bool        `json:"is_exported"`
}

// 🔶 DOC-010: Function parameter analysis - 📝 Parameter representation
type Parameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// 🔶 DOC-010: Token suggestion structure - 💡 Core suggestion data
type TokenSuggestion struct {
	FilePath          string            `json:"file_path"`
	LineNumber        int               `json:"line_number"`
	FunctionName      string            `json:"function_name"`
	FeatureID         string            `json:"feature_id"`
	PriorityIcon      string            `json:"priority_icon"`
	PriorityReason    string            `json:"priority_reason"`
	ActionIcon        string            `json:"action_icon"`
	ActionReason      string            `json:"action_reason"`
	SuggestedToken    string            `json:"suggested_token"`
	Confidence        float64           `json:"confidence"`
	ComplexityLevel   string            `json:"complexity_level"`
	FunctionSignature FunctionSignature `json:"function_signature"`
	Context           map[string]string `json:"context"`
	Timestamp         time.Time         `json:"timestamp"`
}

// 🔶 DOC-010: Analysis results container - 📊 Comprehensive analysis output
type AnalysisResults struct {
	Target            string            `json:"target"`
	FunctionsAnalyzed int               `json:"functions_analyzed"`
	MissingTokens     int               `json:"missing_tokens"`
	FormatViolations  int               `json:"format_violations"`
	Suggestions       []TokenSuggestion `json:"suggestions"`
	ProcessingTime    time.Duration     `json:"processing_time"`
	AnalysisTimestamp time.Time         `json:"analysis_timestamp"`
	ConfidenceStats   ConfidenceStats   `json:"confidence_stats"`
}

// 🔶 DOC-010: Confidence statistics - 📈 Quality metrics
type ConfidenceStats struct {
	AverageConfidence float64 `json:"average_confidence"`
	MinConfidence     float64 `json:"min_confidence"`
	MaxConfidence     float64 `json:"max_confidence"`
	HighConfidence    int     `json:"high_confidence_count"`   // >0.8
	MediumConfidence  int     `json:"medium_confidence_count"` // 0.5-0.8
	LowConfidence     int     `json:"low_confidence_count"`    // <0.5
}

// 🔶 DOC-010: Token validation violation - 🛡️ Standards compliance tracking
type TokenViolation struct {
	FilePath      string `json:"file_path"`
	LineNumber    int    `json:"line_number"`
	ViolationType string `json:"violation_type"`
	CurrentToken  string `json:"current_token"`
	SuggestedFix  string `json:"suggested_fix"`
	Severity      string `json:"severity"`
	RuleID        string `json:"rule_id"`
	Description   string `json:"description"`
}

// 🔶 DOC-010: Batch processing results - 🚀 Mass analysis output
type BatchResults struct {
	Directory         string                `json:"directory"`
	FilesProcessed    int                   `json:"files_processed"`
	TotalFunctions    int                   `json:"total_functions"`
	TotalSuggestions  int                   `json:"total_suggestions"`
	TotalViolations   int                   `json:"total_violations"`
	ProcessingTime    time.Duration         `json:"processing_time"`
	PriorityBreakdown PriorityBreakdown     `json:"priority_breakdown"`
	ActionBreakdown   ActionBreakdown       `json:"action_breakdown"`
	TopSuggestions    []TokenSuggestion     `json:"top_suggestions"`
	FileResults       map[string]FileResult `json:"file_results"`
}

// 🔶 DOC-010: Priority categorization - ⭐ Priority distribution tracking
type PriorityBreakdown struct {
	Critical int `json:"critical"` // ⭐
	High     int `json:"high"`     // 🔺
	Medium   int `json:"medium"`   // 🔶
	Low      int `json:"low"`      // 🔻
}

// 🔶 DOC-010: Action categorization - 🔧 Action type distribution
type ActionBreakdown struct {
	Analysis      int `json:"analysis"`      // 🔍
	Documentation int `json:"documentation"` // 📝
	Configuration int `json:"configuration"` // 🔧
	Protection    int `json:"protection"`    // 🛡️
}

// 🔶 DOC-010: Individual file results - 📁 Per-file analysis data
type FileResult struct {
	FilePath         string            `json:"file_path"`
	FunctionsFound   int               `json:"functions_found"`
	SuggestionsCount int               `json:"suggestions_count"`
	ViolationsCount  int               `json:"violations_count"`
	Suggestions      []TokenSuggestion `json:"suggestions"`
	Violations       []TokenViolation  `json:"violations"`
	ProcessingTime   time.Duration     `json:"processing_time"`
}

// 🔶 DOC-010: Feature mapping entry - 🎯 Feature tracking integration
type FeatureMapping struct {
	FeatureID   string `json:"feature_id"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// 🔶 DOC-010: Analysis configuration - ⚙️ Suggestion engine settings
type AnalysisConfig struct {
	// 🔶 DOC-010: Priority assignment settings - ⭐ Priority determination
	PriorityRules struct {
		CriticalPatterns []string `json:"critical_patterns"`
		HighPatterns     []string `json:"high_patterns"`
		MediumPatterns   []string `json:"medium_patterns"`
		LowPatterns      []string `json:"low_patterns"`
	} `json:"priority_rules"`

	// 🔶 DOC-010: Action assignment settings - 🔧 Action determination
	ActionRules struct {
		AnalysisPatterns      []string `json:"analysis_patterns"`
		DocumentationPatterns []string `json:"documentation_patterns"`
		ConfigurationPatterns []string `json:"configuration_patterns"`
		ProtectionPatterns    []string `json:"protection_patterns"`
	} `json:"action_rules"`

	// 🔶 DOC-010: Feature tracking integration - 🎯 Feature ID mapping
	FeatureTrackingFile string                    `json:"feature_tracking_file"`
	FeatureMappings     map[string]FeatureMapping `json:"feature_mappings"`

	// 🔶 DOC-010: Confidence calculation settings - 📊 Quality thresholds
	ConfidenceWeights struct {
		SignatureMatch   float64 `json:"signature_match"`
		PatternMatch     float64 `json:"pattern_match"`
		ContextMatch     float64 `json:"context_match"`
		FeatureMapping   float64 `json:"feature_mapping"`
		ComplexityFactor float64 `json:"complexity_factor"`
	} `json:"confidence_weights"`

	// 🔶 DOC-010: Validation settings - 🛡️ Quality control
	ValidationRules struct {
		RequiredTokenFormat string   `json:"required_token_format"`
		ValidPriorityIcons  []string `json:"valid_priority_icons"`
		ValidActionIcons    []string `json:"valid_action_icons"`
		MinConfidence       float64  `json:"min_confidence"`
		MaxSuggestions      int      `json:"max_suggestions"`
	} `json:"validation_rules"`

	// 🔶 DOC-010: Output formatting settings - 📝 Display preferences
	OutputSettings struct {
		VerboseOutput    bool `json:"verbose_output"`
		IncludeContext   bool `json:"include_context"`
		SortByConfidence bool `json:"sort_by_confidence"`
		GroupByPriority  bool `json:"group_by_priority"`
		ShowStatistics   bool `json:"show_statistics"`
	} `json:"output_settings"`
}

// 🔶 DOC-010: Function complexity analysis - 📊 Complexity metrics
type ComplexityMetrics struct {
	CyclomaticComplexity int     `json:"cyclomatic_complexity"`
	LinesOfCode          int     `json:"lines_of_code"`
	ParameterCount       int     `json:"parameter_count"`
	ReturnValueCount     int     `json:"return_value_count"`
	CallDepth            int     `json:"call_depth"`
	ComplexityLevel      string  `json:"complexity_level"` // "LOW", "MEDIUM", "HIGH", "CRITICAL"
	ComplexityScore      float64 `json:"complexity_score"`
}

// 🔶 DOC-010: Pattern matching result - 🔍 Pattern recognition output
type PatternMatch struct {
	Pattern     string  `json:"pattern"`
	MatchType   string  `json:"match_type"` // "function_name", "parameter", "return_type", "context"
	Confidence  float64 `json:"confidence"`
	Description string  `json:"description"`
}

// 🔶 DOC-010: Context analysis result - 🔍 Surrounding code analysis
type ContextAnalysis struct {
	SurroundingLines   []string       `json:"surrounding_lines"`
	ImportStatements   []string       `json:"import_statements"`
	UsedTypes          []string       `json:"used_types"`
	CalledFunctions    []string       `json:"called_functions"`
	ErrorHandling      bool           `json:"error_handling"`
	ResourceManagement bool           `json:"resource_management"`
	PatternMatches     []PatternMatch `json:"pattern_matches"`
}

// 🔶 DOC-010: Suggestion generation metadata - 📊 Generation tracking
type SuggestionMetadata struct {
	GenerationRules        []string           `json:"generation_rules"`
	AppliedPatterns        []string           `json:"applied_patterns"`
	ConfidenceFactors      map[string]float64 `json:"confidence_factors"`
	AlternativeSuggestions []string           `json:"alternative_suggestions"`
	ValidationResults      []string           `json:"validation_results"`
	GenerationTime         time.Duration      `json:"generation_time"`
}

// 🔶 DOC-010: Default analysis configuration - ⚙️ Standard configuration
func DefaultAnalysisConfig() *AnalysisConfig {
	return &AnalysisConfig{
		PriorityRules: struct {
			CriticalPatterns []string `json:"critical_patterns"`
			HighPatterns     []string `json:"high_patterns"`
			MediumPatterns   []string `json:"medium_patterns"`
			LowPatterns      []string `json:"low_patterns"`
		}{
			CriticalPatterns: []string{"main", "init", "archive", "backup", "create"},
			HighPatterns:     []string{"config", "load", "save", "validate", "verify", "generate"},
			MediumPatterns:   []string{"format", "parse", "convert", "transform", "helper"},
			LowPatterns:      []string{"test", "mock", "example", "util", "debug"},
		},
		ActionRules: struct {
			AnalysisPatterns      []string `json:"analysis_patterns"`
			DocumentationPatterns []string `json:"documentation_patterns"`
			ConfigurationPatterns []string `json:"configuration_patterns"`
			ProtectionPatterns    []string `json:"protection_patterns"`
		}{
			AnalysisPatterns:      []string{"get", "find", "search", "discover", "detect", "analyze", "check", "parse"},
			DocumentationPatterns: []string{"format", "print", "write", "update", "log", "output", "render", "display"},
			ConfigurationPatterns: []string{"set", "config", "init", "setup", "create", "build", "generate", "make"},
			ProtectionPatterns:    []string{"protect", "secure", "validate", "verify", "guard", "ensure", "handle", "recover"},
		},
		FeatureTrackingFile: "docs/context/feature-tracking.md",
		ConfidenceWeights: struct {
			SignatureMatch   float64 `json:"signature_match"`
			PatternMatch     float64 `json:"pattern_match"`
			ContextMatch     float64 `json:"context_match"`
			FeatureMapping   float64 `json:"feature_mapping"`
			ComplexityFactor float64 `json:"complexity_factor"`
		}{
			SignatureMatch:   0.3,
			PatternMatch:     0.25,
			ContextMatch:     0.2,
			FeatureMapping:   0.15,
			ComplexityFactor: 0.1,
		},
		ValidationRules: struct {
			RequiredTokenFormat string   `json:"required_token_format"`
			ValidPriorityIcons  []string `json:"valid_priority_icons"`
			ValidActionIcons    []string `json:"valid_action_icons"`
			MinConfidence       float64  `json:"min_confidence"`
			MaxSuggestions      int      `json:"max_suggestions"`
		}{
			RequiredTokenFormat: `^// [⭐🔺🔶🔻] [A-Z]+-[0-9]+: .+ - [🔍📝🔧🛡️] .+$`,
			ValidPriorityIcons:  []string{"⭐", "🔺", "🔶", "🔻"},
			ValidActionIcons:    []string{"🔍", "📝", "🔧", "🛡️"},
			MinConfidence:       0.5,
			MaxSuggestions:      100,
		},
		OutputSettings: struct {
			VerboseOutput    bool `json:"verbose_output"`
			IncludeContext   bool `json:"include_context"`
			SortByConfidence bool `json:"sort_by_confidence"`
			GroupByPriority  bool `json:"group_by_priority"`
			ShowStatistics   bool `json:"show_statistics"`
		}{
			VerboseOutput:    false,
			IncludeContext:   true,
			SortByConfidence: true,
			GroupByPriority:  false,
			ShowStatistics:   true,
		},
	}
}
