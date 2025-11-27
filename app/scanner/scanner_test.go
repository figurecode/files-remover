package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/figurecode/files-remover/conf"
	"github.com/stretchr/testify/assert"
)

func TestResolvePath(t *testing.T) {
	type tCases []struct {
		input    string
		expected string
	}

	checkerPath := func(t testing.TB, tCases tCases) {
		t.Helper()

		for _, tc := range tCases {
			result, err := ResolvePath(tc.input)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		}
	}

	t.Run("resolv absolute path", func(t *testing.T) {
		tCases := tCases{
			{"/home/user/data", "/home/user/data"},
			{"/var/log/", "/var/log/"},
			{"/tmp", "/tmp"},
			{"/", "/"},
		}

		checkerPath(t, tCases)
	})

	t.Run("resolv relative path", func(t *testing.T) {
		wd, _ := os.Getwd()

		tCases := tCases{
			{"log", filepath.Join(wd, "log")},
			{"log/catalog/import", filepath.Join(wd, "log/catalog/import")},
			{".", wd},
			{"..", filepath.Join(wd, "..")},
			{"", wd},
			{"./", wd},
			{"../parent", filepath.Join(wd, "../parent")},
		}

		checkerPath(t, tCases)
	})
}

func Test_match(t *testing.T) {
	tCests := []struct {
		name  string
		cfg   conf.Config
		fName string
		want  bool
	}{
		{
			name:  "exact match",
			cfg:   conf.Config{FilesName: map[string]bool{"document.pdf": true}},
			fName: "document.pdf",
			want:  true,
		},
		{
			name: "partial match",
			cfg: conf.Config{
				FileNameSep: "-",
				FilesName:   map[string]bool{"66f3c59b27ea50223262041": true},
			},
			fName: "66f3c59b27ea50223262041-5f8033af82716c6a4406628341d86046.pdf",
			want:  true,
		},
		{
			name: "regarding the expansion",
			cfg: conf.Config{
				FileNameSep: "-",
				FilesName:   map[string]bool{"5f8033af82716c6a4406628341d86046": true},
			},
			fName: "66f3c59b27ea50223262041-5f8033af82716c6a4406628341d86046.pdf",
			want:  true,
		},
	}

	for _, tt := range tCests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, match(tt.cfg, tt.fName))
		})
	}
}
