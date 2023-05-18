package utils

import (
	"testing"
)

func TestSplitInt(t *testing.T) {
	sa1 := []string{
		"",
		" ",
		", ",
		" , ,,",
	}

	for _, s := range sa1 {
		if r := SplitAsInt[int](s, ","); len(r) != 0 {
			t.Errorf("SplitInt(%v) => %v, wants []", s, r)
		}
	}

	sa2 := []string{
		"6",
		"7,",
		" 6, 899883, ",
		"7723, 234, 2323, 477 , 66",
	}

	for _, s := range sa2 {
		t.Logf("SplitInt(%v) => %v", s, SplitAsInt[int64](s, ","))
	}
}
