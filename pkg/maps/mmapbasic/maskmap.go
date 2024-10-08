package mmapbasic

import (
	"strings"

	"github.com/berybox/KEC/pkg/maps"
)

const (
	// MaskChar Default "blank" character to use before placing KMers
	MaskChar = '\u0000'
)

// MaskMap Basic Mask Map implementation used in original KEC
type MaskMap struct {
	mapp map[string][]byte
}

// NewMaskMap Creates new MaskMap
func NewMaskMap() *MaskMap {
	var mm MaskMap
	mm.mapp = make(map[string][]byte)
	return &mm
}

// Add Adds a name of a sequence and a blank space, which will be used for placing surviving KMers by Unmask
func (mm *MaskMap) Add(names []maps.NameSize) {
	for _, name := range names {
		mm.mapp[name.Name] = []byte(strings.Repeat(string(MaskChar), name.Size))
	}
}

// Unmask Puts KMer in the position to reconstruct original sequence
func (mm *MaskMap) Unmask(kmer string, namepos maps.NamePos) {
	l := len(mm.mapp[namepos.Name])
	if l < namepos.Pos + len(kmer) {
		l = namepos.Pos + len(kmer)
		tmp := make([]byte, l, 2*l)
		copy(tmp, mm.mapp[namepos.Name])
		mm.mapp[namepos.Name] = tmp
	}
	copy(mm.mapp[namepos.Name][namepos.Pos:], []byte(kmer))
}

// Walk Apply MMapWalkFunc to each reconstructed sequence
func (mm *MaskMap) Walk(f maps.MMapWalkFunc) {
	for name, seq := range mm.mapp {
		seqs := strings.FieldsFunc(string(seq), func(c rune) bool { return c == MaskChar })
		f(name, seqs)
	}
}
