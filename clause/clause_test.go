package clause

import "testing"

func testSelect(t *testing.T) {
	var clause Clause
	clause.Set(LIMIT, 3)
	clause.Set(SELECT, "User", []string{"*"})
	clause.Set(WHERE, "Name = ?", "Tom")
	clause.Set(ORDERBY, "Age ASC")
	sql, vars := clause.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Log(sql)
	if sql != "select * from User where Name = ? order by Age ASC limit ?" {
		t.Fatalf("expect select * from User where Name = ? order by Age ASC limit ?, but got %s", sql)
	}
	if len(vars) != 2 || vars[0] != "Tom" || vars[1] != 3 {
		t.Fatalf("expect vars to be [Tom 3], but got %v", vars)
	}
}

func TestClause_Build(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		testSelect(t)
	})
}
