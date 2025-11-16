package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/figurecode/files-remover/conf"
)

// Принимает параметры:
// * -d - обязательно, путь к каталогу, в котором будет происходить поиск
// * -e - опционально, поддиректории, которые следует исключить из поиска
// * -m - опционально, по умолчанию true, флаг режима, боевой или демо. В демо выводим только общую информацию без удаления файлов
// *  имя файла, которое будем искать, или через флаг передать путь к файлу с именами файлов
func main() {
	var scanDir string
	var excDir string
	var isDemo string
	var filesName []string

	flag.StringVar(&scanDir, "d", "", "Путь к директорию, в которой будет происходить поиск")
	flag.StringVar(&excDir, "e", "", "Название поддиректорий, которые нужно исключить из поиска, через запятаю")
	flag.StringVar(&isDemo, "m", "true", "Включить демо-режим")
	flag.Parse()

	filesName = flag.Args();

	cfg, err := conf.New(
		conf.WithDir(scanDir),
		conf.WithExcludeDir(excDir),
		conf.WithIsDemo("true"),
		conf.WithFilesName(filesName),
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error configuration: %v\n", err)
	}

	fmt.Printf("%v\n", cfg)

	err = filepath.WalkDir(cfg.Dir, func(path string, d os.DirEntry, err error) error {

		fmt.Println(d.IsDir(), d.Name(), slices.Contains(cfg.ExcDir, path))

		if d.IsDir() && slices.Contains(cfg.ExcDir, d.Name()) {
			return filepath.SkipDir
		}

		if d.IsDir() {
			fmt.Fprintf(cfg.OutStream, "--- %s\n", path)	
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(cfg.ErrStream, "Error traversing directory %q: %v\n", cfg.Dir, err)

		os.Exit(1)
	}

	fmt.Println("Ok")
	os.Exit(0)
}
