package conf

import (
	"errors"
	"io"
	"os"
)

type Config struct {
	Dir                  string    // директория в которой будем искать
	ExcDir               []string  // названия поддерикторий, которые нужно исключить из обхода
	FilesName            []string  // названия файлов, которые будем искать и удалять
	IsDemo               bool      // демо-режим работы приложения, выводит только информацию
	ErrStream, OutStream io.Writer // стандартный вывод ошибок и результатов
}

type Option func(*Config) error

func WithErrStream(errStream io.Writer) Option {
	return func(c *Config) error {
		c.ErrStream = errStream

		return nil
	}
}

func WithOutStream(outStream io.Writer) Option {
	return func(c *Config) error {
		c.OutStream = outStream

		return nil
	}
}

func WithDir(dir string) Option {
	return func(c *Config) error {
		c.Dir = dir

		return nil
	}
}

func WithFilesName(fNames []string) Option {
	return func(c *Config) error {
		if len(fNames) == 0 {
			return errors.New("The file name list cannot be empty.")
		}

		c.FilesName = fNames

		return nil
	}
}

func WithIsDemo(isDemo bool) Option {
	return func(c *Config) error {
		c.IsDemo = isDemo

		return nil
	}
}

func NewConfig(opts ...Option) (Config, error) {
	c := Config{
		Dir:       "",
		FilesName: make([]string, 0),
		IsDemo:    true,
		ErrStream: os.Stderr,
		OutStream: os.Stdout,
	}

	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return Config{}, err
		}
	}

	return c, nil
}
