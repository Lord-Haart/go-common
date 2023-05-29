package dbhelper

import (
	"context"
	"testing"

	"github.com/Lord-Haart/go-common/utils"
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

func TestRowsAffected(t *testing.T) {
	if r0, err := Exec[int](context.TODO(), "INSERT INTO user (user_name, nick_name, password) VALUE (:1, :2, :3)", "admin2", "管理员", utils.Sha256Salt("123456")); err != nil {
		t.Fatal(err)
	} else if r0 != 1 {
		t.Errorf("InsertOrUpdateUser => %v, want 1", r0)
	}
}

func TestUpdateWithTx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	ctx = BeginTx(ctx, false)

	defer CloseTx(ctx)

	if r0, err := Exec[int](ctx, "INSERT INTO user (user_name, nick_name, password) VALUE (:1, :2, :3)", "admin99", "管理员", utils.Sha256Salt("123456")); err != nil {
		t.Fatal(err)
	} else if r0 != 1 {
		t.Errorf("InsertOrUpdateUser => %v, want 1", r0)
	}

	CloseTx(ctx)
}

// func TestRowsAffected2(t *testing.T) {
// 	if r0 := db.InsertOrUpdateWorkTime(context.TODO(), "xxxxxx001", "708513b8257fd6f547c7598b4c", "1080875157529746",
// 		2, "问题修改2", time.Date(2023, time.February, 7, 0, 0, 0, 0, time.Local), time.Date(2023, time.February, 7, 0, 0, 0, 0, time.Local)); r0 != 1 {
// 		t.Errorf("InsertOrUpdateWorkTime => %v, want 1", r0)
// 	}
// }
