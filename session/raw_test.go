package session

import (
	"database/sql"
	"geemod/dialect"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var TestDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	TestDB, err = sql.Open("sqlite3", "../gee.db")
	if err != nil {
		panic(err)
	}
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSeesion() *Session {
	dialect, ok := dialect.GetDialect("sqlite3")
	if !ok {
		panic("dialect sqlite3 not found")
	}
	return New(TestDB, dialect)
}

func TestSessionExec(t *testing.T) {
	s := NewSeesion()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	if count, err := result.RowsAffected(); err != nil || count != 2 {
		t.Fatal("expect 2, but got", count)
	}
}

func TestSessionQueryRaws(t *testing.T) {
	s := NewSeesion()
	rows, err := s.Raw("select Name from User").QueryRows()
	if err != nil {
		t.Fatal("expect nil, but got", err)
	}
	for rows.Next() {
		var name string
		rows.Scan(&name)
		t.Logf("get name:%s", name)
	}
}

type User struct {
	ID   int `geeorm:"Primary Key"`
	Name string
}

func TestSessionCreate(t *testing.T) {
	s := NewSeesion()
	s.Model(&User{})
	if err := s.CreateTable(); err != nil {
		t.Fatal(err.Error())
	}
	if has, err := s.HasTable(); err != nil || !has {
		t.Fatal("expect true, but got", has)
	}
}
func TestSessionDrop(t *testing.T) {
	s := NewSeesion()
	s.Model(&User{})
	if err := s.DropTable(); err != nil {
		t.Fatal(err.Error())
	}
	if has, err := s.HasTable(); err != nil || has {
		t.Fatal("expect false, but got", has)
	}
}
