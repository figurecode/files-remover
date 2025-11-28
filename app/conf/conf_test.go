package conf

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("correct config", func(t *testing.T) {
		filesName := []string{"file1", "file2"}
		cfg, err := New(
			"/",
			filesName,
		)

		assert.NoError(t, err)
		assert.Equal(t, "/", cfg.Dir)

		for i := range filesName {
			if _, ok := cfg.FilesName[filesName[i]]; !ok {
				t.Errorf("The files map does not contain the entire data set. Absent %s, list %q, got map: %v\n", filesName[i], filesName, cfg.FilesName)

				break
			}
		}

		assert.True(t, cfg.IsDemo)
		assert.Len(t, cfg.ExcDirs, 0)
	})

	t.Run("not correct config", func(t *testing.T) {
		cfg, err := New("", nil)

		assert.Error(t, err)
		assert.Empty(t, cfg)
	})

	t.Run("check empty Dir", func(t *testing.T) {
		_, err := New("", []string{"file1", "file2"})

		assert.ErrorIs(t, err, errMessDirIsNotSpecified)
	})

	t.Run("check empty Files", func(t *testing.T) {
		_, err := New("/", []string{})

		assert.ErrorIs(t, err, errMessFileListIsEmpty)
	})
}

func TestWithErrStream(t *testing.T) {

	t.Run("set ErrStream", func(t *testing.T) {
		cfg := &Config{}
		opt := WithErrStream(&bytes.Buffer{})
		err := opt(cfg)

		assert.NoError(t, err)
		assert.Equal(t, &bytes.Buffer{}, cfg.ErrStream)
	})

	t.Run("set nil ErrStream", func(t *testing.T) {
		cfg := &Config{}
		opt := WithErrStream(nil)
		err := opt(cfg)

		assert.ErrorIs(t, err, errMessErrStreamIsNil)
	})
}

func TestWithOutStream(t *testing.T) {

	t.Run("set OutStream", func(t *testing.T) {
		cfg := &Config{}
		opt := WithOutStream(&bytes.Buffer{})
		err := opt(cfg)

		assert.NoError(t, err)
		assert.Equal(t, &bytes.Buffer{}, cfg.OutStream)
	})

	t.Run("set nil OutStream", func(t *testing.T) {
		cfg := &Config{}
		opt := WithOutStream(nil)
		err := opt(cfg)

		assert.ErrorIs(t, err, errMessOutStreamIsNil)
	})
}

func TestWithExcludeDir(t *testing.T) {
	t.Run("set ExcDir", func(t *testing.T) {
		got := "exclude1, exclude2"
		want := []string{"exclude1", "exclude2"}

		cfg := &Config{}
		opt := WithExcludeDir(got)
		err := opt(cfg)

		assert.NoError(t, err)
		assert.Len(t, cfg.ExcDirs, len(want))
		assert.ElementsMatch(t, cfg.ExcDirs, want)

	})

	t.Run("trim ExcDir value", func(t *testing.T) {
		got := " exclude1, exclude2 "
		want := []string{"exclude1", "exclude2"}

		cfg := &Config{}
		opt := WithExcludeDir(got)
		err := opt(cfg)

		assert.NoError(t, err)
		assert.Len(t, cfg.ExcDirs, len(want))
		assert.ElementsMatch(t, cfg.ExcDirs, want)
	})

	t.Run("check empty ExcDir", func(t *testing.T) {
		got := ""
		want := make([]string, 0)

		cfg := &Config{}
		opt := WithExcludeDir(got)
		err := opt(cfg)

		assert.NoError(t, err)
		assert.Len(t, cfg.ExcDirs, len(want))
		assert.ElementsMatch(t, cfg.ExcDirs, want)
	})
}

func TestWithIsDemo(t *testing.T) {

	t.Run("set IsDemo", func(t *testing.T) {
		cfg := &Config{}

		opt := WithIsDemo("true")
		err := opt(cfg)

		assert.NoError(t, err)
		assert.Equal(t, true, cfg.IsDemo)
	})

	t.Run("set IsDemo false", func(t *testing.T) {
		cfg := &Config{}

		opt := WithIsDemo("")
		err := opt(cfg)

		assert.NoError(t, err)
		assert.Equal(t, false, cfg.IsDemo)
	})
}
