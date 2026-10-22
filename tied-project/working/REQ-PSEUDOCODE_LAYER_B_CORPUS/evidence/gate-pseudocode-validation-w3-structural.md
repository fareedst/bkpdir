# W3 pseudo-code gate — structural Layer B (2026-10-21)

**Wave:** `w3` — 27 Tier C `lead_check_skip` IMPLs per `tied-project/spec/layer-b-wave-registry.yaml`  
**Run:** `scripts/run-pseudocode-layer-b-batch.sh --mode structural --wave w3`

```bash
TIED_BASE_PATH="$PWD/tied-project" scripts/run-pseudocode-layer-b-batch.sh \
  --mode structural --wave w3 \
  --out-dir tied-project/working/REQ-PSEUDOCODE_LAYER_B_CORPUS/pseudocode-validate-w3-structural
```

**Result:** **27/27 ok** after contract hydration (`scripts/w3-hydrate-sidecar-contracts.py` one-shot helper under working folder).

**Full enforce (verification-gate bar):**

```bash
TIED_BASE_PATH="$PWD/tied-project" scripts/run-pseudocode-layer-b-batch.sh \
  --mode full --wave w3 --enforce \
  --out-dir tied-project/working/REQ-PSEUDOCODE_LAYER_B_CORPUS/pseudocode-validate-w3-full
```

**Result:** **27/27 ok** — COVER waived via `lead_check_skip`; structural SHAPE + meaningful How + PRE/POST/EFFECTS on Active PROCEDURE blocks.

**Layer C:** not in scope for W3 (doc/process Tier C); no PSA gate_mode receipts required this wave.
