package kecoriginal

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/berybox/KEC/pkg/filedriver"
	"github.com/berybox/KEC/pkg/filedriver/fastaoriginal"
	"github.com/berybox/KEC/pkg/maps"
	"github.com/berybox/KEC/pkg/utils"

	"github.com/berybox/KEC/pkg/keccore"
)

// Excluder KEC implementation for KMer exclusion used in original KEC
type Excluder struct {
	Options            keccore.Options
	TargetMap          maps.TargetMap
	NontargetMap       maps.NontargetMap
	MaskMap            maps.MaskMap
	targetFilenames    []string
	nontargetFilenames []string
}

// Init Set parameters
func (KEC *Excluder) Init(opts keccore.Options, mapset keccore.MapSet) error {
	KEC.Options = opts
	KEC.TargetMap = mapset.TargetMap
	KEC.NontargetMap = mapset.NontargetMap
	KEC.MaskMap = mapset.MaskMap

	return nil
}

// AddTargetFilename Adds filename of target organism to the queue. Analysis will be performed after calling "Run"
func (KEC *Excluder) AddTargetFilename(filename string) error {
	KEC.targetFilenames = append(KEC.targetFilenames, filename)
	return nil
}

// AddNontargetFilename Adds filename of nontarget organism to the queue. Analysis will be performed after calling "Run"
func (KEC *Excluder) AddNontargetFilename(filename string) error {
	KEC.nontargetFilenames = append(KEC.nontargetFilenames, filename)
	return nil
}

// Run Main KEC analysis
func (KEC *Excluder) Run() error {
	var FileReader filedriver.FileReader
	FileReader = fastaoriginal.FileReader{}

	var actionTime time.Time

	//Add Nontarget
	for _, fn := range KEC.nontargetFilenames {
		actionTime = time.Now()
		err := FileReader.ReadNontarget(fn, KEC.Options.K, KEC.NontargetMap.Add)
		if err != nil {
			return err
		}
		fmt.Fprintf(KEC.Options.Log, "Reading NONTARGET file: \"%s\" took %s\n", filepath.Base(fn), time.Since(actionTime))
	}

	//Add Target
	var namesizes []maps.NameSize
	for _, fn := range KEC.targetFilenames {
		actionTime = time.Now()
		var err error
		namesizes, err = FileReader.ReadTarget(fn, KEC.Options.K, func(kmer string, namepos maps.NamePos) {
			if KEC.Options.ReverseComplement {
				if !KEC.NontargetMap.Contains(kmer) && !KEC.NontargetMap.Contains(utils.ReverseString(utils.DNAComplement(kmer))) {
					KEC.TargetMap.Add(kmer, namepos)
				}
			} else {
				if !KEC.NontargetMap.Contains(kmer) {
					KEC.TargetMap.Add(kmer, namepos)
				}
			}
		})
		if err != nil {
			return err
		}
		fmt.Fprintf(KEC.Options.Log, "Reading TARGET file: \"%s\" took %s\n", filepath.Base(fn), time.Since(actionTime))
	}

	KEC.MaskMap.Add(namesizes)

	//Unmask ramaining target kmers
	actionTime = time.Now()
	KEC.TargetMap.Walk(func(kmer string, namepos maps.NamePos) {
		KEC.MaskMap.Unmask(kmer, namepos)
	})
	fmt.Fprintf(KEC.Options.Log, "Merging k-mers (%d nt) took %s\n", KEC.Options.K, time.Since(actionTime))

	//Write output
	var numSeq int
	KEC.MaskMap.Walk(func(name string, seqs []string) {
		for i, seq := range seqs {
			if len(seq) >= KEC.Options.MinSize && (len(seq) <= KEC.Options.MaxSize || KEC.Options.MaxSize == 0) {
				fmt.Fprintf(KEC.Options.Output, ">KEC no. %d; len = %d; orig: %s\n%s\n\n", i+1, len(seq), name, seq)
				numSeq++
			}
		}
	})

	fmt.Fprintf(KEC.Options.Log, "Found %d unique sequences\n", numSeq)

	return nil
}
