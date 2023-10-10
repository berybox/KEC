package utils

import (
	"os"
	"path"
	"path/filepath"

	"golang.org/x/exp/slices"
)

// FileListExt List files in path with specified extension. If path is not directory, path will be returned
func FileListExt(path string, extensions []string) ([]string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return []string{}, err
	}

	if fileInfo.IsDir() {
		return ReadDirExt(path, extensions)
	}

	return []string{path}, nil
}

// ReadDirExt List files in path with specified extension
func ReadDirExt(dir string, extensions []string) ([]string, error) {
	var ret []string

	file, err := os.Open(dir)
	if err != nil {
		return []string{}, err
	}
	defer file.Close()

	fileInfo, err := file.Readdir(-1)
	if err != nil {
		return []string{}, err
	}

	for _, f := range fileInfo {
		if f.IsDir() == false && slices.Contains(extensions, filepath.Ext(f.Name())) {
			retItem := filepath.FromSlash(path.Join(dir, f.Name()))
			retItem, err = filepath.EvalSymlinks(retItem)
			if err == nil && retItem != "" {
				ret = append(ret, retItem)
			}
		}
	}

	return ret, nil
}
