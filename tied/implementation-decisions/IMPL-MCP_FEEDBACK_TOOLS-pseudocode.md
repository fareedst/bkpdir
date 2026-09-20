# [IMPL-MCP_FEEDBACK_TOOLS] [ARCH-FEEDBACK_STORAGE] [REQ-FEEDBACK_TO_TIED]

## Summary contract

MCP feedback tools load, append, and export client feedback entries to tied/feedback.yaml for upstream TIED reporting.

INPUT: basePath, entry type/title/description, export format
OUTPUT: feedback.yaml updates, markdown or JSON export strings
DATA: FeedbackFile with entries array, generated id and created_at timestamps

## GET_FEEDBACK_PATH

SPEC-ID: IMPL-MCP_FEEDBACK_TOOLS::GET_FEEDBACK_PATH
STEP T001: Resolve optional basePath to {base}/feedback.yaml defaulting to tied/feedback.yaml

- [IMPL-MCP_FEEDBACK_TOOLS] [ARCH-FEEDBACK_STORAGE] [REQ-FEEDBACK_TO_TIED] — How: resolve optional basePath to {base}/feedback.yaml defaulting to tied/feedback.yaml.

PROCEDURE GET_FEEDBACK_PATH(basePath):
  IF basePath provided THEN RETURN join(basePath, "feedback.yaml")
  RETURN default tied feedback path

## LOAD_FEEDBACK

SPEC-ID: IMPL-MCP_FEEDBACK_TOOLS::LOAD_FEEDBACK
STEP T001: Read feedback YAML or return empty entries structure when file is missing

- [IMPL-MCP_FEEDBACK_TOOLS] [ARCH-FEEDBACK_STORAGE] [REQ-FEEDBACK_TO_TIED] — How: read feedback YAML or return empty entries structure when file is missing.

PROCEDURE LOAD_FEEDBACK(basePath):
  path = GET_FEEDBACK_PATH(basePath)
  IF file missing THEN RETURN { entries: [] }
  PARSE yaml INTO FeedbackFile
  RETURN FeedbackFile

## APPEND_ENTRY

SPEC-ID: IMPL-MCP_FEEDBACK_TOOLS::APPEND_ENTRY
STEP T001: Validate type/title/description, assign id and created_at, append entry, write feedback.yaml

- [IMPL-MCP_FEEDBACK_TOOLS] [ARCH-FEEDBACK_STORAGE] [REQ-FEEDBACK_TO_TIED] — How: validate type/title/description, assign id and created_at, append entry, write feedback.yaml.

PROCEDURE APPEND_ENTRY(params):
  VALIDATE required fields
  id = GENERATE timestamp-based id
  APPEND entry TO feedback.entries
  WRITE feedback.yaml

## EXPORT_ENTRIES

SPEC-ID: IMPL-MCP_FEEDBACK_TOOLS::EXPORT_ENTRIES
STEP T001: Format all entries as markdown report or JSON string for MCP export tool

- [IMPL-MCP_FEEDBACK_TOOLS] [ARCH-FEEDBACK_STORAGE] [REQ-FEEDBACK_TO_TIED] — How: format all entries as markdown report or JSON string for MCP export tool.

PROCEDURE EXPORT_ENTRIES(entries, format):
  IF format == markdown THEN RETURN exportMarkdown(entries)
  RETURN exportJson(entries)

## MCP_HANDLER_ADD

SPEC-ID: IMPL-MCP_FEEDBACK_TOOLS::MCP_HANDLER_ADD
STEP T001: Parse tied_feedback_add args, call appendEntry, return ok/id/created_at and optional markdown snippet

- [IMPL-MCP_FEEDBACK_TOOLS] [ARCH-FEEDBACK_STORAGE] [REQ-FEEDBACK_TO_TIED] — How: parse tied_feedback_add args, call appendEntry, return ok/id/created_at and optional markdown snippet.

PROCEDURE MCP_HANDLER_ADD(args):
  PARSE args INTO entry params
  result = APPEND_ENTRY(params)
  RETURN MCP success envelope with optional report_snippet

## MCP_HANDLER_EXPORT

SPEC-ID: IMPL-MCP_FEEDBACK_TOOLS::MCP_HANDLER_EXPORT
STEP T001: Load entries then return exportMarkdown or exportJson string from tied_feedback_export handler

- [IMPL-MCP_FEEDBACK_TOOLS] [ARCH-FEEDBACK_STORAGE] [REQ-FEEDBACK_TO_TIED] — How: load entries then return exportMarkdown or exportJson string from tied_feedback_export handler.

PROCEDURE MCP_HANDLER_EXPORT(args):
  entries = LOAD_FEEDBACK(args.basePath).entries
  RETURN EXPORT_ENTRIES(entries, args.format)
