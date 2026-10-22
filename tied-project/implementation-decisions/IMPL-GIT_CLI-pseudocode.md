# [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]

## Summary contract

pkg/git executes git subprocess commands with configurable working directory, detects repositories, extracts branch/hash/clean status, and aggregates Info for archive naming and legacy root wrappers.

INPUT: Config, working directory, git subcommand args
OUTPUT: trimmed command output, Info struct, errors as GitError
DATA: Command path legacy fallbacks, SubmoduleInfo list

## DEFAULT_CONFIG

SPEC-ID: IMPL-GIT_CLI::DEFAULT_CONFIG
STEP T001: RETURN Config WITH Enabled, Command git, WorkingDirectory, defaults

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: return Config with Enabled, Command git, WorkingDirectory, and submodule/status defaults.

PROCEDURE DEFAULT_CONFIG():
  # [IMPL-GIT_CLI] — Layer B procedure token link
Contract:
INPUT: context for DEFAULT_CONFIG
OUTPUT: result of DEFAULT_CONFIG
PRE: true
POST: true
EFFECTS: pure
  RETURN Config{Enabled=true, Command="git", WorkingDirectory=".", IncludeSubmoduleInfo=false, CheckDirtyStatus=false}

## EXECUTE_GIT_COMMAND

SPEC-ID: IMPL-GIT_CLI::EXECUTE_GIT_COMMAND
STEP T001: RESOLVE git binary WITH legacy fallbacks
STEP T002: RUN in WorkingDirectory AND return trimmed stdout OR GitError

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: resolve Command/GitCommand/git binary, run in WorkingDirectory, return trimmed stdout or GitError.

PROCEDURE EXECUTE_GIT_COMMAND(args):
  # [IMPL-GIT_CLI] — Layer B procedure token link
Contract:
INPUT: args
OUTPUT: result of EXECUTE_GIT_COMMAND
PRE: true
POST: true
EFFECTS: pure
FAILURE_MODES: ErrGitCommandFailed
  cmd = exec.Command(gitCmd, args...); cmd.Dir = WorkingDirectory
  RETURN TrimSpace(output) OR GitError

## IS_REPOSITORY

SPEC-ID: IMPL-GIT_CLI::IS_REPOSITORY
STEP T001: EXECUTE rev-parse --is-inside-work-tree AND compare to true

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: git rev-parse --is-inside-work-tree equals true without error.

PROCEDURE IS_REPOSITORY():
  # [IMPL-GIT_CLI] — Layer B procedure token link
Contract:
INPUT: context for IS_REPOSITORY
OUTPUT: result of IS_REPOSITORY
PRE: true
POST: true
EFFECTS: pure
  out = EXECUTE_GIT_COMMAND("rev-parse", "--is-inside-work-tree")
  RETURN err == nil AND out == "true"

## GET_BRANCH

SPEC-ID: IMPL-GIT_CLI::GET_BRANCH
STEP T001: REQUIRE IsRepository THEN rev-parse --abbrev-ref HEAD

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: rev-parse --abbrev-ref HEAD when inside a repository.

PROCEDURE GET_BRANCH():
  # [IMPL-GIT_CLI] — Layer B procedure token link
Contract:
INPUT: context for GET_BRANCH
OUTPUT: result of GET_BRANCH
PRE: true
POST: true
EFFECTS: pure
  REQUIRE IsRepository; RETURN EXECUTE_GIT_COMMAND("rev-parse", "--abbrev-ref", "HEAD")

## GET_SHORT_HASH

SPEC-ID: IMPL-GIT_CLI::GET_SHORT_HASH
STEP T001: REQUIRE IsRepository THEN rev-parse --short HEAD

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: rev-parse --short HEAD when inside a repository.

PROCEDURE GET_SHORT_HASH():
  # [IMPL-GIT_CLI] — Layer B procedure token link
Contract:
INPUT: context for GET_SHORT_HASH
OUTPUT: result of GET_SHORT_HASH
PRE: true
POST: true
EFFECTS: pure
  REQUIRE IsRepository; RETURN EXECUTE_GIT_COMMAND("rev-parse", "--short", "HEAD")

## IS_WORKING_DIRECTORY_CLEAN

SPEC-ID: IMPL-GIT_CLI::IS_WORKING_DIRECTORY_CLEAN
STEP T001: EXECUTE status --porcelain AND return len(out) == 0

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: git status --porcelain empty means clean working tree.

PROCEDURE IS_WORKING_DIRECTORY_CLEAN():
  # [IMPL-GIT_CLI] — Layer B procedure token link
Contract:
INPUT: context for IS_WORKING_DIRECTORY_CLEAN
OUTPUT: result of IS_WORKING_DIRECTORY_CLEAN
PRE: true
POST: true
EFFECTS: pure
  out = EXECUTE_GIT_COMMAND("status", "--porcelain"); RETURN len(out) == 0

## GET_INFO

SPEC-ID: IMPL-GIT_CLI::GET_INFO
STEP T001: IF NOT IsRepository RETURN Info{IsRepo false}
STEP T002: FILL Branch, Hash FROM subcommands

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: populate Info with IsRepo, Branch, Hash, and IsClean from branch/hash/status helpers.

PROCEDURE GET_INFO():
  # [IMPL-GIT_CLI] — Layer B procedure token link
Contract:
INPUT: context for GET_INFO
OUTPUT: result of GET_INFO
PRE: true
POST: true
EFFECTS: pure
  IF NOT IsRepository: RETURN Info{IsRepo: false}
  FILL Branch, Hash, IsClean from subcommands

## GET_INFO_WITH_STATUS

SPEC-ID: IMPL-GIT_CLI::GET_INFO_WITH_STATUS
STEP T001: GET_INFO THEN optionally check dirty status
STEP T002: IF configured INCLUDE submodule listing

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: extend GetInfo with optional clean check and submodule listing when configured.

PROCEDURE GET_INFO_WITH_STATUS():
  # [IMPL-GIT_CLI] — Layer B procedure token link
Contract:
INPUT: context for GET_INFO_WITH_STATUS
OUTPUT: result of GET_INFO_WITH_STATUS
PRE: true
POST: true
EFFECTS: pure
  info = GET_INFO()
  IF config.CheckDirtyStatus: info.IsClean = IS_WORKING_DIRECTORY_CLEAN()
  IF config.IncludeSubmoduleInfo: info.Submodules = GET_SUBMODULES()
  RETURN info

## IS_SUBMODULE

SPEC-ID: IMPL-GIT_CLI::IS_SUBMODULE
STEP T001: EXECUTE rev-parse --show-superproject-working-tree AND check non-empty

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: rev-parse --show-superproject-working-tree non-empty means submodule checkout.

PROCEDURE IS_SUBMODULE():
  # [IMPL-GIT_CLI] — Layer B procedure token link
Contract:
INPUT: context for IS_SUBMODULE
OUTPUT: result of IS_SUBMODULE
PRE: true
POST: true
EFFECTS: pure
  out = EXECUTE_GIT_COMMAND("rev-parse", "--show-superproject-working-tree")
  RETURN len(out) > 0

## GET_SUBMODULES

SPEC-ID: IMPL-GIT_CLI::GET_SUBMODULES
STEP T001: PARSE submodule status output INTO SubmoduleInfo records

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: parse git submodule status lines into SubmoduleInfo records with path and URL.

PROCEDURE GET_SUBMODULES():
  # [IMPL-GIT_CLI] — Layer B procedure token link
Contract:
INPUT: context for GET_SUBMODULES
OUTPUT: result of GET_SUBMODULES
PRE: true
POST: true
EFFECTS: pure
  PARSE submodule status output INTO []SubmoduleInfo

## LEGACY_IS_GIT_REPOSITORY_WRAPPER

SPEC-ID: IMPL-GIT_CLI::LEGACY_IS_GIT_REPOSITORY_WRAPPER
STEP T001: CREATE ephemeral Repo WITH cwd AND delegate to IsRepository

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: delegate to pkg/git IsGitRepository wrapper around rev-parse --is-inside-work-tree.

PROCEDURE LEGACY_IS_GIT_REPOSITORY_WRAPPER(cwd):
  # [IMPL-GIT_CLI] — Layer B procedure token link
Contract:
INPUT: cwd
OUTPUT: result of LEGACY_IS_GIT_REPOSITORY_WRAPPER
PRE: true
POST: true
EFFECTS: pure
  repo = NEW pkg/git Repo with cwd
  RETURN repo.IsRepository()
