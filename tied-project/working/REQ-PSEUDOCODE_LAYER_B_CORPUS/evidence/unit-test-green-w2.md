# unit-test-green — W2 brownfield baseline

**Slice:** `REQ-PSEUDOCODE_LAYER_B_CORPUS` wave **w2**  
**Expectation:** No production code changes; tests already green from RED lock.

## Commands

| Command | Result |
|---------|--------|
| `go test -count=1 -run 'TestIMPL_CONFIG\|TestHandleArchiveError' .` | pass |
| `SKIP_TIED_MCP=1 scripts/run-spec-verification.sh` | pass (L0 fleet) |
| `tied_validate_consistency` | pass |

**Files:** `config_pseudocode_w2_test.go`, `errors_test.go`

**Next:** `verification-gate` — manifest `verification-evidence-manifest-w2.v1.json`
