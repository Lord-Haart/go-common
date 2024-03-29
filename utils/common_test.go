package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestSha256Salt(t *testing.T) {
	t.Logf(Sha256Salt("123456"))
	t.Logf(Sha256Salt("778899"))
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
	return time.Time(t1).In(time.UTC) == time.Time(t2).In(time.UTC)
}

func TestTimestampJSON(t *testing.T) {
	testcases1 := []struct {
		param1 Timestamp
		result int64
	}{
		{
			param1: Timestamp(time.Date(2003, time.January, 2, 13, 54, 35, 0, time.UTC)),
			result: 1041515675,
		},
		{
			param1: Timestamp(time.Date(2025, time.December, 7, 8, 1, 3, 0, time.UTC)),
			result: 1765094463,
		},
	}

	for _, testcase := range testcases1 {
		if b, err := json.Marshal(testcase.param1); err != nil {
			t.Fatal(err)
		} else if r, err := strconv.ParseInt(string(b), 10, 64); err != nil {
			t.Fatal(err)
		} else if r != testcase.result {
			t.Errorf("Marshal(%v) => %v, wants %v", time.Time(testcase.param1), r, testcase.result)
		}
	}

	testcases2 := []struct {
		param1 int64
		result Timestamp
	}{
		{
			param1: 1041515675,
			result: Timestamp(time.Date(2003, time.January, 2, 13, 54, 35, 0, time.UTC)),
		},
		{
			param1: 1765094463,
			result: Timestamp(time.Date(2025, time.December, 7, 8, 1, 3, 0, time.UTC)),
		},
	}

	for _, testcase := range testcases2 {
		var r Timestamp
		if err := json.Unmarshal([]byte(strconv.FormatInt(testcase.param1, 10)), &r); err != nil {
			t.Fatal(err)
		} else {
			if !timestampIsEqual(r, testcase.result) {
				t.Errorf("Unmarshal(%v) => %v, wants %v", testcase.param1, time.Time(r), time.Time(testcase.result))
			}
		}
	}

	type Rec struct {
		Key  string    `json:"key,omitempty"`
		Time Timestamp `json:"time,omitempty"`
	}

	testcases3 := []struct {
		param1 string
		result Rec
	}{
		{
			param1: `{ "key": "item2", "time": 1041515675 }`,
			result: Rec{Key: "item2", Time: Timestamp(time.Date(2003, time.January, 2, 13, 54, 35, 0, time.UTC))},
		},
		{
			param1: `{ "time": 1765094463, "key":    "item3"}`,
			result: Rec{Key: "item3", Time: Timestamp(time.Date(2025, time.December, 7, 8, 1, 3, 0, time.UTC))},
		},
		{
			param1: `{   "key": "item3"   }`,
			result: Rec{Key: "item3", Time: Timestamp{}},
		},
		{
			param1: `{   "key": "item4", "time": null   }`,
			result: Rec{Key: "item4", Time: Timestamp{}},
		},
		{
			param1: `{   "key": "item4", "time": ""   }`,
			result: Rec{Key: "item4", Time: Timestamp{}},
		},
	}

	for _, testcase := range testcases3 {
		var r Rec
		if err := json.Unmarshal([]byte(testcase.param1), &r); err != nil {
			t.Fatal(err)
		} else {
			if r.Key != testcase.result.Key || !timestampIsEqual(r.Time, testcase.result.Time) {
				t.Errorf("Unmarshal(%v) => %v, wants %v", testcase.param1, r, testcase.result)
			}
		}
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

func TestBooleanJSON(t *testing.T) {
	fmt.Printf("%v\n", 100_000_000)

	testcases1 := []struct {
		param1 Boolean
		result string
	}{
		{
			param1: Boolean{Valid: true, Bool: true},
			result: `true`,
		},
		{
			param1: Boolean{Valid: true, Bool: false},
			result: `false`,
		},
		{
			param1: Boolean{},
			result: `null`,
		},
	}

	for _, testcase := range testcases1 {
		if b, err := json.Marshal(testcase.param1); err != nil {
			t.Fatal(err)
		} else if r := string(b); r != testcase.result {
			t.Errorf("Marshal(%v) => %v, wants %v", testcase.param1, r, testcase.result)
		}
	}

	testcases2 := []struct {
		param1 string
		result Boolean
	}{
		{
			param1: `true`,
			result: Boolean{Valid: true, Bool: true},
		},
		{
			param1: `false`,
			result: Boolean{Valid: true, Bool: false},
		},
		{
			param1: `null`,
			result: Boolean{},
		},
	}

	for _, testcase := range testcases2 {
		var r Boolean
		if err := json.Unmarshal([]byte(testcase.param1), &r); err != nil {
			t.Fatal(err)
		} else {
			if r != testcase.result {
				t.Errorf("Unmarshal(%v) => %v, wants %v", testcase.param1, r, testcase.result)
			}
		}
	}

	type Rec struct {
		Key    string  `json:"key,omitempty"`
		Status Boolean `json:"status,omitempty"`
	}

	testcases3 := []struct {
		param1 string
		result Rec
	}{
		{
			param1: `{ "key": "item2", "status": true }`,
			result: Rec{Key: "item2", Status: Boolean{Valid: true, Bool: true}},
		},
		{
			param1: `{ "status":    false   , "key":    "item3"}`,
			result: Rec{Key: "item3", Status: Boolean{Valid: true, Bool: false}},
		},
		{
			param1: `{ "status": 0 , "key":    "item3"}`,
			result: Rec{Key: "item3", Status: Boolean{Valid: true, Bool: false}},
		},
		{
			param1: `{ "status": "yes" , "key":    "item3"}`,
			result: Rec{Key: "item3", Status: Boolean{Valid: true, Bool: true}},
		},
		{
			param1: `{ "status": "n" , "key":    "item3"}`,
			result: Rec{Key: "item3", Status: Boolean{Valid: true, Bool: false}},
		},
		{
			param1: `{   "key": "item3"   }`,
			result: Rec{Key: "item3", Status: Boolean{}},
		},
		{
			param1: `{   "key": "item4", "time": null   }`,
			result: Rec{Key: "item4", Status: Boolean{}},
		},
		{
			param1: `{   "key": "item4", "time": ""   }`,
			result: Rec{Key: "item4", Status: Boolean{}},
		},
	}

	for _, testcase := range testcases3 {
		var r Rec
		if err := json.Unmarshal([]byte(testcase.param1), &r); err != nil {
			t.Fatal(err)
		} else {
			if r != testcase.result {
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
