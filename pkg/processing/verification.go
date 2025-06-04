// ⭐ EXTRACT-007: Verification systems - Data integrity checking patterns extracted from verify.go - 🔧
package processing

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"time"
)

// VerificationProviderInterface defines the interface for data integrity checking
type VerificationProviderInterface interface {
	Calculate(data io.Reader) (string, error)
	Verify(data io.Reader, expected string) (bool, error)
	GetAlgorithm() string
	GetDisplayName() string
}

// VerificationResult represents the result of a verification operation
type VerificationResult struct {
	Algorithm     string        `json:"algorithm"`
	Checksum      string        `json:"checksum"`
	Expected      string        `json:"expected,omitempty"`
	IsValid       bool          `json:"is_valid"`
	VerifiedAt    time.Time     `json:"verified_at"`
	Duration      time.Duration `json:"duration"`
	BytesVerified int64         `json:"bytes_verified"`
	Errors        []string      `json:"errors,omitempty"`
}

// VerificationStatus represents verification status (extracted from verify.go)
type VerificationStatus struct {
	VerifiedAt   time.Time         `json:"verified_at"`
	IsVerified   bool              `json:"is_verified"`
	HasChecksums bool              `json:"has_checksums"`
	Algorithm    string            `json:"algorithm,omitempty"`
	Checksums    map[string]string `json:"checksums,omitempty"`
	Errors       []string          `json:"errors,omitempty"`
}

// BaseVerificationProvider provides common verification functionality
type BaseVerificationProvider struct {
	algorithm   string
	displayName string
	hasher      func() hash.Hash
}

// NewSHA256Verifier creates a SHA-256 verification provider (extracted from verify.go)
func NewSHA256Verifier() VerificationProviderInterface {
	return &BaseVerificationProvider{
		algorithm:   "sha256",
		displayName: "SHA-256",
		hasher:      sha256.New,
	}
}

// NewSHA512Verifier creates a SHA-512 verification provider
func NewSHA512Verifier() VerificationProviderInterface {
	return &BaseVerificationProvider{
		algorithm:   "sha512",
		displayName: "SHA-512",
		hasher:      sha512.New,
	}
}

// NewMD5Verifier creates an MD5 verification provider
func NewMD5Verifier() VerificationProviderInterface {
	return &BaseVerificationProvider{
		algorithm:   "md5",
		displayName: "MD5",
		hasher:      md5.New,
	}
}

// Calculate computes the checksum for the provided data
func (bvp *BaseVerificationProvider) Calculate(data io.Reader) (string, error) {
	if data == nil {
		return "", NewProcessingError("INVALID_INPUT", "Calculate", "data reader cannot be nil")
	}

	hasher := bvp.hasher()

	// Copy data to hasher
	_, err := io.Copy(hasher, data)
	if err != nil {
		return "", NewProcessingError("CALCULATION_FAILED", "Calculate", fmt.Sprintf("failed to read data: %v", err))
	}

	// Get checksum
	checksum := hex.EncodeToString(hasher.Sum(nil))

	return checksum, nil
}

// Verify compares the checksum of data against an expected value
func (bvp *BaseVerificationProvider) Verify(data io.Reader, expected string) (bool, error) {
	if expected == "" {
		return false, NewProcessingError("INVALID_EXPECTED", "Verify", "expected checksum cannot be empty")
	}

	calculated, err := bvp.Calculate(data)
	if err != nil {
		return false, err
	}

	return calculated == expected, nil
}

// GetAlgorithm returns the algorithm name
func (bvp *BaseVerificationProvider) GetAlgorithm() string {
	return bvp.algorithm
}

// GetDisplayName returns the human-readable algorithm name
func (bvp *BaseVerificationProvider) GetDisplayName() string {
	return bvp.displayName
}

// VerificationManager manages multiple verification providers
type VerificationManager struct {
	providers        map[string]VerificationProviderInterface
	defaultAlgorithm string
}

// NewVerificationManager creates a new verification manager with default providers
func NewVerificationManager() *VerificationManager {
	vm := &VerificationManager{
		providers:        make(map[string]VerificationProviderInterface),
		defaultAlgorithm: "sha256",
	}

	// Register default providers
	vm.RegisterProvider(NewSHA256Verifier())
	vm.RegisterProvider(NewSHA512Verifier())
	vm.RegisterProvider(NewMD5Verifier())

	return vm
}

// RegisterProvider adds a verification provider
func (vm *VerificationManager) RegisterProvider(provider VerificationProviderInterface) {
	vm.providers[provider.GetAlgorithm()] = provider
}

// GetProvider returns a verification provider by algorithm
func (vm *VerificationManager) GetProvider(algorithm string) (VerificationProviderInterface, error) {
	if algorithm == "" {
		algorithm = vm.defaultAlgorithm
	}

	provider, exists := vm.providers[algorithm]
	if !exists {
		return nil, NewProcessingError("UNKNOWN_ALGORITHM", "GetProvider", fmt.Sprintf("unsupported algorithm: %s", algorithm))
	}

	return provider, nil
}

// GetSupportedAlgorithms returns the list of supported algorithms
func (vm *VerificationManager) GetSupportedAlgorithms() []string {
	algorithms := make([]string, 0, len(vm.providers))
	for algorithm := range vm.providers {
		algorithms = append(algorithms, algorithm)
	}
	return algorithms
}

// VerifyWithAlgorithm verifies data using the specified algorithm
func (vm *VerificationManager) VerifyWithAlgorithm(data io.Reader, expected, algorithm string) (*VerificationResult, error) {
	start := time.Now()

	provider, err := vm.GetProvider(algorithm)
	if err != nil {
		return nil, err
	}

	// Calculate checksum
	checksum, err := provider.Calculate(data)
	if err != nil {
		return &VerificationResult{
			Algorithm:  algorithm,
			IsValid:    false,
			VerifiedAt: time.Now(),
			Duration:   time.Since(start),
			Errors:     []string{err.Error()},
		}, err
	}

	// Compare checksums
	isValid := checksum == expected

	return &VerificationResult{
		Algorithm:  algorithm,
		Checksum:   checksum,
		Expected:   expected,
		IsValid:    isValid,
		VerifiedAt: time.Now(),
		Duration:   time.Since(start),
	}, nil
}

// GenerateChecksums creates checksums for multiple data sources (extracted from verify.go)
func (vm *VerificationManager) GenerateChecksums(fileMap map[string]io.Reader, algorithm string) (map[string]string, error) {
	provider, err := vm.GetProvider(algorithm)
	if err != nil {
		return nil, err
	}

	checksums := make(map[string]string)

	for name, reader := range fileMap {
		if reader == nil {
			continue
		}

		checksum, err := provider.Calculate(reader)
		if err != nil {
			return nil, NewProcessingError("CHECKSUM_GENERATION", "GenerateChecksums", fmt.Sprintf("failed to calculate checksum for %s: %v", name, err))
		}

		checksums[name] = checksum
	}

	return checksums, nil
}

// VerifyChecksums verifies multiple checksums (extracted from verify.go)
func (vm *VerificationManager) VerifyChecksums(fileMap map[string]io.Reader, expectedChecksums map[string]string, algorithm string) (*VerificationStatus, error) {
	status := &VerificationStatus{
		VerifiedAt:   time.Now(),
		IsVerified:   true,
		HasChecksums: len(expectedChecksums) > 0,
		Algorithm:    algorithm,
		Checksums:    make(map[string]string),
		Errors:       []string{},
	}

	if !status.HasChecksums {
		return status, nil
	}

	provider, err := vm.GetProvider(algorithm)
	if err != nil {
		status.IsVerified = false
		status.Errors = append(status.Errors, err.Error())
		return status, err
	}

	// Verify each file
	for name, reader := range fileMap {
		if reader == nil {
			continue
		}

		expected, hasExpected := expectedChecksums[name]
		if !hasExpected {
			status.Errors = append(status.Errors, fmt.Sprintf("no expected checksum for %s", name))
			continue
		}

		calculated, err := provider.Calculate(reader)
		if err != nil {
			status.IsVerified = false
			status.Errors = append(status.Errors, fmt.Sprintf("failed to calculate checksum for %s: %v", name, err))
			continue
		}

		status.Checksums[name] = calculated

		if calculated != expected {
			status.IsVerified = false
			status.Errors = append(status.Errors, fmt.Sprintf("checksum mismatch for %s: expected %s, got %s", name, expected, calculated))
		}
	}

	return status, nil
}

// SerializeChecksums serializes checksums to JSON (extracted from verify.go)
func SerializeChecksums(checksums map[string]string) ([]byte, error) {
	if checksums == nil {
		checksums = make(map[string]string)
	}

	data, err := json.Marshal(checksums)
	if err != nil {
		return nil, NewProcessingError("SERIALIZATION_FAILED", "SerializeChecksums", fmt.Sprintf("failed to serialize checksums: %v", err))
	}

	return data, nil
}

// DeserializeChecksums deserializes checksums from JSON (extracted from verify.go)
func DeserializeChecksums(data []byte) (map[string]string, error) {
	if len(data) == 0 {
		return make(map[string]string), nil
	}

	var checksums map[string]string
	err := json.Unmarshal(data, &checksums)
	if err != nil {
		return nil, NewProcessingError("DESERIALIZATION_FAILED", "DeserializeChecksums", fmt.Sprintf("failed to deserialize checksums: %v", err))
	}

	if checksums == nil {
		checksums = make(map[string]string)
	}

	return checksums, nil
}

// CreateVerificationStatus creates a new verification status (extracted from verify.go)
func CreateVerificationStatus(isVerified bool, hasChecksums bool, algorithm string) *VerificationStatus {
	return &VerificationStatus{
		VerifiedAt:   time.Now(),
		IsVerified:   isVerified,
		HasChecksums: hasChecksums,
		Algorithm:    algorithm,
		Checksums:    make(map[string]string),
		Errors:       []string{},
	}
}

// AddVerificationError adds an error to the verification status
func (vs *VerificationStatus) AddError(message string) {
	vs.IsVerified = false
	vs.Errors = append(vs.Errors, message)
}

// GetVerificationSummary returns a human-readable summary
func (vs *VerificationStatus) GetVerificationSummary() string {
	if vs.IsVerified {
		if vs.HasChecksums {
			return fmt.Sprintf("Verified successfully using %s (%d checksums)", vs.Algorithm, len(vs.Checksums))
		}
		return "Verified successfully (no checksums)"
	}

	errorCount := len(vs.Errors)
	if errorCount > 0 {
		return fmt.Sprintf("Verification failed with %d error(s)", errorCount)
	}

	return "Verification failed"
}
