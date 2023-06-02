package dbhelper

import (
	"strconv"
	"strings"
)

type (
	SqlBuilder struct {
		texts []string
	}

	DynamicSqlBuilder struct {
		texts   []string
		prefix  string
		suffix  string
		joint   string
		builder *SqlBuilder
	}
)

func NewSqlBuilder(sql string) *SqlBuilder {
	return &SqlBuilder{
		texts: []string{sql},
	}
}

func (b *SqlBuilder) String() string {
	return strings.Join(b.texts, "\n")
}

func (b *SqlBuilder) Append(sql string) *SqlBuilder {
	b.texts = append(b.texts, "  "+sql)
	return b
}

func (b *SqlBuilder) append0(sql string) {
	b.texts = append(b.texts, sql)
}

func (b *SqlBuilder) AppendIf(sql string, p bool) *SqlBuilder {
	if p {
		b.texts = append(b.texts, sql)
	}
	return b
}

func (b *SqlBuilder) Dynamic(prefix, suffix, joint string) *DynamicSqlBuilder {
	return &DynamicSqlBuilder{
		texts:   []string{},
		prefix:  prefix,
		suffix:  suffix,
		joint:   joint,
		builder: b,
	}
}

func (b *SqlBuilder) Where() *DynamicSqlBuilder {
	return b.Dynamic("WHERE ", "", "\n  AND ")
}

func (b *SqlBuilder) OrderBy(sql ...string) *SqlBuilder {
	if len(sql) > 0 {
		b.append0("ORDER BY " + strings.Join(sql, ", "))
	}
	return b
}

func (b *SqlBuilder) Limit(startRowIndex, maximumRows int) *SqlBuilder {
	var sql string
	if startRowIndex <= 0 {
		sql = "LIMIT " + strconv.FormatInt(int64(maximumRows), 10)
	} else {
		sql = "LIMIT " + strconv.FormatInt(int64(maximumRows), 10) + " OFFSET " + strconv.FormatInt(int64(startRowIndex), 10)
	}
	b.append0(sql)
	return b
}

func (d *DynamicSqlBuilder) Append(sql string) *DynamicSqlBuilder {
	d.texts = append(d.texts, sql)
	return d
}

func (d *DynamicSqlBuilder) AppendIf(sql string, p bool) *DynamicSqlBuilder {
	if p {
		d.texts = append(d.texts, sql)
	}
	return d
}

func (d *DynamicSqlBuilder) End() *SqlBuilder {
	if len(d.texts) > 0 {
		d.builder.append0(d.prefix + strings.Join(d.texts, d.joint) + d.suffix)
	}
	return d.builder
}
