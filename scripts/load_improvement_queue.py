"""Load improvement queue YAML (skips leading # comment header)."""
from __future__ import annotations

from pathlib import Path

import yaml

REPO = Path(__file__).resolve().parents[1]
DEFAULT_QUEUE = REPO / "tied/docs/impl-pseudocode-improvement-queue.yaml"


def load_queue(path: Path | None = None) -> dict:
    p = path or DEFAULT_QUEUE
    lines = [ln for ln in p.read_text().splitlines() if not ln.strip().startswith("#")]
    return yaml.safe_load("\n".join(lines)) or {}


def save_queue(data: dict, path: Path | None = None) -> None:
    p = path or DEFAULT_QUEUE
    old = p.read_text() if p.exists() else ""
    header_lines = [ln for ln in old.splitlines() if ln.strip().startswith("#")]
    if not header_lines:
        header_lines = ["# IMPL Logic Alignment improvement queue"]
    body = yaml.dump(data, sort_keys=False, allow_unicode=True)
    p.write_text("\n".join(header_lines) + "\n" + body)
