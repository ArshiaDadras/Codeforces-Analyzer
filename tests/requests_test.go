package tests

import (
	"os"
	"strings"
	"testing"

	codeforces "github.com/ArshiaDadras/Codeforces-Analyzer/internal/codeforces"
	"github.com/joho/godotenv"
)

var envLoadError error = nil

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		envLoadError = err
	}

	os.Exit(m.Run())
}

func TestGetComments(t *testing.T) {
	blogEntry := codeforces.BlogEntry{ID: 62865}
	comments, err := blogEntry.GetComments()
	if err != nil {
		t.Fatal(err)
	}

	if len(comments) < 123 {
		t.Error("Invalid number of comments")
	}
}

func TestGetView(t *testing.T) {
	blogEntry, err := codeforces.GetBlogEntry(62865)
	if err != nil {
		t.Fatal(err)
	}

	if blogEntry.ID != 62865 {
		t.Error("BlogEntry ID mismatch")
	}
}

func TestGetHacks(t *testing.T) {
	contest := codeforces.Contest{ID: 1923}
	hacks, err := contest.GetHacks()
	if err != nil {
		t.Fatal(err)
	}

	if len(hacks) != 2191 {
		t.Error("Invalid number of hacks")
	}
}

func TestGetContestList(t *testing.T) {
	contests, err := codeforces.GetContestList(false)
	if err != nil {
		t.Fatal(err)
	}

	if len(contests) == 0 {
		t.Error("No contests found")
	}
}

func TestGetRatingChanges(t *testing.T) {
	contest := codeforces.Contest{ID: 1923}
	ratingChanges, err := contest.GetRatingChanges()
	if err != nil {
		t.Fatal(err)
	}

	if len(ratingChanges) != 16715 {
		t.Error("Invalid number of rating changes")
	}
}

func TestGetStandings(t *testing.T) {
	contest := codeforces.Contest{ID: 1923}
	standings, err := contest.GetStandings(10, 10, []string{}, 0, false)
	if err != nil {
		t.Fatal(err)
	}

	if standings.Contest.ID != contest.ID {
		t.Error("Contest ID mismatch")
	}
	if len(standings.Rows) != 10 {
		t.Error("Invalid number of standings")
	}
	for index, row := range standings.Rows {
		if row.Rank > index+10 {
			t.Error("Invalid standings")
		}
	}
}

func TestGetStatus(t *testing.T) {
	contest := codeforces.Contest{ID: 1923}
	status, err := contest.GetStatus(10, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if len(status) != 10 {
		t.Error("Invalid number of submissions")
	}
}

func TestGetProblems(t *testing.T) {
	problems, problemStatistics, err := codeforces.GetProblems([]string{}, "")
	if err != nil {
		t.Fatal(err)
	}

	if len(problems) < 9297 {
		t.Error("Invalid number of problems")
	}
	if len(problemStatistics) < 9297 {
		t.Error("Invalid number of problem statistics")
	}
}

func TestGetRecentStatus(t *testing.T) {
	status, err := codeforces.GetRecentStatus(10, "")
	if err != nil {
		t.Fatal(err)
	}

	if len(status) != 10 {
		t.Error("Invalid number of submissions")
	}
}

func TestGetRecentActions(t *testing.T) {
	actions, err := codeforces.GetRecentActions(10)
	if err != nil {
		t.Fatal(err)
	}

	if len(actions) != 10 {
		t.Error("Invalid number of recent actions")
	}
}

func TestGetBlogEntries(t *testing.T) {
	user := codeforces.User{Handle: "MikeMirzayanov"}
	blogEntries, err := user.GetBlogEntries()
	if err != nil {
		t.Fatal(err)
	}

	if len(blogEntries) < 418 {
		t.Error("Invalid number of blog entries")
	}
}

func TestGetFriends(t *testing.T) {
	if envLoadError != nil {
		t.Skipf(`Error loading .env file: "%v"`, envLoadError)
		return
	}
	if os.Getenv("CF_HANDLE") == "" {
		t.Skip("CF_HANDLE not found in .env file")
		return
	}

	user := codeforces.User{Handle: os.Getenv("CF_HANDLE")}
	friends, err := user.GetFriends(false)
	if err != nil {
		if strings.Contains(err.Error(), "You have to be authenticated to use this method") {
			t.Skip("Authentication required")
			return
		}
		t.Fatal(err)
	}

	if len(friends) < 1 {
		t.Error("Invalid number of friends")
	}
}

func TestGetInfo(t *testing.T) {
	user := codeforces.User{Handle: "MikeMirzayanov"}
	info, err := user.GetInfo()
	if err != nil {
		t.Fatal(err)
	}

	if info.Handle != user.Handle {
		t.Error("User handle mismatch")
	}
}

func TestGetRatedList(t *testing.T) {
	contest := codeforces.Contest{ID: 1923}
	users, err := contest.GetRatedList(false, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(users) < 1 {
		t.Error("Invalid number of rated users")
	}
}

func TestGetRating(t *testing.T) {
	user := codeforces.User{Handle: "ArshiaDadras"}
	ratingChanges, err := user.GetRating()
	if err != nil {
		t.Fatal(err)
	}

	if len(ratingChanges) < 75 {
		t.Error("Invalid number of rating changes")
	}
}
