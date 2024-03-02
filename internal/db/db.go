package db

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sqlx.DB
}

func New(dbName string) *DB {
	rootDir, _ := os.Getwd()
	db, err := sqlx.Connect("sqlite3", filepath.Join(rootDir, dbName))
	if err != nil {
		log.Fatal(err)
	}
	return &DB{db}
}

func (db *DB) RunMigrations() error {
	sqlCreate := `
	CREATE TABLE IF NOT EXISTS blogger (
		id INTEGER PRIMARY KEY, 
		first_name TEXT, 
		last_name TEXT,
		email TEXT,
		age INTEGER,
		gender TEXT,
		bio TEXT,
		create_date DATETIME
	);

	CREATE TABLE IF NOT EXISTS post (
		id INTEGER PRIMARY KEY, 
		author_id INTEGER,
		title TEXT,
		content TEXT,
		likes INTEGER,
		dislikes INTEGER,
		create_date DATETIME
	);

	CREATE TABLE IF NOT EXISTS user (
		id INTEGER PRIMARY KEY, 
		username TEXT,
		password TEXT,
		create_date DATETIME
	);

	CREATE TABLE IF NOT EXISTS comment (
		id INTEGER PRIMARY KEY, 
		user_id INTEGER,
		post_id INTEGER,
		content TEXT,
		likes INTEGER,
		dislikes INTEGER,
		create_date DATETIME
	);
	`

	_, err := db.Exec(sqlCreate)
	return err
}

func Teardown(dbName string) {
	rootDir, _ := os.Getwd()
	err := os.Remove(filepath.Join(rootDir, dbName))
	if err != nil {
		log.Fatal(err)
	}
}
