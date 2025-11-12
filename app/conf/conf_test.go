package conf

import (
	"bytes"
	"testing"
)

func TestWithErrStream(t *testing.T) {
	stderrBuf := &bytes.Buffer{}
	cfg := &Config{}

	opt := WithErrStream(stderrBuf)
	err := opt(cfg)

	if err != nil {
		t.Fatalf("Error set ErrStream %v\n", err)
	}

	if cfg.ErrStream != stderrBuf {
		t.Errorf("got ErrStream %p, want %p\n", cfg.ErrStream, stderrBuf)
	}
}

func TestWithErrStreamNil(t *testing.T) {
	cfg := &Config{}

	opt := WithErrStream(nil)
	err := opt(cfg)

	if err == nil {
		t.Fatalf("WithErrStream not returned error\n")
	}
}

func TestWithOutStream(t *testing.T) {
	stdoutBuf := &bytes.Buffer{}
	cfg := &Config{}

	opt := WithOutStream(stdoutBuf)
	err := opt(cfg)

	if err != nil {
		t.Fatalf("Error set OutStream %v\n", err)
	}

	if cfg.OutStream != stdoutBuf {
		t.Errorf("got OutStream %p, want %p\n", cfg.OutStream, stdoutBuf)
	}
}

func TestWithOutStreamNil(t *testing.T) {
	cfg := &Config{}

	opt := WithOutStream(nil)
	err := opt(cfg)

	if err == nil {
		t.Fatalf("WithOutStream not returned error\n")
	}
}

func TestWithDir(t *testing.T) {
	t.Run("set Dir", func(t *testing.T) {
		cfg := &Config{}

		got := "/tmp"
		want := "/tmp"

		opt := WithDir(got)
		err := opt(cfg)

		if err != nil {
			t.Fatalf("Error set Dir %v\n", err)
		}

		if cfg.Dir != want {
			t.Errorf("got Dir %s, want %s\n", cfg.Dir, want)
		}
	})

	t.Run("trim Dir value", func(t *testing.T) {
		cfg := &Config{}

		got := " /tmp "
		want := "/tmp"

		opt := WithDir(got)
		err := opt(cfg)

		if err != nil {
			t.Fatalf("Error set Dir %v\n", err)
		}

		if cfg.Dir != want {
			t.Errorf("got Dir %s, want %s\n", cfg.Dir, want)
		}
	})

	t.Run("check empty Dir", func(t *testing.T) {
		cfg := &Config{}

		got := ""

		opt := WithDir(got)
		err := opt(cfg)

		if err == nil {
			t.Errorf("WithDir not returned error\n")
		}
	})
}

func TestWithExcludeDir(t *testing.T) {
	t.Run("set ExcDir", func(t *testing.T) {
		cfg := &Config{}

		got := "exclude1 exclude2"
		want := []string{"exclude1", "exclude2"}

		opt := WithExcludeDir(got)
		err := opt(cfg)

		if err != nil {
			t.Fatalf("Error set ExcDir %v\n", err)
		}

		if len(cfg.ExcDir) != len(want) {
			t.Errorf("got ExcDir %q, want %q\n", cfg.ExcDir, want)
		}

		for i := range want {
			if cfg.ExcDir[i] != want[i] {
				t.Errorf("got ExcDir %q, want %q\n", cfg.ExcDir, want)

				break
			}
		}
	})

	t.Run("trim ExcDir value", func(t *testing.T) {
		cfg := &Config{}

		got := " exclude1 exclude2 "
		want := []string{"exclude1", "exclude2"}

		opt := WithExcludeDir(got)
		err := opt(cfg)

		if err != nil {
			t.Fatalf("Error set ExcDir %v\n", err)
		}

		if len(cfg.ExcDir) != len(want) {
			t.Errorf("got ExcDir %q, want %q\n", cfg.ExcDir, want)
		}

		for i := range want {
			if cfg.ExcDir[i] != want[i] {
				t.Errorf("got ExcDir %q, want %q\n", cfg.ExcDir, want)

				break
			}
		}
	})

	t.Run("check empty ExcDir", func(t *testing.T) {
		cfg := &Config{}

		got := ""
		want := make([]string, 0)

		opt := WithExcludeDir(got)
		err := opt(cfg)

		if err != nil {
			t.Fatalf("Error set ExcDir %v\n", err)
		}

		if len(cfg.ExcDir) != len(want) {
			t.Errorf("got ExcDir %q, want %q\n", cfg.ExcDir, want)
		}
	})
}

func TestWithFilesName(t *testing.T) {
	t.Run("set FilesName", func(t *testing.T) {
		cfg := &Config{}

		got := []string{"file1", "file2"}
		want := []string{"file1", "file2"}

		opt := WithFilesName(got)
		err := opt(cfg)

		if err != nil {
			t.Fatalf("Error set FilesName %v\n", err)
		}

		if len(cfg.FilesName) != len(want) {
			t.Errorf("got FilesName %q, want %q\n", cfg.FilesName, want)
		}

		for i := range want {
			if cfg.FilesName[i] != want[i] {
				t.Errorf("got FilesName %q, want %q\n", cfg.FilesName, want)

				break
			}
		}
	})

	t.Run("check empty FilesName", func(t *testing.T) {
		cfg := &Config{}

		got := make([]string, 0)

		opt := WithFilesName(got)
		err := opt(cfg)

		if err == nil {
			t.Errorf("WithFilesName not returned error\n")
		}
	})
}

func TestWithIsDemo(t *testing.T) {

	check := func(t *testing.T, got string, want bool) {
		t.Helper()

		cfg := &Config{}

		opt := WithIsDemo(got)
		err := opt(cfg)

		if err != nil {
			t.Fatalf("Error set IsDemo %v\n", err)
		}

		if cfg.IsDemo != want {
			t.Errorf("got IsDemo %v, want %v\n", cfg.IsDemo, want)
		}
	}

	t.Run("set IsDemo", func(t *testing.T) {
		got := "true"
		want := true

		check(t, got, want)
	})

	t.Run("set IsDemo false", func(t *testing.T) {
		got := ""
		want := false

		check(t, got, want)
	})
}
