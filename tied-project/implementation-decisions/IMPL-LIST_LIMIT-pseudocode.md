# [IMPL-LIST_LIMIT] [ARCH-LIST_LIMIT] [REQ-LIST_LIMIT]

## Summary contract

Limits archive and file-backup listing output to the N most recent items after sorting, with 0 meaning show all.

INPUT: listLimit flag value, cfg, formatter, optional file path
OUTPUT: printed list lines truncated to limit
DATA: archives or backups slice, sort order (most recent first)

## REGISTER_LIMIT_FLAG

SPEC-ID: IMPL-LIST_LIMIT::REGISTER_LIMIT_FLAG
STEP T001: Register --limit/-n persistent flag with default 10 so list subcommand and --list share one limit variable

- [IMPL-LIST_LIMIT] [ARCH-LIST_LIMIT] [REQ-LIST_LIMIT] — How: register --limit/-n persistent flag with default 10 so list subcommand and --list share one limit variable.

PROCEDURE REGISTER_LIMIT_FLAG():
  rootCmd.PersistentFlags().IntVarP(&listLimit, "limit", "n", 10, "Limit the number of items to display (0 = show all)")

## HANDLE_LIST_COMMAND

SPEC-ID: IMPL-LIST_LIMIT::HANDLE_LIST_COMMAND
STEP T001: Load config and formatter then pass listLimit to ListArchivesEnhanced for list subcommand

- [IMPL-LIST_LIMIT] [ARCH-LIST_LIMIT] [REQ-LIST_LIMIT] — How: load config and formatter then pass listLimit to ListArchivesEnhanced for list subcommand.

PROCEDURE HANDLE_LIST_COMMAND():
  cwd = GET_WORKING_DIRECTORY()
  cfg = LOAD_CONFIG(cwd)
  formatter = NEW_OUTPUT_FORMATTER(cfg)
  CALL ListArchivesEnhanced(cfg, formatter, listLimit)

## LIST_ARCHIVES_ENHANCED

SPEC-ID: IMPL-LIST_LIMIT::LIST_ARCHIVES_ENHANCED
STEP T001: Sort archives most-recent-first then truncate to limit when limit > 0 before formatted output

- [IMPL-LIST_LIMIT] [ARCH-LIST_LIMIT] [REQ-LIST_LIMIT] [REQ-OUTPUT_FORMATTING] — How: sort archives most-recent-first then truncate to limit when limit > 0 before formatted output.

PROCEDURE LIST_ARCHIVES_ENHANCED(cfg, formatter, limit):
  archiveDir = RESOLVE_ARCHIVE_DIR(cfg)
  archives = LIST_ARCHIVES(archiveDir)
  IF archives is empty THEN
    PRINT no archives message
    RETURN
  END IF
  SORT archives BY CreationTime descending
  IF limit > 0 AND len(archives) > limit THEN
    archives = archives[:limit]
  END IF
  FOR EACH archive IN archives:
    PRINT formatter.FormatListArchiveWithExtraction(archive.Path, creationTime)

## HANDLE_LIST_FILE_BACKUPS_COMMAND

SPEC-ID: IMPL-LIST_LIMIT::HANDLE_LIST_FILE_BACKUPS_COMMAND
STEP T001: Resolve file path from --list or args and pass listLimit to ListFileBackupsEnhanced

- [IMPL-LIST_LIMIT] [ARCH-LIST_LIMIT] [REQ-LIST_LIMIT] — How: resolve file path from --list or args and pass listLimit to ListFileBackupsEnhanced.

PROCEDURE HANDLE_LIST_FILE_BACKUPS_COMMAND(args):
  filePath = listFile OR args[0]
  cfg = LOAD_CONFIG(cwd)
  formatter = NEW_OUTPUT_FORMATTER(cfg)
  CALL ListFileBackupsEnhanced(cfg, formatter, filePath, listLimit)

## LIST_FILE_BACKUPS_ENHANCED

SPEC-ID: IMPL-LIST_LIMIT::LIST_FILE_BACKUPS_ENHANCED
STEP T001: Sort backups most-recent-first then truncate to limit when limit > 0 before formatted output

- [IMPL-LIST_LIMIT] [ARCH-LIST_LIMIT] [REQ-LIST_LIMIT] [REQ-OUTPUT_FORMATTING] — How: sort backups most-recent-first then truncate to limit when limit > 0 before formatted output.

PROCEDURE LIST_FILE_BACKUPS_ENHANCED(cfg, formatter, filePath, limit):
  backupDir = RESOLVE_BACKUP_DIR(cfg, filePath)
  backups = LIST_FILE_BACKUPS(backupDir, baseFilename)
  IF backups is empty THEN
    PRINT no backups message
    RETURN
  END IF
  IF limit > 0 AND len(backups) > limit THEN
    backups = backups[:limit]
  END IF
  FOR EACH backup IN backups:
    PRINT formatted backup line

## COMMAND_HANDLER_LIST

SPEC-ID: IMPL-LIST_LIMIT::COMMAND_HANDLER_LIST
STEP T001: Apply hardcoded default limit of 10 for CommandHandler archive and file-backup listing paths

- [IMPL-LIST_LIMIT] [ARCH-LIST_LIMIT] [REQ-LIST_LIMIT] — How: apply hardcoded default limit of 10 for CommandHandler archive and file-backup listing paths.

PROCEDURE COMMAND_HANDLER_LIST():
  HandleListArchives: sort archives, limit := 10, truncate, print
  HandleListFileBackups: CALL ListFileBackupsEnhanced(..., limit=10)
