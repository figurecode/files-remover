package conf

import (
	"errors"
	"io"
	"os"
	"strings"
)

var errMessErrStreamIsNil = errors.New("errStream cannot be nil")
var errMessOutStreamIsNil = errors.New("outStream cannot be nil")
var errMessDirIsNotSpecified = errors.New("search directory not specified")
var errMessFileListIsEmpty = errors.New("the file name list cannot be empty")

type Config struct {
	Dir                  string
	FilesName            map[string]bool
	ExcDirs              []string
	FileNameSep          string
	IsDemo               bool
	ErrStream, OutStream io.Writer
}

type Option func(*Config) error

func (c Config) validate() error {
	if c.Dir == "" {
		return errMessDirIsNotSpecified
	}

	if len(c.FilesName) == 0 {
		return errMessFileListIsEmpty
	}

	return nil
}

func WithErrStream(errStream io.Writer) Option {
	return func(c *Config) error {
		if errStream == nil {
			return errMessErrStreamIsNil
		}

		c.ErrStream = errStream

		return nil
	}
}

func WithOutStream(outStream io.Writer) Option {
	return func(c *Config) error {
		if outStream == nil {
			return errMessOutStreamIsNil
		}

		c.OutStream = outStream

		return nil
	}
}

func WithExcludeDir(excDir string) Option {
	return func(c *Config) error {
		if excDir != "" {
			d := strings.Split(excDir, ",")

			for _, v := range d {
				c.ExcDirs = append(c.ExcDirs, strings.TrimSpace(v))
			}
		}

		return nil
	}
}

func WithFileNameSep(sep string) Option {
	return func(c *Config) error {
		c.FileNameSep = sep

		return nil
	}
}

func WithIsDemo(isDemo string) Option {
	return func(c *Config) error {
		c.IsDemo = isDemo == "true"

		return nil
	}
}

func New(dir string, fNames []string, opts ...Option) (Config, error) {
	c := Config{
		Dir:         strings.TrimSpace(dir),
		FilesName:   make(map[string]bool, 0),
		ExcDirs:     make([]string, 0),
		FileNameSep: "",
		IsDemo:      true,
		ErrStream:   os.Stderr,
		OutStream:   os.Stdout,
	}

	for _, v := range fNames {
		c.FilesName[v] = true
	}

	err := c.validate()

	if err != nil {
		return Config{}, err
	}

	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return Config{}, err
		}
	}

	return c, nil
}
