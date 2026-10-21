#!/usr/bin/env python3
"""Read-only Layer B corpus gap audit (Phase 0). [REQ-PSEUDOCODE_FORMAL_VERIFICATION]"""
from __future__ import annotations

import json
import re
import subprocess
import sys
from dataclasses import asdict, dataclass
from pathlib import Path

from formal_spec_registry import iter_tokens, registry_path

ROOT = Path(__file__).resolve().parent.parent
DEFAULT_OUT = (
    ROOT
    / "tied-project/working/REQ-PSEUDOCODE_LAYER_B_CORPUS/evidence/corpus-layer-b-gap-report.json"
)

PILOT_W0 = [
    "IMPL-ATOMIC_OPS",
    "IMPL-CONFIG_STRUCT",
    "IMPL-CFG_006",
    "IMPL-FILE_OPERATIONS",
    "IMPL-ZIP_FORMAT",
    "IMPL-EXCLUSION_PATTERNS",
]
PILOT_EXTRA = [
    "IMPL-DIFF_COMMAND",
    "IMPL-CLI_FRAMEWORK",
    "IMPL-DATA_MODELS",
    "IMPL-DOC_ENHANCEMENT",
    "IMPL-EXTRACTION_PRINCIPLES",
]

_RE_GRAMMAR = re.compile(r"^Grammar-Version:\s*(\S+)", re.M)
_RE_PROC = re.compile(r"^procedure\s+([A-Z0-9_]+)\s*:", re.M)
_RE_PRE_TRUE = re.compile(r"^\s*PRE:\s*true\s*$", re.M)
_RE_POST_TRUE = re.compile(r"^\s*POST:\s*true\s*$", re.M)
_RE_ERROR_STEP = re.compile(r"^\s*ERROR\s+([A-Z0-9_]+)", re.M)
_RE_ASYNC = re.compile(r"\b(AWAIT|Async|SEND)\b")


def sidecar_path(token: str, tied_base: Path) -> Path | None:
    rel = tied_base / f"implementation-decisions/{token}-pseudocode.md"
    if rel.is_file():
        return rel
    alt = ROOT / "tied" / f"implementation-decisions/{token}-pseudocode.md"
    return alt if alt.is_file() else None


@dataclass
class SidecarMetrics:
    token: str
    path: str | None
    grammar_version: str | None
    active_procedures: int
    pre_true_count: int
    post_true_count: int
    error_step_names: list[str]
    async_in_scope: bool


def analyze_sidecar(token: str, tied_base: Path) -> SidecarMetrics:
    path = sidecar_path(token, tied_base)
    if path is None:
        return SidecarMetrics(token, None, None, 0, 0, 0, [], False)
    text = path.read_text(encoding="utf-8")
    gv = _RE_GRAMMAR.search(text)
    errors = sorted(set(_RE_ERROR_STEP.findall(text)))
    try:
        path_display = str(path.resolve().relative_to(ROOT.resolve()))
    except ValueError:
        path_display = str(path)
    return SidecarMetrics(
        token=token,
        path=path_display,
        grammar_version=gv.group(1) if gv else None,
        active_procedures=len(_RE_PROC.findall(text)),
        pre_true_count=len(_RE_PRE_TRUE.findall(text)),
        post_true_count=len(_RE_POST_TRUE.findall(text)),
        error_step_names=errors,
        async_in_scope=bool(_RE_ASYNC.search(text)),
    )


def run_specctl_validate(sidecar: Path) -> dict:
    proc = subprocess.run(
        ["go", "run", "./cmd/specctl", "validate", str(sidecar)],
        cwd=ROOT,
        capture_output=True,
        text=True,
    )
    return {"exit_code": proc.returncode, "stdout": proc.stdout[-2000:], "stderr": proc.stderr[-2000:]}


def main(argv: list[str] | None = None) -> int:
    args = argv if argv is not None else sys.argv[1:]
    out = DEFAULT_OUT
    pilot_only = False
    if "--out" in args:
        out = Path(args[args.index("--out") + 1])
    if "--pilot-only" in args:
        pilot_only = True

    env_base = __import__("os").environ.get("TIED_BASE_PATH")
    default_base = ROOT / "tied-project"
    if env_base:
        candidate = Path(env_base)
        tied_base = (
            candidate
            if (candidate / "implementation-decisions").is_dir()
            else default_base
        )
    else:
        tied_base = default_base
    reg = registry_path(ROOT / "tied-project/spec/formal-spec-registry.yaml")
    tokens = iter_tokens(reg)
    if pilot_only:
        tokens = [t for t in PILOT_W0 + PILOT_EXTRA if t in tokens]

    metrics = [asdict(analyze_sidecar(t, tied_base)) for t in tokens]
    pilot_set = set(PILOT_W0 + PILOT_EXTRA)
    pilot_metrics = [m for m in metrics if m["token"] in pilot_set]

    report = {
        "schema": "corpus-layer-b-gap-report.v1",
        "corpus_token_count": len(iter_tokens(reg)),
        "metrics": metrics,
        "pilot_tokens": sorted(pilot_set),
        "pilot_metrics": pilot_metrics,
        "aggregates": {
            "grammar_v2_count": sum(1 for m in metrics if m.get("grammar_version") == "v2"),
            "pre_true_total": sum(m.get("pre_true_count", 0) for m in metrics),
            "post_true_total": sum(m.get("post_true_count", 0) for m in metrics),
            "async_in_scope_count": sum(1 for m in metrics if m.get("async_in_scope")),
        },
    }
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps(report, indent=2) + "\n", encoding="utf-8")
    try:
        shown = out.relative_to(ROOT)
    except ValueError:
        shown = out
    print(f"DIAGNOSTIC: wrote {shown} ({len(metrics)} IMPL metrics)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
