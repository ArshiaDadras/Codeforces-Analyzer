package codeforces

type User struct {
	Handle                  string `json:"handle"`
	Email                   string `json:"email"`
	VkId                    string `json:"vkId"`
	OpenID                  string `json:"openId"`
	FirstName               string `json:"firstName"`
	LastName                string `json:"lastName"`
	Country                 string `json:"country"`
	City                    string `json:"city"`
	Organization            string `json:"organization"`
	Contributions           int    `json:"contribution"`
	Rank                    string `json:"rank"`
	Rating                  int    `json:"rating"`
	MaxRank                 string `json:"maxRank"`
	MaxRating               int    `json:"maxRating"`
	LastOnlineTimeSeconds   int    `json:"lastOnlineTimeSeconds"`
	RegistrationTimeSeconds int    `json:"registrationTimeSeconds"`
	FriendOfCount           int    `json:"friendOfCount"`
	Avatar                  string `json:"avatar"`
	TitlePhoto              string `json:"titlePhoto"`
}

type BlogEntry struct {
	ID                      int      `json:"id"`
	OriginalLocale          string   `json:"originalLocale"`
	CreationTimeSeconds     int      `json:"creationTimeSeconds"`
	AuthorHandle            string   `json:"authorHandle"`
	Title                   string   `json:"title"`
	Content                 string   `json:"content"`
	Locale                  string   `json:"locale"`
	ModificationTimeSeconds int      `json:"modificationTimeSeconds"`
	AllowViewHistory        bool     `json:"allowViewHistory"`
	Tags                    []string `json:"tags"`
	Rating                  int      `json:"rating"`
}

type Comment struct {
	ID                  int    `json:"id"`
	CreationTimeSeconds int    `json:"creationTimeSeconds"`
	CommentatorHandle   string `json:"commentatorHandle"`
	Locale              string `json:"locale"`
	Text                string `json:"text"`
	ParentCommentId     int    `json:"parentCommentId"`
	Rating              int    `json:"rating"`
}

type RecentAction struct {
	TimeSeconds int       `json:"timeSeconds"`
	BlogEntry   BlogEntry `json:"blogEntry"`
	Comment     Comment   `json:"comment"`
}

type RatingChange struct {
	ContestID               int    `json:"contestId"`
	ContestName             string `json:"contestName"`
	Handle                  string `json:"handle"`
	Rank                    int    `json:"rank"`
	RatingUpdateTimeSeconds int    `json:"ratingUpdateTimeSeconds"`
	OldRating               int    `json:"oldRating"`
	NewRating               int    `json:"newRating"`
}

type Contest struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Type                string `json:"type"`
	Phase               string `json:"phase"`
	Frozen              bool   `json:"frozen"`
	DurationSeconds     int    `json:"durationSeconds"`
	StartTimeSeconds    int    `json:"startTimeSeconds"`
	RelativeTimeSeconds int    `json:"relativeTimeSeconds"`
	PreparedBy          string `json:"preparedBy"`
	WebsiteURL          string `json:"websiteUrl"`
	Description         string `json:"description"`
	Difficulty          int    `json:"difficulty"`
	Kind                string `json:"kind"`
	IcpcRegion          string `json:"icpcRegion"`
	Country             string `json:"country"`
	City                string `json:"city"`
	Season              string `json:"season"`
}

type Party struct {
	ContestID        int    `json:"contestId"`
	Members          []User `json:"members"`
	ParticipantType  string `json:"participantType"`
	TeamID           int    `json:"teamId"`
	TeamName         string `json:"teamName"`
	Ghost            bool   `json:"ghost"`
	Room             int    `json:"room"`
	StartTimeSeconds int    `json:"startTimeSeconds"`
}

type Member struct {
	Handle string `json:"handle"`
	Name   string `json:"name"`
}

type Problem struct {
	ContestID      int      `json:"contestId"`
	ProblemsetName string   `json:"problemsetName"`
	Index          string   `json:"index"`
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	Points         float64  `json:"points"`
	Rating         int      `json:"rating"`
	Tags           []string `json:"tags"`
	SolvedCount    int      `json:"solvedCount"`
}

type ProblemStatistics struct {
	ContestID   int    `json:"contestId"`
	Index       string `json:"index"`
	SolvedCount int    `json:"solvedCount"`
}

type ProblemResult struct {
	Points                    float64 `json:"points"`
	Penalty                   int     `json:"penalty"`
	RejectedAttemptCount      int     `json:"rejectedAttemptCount"`
	Type                      string  `json:"type"`
	BestSubmissionTimeSeconds int     `json:"bestSubmissionTimeSeconds"`
}

type Submission struct {
	ID                  int     `json:"id"`
	ContestID           int     `json:"contestId"`
	CreationTimeSeconds int     `json:"creationTimeSeconds"`
	RelativeTimeSeconds int     `json:"relativeTimeSeconds"`
	Problem             Problem `json:"problem"`
	Author              User    `json:"author"`
	ProgrammingLanguage string  `json:"programmingLanguage"`
	Verdict             string  `json:"verdict"`
	Testset             string  `json:"testset"`
	PassedTestCount     int     `json:"passedTestCount"`
	TimeConsumedMillis  int     `json:"timeConsumedMillis"`
	MemoryConsumedBytes int     `json:"memoryConsumedBytes"`
	Points              float64 `json:"points"`
}

type Hack struct {
	ID                  int     `json:"id"`
	CreationTimeSeconds int     `json:"creationTimeSeconds"`
	Hacker              Party   `json:"hacker"`
	Defender            Party   `json:"defender"`
	Verdict             string  `json:"verdict"`
	Problem             Problem `json:"problem"`
	Test                string  `json:"test"`
	JudgeProtocol       struct {
		Manual   string `json:"manual"`
		Protocol string `json:"protocol"`
		Verdict  string `json:"verdict"`
	} `json:"judgeProtocol"`
}

type RanklistRow struct {
	Party                     Party           `json:"party"`
	Rank                      int             `json:"rank"`
	Points                    float64         `json:"points"`
	Penalty                   int             `json:"penalty"`
	SuccessfulHackCount       int             `json:"successfulHackCount"`
	UnsuccessfulHackCount     int             `json:"unsuccessfulHackCount"`
	ProblemResults            []ProblemResult `json:"problemResults"`
	LastSubmissionTimeSeconds int             `json:"lastSubmissionTimeSeconds"`
}

type Standings struct {
	Contest  Contest       `json:"contest"`
	Problems []Problem     `json:"problems"`
	Rows     []RanklistRow `json:"rows"`
}
