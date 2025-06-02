# Cross-Reference Template for Enhanced Documentation

## Semantic Linking Strategy

### Feature Reference Format
Every feature mention should include:
```markdown
**Feature**: [ARCH-001: Archive Naming](feature-tracking.md#ARCH-001) 
**Spec**: [Archive Naming Convention](specification.md#archive-naming-convention)
**Requirements**: [Archive Naming](requirements.md#archive-naming-requirements)  
**Architecture**: [NamingService](architecture.md#naming-service)
**Tests**: [TestGenerateArchiveName](testing.md#testgeneratearchivename)
**Code**: `// ARCH-001: Archive naming` tokens
```

## Example Implementation

### ARCH-001: Archive Naming Convention
**Feature**: [ARCH-001: Archive Naming](feature-tracking.md#ARCH-001) 
**Spec**: [Archive Naming Convention](specification.md#archive-naming-convention)
**Requirements**: [Archive Naming](requirements.md#archive-naming-requirements)  
**Architecture**: [NamingService](architecture.md#naming-service)
**Tests**: [TestGenerateArchiveName](testing.md#testgeneratearchivename)
**Code**: `// ARCH-001: Archive naming` tokens

This feature implements the core archive naming convention that generates unique, sortable archive filenames with timestamps, Git information, and optional notes.

### Change Impact Tracking
```markdown
## Change Impact Matrix
| Change Type | Documents to Update | Validation Required |
|-------------|-------------------|-------------------|
| New Feature | spec, requirements, architecture, testing, feature-tracking | Full validation |
| Requirement Change | requirements, architecture, testing | Implementation validation |
| Implementation Detail | architecture, testing | Test validation |
| Bug Fix | requirements, testing | Regression validation |
```

### Bi-directional Linking
Each document should contain:
- **Forward Links**: "This requirement is implemented by [Architecture Component](architecture.md#component)"
- **Backward Links**: "This component implements [Requirement REQ-001](requirements.md#REQ-001)"
- **Sibling Links**: "Related to [Feature FILE-001](feature-tracking.md#FILE-001)" 

## Additional Cross-Reference Example

### TEST-INFRA-001-A: Archive Corruption Testing Framework
**Feature**: [TEST-INFRA-001-A: Archive Corruption Testing Framework](feature-tracking.md#TEST-INFRA-001-A)
**Spec**: [Testing Infrastructure - Archive Corruption Testing Framework](specification.md#archive-corruption-testing-framework)  
**Requirements**: [Testing Infrastructure Requirements - Archive Corruption Testing Framework](requirements.md#archive-corruption-testing-framework)
**Architecture**: [Testing Infrastructure Architecture - Archive Corruption Testing Framework](architecture.md#archive-corruption-testing-framework)
**Tests**: [TestCorruptionType, TestCreateTestArchive, TestCorruptionCRC, TestCorruptionReproducibility, BenchmarkCorruptionCRC](testing.md#archive-corruption-testing-framework)
**Code**: `// TEST-INFRA-001-A: Archive corruption testing framework` tokens

This example demonstrates a completed feature with comprehensive cross-references across all documentation layers. The Archive Corruption Testing Framework shows how features can have rich documentation coverage including:

- **Feature tracking** with completion status and implementation notes
- **User-facing specification** with behavior description and examples
- **Technical requirements** with design decisions and implementation constraints  
- **System architecture** with data models and integration patterns
- **Comprehensive testing** with detailed test coverage and performance benchmarks
- **Implementation tokens** linking documentation to actual code

The cross-reference density (8+ references across documentation layers) demonstrates the value of rich linking for LLM consumption and change impact analysis. 