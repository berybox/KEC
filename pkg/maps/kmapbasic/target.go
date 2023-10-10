package kmapbasic

import (
	"github.com/berybox/KEC/pkg/maps"
)

// Target Basic target map implementation from original KEC
type Target struct {
	mapp map[string][]maps.NamePos
}

// NewTarget Creates new Target map
func NewTarget() *Target {
	var t Target
	t.mapp = make(map[string][]maps.NamePos)
	return &t
}

// Add Adds a KMer to the target map
func (t *Target) Add(kmer string, namepos maps.NamePos) {
	t.mapp[kmer] = append(t.mapp[kmer], namepos)
}

// Delete Deletes a KMer from the target map
func (t *Target) Delete(kmer string) {
	delete(t.mapp, kmer)
}

// Walk Apply TargetWalkFunc to each KMer
func (t *Target) Walk(f maps.TargetWalkFunc) {
	for kmer, pozs := range t.mapp {
		for _, poz := range pozs {
			f(kmer, poz)
		}
	}
}
