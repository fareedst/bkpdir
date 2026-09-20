# [IMPL-GIT_DIRTY_CONFIG] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]

## Summary contract

Gate `-dirty` suffix in archive filenames on ShowGitDirtyStatus config, propagate the setting through ArchiveConfig and config merge helpers.

INPUT: ArchiveConfig git cleanliness, cfg.ShowGitDirtyStatus
OUTPUT: archive basename with optional -dirty segment
DATA: GitBranch, GitHash, GitIsClean, legacy and nested Git.ShowDirtyStatus

## GENERATE_FULL_ARCHIVE_NAME_FROM_CONFIG

SPEC-ID: IMPL-GIT_DIRTY_CONFIG::GENERATE_FULL_ARCHIVE_NAME_FROM_CONFIG
STEP T001: APPEND -dirty WHEN repo dirty AND ShowGitDirtyStatus enabled

- [IMPL-GIT_DIRTY_CONFIG] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: append -dirty to full archive name only when repo is dirty and ShowGitDirtyStatus is enabled.

PROCEDURE GENERATE_FULL_ARCHIVE_NAME_FROM_CONFIG(cfg):
  name = prefix + timestamp + optional =branch=hash
  IF git metadata present AND NOT GitIsClean AND ShowGitDirtyStatus: name += "-dirty"
  APPEND optional note; RETURN name + ".zip"

## GENERATE_INCREMENTAL_ARCHIVE_NAME

SPEC-ID: IMPL-GIT_DIRTY_CONFIG::GENERATE_INCREMENTAL_ARCHIVE_NAME
STEP T001: APPLY conditional -dirty suffix on incremental archive names

- [IMPL-GIT_DIRTY_CONFIG] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: apply same conditional -dirty suffix on incremental update archive names.

PROCEDURE GENERATE_INCREMENTAL_ARCHIVE_NAME(cfg):
  name = base + "_update=" + timestamp + optional git segments
  IF dirty AND ShowGitDirtyStatus: name += "-dirty"
  RETURN name + ".zip"

## GENERATE_FULL_ARCHIVE_NAME

SPEC-ID: IMPL-GIT_DIRTY_CONFIG::GENERATE_FULL_ARCHIVE_NAME
STEP T001: COPY ShowGitDirtyStatus AND populate git metadata before name build

- [IMPL-GIT_DIRTY_CONFIG] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: copy cfg.ShowGitDirtyStatus into ArchiveConfig before delegating to name builder with git info when enabled.

PROCEDURE GENERATE_FULL_ARCHIVE_NAME(cfg, cwd, note):
  archiveConfig.ShowGitDirtyStatus = cfg.ShowGitDirtyStatus
  IF IncludeGitInfo: populate branch, hash, isClean from repository
  RETURN GenerateArchiveName(archiveConfig)

## MERGE_SHOW_GIT_DIRTY_STATUS

SPEC-ID: IMPL-GIT_DIRTY_CONFIG::MERGE_SHOW_GIT_DIRTY_STATUS
STEP T001: MERGE show_git_dirty_status WITH CFG-001 explicit-set precedence

- [IMPL-GIT_DIRTY_CONFIG] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: merge show_git_dirty_status in mergeBasicSettings and mergeGitSettings respecting CFG-001 explicit-set precedence.

PROCEDURE MERGE_SHOW_GIT_DIRTY_STATUS(dst, src, inheritContext, explicitlySetFields):
  IN inheritance: set when explicitly in source or differs from default unless earlier file set it
  IN sequential: set only when earlier did not set and source explicitly provides non-default value
  SYNC dst.Git.ShowDirtyStatus with legacy ShowGitDirtyStatus
