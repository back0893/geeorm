package schema

import (
	"geemod/dialect"
	"testing"
)

func TestParse(t *testing.T) {
	type User struct {
		ID   int `geeorm:"PRIMARY KEY"`
		Name string
	}

	dialect, ok := dialect.GetDialect("sqlite3")
	if !ok {
		panic("sqlite3 dialect not found")
	}
	schema := Parse(&User{}, dialect)
	if schema.Name != "User" {
		t.Fatalf("expect schema name User, but got %s", schema.Name)
	}
	if schema.GetField("ID").Type != "integer" {
		t.Fatalf("expect ID type integer, but got %s", schema.GetField("ID").Type)
	}
	if schema.GetField("ID").Tag != "PRIMARY KEY" {
		t.Fatalf("expect ID tag PRIMARY KEY, but got %s", schema.GetField("ID").Tag)
	}
	if schema.GetField("Name").Type != "text" {
		t.Fatalf("expect Name type text, but got %s", schema.GetField("Name").Type)
	}
}
