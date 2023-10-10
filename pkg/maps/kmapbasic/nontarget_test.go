package kmapbasic_test

import (
	"testing"

	"github.com/berybox/KEC/pkg/maps/kmapbasic"
)

var (
	testKmersA = []string{
		"ACTGACGATGC",
		"ACTGACAAGTG",
		"NNNNNNNNNNN",
		"           ",
		"",
		"\n",
	}

	testKmersB = []string{
		"ACTGACGATGC ",
		" ACTGACGATGC",
		"\nACTGACGATGC",
		"ACTGACGATGC\n",
		" ",
		"\n\n",
	}
)

func TestNontargetMap(t *testing.T) {
	nontargetMap := kmapbasic.NewNontarget()

	for _, k := range testKmersA {
		nontargetMap.Add(k)
	}

	for i, k := range testKmersA {
		if !nontargetMap.Contains(k) {
			t.Fatalf("KMer %s [%d] not found", k, i)
		}
	}

	for i, k := range testKmersB {
		if nontargetMap.Contains(k) {
			t.Fatalf("KMer %s [%d] should not be present", k, i)
		}
	}
}
