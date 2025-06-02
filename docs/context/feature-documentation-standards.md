# Feature Documentation Standards

## Feature Tracking Format

Each feature must be documented across all relevant layers with the following structure:

- **Immutable ID**: A unique identifier following the pattern `PREFIX-###` (e.g., `ARCH-001`, `CMD-001`, `CFG-001`)
- **Specification Reference**: Link to the relevant section in the specification document that describes what the feature does
- **Requirements Reference**: Link to the specific requirements that this feature implements
- **Architecture Reference**: Link to the architectural components that implement this feature
- **Test Reference**: Link to the test coverage for this feature
- **Implementation Tokens**: Code markers that link the implementation back to this feature

## Implementation Token Format

Implementation tokens should be placed in the code as comments with the following format:

```go
// FEATURE-ID: Brief description of the implementation
```

Example:
```go
// ARCH-001: Archive naming convention implementation
func generateArchiveName() string {
    // ...
}
```

## Cross-Referencing

All features should be cross-referenced across these documents:
1. `specification.md` - What the feature does
2. `requirements.md` - Why the feature exists
3. `architecture.md` - How the feature is implemented
4. `testing.md` - How the feature is verified
5. `feature-tracking.md` - The master registry of all features

## Versioning

- When modifying an existing feature, append a modification suffix (e.g., `ARCH-001-MOD-001`)
- Document all modifications in the relevant files
- Keep the original feature ID in the implementation tokens for traceability
