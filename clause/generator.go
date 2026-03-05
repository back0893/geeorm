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
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
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

func _update(values ...any) (string, []any) {
	//update $table_name set $col1=$val1,$col2=$val2 ...
	tableName := values[0].(string)
	var sqlBuiler strings.Builder
	sqlBuiler.WriteString(fmt.Sprintf("update %s set ", tableName))
	var vars []any
	updates := make([]string, 0)
	for key, val := range values[1].(map[string]any) {
		updates = append(updates, fmt.Sprintf("%s=?", key))
		vars = append(vars, val)
	}
	sqlBuiler.WriteString(strings.Join(updates, ","))
	return sqlBuiler.String(), vars
}

func _delete(values ...any) (string, []any) {
	//deleet $tableName
	return fmt.Sprintf("delete from %s", values[0].(string)), nil
}
func _count(values ...any) (string, []any) {
	//select count($field) from $table_name
	return _select(values[0].(string), []string{fmt.Sprintf("count(%s)", values[1].(string))})
}
