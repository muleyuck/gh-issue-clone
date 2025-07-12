package variables

import (
	"reflect"
	"testing"

	graphql "github.com/cli/shurcooL-graphql"
	"github.com/muleyuck/gh-issue-clone/internal/types"
)

func TestCreateIssueInput(t *testing.T) {
	testRepositoryId := graphql.ID("repo-id")
	testMilestoneID := graphql.ID("milestone-id")
	testParentID := graphql.ID("parent-id")
	testIssueTypeID := graphql.ID("issue-type-id")
	testTemplateName := graphql.String("Test Template")

	getIssueQuery := types.GetIssueQuery{
		Repository: types.RepositoryQuery{
			Id: testRepositoryId,
			Issue: types.IssueQuery{
				Title: "Test Issue",
				Body:  "Test Body",
				Assinees: struct {
					Nodes []struct{ Id graphql.ID }
				}{
					Nodes: []struct{ Id graphql.ID }{
						{Id: "test-assignee-id"},
					},
				},
				Labels: struct {
					Nodes []struct{ Id graphql.ID }
				}{
					Nodes: []struct{ Id graphql.ID }{
						{Id: "test-label-id"},
					},
				},
				IssueType: struct{ Id graphql.ID }{
					Id: testIssueTypeID,
				},
				Milestone: struct{ Id graphql.ID }{
					Id: testMilestoneID,
				},
				Parent: struct{ Id graphql.ID }{
					Id: testParentID,
				},
			},
		},
	}

	tests := []struct {
		name     string
		query    types.GetIssueQuery
		template *types.IssueTemplate
		expected map[string]any
	}{
		{
			name:     "Without template",
			query:    getIssueQuery,
			template: nil,
			expected: map[string]any{
				"input": types.CreateIssueInput{
					AssigneeIds:   []graphql.ID{"test-assignee-id"},
					Body:          "Test Body",
					IssueTypeId:   testIssueTypeID,
					LabelIds:      []graphql.ID{"test-label-id"},
					MilestoneId:   testMilestoneID,
					ParentIssueId: testParentID,
					RepositoryId:  testRepositoryId,
					Title:         "Test Issue",
				},
			},
		},
		{
			name:  "With template",
			query: getIssueQuery,
			template: &types.IssueTemplate{
				Name:  testTemplateName,
				Title: "Template Title",
				Body:  "Template Body",
				Assinees: struct {
					Nodes []struct{ Id graphql.ID }
				}{
					Nodes: []struct{ Id graphql.ID }{
						{Id: "template-assignee-id"},
					},
				},
				Labels: struct {
					Nodes []struct{ Id graphql.ID }
				}{
					Nodes: []struct{ Id graphql.ID }{
						{Id: "template-label-id"},
					},
				},
			},
			expected: map[string]any{
				"input": types.CreateIssueInput{
					AssigneeIds:   []graphql.ID{"template-assignee-id"},
					Body:          "Template Body",
					IssueTypeId:   testIssueTypeID,
					LabelIds:      []graphql.ID{"template-label-id"},
					MilestoneId:   testMilestoneID,
					ParentIssueId: testParentID,
					RepositoryId:  testRepositoryId,
					Title:         "Template Title",
				},
			},
		},
		{
			name:  "With template without title",
			query: getIssueQuery,
			template: &types.IssueTemplate{
				Name:  testTemplateName,
				Title: "",
				Body:  "Template Body",
				Assinees: struct {
					Nodes []struct{ Id graphql.ID }
				}{
					Nodes: []struct{ Id graphql.ID }{
						{Id: "template-assignee-id"},
					},
				},
				Labels: struct {
					Nodes []struct{ Id graphql.ID }
				}{
					Nodes: []struct{ Id graphql.ID }{
						{Id: "template-label-id"},
					},
				},
			},
			expected: map[string]any{
				"input": types.CreateIssueInput{
					AssigneeIds:   []graphql.ID{"template-assignee-id"},
					Body:          "Template Body",
					IssueTypeId:   testIssueTypeID,
					LabelIds:      []graphql.ID{"template-label-id"},
					MilestoneId:   testMilestoneID,
					ParentIssueId: testParentID,
					RepositoryId:  testRepositoryId,
					Title:         "Test Issue",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateIssueInput(tt.query, tt.template)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("%s: createIssueInput() = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}
