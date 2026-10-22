# persist-implementation-records — W2

**Slice:** `REQ-PSEUDOCODE_LAYER_B_CORPUS` wave **w2**

| IMPL | Action | Result |
|------|--------|--------|
| `IMPL-CONFIGURABLE_STRINGS` | `impl_detail_set_essence_pseudocode` (sidecar path) + index `metadata.last_updated` | ok |
| `IMPL-CONFIG_DISPLAY_FLATTENING` | same | ok |

**Validation:** `tied_validate_consistency` → `ok: true`; `tied_project_lint` (both IMPLs, `include_pseudocode: true`) → pseudocode sections ok.

**Sidecars (canonical):** `tied-project/implementation-decisions/IMPL-*-pseudocode.md`  
**Detail YAML:** metadata refreshed; essence remains sidecar-projected per TIED IMPL schema.

**Preload contract:** no `./tied-project/agent-preload-contract.yaml` in repo — N/A (paths unchanged in `code_locations`).

**Next:** `unit-test-red` — COVER-001/002 tests tied to STEP lines for W2 formal-without-oracle bar.
