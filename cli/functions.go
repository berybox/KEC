package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

var (
	errorPrefix          = "ERROR: "
	validFastaExtensions = []string{".fasta", ".fna", ".ffn", ".faa", ".frn"}
)

func fileInfo(filename string) byte {
	info, err := os.Stat(filename)
	if err != nil {
		return 0
	}
	switch mode := info.Mode(); {
	case mode.IsRegular():
		return 1
	case mode.IsDir():
		return 2
	}
	return 0
}

func showErrorString(errorText string, exit bool) {
	showError(fmt.Errorf(errorText), exit)
}

func showError(errorText error, exit bool) {
	fmt.Fprintf(os.Stderr, errorPrefix+errorText.Error()+"\n")
	if exit {
		os.Exit(1)
	}
}

func showMsg(msg string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, msg, a...)
}

func fileList(path string) []string {
	if fileInfo(path) == 2 {
		return readDir(path)
	}
	return []string{path}
}

func readDir(dir string) []string {
	var ret []string

	file, err := os.Open(dir)
	if err != nil {
		showError(err, false)
		return ret
	}
	defer file.Close()

	fileInfo, err := file.Readdir(-1)
	if err != nil {
		showError(err, false)
		return ret
	}

	for _, f := range fileInfo {
		//if f.IsDir() == false && filepath.Ext(f.Name()) == ".fasta" {
		if f.IsDir() == false && containsString(validFastaExtensions, filepath.Ext(f.Name())) {
			retItem := filepath.FromSlash(path.Join(dir, f.Name()))
			retItem, err = filepath.EvalSymlinks(retItem)
			if err == nil && retItem != "" {
				ret = append(ret, retItem)
			}
		}
	}

	return ret
}

func containsString(str []string, s string) bool {
	for _, a := range str {
		if a == s {
			return true
		}
	}
	return false
}
