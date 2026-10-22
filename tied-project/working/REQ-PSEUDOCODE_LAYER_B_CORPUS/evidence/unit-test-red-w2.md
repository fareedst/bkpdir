# unit-test-red — W2 COVER-001/002

**Slice:** `REQ-PSEUDOCODE_LAYER_B_CORPUS` wave **w2**  
**Pattern:** Brownfield baseline lock — tests assert POST/FAILURE_MODES from sidecar STEP blocks; first run **green**.

## Tests added/aligned

| SPEC block | STEP | Test | COVER |
|------------|------|------|-------|
| `DEFAULT_CONFIG_ERROR_FORMATS` | T001–T002 | `TestIMPL_CONFIGURABLE_STRINGS_DEFAULT_CONFIG_ERROR_FORMATS_COVER001` | 001 |
| `GET_ERROR_FORMAT_STRINGS` | T001 | `TestIMPL_CONFIGURABLE_STRINGS_GET_ERROR_FORMAT_STRINGS_COVER001` | 001 |
| `FORMAT_DISK_FULL_ERROR` | T001 | `TestIMPL_CONFIGURABLE_STRINGS_FORMAT_DISK_FULL_ERROR_COVER001` | 001 |
| `HANDLE_ARCHIVE_ERROR` | T001–T002 | `TestHandleArchiveError` (subtests `COVER001`/`COVER002` + FAILURE_MODES names) | 001 + 002 |
| `REFLECTCONFIGFIELDS_FLAT_DISPLAY` | T001 | `TestIMPL_CONFIG_DISPLAY_FLATTENING_reflectConfigFields_COVER001` | 001 |
| `GETALLCONFIGVALUESWITHSOURCES` | T001 | `TestIMPL_CONFIG_DISPLAY_FLATTENING_GetAllConfigValuesWithSources_COVER001` | 001 |
| `DETERMINEFIELDCATEGORY` | T001 | `TestIMPL_CONFIG_DISPLAY_FLATTENING_determineFieldCategory_COVER001` | 001 |
| `MERGEGITSETTINGS_LEGACY_FLAT_KEYS` | T001 | `TestIMPL_CONFIG_DISPLAY_FLATTENING_mergeGitSettings_COVER001` | 001 |

**Files:** `config_pseudocode_w2_test.go`, `errors_test.go`

**Command:** `go test -count=1 -run 'TestIMPL_CONFIG|TestHandleArchiveError' .` → pass

**Next:** `unit-test-green` (no production delta expected) → `verification-gate` with `--mode full --wave w2 --enforce`.
