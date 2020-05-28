package postgres

type OrderDir string

const (
	OrderDirASC  OrderDir = "ASC"
	OrderDirDESC OrderDir = "DESC"
)

type QueryBuilder interface {
	Select(columns ...string) QueryBuilder
	From(table string) QueryBuilder
	Where(expression ...WhereExpression) QueryBuilder
	OrderBy(field string) QueryBuilder
	OrderDir(direction OrderDir) QueryBuilder
	Limit(limit int) QueryBuilder
	Offset(offset int) QueryBuilder
	Build() (sql string, args []interface{})
	BuildCount() (sql string, args []interface{})
}

type WhereExpression interface {
	Build(numb int) string
	GetArg() interface{}
}
