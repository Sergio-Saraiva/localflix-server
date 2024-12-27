package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type AppDatabase struct {
	Db *sql.DB
}

func NewAppDatabase() *AppDatabase {
	db, err := sql.Open("sqlite3", ".file:locked.db?cache=shared")
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		panic(err)
	}

	db.Exec("CREATE TABLE IF NOT EXISTS categories (id INTEGER PRIMARY KEY, name TEXT)")
	db.Exec("CREATE TABLE IF NOT EXISTS folders (id INTEGER PRIMARY KEY, path TEXT, category_id INTEGER, FOREIGN KEY(category_id) REFERENCES categories(id))")
	return &AppDatabase{
		Db: db,
	}
}
