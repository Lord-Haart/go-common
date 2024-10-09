package encryption

import "testing"

func TestDecrypt(t *testing.T) {
	testcases := []struct {
		p1 string
		p2 []byte
		r1 string
	}{
		{p1: "H/kW2klYNBA9KOO1jw4X1BRKAqykmeFhsXWfy232SvGY1fvAM9hur+aBzpaz3s93", p2: []byte("chem123456"), r1: "12"},
		{p1: "69JEWPKNK18cDQbW37uhu2xE9xRJ/4Eml9HILlGgPOo8NYv0I3eJ7aTfvIugVl6d", p2: []byte("chem123456"), r1: "759"},
	}

	for _, testcase := range testcases {
		if r, err := decryptWithPBEHmacSHA512AndAES_256(testcase.p1, testcase.p2); err != nil {
			t.Fatalf("Decrypt(): %#v", err)
		} else if r != testcase.r1 {
			t.Errorf("Decrypt(%#v, %#v) => %#v, want %#v", testcase.p1, testcase.p2, r, testcase.r1)
		}
	}
}

func TestEncrypt(t *testing.T) {
	testcases := []struct {
		p1 string
		p2 []byte
	}{
		{p1: "12", p2: []byte("chem123456")},
		{p1: "759", p2: []byte("chem123456")},
	}

	for _, testcase := range testcases {
		if r, err := encryptWithPBEHmacSHA512AndAES_256(testcase.p1, testcase.p2); err != nil {
			t.Fatalf("Encrypt(): %#v", err)
		} else {
			t.Logf("Encrypt(%#v, %#v) => %#v", testcase.p1, testcase.p2, r)
		}
	}
}
