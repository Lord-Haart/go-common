package utils

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)

func TestSha256Salt(t *testing.T) {
	t.Logf("%v", Sha256Salt("123456"))
	t.Logf("%v", Sha256Salt("778899"))
}

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

func timestampIsEqual(t1, t2 Timestamp) bool {
	if !t1.Valid {
		return !t2.Valid
	} else {
		return t2.Valid && t1.V.In(time.UTC) == t2.V.In(time.UTC)
	}
}

func TestISODateJSON(t *testing.T) {
	testcases1 := []struct {
		param1 ISODate
		result string
	}{
		{
			param1: ISODate(time.Date(2003, time.January, 2, 0, 0, 0, 0, time.Local)),
			result: `"20030102"`,
		},
		{
			param1: ISODate(time.Date(2025, time.December, 7, 0, 0, 0, 0, time.Local)),
			result: `"20251207"`,
		},
		{
			param1: ISODate{},
			result: `null`,
		},
	}

	for _, testcase := range testcases1 {
		if b, err := json.Marshal(testcase.param1); err != nil {
			t.Fatal(err)
		} else if r := string(b); r != testcase.result {
			t.Errorf("Marshal(%v) => %v, wants %v", time.Time(testcase.param1), r, testcase.result)
		}
	}

	testcases2 := []struct {
		param1 string
		result ISODate
	}{
		{
			param1: `"20030102"`,
			result: ISODate(time.Date(2003, time.January, 2, 0, 0, 0, 0, time.Local)),
		},
		{
			param1: `"20231207"`,
			result: ISODate(time.Date(2023, time.December, 7, 0, 0, 0, 0, time.Local)),
		},
	}

	for _, testcase := range testcases2 {
		var r ISODate
		if err := json.Unmarshal([]byte(testcase.param1), &r); err != nil {
			t.Fatal(err)
		} else {
			if r != testcase.result {
				t.Errorf("Unmarshal(%v) => %v, wants %v", testcase.param1, r, testcase.result)
			}
		}
	}

	type Rec struct {
		Key  string  `json:"key,omitempty"`
		Date ISODate `json:"date,omitempty"`
	}

	testcases3 := []struct {
		param1 string
		result Rec
	}{
		{
			param1: `{ "key": "item2", "date": "20030102" }`,
			result: Rec{Key: "item2", Date: ISODate(time.Date(2003, time.January, 2, 0, 0, 0, 0, time.Local))},
		},
		{
			param1: `{ "date": "20251207" , "key":    "item3"}`,
			result: Rec{Key: "item3", Date: ISODate(time.Date(2025, time.December, 7, 0, 0, 0, 0, time.Local))},
		},
		{
			param1: `{   "key": "item3"   }`,
			result: Rec{Key: "item3", Date: ISODate{}},
		},
		{
			param1: `{   "key": "item4", "time": null   }`,
			result: Rec{Key: "item4", Date: ISODate{}},
		},
		{
			param1: `{   "key": "item4", "time": ""   }`,
			result: Rec{Key: "item4", Date: ISODate{}},
		},
	}

	for _, testcase := range testcases3 {
		var r Rec
		if err := json.Unmarshal([]byte(testcase.param1), &r); err != nil {
			t.Fatal(err)
		} else {
			if r.Key != testcase.result.Key || r.Date != testcase.result.Date {
				t.Errorf("Unmarshal(%v) => %v, wants %v", testcase.param1, r, testcase.result)
			}
		}
	}
}

func TestSet(t *testing.T) {
	st1 := NewSet("a", "", "a", "b")

	t.Logf("%v", st1)

	if st1.Len() != 3 {
		t.Errorf("len(s) => %v, wants %v", st1.Len(), 3)
	}

	t.Logf("%v", st1.AllKeys())
}

func TestRollingFileWriter(t *testing.T) {
	log.SetOutput(DefaultRollingFileWriter)

	log.Printf("[DEBUG] xxxxx\n")
	log.Printf("[DEBUG] 12345\n")
	log.Printf("[DEBUG] 67890\n")
}

func TestHasPrefixFold(t *testing.T) {
	testcases := []struct {
		p1 string
		p2 string
		r  bool
	}{
		{"", "", true},
		{"", "a", false},
		{"a", "", true},
		{"a", "a", true},
		{"a", "b", false},
		{"abc", "A", true},
		{"abc", "aB", true},
		{"abc", "aBc", true},
		{"Ahaha", "Bear ", false},
		{"Bear Ahaha", "Bear ", true},
	}

	for _, testcase := range testcases {
		if r := HasPrefixFold(testcase.p1, testcase.p2); r != testcase.r {
			t.Errorf("HasPrefixFold(%q, %q) = %v, want %v", testcase.p1, testcase.p2, r, testcase.r)
		}
	}
}

func TestTrimPrefixFold(t *testing.T) {
	testcases := []struct {
		p1 string
		p2 string
		r  string
	}{
		{"", "", ""},
		{"", "a", ""},
		{"a", "", "a"},
		{"a", "a", ""},
		{"a", "b", "a"},
		{"abc", "A", "bc"},
		{"abc", "aB", "c"},
		{"abc", "aBc", ""},
		{"Ahaha", "Bear ", "Ahaha"},
		{"Bear Ahaha", "Bear ", "Ahaha"},
	}

	for _, testcase := range testcases {
		if r := TrimPrefixFold(testcase.p1, testcase.p2); r != testcase.r {
			t.Errorf("TrimPrefixFold(%q, %q) = %q, want %q", testcase.p1, testcase.p2, r, testcase.r)
		}
	}
}
