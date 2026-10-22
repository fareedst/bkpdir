# Verification — REQ-IMMUTABLE_FILE_EXCLUSION (multi-segment any-depth excludes)

- `go test ./...` — pass (2026-10-21)
- `TestMultiSegmentDirPatternAnyDepth_REQ_IMMUTABLE_FILE_EXCLUSION` — pass
- `tied_validate_consistency` — index valid
- Manual stdd backup spot-check — deferred to sponsor (requires stdd tree + full archive run)
