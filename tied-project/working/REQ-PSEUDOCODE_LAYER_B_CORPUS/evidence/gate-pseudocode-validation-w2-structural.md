# gate-pseudocode-validation — W2 structural baseline

**Slice:** `REQ-PSEUDOCODE_LAYER_B_CORPUS` wave **w2** (2 IMPLs)  
**Run:** `scripts/run-pseudocode-layer-b-batch.sh --mode structural --wave w2`  
**Reports:** `pseudocode-validate-w2/`

| IMPL | ok | error count | Primary codes |
|------|----|-------------|---------------|
| `IMPL-CONFIGURABLE_STRINGS` | **true** (post-edit) | 0 | contracts + `FAILURE_MODES: none` on seed/map blocks; `HANDLE_ARCHIVE_ERROR` failure ladder |
| `IMPL-CONFIG_DISPLAY_FLATTENING` | **true** (post-edit) | 0 | contracts on four procedures; `DATA_TRANSITION` on `mergeGitSettings` |

**Disposition:** W2 structural gate **green** after sidecar contract pass (`run-pseudocode-layer-b-batch.sh --mode structural --wave w2`).

**Next:** ~~Layer C static analysis~~ (done — see `gate-pseudocode-validation-w2-layer-c.md`) → persist-implementation-records → COVER-001/002 tests → `--mode full --enforce` per W2 bar.
