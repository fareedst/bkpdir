# Governance Change Workflow

Overview: Lightweight, plan-driven governance change lifecycle to move from policy ideas to implemented decisions with traceability back to REQ tokens.

## Phases
- Proposal: A governance change proposal is drafted, linked to a REQ token, and aligned with cross-links to ARCH/IMPL decisions.
- Review: Stakeholders review the proposal; ensure alignment with existing tokens and architecture decisions.
- Approval: Obtain necessary approvals to implement changes in architecture or implementation documents.
- Implementation: Update ARCH/IMPL decisions and code comments to reflect changes; update REQ tokens if needed.
- Validation: Run token validation and cross-link checks; update docs and tests as needed.
- Deployment/Migration: Record migration actions and status in the governance pages and tasks.

## Artifacts
- Cross-links to REQ tokens in `docs/governance/cross-links.md`.
- Architecture decisions in `stdd/architecture-decisions.md`.
- Implementation decisions in `stdd/implementation-decisions.md`.
- Token registry updates in `stdd/semantic-tokens.md`.

## Cadence
- Quarterly governance reviews; escalate any high-risk changes earlier as needed.

