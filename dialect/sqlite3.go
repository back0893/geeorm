package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type sqlite3Dialect struct{}

var _ Dialect = (*sqlite3Dialect)(nil)

func (d *sqlite3Dialect) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.String:
		return "text"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Uint64, reflect.Int64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.Bool:
		return "bool"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("%s unsupported type %s", typ.Type().Name(), typ.Kind()))
}

func (d *sqlite3Dialect) TableExistSQL(tableName string) (string, []any) {
	args := []any{tableName}
	return "SELECT name FROM sqlite_master WHERE type='table' AND name = ?", args
}
