// [REQ-CONFIGURATION]
package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestConfigOutputGrouping(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		wantContains   []string
		wantNotContain []string
	}{
		{
			name: "default grouped output",
			args: []string{"config"},
			wantContains: []string{
				"## Basic Settings",
				"## Archive Settings",
				"Configuration Sources:",
			},
			wantNotContain: []string{
				"config: ", // Flat output header
			},
		},
		{
			name: "flat output flag",
			args: []string{"config", "--flat"},
			wantContains: []string{
				"config: ",
				"archive_dir_path:",
			},
			wantNotContain: []string{
				"## Basic Settings",
			},
		},
		{
			name: "format flat",
			args: []string{"config", "--format", "flat"},
			wantContains: []string{
				"config: ",
				"archive_dir_path:",
			},
			wantNotContain: []string{
				"## Basic Settings",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			tmpDir := t.TempDir()
			origWd, _ := os.Getwd()
			os.Chdir(tmpDir)
			defer os.Chdir(origWd)

			// Create dummy config file
			cfgPath := ".bkpdir.yml" // Use relative path in tmpDir
			os.WriteFile(cfgPath, []byte("archive_dir_path: ./archives\n"), 0644)

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			rootCmd := createTestRootCmd()
			// We don't need SetOut/SetErr for cobra if we capture os.Stdout/os.Stderr,
			// but let's keep them to avoid double printing if cobra uses them.
			rootCmd.SetOut(w)
			rootCmd.SetErr(w)
			rootCmd.SetArgs(tt.args)

			var buf bytes.Buffer
			done := make(chan struct{})
			go func() {
				_, _ = buf.ReadFrom(r)
				done <- struct{}{}
			}()

			err := rootCmd.Execute()

			// Close writer and restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Wait for read to finish
			<-done

			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			output := buf.String()

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Output missing expected string %q", want)
				}
			}
			for _, notWant := range tt.wantNotContain {
				if strings.Contains(output, notWant) {
					t.Errorf("Output contains unexpected string %q", notWant)
				}
			}
		})
	}
}
