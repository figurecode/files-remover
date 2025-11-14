package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/figurecode/files-remover/conf"
)

// Принимает параметры:
// * путь к каталогу, в котором будет происходить поиск
// * поддиректории, которые следует исключить из поиска
// * имя файла, которое будем искать, или через флаг передать путь к файлу с именами файлов
// * флаг режима, боевой или демо. В демо выводим только общую информацию без удаления файлов
func main() {
	var scanDir string
	var excDir string
	var isDemo string
	var filesName []string

	flag.StringVar(&scanDir, "d", "", "Путь к директорию, в которой будет происходить поиск")
	flag.StringVar(&excDir, "e", "", "Название поддиректорий, которые нужно исключить из поиска")
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
}
