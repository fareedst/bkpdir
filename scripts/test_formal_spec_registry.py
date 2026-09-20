#!/usr/bin/env python3
"""Tests for formal_spec_registry parser. [REQ-PSEUDOCODE_FORMAL_VERIFICATION]"""
from __future__ import annotations

import unittest
from pathlib import Path

from formal_spec_registry import count_tokens, iter_tokens

FIXTURE = Path(__file__).resolve().parent / "testdata/formal_spec_registry/quoted_and_unquoted.yaml"
REPO_REGISTRY = Path(__file__).resolve().parent.parent / "tied/spec/formal-spec-registry.yaml"


class FormalSpecRegistryTest(unittest.TestCase):
    def test_fixture_quoted_and_unquoted(self):
        self.assertEqual(
            iter_tokens(FIXTURE),
            ["IMPL-ALPHA", "IMPL-BETA", "IMPL-GAMMA", "IMPL-DELTA"],
        )
        self.assertEqual(count_tokens(FIXTURE), 4)

    def test_repo_registry_has_73_tokens(self):
        n = count_tokens(REPO_REGISTRY)
        self.assertEqual(n, 73)
        tokens = iter_tokens(REPO_REGISTRY)
        self.assertTrue(all(t.startswith("IMPL-") for t in tokens))
        self.assertEqual(len(tokens), len(set(tokens)))


if __name__ == "__main__":
    unittest.main()
