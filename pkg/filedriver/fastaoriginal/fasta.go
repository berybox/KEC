package fastaoriginal

const (
	// NameSymbol Fasta symbol to specify name of sequence
	NameSymbol = ">"

	fastaReadAtBuffSize   int    = 256
	fastaEntryPosBuffSize int    = 64 * 1024
	newlineSymbolS1       string = "\r"
	newlineSymbolS2       string = "\n"
	descriptionSymbolB    byte   = '>'
	descriptionSymbolS    string = ">"
	commentSymbol         string = ";"
	newlineSymbolB        byte   = '\n'
)

// FileReader Fasta file reader used in original KEC implementation
type FileReader struct{}
