#!/bin/bash
# [REQ-BKPDIR_CLI_OUTPUT_STANDARD] Shell adapter for operator output stream contract.
# Source from scripts that implement usage()/help paths.

_cli_out_help_stream_init() {
  _CLI_HELP_FD=2
  if [[ "${TIED_CLI_HELP_STDOUT:-}" == "1" || "${CLI_OUT_STDOUT_FLAG:-}" == "1" ]]; then
    _CLI_HELP_FD=1
  fi
}

_cli_out_plain_init() {
  _CLI_PLAIN=0
  if [[ -n "${TIED_CLI_PLAIN:-}" && "${TIED_CLI_PLAIN}" != "0" ]]; then
    _CLI_PLAIN=1
  fi
  if [[ -n "${NO_COLOR:-}" ]]; then
    _CLI_PLAIN=1
  fi
  if [[ "${_CLI_HELP_FD}" -eq 2 ]] && ! [[ -t 2 ]]; then
    _CLI_PLAIN=1
  fi
  if [[ "${_CLI_HELP_FD}" -eq 1 ]] && ! [[ -t 1 ]]; then
    _CLI_PLAIN=1
  fi
}

_cli_out_strip_ansi() {
  sed $'s/\033\\[[0-9;]*m//g'
}

# Read heredoc from stdin; write to selected help fd (stderr by default).
_cli_out_write_heredoc() {
  _cli_out_help_stream_init
  _cli_out_plain_init
  local body
  body="$(cat)"
  if [[ "${_CLI_PLAIN}" -eq 1 ]]; then
    body="$(printf '%s' "$body" | _cli_out_strip_ansi)"
  fi
  printf '%s' "$body" >&${_CLI_HELP_FD}
}

# usage: _cli_out_emit_usage "line1" "line2" ...
_cli_out_emit_usage() {
  _cli_out_help_stream_init
  _cli_out_plain_init
  for line in "$@"; do
    if [[ "${_CLI_PLAIN}" -eq 1 ]]; then
      echo "$line" | _cli_out_strip_ansi >&${_CLI_HELP_FD}
    else
      echo -e "$line" >&${_CLI_HELP_FD}
    fi
  done
}
