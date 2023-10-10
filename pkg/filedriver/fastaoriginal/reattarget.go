package fastaoriginal

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/berybox/KEC/pkg/maps"
)

// ReadTarget Reads KMers of size K inside file as a target and applies TargetWalkFunc to each KMer
func (fr FileReader) ReadTarget(filename string, k int, f maps.TargetWalkFunc) ([]maps.NameSize, error) {
	var ret []maps.NameSize
	file, err := os.Open(filename)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	//get entry positions in entire file
	fastaEntryPozs := getFastaEntryPositions(file)
	fileinfo, err := file.Stat()
	if err != nil {
		return ret, err
	}
	fastaEntryPozs = append(fastaEntryPozs, int(fileinfo.Size()))

	//fill targetSeq
	for i := 0; i < len(fastaEntryPozs)-1; i++ {
		name := fastaNameAtPos(file, fastaEntryPozs[i])
		name = cleanFastaName(name)
		seq := fileSubstring(file, fastaEntryPozs[i]+len(name)+1, fastaEntryPozs[i+1])

		ret = append(ret, maps.NameSize{Name: name, Size: len(seq)})

		//create Kmers while in memory
		toKmerPosT(name, seq, k, f)
	}

	return ret, nil
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

func fastaNameAtPos(file *os.File, pos int) string {
	buf := make([]byte, fastaReadAtBuffSize)
	var ret []byte
	for {
		file.ReadAt(buf, int64(pos))
		newlinePos := bytes.IndexByte(buf, newlineSymbolB)
		if newlinePos != -1 {
			ret = append(ret, buf[:newlinePos]...)
			break
		} else {
			ret = append(ret, buf...)
			pos += fastaReadAtBuffSize
		}
	}
	return string(ret)
}

func cleanFastaName(name string) string {
	ret := strings.Trim(name, "\r\n")
	ret = strings.TrimPrefix(ret, descriptionSymbolS)
	return ret
}

func fileSubstring(file *os.File, start, stop int) string {
	buf := make([]byte, stop-start)
	file.ReadAt(buf, int64(start))
	ret := strings.ReplaceAll(string(buf), newlineSymbolS1, "")
	ret = strings.ReplaceAll(ret, newlineSymbolS2, "")
	return ret
}

func toKmerPosT(name, seq string, k int, f maps.TargetWalkFunc) {
	if len(seq) >= k {
		maxi := len(seq) - k + 1
		for i := 0; i < maxi; i++ {
			f(seq[i:i+k], maps.NamePos{Name: name, Pos: i})
		}
	}
}
