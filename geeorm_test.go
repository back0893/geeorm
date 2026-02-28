package geeorm

import (
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
