#!/usr/bin/env python3
"""Unit tests for lead_hygiene audit/strip. [REQ-PSEUDOCODE_FORMAL_VERIFICATION]"""
from __future__ import annotations

import unittest
from pathlib import Path

from lead_hygiene import (
    audit_file,
    audit_redundant_banners,
    load_impl_code_files,
    strip_package_leads,
    strip_redundant_banners,
)

FIXTURES = Path(__file__).resolve().parent / "testdata/lead_hygiene"


class LeadHygieneTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.code_map = load_impl_code_files()

    def test_strip_removes_all_package_leads_from_test_files(self):
        content = (FIXTURES / "pkg_config_prefix.go").read_text()
        new_content, removed = strip_package_leads(
            content, "pkg/config/config_test.go", self.code_map
        )
        self.assertEqual(removed, 2)
        self.assertNotIn("IMPL-CFG_QUOTED_KEY_PREFIX", new_content)
        self.assertNotIn("IMPL-ZIP_FORMAT", new_content)
        self.assertIn("func TestExample", new_content)

    def test_strip_redundant_banners_keeps_block_leads(self):
        content = (
            "// [IMPL-SAMPLE] paraphrased banner\n"
            "package p\n\n"
            "// - [IMPL-SAMPLE] [ARCH-X] [REQ-Y] — How: literal block lead.\n"
            "func F() {}\n"
        )
        new_content, removed = strip_redundant_banners(content)
        self.assertEqual(removed, 1)
        self.assertNotIn("paraphrased banner", new_content)
        self.assertIn("// - [IMPL-SAMPLE]", new_content)

    def test_audit_redundant_banners(self):
        content = (
            "// [IMPL-SAMPLE] banner\n"
            "// - [IMPL-SAMPLE] [ARCH-X] — How: block.\n"
        )
        audit = audit_redundant_banners("sample.go", content)
        self.assertEqual(audit.redundant_banner_count, 1)
        self.assertEqual(audit.tokens, ["SAMPLE"])

    def test_audit_marks_mega_block_suspect(self):
        content = (FIXTURES / "pkg_config_prefix.go").read_text()
        pad = (
            "// - [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] — How: validate cwd\n"
            * 35
        )
        content = content.replace("import", pad + "import", 1)
        audit = audit_file("config_test.go", content, self.code_map)
        self.assertTrue(audit.suspect)
        self.assertEqual(audit.action, "strip_all")


if __name__ == "__main__":
    unittest.main()
