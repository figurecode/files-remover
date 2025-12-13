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

	flag.StringVar(&scanDir, "d", "", "Directory to search in. If not specified, the directory from which the program is run will be used")
	flag.StringVar(&excDir, "e", "", "Excluded subdirectories (comma-separated)")
	flag.StringVar(&isDemo, "m", "true", "Mode: true — demo (dry-run), false — actual deletion (default: true)")
	flag.StringVar(&fileNameSep, "s", "", "Separator in filename (default: empty). If not specified, search is performed by exact full filename including extension")

	if len(os.Args) == 1 || (len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help")) {
		fmt.Printf(`
files-remover — bulk file removal by name pattern

Usage:
	files-remover -d <directory> [flags] <pattern1> [pattern2...]

Flags:
	-d string   Directory to search (if omitted, current working directory is used)
	-e string   Excluded subdirectories (comma-separated)
	-m string   Mode: true — demo/dry-run, false — real deletion (default: true)
	-s string   Filename separator (default: empty)

Examples:
	files-remover -d /tmp temp-log backup-2024-10-12.tgz
	files-remover -d /tmp temp-log -s . backup-2024
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
