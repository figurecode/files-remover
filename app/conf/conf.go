package conf

import (
	"errors"
	"io"
	"os"
	"strings"
)

type Config struct {
	Dir                  string          // директория в которой будем искать
	ExcDir               []string        // названия поддерикторий, которые нужно исключить из обхода
	FilesName            map[string]bool // названия файлов, которые будем искать и удалять
	FileNameSep          string          // разделитель для разбиения названия файла на части
	IsDemo               bool            // демо-режим работы приложения, выводит только информацию
	ErrStream, OutStream io.Writer       // стандартный вывод ошибок и результатов
}

type Option func(*Config) error

func WithErrStream(errStream io.Writer) Option {
	return func(c *Config) error {
		if errStream == nil {
			return errors.New("ErrStream cannot be nil")
		}

		c.ErrStream = errStream

		return nil
	}
}

func WithOutStream(outStream io.Writer) Option {
	return func(c *Config) error {
		if outStream == nil {
			return errors.New("OutStream cannot be nil")
		}

		c.OutStream = outStream

		return nil
	}
}

func WithDir(dir string) Option {
	return func(c *Config) error {
		c.Dir = strings.TrimSpace(dir)

		if len(c.Dir) == 0 {
			return errors.New("Search directory not specified")
		}

		return nil
	}
}

func WithExcludeDir(excDir string) Option {
	return func(c *Config) error {
		if len(excDir) != 0 {
			d := strings.Split(excDir, ",")

			for _, v := range d {
				c.ExcDir = append(c.ExcDir, strings.TrimSpace(v))
			}
		}

		return nil
	}
}

func WithFilesName(fNames []string) Option {
	return func(c *Config) error {
		if len(fNames) == 0 {
			return errors.New("The file name list cannot be empty.")
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
		ExcDir:      make([]string, 0),
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
