package session

import (
	"database/sql"
	"fmt"
	"geemod/dialect"
	"geemod/log"
	"geemod/schema"
	"reflect"
	"strings"
)

type Session struct {
	db        *sql.DB
	builder   strings.Builder
	sqlValues []any

	dialect  dialect.Dialect
	refTable *schema.Schema
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		builder: strings.Builder{},
		dialect: dialect,
	}
}

func (s *Session) Clear() {
	s.builder.Reset()
	s.sqlValues = nil
}
func (s *Session) Raw(sql string, values ...any) *Session {
	s.builder.WriteString(sql)
	s.builder.WriteString(" ")
	s.sqlValues = append(s.sqlValues, values...)
	return s
}

func (s *Session) Exec() (sql.Result, error) {
	defer s.Clear()
	log.Infof("sql:%s", s.builder.String())
	result, err := s.db.Exec(s.builder.String(), s.sqlValues...)
	if err != nil {
		log.Errorf("exec failed, err:%v", err)
		return nil, err
	}
	return result, nil
}

func (s *Session) QueryRow() (*sql.Row, error) {
	defer s.Clear()
	log.Infof("sql:%s", s.builder.String())
	row := s.db.QueryRow(s.builder.String(), s.sqlValues...)
	if row.Err() != nil {
		log.Errorf("query failed, err:%v", row.Err())
		return nil, row.Err()
	}
	return row, nil
}

func (s *Session) QueryRows() (*sql.Rows, error) {
	defer s.Clear()
	log.Infof("sql:%s", s.builder.String())
	rows, err := s.db.Query(s.builder.String(), s.sqlValues...)
	if err != nil {
		log.Errorf("query failed, err:%v", err)
		return nil, err
	}
	return rows, nil
}

func (s *Session) Model(model any) *Session {
	if s.refTable == nil || reflect.TypeOf(model) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(model, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Errorf("refTable is nil")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	schmea := s.refTable
	cols := make([]string, 0)
	for _, field := range schmea.Fields {
		cols = append(cols, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	sql := fmt.Sprintf("create table %s (%s);", schmea.Name, strings.Join(cols, ","))
	_, err := s.Raw(sql).Exec()
	if err != nil {
		log.Errorf("create table failed, err:%v", err)
		return err
	}
	return nil
}
func (s *Session) DropTable() error {
	sql := fmt.Sprintf("drop table if exists %s;", s.refTable.Name)
	_, err := s.Raw(sql).Exec()
	if err != nil {
		log.Errorf("create table failed, err:%v", err)
		return err
	}
	return nil
}
func (s *Session) HasTable() (bool, error) {
	sql, values := s.dialect.TableExistSQL(s.refTable.Name)
	row, err := s.Raw(sql, values...).QueryRow()
	if err != nil {
		log.Errorf("query table failed, err:%v", err)
		return false, err
	}
	tableName := ""
	row.Scan(&tableName)
	return tableName == s.refTable.Name, nil
}
