package fastaoriginal

import (
	"os"

	"github.com/berybox/KEC/pkg/maps"
)

// ReadNontarget Reads KMers of size K inside file as a nontarget and applies NontargetWalkFunc to each KMer
func (fr FileReader) ReadNontarget(filename string, k int, f maps.NontargetWalkFunc) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	//get entry positions in entire file
	fastaEntryPozs := getFastaEntryPositions(file)
	fileinfo, err := file.Stat()
	if err != nil {
		return err
	}
	fastaEntryPozs = append(fastaEntryPozs, int(fileinfo.Size()))

	//get sequences, names are there to get length
	for i := 0; i < len(fastaEntryPozs)-1; i++ {
		name := fastaNameAtPos(file, fastaEntryPozs[i])
		seq := fileSubstring(file, fastaEntryPozs[i]+len(name), fastaEntryPozs[i+1])

		//create Kmers while in memory
		if len(seq) >= k {
			maxi := len(seq) - k + 1
			for i := 0; i < maxi; i++ {
				f(seq[i : i+k])
			}
		}

	}

	return nil
}
