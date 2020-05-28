package postgres

import (
	"bytes"
	"fmt"
	"strings"
)

type queryBuilder struct {
	columns  []string
	table    string
	where    []WhereExpression
	orderBy  string
	orderDir OrderDir
	limit    int
	offset   int
}

func NewQueryBuilder() QueryBuilder {
	return &queryBuilder{}
}

func (b *queryBuilder) Select(columns ...string) QueryBuilder {
	b.columns = columns
	return b
}

func (b *queryBuilder) From(table string) QueryBuilder {
	b.table = table
	return b
}

func (b *queryBuilder) Where(expressions ...WhereExpression) QueryBuilder {
	b.where = expressions
	return b
}

func (b *queryBuilder) OrderBy(field string) QueryBuilder {
	b.orderBy = field
	return b
}

func (b *queryBuilder) OrderDir(direction OrderDir) QueryBuilder {
	b.orderDir = direction
	return b
}

func (b *queryBuilder) Limit(limit int) QueryBuilder {
	b.limit = limit
	return b
}

func (b *queryBuilder) Offset(offset int) QueryBuilder {
	b.offset = offset
	return b
}

func (b *queryBuilder) Build() (string, []interface{}) {
	buff := &bytes.Buffer{}

	// select columns
	buff.WriteString("SELECT ")
	buff.WriteString(strings.Join(b.columns, ", "))

	// from table
	buff.WriteString(fmt.Sprintf(" FROM %s ", b.table))

	// where
	args := make([]interface{}, 0, len(b.where))

	if len(b.where) > 0 {
		buff.WriteString(" WHERE ")

		where := make([]string, 0, len(b.where))
		for i, exp := range b.where {
			args = append(args, exp.GetArg())
			where = append(where, exp.Build(i+1))
		}

		buff.WriteString(strings.Join(where, " AND "))
	}

	// ordering
	if b.orderBy != "" {
		exp := fmt.Sprintf(" ORDER BY %s ORDER DIR %s ", b.orderBy, b.orderDir)
		buff.WriteString(exp)
	}

	// limits
	if b.limit != 0 {
		exp := fmt.Sprintf(" LIMIT %d OFFSET %d", b.limit, b.offset)
		buff.WriteString(exp)
	}

	return buff.String(), args
}

func (b *queryBuilder) BuildCount() (string, []interface{}) {
	buff := &bytes.Buffer{}

	// select columns
	buff.WriteString("SELECT ")
	buff.WriteString(" count(id) ")

	// from table
	buff.WriteString(fmt.Sprintf(" FROM %s ", b.table))

	// where
	args := make([]interface{}, 0, len(b.where))

	if len(b.where) > 0 {
		buff.WriteString(" WHERE ")

		where := make([]string, 0, len(b.where))
		for i, exp := range b.where {
			args = append(args, exp.GetArg())
			where = append(where, exp.Build(i+1))
		}

		buff.WriteString(strings.Join(where, " AND "))
	}

	return buff.String(), args
}
