package kecoriginal

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/berybox/KEC/pkg/filedriver"
	"github.com/berybox/KEC/pkg/filedriver/fastaoriginal"
	"github.com/berybox/KEC/pkg/maps"

	"github.com/berybox/KEC/pkg/keccore"
)

// Includer KEC implementation for KMer inclusion used in original KEC
type Includer struct {
	Options            keccore.Options
	ReferenceMap       maps.TargetMap
	PoolMap            maps.NontargetMap
	MaskMap            maps.MaskMap
	referenceFilenames []string
	poolFilenames      []string
}

// Init Set parameters
func (KEC *Includer) Init(opts keccore.Options, mapset keccore.MapSet) error {
	KEC.Options = opts
	KEC.ReferenceMap = mapset.TargetMap
	KEC.PoolMap = mapset.NontargetMap
	KEC.MaskMap = mapset.MaskMap

	return nil
}

// AddReferenceFilename Adds filename of reference organism to the queue. Analysis will be performed after calling "Run"
func (KEC *Includer) AddReferenceFilename(filename string) error {
	KEC.referenceFilenames = append(KEC.referenceFilenames, filename)
	return nil
}

// AddPoolFilename Adds filename of pool organism to the queue. Analysis will be performed after calling "Run"
func (KEC *Includer) AddPoolFilename(filename string) error {
	KEC.poolFilenames = append(KEC.poolFilenames, filename)
	return nil
}

// Run Main KEC analysis
func (KEC *Includer) Run() error {
	var FileReader filedriver.FileReader
	FileReader = fastaoriginal.FileReader{}

	var actionTime time.Time

	//Add Reference
	var namesizes []maps.NameSize
	for _, fn := range KEC.referenceFilenames {
		actionTime = time.Now()
		var err error
		namesizes, err = FileReader.ReadTarget(fn, KEC.Options.K, KEC.ReferenceMap.Add)
		if err != nil {
			return err
		}
		fmt.Fprintf(KEC.Options.Log, "Added Reference: %s, took %s\n", filepath.Base(fn), time.Since(actionTime))
	}

	//Add Pool
	for _, fn := range KEC.poolFilenames {
		actionTime = time.Now()

		err := FileReader.ReadNontarget(fn, KEC.Options.K, KEC.PoolMap.Add)
		KEC.ReferenceMap.Walk(func(kmer string, _ maps.NamePos) {
			if !KEC.PoolMap.Contains(kmer) {
				KEC.ReferenceMap.Delete(kmer)
			}
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(KEC.Options.Log, "Added Pool: %s, took %s\n", filepath.Base(fn), time.Since(actionTime))
	}

	KEC.MaskMap.Add(namesizes)

	//Unmask ramaining target kmers
	actionTime = time.Now()
	KEC.ReferenceMap.Walk(func(kmer string, namepos maps.NamePos) {
		KEC.MaskMap.Unmask(kmer, namepos)
	})
	fmt.Fprintf(KEC.Options.Log, "Unmasking took %s\n", time.Since(actionTime))

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

	fmt.Fprintf(KEC.Options.Log, "Found %d sequences\n", numSeq)

	return nil
}
