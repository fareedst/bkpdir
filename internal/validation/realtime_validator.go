// 🔶 DOC-012: Real-time validation service - ⚡ Live validation with sub-second response times
package validation

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// RealTimeValidator provides live validation services with sub-second response times
type RealTimeValidator struct {
	gateway        *AIValidationGateway
	cache          *ValidationCache
	subscribers    map[string]*ValidationSubscriber
	mu             sync.RWMutex
	metricsTracker *PerformanceMetrics
}

// ValidationCache provides intelligent caching for performance optimization
type ValidationCache struct {
	entries map[string]*CacheEntry
	mu      sync.RWMutex
	ttl     time.Duration
}

// CacheEntry represents a cached validation result
type CacheEntry struct {
	Response  *ValidationResponse
	Timestamp time.Time
	Hash      string
}

// ValidationSubscriber represents a real-time validation subscriber
type ValidationSubscriber struct {
	ID       string
	Channel  chan *RealTimeValidationUpdate
	FileSet  map[string]bool
	LastSeen time.Time
}

// RealTimeValidationUpdate represents a real-time validation update
type RealTimeValidationUpdate struct {
	Type            string                     `json:"type"` // "validation", "suggestion", "status"
	File            string                     `json:"file"`
	Status          string                     `json:"status"`
	Errors          []AIOptimizedError         `json:"errors"`
	Warnings        []AIOptimizedWarning       `json:"warnings"`
	Suggestions     []IntelligentSuggestion    `json:"suggestions"`
	StatusIndicator *ValidationStatusIndicator `json:"status_indicator"`
	ProcessingTime  time.Duration              `json:"processing_time"`
	Timestamp       time.Time                  `json:"timestamp"`
}

// IntelligentSuggestion represents an intelligent correction suggestion
type IntelligentSuggestion struct {
	SuggestionID string            `json:"suggestion_id"`
	Type         string            `json:"type"` // "icon_fix", "token_format", "priority_correct"
	Original     string            `json:"original"`
	Suggested    string            `json:"suggested"`
	Confidence   float64           `json:"confidence"`
	Reason       string            `json:"reason"`
	FileLocation *FileLocation     `json:"file_location"`
	AutoApply    bool              `json:"auto_apply"`
	Context      map[string]string `json:"context"`
}

// ValidationStatusIndicator represents visual status indicators
type ValidationStatusIndicator struct {
	OverallStatus   string            `json:"overall_status"` // "pass", "warning", "error", "processing"
	FileStatus      map[string]string `json:"file_status"`
	ErrorCount      int               `json:"error_count"`
	WarningCount    int               `json:"warning_count"`
	ComplianceLevel string            `json:"compliance_level"` // "excellent", "good", "needs_work", "poor"
	VisualElements  *VisualElements   `json:"visual_elements"`
}

// VisualElements represents visual feedback elements
type VisualElements struct {
	StatusColor   string `json:"status_color"`
	StatusIcon    string `json:"status_icon"`
	ProgressBar   int    `json:"progress_bar"` // 0-100 percentage
	Tooltip       string `json:"tooltip"`
	BadgeText     string `json:"badge_text"`
	AnimationType string `json:"animation_type"`
}

// PerformanceMetrics tracks real-time validation performance
type PerformanceMetrics struct {
	mu                   sync.RWMutex
	totalValidations     int64
	averageResponseTime  time.Duration
	cacheHitRate         float64
	activeSubscribers    int
	validationsPerSecond float64
	lastMetricsUpdate    time.Time
}

// NewRealTimeValidator creates a new real-time validator
func NewRealTimeValidator() *RealTimeValidator {
	return &RealTimeValidator{
		gateway:     NewAIValidationGateway(),
		cache:       NewValidationCache(5 * time.Minute), // 5-minute TTL
		subscribers: make(map[string]*ValidationSubscriber),
		metricsTracker: &PerformanceMetrics{
			lastMetricsUpdate: time.Now(),
		},
	}
}

// NewValidationCache creates a new validation cache
func NewValidationCache(ttl time.Duration) *ValidationCache {
	cache := &ValidationCache{
		entries: make(map[string]*CacheEntry),
		ttl:     ttl,
	}

	// 🔶 DOC-012: Performance optimization - ⚡ Automatic cache cleanup
	go cache.cleanupWorker()

	return cache
}

// ValidateRealtimeFile provides real-time validation for a single file
func (rtv *RealTimeValidator) ValidateRealtimeFile(ctx context.Context, filePath string, content string) (*RealTimeValidationUpdate, error) {
	startTime := time.Now()

	// 🔶 DOC-012: Performance optimization - 🔍 Cache lookup for speed
	if cached := rtv.cache.Get(filePath, content); cached != nil {
		rtv.metricsTracker.RecordCacheHit()
		return rtv.createUpdateFromResponse(filePath, cached, time.Since(startTime)), nil
	}

	// 🔶 DOC-012: Live validation processing - 🔧 Real-time validation execution
	request := ValidationRequest{
		SourceFiles:    []string{filePath},
		ValidationMode: "realtime", // Special mode for real-time validation
		RequestContext: &AIRequestContext{
			AssistantID: "realtime-validator",
			SessionID:   fmt.Sprintf("rt-%d", time.Now().UnixNano()),
			Timestamp:   time.Now(),
		},
	}

	response, err := rtv.gateway.ProcessValidationRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("🔶 DOC-012: Real-time validation failed: %w", err)
	}

	// Cache the result for performance optimization
	rtv.cache.Set(filePath, content, response)

	// Generate intelligent suggestions
	suggestions := rtv.generateIntelligentSuggestions(filePath, content, response)

	// Create status indicator
	statusIndicator := rtv.createStatusIndicator(response)

	processingTime := time.Since(startTime)
	rtv.metricsTracker.RecordValidation(processingTime)

	update := &RealTimeValidationUpdate{
		Type:            "validation",
		File:            filePath,
		Status:          response.Status,
		Errors:          response.Errors,
		Warnings:        response.Warnings,
		Suggestions:     suggestions,
		StatusIndicator: statusIndicator,
		ProcessingTime:  processingTime,
		Timestamp:       time.Now(),
	}

	// Notify subscribers
	rtv.notifySubscribers(filePath, update)

	return update, nil
}

// SubscribeToValidation subscribes to real-time validation updates for specific files
func (rtv *RealTimeValidator) SubscribeToValidation(subscriberID string, files []string) <-chan *RealTimeValidationUpdate {
	rtv.mu.Lock()
	defer rtv.mu.Unlock()

	// 🔶 DOC-012: Editor integration - 🔧 Real-time subscription management
	subscriber := &ValidationSubscriber{
		ID:       subscriberID,
		Channel:  make(chan *RealTimeValidationUpdate, 100),
		FileSet:  make(map[string]bool),
		LastSeen: time.Now(),
	}

	for _, file := range files {
		subscriber.FileSet[file] = true
	}

	rtv.subscribers[subscriberID] = subscriber
	rtv.metricsTracker.IncrementActiveSubscribers()

	// Start cleanup routine for inactive subscribers
	go rtv.cleanupInactiveSubscriber(subscriberID)

	return subscriber.Channel
}

// UnsubscribeFromValidation unsubscribes from real-time validation updates
func (rtv *RealTimeValidator) UnsubscribeFromValidation(subscriberID string) {
	rtv.mu.Lock()
	defer rtv.mu.Unlock()

	if subscriber, exists := rtv.subscribers[subscriberID]; exists {
		close(subscriber.Channel)
		delete(rtv.subscribers, subscriberID)
		rtv.metricsTracker.DecrementActiveSubscribers()
	}
}

// GetValidationStatusIndicator provides current validation status for display
func (rtv *RealTimeValidator) GetValidationStatusIndicator(files []string) *ValidationStatusIndicator {
	// 🔶 DOC-012: Status indicators - 📊 Visual compliance status
	fileStatus := make(map[string]string)
	totalErrors := 0
	totalWarnings := 0

	for _, file := range files {
		if cached := rtv.cache.GetByFile(file); cached != nil {
			fileStatus[file] = cached.Status
			totalErrors += len(cached.Errors)
			totalWarnings += len(cached.Warnings)
		} else {
			fileStatus[file] = "unknown"
		}
	}

	overallStatus := rtv.determineOverallStatus(fileStatus)
	complianceLevel := rtv.calculateComplianceLevel(totalErrors, totalWarnings, len(files))

	return &ValidationStatusIndicator{
		OverallStatus:   overallStatus,
		FileStatus:      fileStatus,
		ErrorCount:      totalErrors,
		WarningCount:    totalWarnings,
		ComplianceLevel: complianceLevel,
		VisualElements:  rtv.createVisualElements(overallStatus, complianceLevel),
	}
}

// StartValidationServer starts the real-time validation HTTP server
func (rtv *RealTimeValidator) StartValidationServer(port int) error {
	// 🔶 DOC-012: Live validation service - 🔧 HTTP API for real-time validation
	http.HandleFunc("/validate", rtv.handleValidationRequest)
	http.HandleFunc("/subscribe", rtv.handleSubscriptionRequest)
	http.HandleFunc("/status", rtv.handleStatusRequest)
	http.HandleFunc("/suggestions", rtv.handleSuggestionsRequest)
	http.HandleFunc("/metrics", rtv.handleMetricsRequest)

	serverAddr := fmt.Sprintf(":%d", port)
	fmt.Printf("🔶 DOC-012: Real-time validation server starting on %s\n", serverAddr)

	return http.ListenAndServe(serverAddr, nil)
}

// generateIntelligentSuggestions creates intelligent correction suggestions
func (rtv *RealTimeValidator) generateIntelligentSuggestions(filePath string, content string, response *ValidationResponse) []IntelligentSuggestion {
	// 🔶 DOC-012: Intelligent corrections - 📝 Smart suggestions based on validation results
	var suggestions []IntelligentSuggestion

	for _, err := range response.Errors {
		if suggestion := rtv.createSuggestionFromError(filePath, content, err); suggestion != nil {
			suggestions = append(suggestions, *suggestion)
		}
	}

	for _, warning := range response.Warnings {
		if suggestion := rtv.createSuggestionFromWarning(filePath, content, warning); suggestion != nil {
			suggestions = append(suggestions, *suggestion)
		}
	}

	// Add proactive suggestions based on content analysis
	proactiveSuggestions := rtv.generateProactiveSuggestions(filePath, content)
	suggestions = append(suggestions, proactiveSuggestions...)

	return suggestions
}

// createSuggestionFromError creates a suggestion from a validation error
func (rtv *RealTimeValidator) createSuggestionFromError(filePath string, content string, error AIOptimizedError) *IntelligentSuggestion {
	// 🔶 DOC-012: Intelligent corrections - 🔧 Error-based suggestions
	switch error.Category {
	case "icon_format":
		return rtv.createIconFormatSuggestion(filePath, error)
	case "token_format":
		return rtv.createTokenFormatSuggestion(filePath, error)
	case "priority_mismatch":
		return rtv.createPriorityMismatchSuggestion(filePath, error)
	default:
		return rtv.createGenericSuggestion(filePath, error)
	}
}

// createSuggestionFromWarning creates a suggestion from a validation warning
func (rtv *RealTimeValidator) createSuggestionFromWarning(filePath string, content string, warning AIOptimizedWarning) *IntelligentSuggestion {
	// 🔶 DOC-012: Intelligent corrections - 📝 Warning-based suggestions
	return &IntelligentSuggestion{
		SuggestionID: fmt.Sprintf("warn-%s-%d", filepath.Base(filePath), time.Now().UnixNano()),
		Type:         "improvement",
		Original:     warning.Message,
		Suggested:    rtv.generateSuggestionFromWarning(warning),
		Confidence:   0.75,
		Reason:       "Addressing validation warning",
		FileLocation: warning.FileReference,
		AutoApply:    false,
		Context:      map[string]string{"category": warning.Category},
	}
}

// Helper methods for cache operations
func (vc *ValidationCache) Get(filePath string, content string) *ValidationResponse {
	vc.mu.RLock()
	defer vc.mu.RUnlock()

	key := vc.generateKey(filePath, content)
	if entry, exists := vc.entries[key]; exists {
		if time.Since(entry.Timestamp) < vc.ttl {
			return entry.Response
		}
		// Entry expired, remove it
		delete(vc.entries, key)
	}
	return nil
}

// Additional helper methods to complete the implementation

// Set stores a validation response in the cache
func (vc *ValidationCache) Set(filePath string, content string, response *ValidationResponse) {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	key := vc.generateKey(filePath, content)
	vc.entries[key] = &CacheEntry{
		Response:  response,
		Timestamp: time.Now(),
		Hash:      vc.generateContentHash(content),
	}
}

// GetByFile gets cached results for a specific file (used for status indicators)
func (vc *ValidationCache) GetByFile(filePath string) *ValidationResponse {
	vc.mu.RLock()
	defer vc.mu.RUnlock()

	// Find most recent entry for this file
	var latestEntry *CacheEntry
	var latestTime time.Time

	for key, entry := range vc.entries {
		if strings.Contains(key, filePath) {
			if entry.Timestamp.After(latestTime) {
				latestTime = entry.Timestamp
				latestEntry = entry
			}
		}
	}

	if latestEntry != nil && time.Since(latestEntry.Timestamp) < vc.ttl {
		return latestEntry.Response
	}
	return nil
}

// generateKey creates a cache key from file path and content
func (vc *ValidationCache) generateKey(filePath string, content string) string {
	hash := vc.generateContentHash(content)
	return fmt.Sprintf("%s:%s", filePath, hash[:8]) // Use first 8 chars of hash
}

// generateContentHash creates a hash of the file content
func (vc *ValidationCache) generateContentHash(content string) string {
	hasher := sha256.New()
	hasher.Write([]byte(content))
	return hex.EncodeToString(hasher.Sum(nil))
}

// cleanupWorker periodically cleans expired cache entries
func (vc *ValidationCache) cleanupWorker() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		vc.mu.Lock()
		now := time.Now()
		for key, entry := range vc.entries {
			if now.Sub(entry.Timestamp) > vc.ttl {
				delete(vc.entries, key)
			}
		}
		vc.mu.Unlock()
	}
}

// Performance metrics methods
func (pm *PerformanceMetrics) RecordCacheHit() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	// Simplified cache hit rate calculation
	pm.cacheHitRate = (pm.cacheHitRate + 1.0) / 2.0
}

func (pm *PerformanceMetrics) RecordValidation(duration time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.totalValidations++
	if pm.averageResponseTime == 0 {
		pm.averageResponseTime = duration
	} else {
		pm.averageResponseTime = (pm.averageResponseTime + duration) / 2
	}

	// Calculate validations per second
	elapsed := time.Since(pm.lastMetricsUpdate)
	if elapsed > time.Second {
		pm.validationsPerSecond = float64(pm.totalValidations) / elapsed.Seconds()
		pm.lastMetricsUpdate = time.Now()
	}
}

func (pm *PerformanceMetrics) IncrementActiveSubscribers() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.activeSubscribers++
}

func (pm *PerformanceMetrics) DecrementActiveSubscribers() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.activeSubscribers--
}

// Real-time validator helper methods
func (rtv *RealTimeValidator) createUpdateFromResponse(filePath string, response *ValidationResponse, processingTime time.Duration) *RealTimeValidationUpdate {
	// 🔶 DOC-012: Performance optimization - 🔧 Cached response processing
	suggestions := rtv.generateIntelligentSuggestions(filePath, "", response)
	statusIndicator := rtv.createStatusIndicator(response)

	return &RealTimeValidationUpdate{
		Type:            "validation",
		File:            filePath,
		Status:          response.Status,
		Errors:          response.Errors,
		Warnings:        response.Warnings,
		Suggestions:     suggestions,
		StatusIndicator: statusIndicator,
		ProcessingTime:  processingTime,
		Timestamp:       time.Now(),
	}
}

func (rtv *RealTimeValidator) createStatusIndicator(response *ValidationResponse) *ValidationStatusIndicator {
	// 🔶 DOC-012: Status indicators - 📊 Visual feedback creation
	errorCount := len(response.Errors)
	warningCount := len(response.Warnings)

	var overallStatus string
	switch response.Status {
	case "pass":
		overallStatus = "pass"
	case "warning":
		overallStatus = "warning"
	case "fail":
		overallStatus = "error"
	default:
		overallStatus = "unknown"
	}

	complianceLevel := rtv.calculateComplianceLevel(errorCount, warningCount, 1)

	return &ValidationStatusIndicator{
		OverallStatus:   overallStatus,
		FileStatus:      map[string]string{"current": overallStatus},
		ErrorCount:      errorCount,
		WarningCount:    warningCount,
		ComplianceLevel: complianceLevel,
		VisualElements:  rtv.createVisualElements(overallStatus, complianceLevel),
	}
}

func (rtv *RealTimeValidator) notifySubscribers(filePath string, update *RealTimeValidationUpdate) {
	// 🔶 DOC-012: Editor integration - 📝 Real-time notification system
	rtv.mu.RLock()
	defer rtv.mu.RUnlock()

	for _, subscriber := range rtv.subscribers {
		if subscriber.FileSet[filePath] || len(subscriber.FileSet) == 0 {
			select {
			case subscriber.Channel <- update:
				subscriber.LastSeen = time.Now()
			default:
				// Channel is full, skip this update to prevent blocking
			}
		}
	}
}

func (rtv *RealTimeValidator) cleanupInactiveSubscriber(subscriberID string) {
	// Clean up inactive subscribers after 30 minutes
	time.Sleep(30 * time.Minute)

	rtv.mu.Lock()
	defer rtv.mu.Unlock()

	if subscriber, exists := rtv.subscribers[subscriberID]; exists {
		if time.Since(subscriber.LastSeen) > 30*time.Minute {
			close(subscriber.Channel)
			delete(rtv.subscribers, subscriberID)
			rtv.metricsTracker.DecrementActiveSubscribers()
		}
	}
}

func (rtv *RealTimeValidator) determineOverallStatus(fileStatus map[string]string) string {
	// 🔶 DOC-012: Status indicators - 🔍 Overall status calculation
	hasErrors := false
	hasWarnings := false

	for _, status := range fileStatus {
		switch status {
		case "error", "fail":
			hasErrors = true
		case "warning":
			hasWarnings = true
		}
	}

	if hasErrors {
		return "error"
	}
	if hasWarnings {
		return "warning"
	}
	return "pass"
}

func (rtv *RealTimeValidator) calculateComplianceLevel(errorCount, warningCount, fileCount int) string {
	// 🔶 DOC-012: Status indicators - 📊 Compliance level calculation
	if errorCount == 0 && warningCount == 0 {
		return "excellent"
	}
	if errorCount == 0 && warningCount <= fileCount {
		return "good"
	}
	if errorCount <= fileCount {
		return "needs_work"
	}
	return "poor"
}

func (rtv *RealTimeValidator) createVisualElements(overallStatus, complianceLevel string) *VisualElements {
	// 🔶 DOC-012: Status indicators - 🎨 Visual element creation
	elements := &VisualElements{}

	switch overallStatus {
	case "pass":
		elements.StatusColor = "#28a745" // Green
		elements.StatusIcon = "✅"
		elements.AnimationType = "pulse-success"
	case "warning":
		elements.StatusColor = "#ffc107" // Yellow
		elements.StatusIcon = "⚠️"
		elements.AnimationType = "pulse-warning"
	case "error":
		elements.StatusColor = "#dc3545" // Red
		elements.StatusIcon = "❌"
		elements.AnimationType = "pulse-error"
	default:
		elements.StatusColor = "#6c757d" // Gray
		elements.StatusIcon = "❓"
		elements.AnimationType = "none"
	}

	// Set progress bar based on compliance level
	switch complianceLevel {
	case "excellent":
		elements.ProgressBar = 100
		elements.BadgeText = "Perfect"
	case "good":
		elements.ProgressBar = 80
		elements.BadgeText = "Good"
	case "needs_work":
		elements.ProgressBar = 60
		elements.BadgeText = "Issues"
	case "poor":
		elements.ProgressBar = 30
		elements.BadgeText = "Poor"
	}

	elements.Tooltip = fmt.Sprintf("Compliance: %s", complianceLevel)

	return elements
}

// Intelligent suggestion creation methods
func (rtv *RealTimeValidator) createIconFormatSuggestion(filePath string, error AIOptimizedError) *IntelligentSuggestion {
	// 🔶 DOC-012: Intelligent corrections - 🔧 Icon format suggestions
	return &IntelligentSuggestion{
		SuggestionID: fmt.Sprintf("icon-%s-%d", filepath.Base(filePath), time.Now().UnixNano()),
		Type:         "icon_fix",
		Original:     error.Message,
		Suggested:    rtv.generateIconFormatFix(error),
		Confidence:   0.9,
		Reason:       "Icon format standardization",
		FileLocation: error.FileReference,
		AutoApply:    true,
		Context:      error.Context,
	}
}

func (rtv *RealTimeValidator) createTokenFormatSuggestion(filePath string, error AIOptimizedError) *IntelligentSuggestion {
	// 🔶 DOC-012: Intelligent corrections - 📝 Token format suggestions
	return &IntelligentSuggestion{
		SuggestionID: fmt.Sprintf("token-%s-%d", filepath.Base(filePath), time.Now().UnixNano()),
		Type:         "token_format",
		Original:     error.Message,
		Suggested:    rtv.generateTokenFormatFix(error),
		Confidence:   0.85,
		Reason:       "Token format standardization",
		FileLocation: error.FileReference,
		AutoApply:    true,
		Context:      error.Context,
	}
}

func (rtv *RealTimeValidator) createPriorityMismatchSuggestion(filePath string, error AIOptimizedError) *IntelligentSuggestion {
	// 🔶 DOC-012: Intelligent corrections - ⭐ Priority correction suggestions
	return &IntelligentSuggestion{
		SuggestionID: fmt.Sprintf("priority-%s-%d", filepath.Base(filePath), time.Now().UnixNano()),
		Type:         "priority_correct",
		Original:     error.Message,
		Suggested:    rtv.generatePriorityFix(error),
		Confidence:   0.8,
		Reason:       "Priority icon alignment",
		FileLocation: error.FileReference,
		AutoApply:    false, // Priority changes need manual review
		Context:      error.Context,
	}
}

func (rtv *RealTimeValidator) createGenericSuggestion(filePath string, error AIOptimizedError) *IntelligentSuggestion {
	// 🔶 DOC-012: Intelligent corrections - 🔧 Generic suggestions
	return &IntelligentSuggestion{
		SuggestionID: fmt.Sprintf("generic-%s-%d", filepath.Base(filePath), time.Now().UnixNano()),
		Type:         "improvement",
		Original:     error.Message,
		Suggested:    fmt.Sprintf("Review %s validation error", error.Category),
		Confidence:   0.5,
		Reason:       "General validation improvement",
		FileLocation: error.FileReference,
		AutoApply:    false,
		Context:      error.Context,
	}
}

func (rtv *RealTimeValidator) generateProactiveSuggestions(filePath string, content string) []IntelligentSuggestion {
	// 🔶 DOC-012: Intelligent corrections - 🔍 Proactive suggestion generation
	var suggestions []IntelligentSuggestion

	// Look for potential improvements in token format
	if strings.Contains(content, "//") && strings.Contains(content, ":") {
		tokenPattern := regexp.MustCompile(`//\s*([A-Z]+-[0-9]+):(.*)`)
		matches := tokenPattern.FindAllStringSubmatch(content, -1)

		for _, match := range matches {
			if len(match) >= 3 {
				// Check if token could be enhanced with icons
				if !strings.Contains(match[0], "🔶") && !strings.Contains(match[0], "⭐") {
					suggestions = append(suggestions, IntelligentSuggestion{
						SuggestionID: fmt.Sprintf("proactive-%s-%d", filepath.Base(filePath), time.Now().UnixNano()),
						Type:         "icon_enhancement",
						Original:     match[0],
						Suggested:    fmt.Sprintf("// 🔶 %s:%s", match[1], match[2]),
						Confidence:   0.7,
						Reason:       "Add priority icon for better visibility",
						FileLocation: nil, // Would need line number detection
						AutoApply:    false,
						Context:      map[string]string{"type": "enhancement"},
					})
				}
			}
		}
	}

	return suggestions
}

// Fix generation helper methods
func (rtv *RealTimeValidator) generateIconFormatFix(error AIOptimizedError) string {
	// 🔶 DOC-012: Intelligent corrections - 🔧 Icon format fixes
	if strings.Contains(error.Message, "missing priority icon") {
		return "Add appropriate priority icon (⭐🔺🔶🔻) based on feature importance"
	}
	if strings.Contains(error.Message, "missing action icon") {
		return "Add appropriate action icon (🔍📝🔧🛡️) based on function behavior"
	}
	return "Review icon format requirements in documentation"
}

func (rtv *RealTimeValidator) generateTokenFormatFix(error AIOptimizedError) string {
	// 🔶 DOC-012: Intelligent corrections - 📝 Token format fixes
	if strings.Contains(error.Message, "token format") {
		return "Use format: // 🔶 FEATURE-ID: Description - 🔧 Action description"
	}
	return "Review token format requirements in DOC-008 documentation"
}

func (rtv *RealTimeValidator) generatePriorityFix(error AIOptimizedError) string {
	// 🔶 DOC-012: Intelligent corrections - ⭐ Priority alignment fixes
	return "Review feature priority in feature-tracking.md and use appropriate icon: ⭐ Critical, 🔺 High, 🔶 Medium, 🔻 Low"
}

func (rtv *RealTimeValidator) generateSuggestionFromWarning(warning AIOptimizedWarning) string {
	// 🔶 DOC-012: Intelligent corrections - 📝 Warning-based improvements
	switch warning.Category {
	case "token_consistency":
		return "Consider standardizing token format for better consistency"
	case "icon_usage":
		return "Review icon usage guidelines for optimal visual feedback"
	default:
		return fmt.Sprintf("Address %s warning for improved code quality", warning.Category)
	}
}

// HTTP handlers for the real-time validation server
func (rtv *RealTimeValidator) handleValidationRequest(w http.ResponseWriter, r *http.Request) {
	// 🔶 DOC-012: Live validation service - 🔧 HTTP validation endpoint
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		FilePath string `json:"file_path"`
		Content  string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update, err := rtv.ValidateRealtimeFile(ctx, request.FilePath, request.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(update)
}

func (rtv *RealTimeValidator) handleSubscriptionRequest(w http.ResponseWriter, r *http.Request) {
	// 🔶 DOC-012: Editor integration - 🔧 WebSocket subscription endpoint
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "WebSocket subscription endpoint - implement WebSocket upgrade here",
		"status":  "placeholder",
	})
}

func (rtv *RealTimeValidator) handleStatusRequest(w http.ResponseWriter, r *http.Request) {
	// 🔶 DOC-012: Status indicators - 📊 Status endpoint
	files := r.URL.Query()["files"]
	if len(files) == 0 {
		http.Error(w, "No files specified", http.StatusBadRequest)
		return
	}

	indicator := rtv.GetValidationStatusIndicator(files)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(indicator)
}

func (rtv *RealTimeValidator) handleSuggestionsRequest(w http.ResponseWriter, r *http.Request) {
	// 🔶 DOC-012: Intelligent corrections - 📝 Suggestions endpoint
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Intelligent suggestions endpoint",
		"status":  "available",
	})
}

func (rtv *RealTimeValidator) handleMetricsRequest(w http.ResponseWriter, r *http.Request) {
	// 🔶 DOC-012: Performance optimization - 📊 Metrics endpoint
	rtv.metricsTracker.mu.RLock()
	defer rtv.metricsTracker.mu.RUnlock()

	metrics := map[string]interface{}{
		"total_validations":      rtv.metricsTracker.totalValidations,
		"average_response_time":  rtv.metricsTracker.averageResponseTime.String(),
		"cache_hit_rate":         rtv.metricsTracker.cacheHitRate,
		"active_subscribers":     rtv.metricsTracker.activeSubscribers,
		"validations_per_second": rtv.metricsTracker.validationsPerSecond,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
