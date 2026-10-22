# Adversarial inquiry waiver — incremental anchor mode

**Request:** `REQ-INCREMENTAL_ANCHOR_MODE`  
**CITDP:** `CITDP-INCREMENTAL-ANCHOR-MODE.yaml`  
**depth_tier:** `minimal`  
**gate_policy:** `advisory`

## Rationale

Local filesystem archive semantics (default full-anchored incrementals, chained opt-in, same-minute
path disambiguation). No new network or auth surface. Counterexamples and falsification questions
are recorded in the CITDP `risk_analysis.adversarial_inquiry` block.
