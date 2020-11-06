package KEC

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	fastaEntryPosBuffSize int = 64 * 1024
)

var (
	validFastaExtensions = []string{".fasta", ".fna", ".ffn", ".faa", ".frn"}
)

func (kec *KEC) openFile(filename string) {
	var err error

	kec.currentfile, err = os.Open(filename)
	if err != nil {
		panic(err)
	}
}

func (kec *KEC) closeFile() {
	kec.currentfile.Close()
	kec.currentfile = nil
}

func getFastaEntryPositions(r io.Reader) []int {
	var pozs []int

	count := -1

	buf := make([]byte, fastaEntryPosBuffSize)
	zeroBuf := make([]byte, fastaEntryPosBuffSize)

	for {
		copy(buf, zeroBuf)
		readSize, err := r.Read(buf)
		count++
		if err != nil {
			if err != io.EOF {
				fmt.Printf("ERROR: Reading file")
			}
			break
		}

		var buffPosition int
		for {
			i := bytes.IndexByte(buf[buffPosition:], descriptionSymbolB)
			if i == -1 || readSize == buffPosition {
				break
			}
			buffPosition += i + 1
			poz := buffPosition - 1 + (fastaEntryPosBuffSize * count)

			pozs = append(pozs, poz)
		}
	}

	return pozs
}

func (kec *KEC) fileSize() int {
	finfo, err := kec.currentfile.Stat()
	if err != nil {
		panic(err)
	}
	return int(finfo.Size())
}

func cleanFastaName(name string) string {
	//ret := strings.ReplaceAll(name, newlineSymbolS1, "")
	//ret = strings.ReplaceAll(ret, newlineSymbolS2, "")
	ret := strings.Trim(name, "\r\n")
	ret = strings.TrimPrefix(ret, descriptionSymbolS)
	return ret
}

func readDir(dir string) []string {
	file, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileInfo, err := file.Readdir(-1)
	if err != nil {
		panic(err)
	}

	var ret []string
	for _, f := range fileInfo {
		//if f.IsDir() == false && filepath.Ext(f.Name()) == ".fasta" {
		if f.IsDir() == false && containsString(validFastaExtensions, filepath.Ext(f.Name())) {
			retItem := filepath.FromSlash(path.Join(dir, f.Name()))
			ret = append(ret, retItem)
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

func reverseString(str string) string {
	runes := []rune(str)
	//Reverse
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	//Complement
	for i := 0; i < len(runes); i++ {
		switch runes[i] {
		case 'A', 'a':
			runes[i] = 'T'
		case 'T', 't':
			runes[i] = 'A'
		case 'C', 'c':
			runes[i] = 'G'
		case 'G', 'g':
			runes[i] = 'C'
		}
	}

	return string(runes)
}
