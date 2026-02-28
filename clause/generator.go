package clause

import (
	"fmt"
	"strings"
)

// 生成sql的各个部分
type generator func(values ...any) (string, []any)

var generators = map[Type]generator{}

func init() {
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
		vars = value.([]any)
		if bindStr == "" {
			bindStr = genBindVars(len(vars))
		}
		sqlBuilder.WriteString("(")
		sqlBuilder.WriteString(bindStr)
		sqlBuilder.WriteString(")")
		if idx+1 < len(values) {
			sqlBuilder.WriteString(",")
		}
		vars = append(vars, vars...)
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
	sqlBuilder := strings.Builder{}
	sqlBuilder.WriteString("order by ")
	for i := 0; i < len(values); i += 2 {
		sqlBuilder.WriteString(values[i].(string))
		sqlBuilder.WriteString(values[i+1].(string))
	}
	return sqlBuilder.String(), nil
}

func _where(values ...any) (string, []any) {
	//where $cond
	sqlBuilder := strings.Builder{}
	sqlBuilder.WriteString("where ")
	for i := 0; i < len(values); i += 2 {
		sqlBuilder.WriteString(values[i].(string))
		sqlBuilder.WriteString(values[i+1].(string))
	}
	return sqlBuilder.String(), nil
}
