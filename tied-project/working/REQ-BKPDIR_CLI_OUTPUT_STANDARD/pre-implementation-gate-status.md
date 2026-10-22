# Pre-implementation gate — [REQ-BKPDIR_CLI_OUTPUT_STANDARD]

**Phase:** `pre_implementation` (after refine-plan, before RED tests)  
**Depth:** `minimal` · **Gate policy:** `advisory`

## refine-plan gate snapshot (2026-10-22)

Initial `tied_checklist_gate_validate` result: **`allowed: false`** (expected until build-plan completes TIED authoring and Tracker dispositions).

Diagnostics to clear before RED tests:

| Diagnostic | Unblock action |
| --- | --- |
| `missing_required_step:sub-adversarial-inquiry-pass` | Mark `sub-adversarial-inquiry-pass` **not_applicable** with `policy: minimal-depth-no-inquiry` on Tracker (see PLAN-BKPDIR-TWO-FOLDER-CLEANUP pattern), or run structural pass if depth rises |
| Minimal adversarial CITDP fields | Satisfied in `CITDP-REQ-BKPDIR_CLI_OUTPUT_STANDARD.yaml` under `risk_analysis.adversarial_inquiry` |
| Pending Tracker steps through `test-strategy` | Complete during build-plan: author REQ/ARCH/IMPL + sidecar, `gate-pseudocode-validation`, disposition planning steps |
| Stale gate after TIED edits | Re-run gate with `receipt_persistence` into `tied-project/working/REQ-BKPDIR_CLI_OUTPUT_STANDARD/gates/` |

## build-plan entry checklist

1. Author `[REQ-BKPDIR_CLI_OUTPUT_STANDARD]`, `[ARCH-BKPDIR_CLI_OUTPUT_STANDARD]`, `[IMPL-BKPDIR_CLI_OUTPUT_CONFORMANCE]` + token-commented sidecar via tied-cli/MCP.
2. Register tokens in `tied-project/semantic-tokens.yaml`; run `pseudocode_validate`.
3. Update Tracker dispositions for `session-bootstrap` … `test-strategy` and `sub-adversarial-inquiry-pass` (N/A at minimal).
4. Re-run `tied_checklist_gate_validate` with `allowed: true` before `pkg/clioutput` RED tests.
