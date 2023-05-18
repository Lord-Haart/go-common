package utils

import "time"

var (
	workingDays = NewSet(
		20230423, 20230506, /*五一*/
		20230625,           /*端午*/
		20231007, 20231008, /*国庆*/
	)

	restDays = NewSet(
		20230405,                                         /*清明*/
		20230429, 20230430, 20230501, 20230502, 20230503, /*五一*/
		20230622, 20230623, 20230624, /*端午*/
		20230929, 20230930, 20231001, 20231002, 20231003, 20231004, 20231005, 20231006, /*国庆*/
	)
)

// IsWorkingDay 判断指定日是否是工作日。
func IsWorkingDay(t time.Time) bool {
	y, m, d := t.Date()
	td := y*10000 + int(m)*100 + d

	// 如果是调休的工作日，那么就是工作日。
	// 否则如果是节假日，那么就是休息日。
	// 否则如果是周六周日就是休息日。
	// 否则是工作日。
	if workingDays.Contains(td) {
		return true
	} else if restDays.Contains(td) {
		return false
	} else {
		wd := t.Weekday()
		return wd != time.Sunday && wd != time.Saturday
	}
}
