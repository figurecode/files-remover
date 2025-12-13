package remover

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	t.Run("Remove real files", func(t *testing.T) {
		tmpDir := t.TempDir()

		files := map[string]int64{
			filepath.Join(tmpDir, "info-1.log"): 100,
			filepath.Join(tmpDir, "info-2.log"): 200,
		}

		for path := range files {
			if err := os.WriteFile(path, []byte("important"), 0644); err != nil {
				t.Fatal(err)
			}
		}

		if err := Execute(files); err != nil {
			t.Fatalf("Execute() return error: %v", err)
		}

		for path := range files {
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				t.Errorf("File %q was not deleted", path)
			}
		}
	})

	t.Run("File already missing", func(t *testing.T) {
		files := map[string]int64{
			"/tmp/file/does-not/exist/really.log": 12345,
		}

		if err := Execute(files); err != nil {
			t.Fatalf("Execute() returned error on missing file: %v", err)
		}
	})

	t.Run("Empty files map", func(t *testing.T) {
		files := map[string]int64{}

		if err := Execute(files); err != nil {
			t.Fatalf("Execute() on empty map returned error: %v", err)
		}
	})
}

func TestDebugRemover(t *testing.T) {
	t.Run("Does not remove files", func(t *testing.T) {
		tmpDir := t.TempDir()

		files := map[string]int64{
			filepath.Join(tmpDir, "info-1.log"): 100,
			filepath.Join(tmpDir, "info-2.log"): 200,
		}

		for path := range files {
			if err := os.WriteFile(path, []byte("important"), 0644); err != nil {
				t.Fatal(err)
			}
		}

		var buf bytes.Buffer
		if err := DebugRemover(files, &buf); err != nil {
			t.Fatal(err)
		}

		for path := range files {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Fatalf("DebugRemover() deleted file %s", path)
			}
		}
	})

	t.Run("Output files and size", func(t *testing.T) {
		files := map[string]int64{
			"/tmp/log/fake/info1.log": 1024,
			"/tmp/log/fake/info2.log": 2048,
		}

		var buf bytes.Buffer
		if err := DebugRemover(files, &buf); err != nil {
			t.Fatalf("DebugRemover() error %v", err)
		}

		got := buf.String()
		want := []string{
			"2 files will be deleted in total",
			"3.0 KB of disk space will be freed",
			"Files to be deleted:",
			"---------------------------------",
			"PATH: /tmp/log/fake/info1.log",
			"---------------------------------",
			"PATH: /tmp/log/fake/info2.log",
			"",
			"END",
		}

		for _, w := range want {
			if !strings.Contains(got, w) {
				t.Errorf("output missing %q\nfull output:\n%s", w, got)
			}
		}
	})

	t.Run("Empty files map", func(t *testing.T) {
		var buf bytes.Buffer

		files := make(map[string]int64)

		if err := DebugRemover(files, &buf); err != nil {
			t.Fatalf("DebugRemover() on empty map returned error: %v", err)
		}

		expected := "0 files will be deleted in total\n0 B of disk space will be freed\n\nFiles to be deleted:\n\nEND\n"
		if got := buf.String(); got != expected {
			t.Errorf("wrong output for empty map\ngot:  %q\nwant: %q", got, expected)
		}
	})
}
