package filedriver

import (
	"github.com/berybox/KEC/pkg/maps"
)

// FastaExtensions File extensions to use when listing directory
var FastaExtensions = []string{".fasta", ".fna", ".ffn", ".faa", ".frn"}

// FileReader Interface for reading files for KEC
type FileReader interface {
	ReadTarget(string, int, maps.TargetWalkFunc) ([]maps.NameSize, error)
	ReadNontarget(string, int, maps.NontargetWalkFunc) error
}
