package utils

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"
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
			r = Timestamp(time.Time(r).In(time.UTC))
			if r != testcase.result {
				t.Errorf("Unmarshal(%v) => %v, wants %v", testcase.param1, time.Time(r), time.Time(testcase.result))
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
}
