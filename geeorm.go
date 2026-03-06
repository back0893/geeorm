package geeorm

import (
	"database/sql"
	"errors"
	"fmt"
	"geemod/dialect"
	"geemod/session"
	"strings"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (*Engine, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	dialect, ok := dialect.GetDialect(driver)
	if !ok {
		return nil, errors.New(driver + "dialect not found")
	}
	return &Engine{db: db, dialect: dialect}, nil
}
func (e *Engine) Close() error {
	return e.db.Close()
}
func (e *Engine) Session() *session.Session {
	return session.New(e.db, e.dialect)
}

type TxFunc func(*session.Session) error

func (e *Engine) Transaction(fn TxFunc) error {
	s := e.Session()
	if err := s.Begin(); err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			s.Rollback()
		}
	}()
	err := fn(s)
	if err != nil {
		s.Rollback()
		return err
	}
	s.Commit()
	return nil
}

func difference[T comparable](a, b []T) []T {
	m := make(map[T]struct{})
	for _, v := range b {
		m[v] = struct{}{}
	}
	var diff []T
	for _, v := range a {
		if _, ok := m[v]; !ok {
			diff = append(diff, v)
		}
	}
	return diff
}

func (e *Engine) Migrate(value any) error {
	err := e.Transaction(func(s *session.Session) error {
		s = s.Model(value)
		has, err := s.HasTable()
		if err != nil {
			return err
		}
		if !has {
			return s.CreateTable()
		}
		table := s.RefTable()
		rows, err := s.Raw(fmt.Sprintf("SELECT * FROM %s limit 1", table.Name)).QueryRows()
		if err != nil {
			return err
		}
		defer rows.Close()
		columns, err := rows.Columns()
		if err != nil {
			return err
		}
		delFields := difference(columns, table.FieldNames)
		addFields := difference(table.FieldNames, columns)

		for _, addField := range addFields {
			f := s.RefTable().GetField(addField)
			s.Raw(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, addField, f.Type))
			_, err = s.Exec()
			if err != nil {
				return err
			}
		}
		if len(delFields) > 0 {
			tmp := "tmp_" + table.Name
			fieldStr := strings.Join(table.FieldNames, ", ")
			s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmp, fieldStr, table.Name))
			s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
			s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name))
			_, err = s.Exec()
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
