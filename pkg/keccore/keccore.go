package keccore

import (
	"io"

	"github.com/berybox/KEC/pkg/maps"
)

// Excluder Main KEC interface for KMer exclusion
type Excluder interface {
	Init(Options, MapSet) error

	AddTargetFilename(string) error
	AddNontargetFilename(string) error

	Run() error
}

// Includer Main KEC interface for KMer inclusion
type Includer interface {
	Init(Options, MapSet) error

	AddReferenceFilename(string) error
	AddPoolFilename(string) error

	Run() error
}

// Options Parameters to use for KEC run
type Options struct {
	K                 int
	MinSize           int
	MaxSize           int
	ReverseComplement bool
	Output            io.WriteCloser
	Log               io.Writer
}

// MapSet Set of mapping implementations
type MapSet struct {
	TargetMap    maps.TargetMap
	NontargetMap maps.NontargetMap
	MaskMap      maps.MaskMap
}
