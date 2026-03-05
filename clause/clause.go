package clause

import "strings"

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	OFFSET
	ORDERBY
	WHERE
	UPDATE
	DELETE
	COUNT
)

type Clause struct {
	sql     map[Type]string
	sqlVals map[Type][]any
}

func (clause *Clause) Set(type_ Type, values ...any) {
	if clause.sql == nil {
		clause.sql = make(map[Type]string)
		clause.sqlVals = make(map[Type][]any)
	}
	sql, vars := generators[type_](values...)
	clause.sql[type_] = sql
	clause.sqlVals[type_] = vars
}
func (clause *Clause) Build(orders ...Type) (string, []any) {
	var sqls []string
	var vars []any
	for _, order := range orders {
		if sql, ok := clause.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, clause.sqlVals[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
