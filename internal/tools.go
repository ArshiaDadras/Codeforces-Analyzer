package internal

import (
	"encoding/json"
	"regexp"
	"strconv"

	"github.com/ArshiaDadras/Codeforces-Analyzer/internal/codeforces"
)

func UpdateProblemsFromAPI() {
	InitDB()

	problems, problemStatistics, err := codeforces.GetProblems([]string{}, "")
	if err != nil {
		panic(err)
	}

	for i, problem := range problems {
		marshaledTags, err := json.Marshal(problem.Tags)
		if err != nil {
			panic(err)
		}

		var problemID int
		problem.SolvedCount = problemStatistics[i].SolvedCount
		if err := db.QueryRow("SELECT id FROM problems WHERE contest_id = ? AND idx = ?", problem.ContestID, problem.Index).Scan(&problemID); err != nil {
			if _, err = db.Exec("INSERT INTO problems (contest_id, problemset_name, idx, name, type, points, rating, tags, solved_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", problem.ContestID, problem.ProblemsetName, problem.Index, problem.Name, problem.Type, problem.Points, problem.Rating, marshaledTags, problem.SolvedCount); err != nil {
				panic(err)
			}
		} else {
			if _, err = db.Exec("UPDATE problems SET problemset_name = ?, name = ?, type = ?, points = ?, rating = ?, tags = ?, solved_count = ? WHERE contest_id = ? AND idx = ?", problem.ProblemsetName, problem.Name, problem.Type, problem.Points, problem.Rating, marshaledTags, problem.SolvedCount, problem.ContestID, problem.Index); err != nil {
				panic(err)
			}
		}
	}
}

func CrawlBlogEntry(blogID int) {
	InitDB()

	blog, err := codeforces.GetBlogEntry(blogID)
	if err != nil {
		return
	}
	marshaledTags, err := json.Marshal(blog.Tags)
	if err != nil {
		panic(err)
	}

	var lastModificationTime int
	if err = db.QueryRow("SELECT modification_time FROM blog_entries WHERE id = ?", blog.ID).Scan(&lastModificationTime); err != nil {
		lastModificationTime = 0
	}

	if lastModificationTime > 0 {
		if lastModificationTime >= blog.ModificationTimeSeconds {
			return
		}

		if _, err = db.Exec("UPDATE blog_entries SET original_locale = ?, creation_time = ?, author_handle = ?, title = ?, content = ?, locale = ?, modification_time = ?, allow_view_history = ?, tags = ?, rating = ? WHERE id = ?", blog.OriginalLocale, blog.CreationTimeSeconds, blog.AuthorHandle, blog.Title, blog.Content, blog.Locale, blog.ModificationTimeSeconds, blog.AllowViewHistory, marshaledTags, blog.Rating, blog.ID); err != nil {
			panic(err)
		}
	} else {
		if _, err = db.Exec("INSERT INTO blog_entries (id, original_locale, creation_time, author_handle, title, content, locale, modification_time, allow_view_history, tags, rating) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", blog.ID, blog.OriginalLocale, blog.CreationTimeSeconds, blog.AuthorHandle, blog.Title, blog.Content, blog.Locale, blog.ModificationTimeSeconds, blog.AllowViewHistory, marshaledTags, blog.Rating); err != nil {
			panic(err)
		}
	}

	// TODO: Analyze current blog problems

	r := regexp.MustCompile(`https://codeforces.com/blog/entry/(\d+)`)
	matches := r.FindAllStringSubmatch(blog.Content, -1)
	for _, match := range matches {
		nextBlogID, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		CrawlBlogEntry(nextBlogID)
	}
}
