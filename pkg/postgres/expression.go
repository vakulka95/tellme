package postgres

import (
	"fmt"
	"strings"
)

type expression struct {
	field string
	value Value
	build BuildOperatorExpression
}

func NewExpression(field string, value Value, build BuildOperatorExpression) *expression {
	return &expression{field: field, value: value, build: build}
}

func (e *expression) Build(n int) string {
	return fmt.Sprintf("( %s )", e.build(e.field, n))
}

func (e *expression) Arg() interface{} {
	return e.value.Arg()
}

func (e *expression) Valid() bool {
	return e.value.Valid()
}

type orExpression struct {
	fields []string
	value  Value
	build  BuildOperatorExpression
}

func NewOrExpression(value Value, build BuildOperatorExpression, fields ...string) *orExpression {
	return &orExpression{fields: fields, value: value, build: build}
}

func (e *orExpression) Build(n int) string {
	fields := make([]string, 0, len(e.fields))
	for _, f := range e.fields {
		fields = append(fields, e.build(f, n))
	}
	return fmt.Sprintf("( %s )", strings.Join(fields, " OR "))
}

func (e *orExpression) Arg() interface{} {
	return e.value.Arg()
}

func (e *orExpression) Valid() bool {
	return e.value.Valid()
}
