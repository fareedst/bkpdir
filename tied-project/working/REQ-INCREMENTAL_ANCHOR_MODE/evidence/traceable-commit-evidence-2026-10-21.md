# Traceable commit — REQ-INCREMENTAL_ANCHOR_MODE

**Tokens:** REQ-INCREMENTAL_ANCHOR_MODE (new), LEAP on REQ-INCREMENTAL_DUPLICATE_PREVENTION, REQ-DIFF_COMMAND, ARCH/IMPL peers.

**Verification:** `go test -run 'TestIncremental|TestReconstruct|TestIncCommand|TestCalculateDiff|IncrementalAnchor|Disambiguation' .` (log under `evidence/verification-go-test-2026-10-21.txt`).

**Gates:** `pre_implementation` + `verification` receipts under `gates/`; envelope `request-evidence-envelope.v1.json` with `fail_on_error_gaps: true`.

**Product:** Default `incremental_anchor_mode: full`; chained opt-in; `-seqN` path disambiguation; user spec + vocab.
