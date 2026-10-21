# [IMPL-CFG_INHERITANCE_PATH_RESOLUTION] [ARCH-CFG_005] [REQ-CFG_005] [REQ-CONFIGURATION]

## Summary contract

Resolve inherited config file paths relative to parent directories, expand home and environment variables, detect cycles, and build ordered inheritance chains.

INPUT: configPath, basePath, inherit metadata
OUTPUT: resolved absolute paths, inheritance chain file list
DATA: inheritanceChain, pathResolver, visited map

## EXPAND_PATH

SPEC-ID: IMPL-CFG_INHERITANCE_PATH_RESOLUTION::EXPAND_PATH
STEP T001: EXPAND tilde to home and environment variables in path

- [IMPL-CFG_INHERITANCE_PATH_RESOLUTION] [ARCH-CFG_005] [REQ-CFG_005] [REQ-CONFIGURATION] — How: expand tilde-prefixed paths to home directory and expand environment variables in path text.

PROCEDURE EXPAND_PATH(path):
  IF path empty THEN ERROR
  IF path starts with "~/ THEN JOIN home directory with remainder
  RETURN EXPAND_ENV(path)

## RESOLVE_PATH

SPEC-ID: IMPL-CFG_INHERITANCE_PATH_RESOLUTION::RESOLVE_PATH
STEP T001: EXPAND path then return absolute or join with base directory

- [IMPL-CFG_INHERITANCE_PATH_RESOLUTION] [ARCH-CFG_005] [REQ-CFG_005] [REQ-CONFIGURATION] — How: expand path then return absolute clean path or join with base directory of parent config file.

PROCEDURE RESOLVE_PATH(path, basePath):
  expanded = EXPAND_PATH(path)
  IF expanded is absolute THEN RETURN CLEAN(expanded)
  IF basePath empty THEN RETURN expanded
  baseDir = directory of basePath when base is file ELSE basePath when directory
  RETURN CLEAN(JOIN(baseDir, expanded))

## BUILD_CHAIN_RECURSIVE

SPEC-ID: IMPL-CFG_INHERITANCE_PATH_RESOLUTION::BUILD_CHAIN_RECURSIVE
STEP T001: RESOLVE path and detect circular dependency
STEP T002: RECURSE inherit parents then APPEND file to chain

- [IMPL-CFG_INHERITANCE_PATH_RESOLUTION] [ARCH-CFG_005] [REQ-CFG_005] [REQ-CONFIGURATION] — How: resolve path, detect circular visit, load inherit list, recurse parents with parent directory as base, append file to chain.

PROCEDURE BUILD_CHAIN_RECURSIVE(configPath, basePath, chain):
  resolved = RESOLVE_PATH(configPath, basePath)
  IF resolved in chain.visited THEN ERROR circular dependency
  MARK visited, VALIDATE file exists
  parents = LOAD inherit list FROM resolved file
  FOR EACH parent IN parents:
    BUILD_CHAIN_RECURSIVE(parent, PARENT_DIR(resolved), chain)
  APPEND resolved TO chain.files
