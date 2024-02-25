package codeforces

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

func SortedParams(url string) string {
	params := strings.Split(strings.Split(url, "?")[1], "&")
	sort.Strings(params)

	return strings.Split(url, "?")[0] + "?" + strings.Join(params, "&")
}

func GenerateAPISig(url, secret string) string {
	rnd := rand.Intn(1000000-100000) + 100000
	url = SortedParams(url[strings.Index(url, "/api/")+5:])

	SHA512 := sha512.New()
	SHA512.Write([]byte(fmt.Sprintf("%d/%s#%s", rnd, url, secret)))

	return fmt.Sprintf("%d%x", rnd, SHA512.Sum(nil))
}

func GetRequest(url string) ([]byte, error) {
	public := os.Getenv("CF_PUBLIC_KEY")
	secret := os.Getenv("CF_SECRET_KEY")
	if public != "" && secret != "" {
		if url[len(url)-1] != '?' {
			url += "&"
		}
		url += fmt.Sprintf("apiKey=%s&time=%d", public, time.Now().Unix())
		url += fmt.Sprintf("&apiSig=%s", GenerateAPISig(url, secret))
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := struct {
		Status  string          `json:"status"`
		Result  json.RawMessage `json:"result"`
		Comment string          `json:"comment"`
	}{}
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.Status != "OK" {
		return nil, fmt.Errorf(`codeforces API returned status %s with error message "%s"`, response.Status, response.Comment)
	}

	return response.Result, nil
}

func (blogEntry *BlogEntry) GetComments() ([]Comment, error) {
	resp, err := GetRequest(fmt.Sprintf("https://codeforces.com/api/blogEntry.comments?blogEntryId=%d", blogEntry.ID))
	if err != nil {
		return nil, err
	}

	comments := []Comment{}
	if err = json.Unmarshal(resp, &comments); err != nil {
		return nil, err
	}

	return comments, nil
}

func (blogEntry *BlogEntry) GetView() (BlogEntry, error) {
	resp, err := GetRequest(fmt.Sprintf("https://codeforces.com/api/blogEntry.view?blogEntryId=%d", blogEntry.ID))
	if err != nil {
		return BlogEntry{}, err
	}

	blogEntryResult := BlogEntry{}
	if err = json.Unmarshal(resp, &blogEntryResult); err != nil {
		return BlogEntry{}, err
	}

	return blogEntryResult, nil
}

func (contest *Contest) GetHacks() ([]Hack, error) {
	resp, err := GetRequest(fmt.Sprintf("https://codeforces.com/api/contest.hacks?contestId=%d", contest.ID))
	if err != nil {
		return nil, err
	}

	hacks := []Hack{}
	if err = json.Unmarshal(resp, &hacks); err != nil {
		return nil, err
	}

	return hacks, nil
}

func GetContestList(gym bool) ([]Contest, error) {
	resp, err := GetRequest(fmt.Sprintf("https://codeforces.com/api/contest.list?gym=%t", gym))
	if err != nil {
		return nil, err
	}

	contests := []Contest{}
	if err = json.Unmarshal(resp, &contests); err != nil {
		return nil, err
	}

	return contests, nil
}

func (contest *Contest) GetRatingChanges() ([]RatingChange, error) {
	resp, err := GetRequest(fmt.Sprintf("https://codeforces.com/api/contest.ratingChanges?contestId=%d", contest.ID))
	if err != nil {
		return nil, err
	}

	ratingChanges := []RatingChange{}
	if err = json.Unmarshal(resp, &ratingChanges); err != nil {
		return nil, err
	}

	return ratingChanges, nil
}

func (contest *Contest) GetStandings(from, count int, handles []string, room int, showUnofficial bool) (Standings, error) {
	url := fmt.Sprintf("https://codeforces.com/api/contest.standings?contestId=%d&showUnofficial=%t", contest.ID, showUnofficial)
	if from > 1 {
		url += fmt.Sprintf("&from=%d", from)
	}
	if count > 0 {
		url += fmt.Sprintf("&count=%d", count)
	}
	if len(handles) > 0 {
		url += fmt.Sprintf("&handles=%s", strings.Join(handles, ";"))
	}
	if room > 0 {
		url += fmt.Sprintf("&room=%d", room)
	}

	resp, err := GetRequest(url)
	if err != nil {
		return Standings{}, err
	}

	standings := Standings{}
	if err = json.Unmarshal(resp, &standings); err != nil {
		return Standings{}, err
	}

	return standings, nil
}

func (contest *Contest) GetStatus(from, count int, handle string) ([]Submission, error) {
	url := fmt.Sprintf("https://codeforces.com/api/contest.status?contestId=%d", contest.ID)
	if from > 1 {
		url += fmt.Sprintf("&from=%d", from)
	}
	if count > 0 {
		url += fmt.Sprintf("&count=%d", count)
	}
	if handle != "" {
		url += fmt.Sprintf("&handle=%s", handle)
	}

	resp, err := GetRequest(url)
	if err != nil {
		return nil, err
	}

	submissions := []Submission{}
	if err = json.Unmarshal(resp, &submissions); err != nil {
		return nil, err
	}

	return submissions, nil
}

func GetProblems(tags []string, problemsetName string) ([]Problem, []ProblemStatistics, error) {
	url := "https://codeforces.com/api/problemset.problems?"
	if len(tags) > 0 {
		url += fmt.Sprintf("tags=%s", strings.Join(tags, ";"))
	}
	if problemsetName != "" {
		if url[len(url)-1] != '?' {
			url += "&"
		}
		url += fmt.Sprintf("problemsetName=%s", problemsetName)
	}

	resp, err := GetRequest(url)
	if err != nil {
		return nil, nil, err
	}

	response := struct {
		Problems          []Problem           `json:"problems"`
		ProblemStatistics []ProblemStatistics `json:"problemStatistics"`
	}{}
	if err = json.Unmarshal(resp, &response); err != nil {
		return nil, nil, err
	}

	return response.Problems, response.ProblemStatistics, nil
}

func GetRecentStatus(maxCount int, problemsetName string) ([]Submission, error) {
	url := "https://codeforces.com/api/recentActions?"
	if maxCount > 0 {
		url += fmt.Sprintf("maxCount=%d", maxCount)
	}
	if problemsetName != "" {
		if url[len(url)-1] != '?' {
			url += "&"
		}
		url += fmt.Sprintf("problemsetName=%s", problemsetName)
	}

	resp, err := GetRequest(url)
	if err != nil {
		return nil, err
	}

	submissions := []Submission{}
	if err = json.Unmarshal(resp, &submissions); err != nil {
		return nil, err
	}

	return submissions, nil
}

func GetRecentActions(maxCount int) ([]RecentAction, error) {
	url := "https://codeforces.com/api/recentActions?"
	if maxCount > 0 {
		url += fmt.Sprintf("maxCount=%d", maxCount)
	}

	resp, err := GetRequest(url)
	if err != nil {
		return nil, err
	}

	actions := []RecentAction{}
	if err = json.Unmarshal(resp, &actions); err != nil {
		return nil, err
	}

	return actions, nil
}

func (user *User) GetBlogEntries() ([]BlogEntry, error) {
	resp, err := GetRequest(fmt.Sprintf("https://codeforces.com/api/user.blogEntries?handle=%s", user.Handle))
	if err != nil {
		return nil, err
	}

	blogEntries := []BlogEntry{}
	if err = json.Unmarshal(resp, &blogEntries); err != nil {
		return nil, err
	}

	return blogEntries, nil
}

func (user *User) GetFriends(onlyOnline bool) ([]string, error) {
	resp, err := GetRequest(fmt.Sprintf("https://codeforces.com/api/user.friends?onlyOnline=%t&handle=%s", onlyOnline, user.Handle))
	if err != nil {
		return nil, err
	}

	friends := []string{}
	if err = json.Unmarshal(resp, &friends); err != nil {
		return nil, err
	}

	return friends, nil
}

func (user *User) GetInfo() (User, error) {
	resp, err := GetRequest(fmt.Sprintf("https://codeforces.com/api/user.info?handles=%s", user.Handle))
	if err != nil {
		return User{}, err
	}

	users := []User{}
	if err = json.Unmarshal(resp, &users); err != nil {
		return User{}, err
	}

	return users[0], nil
}

func GetUsersInfo(handles []string) ([]User, error) {
	resp, err := GetRequest(fmt.Sprintf("https://codeforces.com/api/user.info?handles=%s", strings.Join(handles, ";")))
	if err != nil {
		return nil, err
	}

	users := []User{}
	if err = json.Unmarshal(resp, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (contest *Contest) GetRatedList(activeOnly, includeRetired bool) ([]User, error) {
	url := fmt.Sprintf("https://codeforces.com/api/user.ratedList?activeOnly=%t&includeRetired=%t", activeOnly, includeRetired)
	if contest.ID > 0 {
		url += fmt.Sprintf("&contestId=%d", contest.ID)
	}

	resp, err := GetRequest(url)
	if err != nil {
		return nil, err
	}

	users := []User{}
	if err = json.Unmarshal(resp, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func GetGlobalRatedList(activeOnly, includeRetired bool) ([]User, error) {
	url := fmt.Sprintf("https://codeforces.com/api/user.ratedList?activeOnly=%t&includeRetired=%t", activeOnly, includeRetired)

	resp, err := GetRequest(url)
	if err != nil {
		return nil, err
	}

	users := []User{}
	if err = json.Unmarshal(resp, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (user *User) GetRating() ([]RatingChange, error) {
	resp, err := GetRequest(fmt.Sprintf("https://codeforces.com/api/user.rating?handle=%s", user.Handle))
	if err != nil {
		return nil, err
	}

	ratingChanges := []RatingChange{}
	if err = json.Unmarshal(resp, &ratingChanges); err != nil {
		return nil, err
	}

	return ratingChanges, nil
}
