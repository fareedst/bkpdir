#!/usr/bin/env python3
"""Shared logic for package-level IMPL block-lead audit and strip. [REQ-PSEUDOCODE_FORMAL_VERIFICATION]"""
from __future__ import annotations

import json
import re
from dataclasses import asdict, dataclass
from pathlib import Path

import yaml

REPO = Path(__file__).resolve().parents[1]
IMPL_DIR = REPO / "tied/implementation-decisions"

BLOCK_LEAD_RE = re.compile(r"^\s*// - (\[IMPL-[^\]]+\].*)$")
# D15: paraphrased file-level banners (not literal block leads)
BANNER_IMPL_RE = re.compile(r"^\s*// (?!\- )\[IMPL-([A-Z0-9_]+)\]")
PACKAGE_RE = re.compile(r"^package\s+(\w+)")
FUNC_TEST_RE = re.compile(r"^func\s+(?:Test|Benchmark)\w*\(")
IMPORT_RE = re.compile(r"^import\s+[\(\w]")

FINGERPRINT_SUBSTR = "[IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] — How: validate cwd"

SUSPECT_MIN_PACKAGE_LEADS = 30


@dataclass
class FileAudit:
    path: str
    package_level_lead_count: int
    suspect: bool
    in_check_leads_pool: bool
    legitimate_prefix_count: int
    action: str


@dataclass
class BannerAudit:
    path: str
    redundant_banner_count: int
    tokens: list[str]


def _impl_token_from_lead(lead_body: str) -> str | None:
    m = re.search(r"\[IMPL-([A-Z0-9_]+)\]", lead_body)
    return f"IMPL-{m.group(1)}" if m else None


def load_impl_code_files() -> dict[str, set[str]]:
    out: dict[str, set[str]] = {}
    for ypath in IMPL_DIR.glob("IMPL-*.yaml"):
        if ypath.name.endswith("-pseudocode.md"):
            continue
        try:
            data = yaml.safe_load(ypath.read_text())
        except Exception:
            continue
        if not isinstance(data, dict):
            continue
        for token, rec in data.items():
            if not token.startswith("IMPL-"):
                continue
            files: set[str] = set()
            loc = rec.get("code_locations") or {}
            for f in loc.get("files") or []:
                if isinstance(f, str):
                    files.add(f.replace("\\", "/"))
                elif isinstance(f, dict):
                    p = f.get("path") or f.get("file")
                    if p:
                        files.add(str(p).replace("\\", "/"))
            for fn in loc.get("functions") or []:
                if isinstance(fn, dict) and fn.get("file"):
                    files.add(str(fn["file"]).replace("\\", "/"))
            out[token] = files
    return out


def in_check_leads_pool(rel: str) -> bool:
    p = rel.replace("\\", "/")
    if p.startswith("pkg/"):
        return True
    if "/" not in p and p.endswith(".go"):
        return True
    return False


def lead_belongs_in_test_package(token: str, test_rel: str, code_map: dict[str, set[str]]) -> bool:
    files = code_map.get(token, set())
    if not files:
        return False
    test_path = Path(test_rel)
    test_dir = test_path.parent
    test_dir_s = "." if str(test_dir) in ("", ".") else str(test_dir).replace("\\", "/")
    for cf in files:
        cf_path = Path(cf)
        code_dir = cf_path.parent
        code_dir_s = "." if str(code_dir) in ("", ".") else str(code_dir).replace("\\", "/")
        if code_dir_s == test_dir_s:
            return True
    return False


def split_go_file(content: str) -> tuple[list[str], list[str], list[str]]:
    """Return (prefix_before_package, package_region_lines, rest). package_region is after package until import/func/non-lead."""
    lines = content.splitlines(keepends=True)
    pkg_idx = None
    for i, line in enumerate(lines):
        if PACKAGE_RE.match(line.strip()):
            pkg_idx = i
            break
    if pkg_idx is None:
        return lines, [], []

    region: list[str] = []
    rest_start = pkg_idx + 1
    i = pkg_idx + 1
    while i < len(lines):
        stripped = lines[i].strip()
        if stripped == "":
            i += 1
            continue
        if BLOCK_LEAD_RE.match(lines[i]):
            region.append(lines[i])
            i += 1
            continue
        break
    rest_start = i
    prefix = lines[: pkg_idx + 1]
    rest = lines[rest_start:]
    return prefix, region, rest


def audit_file(rel: str, content: str, code_map: dict[str, set[str]]) -> FileAudit:
    _, region, _ = split_go_file(content)
    legit = 0
    for line in region:
        m = BLOCK_LEAD_RE.match(line)
        if not m:
            continue
        tok = _impl_token_from_lead(m.group(1))
        if tok and lead_belongs_in_test_package(tok, rel, code_map):
            legit += 1
    count = len(region)
    has_fp = any(FINGERPRINT_SUBSTR in ln for ln in region)
    suspect = count >= SUSPECT_MIN_PACKAGE_LEADS or has_fp
    if not region:
        action = "keep"
    elif rel.endswith("_test.go") and count > 0:
        action = "strip_all"
    elif not rel.endswith("_test.go") and count > 0 and suspect:
        action = "review"
    elif legit == count:
        action = "keep"
    elif legit > 0:
        action = "strip_aggregate_keep_prefix"
    else:
        action = "strip_all"
    return FileAudit(
        path=rel,
        package_level_lead_count=count,
        suspect=suspect,
        in_check_leads_pool=in_check_leads_pool(rel),
        legitimate_prefix_count=legit,
        action=action,
    )


def block_lead_impl_tokens(content: str) -> set[str]:
    tokens: set[str] = set()
    for line in content.splitlines():
        m = BLOCK_LEAD_RE.match(line)
        if not m:
            continue
        tok = _impl_token_from_lead(m.group(1))
        if tok:
            tokens.add(tok.removeprefix("IMPL-"))
    return tokens


def audit_redundant_banners(rel: str, content: str) -> BannerAudit:
    """D15: banners whose IMPL token also appears on a literal // - block lead in the same file."""
    block_tokens = block_lead_impl_tokens(content)
    redundant_tokens: set[str] = set()
    count = 0
    for line in content.splitlines():
        m = BANNER_IMPL_RE.match(line)
        if m and m.group(1) in block_tokens:
            redundant_tokens.add(m.group(1))
            count += 1
    return BannerAudit(
        path=rel,
        redundant_banner_count=count,
        tokens=sorted(redundant_tokens),
    )


def strip_redundant_banners(content: str) -> tuple[str, int]:
    block_tokens = block_lead_impl_tokens(content)
    if not block_tokens:
        return content, 0
    out: list[str] = []
    removed = 0
    for line in content.splitlines(keepends=True):
        m = BANNER_IMPL_RE.match(line)
        if m and m.group(1) in block_tokens:
            removed += 1
            continue
        out.append(line)
    return "".join(out), removed


def strip_package_leads(content: str, test_rel: str, code_map: dict[str, set[str]]) -> tuple[str, int]:
    prefix, region, rest = split_go_file(content)
    if not region:
        return content, 0

    # Test files: package-level block leads are never valid (D13 = before Test/Benchmark only).
    if test_rel.endswith("_test.go"):
        removed = len(region)
        return "".join(prefix + rest), removed

    kept: list[str] = []
    removed = 0
    for line in region:
        m = BLOCK_LEAD_RE.match(line)
        if not m:
            kept.append(line)
            continue
        tok = _impl_token_from_lead(m.group(1))
        if tok and lead_belongs_in_test_package(tok, test_rel, code_map):
            kept.append(line)
        else:
            removed += 1
    if not removed:
        return content, 0
    return "".join(prefix + kept + rest), removed


def iter_go_files(wave: str | None) -> list[Path]:
    waves = {
        "a": lambda p: p.name.endswith("_test.go") and len(p.relative_to(REPO).parts) == 1,
        "b": lambda p: "pkg" in p.parts and p.name.endswith("_test.go"),
        "c": lambda p: p.name.endswith("_test.go")
        and len(p.relative_to(REPO).parts) > 0
        and p.relative_to(REPO).parts[0] in ("internal", "test", "cmd", "tools"),
    }
    pred = waves.get(wave or "", lambda p: p.suffix == ".go")
    out: list[Path] = []
    for path in REPO.rglob("*.go"):
        rel = path.relative_to(REPO)
        if "vendor" in rel.parts or ".git" in rel.parts:
            continue
        if pred(path):
            out.append(path)
    return sorted(out)


def run_banner_audit(output_json: Path | None = None) -> list[BannerAudit]:
    results: list[BannerAudit] = []
    for path in REPO.rglob("*.go"):
        rel = str(path.relative_to(REPO)).replace("\\", "/")
        if "vendor" in rel.split("/") or rel.startswith("scripts/testdata/"):
            continue
        try:
            content = path.read_text(errors="replace")
        except OSError:
            continue
        audit = audit_redundant_banners(rel, content)
        if audit.redundant_banner_count:
            results.append(audit)
    results.sort(key=lambda r: (-r.redundant_banner_count, r.path))
    if output_json:
        output_json.parent.mkdir(parents=True, exist_ok=True)
        output_json.write_text(json.dumps([asdict(r) for r in results], indent=2) + "\n")
    return results


def run_audit(output_json: Path | None = None) -> list[FileAudit]:
    code_map = load_impl_code_files()
    results: list[FileAudit] = []
    for path in REPO.rglob("*.go"):
        rel = str(path.relative_to(REPO)).replace("\\", "/")
        if "vendor" in rel.split("/") or rel.startswith("scripts/testdata/"):
            continue
        try:
            content = path.read_text(errors="replace")
        except OSError:
            continue
        results.append(audit_file(rel, content, code_map))
    results.sort(key=lambda r: (-int(r.suspect), -r.package_level_lead_count, r.path))
    if output_json:
        output_json.parent.mkdir(parents=True, exist_ok=True)
        output_json.write_text(json.dumps([asdict(r) for r in results], indent=2) + "\n")
    return results
