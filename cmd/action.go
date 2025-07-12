package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/muleyuck/gh-issue-clone/internal/types"
	"github.com/muleyuck/gh-issue-clone/internal/url"
	"github.com/muleyuck/gh-issue-clone/internal/variables"
	"github.com/urfave/cli/v3"
)

func CloneIssue(ctx context.Context, c *cli.Command) error {
	issueURL := c.StringArg("url")
	g, err := url.DivideURL(issueURL)
	if err != nil {
		return err
	}

	fmt.Printf("Fetching issue details from %s/%s#%d...\n", g.Owner, g.Repo, g.IssueNumber)

	ops := api.ClientOptions{
		Timeout: 30 * time.Second,
	}
	client, err := api.NewGraphQLClient(ops)
	if err != nil {
		return err
	}
	var query types.GetIssueQuery
	err = client.Query("GetIssue", &query, variables.GetIssueInput(g.Owner, g.Repo, g.IssueNumber))
	if err != nil {
		return err
	}
	if query.Repository.Id == nil {
		return fmt.Errorf(" Repository not found or you don't have access to it")
	}

	templateName := c.String("template")
	template := variables.FindTemplateByName(query.Repository.IssueTemplates, templateName)

	fmt.Println("Creating new issue from fetched issue details...")

	var createMutation types.CreateIssueMutation
	err = client.Mutate("CreateIssue", &createMutation, variables.CreateIssueInput(query, template))
	if err != nil {
		return err
	}

	var addMutation types.AddProjectV2ItemByIdMutation
	var updateMutation types.UpdateProjectV2ItemFieldValueMutation
	var deleteMutation types.DeleteIssueMutation

	projectItems := query.Repository.Issue.ProjectItems.Nodes
	for _, projectItem := range projectItems {
		fmt.Printf("Found relevant project: %s. Add the issue to the project.\n", projectItem.Project.Title)
		issueId := createMutation.CreateIssue.Issue.Id
		errAdd := client.Mutate(
			"AddProjectV2ItemById",
			&addMutation,
			variables.AddProjectV2ItemByIdInput(projectItem.Project.Id, issueId),
		)
		if errAdd != nil {
			fmt.Println("✗ An error occurred while adding to project. delete the issue.")
			err = client.Mutate("DeleteIssue", &deleteMutation, variables.DeleteIssueInput(issueId))
			if err != nil {
				return err
			}
			return errAdd
		}

		for _, projectField := range projectItem.FieldValues.Nodes {
			input := variables.UpdateProjectV2ItemFieldValueInput(
				projectItem.Project.Id,
				addMutation.AddProjectV2ItemById.Item.Id,
				projectField,
			)
			if input == nil {
				continue
			}
			err = client.Mutate("UpdateProjectV2ItemFieldValue", &updateMutation, input)
			if err != nil {
				fieldName := projectField.ProjectV2ItemFieldValueCommon.Field.ProjectV2FieldCommon.Name
				log.Printf("Fail to update project field[%s]: %v\n", fieldName, err)
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
}
