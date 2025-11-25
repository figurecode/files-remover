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
	Dir                  string          // директория в которой будем искать
	ExcDirs              []string        // названия поддерикторий, которые нужно исключить из обхода
	FilesName            map[string]bool // названия файлов, которые будем искать и удалять
	FileNameSep          string          // разделитель для разбиения названия файла на части
	IsDemo               bool            // демо-режим работы приложения, выводит только информацию
	ErrStream, OutStream io.Writer       // стандартный вывод ошибок и результатов
}

type Option func(*Config) error

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

func WithDir(dir string) Option {
	return func(c *Config) error {
		c.Dir = strings.TrimSpace(dir)

		if len(c.Dir) == 0 {
			return ErrMessDirIsNotSpecified
		}

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

func WithFilesName(fNames []string) Option {
	return func(c *Config) error {
		if len(fNames) == 0 {
			return ErrMessFileListIsEmpty
		}

		if c.FilesName == nil {
			c.FilesName = make(map[string]bool)
		}

		for _, v := range fNames {
			if !c.FilesName[v] {
				c.FilesName[v] = true
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

func New(opts ...Option) (Config, error) {
	c := Config{
		Dir:         "",
		ExcDirs:     make([]string, 0),
		FilesName:   make(map[string]bool),
		FileNameSep: "-",
		IsDemo:      true,
		ErrStream:   os.Stderr,
		OutStream:   os.Stdout,
	}

	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return Config{}, err
		}
	}

	return c, nil
}
