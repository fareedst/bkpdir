#!/usr/bin/env python3
"""Ensure L4 mutation tracking fields exist on improvement-queue tokens.

Does NOT set mutation_verified — use after run-mutation-pilot.sh / run-mutation-waves.sh.
"""
from __future__ import annotations

from pathlib import Path

import yaml

REPO = Path(__file__).resolve().parents[1]
QUEUE = REPO / "tied/docs/impl-pseudocode-improvement-queue.yaml"
ORACLE_REG = REPO / "tied/spec/oracle-registry.yaml"


def pkg_for_impl(impl: str) -> list[str]:
    data = yaml.safe_load(ORACLE_REG.read_text())
    pkgs: list[str] = []
    for o in data.get("oracles", []):
        if o.get("impl") == impl:
            pkgs.extend(o.get("pkg_paths") or [])
    return list(dict.fromkeys(pkgs))


def main() -> int:
    q = yaml.safe_load(QUEUE.read_text())
    n = 0
    for t in q["tokens"]:
        if "mutation_verified" not in t:
            t["mutation_verified"] = False
            n += 1
        t.setdefault("mutation_score", None)
        t.setdefault("mutation_pkg", "")
        pkgs = pkg_for_impl(t["token"])
        if pkgs and not t.get("mutation_pkg"):
            t["mutation_pkg"] = pkgs[0]
    QUEUE.write_text(yaml.dump(q, default_flow_style=False, sort_keys=False, allow_unicode=True, width=120))
    print(f"DEBUG: seeded mutation_verified/mutation_score on {len(q['tokens'])} tokens ({n} new flags)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
