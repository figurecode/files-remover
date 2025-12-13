package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/figurecode/files-remover/conf"
	"github.com/stretchr/testify/assert"
)

func TestResolvePath(t *testing.T) {
	wd, _ := os.Getwd()
	tests := []struct {
		got  string
		want string
	}{
		{
			got:  "/home/user/data",
			want: "/home/user/data",
		},
		{
			got:  "/var/log/",
			want: "/var/log/",
		},
		{
			got:  "/tmp",
			want: "/tmp",
		},
		{
			got:  "/",
			want: "/",
		},
		{
			got:  "log",
			want: filepath.Join(wd, "log"),
		},
		{
			got:  "log/catalog/import",
			want: filepath.Join(wd, "log/catalog/import"),
		},
		{
			got:  ".",
			want: wd,
		},
		{
			got:  "..",
			want: filepath.Join(wd, ".."),
		},
		{
			got:  "",
			want: wd,
		},
		{
			got:  "./",
			want: wd,
		},
		{
			got:  "../parent",
			want: filepath.Join(wd, "../parent"),
		},
	}

	for _, tt := range tests {
		result, err := ResolvePath(tt.got)

		assert.NoError(t, err)
		assert.Equal(t, tt.want, result)
	}
}

func TestScanDir(t *testing.T) {
	t.Run("search by part of file name, separator specified", func(t *testing.T) {
		tmpDir := t.TempDir()

		createFiles(t, tmpDir, map[string]int64{
			"exact-match.pdf":               1024,
			"hash-part1-otherpart.pdf":      2048,
			"hash-part2.pdf":                3072,
			"ignored.txt":                   100,
			"node_modules/exact-match.pdf":  9999,
			"logs/hash-part1-otherpart.pdf": 4096,
		})

		cfg, err := conf.New(
			tmpDir,
			[]string{"hash", "hash"},
			conf.WithFileNameSep("-"),
		)
		assert.NoError(t, err)

		files, err := ScanDir(cfg)
		assert.NoError(t, err)
		assert.Len(t, files, 3)

		expected := map[string]int64{
			filepath.Join(tmpDir, "hash-part1-otherpart.pdf"):      2048,
			filepath.Join(tmpDir, "hash-part2.pdf"):                3072,
			filepath.Join(tmpDir, "logs/hash-part1-otherpart.pdf"): 4096,
		}

		for path, size := range expected {
			foundSize, ok := files[path]
			assert.True(t, ok, "expected file not found: %s", path)
			assert.Equal(t, size, foundSize)
		}
	})

	t.Run("search by file name, separator specified", func(t *testing.T) {
		tmpDir := t.TempDir()

		createFiles(t, tmpDir, map[string]int64{
			"exact-match.pdf":               1024,
			"hash-part1-otherpart.pdf":      2048,
			"hash-part2.pdf":                3072,
			"ignored.txt":                   100,
			"node_modules/exact-match.pdf":  9999,
			"logs/hash-part1-otherpart.pdf": 4096,
		})

		cfg, err := conf.New(
			tmpDir,
			[]string{"hash-part1-otherpart"},
			conf.WithFileNameSep("."),
		)
		assert.NoError(t, err)

		files, err := ScanDir(cfg)
		assert.NoError(t, err)
		assert.Len(t, files, 2)

		expected := map[string]int64{
			filepath.Join(tmpDir, "hash-part1-otherpart.pdf"):      2048,
			filepath.Join(tmpDir, "logs/hash-part1-otherpart.pdf"): 4096,
		}

		for path, size := range expected {
			foundSize, ok := files[path]
			assert.True(t, ok, "expected file not found: %s", path)
			assert.Equal(t, size, foundSize)
		}
	})

	t.Run("search by full file name", func(t *testing.T) {
		tmpDir := t.TempDir()

		createFiles(t, tmpDir, map[string]int64{
			"exact-match.pdf":               1024,
			"hash-part1-otherpart.pdf":      2048,
			"hash-part2.pdf":                3072,
			"ignored.txt":                   100,
			"node_modules/exact-match.pdf":  9999,
			"logs/hash-part1-otherpart.pdf": 4096,
		})

		cfg, err := conf.New(
			tmpDir,
			[]string{"exact-match.pdf"},
		)
		assert.NoError(t, err)

		files, err := ScanDir(cfg)
		assert.NoError(t, err)
		assert.Len(t, files, 2)

		expected := map[string]int64{
			filepath.Join(tmpDir, "exact-match.pdf"):              1024,
			filepath.Join(tmpDir, "node_modules/exact-match.pdf"): 9999,
		}

		for path, size := range expected {
			foundSize, ok := files[path]
			assert.True(t, ok, "expected file not found: %s", path)
			assert.Equal(t, size, foundSize)
		}
	})

	t.Run("exclude directory from search", func(t *testing.T) {
		tmpDir := t.TempDir()

		createFiles(t, tmpDir, map[string]int64{
			"exact-match.pdf":               1024,
			"hash-part1-otherpart.pdf":      2048,
			"hash-part2.pdf":                3072,
			"ignored.txt":                   100,
			"node_modules/exact-match.pdf":  9999,
			"logs/hash-part1-otherpart.pdf": 4096,
		})

		cfg, err := conf.New(
			tmpDir,
			[]string{"exact", "match"},
			conf.WithExcludeDir("node_modules"),
			conf.WithFileNameSep("-"),
		)
		assert.NoError(t, err)

		files, err := ScanDir(cfg)
		assert.NoError(t, err)
		assert.Len(t, files, 1)

		expected := map[string]int64{
			filepath.Join(tmpDir, "exact-match.pdf"): 1024,
		}

		for path, size := range expected {
			foundSize, ok := files[path]
			assert.True(t, ok, "expected file not found: %s", path)
			assert.Equal(t, size, foundSize)
		}

		unexpected := []string{
			filepath.Join(tmpDir, "node_modules/exact-match.pdf"),
			filepath.Join(tmpDir, "ignored.txt"),
			filepath.Join(tmpDir, "hash-part2.pdf"),
			filepath.Join(tmpDir, "hash-part1-otherpart.pdf"),
			filepath.Join(tmpDir, "logs/hash-part1-otherpart.pdf"),
		}
		for _, path := range unexpected {
			_, ok := files[path]
			assert.False(t, ok, "unexpected file found: %s", path)
		}
	})
}

func Test_match(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		names    map[string]bool
		sep      string
		want     bool
	}{
		{
			name:     "exact match, no separator",
			filename: "document.pdf",
			names:    map[string]bool{"document.pdf": true},
			sep:      "",
			want:     true,
		},
		{
			name:     "no match, no separator",
			filename: "other.txt",
			names:    map[string]bool{"document.pdf": true},
			sep:      "",
			want:     false,
		},
		{
			name:     "partial match with dash",
			filename: "66f3c59b27ea50223262041-5f8033af82716c6a4406628341d86046.pdf",
			names:    map[string]bool{"66f3c59b27ea50223262041": true},
			sep:      "-",
			want:     true,
		},
		{
			name:     "partial match in second part",
			filename: "prefix-5f8033af82716c6a4406628341d86046.pdf",
			names:    map[string]bool{"5f8033af82716c6a4406628341d86046": true},
			sep:      "-",
			want:     true,
		},
		{
			name:     "no match with separator",
			filename: "wrong-hash-file.pdf",
			names:    map[string]bool{"correct-hash": true},
			sep:      "-",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := match(tt.filename, tt.names, tt.sep)
			assert.Equal(t, tt.want, got)
		})
	}
}

func createFiles(t *testing.T, baseDir string, filesPathAndSize map[string]int64) {
	t.Helper()

	for fPath, size := range filesPathAndSize {
		realFPath := filepath.Join(baseDir, fPath)

		dir := filepath.Dir(realFPath)
		if dir != baseDir {
			assert.NoError(t, os.MkdirAll(dir, 0o755))
		}

		f, err := os.Create(realFPath)
		assert.NoError(t, err)

		if size > 0 {
			_, err := f.Write(make([]byte, size))
			assert.NoError(t, err)
		}

		assert.NoError(t, f.Close())
	}
}
