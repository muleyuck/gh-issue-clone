package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/muleyuck/gh-issue-clone/versions"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:    "issue-clone",
		Usage:   "clone GitHub issues from a given issue.",
		Version: versions.AppVersion,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "template", Aliases: []string{"t"}, Usage: "Issue Template Name"},
		},
		Arguments: []cli.Argument{
			&cli.StringArg{Name: "url"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			issueURL := cmd.StringArg("url")
			if len(issueURL) <= 0 {
				return fmt.Errorf(" No URL provided. Usage: gh issue-clone <issue-url>")
			}
			templateName := cmd.String("template")
			// Extract owner, repo, and issue number from URL
			re := regexp.MustCompile(`github\.com/([^/]+)/([^/]+)/issues/([0-9]+)`)
			matches := re.FindStringSubmatch(issueURL)

			if matches == nil || len(matches) != 4 {
				return fmt.Errorf(" Invalid GitHub issue URL format. Expected: https://github.com/owner/repo/issues/number")
			}

			owner := matches[1]
			repo := matches[2]
			issueNumberStr := matches[3]

			issueNumber, err := strconv.Atoi(issueNumberStr)
			if err != nil {
				return err
			}

			fmt.Printf("Fetching issue details from %s/%s#%d...\n", owner, repo, issueNumber)

			ops := api.ClientOptions{
				Timeout: 30 * time.Second,
			}
			client, err := api.NewGraphQLClient(ops)
			if err != nil {
				return err
			}
			var query GetIssueQuery
			err = client.Query("GetIssue", &query, getIssueInput(owner, repo, issueNumber))
			if err != nil {
				return err
			}
			if query.Repository.Id == nil {
				return fmt.Errorf(" Repository not found or you don't have access to it")
			}

			template := FindByName(query.Repository.IssueTemplates, templateName)

			fmt.Println("Creating new issue from fetched issue details...")

			var createMutation CreateIssueMutation
			err = client.Mutate("CreateIssue", &createMutation, createIssueInput(query, template))
			if err != nil {
				return err
			}

			var addMutation AddProjectV2ItemByIdMutation
			var updateMutation UpdateProjectV2ItemFieldValueMutation
			var deleteMutation DeleteIssueMutation

			projectItems := query.Repository.Issue.ProjectItems.Nodes
			for _, projectItem := range projectItems {
				fmt.Printf("Found relevant project: %s. Add the issue to the project.\n", projectItem.Project.Title)
				issueId := createMutation.CreateIssue.Issue.Id
				err = client.Mutate("AddProjectV2ItemById", &addMutation, addProjectV2ItemByIdInput(projectItem.Project.Id, issueId))
				if err != nil {
					fmt.Println("✗ An error occurred while adding to project. delete the issue.")
					err = client.Mutate("DeleteIssue", &deleteMutation, deleteIssueInput(issueId))
					if err != nil {
						return err
					}
					return err
				}

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
			fmt.Println("----------------------------------------------------------------------------")
			fmt.Println("✓ Issue cloned successfully")
			fmt.Printf("ID:    #%d\n", createMutation.CreateIssue.Issue.Number)
			fmt.Printf("Title: %s\n", createMutation.CreateIssue.Issue.Title)
			fmt.Printf("URL:   %s\n", createMutation.CreateIssue.Issue.Url)
			fmt.Println("----------------------------------------------------------------------------")
			return nil
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Printf("✗ Error:%s\n", err)
		os.Exit(1)
	}
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
