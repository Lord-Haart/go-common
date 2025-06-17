package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestBooleanJSON(t *testing.T) {
	fmt.Printf("%v\n", 100_000_000)

	testcases1 := []struct {
		param1 Boolean
		result string
	}{
		{
			param1: Boolean{Valid: true, V: true},
			result: `true`,
		},
		{
			param1: Boolean{Valid: true, V: false},
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
			result: Boolean{Valid: true, V: true},
		},
		{
			param1: `false`,
			result: Boolean{Valid: true, V: false},
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
			result: Rec{Key: "item2", Status: Boolean{Valid: true, V: true}},
		},
		{
			param1: `{ "status":    false   , "key":    "item3"}`,
			result: Rec{Key: "item3", Status: Boolean{Valid: true, V: false}},
		},
		{
			param1: `{ "status": 0 , "key":    "item3"}`,
			result: Rec{Key: "item3", Status: Boolean{Valid: true, V: false}},
		},
		{
			param1: `{ "status": "yes" , "key":    "item3"}`,
			result: Rec{Key: "item3", Status: Boolean{Valid: true, V: true}},
		},
		{
			param1: `{ "status": "n" , "key":    "item3"}`,
			result: Rec{Key: "item3", Status: Boolean{Valid: true, V: false}},
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

func TestStringJSON(t *testing.T) {
	testcases1 := []struct {
		param1 String
		result string
	}{
		{
			param1: String{Valid: true, V: "aaa"},
			result: `"aaa"`,
		},
		{
			param1: String{Valid: true, V: "ddca"},
			result: `"ddca"`,
		},
		{
			param1: String{},
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
		result String
	}{
		{
			param1: `"aaa"`,
			result: String{Valid: true, V: "aaa"},
		},
		{
			param1: `"ddca"`,
			result: String{Valid: true, V: "ddca"},
		},
		{
			param1: `null`,
			result: String{},
		},
	}

	for _, testcase := range testcases2 {
		var r String
		if err := json.Unmarshal([]byte(testcase.param1), &r); err != nil {
			t.Fatal(err)
		} else {
			if r != testcase.result {
				t.Errorf("Unmarshal(%v) => %v, wants %v", testcase.param1, r, testcase.result)
			}
		}
	}

	type Rec struct {
		Key    string `json:"key,omitempty"`
		Status String `json:"status,omitempty"`
	}

	testcases3 := []struct {
		param1 string
		result Rec
	}{
		{
			param1: `{ "key": "item2", "status": "aaa" }`,
			result: Rec{Key: "item2", Status: String{Valid: true, V: "aaa"}},
		},
		{
			param1: `{ "status":    "ddca"   , "key":    "item3"}`,
			result: Rec{Key: "item3", Status: String{Valid: true, V: "ddca"}},
		},
		{
			param1: `{ "status": "0" , "key":    "item3"}`,
			result: Rec{Key: "item3", Status: String{Valid: true, V: "0"}},
		},
		{
			param1: `{ "status": "yes" , "key":    "item3"}`,
			result: Rec{Key: "item3", Status: String{Valid: true, V: "yes"}},
		},
		{
			param1: `{ "status": "n" , "key":    "item3"}`,
			result: Rec{Key: "item3", Status: String{Valid: true, V: "n"}},
		},
		{
			param1: `{   "key": "item3"   }`,
			result: Rec{Key: "item3", Status: String{}},
		},
		{
			param1: `{   "key": "item4", "time": null   }`,
			result: Rec{Key: "item4", Status: String{}},
		},
		{
			param1: `{   "key": "item4", "time": ""   }`,
			result: Rec{Key: "item4", Status: String{}},
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

func TestTimestampJSON(t *testing.T) {
	testcases1 := []struct {
		param1 Timestamp
		result float64
	}{
		{
			param1: Timestamp{Valid: true, V: time.Date(2003, time.January, 2, 13, 54, 35, 0, time.UTC)},
			result: 1041515675,
		},
		{
			param1: Timestamp{Valid: true, V: time.Date(2025, time.December, 7, 8, 1, 3, 0, time.UTC)},
			result: 1765094463,
		},
		{
			param1: Timestamp{Valid: true, V: time.Date(2025, time.December, 7, 8, 1, 3, 775_000000, time.UTC)},
			result: 1765094463.775,
		},
	}

	for _, testcase := range testcases1 {
		if b, err := json.Marshal(testcase.param1); err != nil {
			t.Fatal(err)
		} else if r, err := strconv.ParseFloat(string(b), 64); err != nil {
			t.Fatal(err)
		} else if r != testcase.result {
			t.Errorf("Marshal(%v) => %v, wants %v", testcase.param1, r, testcase.result)
		}
	}

	testcases2 := []struct {
		param1 float64
		result Timestamp
	}{
		{
			param1: 1041515675,
			result: Timestamp{Valid: true, V: time.Date(2003, time.January, 2, 13, 54, 35, 0, time.UTC)},
		},
		{
			param1: 1765094463.094,
			result: Timestamp{Valid: true, V: time.Date(2025, time.December, 7, 8, 1, 3, 94_000000, time.UTC)},
		},
	}

	for _, testcase := range testcases2 {
		var r Timestamp
		if err := json.Unmarshal([]byte(strconv.FormatFloat(testcase.param1, 'f', 3, 64)), &r); err != nil {
			t.Fatal(err)
		} else {
			if !timestampIsEqual(r, testcase.result) {
				t.Errorf("Unmarshal(%v) => %v, wants %v", testcase.param1, r, testcase.result)
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
			result: Rec{Key: "item2", Time: Timestamp{Valid: true, V: time.Date(2003, time.January, 2, 13, 54, 35, 0, time.UTC)}},
		},
		{
			param1: `{ "time": 1765094463, "key":    "item3"}`,
			result: Rec{Key: "item3", Time: Timestamp{Valid: true, V: time.Date(2025, time.December, 7, 8, 1, 3, 0, time.UTC)}},
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

func TestIntegerMerge(t *testing.T) {
	testcases := []struct {
		p1 Integer
		p2 Integer
		r  Integer
	}{
		{Integer{true, 1}, Integer{true, 100}, Integer{true, 100}},
		{Integer{true, 1}, Integer{true, 1000}, Integer{true, 1000}},
		{Integer{false, 2}, Integer{true, 7}, Integer{true, 7}},
		{Integer{false, 3}, Integer{true, 8}, Integer{true, 8}},
		{Integer{false, 4}, Integer{true, 9}, Integer{true, 9}},
		{Integer{true, 2}, Integer{false, 94}, Integer{true, 2}},
		{Integer{true, 4}, Integer{false, 102}, Integer{true, 4}},
	}

	for _, testcase := range testcases {
		s := testcase.p1
		s.Merge(testcase.p2)
		if s != testcase.r {
			t.Errorf("Integer.Merge(%v, %v) = %v, want %v", testcase.p1, testcase.p2, s, testcase.r)
		}
	}
}

func TestStringMerge(t *testing.T) {
	testcases := []struct {
		p1 String
		p2 String
		r  String
	}{
		{String{true, "a"}, String{true, "ooo"}, String{true, "ooo"}},
		{String{true, "b"}, String{true, "xxxx"}, String{true, "xxxx"}},
		{String{false, "a"}, String{true, "ooo"}, String{true, "ooo"}},
		{String{false, "b"}, String{true, "xxxxx"}, String{true, "xxxxx"}},
		{String{false, "c"}, String{true, "kkkkk"}, String{true, "kkkkk"}},
		{String{true, "a"}, String{false, "oooo"}, String{true, "a"}},
		{String{true, "b"}, String{false, "kkk"}, String{true, "b"}},
	}

	for _, testcase := range testcases {
		s := testcase.p1
		s.Merge(testcase.p2)
		if s != testcase.r {
			t.Errorf("String.Merge(%v, %v) = %v, want %v", testcase.p1, testcase.p2, s, testcase.r)
		}
	}
}

func TestBoolMerge(t *testing.T) {
	testcases := []struct {
		p1 Boolean
		p2 Boolean
		r  Boolean
	}{
		{Boolean{true, true}, Boolean{true, false}, Boolean{true, false}},
		{Boolean{true, false}, Boolean{true, true}, Boolean{true, true}},
		{Boolean{false, false}, Boolean{true, false}, Boolean{true, false}},
		{Boolean{false, true}, Boolean{true, true}, Boolean{true, true}},
		{Boolean{false, false}, Boolean{true, true}, Boolean{true, true}},
		{Boolean{true, true}, Boolean{false, false}, Boolean{true, true}},
		{Boolean{true, false}, Boolean{false, true}, Boolean{true, false}},
	}

	for _, testcase := range testcases {
		s := testcase.p1
		s.Merge(testcase.p2)
		if s != testcase.r {
			t.Errorf("Boolean.Merge(%v, %v) = %v, want %v", testcase.p1, testcase.p2, s, testcase.r)
		}
	}
}
