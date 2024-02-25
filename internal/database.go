package internal

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB = nil

func InitDB() {
	if db != nil {
		return
	}

	var err error
	db, err = sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		panic(err)
	}
}

func CreateTables() {
	InitDB()

	executionCommands := []string{
		`CREATE TABLE IF NOT EXISTS blog_entries (
			id INTEGER PRIMARY KEY,
			original_locale TEXT,
			creation_time INTEGER,
			author_handle TEXT,
			title TEXT,
			content TEXT,
			locale TEXT,
			modification_time INTEGER,
			allow_view_history BOOLEAN,
			tags TEXT,
			rating INTEGER NULL
		)`,
		`CREATE TABLE IF NOT EXISTS problems (
			id INTEGER PRIMARY KEY,
			contest_id INTEGER NULL,
			problemset_name TEXT NULL,
			idx TEXT,
			name TEXT,
			type TEXT,
			points REAL,
			rating INTEGER NULL,
			tags TEXT,
			solved_count INTEGER NULL
		)`,
		"CREATE INDEX IF NOT EXISTS idx_blog_entries_title ON blog_entries (title)",
		"CREATE INDEX IF NOT EXISTS idx_blog_entries_tags ON blog_entries (tags)",
		"CREATE INDEX IF NOT EXISTS idx_blog_entries_rating ON blog_entries (rating)",
		"CREATE INDEX IF NOT EXISTS idx_problems_contest_id ON problems (contest_id)",
		"CREATE INDEX IF NOT EXISTS idx_problems_idx ON problems (idx)",
		"CREATE INDEX IF NOT EXISTS idx_problems_rating ON problems (rating)",
		"CREATE INDEX IF NOT EXISTS idx_problems_tags ON problems (tags)",
		"PRAGMA foreign_keys = ON",
		"VACUUM",
		"ANALYZE",
		"REINDEX",
	}

	for _, command := range executionCommands {
		fmt.Println(command)
		_, err := db.Exec(command)
		if err != nil {
			panic(err)
		}
	}
}
