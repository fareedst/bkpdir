#!/usr/bin/env python3
"""Select IMPL tokens for Layer B CI enforcement from git diff + wave registry.

[REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-SPEC_CTL]
"""
from __future__ import annotations

import argparse
import json
import re
import subprocess
import sys
from pathlib import Path

import yaml

REPO = Path(__file__).resolve().parents[1]
DEFAULT_REGISTRY = REPO / "tied-project/spec/layer-b-wave-registry.yaml"
SIDECAR_RE = re.compile(
    r"(?:^|/)implementation-decisions/(IMPL-[A-Z0-9_]+)-pseudocode\.md$"
)


def load_wave_registry(path: Path) -> dict:
    data = yaml.safe_load(path.read_text())
    if not isinstance(data, dict):
        raise ValueError(f"invalid registry: {path}")
    return data


def token_from_sidecar_path(path: str) -> str | None:
    norm = path.replace("\\", "/")
    m = SIDECAR_RE.search(norm)
    return m.group(1) if m else None


def git_changed_paths(base_ref: str, head_ref: str = "HEAD") -> list[str]:
    proc = subprocess.run(
        ["git", "diff", "--name-only", f"{base_ref}...{head_ref}"],
        cwd=REPO,
        text=True,
        capture_output=True,
    )
    if proc.returncode != 0:
        proc = subprocess.run(
            ["git", "diff", "--name-only", base_ref, head_ref],
            cwd=REPO,
            text=True,
            capture_output=True,
        )
    if proc.returncode != 0:
        return []
    return [ln.strip() for ln in proc.stdout.splitlines() if ln.strip()]


def changed_impl_tokens(base_ref: str, head_ref: str = "HEAD") -> list[str]:
    out: list[str] = []
    seen: set[str] = set()
    for path in git_changed_paths(base_ref, head_ref):
        tok = token_from_sidecar_path(path)
        if tok and tok not in seen:
            seen.add(tok)
            out.append(tok)
    return sorted(out)


def wave_for_token(registry: dict, token: str) -> tuple[str | None, bool]:
    waves = registry.get("waves") or {}
    for wave_id, wave in waves.items():
        tokens = wave.get("tokens") or []
        if token in tokens:
            flag = wave.get("enforce_in_ci", False)
            enforce = flag is True or str(flag).lower() == "true"
            return wave_id, enforce
    return None, False


def tokens_for_wave(registry: dict, wave_id: str) -> list[str]:
    wave = (registry.get("waves") or {}).get(wave_id)
    if not wave:
        raise KeyError(f"unknown wave {wave_id!r}")
    return list(wave.get("tokens") or [])


def select_enforcement_tokens(
    registry: dict,
    *,
    mode: str,
    wave_id: str | None = None,
    changed: list[str] | None = None,
) -> list[str]:
    """Return sorted unique IMPL tokens to enforce with pseudocode_validate full mode."""
    if mode == "wave":
        if not wave_id:
            raise ValueError("wave mode requires wave_id")
        return sorted(tokens_for_wave(registry, wave_id))

    if mode == "force_changed":
        return sorted(changed or [])

    if mode == "if_changed":
        selected: list[str] = []
        for token in changed or []:
            _, enforce = wave_for_token(registry, token)
            if enforce:
                selected.append(token)
        return sorted(set(selected))

    raise ValueError(f"unknown mode {mode!r}")


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--registry", type=Path, default=DEFAULT_REGISTRY)
    parser.add_argument("--base-ref", default="HEAD")
    parser.add_argument("--head-ref", default="HEAD")
    parser.add_argument(
        "--mode",
        choices=("if_changed", "wave", "force_changed"),
        default="if_changed",
    )
    parser.add_argument("--wave", default="")
    parser.add_argument("--dry-run", action="store_true")
    parser.add_argument("--json", action="store_true")
    args = parser.parse_args()

    if not args.registry.is_file():
        print(f"DIAGNOSTIC: missing registry {args.registry}", file=sys.stderr)
        return 2

    registry = load_wave_registry(args.registry)
    changed = changed_impl_tokens(args.base_ref, args.head_ref)
    wave_id = args.wave.strip().lower() or None
    if wave_id and not wave_id.startswith("w"):
        wave_id = f"w{wave_id}"

    mode = args.mode
    if mode == "if_changed" and wave_id:
        mode = "wave"

    tokens = select_enforcement_tokens(
        registry,
        mode=mode,
        wave_id=wave_id,
        changed=changed,
    )

    payload = {
        "mode": mode,
        "wave": wave_id,
        "changed_sidecar_tokens": changed,
        "enforce_tokens": tokens,
        "registry": str(args.registry.relative_to(REPO)),
    }
    if args.json or args.dry_run:
        print(json.dumps(payload, indent=2))
    else:
        for t in tokens:
            print(t)

    if args.dry_run:
        return 0
    return 0 if tokens else 0


if __name__ == "__main__":
    raise SystemExit(main())
