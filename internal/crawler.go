package internal

import (
	"database/sql"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/ArshiaDadras/Codeforces-Analyzer/internal/codeforces"
)

const CodeforcesUrl = `((http|https)://)?(www.)?codeforces\.com`
const problemUrlRegex = CodeforcesUrl + `/(([A-Za-z/]+/problem/\d+/[A-Za-z\d]+)|(contest/\d+/problem/[A-Za-z\d]+)|(gym/\d+/problem/[A-Za-z\d]+))`
const blogUrlRegex = CodeforcesUrl + `/blog/entry/(\d+)`

func UpdateProblemsFromAPI() error {
	log.Println("Updating problems from API...")

	problems, problemStatistics, err := codeforces.GetProblems([]string{}, "")
	if err != nil {
		return err
	}

	for i, problem := range problems {
		problem.SolvedCount = problemStatistics[i].SolvedCount
		if err := SaveProblem(problem); err != nil {
			return err
		}
	}

	return nil
}

func FindTagsForProblem(problemUrl string, content string) []string {
	return []string{}
}

func AnalyzeProblem(problemUrl string, blogID int, content string) error {
	log.Printf("Analyzing problem %s...\n", problemUrl)

	data := strings.Split(problemUrl, "/")
	data = data[len(data)-4:]
	if data[0] == "gym" || data[0] == "contest" {
		data[1], data[2] = data[2], data[1]
	}
	data = append(data[:1], data[2:]...)

	problemID, err := strconv.Atoi(data[1])
	if err != nil {
		return err
	}

	referenced := &codeforces.ReferencedProblem{
		BlogID:      blogID,
		ProblemType: data[0],
		ProblemID:   problemID,
		Index:       data[2],
		Tags:        FindTagsForProblem(problemUrl, content),
	}

	return SaveReferencedProblem(referenced)
}

func AnalyzeProblemsOnBlog(blog *codeforces.BlogEntry) []int {
	r := regexp.MustCompile(problemUrlRegex)
	for _, match := range r.FindAllStringSubmatch(blog.Content, -1) {
		if err := AnalyzeProblem(match[0], blog.ID, blog.Content); err != nil {
			continue
		}
	}

	blogIDs := make([]int, 0)
	r = regexp.MustCompile(blogUrlRegex)
	for _, match := range r.FindAllStringSubmatch(blog.Content, -1) {
		blogID, err := strconv.Atoi(match[len(match)-1])
		if err != nil {
			continue
		}
		blogIDs = append(blogIDs, blogID)
	}

	return blogIDs
}

func AnalyzeProblemsOnComments(blog *codeforces.BlogEntry) []int {
	r := regexp.MustCompile(problemUrlRegex)
	for _, comment := range blog.Comments {
		for _, match := range r.FindAllStringSubmatch(comment.Text, -1) {
			if err := AnalyzeProblem(match[0], blog.ID, blog.Content+`<div class="comment">`+comment.Text+`</div>`); err != nil {
				continue
			}
		}
	}

	blogIDs := make([]int, 0)
	r = regexp.MustCompile(blogUrlRegex)
	for _, comment := range blog.Comments {
		for _, match := range r.FindAllStringSubmatch(comment.Text, -1) {
			blogID, err := strconv.Atoi(match[len(match)-1])
			if err != nil {
				continue
			}
			blogIDs = append(blogIDs, blogID)
		}
	}

	return blogIDs
}

func CrawlBlogEntry(blogID int) error {
	log.Printf("Crawling blog %d...\n", blogID)

	blog, err := codeforces.GetBlogEntry(blogID)
	if err != nil {
		return err
	}
	if strings.Contains(strings.ToLower(blog.Title), "editorial") {
		log.Printf("Skipping blog %d because it's an editorial...\n", blogID)
		return nil
	}
	lastVersion, err := GetBlogEntry(blogID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	nextBlogs := make([]int, 0)
	if lastVersion == nil || lastVersion.ModificationTimeSeconds < blog.ModificationTimeSeconds {
		nextBlogs = AnalyzeProblemsOnBlog(blog)
	}
	if lastVersion == nil || len(lastVersion.Comments) < len(blog.Comments) {
		nextBlogs = append(nextBlogs, AnalyzeProblemsOnComments(blog)...)
	}

	if err := SaveBlogEntry(blog); err != nil {
		return err
	}

	for _, nextBlogID := range nextBlogs {
		if nextBlogID == blogID {
			continue
		}

		err := CrawlBlogEntry(nextBlogID)
		if err != nil {
			log.Printf("Error crawling blog %d: %s\n", nextBlogID, err)
		}
	}

	return nil
}
