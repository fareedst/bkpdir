// 🔺 DOC-011: Compliance monitoring system - 📊 AI assistant behavior tracking
package validation

import (
	"fmt"
	"sync"
	"time"
)

// ComplianceMonitor tracks AI assistant validation behavior and adherence
type ComplianceMonitor struct {
	behaviorTracker     *AIBehaviorTracker
	adherenceCalculator *AdherenceMetrics
	auditTrailManager   *AuditTrailManager
	dashboardGenerator  *ComplianceDashboard
	mu                  sync.RWMutex
}

// ValidationEvent represents a validation event for tracking
type ValidationEvent struct {
	AssistantID     string        `json:"assistant_id"`
	SessionID       string        `json:"session_id"`
	Status          string        `json:"status"` // "pass", "fail", "warning"
	ErrorCount      int           `json:"error_count"`
	WarningCount    int           `json:"warning_count"`
	ComplianceScore float64       `json:"compliance_score"`
	Timestamp       time.Time     `json:"timestamp"`
	ValidationMode  string        `json:"validation_mode"`
	ProcessingTime  time.Duration `json:"processing_time"`
}

// AIBehaviorTracker tracks patterns in AI assistant validation behavior
type AIBehaviorTracker struct {
	validationEvents  chan ValidationEvent
	behaviorPatterns  *PatternAnalyzer
	complianceScoring *ComplianceScorer
	eventHistory      []ValidationEvent
	mu                sync.RWMutex
}

// PatternAnalyzer analyzes AI assistant behavior patterns
type PatternAnalyzer struct {
	patterns map[string]*BehaviorPattern
	mu       sync.RWMutex
}

// BehaviorPattern represents a detected behavior pattern
type BehaviorPattern struct {
	AssistantID    string    `json:"assistant_id"`
	PatternType    string    `json:"pattern_type"`
	Frequency      float64   `json:"frequency"`
	LastOccurrence time.Time `json:"last_occurrence"`
	TrendDirection string    `json:"trend_direction"` // "improving", "declining", "stable"
	Confidence     float64   `json:"confidence"`
}

// ComplianceScorer calculates compliance scores
type ComplianceScorer struct {
	weights map[string]float64
	mu      sync.RWMutex
}

// AdherenceMetrics calculates adherence metrics
type AdherenceMetrics struct {
	assistantMetrics map[string]*AssistantMetrics
	mu               sync.RWMutex
}

// AssistantMetrics holds metrics for a specific AI assistant
type AssistantMetrics struct {
	AssistantID      string    `json:"assistant_id"`
	TotalValidations int       `json:"total_validations"`
	SuccessfulPasses int       `json:"successful_passes"`
	FailureCount     int       `json:"failure_count"`
	WarningCount     int       `json:"warning_count"`
	AverageScore     float64   `json:"average_score"`
	LastValidation   time.Time `json:"last_validation"`
	ComplianceRating string    `json:"compliance_rating"`
}

// AuditTrailManager manages audit trails for compliance tracking
type AuditTrailManager struct {
	auditEntries []AuditEntry
	mu           sync.RWMutex
}

// AuditEntry represents an audit trail entry
type AuditEntry struct {
	EntryID     string                 `json:"entry_id"`
	AssistantID string                 `json:"assistant_id"`
	Action      string                 `json:"action"`
	Details     map[string]interface{} `json:"details"`
	Timestamp   time.Time              `json:"timestamp"`
	Severity    string                 `json:"severity"`
}

// ComplianceDashboard generates compliance dashboards
type ComplianceDashboard struct {
	Metrics     map[string]*AssistantMetrics `json:"metrics"`
	Patterns    map[string]*BehaviorPattern  `json:"patterns"`
	Summary     *ComplianceSummary           `json:"summary"`
	GeneratedAt time.Time                    `json:"generated_at"`
}

// ComplianceSummary provides overall compliance summary
type ComplianceSummary struct {
	TotalAssistants    int      `json:"total_assistants"`
	OverallCompliance  float64  `json:"overall_compliance"`
	TrendDirection     string   `json:"trend_direction"`
	TopIssues          []string `json:"top_issues"`
	RecommendedActions []string `json:"recommended_actions"`
}

// NewComplianceMonitor creates a new compliance monitor
func NewComplianceMonitor() *ComplianceMonitor {
	return &ComplianceMonitor{
		behaviorTracker:     NewAIBehaviorTracker(),
		adherenceCalculator: NewAdherenceMetrics(),
		auditTrailManager:   NewAuditTrailManager(),
		dashboardGenerator:  NewComplianceDashboard(),
	}
}

// NewAIBehaviorTracker creates a new AI behavior tracker
func NewAIBehaviorTracker() *AIBehaviorTracker {
	return &AIBehaviorTracker{
		validationEvents:  make(chan ValidationEvent, 1000),
		behaviorPatterns:  NewPatternAnalyzer(),
		complianceScoring: NewComplianceScorer(),
		eventHistory:      make([]ValidationEvent, 0),
	}
}

// NewPatternAnalyzer creates a new pattern analyzer
func NewPatternAnalyzer() *PatternAnalyzer {
	return &PatternAnalyzer{
		patterns: make(map[string]*BehaviorPattern),
	}
}

// NewComplianceScorer creates a new compliance scorer
func NewComplianceScorer() *ComplianceScorer {
	return &ComplianceScorer{
		weights: map[string]float64{
			"pass_rate":       0.4,
			"error_frequency": 0.3,
			"warning_rate":    0.2,
			"consistency":     0.1,
		},
	}
}

// NewAdherenceMetrics creates a new adherence metrics calculator
func NewAdherenceMetrics() *AdherenceMetrics {
	return &AdherenceMetrics{
		assistantMetrics: make(map[string]*AssistantMetrics),
	}
}

// NewAuditTrailManager creates a new audit trail manager
func NewAuditTrailManager() *AuditTrailManager {
	return &AuditTrailManager{
		auditEntries: make([]AuditEntry, 0),
	}
}

// NewComplianceDashboard creates a new compliance dashboard
func NewComplianceDashboard() *ComplianceDashboard {
	return &ComplianceDashboard{
		Metrics:  make(map[string]*AssistantMetrics),
		Patterns: make(map[string]*BehaviorPattern),
	}
}

// RecordValidationEvent records a validation event for tracking
func (cm *ComplianceMonitor) RecordValidationEvent(event ValidationEvent) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 🔺 DOC-011: Event recording - 📊 Track AI assistant validation behavior
	cm.behaviorTracker.RecordEvent(event)
	cm.adherenceCalculator.UpdateMetrics(event)

	// Create audit entry
	auditEntry := AuditEntry{
		EntryID:     fmt.Sprintf("validation-%d", time.Now().Unix()),
		AssistantID: event.AssistantID,
		Action:      "validation_request",
		Details: map[string]interface{}{
			"status":           event.Status,
			"error_count":      event.ErrorCount,
			"warning_count":    event.WarningCount,
			"compliance_score": event.ComplianceScore,
			"validation_mode":  event.ValidationMode,
		},
		Timestamp: event.Timestamp,
		Severity:  cm.determineSeverity(event),
	}

	cm.auditTrailManager.AddEntry(auditEntry)
}

// RecordEvent records a validation event in the behavior tracker
func (bt *AIBehaviorTracker) RecordEvent(event ValidationEvent) {
	bt.mu.Lock()
	defer bt.mu.Unlock()

	// 🔺 DOC-011: Behavior tracking - 🔍 Store validation events for analysis
	bt.eventHistory = append(bt.eventHistory, event)

	// Analyze patterns
	bt.behaviorPatterns.AnalyzeEvent(event)

	// Send to processing channel (non-blocking)
	select {
	case bt.validationEvents <- event:
	default:
		// Channel full, skip this event to prevent blocking
	}
}

// AnalyzeEvent analyzes an event for behavior patterns
func (pa *PatternAnalyzer) AnalyzeEvent(event ValidationEvent) {
	pa.mu.Lock()
	defer pa.mu.Unlock()

	// 🔺 DOC-011: Pattern analysis - 📝 Detect AI assistant behavior patterns
	patternKey := fmt.Sprintf("%s_%s", event.AssistantID, event.Status)

	if pattern, exists := pa.patterns[patternKey]; exists {
		// Update existing pattern
		pattern.Frequency = pa.updateFrequency(pattern.Frequency, event)
		pattern.LastOccurrence = event.Timestamp
		pattern.TrendDirection = pa.calculateTrend(pattern, event)
		pattern.Confidence = pa.calculateConfidence(pattern)
	} else {
		// Create new pattern
		pa.patterns[patternKey] = &BehaviorPattern{
			AssistantID:    event.AssistantID,
			PatternType:    event.Status,
			Frequency:      1.0,
			LastOccurrence: event.Timestamp,
			TrendDirection: "new",
			Confidence:     0.1, // Low confidence for new patterns
		}
	}
}

// UpdateMetrics updates adherence metrics for an AI assistant
func (am *AdherenceMetrics) UpdateMetrics(event ValidationEvent) {
	am.mu.Lock()
	defer am.mu.Unlock()

	// 🔺 DOC-011: Metrics update - 📊 Calculate AI assistant adherence metrics
	if metrics, exists := am.assistantMetrics[event.AssistantID]; exists {
		// Update existing metrics
		metrics.TotalValidations++
		if event.Status == "pass" {
			metrics.SuccessfulPasses++
		} else if event.Status == "fail" {
			metrics.FailureCount++
		}
		metrics.WarningCount += event.WarningCount
		metrics.AverageScore = am.calculateAverageScore(metrics, event.ComplianceScore)
		metrics.LastValidation = event.Timestamp
		metrics.ComplianceRating = am.calculateComplianceRating(metrics)
	} else {
		// Create new metrics
		successCount := 0
		failureCount := 0
		if event.Status == "pass" {
			successCount = 1
		} else if event.Status == "fail" {
			failureCount = 1
		}

		am.assistantMetrics[event.AssistantID] = &AssistantMetrics{
			AssistantID:      event.AssistantID,
			TotalValidations: 1,
			SuccessfulPasses: successCount,
			FailureCount:     failureCount,
			WarningCount:     event.WarningCount,
			AverageScore:     event.ComplianceScore,
			LastValidation:   event.Timestamp,
			ComplianceRating: am.calculateInitialRating(event.ComplianceScore),
		}
	}
}

// GenerateDashboard generates a compliance dashboard
func (cm *ComplianceMonitor) GenerateDashboard() *ComplianceDashboard {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// 🔺 DOC-011: Dashboard generation - 📝 Create compliance visualization
	dashboard := &ComplianceDashboard{
		Metrics:     cm.adherenceCalculator.GetAllMetrics(),
		Patterns:    cm.behaviorTracker.behaviorPatterns.GetAllPatterns(),
		Summary:     cm.generateComplianceSummary(),
		GeneratedAt: time.Now(),
	}

	return dashboard
}

// AddEntry adds an audit trail entry
func (atm *AuditTrailManager) AddEntry(entry AuditEntry) {
	atm.mu.Lock()
	defer atm.mu.Unlock()

	atm.auditEntries = append(atm.auditEntries, entry)
}

// Helper methods

func (cm *ComplianceMonitor) determineSeverity(event ValidationEvent) string {
	if event.Status == "fail" || event.ErrorCount > 5 {
		return "high"
	} else if event.Status == "warning" || event.WarningCount > 10 {
		return "medium"
	}
	return "low"
}

func (pa *PatternAnalyzer) updateFrequency(currentFreq float64, event ValidationEvent) float64 {
	// Simple frequency update - can be enhanced with more sophisticated algorithms
	return currentFreq + 1.0
}

func (pa *PatternAnalyzer) calculateTrend(pattern *BehaviorPattern, event ValidationEvent) string {
	// Simplified trend calculation
	if event.ComplianceScore > 0.8 {
		return "improving"
	} else if event.ComplianceScore < 0.5 {
		return "declining"
	}
	return "stable"
}

func (pa *PatternAnalyzer) calculateConfidence(pattern *BehaviorPattern) float64 {
	// Confidence increases with frequency, capped at 1.0
	confidence := pattern.Frequency / 10.0
	if confidence > 1.0 {
		confidence = 1.0
	}
	return confidence
}

func (am *AdherenceMetrics) calculateAverageScore(metrics *AssistantMetrics, newScore float64) float64 {
	totalScore := metrics.AverageScore*float64(metrics.TotalValidations-1) + newScore
	return totalScore / float64(metrics.TotalValidations)
}

func (am *AdherenceMetrics) calculateComplianceRating(metrics *AssistantMetrics) string {
	passRate := float64(metrics.SuccessfulPasses) / float64(metrics.TotalValidations)

	if passRate >= 0.95 && metrics.AverageScore >= 0.9 {
		return "excellent"
	} else if passRate >= 0.85 && metrics.AverageScore >= 0.8 {
		return "good"
	} else if passRate >= 0.70 && metrics.AverageScore >= 0.6 {
		return "satisfactory"
	}
	return "needs_improvement"
}

func (am *AdherenceMetrics) calculateInitialRating(score float64) string {
	if score >= 0.9 {
		return "excellent"
	} else if score >= 0.8 {
		return "good"
	} else if score >= 0.6 {
		return "satisfactory"
	}
	return "needs_improvement"
}

func (am *AdherenceMetrics) GetAllMetrics() map[string]*AssistantMetrics {
	am.mu.RLock()
	defer am.mu.RUnlock()

	result := make(map[string]*AssistantMetrics)
	for k, v := range am.assistantMetrics {
		result[k] = v
	}
	return result
}

func (pa *PatternAnalyzer) GetAllPatterns() map[string]*BehaviorPattern {
	pa.mu.RLock()
	defer pa.mu.RUnlock()

	result := make(map[string]*BehaviorPattern)
	for k, v := range pa.patterns {
		result[k] = v
	}
	return result
}

func (cm *ComplianceMonitor) generateComplianceSummary() *ComplianceSummary {
	metrics := cm.adherenceCalculator.GetAllMetrics()

	totalAssistants := len(metrics)
	totalCompliance := 0.0
	issueCount := make(map[string]int)

	for _, metric := range metrics {
		totalCompliance += metric.AverageScore
		if metric.ComplianceRating == "needs_improvement" {
			issueCount["poor_compliance"]++
		}
		if metric.FailureCount > 5 {
			issueCount["high_failure_rate"]++
		}
	}

	overallCompliance := 0.0
	if totalAssistants > 0 {
		overallCompliance = totalCompliance / float64(totalAssistants)
	}

	return &ComplianceSummary{
		TotalAssistants:    totalAssistants,
		OverallCompliance:  overallCompliance,
		TrendDirection:     "stable", // Simplified for now
		TopIssues:          cm.extractTopIssues(issueCount),
		RecommendedActions: cm.generateRecommendations(overallCompliance),
	}
}

func (cm *ComplianceMonitor) extractTopIssues(issueCount map[string]int) []string {
	var issues []string
	for issue, count := range issueCount {
		if count > 0 {
			issues = append(issues, fmt.Sprintf("%s (%d occurrences)", issue, count))
		}
	}
	return issues
}

func (cm *ComplianceMonitor) generateRecommendations(overallCompliance float64) []string {
	if overallCompliance < 0.6 {
		return []string{
			"Implement additional AI assistant training",
			"Review and update validation guidelines",
			"Consider automated remediation tools",
		}
	} else if overallCompliance < 0.8 {
		return []string{
			"Focus on consistent validation practices",
			"Enhance error reporting for better AI understanding",
		}
	}
	return []string{
		"Maintain current compliance practices",
		"Consider advanced validation features",
	}
}
