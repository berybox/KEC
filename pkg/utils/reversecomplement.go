package utils

// ReverseString reverse chracter sequence. This seems to be fastest so far
func ReverseString(s string) string {
	chars := []rune(s)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}

// DNAComplement very basic but fast DNA complement
func DNAComplement(s string) string {
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		switch runes[i] {
		case 'A', 'a':
			runes[i] = 'T'
		case 'T', 't':
			runes[i] = 'A'
		case 'C', 'c':
			runes[i] = 'G'
		case 'G', 'g':
			runes[i] = 'C'
		}
	}
	return string(runes)
}
