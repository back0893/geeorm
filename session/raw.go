package session

import (
	"database/sql"
	"errors"
	"fmt"
	"geemod/clause"
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

	clause clause.Clause
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
	s.clause = clause.Clause{}
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

func (s *Session) Insert(values ...any) (int64, error) {
	var vars []any
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		CallHookMethod(value, func(method IBeforeInsert) {
			method.BeforeInsert(s)
		})
		vars = append(vars, s.refTable.RecordValues(value))
	}
	s.clause.Set(clause.VALUES, vars...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	for _, value := range values {
		CallHookMethod(value, func(method IAfterInsert) {
			method.AfterInsert(s)
		})
	}

	return result.RowsAffected()
}

func (s *Session) Find(value any) error {
	CallHookMethod(s.refTable.Model, func(method IBeforeQuery) {
		method.BeforeQuery(s)
	})
	//select $fileds from $tableName
	descSlice := reflect.ValueOf(value)
	if descSlice.Kind() != reflect.Ptr {
		return errors.New("not pointer")
	}
	descSlice = descSlice.Elem()
	if descSlice.Kind() != reflect.Slice {
		return errors.New("not slice")
	}
	descType := descSlice.Type().Elem()
	srcdescType := descType
	if descType.Kind() == reflect.Ptr {
		descType = descType.Elem()
	}
	schema := s.Model(reflect.New(descType).Elem().Interface()).RefTable()

	s.clause.Set(clause.SELECT, schema.Name, schema.FieldNames)

	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}
	for rows.Next() {
		record := reflect.New(descType).Elem()
		var values []any
		for _, filed := range schema.FieldNames {
			values = append(values, record.FieldByName(filed).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		if srcdescType.Kind() == reflect.Ptr {
			descSlice.Set(reflect.Append(descSlice, record.Addr()))
			CallHookMethod(record.Addr().Interface(), func(method IAfterQuery) {
				method.AfterQuery(s)
			})
		} else {
			descSlice.Set(reflect.Append(descSlice, record))
			CallHookMethod(record.Interface(), func(method IAfterQuery) {
				method.AfterQuery(s)
			})
		}

	}
	return rows.Close()
}
func (s *Session) First(value any) error {
	dst := reflect.ValueOf(value)
	if dst.Kind() != reflect.Ptr {
		return errors.New("not pointer")
	}
	descSlice := reflect.New(reflect.SliceOf(dst.Type())).Elem()
	if err := s.Limit(1).Find(descSlice.Addr().Interface()); err != nil {
		return err
	}
	if descSlice.Len() == 0 {
		return errors.New("no result")
	}
	dst.CanSet()
	dst.Elem().Set(descSlice.Index(0).Elem())
	return nil
}
