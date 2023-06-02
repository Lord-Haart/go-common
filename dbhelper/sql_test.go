package dbhelper

import "testing"

func TestSqlBuilder1(t *testing.T) {
	b0 := NewSqlBuilder("SELECT *").
		Append("FROM bb").
		AppendIf("INNER JOIN cc", true).
		AppendIf("LEFT JOIN dd", false).
		Where().
		AppendIf("b = :2", true).
		AppendIf("c = :3", false).
		Append("d = :5").
		End().
		OrderBy("name DESC", "create_time").
		Limit(0, 10)

	t.Logf("sql: %s", b0)
}
