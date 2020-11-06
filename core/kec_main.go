package KEC

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

const (
	fastaReadAtBuffSize int = 256
	fastaLineLen        int = 70

	descriptionSymbolB byte   = '>'
	descriptionSymbolS string = ">"
	commentSymbol      string = ";"
	newlineSymbolB     byte   = '\n'
	newlineSymbolS1    string = "\r"
	newlineSymbolS2    string = "\n"
	maskChar           string = " "
)

//targetPos structure to hold target Kmer position
type targetPos struct {
	Name string
	Pos  int
}

//KEC main object holding all Kmers
type KEC struct {
	k           int
	minSize     int
	maxSize     int
	include     bool
	reverse     bool
	targetMask  map[string][]byte
	targetK     map[string][]targetPos
	nontargetK  map[string]struct{}
	crossrefSeq map[string]string
	crossrefK   map[string][]targetPos
	currentfile *os.File
}

//New creates a new KEC object
func New(k, minSize, maxSize int, include bool, reverse bool) *KEC {
	return &KEC{
		k:           k,
		minSize:     minSize,
		maxSize:     maxSize,
		include:     include,
		reverse:     reverse,
		targetMask:  make(map[string][]byte),
		targetK:     make(map[string][]targetPos),
		nontargetK:  make(map[string]struct{}),
		crossrefSeq: make(map[string]string),
		crossrefK:   make(map[string][]targetPos),
		currentfile: nil,
	}
}

//NumCrossRefSeq returns number of resulting sequences
func (kec *KEC) NumCrossRefSeq() int {
	return len(kec.crossrefSeq)
}

//NumCrossRefK returns number of crossreferenced K-mers (returns 0 of target sequences were added directly - e.g. AddTargetFastaCR)
func (kec *KEC) NumCrossRefK() int {
	return len(kec.crossrefK)
}

//AddTargetFasta adds sequences from fasta formatted file (single threaded)
func (kec *KEC) AddTargetFasta(filename string) {
	kec.addTargetFastaPrototype(filename, kec.addTargetKmers)
}

//AddTargetFastaCR adds sequences from fasta formatted file and crossreference direcrtly
func (kec *KEC) AddTargetFastaCR(filename string) {
	kec.addTargetFastaPrototype(filename, kec.addTargetKmersCR)
}

func (kec *KEC) addTargetFastaPrototype(filename string, kmerFunc func(name, seq string)) {
	kec.openFile(filename)
	defer kec.closeFile()

	//get entry positions in entire file
	fastaEntryPozs := getFastaEntryPositions(kec.currentfile)
	fastaEntryPozs = append(fastaEntryPozs, kec.fileSize())

	//fill targetSeq
	for i := 0; i < len(fastaEntryPozs)-1; i++ {
		name := kec.fastaNameAtPos(fastaEntryPozs[i])
		name = cleanFastaName(name)
		seq := kec.fileSubstring(fastaEntryPozs[i]+len(name)+1, fastaEntryPozs[i+1])
		if len(seq) < kec.minSize { //can't result in large enough seq
			continue
		}
		//fmt.Printf("N:%s:N\nS:%s:S\n#####\n", name, seq)
		kec.targetMask[name] = []byte(strings.Repeat(maskChar, len(seq)))

		//create Kmers while in memory
		kmerFunc(name, seq)
	}
}

func (kec *KEC) addTargetKmers(name, seq string) {
	if len(seq) >= kec.k {
		maxi := len(seq) - kec.k + 1
		for i := 0; i < maxi; i++ {
			kmer := seq[i : i+kec.k]
			tpos := targetPos{Name: name, Pos: i}
			kec.targetK[kmer] = append(kec.targetK[kmer], tpos)
		}
	}
}

func (kec *KEC) addTargetKmersCR(name, seq string) {
	if len(seq) >= kec.k {
		maxi := len(seq) - kec.k + 1
		for i := 0; i < maxi; i++ {
			kmer := seq[i : i+kec.k]
			if _, ok := kec.nontargetK[kmer]; ok == kec.include {
				copy(kec.targetMask[name][i:], kmer)
			}
		}
	}
}

//AddNontargetFasta adds sequences from fasta formatted file (single threaded)
func (kec *KEC) AddNontargetFasta(filename string) {
	kec.openFile(filename)
	defer kec.closeFile()

	//get entry positions in entire file
	fastaEntryPozs := getFastaEntryPositions(kec.currentfile)
	fastaEntryPozs = append(fastaEntryPozs, kec.fileSize())

	//get sequences, names are there to get length
	for i := 0; i < len(fastaEntryPozs)-1; i++ {
		name := kec.fastaNameAtPos(fastaEntryPozs[i])
		seq := kec.fileSubstring(fastaEntryPozs[i]+len(name), fastaEntryPozs[i+1])
		//fmt.Printf("N:%s:N\nS:%s:S\n#####\n", name, seq)

		//create Kmers while in memory
		if len(seq) >= kec.k {
			maxi := len(seq) - kec.k + 1
			for i := 0; i < maxi; i++ {
				kmer := seq[i : i+kec.k]
				kec.nontargetK[kmer] = struct{}{}
			}

			if kec.reverse {
				revSeq := reverseString(seq)
				for i := 0; i < maxi; i++ {
					kmer := revSeq[i : i+kec.k]
					kec.nontargetK[kmer] = struct{}{}
				}
			}
		}
	}
}

//CrossReference target Kmers to nontarget Kmers if needed
func (kec *KEC) CrossReference() {
	for kmer, tpos := range kec.targetK {
		if _, ok := kec.nontargetK[kmer]; ok == kec.include {
			kec.crossrefK[kmer] = tpos
		}
	}

	//fill the masked map with crossreferenced Kmers
	for kmer, pos := range kec.crossrefK {
		for _, p := range pos {
			copy(kec.targetMask[p.Name][p.Pos:], kmer)
		}
	}
}

//MergeCrossRef recreate sequences from crossreferenced Kmers
func (kec *KEC) MergeCrossRef() {
	i := 0
	for targetName, maskedSeq := range kec.targetMask {
		seqs := strings.Fields(string(maskedSeq))
		for _, seq := range seqs {
			if len(seq) >= kec.minSize && (len(seq) <= kec.maxSize || kec.maxSize == 0) {
				entryName := fmt.Sprintf("KEC no. %d; len = %d; orig: %s", i+1, len(seq), targetName)
				kec.crossrefSeq[entryName] = seq
				i++
			}
		}
	}
}

//SaveCrossRefFasta save resulting sequences to a file formated in fasta
func (kec *KEC) SaveCrossRefFasta(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for key, val := range kec.crossrefSeq {
		file.WriteString(descriptionSymbolS + key + newlineSymbolS2)
		seqMod := len(val) % fastaLineLen
		for i := 0; i < len(val)-seqMod; i += fastaLineLen {
			file.WriteString(val[i:i+fastaLineLen] + newlineSymbolS2)
		}
		if seqMod > 0 {
			file.WriteString(val[len(val)-seqMod:] + newlineSymbolS2)
		}
		file.WriteString(newlineSymbolS2)
	}
}

func (kec *KEC) fastaNameAtPos(pos int) string {
	buf := make([]byte, fastaReadAtBuffSize)
	var ret []byte
	for {
		kec.currentfile.ReadAt(buf, int64(pos))
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

func (kec *KEC) fileSubstring(start, stop int) string {
	buf := make([]byte, stop-start)
	kec.currentfile.ReadAt(buf, int64(start))
	ret := strings.ReplaceAll(string(buf), newlineSymbolS1, "")
	ret = strings.ReplaceAll(ret, newlineSymbolS2, "")
	return ret
}
