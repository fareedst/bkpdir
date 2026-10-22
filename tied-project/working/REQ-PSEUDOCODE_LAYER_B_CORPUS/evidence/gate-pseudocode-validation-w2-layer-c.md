# sub-pseudocode-static-analysis-pass — W2 Layer C

**Slice:** `REQ-PSEUDOCODE_LAYER_B_CORPUS` wave **w2**  
**Tool:** `pseudocode_analyze` via `tied-cli.sh` with `gate_mode: true`  
**Checklist:** `tied-bundle/docs/pseudocode-static-analysis-checklist.yaml`

| IMPL | ok | gate_mode_applied | Report |
|------|----|-------------------|--------|
| `IMPL-CONFIGURABLE_STRINGS` | true | true | `pseudocode-analysis/IMPL-CONFIGURABLE_STRINGS.v1.json` |
| `IMPL-CONFIG_DISPLAY_FLATTENING` | true | true | `pseudocode-analysis/IMPL-CONFIG_DISPLAY_FLATTENING.v1.json` |

**Run copy:** `pseudocode-analysis-w2/*-layer-c.json` (same payload as `.v1.json` names above).

**Disposition:** PSA gate satisfied for both W2 sidecars; `input_identity.hash` stable until next sidecar edit.

**Next:** ~~persist-implementation-records~~ (done — `persist-implementation-records-w2.md`) → W2 COVER-linked RED tests → `--mode full --wave w2 --enforce`.
