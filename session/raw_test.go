package session

import (
	"database/sql"
	"geemod/dialect"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var TestDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	TestDB, err = sql.Open("sqlite3", "./gee.db")
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
	Age  int
}

func (u *User) BeforeInsert(s *Session) {
	log.Printf("insert id:%d name:%s,age:%d", u.ID, u.Name, u.Age)
}
func (u *User) AfterInsert(s *Session) {
	log.Printf("after insert id:%d", u.ID)
}
func (u *User) BeforeQuery(s *Session) {
	log.Printf("before query id:%d", u.ID)
}
func (u *User) AfterQuery(s *Session) {
	log.Printf("after query id:%d", u.ID)
}
func (u *User) BeforeDelete(s *Session) {
	log.Printf("before delete id:%d", u.ID)
}
func (u *User) AfterDelete(s *Session) {
	log.Printf("after delete id:%d", u.ID)
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

func testRecordInit(t *testing.T) *Session {
	var (
		user1 = &User{1, "Tom", 18}
		user2 = &User{2, "Sam", 25}
	)
	t.Helper()
	s := NewSeesion()
	s.Model(&User{})
	err1 := s.DropTable()
	err2 := s.CreateTable()
	_, err3 := s.Insert(user1, user2)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("failed init test records")
	}
	return s
}

func TestInsert(t *testing.T) {
	s := testRecordInit(t)
	var user3 = &User{3, "Jack", 25}
	affected, err := s.Insert(user3)
	if err != nil || affected != 1 {
		t.Fatal("failed to create record")
	}
}

func TestFind(t *testing.T) {
	s := testRecordInit(t)
	var users []*User
	if err := s.Find(&users); err != nil {
		t.Fatalf("failed to find records %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("failed to find records %d", len(users))
	}
	if users[0].Name != "Tom" || users[0].Age != 18 || users[0].ID != 1 {
		t.Fatalf("failed to find records %v", users[0])
	}
	if users[1].Name != "Sam" || users[1].Age != 25 || users[1].ID != 2 {
		t.Fatalf("failed to find records %v", users[1])
	}

}

func TestFirst(t *testing.T) {
	s := testRecordInit(t)
	var user User
	if err := s.First(&user); err != nil {
		t.Fatalf("failed to find records %v", err)
	}
	if user.Name != "Tom" || user.Age != 18 || user.ID != 1 {
		t.Fatalf("failed to find records %v", user)
	}
}

func TestOrderBy(t *testing.T) {
	s := testRecordInit(t)
	var users []*User
	if err := s.OrderBy("ID", true).Limit(2).Find(&users); err != nil {
		t.Fatalf("failed to find records %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("failed to find records %d", len(users))
	}
	if users[0].Name != "Sam" || users[0].Age != 25 || users[0].ID != 2 {
		t.Fatalf("failed to find records %v", users[0])
	}
	if users[1].Name != "Tom" || users[1].Age != 18 || users[1].ID != 1 {
		t.Fatalf("failed to find records %v", users[1])
	}

}
func TestCount(t *testing.T) {
	s := testRecordInit(t)
	count, err := s.Model(&User{}).Where("ID=?", 1).Count("*")
	if err != nil {
		t.Fatalf("failed to find records %v", err)
	}
	if count != 1 {
		t.Fatalf("failed to find records %d", count)
	}
}
func TestDelete(t *testing.T) {
	s := testRecordInit(t)
	count, err := s.Model(&User{}).Where("ID=?", 1).Delete()
	if err != nil {
		t.Fatalf("failed to find records %v", err)
	}
	if count != 1 {
		t.Fatalf("failed to find records %d", count)
	}
}
