package dialect

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestMain(t *testing.M) {
	RegisterDialect("sqlite3", &sqlite3Dialect{})
	r := t.Run()
	os.Exit(r)
}
func TestSqlite3(t *testing.T) {
	dialect, ok := GetDialect("sqlite3")
	if !ok {
		t.Fatal("sqlite3 dialect not found")
	}
	if v := dialect.DataTypeOf(reflect.ValueOf(int32(1))); v != "integer" {
		t.Fatalf("expect INTEGER, but got %s", v)
	}
	if v := dialect.DataTypeOf(reflect.ValueOf(time.Now())); v != "datetime" {
		t.Fatalf("expect INTEGER, but got %s", v)
	}
	if v := dialect.DataTypeOf(reflect.ValueOf([]int{1})); v != "blob" {
		t.Fatalf("expect blob, but got %s", v)
	}
}
