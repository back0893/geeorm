package dialect

import "reflect"

var dialectmap = map[string]Dialect{}

func init() {
	RegisterDialect("sqlite3", &sqlite3Dialect{})
}

type Dialect interface {
	DataTypeOf(typ reflect.Value) string
	TableExistSQL(tableName string) (string, []any)
}

func RegisterDialect(dirverName string, dialect Dialect) {
	dialectmap[dirverName] = dialect
}

func GetDialect(dirverName string) (dialect Dialect, ok bool) {
	dialect, ok = dialectmap[dirverName]
	return
}
