# Traceable commit evidence — appearance profile close-out

**REQ:** REQ-BKPDIR_CLI_OUTPUT_STANDARD  
**Run ID:** cli-output-profile-close-20261022

## Gates

- Verification runner: `merged_decision.allowed: true`
- Close-out runner: `evidence/close-out-gates-run-profile.json` — blocking_gap_count 0
- Pseudocode Layer B: `IMPL-BKPDIR_CLI_OUTPUT_CONFORMANCE` ok (gate_mode)

## Hygiene

- `evidence/git-hygiene-receipt-2026-10-22-profile.md` (preview only, apply false)

## Tests

- `evidence/go-test-full-2026-10-22-profile.txt`
- `evidence/go-test-profile-2026-10-22.txt`
- Manifest: `verification-evidence-manifest.v1.json`

## Vocab

- bkpdir-cli.md profile paragraph pre-existing; no new semantic tokens.
