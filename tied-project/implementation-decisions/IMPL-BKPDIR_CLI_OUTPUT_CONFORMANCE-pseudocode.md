# [IMPL-BKPDIR_CLI_OUTPUT_CONFORMANCE] [ARCH-BKPDIR_CLI_OUTPUT_STANDARD] [REQ-BKPDIR_CLI_OUTPUT_STANDARD]

## Summary contract

Shared operator output policy for bkpdir and adopted repo CLIs: resolve help/diagnostic streams, plain mode, and profile adapter; apply to Cobra roots; peel global texture flags before auto-detection; conformance harness for subprocess tests.

INPUT: argv slice, environment, optional TTY probes on selected writers
OUTPUT: resolved Policy, peeled argv, help/diagnostic writers
DATA: ProfileVersion, HelpDestination, PlainMode

## RESOLVE_POLICY

- [IMPL-BKPDIR_CLI_OUTPUT_CONFORMANCE] [ARCH-BKPDIR_CLI_OUTPUT_STANDARD] [REQ-BKPDIR_CLI_OUTPUT_STANDARD] — How: merge env (TIED_CLI_PLAIN, TIED_CLI_HELP_STDOUT, NO_COLOR) with peeled flags; default help to stderr; plain when any trigger or non-TTY help writer.

PROCEDURE RESOLVE_POLICY(env, peeledFlags, helpWriterIsTTY):
Contract:
INPUT: env, peeledFlags, helpWriterIsTTY
OUTPUT: Policy
PRE: true
POST: Policy.HelpDestination is stdout IFF explicit override else stderr
EFFECTS: pure
  SET helpDest = stderr
  IF peeledFlags.stdout OR env.TIED_CLI_HELP_STDOUT THEN helpDest = stdout
  SET plain = peeledFlags.plain OR env.TIED_CLI_PLAIN OR env.NO_COLOR non-empty OR NOT helpWriterIsTTY
  RETURN Policy{HelpDestination: helpDest, Plain: plain}

## PEEL_TEXTURE_ARGS

- [IMPL-BKPDIR_CLI_OUTPUT_CONFORMANCE] [ARCH-BKPDIR_CLI_OUTPUT_STANDARD] [REQ-BKPDIR_CLI_OUTPUT_STANDARD] — How: remove global texture tokens from argv before Cobra or auto-detect parsing; never consume -n as plain/limit conflict.

PROCEDURE PEEL_TEXTURE_ARGS(argv):
Contract:
INPUT: argv
OUTPUT: remaining, peeledFlags
PRE: true
POST: remaining excludes --plain, -p, --no-color, --stdout and their-only forms
EFFECTS: pure
  FOR each token in argv
    IF token in texture globals THEN record on peeledFlags AND skip
    ELSE append to remaining
  RETURN remaining, peeledFlags

## APPLY_TO_ROOT

- [IMPL-BKPDIR_CLI_OUTPUT_CONFORMANCE] [ARCH-BKPDIR_CLI_OUTPUT_STANDARD] [REQ-BKPDIR_CLI_OUTPUT_STANDARD] — How: register persistent texture flags, SetOut stdout for machine data, SetErr stderr, custom HelpFunc writes to Policy help writer; failures stay on stderr.

PROCEDURE APPLY_TO_ROOT(cmd, policy):
Contract:
INPUT: cmd, policy
OUTPUT: configured cmd
PRE: cmd != nil
POST: cmd.HelpFunc writes only to policy help writer
EFFECTS: mutates cmd
DATA_TRANSITION: cmd.Out/Err/HelpFunc and persistent flags
  REGISTER persistent --plain, --stdout on cmd
  SET cmd.Out = stdout
  SET cmd.Err = stderr
  SET cmd.HelpFunc = lambda that renders help to policy.HelpWriter()
  RETURN cmd

## EXECUTE_WITH_POLICY

- [IMPL-BKPDIR_CLI_OUTPUT_CONFORMANCE] [ARCH-BKPDIR_CLI_OUTPUT_STANDARD] [REQ-BKPDIR_CLI_OUTPUT_STANDARD] — How: peel argv, resolve policy, apply to root, delegate to existing auto-detection router with peeled args.

PROCEDURE EXECUTE_WITH_POLICY(rootCmd, argv, env):
Contract:
INPUT: rootCmd, argv, env
OUTPUT: error
PRE: rootCmd != nil
POST: execution errors returned to caller for stderr print; help obeys policy
EFFECTS: IO
FAILURE_MODES: executeWithAutoDetection returns error on Cobra failure
  remaining, flags = PEEL_TEXTURE_ARGS(argv)
  policy = RESOLVE_POLICY(env, flags, TTY(policy.HelpWriter()))
  APPLY_TO_ROOT(rootCmd, policy)
  RETURN executeWithAutoDetection(rootCmd, remaining)

## PRINT_EXEC_ERROR

SPEC-ID: IMPL-BKPDIR_CLI_OUTPUT_CONFORMANCE::PRINT_EXEC_ERROR
STEP T001: FPRINT execution failures to stderr before non-zero exit

- [IMPL-BKPDIR_CLI_OUTPUT_CONFORMANCE] [ARCH-BKPDIR_CLI_OUTPUT_STANDARD] [REQ-BKPDIR_CLI_OUTPUT_STANDARD] — How: main prints execution failures to stderr before non-zero exit.

PROCEDURE PRINT_EXEC_ERROR(err):
  # [IMPL-BKPDIR_CLI_OUTPUT_CONFORMANCE] — Layer B procedure token link
Contract:
INPUT: err
OUTPUT: none
PRE: err != nil
POST: bytes written only to stderr
EFFECTS: IO
  FPRINT stderr, err
