# 🤖 AI Documentation Templates

> **[LOW] DOC-013**: AI-First Documentation Templates  
> **Implementation Tokens**: `// DOC-013: AI documentation templates`  
> **Purpose**: Standardized templates for AI assistants to create consistent, machine-readable documentation and code comments.

## 📑 Template Overview

### 🎯 Template Categories

This document provides AI assistants with standardized templates for:
- **Feature Documentation**: Complete feature specification and implementation documentation
- **Code Comments**: AI-optimized code comment structures with implementation tokens
- **Technical Documentation**: Architecture, requirements, and specification templates
- **Testing Documentation**: Test coverage and validation documentation templates

### 🔗 Integration Standards

All templates follow:
- **DOC-007**: Source Code Icon Integration standards
- **DOC-008**: Icon Validation and Enforcement requirements  
- **DOC-009**: Mass Implementation Token Standardization format
- **AI-Assistant-Compliance**: Token referencing and cross-reference requirements

## [ACTION:format-processing] Feature Documentation Templates

### 🚀 Complete Feature Documentation Template
```markdown
# [Feature Name]

> **[PRIORITY_ICON] FEATURE-ID**: [Brief one-sentence feature description]
> 
> **Implementation Tokens**: `// FEATURE-ID: [Implementation description]`
> 
> **Priority**: [[CRITICAL][HIGH][MEDIUM][LOW]] [CRITICAL|HIGH|MEDIUM|LOW]  
> **Status**: [[ACTION:format-processing] Not Started|[ACTION:migration] In Progress|✅ Completed]  
> **Dependencies**: [List major dependencies or "None"]

## 📑 Feature Overview

### 🎯 Purpose
[Clear, single-paragraph statement of feature purpose optimized for AI comprehension. Should answer: What does this feature do and why is it needed?]

### 🔗 Dependencies  
- **DEPENDENCY-001**: [Description] - See [link](path/to/dependency.md)
- **DEPENDENCY-002**: [Description] - See [link](path/to/dependency.md)
- **None**: [If no dependencies exist]

### [ACTION:core-functionality] Integration Points
- **[System/Component 1]**: [How this feature integrates]
- **[System/Component 2]**: [How this feature integrates]
- **Standalone**: [If no integration points exist]

## 🏗️ Technical Implementation

### 📋 Requirements
- **R-[FEATURE-ID]-001**: [Specific, testable requirement with clear acceptance criteria]
- **R-[FEATURE-ID]-002**: [Specific, testable requirement with clear acceptance criteria]
- **R-[FEATURE-ID]-003**: [Additional requirements as needed]

### [ACTION:core-functionality] Architecture
[Brief description of technical architecture and design decisions]

#### **Core Components**
- **Component 1**: [Purpose and responsibilities]
- **Component 2**: [Purpose and responsibilities]

#### **Data Flow**
```
[Input] → [Processing Step 1] → [Processing Step 2] → [Output]
```

### 💻 Code Implementation
```go
// [PRIORITY_ICON] FEATURE-ID: [Implementation description] - [ACTION_ICON] [Context]
func [FeatureFunctionName]([parameters]) ([returnType], error) {
    // [ACTION:discovery] Input validation with clear error context
    if [validation condition] {
        return [zero value], &StructuredError{
            Code:    "[ERROR_CODE]",
            Message: "[human-readable error message]",
            Context: map[string]interface{}{
                "function": "[FeatureFunctionName]",
                "input":    [relevant input values],
            },
        }
    }
    
    // [ACTION:core-functionality] Main feature operation
    [implementation logic]
    
    // ✅ Success return
    return [result], nil
}
```

## 🧪 Testing Requirements

### 📋 Test Coverage
- **Unit Tests**: [List specific test functions required]
- **Integration Tests**: [List integration test scenarios]
- **Performance Tests**: [Any performance requirements]
- **Error Path Tests**: [List error conditions to test]

### [ACTION:core-functionality] Test Implementation
```go
// [LOW] FEATURE-ID: Feature testing implementation - 🧪 Test coverage
func Test[FeatureName](t *testing.T) {
    tests := []struct {
        name     string
        input    [InputType]
        expected [ExpectedType]
        wantErr  bool
    }{
        {
            name:     "[Test case description]",
            input:    [test input],
            expected: [expected output],
            wantErr:  false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := [FeatureFunctionName](tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## ✅ Implementation Checklist

### 📋 Pre-Implementation
- [ ] Feature ID exists in `feature-tracking.md`
- [ ] Dependencies identified and documented
- [ ] Requirements clearly defined and testable
- [ ] Architecture design approved

### [ACTION:core-functionality] Implementation Phase
- [ ] Implementation tokens added to all code
- [ ] Error handling follows structured error patterns
- [ ] Code comments follow AI-first documentation standards
- [ ] Integration points properly implemented

### 🧪 Testing Phase
- [ ] Unit tests implemented and passing
- [ ] Integration tests implemented and passing
- [ ] Error path tests implemented and passing
- [ ] Performance tests implemented (if applicable)

### [ACTION:format-processing] Documentation Phase
- [ ] Feature documentation completed
- [ ] Cross-references updated in related documents
- [ ] `feature-tracking.md` status updated
- [ ] API documentation updated (if applicable)

### ✅ Completion
- [ ] All tests passing
- [ ] All lint checks passing
- [ ] DOC-008 icon validation passing
- [ ] Feature marked complete in `feature-tracking.md`

## 📊 Success Criteria
- **Functionality**: [Feature works as specified in requirements]
- **Performance**: [Any performance criteria, or "No specific requirements"]
- **Integration**: [Successfully integrates with dependent systems]
- **Documentation**: [Complete documentation following AI-first standards]
- **Testing**: [Comprehensive test coverage with all tests passing]
```

### [ACTION:core-functionality] Configuration Feature Template
```markdown
# [Configuration Feature Name]

> **[HIGH] CFG-[XXX]**: [Configuration feature description]
> 
> **Implementation Tokens**: `// CFG-[XXX]: [Configuration implementation]`
> 
> **Priority**: [HIGH] HIGH  
> **Status**: [[ACTION:format-processing][ACTION:migration]✅]

## 📑 Configuration Overview

### 🎯 Configuration Purpose
[Description of what this configuration controls and why it's needed]

### [ACTION:core-functionality] Configuration Parameters
- **Parameter 1**: [Type] - [Description] - Default: [value]
- **Parameter 2**: [Type] - [Description] - Default: [value]

### 📋 Configuration Schema
```yaml
# Configuration example
[parameter_name]: [default_value]  # [description]
[nested_config]:
  [sub_parameter]: [value]  # [description]
```

## 🏗️ Implementation Details

### [ACTION:core-functionality] Configuration Loading
```go
// [HIGH] CFG-[XXX]: Configuration loading implementation - [ACTION:core-functionality] Config management
type [ConfigStruct] struct {
    [Parameter1] [Type] `yaml:"[parameter1]" json:"[parameter1]"`
    [Parameter2] [Type] `yaml:"[parameter2]" json:"[parameter2]"`
}

func (c *[ConfigStruct]) Validate() error {
    // [ACTION:validation] Configuration validation logic
    [validation implementation]
}
```

### 📊 Usage Examples
```go
// [HIGH] CFG-[XXX]: Configuration usage example - [ACTION:format-processing] Implementation guide
config := &[ConfigStruct]{
    [Parameter1]: [defaultValue1],
    [Parameter2]: [defaultValue2],
}

if err := config.Validate(); err != nil {
    return fmt.Errorf("invalid configuration: %w", err)
}
```

## 🧪 Testing Requirements
- [ ] Default value validation
- [ ] Configuration loading from file
- [ ] Environment variable override testing
- [ ] Invalid configuration rejection testing
```

## 💻 Code Comment Templates

### 🚀 Standard Function Comment Template
```go
// [LOW] DOC-013: Standard function comment template - [ACTION:format-processing] Code documentation

// [PRIORITY_ICON] FEATURE-ID: [Function purpose in one clear sentence] - [ACTION_ICON] [Operation category]
//
// Purpose: [Detailed description of what this function accomplishes]
// Parameters: [Description of each parameter and its constraints]
// Returns: [Description of return values and their meaning]
// Dependencies: [List of related features, components, or external dependencies]
// Error Conditions: [Specific scenarios where this function returns errors]
// AI Assistant Notes: [Any special considerations for AI maintenance or modification]
//
// Example Usage:
//   result, err := FunctionName(param1, param2)
//   if err != nil {
//       // Handle error appropriately
//   }
func FunctionName(param1 ParamType1, param2 ParamType2) (ReturnType, error) {
    // [ACTION:discovery] Input validation with structured error responses
    if param1 == nil {
        return ReturnType{}, &StructuredError{
            Code:    "INVALID_INPUT",
            Message: "param1 cannot be nil",
            Context: map[string]interface{}{
                "function": "FunctionName",
                "param":    "param1",
            },
        }
    }
    
    // [ACTION:core-functionality] Main operation with clear operation tracking
    [implementation logic with step-by-step comments]
    
    // ✅ Success path with result validation
    return result, nil
}
```

### [ACTION:validation] Error Handling Function Template  
```go
// [LOW] DOC-013: Error handling function template - [ACTION:validation] Error management

// [PRIORITY_ICON] ERROR-[XXX]: [Error handling purpose] - [ACTION:validation] Error processing
//
// Purpose: [What error conditions this function handles]
// Error Types: [List of specific error types this function processes]
// Recovery Actions: [What recovery or cleanup actions are performed]
// Escalation: [When errors are escalated vs handled locally]
func HandleSpecificError(err error, context ErrorContext) error {
    // [ACTION:discovery] Error type classification
    switch errorType := classifyError(err); errorType {
    case DiskSpaceError:
        // [ACTION:core-functionality] Disk space error handling
        return handleDiskSpaceError(err, context)
        
    case PermissionError:
        // [ACTION:validation] Permission error handling  
        return handlePermissionError(err, context)
        
    case NetworkError:
        // [ACTION:migration] Network error with retry logic
        return handleNetworkErrorWithRetry(err, context)
        
    default:
        // [ACTION:format-processing] Unknown error logging and escalation
        logUnknownError(err, context)
        return fmt.Errorf("unhandled error type: %w", err)
    }
}
```

### [ACTION:core-functionality] Configuration Function Template
```go
// [LOW] DOC-013: Configuration function template - [ACTION:core-functionality] Configuration management

// [HIGH] CFG-[XXX]: [Configuration operation description] - [ACTION:core-functionality] Config processing
//
// Purpose: [What configuration aspect this function manages]
// Configuration Sources: [List of configuration sources: file, env, defaults]  
// Validation: [What validation is performed on configuration values]
// Side Effects: [Any side effects like file creation, environment changes]
func LoadConfigurationValue(key string, defaultValue interface{}) (interface{}, error) {
    // [ACTION:discovery] Configuration source discovery
    sources := []ConfigSource{
        NewFileConfigSource(configPath),
        NewEnvConfigSource(),
        NewDefaultConfigSource(),
    }
    
    // [ACTION:core-functionality] Configuration value resolution with precedence
    for _, source := range sources {
        if value, exists := source.GetValue(key); exists {
            // [ACTION:validation] Configuration value validation
            if err := validateConfigValue(key, value); err != nil {
                continue // Try next source
            }
            return value, nil
        }
    }
    
    // [ACTION:format-processing] Default value usage with validation
    if err := validateConfigValue(key, defaultValue); err != nil {
        return nil, fmt.Errorf("invalid default value for %s: %w", key, err)
    }
    
    return defaultValue, nil
}
```

### 🧪 Test Function Template
```go
// [LOW] DOC-013: Test function template - 🧪 Test implementation

// [LOW] TEST-[XXX]: [Test purpose description] - 🧪 Test coverage
//
// Test Coverage: [What functionality this test validates]
// Test Scenarios: [List of specific scenarios being tested]
// Test Data: [Description of test data used]
// Dependencies: [Any test dependencies or setup requirements]
func TestFunctionName(t *testing.T) {
    // [ACTION:core-functionality] Test setup with cleanup
    testDir := t.TempDir()
    defer func() {
        // 🧹 Test cleanup
        if err := os.RemoveAll(testDir); err != nil {
            t.Logf("cleanup warning: %v", err)
        }
    }()
    
    // 📋 Test cases with comprehensive scenarios
    tests := []struct {
        name        string
        input       InputType
        setupFunc   func(t *testing.T) error
        expected    ExpectedType
        wantErr     bool
        errorCheck  func(t *testing.T, err error)
    }{
        {
            name:  "Successful operation with valid input",
            input: ValidInputData,
            setupFunc: func(t *testing.T) error {
                // [ACTION:core-functionality] Test-specific setup
                return createTestData(testDir)
            },
            expected: ExpectedSuccessResult,
            wantErr:  false,
        },
        {
            name:  "Error handling with invalid input",
            input: InvalidInputData,
            setupFunc: func(t *testing.T) error {
                // No setup needed for error case
                return nil
            },
            expected: ZeroValue,
            wantErr:  true,
            errorCheck: func(t *testing.T, err error) {
                // [ACTION:validation] Specific error validation
                assert.Contains(t, err.Error(), "expected error text")
            },
        },
    }
    
    // 🧪 Test execution with proper error handling
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // [ACTION:core-functionality] Test setup
            if tt.setupFunc != nil {
                if err := tt.setupFunc(t); err != nil {
                    t.Fatalf("test setup failed: %v", err)
                }
            }
            
            // 🏃 Execute function under test
            result, err := FunctionUnderTest(tt.input)
            
            // [ACTION:validation] Validate results
            if tt.wantErr {
                assert.Error(t, err, "expected error but got none")
                if tt.errorCheck != nil {
                    tt.errorCheck(t, err)
                }
                return
            }
            
            // ✅ Validate success case
            assert.NoError(t, err, "unexpected error: %v", err)
            assert.Equal(t, tt.expected, result, "result mismatch")
        })
    }
}
```

## 📊 Technical Documentation Templates

### 🏗️ Architecture Documentation Template
```markdown
# [Component/System] Architecture

> **[PRIORITY_ICON] ARCH-[XXX]**: [Architecture component description]
> 
> **Implementation Tokens**: `// ARCH-[XXX]: [Architecture implementation]`
> 
> **Priority**: [[CRITICAL][HIGH][MEDIUM][LOW]] [Priority Level]

## 📑 Architecture Overview

### 🎯 Component Purpose
[Clear description of what this architectural component provides]

### [ACTION:core-functionality] Design Principles
- **Principle 1**: [Description and rationale]
- **Principle 2**: [Description and rationale]
- **Principle 3**: [Description and rationale]

### 🔗 System Integration
[How this component integrates with the larger system]

## 🏗️ Component Architecture

### 📋 Core Components
```
┌─────────────────────────────────────────────────────────────┐
│                [Component Name]                            │
├─────────────────────────────────────────────────────────────┤
│  Input: [Input description]                               │
│  ├─ [Input Type 1]                                        │
│  ├─ [Input Type 2]                                        │
│  └─ [Input Type 3]                                        │
├─────────────────────────────────────────────────────────────┤
│  Layer 1: [Layer description]                             │
│  ├─ [Component 1]: [Purpose]                              │
│  ├─ [Component 2]: [Purpose]                              │
│  └─ [Component 3]: [Purpose]                              │
├─────────────────────────────────────────────────────────────┤
│  Layer 2: [Layer description]                             │
│  ├─ [Component 4]: [Purpose]                              │
│  └─ [Component 5]: [Purpose]                              │
├─────────────────────────────────────────────────────────────┤
│  Output: [Output description]                             │
│  ├─ [Output Type 1]                                       │
│  └─ [Output Type 2]                                       │
└─────────────────────────────────────────────────────────────┘
```

### [ACTION:core-functionality] Data Flow
```
[Data Source] → [Processing Step 1] → [Processing Step 2] → [Data Destination]
```

### 💻 Implementation Architecture
```go
// [PRIORITY_ICON] ARCH-[XXX]: Core architecture interfaces - 🏗️ Architecture definition
type [ComponentInterface] interface {
    [Method1](param [Type]) ([ReturnType], error)
    [Method2](param [Type]) ([ReturnType], error)
}

type [ConcreteComponent] struct {
    [field1] [Type]
    [field2] [Type]
}

func (c *[ConcreteComponent]) [Method1](param [Type]) ([ReturnType], error) {
    // [PRIORITY_ICON] ARCH-[XXX]: [Method purpose] - [ACTION_ICON] [Operation]
    [implementation]
}
```

## 📊 Performance Characteristics
- **Throughput**: [Performance metrics]
- **Latency**: [Response time characteristics]  
- **Resource Usage**: [Memory, CPU, disk requirements]
- **Scalability**: [Scaling characteristics and limits]

## 🔒 Design Decisions
- **Decision 1**: [Decision description and rationale]
- **Decision 2**: [Decision description and rationale]
- **Trade-offs**: [Key trade-offs made and alternatives considered]
```

### 📋 Requirements Documentation Template
```markdown
# [Feature/System] Requirements

> **[PRIORITY_ICON] REQ-[XXX]**: [Requirements scope description]
> 
> **Implementation Tokens**: `// REQ-[XXX]: [Requirements implementation]`
> 
> **Priority**: [[CRITICAL][HIGH][MEDIUM][LOW]] [Priority Level]

## 📑 Requirements Overview

### 🎯 Scope
[Clear definition of what these requirements cover]

### 🔗 Related Features
- **[FEATURE-ID]**: [Relationship description]
- **[FEATURE-ID]**: [Relationship description]

## 📋 Functional Requirements

### [ACTION:core-functionality] Core Requirements
- **R-[REQ-ID]-001**: [Requirement statement with clear acceptance criteria]
  - **Input**: [What inputs are required]
  - **Processing**: [What processing must occur]
  - **Output**: [What outputs must be produced]
  - **Acceptance**: [How to verify requirement is met]

- **R-[REQ-ID]-002**: [Additional requirement]
  - **Input**: [Required inputs]
  - **Processing**: [Required processing]
  - **Output**: [Required outputs]
  - **Acceptance**: [Verification criteria]

### [ACTION:validation] Validation Requirements
- **R-[REQ-ID]-V01**: [Input validation requirements]
- **R-[REQ-ID]-V02**: [Output validation requirements]
- **R-[REQ-ID]-V03**: [Process validation requirements]

### [ACTION:discovery] Error Handling Requirements
- **R-[REQ-ID]-E01**: [Error detection requirements]
- **R-[REQ-ID]-E02**: [Error recovery requirements]
- **R-[REQ-ID]-E03**: [Error reporting requirements]

## 📊 Non-Functional Requirements

### 🚀 Performance Requirements
- **R-[REQ-ID]-P01**: [Performance requirement with measurable criteria]
- **R-[REQ-ID]-P02**: [Scalability requirement]
- **R-[REQ-ID]-P03**: [Resource usage requirement]

### 🔒 Security Requirements
- **R-[REQ-ID]-S01**: [Security requirement]
- **R-[REQ-ID]-S02**: [Access control requirement]
- **R-[REQ-ID]-S03**: [Data protection requirement]

## ✅ Acceptance Criteria
- [ ] [Specific, testable acceptance criterion]
- [ ] [Specific, testable acceptance criterion]
- [ ] [Specific, testable acceptance criterion]

## 🧪 Verification Methods
- **Unit Testing**: [Unit test requirements]
- **Integration Testing**: [Integration test requirements]
- **Performance Testing**: [Performance test requirements]
- **Security Testing**: [Security test requirements]
```

## [ACTION:core-functionality] AI Assistant Usage Guidelines

### 📋 Template Selection Guide

#### **[ACTION:discovery] When to Use Feature Documentation Template**
- Creating documentation for new features (FEATURE-ID pattern)
- Comprehensive feature specification needed
- Integration with multiple systems required
- Full implementation lifecycle documentation needed

#### **[ACTION:core-functionality] When to Use Configuration Template**
- Adding new configuration parameters (CFG-XXX pattern)
- Configuration validation requirements
- Environment variable override support needed
- YAML/JSON configuration schema definition

#### **🏗️ When to Use Architecture Template**  
- Defining system components (ARCH-XXX pattern)
- Component integration documentation
- Design decision documentation needed
- Performance characteristics documentation

#### **📋 When to Use Requirements Template**
- Formal requirement specification needed (REQ-XXX pattern)
- Acceptance criteria definition required
- Verification method specification needed
- Cross-feature requirement dependencies

### [ACTION:validation] Template Customization Guidelines

#### **[ACTION:core-functionality] Icon Usage Standards**
- **Priority Icons**: Must match feature-tracking.md priority levels
- **Action Icons**: Must match actual function behavior
- **Section Icons**: Use consistently across all templates
- **Cross-Reference Icons**: Include in all related feature links

#### **[ACTION:format-processing] Content Standards**
- **Clear Language**: Use simple, direct sentences for AI comprehension
- **Consistent Terminology**: Use standardized vocabulary across all templates
- **Complete Cross-References**: Include all related feature IDs and links
- **Measurable Criteria**: All requirements must be testable and verifiable

#### **🔗 Cross-Reference Requirements**
- **Feature Links**: Always link to related features in feature-tracking.md
- **Implementation Tokens**: Include in all code examples
- **Document Links**: Reference all related context documents
- **Bidirectional Links**: Ensure links work in both directions

### ✅ Template Validation Checklist

#### **📋 Pre-Use Validation**
- [ ] Template matches intended documentation type
- [ ] All placeholder values identified for replacement
- [ ] Required sections present for documentation scope
- [ ] Icon usage follows DOC-007/DOC-008 standards

#### **[ACTION:core-functionality] Content Validation**  
- [ ] All FEATURE-ID references use correct format
- [ ] Implementation tokens match feature-tracking.md entries
- [ ] Cross-references include valid links
- [ ] Priority icons match actual feature priority

#### **[ACTION:validation] Post-Creation Validation**
- [ ] DOC-008 icon validation passes
- [ ] Cross-reference integrity validated
- [ ] Implementation tokens properly formatted
- [ ] Feature-tracking.md updated with new documentation

---

**🤖 These templates provide AI assistants with consistent, standardized formats for creating high-quality documentation that follows AI-first principles and integrates seamlessly with the existing documentation ecosystem.** 