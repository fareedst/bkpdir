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
  RETURN Config{Enabled=true, Command="git", WorkingDirectory=".", IncludeSubmoduleInfo=false, CheckDirtyStatus=false}

## EXECUTE_GIT_COMMAND

SPEC-ID: IMPL-GIT_CLI::EXECUTE_GIT_COMMAND
STEP T001: RESOLVE git binary WITH legacy fallbacks
STEP T002: RUN in WorkingDirectory AND return trimmed stdout OR GitError

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: resolve Command/GitCommand/git binary, run in WorkingDirectory, return trimmed stdout or GitError.

PROCEDURE EXECUTE_GIT_COMMAND(args):
  cmd = exec.Command(gitCmd, args...); cmd.Dir = WorkingDirectory
  RETURN TrimSpace(output) OR GitError

## IS_REPOSITORY

SPEC-ID: IMPL-GIT_CLI::IS_REPOSITORY
STEP T001: EXECUTE rev-parse --is-inside-work-tree AND compare to true

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: git rev-parse --is-inside-work-tree equals true without error.

PROCEDURE IS_REPOSITORY():
  out = EXECUTE_GIT_COMMAND("rev-parse", "--is-inside-work-tree")
  RETURN err == nil AND out == "true"

## GET_BRANCH

SPEC-ID: IMPL-GIT_CLI::GET_BRANCH
STEP T001: REQUIRE IsRepository THEN rev-parse --abbrev-ref HEAD

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: rev-parse --abbrev-ref HEAD when inside a repository.

PROCEDURE GET_BRANCH():
  REQUIRE IsRepository; RETURN EXECUTE_GIT_COMMAND("rev-parse", "--abbrev-ref", "HEAD")

## GET_SHORT_HASH

SPEC-ID: IMPL-GIT_CLI::GET_SHORT_HASH
STEP T001: REQUIRE IsRepository THEN rev-parse --short HEAD

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: rev-parse --short HEAD when inside a repository.

PROCEDURE GET_SHORT_HASH():
  REQUIRE IsRepository; RETURN EXECUTE_GIT_COMMAND("rev-parse", "--short", "HEAD")

## IS_WORKING_DIRECTORY_CLEAN

SPEC-ID: IMPL-GIT_CLI::IS_WORKING_DIRECTORY_CLEAN
STEP T001: EXECUTE status --porcelain AND return len(out) == 0

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: git status --porcelain empty means clean working tree.

PROCEDURE IS_WORKING_DIRECTORY_CLEAN():
  out = EXECUTE_GIT_COMMAND("status", "--porcelain"); RETURN len(out) == 0

## GET_INFO

SPEC-ID: IMPL-GIT_CLI::GET_INFO
STEP T001: IF NOT IsRepository RETURN Info{IsRepo false}
STEP T002: FILL Branch, Hash FROM subcommands

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: populate Info with IsRepo, Branch, Hash, and IsClean from branch/hash/status helpers.

PROCEDURE GET_INFO():
  IF NOT IsRepository: RETURN Info{IsRepo: false}
  FILL Branch, Hash, IsClean from subcommands

## GET_INFO_WITH_STATUS

SPEC-ID: IMPL-GIT_CLI::GET_INFO_WITH_STATUS
STEP T001: GET_INFO THEN optionally check dirty status
STEP T002: IF configured INCLUDE submodule listing

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: extend GetInfo with optional clean check and submodule listing when configured.

PROCEDURE GET_INFO_WITH_STATUS():
  info = GET_INFO()
  IF config.CheckDirtyStatus: info.IsClean = IS_WORKING_DIRECTORY_CLEAN()
  IF config.IncludeSubmoduleInfo: info.Submodules = GET_SUBMODULES()
  RETURN info

## IS_SUBMODULE

SPEC-ID: IMPL-GIT_CLI::IS_SUBMODULE
STEP T001: EXECUTE rev-parse --show-superproject-working-tree AND check non-empty

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: rev-parse --show-superproject-working-tree non-empty means submodule checkout.

PROCEDURE IS_SUBMODULE():
  out = EXECUTE_GIT_COMMAND("rev-parse", "--show-superproject-working-tree")
  RETURN len(out) > 0

## GET_SUBMODULES

SPEC-ID: IMPL-GIT_CLI::GET_SUBMODULES
STEP T001: PARSE submodule status output INTO SubmoduleInfo records

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: parse git submodule status lines into SubmoduleInfo records with path and URL.

PROCEDURE GET_SUBMODULES():
  PARSE submodule status output INTO []SubmoduleInfo

## LEGACY_IS_GIT_REPOSITORY_WRAPPER

SPEC-ID: IMPL-GIT_CLI::LEGACY_IS_GIT_REPOSITORY_WRAPPER
STEP T001: CREATE ephemeral Repo WITH cwd AND delegate to IsRepository

- [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: delegate to pkg/git IsGitRepository wrapper around rev-parse --is-inside-work-tree.

PROCEDURE LEGACY_IS_GIT_REPOSITORY_WRAPPER(cwd):
  repo = NEW pkg/git Repo with cwd
  RETURN repo.IsRepository()
