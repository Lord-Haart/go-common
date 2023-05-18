package utils

import (
	"testing"
	"time"
)

func TestIsWorkingDay(t *testing.T) {
	testcases := []struct {
		Param1 time.Time
		Result bool
	}{
		{Param1: time.Date(2023, time.April, 15, 0, 0, 0, 0, time.Local), Result: false}, /*周六*/
		{Param1: time.Date(2023, time.April, 16, 0, 0, 0, 0, time.Local), Result: false}, /*周日*/
		{Param1: time.Date(2023, time.April, 22, 0, 0, 0, 0, time.Local), Result: false}, /*周六*/
		{Param1: time.Date(2023, time.April, 23, 0, 0, 0, 0, time.Local), Result: true},  /*五一调休，工作日*/
		{Param1: time.Date(2023, time.April, 24, 0, 0, 0, 0, time.Local), Result: true},  /*周一*/
		{Param1: time.Date(2023, time.April, 27, 0, 0, 0, 0, time.Local), Result: true},  /*周四*/
		{Param1: time.Date(2023, time.April, 28, 0, 0, 0, 0, time.Local), Result: true},  /*周五*/
		{Param1: time.Date(2023, time.April, 29, 0, 0, 0, 0, time.Local), Result: false}, /*五一调休*/
		{Param1: time.Date(2023, time.May, 1, 0, 0, 0, 0, time.Local), Result: false},    /*五一调休*/
		{Param1: time.Date(2023, time.May, 6, 0, 0, 0, 0, time.Local), Result: true},     /*五一调休，工作日*/
	}

	for _, testcase := range testcases {
		if r := IsWorkingDay(testcase.Param1); r != testcase.Result {
			t.Errorf("IsWorkingDay(%#v) => %#v, wants %#v", testcase.Param1, r, testcase.Result)
		}
	}
}
