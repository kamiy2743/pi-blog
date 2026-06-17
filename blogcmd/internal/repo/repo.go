package repo

import (
	"errors"
	"os"
	"path/filepath"
)

func FindRoot() (string, error) {
	excutePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(excutePath)
	for {
		if isRoot(dir) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("repo root が見つかりません")
		}
		dir = parent
	}
}

func isRoot(dir string) bool {
	// blogのスクリプトを探す
	stat, err := os.Stat(filepath.Join(dir, "blog"))
	return err == nil && !stat.IsDir()
}
