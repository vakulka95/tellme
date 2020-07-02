package postgres

import (
	"fmt"
	"strings"
)

type OrderDir string

const (
	OrderDirASC  OrderDir = "ASC"
	OrderDirDESC OrderDir = "DESC"
)

func OrderDirFromString(d string) OrderDir {
	switch strings.ToUpper(d) {
	case "ASC":
		return OrderDirASC
	case "DESC":
		return OrderDirDESC
	default:
		return OrderDirDESC
	}
}

/*
	=	equal	ARRAY[1.1,2.1,3.1]::int[] = ARRAY[1,2,3]	t
	<>	not equal	ARRAY[1,2,3] <> ARRAY[1,2,4]	t
	<	less than	ARRAY[1,2,3] < ARRAY[1,2,4]	t
	>	greater than	ARRAY[1,4,3] > ARRAY[1,2,4]	t
	<=	less than or equal	ARRAY[1,2,3] <= ARRAY[1,2,3]	t
	>=	greater than or equal	ARRAY[1,4,3] >= ARRAY[1,4,3]	t
	@>	contains	ARRAY[1,4,3] @> ARRAY[3,1]	t
	<@	is contained by
	any() any elements of array
*/

type BuildOperatorExpression func(field string, num int) string

func OperatorEqual(field string, num int) string {
	return fmt.Sprintf("%s=$%d", field, num)
}

func OperatorNotEqual(field string, num int) string {
	return fmt.Sprintf("%s<>$%d", field, num)
}

func OperatorLessThen(field string, num int) string {
	return fmt.Sprintf("%s<$%d", field, num)
}

func OperatorGreaterEqual(field string, num int) string {
	return fmt.Sprintf("%s>$%d", field, num)
}

func OperatorLessThenOrEqual(field string, num int) string {
	return fmt.Sprintf("%s<=$%d", field, num)
}

func OperatorGreaterEqualOrEqual(field string, num int) string {
	return fmt.Sprintf("%s>=$%d", field, num)
}

func OperatorContains(field string, num int) string {
	return fmt.Sprintf("%s@>$%d", field, num)
}

func OperatorIsContainedBy(field string, num int) string {
	return fmt.Sprintf("%s<@$%d", field, num)
}

func OperatorAny(field string, num int) string {
	return fmt.Sprintf("%s=any($%d)", field, num)
}

func OperatorLike(field string, num int) string {
	return fmt.Sprintf(` %s LIKE  '%%' || $%d || '%%' `, field, num)
}

type QueryBuilder interface {
	Select(columns ...string) QueryBuilder
	From(table string) QueryBuilder
	Where(expression ...Expression) QueryBuilder
	OrderBy(field string) QueryBuilder
	OrderDir(direction OrderDir) QueryBuilder
	Limit(limit int) QueryBuilder
	Offset(offset int) QueryBuilder
	Build() (sql string, args []interface{})
	BuildCount() (sql string, args []interface{})
}

type Value interface {
	Valid() bool
	Arg() interface{}
}

type Expression interface {
	Build(n int) string
	Arg() interface{}
	Valid() bool
}
