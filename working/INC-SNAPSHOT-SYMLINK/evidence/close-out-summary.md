# Close-out evidence: incremental snapshot symlink hashing

**Change request:** ad-hoc fix (no separate CITDP); traceability via `[IMPL-DIRECTORY_COMPARISON]`.

## Problem

`bkpdir inc` failed with `read …/node_modules/@tied/agentstream: is a directory` when diffing against a full archive.

## Fix

`pkg/fileops/comparison.go`: hash symlink targets via `Readlink` (aligned with zip symlink entries), regression test `TestCreateDirectorySnapshot_SymlinkToDirectory_REQ_DIFF_COMMAND`.

## Verification

- `cd pkg/fileops && go test ./...` — PASS
- Root: `go test -run 'IncCommand|CreateDirectorySnapshot|CalculateDiff'` — PASS
- Manual: `bkpdir inc` in `stdd` repo — incremental archive created (2026-09-24)

## TIED

- `tied/implementation-decisions/IMPL-DIRECTORY_COMPARISON-pseudocode.md` — `FILEOPS_CREATE_DIRECTORY_SNAPSHOT`
- `tied/vocab/bkpdir-archives-backup-diff.md` — symlink snapshot hash row
- `tied_validate_consistency` — ok (2026-09-24 close-out pass)
