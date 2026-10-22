#!/usr/bin/env python3
"""RED/GREEN tests for run-pseudocode-layer-b-batch.sh. [REQ-PSEUDOCODE_FORMAL_VERIFICATION]"""
from __future__ import annotations

import json
import os
import subprocess
import tempfile
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
SCRIPT = ROOT / "scripts/run-pseudocode-layer-b-batch.sh"
FAKE_CLI = ROOT / "scripts/testdata/layer_b_batch/fake_tied_cli_fail.sh"
FIXTURE_SIDE = ROOT / "scripts/testdata/layer_b_batch/minimal-sidecar.md"


class LayerBBatchScriptTest(unittest.TestCase):
    def _run(self, *extra: str, env: dict[str, str] | None = None) -> subprocess.CompletedProcess[str]:
        merged = os.environ.copy()
        merged["TIED_BASE_PATH"] = str(ROOT / "scripts/testdata/layer_b_batch/tied-project")
        merged["TIED_CLI"] = str(FAKE_CLI)
        if env:
            merged.update(env)
        return subprocess.run(
            ["bash", str(SCRIPT), *extra],
            cwd=ROOT,
            env=merged,
            text=True,
            capture_output=True,
        )

    def test_enforce_exits_nonzero_when_validate_errors(self):
        self.assertTrue(FAKE_CLI.exists(), "fake tied-cli fixture missing")
        with tempfile.TemporaryDirectory() as tmp:
            out = Path(tmp)
            proc = self._run(
                "--mode",
                "structural",
                "--out-dir",
                str(out),
                "--token",
                "IMPL-OK",
                "--enforce",
                "IMPL-FAIL",
                "--token",
                "IMPL-FAIL",
            )
            self.assertNotEqual(proc.returncode, 0, proc.stdout + proc.stderr)
            report = out / "IMPL-FAIL.json"
            self.assertTrue(report.is_file(), proc.stderr)
            body = json.loads(report.read_text())
            self.assertFalse(body.get("ok", True))

    def test_without_enforce_continues_on_errors(self):
        with tempfile.TemporaryDirectory() as tmp:
            out = Path(tmp)
            proc = self._run(
                "--mode",
                "structural",
                "--out-dir",
                str(out),
                "--token",
                "IMPL-FAIL",
            )
            self.assertEqual(proc.returncode, 0, proc.stderr)
            self.assertTrue((out / "IMPL-FAIL.json").is_file())

    def test_missing_tokens_is_usage_error(self):
        proc = self._run("--mode", "structural")
        self.assertEqual(proc.returncode, 2, proc.stdout + proc.stderr)


if __name__ == "__main__":
    unittest.main()
