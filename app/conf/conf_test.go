package conf

import (
	"bytes"
	"io"
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

		assert.Equal(t, "/", cfg.Dir)

		for i := range filesName {
			if _, ok := cfg.FilesName[filesName[i]]; !ok {
				t.Errorf("The files map does not contain the entire data set. Absent %s, list %q, got map: %v\n", filesName[i], filesName, cfg.FilesName)

				break
			}
		}

		assert.NoError(t, err)
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

		assert.Error(t, err, ErrMessDirIsNotSpecified.Error())
	})

	t.Run("check empty Files", func(t *testing.T) {
		_, err := New("/", []string{})

		assert.Error(t, err, ErrMessFileListIsEmpty.Error())
	})
}

func TestWithErrStream(t *testing.T) {

	t.Run("set ErrStream", func(t *testing.T) {
		errStreamHelper(t, &bytes.Buffer{}, &bytes.Buffer{}, nil)
	})

	t.Run("set nil ErrStream", func(t *testing.T) {
		errStreamHelper(t, nil, nil, ErrMessErrStreamIsNil)
	})
}

func TestWithOutStream(t *testing.T) {

	t.Run("set OutStream", func(t *testing.T) {
		outStreamHelper(t, &bytes.Buffer{}, &bytes.Buffer{}, nil)
	})

	t.Run("set nil OutStream", func(t *testing.T) {
		outStreamHelper(t, nil, nil, ErrMessOutStreamIsNil)
	})
}

func TestWithExcludeDir(t *testing.T) {
	t.Run("set ExcDir", func(t *testing.T) {
		excDirsHelper(t, "exclude1, exclude2", []string{"exclude1", "exclude2"})
	})

	t.Run("trim ExcDir value", func(t *testing.T) {
		excDirsHelper(t, " exclude1, exclude2 ", []string{"exclude1", "exclude2"})
	})

	t.Run("check empty ExcDir", func(t *testing.T) {
		excDirsHelper(t, "", []string{})
	})
}

func TestWithIsDemo(t *testing.T) {

	t.Run("set IsDemo", func(t *testing.T) {
		isDemoHelper(t, "true", true)
	})

	t.Run("set IsDemo false", func(t *testing.T) {
		isDemoHelper(t, "", false)
	})
}

func errStreamHelper(t testing.TB, got, want io.Writer, wantErr error) {
	t.Helper()

	cfg := &Config{}
	opt := WithErrStream(got)
	err := opt(cfg)

	checkerEqual(t, cfg.ErrStream, want, err, wantErr)
}

func outStreamHelper(t testing.TB, got, want io.Writer, wantErr error) {
	t.Helper()

	cfg := &Config{}
	opt := WithOutStream(got)
	err := opt(cfg)

	checkerEqual(t, cfg.OutStream, want, err, wantErr)
}

func excDirsHelper(t testing.TB, got string, want []string) {
	t.Helper()

	cfg := &Config{}

	opt := WithExcludeDir(got)
	err := opt(cfg)

	checkerSlice(t, cfg.ExcDirs, want, err, nil)
}

func isDemoHelper(t *testing.T, got string, want bool) {
	t.Helper()

	cfg := &Config{}

	opt := WithIsDemo(got)
	err := opt(cfg)

	checkerEqual(t, cfg.IsDemo, want, err, nil)
}

func checkerEqual(t testing.TB, got, want interface{}, gotErr, wantErr error) {
	if wantErr == nil {
		assert.NoError(t, gotErr)
		assert.Equal(t, want, got)
	}

	if wantErr != nil {
		assert.Error(t, gotErr, wantErr.Error())
	}
}

func checkerSlice(t testing.TB, got, want []string, gotErr, wantErr error) {
	assert.NoError(t, gotErr)
	assert.Len(t, got, len(want))
	assert.ElementsMatch(t, got, want)

	if wantErr != nil {
		assert.Error(t, gotErr, wantErr.Error())
	}
}
