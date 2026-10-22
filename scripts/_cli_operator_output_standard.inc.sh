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

# Appearance profile helpers (decorated when TTY + not plain).
_cli_out_profile_rule() {
  _cli_out_help_stream_init
  _cli_out_plain_init
  if [[ "${_CLI_PLAIN}" -eq 1 ]]; then
    printf '%*s\n' 56 '' | tr ' ' '-' >&${_CLI_HELP_FD}
  else
    printf '\033[2m%*s\033[0m\n' 56 '' | tr ' ' $'\342\224\200' >&${_CLI_HELP_FD}
  fi
}

# usage: _cli_out_profile_section "Section title:"
_cli_out_profile_section() {
  _cli_out_help_stream_init
  _cli_out_plain_init
  local title="$1"
  if [[ "${_CLI_PLAIN}" -eq 1 ]]; then
    printf '%s\n' "$title" >&${_CLI_HELP_FD}
  else
    printf '\033[1m\033[36m▶ %s\033[0m\n' "$title" >&${_CLI_HELP_FD}
  fi
  _cli_out_profile_rule
}

# usage: _cli_out_profile_cmd_row "command-name" "description"
_cli_out_profile_cmd_row() {
  _cli_out_help_stream_init
  _cli_out_plain_init
  local name="$1"
  local desc="$2"
  if [[ "${_CLI_PLAIN}" -eq 1 ]]; then
    printf '  %-26s%s\n' "$name" "$desc" >&${_CLI_HELP_FD}
  else
    printf '  \033[1m%s\033[0m' "$name" >&${_CLI_HELP_FD}
    local pad=$((26 - ${#name}))
    if (( pad < 1 )); then pad=1; fi
    printf '%*s\033[2m%s\033[0m\n' "$pad" '' "$desc" >&${_CLI_HELP_FD}
  fi
}

# usage: _cli_out_profile_dim "dim body line"
_cli_out_profile_dim() {
  _cli_out_help_stream_init
  _cli_out_plain_init
  local line="$1"
  if [[ "${_CLI_PLAIN}" -eq 1 ]]; then
    printf '%s\n' "$line" >&${_CLI_HELP_FD}
  else
    printf '\033[2m%s\033[0m\n' "$line" >&${_CLI_HELP_FD}
  fi
}
