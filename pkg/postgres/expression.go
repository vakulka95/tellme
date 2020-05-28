package postgres

import "fmt"

/*
=	equal	ARRAY[1.1,2.1,3.1]::int[] = ARRAY[1,2,3]	t
<>	not equal	ARRAY[1,2,3] <> ARRAY[1,2,4]	t
<	less than	ARRAY[1,2,3] < ARRAY[1,2,4]	t
>	greater than	ARRAY[1,4,3] > ARRAY[1,2,4]	t
<=	less than or equal	ARRAY[1,2,3] <= ARRAY[1,2,3]	t
>=	greater than or equal	ARRAY[1,4,3] >= ARRAY[1,4,3]	t
@>	contains	ARRAY[1,4,3] @> ARRAY[3,1]	t
<@	is contained by
*/

// equal exp
type equal struct {
	fieldName string
	arg       interface{}
}

func NewEqual(fieldName string, arg interface{}) WhereExpression {
	return equal{fieldName: fieldName, arg: arg}
}

func (e equal) Build(n int) string {
	return fmt.Sprintf("%s=$%d", e.fieldName, n)
}

func (e equal) GetArg() interface{} {
	return e.arg
}

// notEqual exp
type notEqual struct {
	fieldName string
	arg       interface{}
}

func NewNotEqual(fieldName string, arg interface{}) WhereExpression {
	return notEqual{fieldName: fieldName, arg: arg}
}

func (e notEqual) Build(n int) string {
	return fmt.Sprintf("%s<>$%d", e.fieldName, n)
}

func (e notEqual) GetArg() interface{} {
	return e.arg
}

// in exp
type in struct {
	fieldName string
	arg       interface{}
}

func NewIn(fieldName string, arg interface{}) WhereExpression {
	return in{fieldName: fieldName, arg: arg}
}

func (e in) Build(n int) string {
	return fmt.Sprintf("%s in $%d", e.fieldName, n)
}

func (e in) GetArg() interface{} {
	return e.arg
}
