# [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY]

## Summary contract

When the user passes a filesystem path as the first argument (not a known subcommand or global flag), classify the path and route to file backup or directory archive without requiring explicit subcommands.

INPUT: argv without program name, root CLI command tree
OUTPUT: executed backup/archive or standard Cobra dispatch
DATA: path, dry-run flag, note text, known command and flag lists

## IS_FILE

SPEC-ID: IMPL-AUTO_DETECTION::IS_FILE
STEP T001: STAT path and return true when mode is regular file

- [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: STAT path and return true only when mode is a regular file.

PROCEDURE IS_FILE(path):
  info, err = STAT(path)
  IF err != nil THEN RETURN false
  RETURN info mode is regular file

## IS_DIRECTORY

SPEC-ID: IMPL-AUTO_DETECTION::IS_DIRECTORY
STEP T001: STAT path and return true when entry is directory

- [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: STAT path and return true only when entry is a directory.

PROCEDURE IS_DIRECTORY(path):
  info, err = STAT(path)
  IF err != nil THEN RETURN false
  RETURN info is directory

## VALIDATE_PATH

SPEC-ID: IMPL-AUTO_DETECTION::VALIDATE_PATH
STEP T001: STAT path and map errors to not-exist, permission, or generic failure

- [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: STAT path and return specific errors for not-exist, permission denied, or other access failures.

PROCEDURE VALIDATE_PATH(path):
  _, err = STAT(path)
  IF err == nil THEN RETURN success
  IF not-exist THEN RETURN error "path does not exist"
  IF permission denied THEN RETURN error "permission denied"
  RETURN error "error accessing path" WITH underlying err

## HANDLE_AUTO_DETECTED_COMMAND

SPEC-ID: IMPL-AUTO_DETECTION::HANDLE_AUTO_DETECTED_COMMAND
STEP T001: REQUIRE path arg and VALIDATE_PATH
STEP T002: ROUTE file to file backup OR directory to archive OR unsupported exit

- [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: require path arg, validate accessibility, route file to file backup handler else directory archive else unsupported type exit.

PROCEDURE HANDLE_AUTO_DETECTED_COMMAND(args):
  IF len(args) == 0 THEN exit "no path provided"
  path = args[0]
  IF VALIDATE_PATH(path) fails THEN exit WITH error
  IF IS_FILE(path) THEN HANDLE_AUTO_DETECTED_FILE_BACKUP(args)
  ELSE IF IS_DIRECTORY(path) THEN HANDLE_AUTO_DETECTED_DIRECTORY_ARCHIVE(args)
  ELSE exit "unsupported file type"

## HANDLE_AUTO_DETECTED_FILE_BACKUP

SPEC-ID: IMPL-AUTO_DETECTION::HANDLE_AUTO_DETECTED_FILE_BACKUP
STEP T001: LOAD config and formatter from cwd
STEP T002: RESOLVE note from flag or second arg
STEP T003: INVOKE CreateFileBackupEnhanced with dry-run

- [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: load config from cwd, build formatter, resolve note from global flag or second arg, invoke enhanced file backup with dry-run.

PROCEDURE HANDLE_AUTO_DETECTED_FILE_BACKUP(args):
  ctx = BACKGROUND_CONTEXT
  cwd = GET_WORKING_DIRECTORY
  cfg = LOAD_CONFIG(cwd)
  formatter = NEW_OUTPUT_FORMATTER(cfg)
  filePath = args[0]
  backupNote = global note OR args[1] when present
  CREATE_FILE_BACKUP_ENHANCED(ctx, cfg, formatter, filePath, backupNote, dryRun)
  ON error EXIT WITH archive error handler code

## HANDLE_AUTO_DETECTED_DIRECTORY_ARCHIVE

SPEC-ID: IMPL-AUTO_DETECTION::HANDLE_AUTO_DETECTED_DIRECTORY_ARCHIVE
STEP T001: CHANGE_DIRECTORY to target and DEFER restore
STEP T002: LOAD config from dot and resolve note
STEP T003: INVOKE CreateFullArchiveWithContext with dry-run

- [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: chdir to target directory, load config from dot, resolve note, run full archive with context and restore cwd on exit.

PROCEDURE HANDLE_AUTO_DETECTED_DIRECTORY_ARCHIVE(args):
  ctx = BACKGROUND_CONTEXT
  dirPath = args[0]
  originalDir = GET_WORKING_DIRECTORY
  CHANGE_DIRECTORY(dirPath)
  DEFER restore originalDir
  cfg = LOAD_CONFIG(".")
  formatter = NEW_OUTPUT_FORMATTER(cfg)
  archiveNote = global note OR args[1] when present
  CREATE_FULL_ARCHIVE_WITH_CONTEXT(ctx, cfg, archiveNote, dryRun)
  ON error EXIT WITH archive error handler code

## NEW_ROOT_COMMAND

SPEC-ID: IMPL-AUTO_DETECTION::NEW_ROOT_COMMAND
STEP T001: REGISTER persistent flags and subcommands
STEP T002: root Run dispatches config, list, auto-detect, or help

- [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: build production root with persistent flags and subcommands; root Run dispatches config/list flags or auto-detect when positional args present.

PROCEDURE NEW_ROOT_COMMAND():
  root = CREATE root command WITH version template and examples
  REGISTER persistent flags dry-run, config, list, limit, debug
  ADD subcommands create, config, template, full, inc, list, diff, backup, version
  root Run:
    IF showConfig THEN handleConfigCommand; RETURN
    IF list flag changed THEN handleListFileBackups; RETURN
    IF positional args non-empty THEN HANDLE_AUTO_DETECTED_COMMAND(args); RETURN
    SHOW help
  RETURN root

## EXECUTE_WITH_AUTO_DETECTION

SPEC-ID: IMPL-AUTO_DETECTION::EXECUTE_WITH_AUTO_DETECTION
STEP T001: IF first arg is known command or flag THEN Cobra Execute
STEP T002: ELSE parse global flags and HANDLE_AUTO_DETECTED_COMMAND

- [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: bypass Cobra when first token is not known command/flag; parse global dry-run and note from argv; else SetArgs and Execute.

PROCEDURE EXECUTE_WITH_AUTO_DETECTION(rootCmd, args):
  IF len(args) == 0 THEN SetArgs and Execute root; RETURN
  firstArg = args[0]
  IF firstArg in knownCommands OR globalFlags OR starts with "-" THEN SetArgs and Execute; RETURN
  filteredArgs, dryRunFlag, noteFlag = PARSE_GLOBAL_FLAGS_FROM(args)
  SET global dryRun and note FROM parsed flags
  HANDLE_AUTO_DETECTED_COMMAND(filteredArgs)
  RETURN success

## EMBEDDED_MINITEST_PATH_TYPE_DETECTION

- [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: integration test validates validatePath, isFile, and isDirectory for file, directory, missing, and absolute paths.

PROCEDURE EMBEDDED_MINITEST_PATH_TYPE_DETECTION():
  FOR EACH fixture path EXPECT validatePath and type detection match expectedOp

## EMBEDDED_MINITEST_BACKWARD_COMPAT

- [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: explicit backup, create, and full subcommands still execute via Cobra without auto-detect path routing.

PROCEDURE EMBEDDED_MINITEST_BACKWARD_COMPAT():
  FOR EACH explicit command argv EXECUTE test root AND EXPECT success

## EMBEDDED_MINITEST_COMPOSITION_SUBCOMMANDS

- [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: production newRootCommand registers fixed subcommand set matching composition contract.

PROCEDURE EMBEDDED_MINITEST_COMPOSITION_SUBCOMMANDS():
  ASSERT sorted subcommand names equal expected list

## EMBEDDED_MINITEST_COMPOSITION_KNOWN_SUBCOMMAND

- [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: first token template routes through Cobra Execute and emits help output.

PROCEDURE EMBEDDED_MINITEST_COMPOSITION_KNOWN_SUBCOMMAND():
  executeWithAutoDetection(root, ["template", "--help"]) EXPECT help text

## EMBEDDED_MINITEST_COMPOSITION_PATH_FILE_DRYRUN

- [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: path-first argv with dry-run exercises auto-detect file backup without writing archives.

PROCEDURE EMBEDDED_MINITEST_COMPOSITION_PATH_FILE_DRYRUN():
  executeWithAutoDetection(root, ["f.txt", "--dry-run"]) EXPECT success in temp workspace
