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

func TestSqlBuilder2(t *testing.T) {
	b0 := NewSqlBuilder("UPDATE foot").
		Set().
		AppendIf("b = :2", true).
		AppendIf("c = :3", false).
		AppendIf("e = :4", true).
		Append("k = :7").
		End().
		Where().
		Append("d = :5").
		End()

	t.Logf("sql: %s", b0)
}

func TestSqlBuilder3(t *testing.T) {
	b0 := NewSqlBuilder("INSERT INTO foo").
		Inserter("`").
		Append("id").
		AppendIf("name", false).
		Append("create_time").
		End()

	t.Logf("sql: %s", b0)
}
