# REQ-PSEUDOCODE_FORMAL_VERIFICATION — evidence layout

**Active (envelope scan root):** top-level files only. Subdirectories are not indexed by `request_evidence_envelope_build`.

| File | Role |
|------|------|
| `request-evidence-envelope.v1.json` | Canonical request envelope (integrated close-out) |
| `verification-evidence-manifest.v1.json` | Scan-root verification manifest (`run_id: scan-root-manifest-20260920`) |
| `evidence-chain-profile.v1.json` | Layer C profile snapshot |
| `activation-pre_implementation.json` | Inquiry activation metadata |
| `*-evidence.md` | Checklist slug stubs referenced by `gates/ledger.jsonl` (hygiene close-out sync) |

**Scan-root slice:** command captures and copies live under `../scan-root/evidence/` (including `manifest-run/`).

**Archived (2026-09-20):** lead-hygiene discovery logs, bulk audit, old manifest captures, duplicate ledger → `archive-lead-hygiene-20260919/`. See `../CLOSE-OUT-HANDOFF.md` for hygiene proof boundary; paths in that doc use the archive prefix.
