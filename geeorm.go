package geeorm

import (
	"database/sql"
	"errors"
	"geemod/dialect"
	"geemod/session"
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
