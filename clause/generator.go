package clause

import (
	"fmt"
	"strings"
)

// 生成sql的各个部分
type generator func(values ...any) (string, []any)

var generators = map[Type]generator{}

func init() {
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[OFFSET] = _offset
	generators[ORDERBY] = _orderby
	generators[WHERE] = _where
}

func _insert(values ...any) (string, []any) {
	//insert into $table_name ($fields)
	tableName := values[0].(string)
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("insert into %s (%s)", tableName, fields), nil
}
func genBindVars(n int) string {
	return strings.Join(strings.Split(strings.Repeat("?", n), ""), ",")
}
func _values(values ...any) (string, []any) {
	//values ($values)
	sqlBuilder := strings.Builder{}
	vars := make([]any, 0)
	bindStr := ""
	sqlBuilder.WriteString("values ")
	for idx, value := range values {
		tmpVars := value.([]any)
		if bindStr == "" {
			bindStr = genBindVars(len(tmpVars))
		}
		sqlBuilder.WriteString("(")
		sqlBuilder.WriteString(bindStr)
		sqlBuilder.WriteString(")")
		if idx+1 < len(values) {
			sqlBuilder.WriteString(",")
		}
		vars = append(vars, tmpVars...)
	}
	return sqlBuilder.String(), vars
}

// tableName []string
func _select(values ...any) (string, []any) {
	//select $fields from $table_name
	tableName := values[0].(string)
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("select %s from %s", fields, tableName), nil
}

func _limit(values ...any) (string, []any) {
	return "limit ?", values
}
func _offset(values ...any) (string, []any) {
	return "offset ?", values
}

func _orderby(values ...any) (string, []any) {
	//order by $fields $desc
	return fmt.Sprintf("order by %s", values[0]), nil
}

func _where(values ...any) (string, []any) {
	//where $cond
	return fmt.Sprintf("where %s", values[0]), values[1:]
}
