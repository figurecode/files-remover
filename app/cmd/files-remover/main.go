package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/figurecode/files-remover/conf"
	"github.com/figurecode/files-remover/remover"
	"github.com/figurecode/files-remover/scanner"
)

func main() {
	var scanDir string
	var excDir string
	var fileNameSep string
	var isDemo string
	var filesName []string

	flag.StringVar(&scanDir, "d", "", "Директория, в которой будет происходить поиск. Если не указан, то используется директория запуска скрипта")
	flag.StringVar(&excDir, "e", "", "Исключаемые поддиректории через запятую")
	flag.StringVar(&isDemo, "m", "true", "Режим: true — демо, false — удаление (по умолчанию true)")
	flag.StringVar(&fileNameSep, "s", "-", "Разделитель в имени файла (по умолчанию '-')")

	if len(os.Args) == 1 || (len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help")) {
		fmt.Printf(`
files-remover — массовая очистка файлов по шаблону имени

Использование:
	files-remover -d <директория> [флаги] <шаблон1> [шаблон2...]

Флаги:
	-d string   Директория для поиска (если не указана, то используется директория запуска скрипта)
	-e string   Исключаемые поддиректории через запятую
	-m string   Режим: true — демо, false — удаление (по умолчанию true)
	-s string   Разделитель в имени файла (по умолчанию "-")

Примеры:
	files-remover -d /tmp temp-log backup-2024
	files-remover -d /var/log -m false -e journal access-2024.log
`)
		os.Exit(0)
	}

	flag.Parse()

	filesName = flag.Args()

	path, err := scanner.ResolvePath(scanDir)
	if err != nil {
		log.Fatalf("Error scan dir: %v\n", err)
	}

	cfg, err := conf.New(
		path,
		filesName,
		conf.WithExcludeDir(excDir),
		conf.WithIsDemo(isDemo),
		conf.WithFileNameSep(fileNameSep),
	)

	if err != nil {
		log.Fatalf("Error configuration: %v\n", err)
	}

	files, err := scanner.ScanDir(cfg)

	if err != nil {
		fmt.Fprintf(cfg.ErrStream, "Error traversing directory %q: %v\n", cfg.Dir, err)

		os.Exit(1)
	}

	if cfg.IsDemo {
		err = remover.DebugRemover(files, cfg.OutStream)
	} else {
		err = remover.Execute(files)
	}

	if err != nil {
		fmt.Fprintf(cfg.ErrStream, "Error remove files %v\n", err)

		os.Exit(1)
	}
}
