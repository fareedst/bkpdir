# Traceable commit — REQ-IMMUTABLE_ARCHIVE_NAMING / -seqN LEAP

**Tokens:** [REQ-IMMUTABLE_ARCHIVE_NAMING], [REQ-INCREMENTAL_ANCHOR_MODE], [ARCH-INCREMENTAL_ANCHOR_MODE], [IMPL-ZIP_FORMAT], [IMPL-INCREMENTAL_ANCHOR_MODE]

**Scope:** Document optional `-seqN` incremental path disambiguator in immutable naming REQ; index + semantic-token summaries; ARCH/IMPL traceability; vocab grammar and user spec; code annotation on `resolveIncrementalArchivePath`.

**Verification:** No behavior change. `tied_project_lint` on touched TIED paths; `tied_validate_consistency` ok.

**Product:** `-seqN` before `.zip` when generated incremental path collides; not a user note; `N` ≥ 2 until unused.
