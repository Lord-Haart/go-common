package dbhelper

import (
	"testing"

	mysql "github.com/go-sql-driver/mysql"
)

func init() {
	// db.InitDb("root:123456@(localhost:3306)/insight")
	if err := InitMySqlDb("localhost", "root", "123456", "insight"); err != nil {
		panic(err)
	}
}

func TestFormatDSN(t *testing.T) {
	cfg := mysql.Config{
		User:   "p030",
		Passwd: "Hgb123!@",
		Net:    "http",
		Addr:   "192.168.1.12:3306",
		DBName: "p030",
	}

	t.Logf("%s", cfg.FormatDSN())
}

// func TestRowsAffected(t *testing.T) {
// 	if r0 := InsertOrUpdateUser(context.TODO(), "admin", "管理员", utils.Sha256Salt("123456")); r0 != 1 {
// 		t.Errorf("InsertOrUpdateUser => %v, want 1", r0)
// 	}
// }

// func TestRowsAffected2(t *testing.T) {
// 	if r0 := db.InsertOrUpdateWorkTime(context.TODO(), "xxxxxx001", "708513b8257fd6f547c7598b4c", "1080875157529746",
// 		2, "问题修改2", time.Date(2023, time.February, 7, 0, 0, 0, 0, time.Local), time.Date(2023, time.February, 7, 0, 0, 0, 0, time.Local)); r0 != 1 {
// 		t.Errorf("InsertOrUpdateWorkTime => %v, want 1", r0)
// 	}
// }
