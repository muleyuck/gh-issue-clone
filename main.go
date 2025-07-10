package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: gh issue-clone <issue-url>")
		return
	}

	issueURL := os.Args[1]

	// TODO: get tenplate from arguments

	// Extract owner, repo, and issue number from URL
	re := regexp.MustCompile(`github\.com/([^/]+)/([^/]+)/issues/([0-9]+)`)
	matches := re.FindStringSubmatch(issueURL)

	if matches == nil || len(matches) != 4 {
		log.Fatal("Invalid GitHub issue URL format. Expected: https://github.com/owner/repo/issues/number")
	}

	owner := matches[1]
	repo := matches[2]
	issueNumberStr := matches[3]

	issueNumber, err := strconv.Atoi(issueNumberStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Fetching issue details from %s/%s#%d...\n", owner, repo, issueNumber)

	ops := api.ClientOptions{
		Timeout: 30 * time.Second,
	}
	client, err := api.NewGraphQLClient(ops)
	if err != nil {
		log.Fatal(err)
	}
	var query GetIssueQuery
	err = client.Query("GetIssue", &query, getIssueInput(owner, repo, issueNumber))
	if err != nil {
		log.Fatal(err)
	}
	if query.Repository.Id == nil {
		log.Fatal("repository id is null. issue should be belong to any repository.")
	}
	// var template struct{}
	// if len(issueTemplate) > 0 {
	// 	m := map[graphql.String]struct{}{}
	// 	for _, template := range query.Repository.IssueTemplates {
	// 		m[template.Name] = struct{}{}
	// 	}
	// 	val, ok := m[graphql.String(issueTemplate)]
	// 	if !ok {
	// 		return
	// 	}
	// 	template = val
	// }

	var createMutation CreateIssueMutation
	err = client.Mutate("CreateIssue", &createMutation, createIssueInput(query))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created issue with ID: %s\n", createMutation.CreateIssue.Issue.Id)

	var addMutation AddProjectV2ItemByIdMutation
	var updateMutation UpdateProjectV2ItemFieldValueMutation

	projectItems := query.Repository.Issue.ProjectItems.Nodes
	for _, projectItem := range projectItems {
		err = client.Mutate("AddProjectV2ItemById", &addMutation, addProjectV2ItemByIdInput(projectItem.Project.Id, createMutation.CreateIssue.Issue.Id))
		if err != nil {
			// TODO: Remove Issue to Rollback
			log.Fatal(err)
		}

		fmt.Printf("Add Project as ID: %s\n", addMutation.AddProjectV2ItemById.Item.Id)
		for _, projectField := range projectItem.FieldValues.Nodes {
			input := updateProjectV2ItemFieldValueInput(projectItem.Project.Id, addMutation.AddProjectV2ItemById.Item.Id, projectField)
			if input == nil {
				continue
			}
			err = client.Mutate("UpdateProjectV2ItemFieldValue", &updateMutation, input)
			if err != nil {
				log.Printf("Fail to update project field[%s]: %v\n", projectField.ProjectV2ItemFieldValueCommon.Field.ProjectV2FieldCommon.Name, err)
				continue
			}
		}
	}
	fmt.Printf("Issue cloned successfully: %s\n", createMutation.CreateIssue.Issue.Url)
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
