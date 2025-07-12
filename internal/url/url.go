package url

import (
	"fmt"
	"regexp"
	"strconv"
)

type githubURL struct {
	Owner       string
	Repo        string
	IssueNumber int
}

func DivideURL(url string) (*githubURL, error) {
	if len(url) <= 0 {
		return nil, fmt.Errorf(" No URL provided. Usage: gh issue-clone <issue-url>")
	}
	// Extract owner, repo, and issue number from URL
	re := regexp.MustCompile(`github\.com/([^/]+)/([^/]+)/issues/([0-9]+)`)
	matches := re.FindStringSubmatch(url)

	if matches == nil || len(matches) != 4 {
		return nil, fmt.Errorf(" Invalid GitHub issue URL format. Expected: https://github.com/owner/repo/issues/number")
	}
	issueNum, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, err
	}

	return &githubURL{
		Owner:       matches[1],
		Repo:        matches[2],
		IssueNumber: issueNum,
	}, nil
}
