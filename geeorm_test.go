package geeorm

import (
	"errors"
	"geemod/session"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDb(t *testing.T) *Engine {
	e, err := NewEngine("sqlite3", "./gee.db")
	if err != nil {
		t.Fatal("expect nil, but got", err)
	}
	return e
}

func TestNewEngine(t *testing.T) {
	engine := OpenDb(t)
	defer engine.Close()
}

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	// t.Run("commit", func(t *testing.T) {
	// 	transactionCommit(t)
	// })
}

func transactionRollback(t *testing.T) {
	engine := OpenDb(t)
	defer engine.Close()
	s := engine.Session()
	_ = s.Model(&User{}).DropTable()
	err := engine.Transaction(func(s *session.Session) error {
		_ = s.Model(&User{}).CreateTable()
		_, err := s.Insert(&User{"Tom", 18})
		if err != nil {
			return err
		}
		return errors.New("err")
	})
	hasTable, _ := s.HasTable()
	if err == nil || hasTable {
		t.Fatal("failed to rollback")
	}
}

func TestEngine_Migrate(t *testing.T) {
	engine := OpenDb(t)
	defer engine.Close()
	s := engine.Session()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text PRIMARY KEY, XXX integer);").Exec()
	_, _ = s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	engine.Migrate(&User{})

	rows, _ := s.Raw("SELECT * FROM User").QueryRows()
	columns, _ := rows.Columns()
	if !reflect.DeepEqual(columns, []string{"Name", "Age"}) {
		t.Fatal("Failed to migrate table User, got columns", columns)
	}
}
