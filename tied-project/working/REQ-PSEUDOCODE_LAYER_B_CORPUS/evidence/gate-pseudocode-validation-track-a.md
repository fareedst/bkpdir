# gate-pseudocode-validation — Track A bootstrap (REQ-PSEUDOCODE_LAYER_B_CORPUS)

**Slice:** `REQ-PSEUDOCODE_LAYER_B_CORPUS` (continuation of `[REQ-PSEUDOCODE_FORMAL_VERIFICATION]`)

**Disposition:** Track A only — no IMPL sidecar edits in this pass.

| Layer | Action | Result |
|-------|--------|--------|
| Layer A | Deferred until sidecar wave (W0) | N/A — corpus unchanged |
| Layer B structural | Deferred to Phase 0 audit + batch driver on pilot set | See `corpus-layer-b-gap-report.json` |
| Layer C `gate_mode` | Unchanged — existing W0 token set in `run_layer_c_gate.sh` | N/A at bootstrap |

**Rationale:** Hybrid plan allows tooling RED before `pre_implementation` gate; pseudo-code validation gate for edited sidecars runs at W0 TIED pass.

**Pre-contract-grammar:** Legacy Active procedures remain untouched until wave edits.
