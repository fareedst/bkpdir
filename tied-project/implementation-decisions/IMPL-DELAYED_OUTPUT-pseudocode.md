# [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING]

## Summary contract

Buffer stdout and stderr messages in OutputCollector until Flush* or Clear, and route OutputFormatter Print methods through the collector when delayed mode is enabled.

INPUT: message content, destination, message type
OUTPUT: deferred writes to os.Stdout/Stderr
DATA: OutputMessage slice, OutputFormatter.collector pointer

## NEW_OUTPUT_COLLECTOR

SPEC-ID: IMPL-DELAYED_OUTPUT::NEW_OUTPUT_COLLECTOR
STEP T001: RETURN empty OutputCollector WITH messages slice

- [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: return an empty OutputCollector ready to append messages.

PROCEDURE NEW_OUTPUT_COLLECTOR():
  # [IMPL-DELAYED_OUTPUT] — Layer B procedure token link
Contract:
INPUT: context for NEW_OUTPUT_COLLECTOR
OUTPUT: result of NEW_OUTPUT_COLLECTOR
PRE: true
POST: true
EFFECTS: pure
  RETURN &OutputCollector{messages: []}

## ADD_STDOUT

SPEC-ID: IMPL-DELAYED_OUTPUT::ADD_STDOUT
STEP T001: APPEND OutputMessage WITH Destination stdout AND messageType

- [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: append OutputMessage with Destination stdout and given message type.

PROCEDURE ADD_STDOUT(content, messageType):
  # [IMPL-DELAYED_OUTPUT] — Layer B procedure token link
Contract:
INPUT: content, messageType
OUTPUT: result of ADD_STDOUT
PRE: true
POST: true
EFFECTS: pure
  APPEND OutputMessage{content, "stdout", messageType}

## ADD_STDERR

SPEC-ID: IMPL-DELAYED_OUTPUT::ADD_STDERR
STEP T001: APPEND OutputMessage WITH Destination stderr AND messageType

- [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: append OutputMessage with Destination stderr and given message type.

PROCEDURE ADD_STDERR(content, messageType):
  # [IMPL-DELAYED_OUTPUT] — Layer B procedure token link
Contract:
INPUT: content, messageType
OUTPUT: result of ADD_STDERR
PRE: true
POST: true
EFFECTS: pure
  APPEND OutputMessage{content, "stderr", messageType}

## GET_MESSAGES

SPEC-ID: IMPL-DELAYED_OUTPUT::GET_MESSAGES
STEP T001: RETURN snapshot of buffered messages WITHOUT flushing

- [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: return a snapshot of all buffered messages without flushing.

PROCEDURE GET_MESSAGES():
  # [IMPL-DELAYED_OUTPUT] — Layer B procedure token link
Contract:
INPUT: context for GET_MESSAGES
OUTPUT: result of GET_MESSAGES
PRE: true
POST: true
EFFECTS: pure
  RETURN messages slice

## FLUSH_ALL

SPEC-ID: IMPL-DELAYED_OUTPUT::FLUSH_ALL
STEP T001: WRITE every message to stdout OR stderr
STEP T002: CLEAR messages buffer

- [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: write every message to stdout or stderr then clear the buffer.

PROCEDURE FLUSH_ALL():
  # [IMPL-DELAYED_OUTPUT] — Layer B procedure token link
Contract:
INPUT: context for FLUSH_ALL
OUTPUT: result of FLUSH_ALL
PRE: true
POST: true
EFFECTS: pure
  FOR EACH msg: write to matching stream
  CLEAR messages

## FLUSH_STDOUT

SPEC-ID: IMPL-DELAYED_OUTPUT::FLUSH_STDOUT
STEP T001: PRINT stdout messages AND retain stderr entries

- [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: write only stdout messages and retain stderr entries in the buffer.

PROCEDURE FLUSH_STDOUT():
  # [IMPL-DELAYED_OUTPUT] — Layer B procedure token link
Contract:
INPUT: context for FLUSH_STDOUT
OUTPUT: result of FLUSH_STDOUT
PRE: true
POST: true
EFFECTS: pure
  KEEP stderr messages; print stdout messages

## FLUSH_STDERR

SPEC-ID: IMPL-DELAYED_OUTPUT::FLUSH_STDERR
STEP T001: PRINT stderr messages AND retain stdout entries

- [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: write only stderr messages and retain stdout entries in the buffer.

PROCEDURE FLUSH_STDERR():
  # [IMPL-DELAYED_OUTPUT] — Layer B procedure token link
Contract:
INPUT: context for FLUSH_STDERR
OUTPUT: result of FLUSH_STDERR
PRE: true
POST: true
EFFECTS: pure
  KEEP stdout messages; print stderr messages

## CLEAR

SPEC-ID: IMPL-DELAYED_OUTPUT::CLEAR
STEP T001: DISCARD all buffered messages WITHOUT printing

- [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: discard all buffered messages without printing.

PROCEDURE CLEAR():
  # [IMPL-DELAYED_OUTPUT] — Layer B procedure token link
Contract:
INPUT: context for CLEAR
OUTPUT: result of CLEAR
PRE: true
POST: true
EFFECTS: pure
  messages = []

## IS_DELAYED_MODE

SPEC-ID: IMPL-DELAYED_OUTPUT::IS_DELAYED_MODE
STEP T001: RETURN formatter.collector != nil

- [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: report whether OutputFormatter has a non-nil collector attached.

PROCEDURE IS_DELAYED_MODE(formatter):
  # [IMPL-DELAYED_OUTPUT] — Layer B procedure token link
Contract:
INPUT: formatter
OUTPUT: result of IS_DELAYED_MODE
PRE: true
POST: true
EFFECTS: pure
  RETURN formatter.collector != nil

## GET_COLLECTOR

SPEC-ID: IMPL-DELAYED_OUTPUT::GET_COLLECTOR
STEP T001: RETURN attached collector pointer

- [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: return the attached collector pointer for tests and flush orchestration.

PROCEDURE GET_COLLECTOR(formatter):
  # [IMPL-DELAYED_OUTPUT] — Layer B procedure token link
Contract:
INPUT: formatter
OUTPUT: result of GET_COLLECTOR
PRE: true
POST: true
EFFECTS: pure
  RETURN formatter.collector

## SET_COLLECTOR

SPEC-ID: IMPL-DELAYED_OUTPUT::SET_COLLECTOR
STEP T001: ATTACH OR detach collector on formatter

- [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: attach or detach delayed-mode collector (nil disables buffering).

PROCEDURE SET_COLLECTOR(formatter, collector):
  # [IMPL-DELAYED_OUTPUT] — Layer B procedure token link
Contract:
INPUT: formatter, collector
OUTPUT: result of SET_COLLECTOR
PRE: true
POST: true
EFFECTS: pure
  formatter.collector = collector

## PRINT_ROUTING

SPEC-ID: IMPL-DELAYED_OUTPUT::PRINT_ROUTING
STEP T001: IF collector present THEN AddStdout OR AddStderr
STEP T002: ELSE direct print to stream

- [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: when collector present, AddStdout/AddStderr instead of immediate fmt.Print for formatted messages.

PROCEDURE PRINT_ROUTING(formatter, message, stream):
  # [IMPL-DELAYED_OUTPUT] — Layer B procedure token link
Contract:
INPUT: formatter, message, stream
OUTPUT: result of PRINT_ROUTING
PRE: true
POST: true
EFFECTS: pure
  IF collector: collector.Add* ELSE direct print

## NEW_OUTPUT_FORMATTER_WITH_COLLECTOR

SPEC-ID: IMPL-DELAYED_OUTPUT::NEW_OUTPUT_FORMATTER_WITH_COLLECTOR
STEP T001: CONSTRUCT AIFormatterAdapter WITH cfg AND pre-wired collector

- [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: construct AIFormatterAdapter with cfg and pre-wired OutputCollector for delayed CLI output.

PROCEDURE NEW_OUTPUT_FORMATTER_WITH_COLLECTOR(cfg, collector):
  # [IMPL-DELAYED_OUTPUT] — Layer B procedure token link
Contract:
INPUT: cfg, collector
OUTPUT: result of NEW_OUTPUT_FORMATTER_WITH_COLLECTOR
PRE: true
POST: true
EFFECTS: pure
  RETURN AIFormatterAdapter{cfg, collector}
