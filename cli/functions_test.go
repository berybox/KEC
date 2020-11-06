package main

import (
	"fmt"
	"testing"
)

func TestFileInfo(t *testing.T) {
	files := []string{"./main.go", "./functions.go", "templates.go", "this_should_produce_error.go", "..", ".", "", "/\\//\\"}
	for _, f := range files {
		fmt.Printf("File %s has status %d\n", f, fileInfo(f))
	}
}

func TestShowErrorString(t *testing.T) {
	showErrorString("Chyba", false)
}
