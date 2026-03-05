package clause

import "testing"

func TestUpdate(t *testing.T) {
	sql, vars := _update("user", map[string]any{"name": "geektutu", "id": 1})
	t.Log(sql)
	if len(vars) != 2 {
		t.Fatal("vars err")
	}
	if vars[0] != "geektutu" || vars[1] != 1 {
		t.Fatal("vars err")
	}
}
