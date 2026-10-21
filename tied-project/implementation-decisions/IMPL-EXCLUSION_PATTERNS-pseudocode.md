# [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-IMMUTABLE_FILE_EXCLUSION] [REQ-CONFIGURATION] — Directory and glob exclusion matching for backup relative paths.

Grammar-Version: v2

## MATCHES_DIRECTORY_PATTERN

- [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-IMMUTABLE_FILE_EXCLUSION] [REQ-CONFIGURATION] How: dispatch trailing-slash patterns through segment, glob, root-anchored, and any-depth multi-segment rules.
- Contract:
  - INPUT: normalizedPath (forward slashes), normalizedPattern (forward slashes, may end with /)
  - PRE: path and pattern are non-empty after normalization; pattern ends with /
  - OUTPUT: boolean excludedForThisPattern
  - POST:
    - on success with match: true only when path is under the directory implied by pattern per rule branch
    - on success without match: false when no branch matches
  - EFFECTS: pure
  - TERMINATION: total
procedure MATCHES_DIRECTORY_PATTERN:
  # [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-IMMUTABLE_FILE_EXCLUSION] How: strip trailing slash to patternBase then branch.
  Contract:
    INPUT: normalizedPath, normalizedPattern
    OUTPUT: boolean
    PRE: normalizedPattern ends with /
    POST: result reflects directory-only semantics (not plain file named patternBase)
    EFFECTS: pure
  SET patternBase TO TrimSuffix(normalizedPattern, "/")
  IF normalizedPattern starts with "**/" THEN
    - APPLY legacy **/dir segment containment (unchanged)
    - RETURN result
  IF patternBase contains no "/" THEN
    - RETURN pathHasDirSegment(normalizedPath, patternBase)
  IF normalizedPattern starts with "/" THEN
    - SET anchoredBase TO TrimPrefix(patternBase, "/")
    - RETURN normalizedPath equals anchoredBase OR normalizedPath starts with anchoredBase + "/"
  RETURN PATH_UNDER_DIR_PREFIX_ANY_DEPTH(normalizedPath, patternBase)

## PATH_UNDER_DIR_PREFIX_ANY_DEPTH

- [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-IMMUTABLE_FILE_EXCLUSION] How: segment-safe prefix match so multi-segment dir patterns apply at any depth without substring false positives.
- Contract:
  - INPUT: normalizedPath, patternBase (no trailing slash, contains at least one /)
  - PRE: patternBase is non-empty and contains /
  - OUTPUT: boolean
  - POST:
    - on success true: normalizedPath equals patternBase OR normalizedPath starts with patternBase + "/" at a segment boundary
    - on success false: path does not contain patternBase as a full path prefix
  - EFFECTS: pure
  - TERMINATION: total
procedure PATH_UNDER_DIR_PREFIX_ANY_DEPTH:
  # [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-IMMUTABLE_FILE_EXCLUSION] How: require "/" + patternBase + "/" in path or path prefix/suffix equality at boundaries.
  Contract:
    INPUT: normalizedPath, patternBase
    OUTPUT: boolean
    PRE: patternBase contains /
    POST: foox.venv/lib style substring matches return false
    EFFECTS: pure
  IF normalizedPath equals patternBase THEN
    - RETURN true
  IF normalizedPath starts with patternBase + "/" THEN
    - RETURN true
  SET needle TO "/" + patternBase + "/"
  IF normalizedPath contains needle THEN
    - RETURN true
  IF normalizedPath ends with "/" + patternBase THEN
    - RETURN true
  RETURN false
