package main

import (
	"fmt"
	geeorm "geemod"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	engine, err := geeorm.NewEngine("sqlite3", "gee.db")
	if err != nil {
		panic(err)
	}
	defer engine.Close()
	s := engine.Session()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}
