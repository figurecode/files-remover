package conf

import (
	"errors"
	"io"
	"os"
	"strings"
)

var ErrMessErrStreamIsNil = errors.New("ErrStream cannot be nil")
var ErrMessOutStreamIsNil = errors.New("OutStream cannot be nil")
var ErrMessDirIsNotSpecified = errors.New("Search directory not specified")
var ErrMessFileListIsEmpty = errors.New("The file name list cannot be empty")

type Config struct {
	Dir                  string          // обязательный параметр, директория в которой будем искать
	FilesName            map[string]bool // обязательный параметр названия файлов, которые будем искать и удалять
	ExcDirs              []string        // названия поддерикторий, которые нужно исключить из обхода
	FileNameSep          string          // разделитель для разбиения названия файла на части
	IsDemo               bool            // демо-режим работы приложения, выводит только информацию
	ErrStream, OutStream io.Writer       // стандартный вывод ошибок и результатов
}

type Option func(*Config) error

func (c Config) validate() error {
	if len(c.Dir) == 0 {
		return ErrMessDirIsNotSpecified
	}

	if len(c.FilesName) == 0 {
		return ErrMessFileListIsEmpty
	}

	return nil
}

func WithErrStream(errStream io.Writer) Option {
	return func(c *Config) error {
		if errStream == nil {
			return ErrMessErrStreamIsNil
		}

		c.ErrStream = errStream

		return nil
	}
}

func WithOutStream(outStream io.Writer) Option {
	return func(c *Config) error {
		if outStream == nil {
			return ErrMessOutStreamIsNil
		}

		c.OutStream = outStream

		return nil
	}
}

func WithExcludeDir(excDir string) Option {
	return func(c *Config) error {
		if len(excDir) != 0 {
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
		ExcDirs:     make([]string, 0),
		FileNameSep: "-",
		IsDemo:      true,
		ErrStream:   os.Stderr,
		OutStream:   os.Stdout,
	}

	c.Dir = strings.TrimSpace(dir)

	if c.FilesName == nil {
		c.FilesName = make(map[string]bool)
	}
	for _, v := range fNames {
		if !c.FilesName[v] {
			c.FilesName[v] = true
		}
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
