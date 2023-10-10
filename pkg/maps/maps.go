package maps

// TargetWalkFunc Function to apply to each target KMer
type TargetWalkFunc func(string, NamePos)

// NontargetWalkFunc Function to apply to each nontarget KMer
type NontargetWalkFunc func(string)

// MMapWalkFunc Function to apply to each output sequence
type MMapWalkFunc func(string, []string)

// TargetMap KMer map interface for target sequences
type TargetMap interface {
	Add(string, NamePos)
	Delete(string)
	Walk(TargetWalkFunc)
}

// NontargetMap KMer map interface for nontarget sequences
type NontargetMap interface {
	Add(string)
	Contains(string) bool
}

// MaskMap Interface for getting result sequences
type MaskMap interface {
	Add([]NameSize)
	Unmask(string, NamePos)
	Walk(MMapWalkFunc)
}
