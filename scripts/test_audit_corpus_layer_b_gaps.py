#!/usr/bin/env python3
"""Tests for audit_corpus_layer_b_gaps.py."""
from __future__ import annotations

import json
import tempfile
import unittest
from pathlib import Path
import sys

sys.path.insert(0, str(Path(__file__).resolve().parent))
from audit_corpus_layer_b_gaps import analyze_sidecar, main

ROOT = Path(__file__).resolve().parent.parent
TIED = ROOT / "tied-project"


class AuditCorpusLayerBGapsTest(unittest.TestCase):
    def test_analyze_exclusion_patterns_has_v2(self):
        m = analyze_sidecar("IMPL-EXCLUSION_PATTERNS", TIED)
        self.assertIsNotNone(m.path)
        self.assertEqual(m.grammar_version, "v2")

    def test_main_writes_report(self):
        with tempfile.TemporaryDirectory() as tmp:
            out = Path(tmp) / "report.json"
            rc = main(["--out", str(out), "--pilot-only"])
            self.assertEqual(rc, 0)
            body = json.loads(out.read_text())
            self.assertEqual(body["schema"], "corpus-layer-b-gap-report.v1")
            self.assertGreater(len(body["pilot_metrics"]), 0)


if __name__ == "__main__":
    unittest.main()
