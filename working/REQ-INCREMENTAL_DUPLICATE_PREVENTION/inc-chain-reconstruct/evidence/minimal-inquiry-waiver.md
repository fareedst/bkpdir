# Adversarial inquiry waiver — incremental chain reconstruction

**Request:** `REQ-INCREMENTAL_DUPLICATE_PREVENTION` (slice `INC-CHAIN-RECONSTRUCT`)  
**CITDP:** `CITDP-INC-CHAIN-RECONSTRUCT.yaml`  
**depth_tier:** `minimal`  
**gate_policy:** `advisory`  

## Rationale

Bug fix in archive state reconstruction and LEAP alignment. No new surface area beyond existing
`bkpdir inc` / `bkpdir diff` comparison paths. Counterexamples and falsification questions are
recorded in the CITDP `risk_analysis.adversarial_inquiry` block.
