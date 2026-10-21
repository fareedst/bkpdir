# Close-out envelope pass — REQ-PSEUDOCODE_LAYER_B_CORPUS

**Run ID:** `layer-b-corpus-20261021-close`  
**Date:** 2026-10-21

## Commands

```bash
node /path/to/stdd/tools/bootstrap/templates/run-close-out-gates.mjs \
  --project-root /Users/fareed/Documents/dev/go/bkpdir \
  --request-token REQ-PSEUDOCODE_LAYER_B_CORPUS \
  --tracker-path tied-project/working/REQ-PSEUDOCODE_LAYER_B_CORPUS/agent-req-implementation-checklist.yaml \
  --citdp-path tied-project/working/REQ-PSEUDOCODE_LAYER_B_CORPUS/CITDP-REQ-PSEUDOCODE_LAYER_B_CORPUS.yaml \
  --phase close_out \
  --run-id layer-b-corpus-20261021-close \
  --envelope-blocking \
  --reconcile
```

```bash
TIED_BASE_PATH=tied-project .cursor/skills/tied-yaml/scripts/tied-cli.sh request_evidence_envelope_validate \
  '{"envelope": <file>, "fail_on_error_gaps": true, "fail_on_process_gaps": true}'
```

## Result

| Check | Outcome |
|-------|---------|
| Envelope rebuild | `tied-project/working/REQ-PSEUDOCODE_LAYER_B_CORPUS/evidence/request-evidence-envelope.v1.json` |
| Blocking error gaps | **0** |
| Advisory warn gaps | **6** (`finding_unresolved` / `warn_not_success` per inquiry phase — expected under `gate_policy: advisory`) |
| Evidence chain profile | Present at `evidence/evidence-chain-profile.v1.json` |
| Adherence reconcile | `process_grade` band **A** (score 100) |

## Not in scope (full close-out still incomplete)

- `tied_checklist_gate_validate` **close_out** `allowed: false` — `traceable-commit` pending, tracker sparse, close_out inquiry findings unresolved (advisory).
- **Commit skipped** per unified close-out policy.
