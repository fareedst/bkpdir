#!/usr/bin/env python3
"""Tests for layer_b_ci_enforcement.py. [REQ-PSEUDOCODE_FORMAL_VERIFICATION]"""
from __future__ import annotations

import sys
import tempfile
import unittest
from pathlib import Path

import yaml

sys.path.insert(0, str(Path(__file__).resolve().parent))
from layer_b_ci_enforcement import (
    select_enforcement_tokens,
    token_from_sidecar_path,
    wave_for_token,
)


class LayerBCIEnforcementTest(unittest.TestCase):
    def test_token_from_sidecar_path(self):
        self.assertEqual(
            token_from_sidecar_path(
                "tied-project/implementation-decisions/IMPL-FOO-pseudocode.md"
            ),
            "IMPL-FOO",
        )
        self.assertIsNone(token_from_sidecar_path("pkg/foo.go"))

    def test_if_changed_respects_enforce_in_ci(self):
        reg = {
            "waves": {
                "w0": {"enforce_in_ci": True, "tokens": ["IMPL-A"]},
                "w1": {"enforce_in_ci": False, "tokens": ["IMPL-B"]},
            }
        }
        self.assertEqual(
            select_enforcement_tokens(reg, mode="if_changed", changed=["IMPL-A", "IMPL-B"]),
            ["IMPL-A"],
        )

    def test_wave_mode_returns_all_wave_tokens(self):
        reg = {"waves": {"w2": {"enforce_in_ci": False, "tokens": ["IMPL-X", "IMPL-Y"]}}}
        self.assertEqual(
            select_enforcement_tokens(reg, mode="wave", wave_id="w2"),
            ["IMPL-X", "IMPL-Y"],
        )

    def test_wave_for_token(self):
        reg = {"waves": {"w3": {"enforce_in_ci": "optional", "tokens": ["IMPL-Z"]}}}
        wave_id, enforce = wave_for_token(reg, "IMPL-Z")
        self.assertEqual(wave_id, "w3")
        self.assertFalse(enforce)

    def test_registry_fixture_roundtrip(self):
        with tempfile.TemporaryDirectory() as tmp:
            path = Path(tmp) / "registry.yaml"
            path.write_text(
                yaml.safe_dump(
                    {
                        "schema": "layer-b-wave-registry.v1",
                        "waves": {
                            "w0": {
                                "enforce_in_ci": True,
                                "tokens": ["IMPL-ATOMIC_OPS"],
                            }
                        },
                    }
                )
            )
            data = yaml.safe_load(path.read_text())
            self.assertEqual(
                select_enforcement_tokens(
                    data, mode="if_changed", changed=["IMPL-ATOMIC_OPS"]
                ),
                ["IMPL-ATOMIC_OPS"],
            )


if __name__ == "__main__":
    unittest.main()
