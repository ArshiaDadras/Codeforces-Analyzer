package internal

import (
	"database/sql"
	"encoding/json"

	"github.com/ArshiaDadras/Codeforces-Analyzer/internal/codeforces"
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
			tags JSON,
			rating INTEGER NULL
		)`,
		`CREATE TABLE IF NOT EXISTS problems (
			contest_id INTEGER NULL,
			problemset_name TEXT NULL,
			idx TEXT,
			name TEXT,
			type TEXT,
			points REAL,
			rating INTEGER NULL,
			tags JSON,
			solved_count INTEGER NULL,
			PRIMARY KEY (contest_id, idx)
		)`,
		`CREATE TABLE IF NOT EXISTS referenced_problems (
			id INTEGER PRIMARY KEY,
			blog_id INTEGER,
			problem_type TEXT,
			problem_id INTEGER,
			idx TEXT,
			tags JSON
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
		_, err := db.Exec(command)
		if err != nil {
			panic(err)
		}
	}
}

func SaveBlogEntry(blog *codeforces.BlogEntry) error {
	marshaledTags, err := json.Marshal(blog.Tags)
	if err != nil {
		return err
	}

	if _, err = db.Exec("INSERT INTO blog_entries (id, original_locale, creation_time, author_handle, title, content, locale, modification_time, allow_view_history, tags, rating) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", blog.ID, blog.OriginalLocale, blog.CreationTimeSeconds, blog.AuthorHandle, blog.Title, blog.Content, blog.Locale, blog.ModificationTimeSeconds, blog.AllowViewHistory, marshaledTags, blog.Rating); err != nil {
		_, err = db.Exec("UPDATE blog_entries SET original_locale = ?, creation_time = ?, author_handle = ?, title = ?, content = ?, locale = ?, modification_time = ?, allow_view_history = ?, tags = ?, rating = ? WHERE id = ?", blog.OriginalLocale, blog.CreationTimeSeconds, blog.AuthorHandle, blog.Title, blog.Content, blog.Locale, blog.ModificationTimeSeconds, blog.AllowViewHistory, marshaledTags, blog.Rating, blog.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func SaveProblem(problem *codeforces.Problem) error {
	marshaledTags, err := json.Marshal(problem.Tags)
	if err != nil {
		return err
	}

	if _, err = db.Exec("INSERT INTO problems (contest_id, problemset_name, idx, name, type, points, rating, tags, solved_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", problem.ContestID, problem.ProblemsetName, problem.Index, problem.Name, problem.Type, problem.Points, problem.Rating, marshaledTags, problem.SolvedCount); err != nil {
		if _, err = db.Exec("UPDATE problems SET problemset_name = ?, name = ?, type = ?, points = ?, rating = ?, tags = ?, solved_count = ? WHERE contest_id = ? AND idx = ?", problem.ProblemsetName, problem.Name, problem.Type, problem.Points, problem.Rating, marshaledTags, problem.SolvedCount, problem.ContestID, problem.Index); err != nil {
			return err
		}
	}

	return nil
}

func mergeTags(currentTags []string, newTags []string) []string {
	for _, tag := range newTags {
		found := false
		for _, currentTag := range currentTags {
			if tag == currentTag {
				found = true
				break
			}
		}

		if !found {
			currentTags = append(currentTags, tag)
		}
	}

	return currentTags
}

func SaveReferencedProblem(referenced *codeforces.ReferencedProblem) error {
	var referencedID int
	var currentMarshaledTags []byte
	if err := db.QueryRow("SELECT id, tags FROM referenced_problems WHERE blog_id = ? AND problem_type = ? AND problem_id = ? AND idx = ?", referenced.BlogID, referenced.ProblemType, referenced.ProblemID, referenced.Index).Scan(&referencedID, &currentMarshaledTags); err == sql.ErrNoRows {
		marshaledTags, err := json.Marshal(referenced.Tags)
		if err != nil {
			return err
		}

		if _, err = db.Exec("INSERT INTO referenced_problems (blog_id, problem_type, problem_id, idx, tags) VALUES (?, ?, ?, ?, ?)", referenced.BlogID, referenced.ProblemType, referenced.ProblemID, referenced.Index, marshaledTags); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		var currentTags []string
		if err := json.Unmarshal(currentMarshaledTags, &currentTags); err != nil {
			return err
		}

		newTags := mergeTags(currentTags, referenced.Tags)

		marshaledTags, err := json.Marshal(newTags)
		if err != nil {
			return err
		}

		if _, err = db.Exec("UPDATE referenced_problems SET tags = ? WHERE id = ?", marshaledTags, referencedID); err != nil {
			return err
		}
	}

	return nil
}

func GetBlogEntry(blogID int) (*codeforces.BlogEntry, error) {
	var marshaledTags []byte
	blog := new(codeforces.BlogEntry)
	if err := db.QueryRow("SELECT original_locale, creation_time, author_handle, title, content, locale, modification_time, allow_view_history, tags, rating FROM blog_entries WHERE id = ?", blogID).Scan(&blog.OriginalLocale, &blog.CreationTimeSeconds, &blog.AuthorHandle, &blog.Title, &blog.Content, &blog.Locale, &blog.ModificationTimeSeconds, &blog.AllowViewHistory, &marshaledTags, &blog.Rating); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(marshaledTags, &blog.Tags); err != nil {
		return nil, err
	}

	return blog, nil
}
