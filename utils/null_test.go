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

func TestTimestampJSON(t *testing.T) {
	testcases1 := []struct {
		param1 Timestamp
		result int64
	}{
		{
			param1: Timestamp{Valid: true, V: time.Date(2003, time.January, 2, 13, 54, 35, 0, time.UTC)},
			result: 1041515675,
		},
		{
			param1: Timestamp{Valid: true, V: time.Date(2025, time.December, 7, 8, 1, 3, 0, time.UTC)},
			result: 1765094463,
		},
	}

	for _, testcase := range testcases1 {
		if b, err := json.Marshal(testcase.param1); err != nil {
			t.Fatal(err)
		} else if r, err := strconv.ParseInt(string(b), 10, 64); err != nil {
			t.Fatal(err)
		} else if r != testcase.result {
			t.Errorf("Marshal(%v) => %v, wants %v", testcase.param1, r, testcase.result)
		}
	}

	testcases2 := []struct {
		param1 int64
		result Timestamp
	}{
		{
			param1: 1041515675,
			result: Timestamp{Valid: true, V: time.Date(2003, time.January, 2, 13, 54, 35, 0, time.UTC)},
		},
		{
			param1: 1765094463,
			result: Timestamp{Valid: true, V: time.Date(2025, time.December, 7, 8, 1, 3, 0, time.UTC)},
		},
	}

	for _, testcase := range testcases2 {
		var r Timestamp
		if err := json.Unmarshal([]byte(strconv.FormatInt(testcase.param1, 10)), &r); err != nil {
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
