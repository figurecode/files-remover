package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/figurecode/files-remover/conf"
	"github.com/figurecode/files-remover/remover"
	"github.com/figurecode/files-remover/scanner"
)

// Принимает параметры:
// * -d - обязательно, путь к каталогу, в котором будет происходить поиск
// * -e - опционально, поддиректории, которые следует исключить из поиска
// * -m - опционально, по умолчанию true, флаг режима, боевой или демо. В демо выводим только общую информацию без удаления файлов
// * -s - опционально, по умолчанию "-", разделитель для разбиения названия файла на части
// *  имя файла, которое будем искать, или через флаг передать путь к файлу с именами файлов
func main() {
	var scanDir string
	var excDir string
	var fileNameSep string
	var isDemo string
	var filesName []string

	flag.StringVar(&scanDir, "d", "", "Путь к директорию, в которой будет происходить поиск")
	flag.StringVar(&excDir, "e", "", "Название поддиректорий, которые нужно исключить из поиска, через запятаю")
	flag.StringVar(&isDemo, "m", "true", "Включить демо-режим")
	flag.StringVar(&fileNameSep, "s", "-", "Разделитель для разбиения названия файла на части")
	flag.Parse()

	filesName = flag.Args()

	path, err := scanner.ResolvePath(scanDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scan dir: %v\n", err)
	}

	cfg, err := conf.New(
		path,
		filesName,
		conf.WithExcludeDir(excDir),
		conf.WithIsDemo("true"),
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error configuration: %v\n", err)
	}

	files, err := scanner.ScanDir(cfg)

	if err != nil {
		fmt.Fprintf(cfg.ErrStream, "Error traversing directory %q: %v\n", cfg.Dir, err)

		os.Exit(1)
	}

	remover := remover.NewRemover(cfg)

	err = remover.Execute(files)

	if err != nil {
		fmt.Fprintf(cfg.ErrStream, "Error remove files %v\n", err)

		os.Exit(1)
	}

	fmt.Println("Ok")
	os.Exit(0)
}
