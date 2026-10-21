# [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS]

## Summary contract

Builder-pattern CLI framework: command and root builders, flag/context/dry-run/version managers, and a CLIApp container that wires managers and runs the root command with optional signal-aware context.

INPUT: AppInfo, command templates, CommandContext, DryRunOperation
OUTPUT: configured root command, exit status, dry-run log lines
DATA: FlagManager, VersionManager, ContextManager, DryRunManager, BuildInfo

## NEWCOMMANDBUILDER

SPEC-ID: IMPL-CLI_FRAMEWORK::NEWCOMMANDBUILDER
STEP T001: CONSTRUCT DefaultCommandBuilder WITH non-nil FlagManager defaulting when nil

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: construct DefaultCommandBuilder with a non-nil FlagManager (defaulting when nil).

PROCEDURE NEWCOMMANDBUILDER(flagMgr):
  IF flagMgr == nil THEN flagMgr = NEWFLAGMANAGER()
  RETURN DefaultCommandBuilder{flagManager: flagMgr}

## DEFAULTCOMMANDBUILDER_NEWCOMMAND

SPEC-ID: IMPL-CLI_FRAMEWORK::DEFAULTCOMMANDBUILDER_NEWCOMMAND
STEP T001: CREATE subcommand WITH use, short, AND long description fields

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: create a subcommand with use, short, and long description fields.

PROCEDURE DEFAULTCOMMANDBUILDER_NEWCOMMAND(name, short, long):
  DECLARE cmd as empty command object
  SET cmd.Use = name
  SET cmd.Short = short
  SET cmd.Long = long
  RETURN cmd

## DEFAULTCOMMANDBUILDER_WITHHANDLER

SPEC-ID: IMPL-CLI_FRAMEWORK::DEFAULTCOMMANDBUILDER_WITHHANDLER
STEP T001: ATTACH RunE handler TO command AND return command for fluent chaining

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: attach RunE handler to command and return command for fluent chaining.

PROCEDURE DEFAULTCOMMANDBUILDER_WITHHANDLER(cmd, handler):
  SET cmd.RunE = handler
  RETURN cmd

## DEFAULTCOMMANDBUILDER_WITHSUBCOMMANDS

SPEC-ID: IMPL-CLI_FRAMEWORK::DEFAULTCOMMANDBUILDER_WITHSUBCOMMANDS
STEP T001: REGISTER each child command ON parent AND return parent

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: register each child command on parent and return parent.

PROCEDURE DEFAULTCOMMANDBUILDER_WITHSUBCOMMANDS(parent, children):
  FOR EACH child IN children: parent.AddCommand(child)
  RETURN parent

## NEWROOTCOMMANDBUILDER

SPEC-ID: IMPL-CLI_FRAMEWORK::NEWROOTCOMMANDBUILDER
STEP T001: CONSTRUCT DefaultRootCommandBuilder WITH default flag AND version managers when nil

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: construct DefaultRootCommandBuilder with default flag and version managers when nil.

PROCEDURE NEWROOTCOMMANDBUILDER(flagMgr, versionMgr):
  IF flagMgr == nil THEN flagMgr = NEWFLAGMANAGER()
  IF versionMgr == nil THEN versionMgr = NEWVERSIONMANAGER()
  RETURN DefaultRootCommandBuilder{flagManager: flagMgr, versionManager: versionMgr}

## DEFAULTROOTCOMMANDBUILDER_NEWROOTCOMMAND

SPEC-ID: IMPL-CLI_FRAMEWORK::DEFAULTROOTCOMMANDBUILDER_NEWROOTCOMMAND
STEP T001: BUILD root command FROM AppInfo WITH formatted version AND help default

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: build root command from AppInfo with formatted version, version template, and help when invoked without subcommand.

PROCEDURE DEFAULTROOTCOMMANDBUILDER_NEWROOTCOMMAND(info):
  cmd = command with Use=info.Name, Short=info.Short, Long=info.Long
  SET cmd.Version = versionManager.FormatVersion(info.Build)
  SET cmd version template from versionManager.CreateVersionTemplate(info.Build)
  SET default Run to show help when no subcommand
  RETURN cmd

## DEFAULTROOTCOMMANDBUILDER_BUILDCOMMAND

SPEC-ID: IMPL-CLI_FRAMEWORK::DEFAULTROOTCOMMANDBUILDER_BUILDCOMMAND
STEP T001: MATERIALIZE command FROM CommandTemplate WITH handler, prerun, postrun, flags, AND subcommands

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: materialize command from CommandTemplate including handler, prerun, postrun, flags, and subcommands.

PROCEDURE DEFAULTROOTCOMMANDBUILDER_BUILDCOMMAND(template):
  cmd = command from template name, short, long, example
  IF template.Handler SET RunE; IF PreRun SET PreRunE; IF PostRun SET PostRunE
  flagManager.AddFlags(cmd, template.FlagSet)
  FOR EACH sub IN template.Subcommands: cmd.AddCommand(sub)
  RETURN cmd

## DEFAULTROOTCOMMANDBUILDER_WITHGLOBALFLAGS

SPEC-ID: IMPL-CLI_FRAMEWORK::DEFAULTROOTCOMMANDBUILDER_WITHGLOBALFLAGS
STEP T001: ADD persistent help AND verbose flags VIA provided OR embedded FlagManager

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: add persistent help and verbose flags via provided or embedded FlagManager.

PROCEDURE DEFAULTROOTCOMMANDBUILDER_WITHGLOBALFLAGS(cmd, flagMgr):
  IF flagMgr != nil THEN flagMgr.AddGlobalFlags(cmd)
  ELSE embedded flagManager.AddGlobalFlags(cmd)
  RETURN cmd

## NEWCLIAPP

SPEC-ID: IMPL-CLI_FRAMEWORK::NEWCLIAPP
STEP T001: WIRE all managers AND builders THEN create root command FROM AppInfo

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: wire all managers and builders then create root command from AppInfo.

PROCEDURE NEWCLIAPP(info):
  versionMgr = NEWVERSIONMANAGER()
  flagMgr = NEWFLAGMANAGER()
  rootBuilder = NEWROOTCOMMANDBUILDER(flagMgr, versionMgr)
  commandBuilder = NEWCOMMANDBUILDER(flagMgr)
  contextMgr = NEWCONTEXTMANAGER()
  dryRunMgr = NEWDryRunMANAGER()
  rootCmd = rootBuilder.NewRootCommand(info)
  RETURN CLIApp with Info, RootCommand, managers, builders

## CLIAPP_ADDCOMMAND

SPEC-ID: IMPL-CLI_FRAMEWORK::CLIAPP_ADDCOMMAND
STEP T001: ATTACH subcommand TO application root WHEN non-nil

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: attach subcommand to application root.

PROCEDURE CLIAPP_ADDCOMMAND(cmd):
  IF cmd == nil THEN
    RETURN
  END IF
  ATTACH cmd as child of app.RootCommand
  RETURN

## CLIAPP_EXECUTE

SPEC-ID: IMPL-CLI_FRAMEWORK::CLIAPP_EXECUTE
STEP T001: EXECUTE root command WITHOUT extra signal context AND propagate errors

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: run root command without extra signal context.

PROCEDURE CLIAPP_EXECUTE():
  IF app.RootCommand == nil THEN
    RETURN error("root command is required")
  END IF
  EXECUTE app.RootCommand with current process args/context
  PROPAGATE execute error to caller

## CLIAPP_EXECUTEWITHCONTEXT

SPEC-ID: IMPL-CLI_FRAMEWORK::CLIAPP_EXECUTEWITHCONTEXT
STEP T001: RUN root WITH signal-handled context AND exit non-zero ON error

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: run root with signal-handled context on command; exit non-zero on error.

PROCEDURE CLIAPP_EXECUTEWITHCONTEXT():
  ctx, cancel = WITHSIGNALHANDLING(nil)
  DEFER cancel()
  SET root command context to ctx
  IF Execute fails THEN write error to stderr AND exit 1
  RETURN nil

## DEFAULTFLAGMANAGER_ADDFLAGS

SPEC-ID: IMPL-CLI_FRAMEWORK::DEFAULTFLAGMANAGER_ADDFLAGS
STEP T001: REGISTER dry-run, note, config, incremental, AND list flags WHEN FlagSet pointers non-nil

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: register dry-run, note, config, incremental, and list flags when FlagSet pointers are non-nil.

PROCEDURE DEFAULTFLAGMANAGER_ADDFLAGS(cmd, flagSet):
  IF flagSet.DryRun != nil THEN AddDryRunFlag(cmd, flagSet.DryRun)
  IF flagSet.Note != nil THEN AddNoteFlag(cmd, flagSet.Note)
  IF flagSet.Config != nil THEN AddConfigFlag(cmd, flagSet.Config)
  IF flagSet.Incremental != nil THEN AddIncrementalFlag(cmd, flagSet.Incremental)
  IF flagSet.List != nil THEN AddListFlag(cmd, flagSet.List)
  RETURN success

## DEFAULTFLAGMANAGER_ADDGLOBALFLAGS

SPEC-ID: IMPL-CLI_FRAMEWORK::DEFAULTFLAGMANAGER_ADDGLOBALFLAGS
STEP T001: REGISTER persistent help AND verbose bool flags ON command

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: add persistent help and verbose flags to command.

PROCEDURE DEFAULTFLAGMANAGER_ADDGLOBALFLAGS(cmd):
  IF cmd == nil THEN
    RETURN error("command is required")
  END IF
  REGISTER persistent bool flag "help" with shorthand "h" default false and help text
  REGISTER persistent bool flag "verbose" with shorthand "v" default false and help text
  RETURN nil

## DEFAULTCONTEXTMANAGER_CREATE

SPEC-ID: IMPL-CLI_FRAMEWORK::DEFAULTCONTEXTMANAGER_CREATE
STEP T001: RETURN cancellable child context defaulting parent TO background

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: return cancellable child context defaulting parent to background.

PROCEDURE DEFAULTCONTEXTMANAGER_CREATE(parent):
  IF parent == nil THEN
    SET parent = BACKGROUND_CONTEXT()
  END IF
  CREATE childCtx, cancel USING cancellable child of parent
  RETURN childCtx, cancel

## DEFAULTCONTEXTMANAGER_WITHTIMEOUT

SPEC-ID: IMPL-CLI_FRAMEWORK::DEFAULTCONTEXTMANAGER_WITHTIMEOUT
STEP T001: PARSE duration string OR return cancel-only context ON parse failure

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: parse duration string; on parse failure return cancel-only context without timeout.

PROCEDURE DEFAULTCONTEXTMANAGER_WITHTIMEOUT(parent, timeout):
  IF parent == nil THEN parent = BACKGROUND
  IF PARSE_DURATION(timeout) fails THEN RETURN WithCancel(parent)
  RETURN WithTimeout(parent, duration)

## DEFAULTCONTEXTMANAGER_HANDLESIGNALS

SPEC-ID: IMPL-CLI_FRAMEWORK::DEFAULTCONTEXTMANAGER_HANDLESIGNALS
STEP T001: LISTEN for INT AND TERM AND invoke cancel WHEN signal received

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: listen for interrupt/terminate and invoke cancel when signal received.

PROCEDURE DEFAULTCONTEXTMANAGER_HANDLESIGNALS(cancel):
  CREATE buffered signal channel with capacity 1
  REGISTER channel for INT and TERM notifications
  START background listener:
    WAIT for first signal event
    IF cancel != nil THEN INVOKE cancel()
  RETURN

## NEWCANCELLABLEOPERATION

SPEC-ID: IMPL-CLI_FRAMEWORK::NEWCANCELLABLEOPERATION
STEP T001: WRAP function AS operation returning canceled WHEN Cancel was called

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: wrap function as operation that returns canceled when Cancel was called.

PROCEDURE NEWCANCELLABLEOPERATION(op):
  RETURN SimpleCancellableOperation{operation: op, cancelled: false}
  ON Execute: IF cancelled RETURN Canceled; ELSE RUN op(ctx)
  ON Cancel: SET cancelled = true

## WITHSIGNALHANDLING

SPEC-ID: IMPL-CLI_FRAMEWORK::WITHSIGNALHANDLING
STEP T001: CREATE cancelable context that cancels ON INT, TERM, OR parent done

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: create cancelable context that cancels on INT/TERM or parent done.

PROCEDURE WITHSIGNALHANDLING(parent):
  ctx, cancel = WithCancel(parent)
  NOTIFY on INT and TERM
  ON signal OR ctx.Done: cancel()
  RETURN ctx, cancel

## DEFAULTDRYRUNMANAGER_EXECUTE

SPEC-ID: IMPL-CLI_FRAMEWORK::DEFAULTDRYRUNMANAGER_EXECUTE
STEP T001: LOG Describe output WHEN DryRun ELSE run operation Execute

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: log Describe output when DryRun; otherwise run operation Execute.

PROCEDURE DEFAULTDRYRUNMANAGER_EXECUTE(ctx, op):
  IF ctx.DryRun THEN
    description = op.Describe()
    CALL Log(ctx, description)
    RETURN nil
  END IF
  RETURN op.Execute(ctx)

## DEFAULTVERSIONMANAGER_FORMATVERSION

SPEC-ID: IMPL-CLI_FRAMEWORK::DEFAULTVERSIONMANAGER_FORMATVERSION
STEP T001: FORMAT version, build date, AND platform INTO display string

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: format version, build date, and platform into display string.

PROCEDURE DEFAULTVERSIONMANAGER_FORMATVERSION(info):
  READ versionValue from info.Version
  READ buildDate from info.Date
  READ targetPlatform from info.Platform
  RETURN string "{versionValue} (compiled {buildDate}) [{targetPlatform}]"

## DEFAULTVERSIONMANAGER_CREATEVERSIONCOMMAND

SPEC-ID: IMPL-CLI_FRAMEWORK::DEFAULTVERSIONMANAGER_CREATEVERSIONCOMMAND
STEP T001: CREATE version subcommand that prints FormatVersion ON run

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: create version subcommand that prints FormatVersion on run.

PROCEDURE DEFAULTVERSIONMANAGER_CREATEVERSIONCOMMAND(info):
  DECLARE cmd as empty command object
  SET cmd.Use = "version"
  SET cmd.Short = "Display version information"
  SET cmd.Long = "Display detailed version information including build date and platform."
  SET cmd.Run to function that prints FormatVersion(info)
  RETURN cmd

## EMBEDDED_MINITEST: command builder

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify NewCommand sets use/short and WithHandler invokes RunE.

## EMBEDDED_MINITEST: root command builder

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify NewRootCommand sets name, short, and non-empty version from AppInfo.

## EMBEDDED_MINITEST: CLI app

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify NewCLIApp preserves AppInfo and AddCommand registers subcommand on root.

## EMBEDDED_MINITEST: dry-run manager

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify dry-run skips Execute and logs [DRY-RUN] prefix; non-dry-run runs operation.

## EMBEDDED_MINITEST: context manager

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify Create yields active context until cancel; WithTimeout ends after duration.

## EMBEDDED_MINITEST: flag manager

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify AddDryRunFlag and AddNoteFlag register expected flags.

## EMBEDDED_MINITEST: version manager

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify FormatVersion, CreateVersionTemplate, and CreateVersionCommand outputs.

## EMBEDDED_MINITEST: cancellable operation

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify Execute runs until Cancel then returns canceled.

## EMBEDDED_MINITEST: signal handling

- [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify WithSignalHandling context is active until manual cancel.
