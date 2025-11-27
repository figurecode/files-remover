package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/figurecode/files-remover/conf"
)

type FoundFiles map[string]int64

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

func ScanDir(cfg conf.Config) (FoundFiles, error) {
	info, err := os.Stat(cfg.Dir)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", cfg.Dir)
	}

	files := make(FoundFiles)

	err = filepath.WalkDir(cfg.Dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			return checkFile(cfg, path, d, &files)
		}

		if len(cfg.ExcDirs) > 0 && slices.Contains(cfg.ExcDirs, d.Name()) {
			return filepath.SkipDir
		}

		return nil
	})

	return files, err
}

func checkFile(cfg conf.Config, path string, d os.DirEntry, files *FoundFiles) error {
	_, fName := filepath.Split(path)

	if !match(cfg, fName) {
		return nil
	}

	fInfo, err := d.Info()

	if err != nil {
		return nil
	}

	(*files)[path] = fInfo.Size()

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
		cleanPart := strings.TrimSuffix(part, filepath.Ext(part))

		if _, ok := cfg.FilesName[cleanPart]; ok {
			return true
		}
	}

	return false
}
