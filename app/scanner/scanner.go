package scanner

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/figurecode/files-remover/conf"
)

type FilesMather map[string]int64

func ResolvePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, path), nil
}

func ScanDir(cfg conf.Config) (FilesMather, error) {
	files := make(FilesMather)

	err := filepath.WalkDir(cfg.Dir, func(path string, d os.DirEntry, err error) error {

		if d.IsDir() && slices.Contains(cfg.ExcDirs, d.Name()) {
			return filepath.SkipDir
		}

		if !d.IsDir() {
			return checkFile(cfg, path, &files)
		}

		return nil
	})

	return files, err
}

func checkFile(cfg conf.Config, path string, files *FilesMather) error {
	_, fName := filepath.Split(path)

	if match(cfg, fName) {
		fInfo, err := os.Stat(path)

		if err != nil {
			return err
		}

		(*files)[path] = fInfo.Size()
	}

	return nil
}

func match(cfg conf.Config, fName string) bool {
	if len(cfg.FileNameSep) == 0 {
		if _, ok := cfg.FilesName[fName]; ok {
			return true
		}

		return false
	}

	parts := strings.Split(fName, cfg.FileNameSep)
	for _, part := range parts {
		if _, ok := cfg.FilesName[part]; ok {
			return true
		}
	}

	return false
}
