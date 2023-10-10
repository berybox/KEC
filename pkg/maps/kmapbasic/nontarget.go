package kmapbasic

// Nontarget Basic nontarget map implementation from original KEC
type Nontarget struct {
	mapp map[string]struct{}
}

// NewNontarget Creates new Nontarget map
func NewNontarget() *Nontarget {
	var nt Nontarget
	nt.mapp = make(map[string]struct{})
	return &nt
}

// Add Adds a KMer to the nontarget map
func (nt *Nontarget) Add(kmer string) {
	nt.mapp[kmer] = struct{}{}
}

// Contains Checks if the map contains KMer
func (nt *Nontarget) Contains(kmer string) bool {
	_, ret := nt.mapp[kmer]
	return ret
}
