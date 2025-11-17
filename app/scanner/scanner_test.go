package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolvePath(t *testing.T) {
	type tCases []struct {
		input    string
		expected string
	}

	chackerPath := func(t testing.TB, tCases tCases) {
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

		chackerPath(t, tCases)
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

		chackerPath(t, tCases)
	})
}
