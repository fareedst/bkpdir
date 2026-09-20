#!/usr/bin/env python3
"""Parse IMPL tokens from tied/spec/formal-spec-registry.yaml. [REQ-PSEUDOCODE_FORMAL_VERIFICATION]"""
from __future__ import annotations

import re
import sys
from pathlib import Path

_IMPL_LINE = re.compile(r'^\s+- (?:"(IMPL-[^"]+)"|(IMPL-[A-Z0-9_]+))\s*$')


def registry_path(explicit: Path | None = None) -> Path:
    if explicit is not None:
        return explicit
    env = __import__("os").environ.get("FORMAL_SPEC_REGISTRY_PATH")
    if env:
        return Path(env)
    return Path("tied/spec/formal-spec-registry.yaml")


def iter_tokens(path: Path | None = None) -> list[str]:
    p = registry_path(path)
    tokens: list[str] = []
    for line in p.read_text(encoding="utf-8").splitlines():
        m = _IMPL_LINE.match(line)
        if m:
            tokens.append(m.group(1) or m.group(2) or "")
    return [t for t in tokens if t]


def count_tokens(path: Path | None = None) -> int:
    return len(iter_tokens(path))


def main(argv: list[str] | None = None) -> int:
    args = argv if argv is not None else sys.argv[1:]
    if len(args) != 1 or args[0] not in ("tokens", "count"):
        print("usage: formal_spec_registry.py tokens|count", file=sys.stderr)
        return 2
    cmd = args[0]
    if cmd == "tokens":
        for token in iter_tokens():
            print(token)
    else:
        print(count_tokens())
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
