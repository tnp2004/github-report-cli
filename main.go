package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

type Commit struct {
	CommitInfo struct {
		Committer struct {
			Name string `json:"name"`
			Date string `json:"date"`
		} `json:"committer"`
		Message string `json:"message"`
	} `json:"commit"`
}

func main() {
	// github tokens for access to private repo (alternative)
	tokens := ""

	// github account owner
	owner := "tnp2004"
	// The repository that you want to watch commits
	repo := "github-report-cli"

	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.Suffix = " fetching..."
	spinner.Start()
	fmt.Println("")

	argLen := len(os.Args[1:])
	if argLen == 1 {
		color.HiRed("owner, repo are needed!!")
		return
	} else if argLen == 2 {
		owner = os.Args[1]
		repo = os.Args[2]
	} else if argLen > 2 {
		color.HiRed("too many args only need owner and repo!!")
		return
	}

	var commits []Commit
	queryCommits(tokens, owner, repo, &commits)
	printCommits(commits, owner, repo)
	spinner.Stop()
}

func printCommits(commits []Commit, owner, repo string) {
	w := fmt.Sprintf("\n\n/| %s : %s |\\", owner, repo)
	color.Cyan(w)
	for i := len(commits) - 1; i >= 0; i-- {
		commitTimeUTC, _ := time.Parse(time.RFC3339, commits[i].CommitInfo.Committer.Date)
		// chage to your time zone
		timezone, _ := time.LoadLocation("Asia/Bangkok")
		commitTimeInTimeZone := commitTimeUTC.In(timezone)
		commitDate := commitTimeInTimeZone.Format("2006-01-02")
		commitTimeFormatted := commitTimeInTimeZone.Format("Monday, 2 January, 2006 15:4:5 PM")
		text := fmt.Sprintf("\n🪅%s 🕒%s \n📌COMMIT %v \n", commits[i].CommitInfo.Committer.Name, commitTimeFormatted, commits[i].CommitInfo.Message)

		today := time.Now().Format("2006-01-02")
		if commitDate == today {
			// Highlight today
			color.HiCyan(text)
		} else {
			fmt.Print(text)
		}
	}
	fmt.Println("")
}

func queryCommits(tokens, owner, repo string, commits *[]Commit) {
	client := http.Client{}
	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", owner, repo)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println(err)
	}

	var res *http.Response
	if len(tokens) != 0 {
		req.Header = http.Header{
			"Authorization": {"Bearer " + tokens},
		}
	}

	res, err = client.Do(req)

	if err != nil {
		fmt.Println("Something went wrong")
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("\nRepository not found")
		return
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(body, commits)
	if err != nil {
		fmt.Println(err)
	}
}
