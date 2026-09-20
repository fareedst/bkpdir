#!/usr/bin/env python3
"""Project sidecar Summary contract into IMPL detail/index YAML summaries.

Sidecar markdown is canonical; YAML `implementation_approach.summary` is a lean projection.
[REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-TIED_FILES]
"""
from __future__ import annotations

import argparse
import json
import subprocess
import sys
from datetime import date
from pathlib import Path

import yaml

REPO = Path(__file__).resolve().parents[1]
TIED = REPO / "tied"
IMPL_DIR = TIED / "implementation-decisions"
INDEX_PATH = TIED / "implementation-decisions.yaml"
TIED_CLI = REPO / ".cursor/skills/tied-yaml/scripts/tied-cli.sh"
SYNC_REASON = "Sidecar-canonical projection sync ([REQ-PSEUDOCODE_FORMAL_VERIFICATION])"
# Methodology-owned IMPL detail files ([PROC-TIED_METHODOLOGY_READONLY]); project index only.
INDEX_ONLY_IMPL_TOKENS = frozenset(
    {"IMPL-MCP_FEEDBACK_TOOLS", "IMPL-MODULE_VALIDATION", "IMPL-TIED_FILES"}
)

sys.path.insert(0, str(REPO / "scripts"))
from formal_spec_registry import iter_tokens  # noqa: E402


def sidecar_path(token: str) -> Path:
    return IMPL_DIR / f"{token}-pseudocode.md"


def extract_summary_contract_body(text: str) -> list[str]:
    lines = text.splitlines()
    start: int | None = None
    for i, line in enumerate(lines):
        if line.strip() == "## Summary contract":
            start = i + 1
            break
    if start is None:
        return []
    body: list[str] = []
    for line in lines[start:]:
        if line.startswith("## "):
            break
        body.append(line)
    return body


def build_projected_summary(token: str, contract_lines: list[str]) -> str:
    lines_out: list[str] = []
    for raw in contract_lines:
        s = raw.strip()
        if not s:
            continue
        if s.startswith("INPUT:") or s.startswith("OUTPUT:") or s.startswith("DATA:"):
            lines_out.append(s)
            continue
        if not lines_out:
            lines_out.append(s)
        if len(lines_out) >= 5:
            break
    pointer = f"Canonical behavior: tied/implementation-decisions/{token}-pseudocode.md"
    if pointer not in lines_out:
        lines_out.append(pointer)
    return "\n".join(lines_out[:8])


def projected_summary_for_token(token: str) -> str | None:
    path = sidecar_path(token)
    if not path.exists():
        return None
    body = extract_summary_contract_body(path.read_text(encoding="utf-8"))
    if not body:
        return None
    return build_projected_summary(token, body)


def load_yaml_summary(path: Path, token: str) -> str | None:
    if not path.exists():
        return None
    data = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        return None
    rec = data.get(token) or data
    if not isinstance(rec, dict):
        return None
    impl = rec.get("implementation_approach") or {}
    summary = impl.get("summary")
    return summary if isinstance(summary, str) else None


def normalize_summary(text: str | None) -> str:
    if not text:
        return ""
    return text.strip().replace("\r\n", "\n")


def drift_for_token(token: str) -> tuple[str | None, str | None, str | None]:
    projected = projected_summary_for_token(token)
    detail = load_yaml_summary(IMPL_DIR / f"{token}.yaml", token)
    index = load_yaml_summary(INDEX_PATH, token)
    return projected, detail, index


def call_tied_cli(tool: str, payload: dict) -> dict:
    if not TIED_CLI.is_file():
        raise RuntimeError(f"tied-cli not found: {TIED_CLI}")
    proc = subprocess.run(
        [str(TIED_CLI), tool, json.dumps(payload)],
        cwd=REPO,
        capture_output=True,
        text=True,
        check=False,
    )
    if proc.returncode != 0:
        raise RuntimeError(f"{tool} failed: {proc.stderr or proc.stdout}")
    try:
        return json.loads(proc.stdout)
    except json.JSONDecodeError as exc:
        raise RuntimeError(f"{tool} non-JSON stdout: {proc.stdout[:500]}") from exc


def apply_token(token: str, summary: str, *, detail_writable: bool) -> None:
    today = date.today().isoformat()
    updates = {
        "implementation_approach": {"summary": summary},
        "metadata": {
            "last_updated": {
                "date": today,
                "reason": SYNC_REASON,
            }
        },
    }
    if detail_writable:
        detail_result = call_tied_cli(
            "yaml_detail_update", {"token": token, "updates": json.dumps(updates)}
        )
        if detail_result.get("ok") is False:
            err = str(detail_result.get("error", ""))
            if "methodology-owned" not in err:
                raise RuntimeError(f"yaml_detail_update failed for {token}: {err}")
    call_tied_cli(
        "yaml_index_update",
        {
            "index": "implementation",
            "token": token,
            "updates": json.dumps(
                {
                    "implementation_approach": {"summary": summary},
                    "metadata": {
                        "last_updated": {
                            "date": today,
                            "reason": SYNC_REASON,
                        }
                    },
                }
            ),
        },
    )


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--check", action="store_true", help="Exit 1 on projection drift")
    parser.add_argument("--apply", action="store_true", help="Write projections via tied-cli")
    parser.add_argument("--token", action="append", default=[], help="Limit to IMPL token(s)")
    args = parser.parse_args()
    if args.check == args.apply:
        parser.error("specify exactly one of --check or --apply")

    tokens = args.token or iter_tokens()
    errors: list[str] = []
    applied = 0
    for token in tokens:
        projected, detail, index = drift_for_token(token)
        if projected is None:
            errors.append(f"{token}: missing sidecar Summary contract")
            continue
        nd = normalize_summary(detail)
        ni = normalize_summary(index)
        np = normalize_summary(projected)
        index_only = token in INDEX_ONLY_IMPL_TOKENS
        if args.check:
            if not index_only and nd != np:
                errors.append(f"{token}: detail summary drift")
            if ni != np:
                errors.append(f"{token}: index summary drift")
            continue
        detail_ok = index_only or nd == np
        if detail_ok and ni == np:
            continue
        apply_token(token, projected, detail_writable=not index_only)
        applied += 1
        print(f"DEBUG: projected {token}")

    if args.apply:
        print(f"DIAGNOSTIC: applied {applied} token projection(s)")
        return 0

    if errors:
        for err in errors:
            print(f"FAIL: {err}", file=sys.stderr)
        print(f"DIAGNOSTIC: {len(errors)} drift/missing error(s)", file=sys.stderr)
        return 1
    print(f"DIAGNOSTIC: all {len(tokens)} projection(s) in sync")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
