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

	InsertSqlBuilder struct {
		quote   string
		pos     int
		cols    []string
		params  []int
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

func (b *SqlBuilder) Inserter(quote string) *InsertSqlBuilder {
	return &InsertSqlBuilder{quote: quote, cols: []string{}, builder: b}
}

func (b *SqlBuilder) Where() *DynamicSqlBuilder {
	return b.Dynamic("WHERE\n ", "", "\n  AND ")
}

func (b *SqlBuilder) Set() *DynamicSqlBuilder {
	return b.Dynamic("SET\n ", "", ",\n  ")
}

func (b *SqlBuilder) OrderBy(sql ...string) *SqlBuilder {
	if len(sql) > 0 {
		b.append0("ORDER BY " + strings.Join(sql, ", "))
	}
	return b
}

// Limit 创建分页参数SQL。
// startRowIndex 起始行索引，从0开始。
// maximumRows 最大行数，如果小于等于0，则不添加LIMIT子句。
func (b *SqlBuilder) Limit(startRowIndex, maximumRows int) *SqlBuilder {
	var sql string
	if maximumRows <= 0 {
		return b
	} else if startRowIndex <= 0 {
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
		dp := d.prefix
		if dp != "" {
			dp = dp + " "
		}
		ds := d.suffix
		if ds != "" {
			ds = " " + ds
		}
		d.builder.append0(dp + strings.Join(d.texts, d.joint) + ds)
	}
	return d.builder
}

func (d *InsertSqlBuilder) Append(col string) *InsertSqlBuilder {
	d.pos++
	d.cols = append(d.cols, col)
	d.params = append(d.params, d.pos)

	return d
}

func (d *InsertSqlBuilder) AppendIf(col string, p bool) *InsertSqlBuilder {
	d.pos++
	if p {
		d.cols = append(d.cols, col)
		d.params = append(d.params, d.pos)
	}

	return d
}

func (d *InsertSqlBuilder) End() *SqlBuilder {
	if len(d.cols) > 0 {
		buf := make([]string, 0, len(d.cols))
		for _, col := range d.cols {
			buf = append(buf, d.quote+col+d.quote)
		}

		d.builder.append0("(" + strings.Join(buf, ",") + ")")
		posList := make([]string, 0, len(d.params))
		for _, pos := range d.params {
			posList = append(posList, ":"+strconv.Itoa(pos))
		}
		d.builder.append0("VALUES (" + strings.Join(posList, ",") + ")")
	}
	return d.builder
}
