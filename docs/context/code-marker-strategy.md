## Code Marker Strategy

### Recommended Tokens
- `// FEATURE-ID: Brief description` - Primary feature implementation
- `// IMMUTABLE-REF: [section]` - Links to immutable requirements
- `// TEST-REF: [test function]` - Links to test coverage
- `// DECISION-REF: DEC-XXX` - Links to implementation decisions

### Example Implementation
```go
// ARCH-001: Archive naming convention implementation
// IMMUTABLE-REF: Archive Naming Convention
// TEST-REF: TestGenerateArchiveName
// DECISION-REF: DEC-001
func GenerateArchiveName(prefix string, timestamp time.Time, gitInfo *GitInfo, note string) string {
    // Implementation here
}
``` 