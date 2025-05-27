// This file is part of bkpdir

// Copyright 2024 BkpDir Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main provides tests for archive creation and management in BkpDir.
// It verifies archive creation, naming, and verification behavior.
package main

import "testing"

func TestGenerateArchiveName(t *testing.T) {
	tests := []struct {
		cfg    ArchiveNameConfig
		expect string
		name   string
	}{
		{
			cfg: ArchiveNameConfig{
				Prefix:        "HOME",
				Timestamp:     "2024-06-01-12-00",
				GitBranch:     "main",
				GitHash:       "abc123",
				Note:          "",
				IsGit:         true,
				IsIncremental: false,
				BaseName:      "",
			},
			expect: "HOME-2024-06-01-12-00=main=abc123.zip",
			name:   "full with git",
		},
		{
			cfg: ArchiveNameConfig{
				Prefix:        "",
				Timestamp:     "2024-06-01-12-00",
				GitBranch:     "",
				GitHash:       "",
				Note:          "note",
				IsGit:         false,
				IsIncremental: false,
				BaseName:      "",
			},
			expect: "2024-06-01-12-00=note.zip",
			name:   "full with note only",
		},
		{
			cfg: ArchiveNameConfig{
				Prefix:        "",
				Timestamp:     "2024-06-01-12-00",
				GitBranch:     "dev",
				GitHash:       "def456",
				Note:          "backup",
				IsGit:         true,
				IsIncremental: true,
				BaseName:      "HOME-2024-06-01-12-00=main=abc123",
			},
			expect: "HOME-2024-06-01-12-00=main=abc123_update=2024-06-01-12-00=dev=def456=backup.zip",
			name:   "incremental with git and note",
		},
		{
			cfg: ArchiveNameConfig{
				Prefix:        "",
				Timestamp:     "2024-06-01-12-00",
				GitBranch:     "",
				GitHash:       "",
				Note:          "",
				IsGit:         false,
				IsIncremental: true,
				BaseName:      "ARCHIVE-2024-06-01-12-00",
			},
			expect: "ARCHIVE-2024-06-01-12-00_update=2024-06-01-12-00.zip",
			name:   "incremental no git, no note",
		},
		{
			cfg: ArchiveNameConfig{
				Prefix:        "",
				Timestamp:     "2024-06-01-12-00",
				GitBranch:     "",
				GitHash:       "",
				Note:          "",
				IsGit:         false,
				IsIncremental: false,
				BaseName:      "",
			},
			expect: "2024-06-01-12-00.zip",
			name:   "full no git, no note",
		},
		{
			cfg: ArchiveNameConfig{
				Prefix:        "",
				Timestamp:     "2024-06-01-12-00",
				GitBranch:     "",
				GitHash:       "",
				Note:          "",
				IsGit:         true,
				IsIncremental: false,
				BaseName:      "",
			},
			expect: "2024-06-01-12-00.zip",
			name:   "full git, no branch/hash",
		},
		{
			cfg: ArchiveNameConfig{
				Prefix:        "",
				Timestamp:     "2024-06-01-12-00",
				GitBranch:     "",
				GitHash:       "",
				Note:          "test note",
				IsGit:         false,
				IsIncremental: false,
				BaseName:      "",
			},
			expect: "2024-06-01-12-00=test note.zip",
			name:   "full with note containing spaces",
		},
		{
			cfg: ArchiveNameConfig{
				Prefix:        "",
				Timestamp:     "2024-06-01-12-00",
				GitBranch:     "",
				GitHash:       "",
				Note:          "test-note",
				IsGit:         false,
				IsIncremental: true,
				BaseName:      "ARCHIVE-2024-06-01-12-00",
			},
			expect: "ARCHIVE-2024-06-01-12-00_update=2024-06-01-12-00=test-note.zip",
			name:   "incremental with note only",
		},
		{
			cfg: ArchiveNameConfig{
				Prefix:        "",
				Timestamp:     "2024-06-01-12-00",
				GitBranch:     "main",
				GitHash:       "abc123",
				Note:          "test note",
				IsGit:         true,
				IsIncremental: false,
				BaseName:      "",
			},
			expect: "2024-06-01-12-00=main=abc123=test note.zip",
			name:   "full with git and note with spaces",
		},
	}
	for _, tt := range tests {
		got := GenerateArchiveName(tt.cfg)
		if got != tt.expect {
			t.Errorf("%s: GenerateArchiveName() = %q, want %q", tt.name, got, tt.expect)
		}
	}
}
