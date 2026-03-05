package session

import (
	"errors"
	"fmt"
	"geemod/clause"
)

func (s *Session) Update(kv ...any) (int64, error) {
	if len(kv) == 0 {
		return 0, errors.New("empty")
	}
	CallHookMethod(s.refTable.Model, func(method IBeforeUpdate) {
		method.BeforeUpdate(s)
	})
	v, ok := kv[0].(map[string]any)
	if !ok {
		v := make(map[string]any)
		for i := 0; i < len(kv); i += 2 {
			v[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.refTable.Name, v)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	CallHookMethod(s.refTable.Model, func(method IAfterUpdate) {
		method.AfterUpdate(s)
	})
	return result.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	CallHookMethod(s.refTable.Model, func(method IBeforeDelete) {
		method.BeforeDelete(s)
	})
	s.clause.Set(clause.DELETE, s.refTable.Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	CallHookMethod(s.refTable.Model, func(method IAfterDelete) {
		method.AfterDelete(s)
	})
	return result.RowsAffected()
}

func (s *Session) Count(v string) (uint64, error) {
	s.clause.Set(clause.COUNT, s.refTable.Name, v)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row, err := s.Raw(sql, vars...).QueryRow()
	if err != nil {
		return 0, err
	}
	var count uint64
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Session) Limit(limit int) *Session {
	s.clause.Set(clause.LIMIT, limit)
	return s
}
func (s *Session) Offset(offset int) *Session {
	s.clause.Set(clause.OFFSET, offset)
	return s
}

func (s *Session) Where(sql string, vars ...any) *Session {
	s.clause.Set(clause.WHERE, append([]any{sql}, vars...)...)
	return s
}

func orderStr(isDesc bool) string {
	if isDesc {
		return "desc"
	}
	return "asc"
}

func (s *Session) OrderBy(filed string, isDesc bool) *Session {
	s.clause.Set(clause.ORDERBY, fmt.Sprintf("%s %s", filed, orderStr(isDesc)))
	return s
}
