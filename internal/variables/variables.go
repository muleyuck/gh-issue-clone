package variables

import (
	"slices"

	graphql "github.com/cli/shurcooL-graphql"
	"github.com/muleyuck/gh-issue-clone/internal/types"
)

func GetIssueInput(owner string, repo string, issueNumber int) map[string]any {
	return map[string]any{
		"owner":       graphql.String(owner),
		"repo":        graphql.String(repo),
		"issueNumber": graphql.Int(issueNumber),
	}
}

func FindTemplateByName(templates []types.IssueTemplate, name string) *types.IssueTemplate {
	if len(name) == 0 {
		return nil
	}
	m := map[graphql.String]types.IssueTemplate{}
	for _, template := range templates {
		m[template.Name] = template
	}
	val, ok := m[graphql.String(name)]
	if !ok {
		return nil
	}
	return &val
}

func generateIdSlice(targets []struct{ Id graphql.ID }) []graphql.ID {
	return slices.Collect(func(yield func(graphql.ID) bool) {
		for _, obj := range targets {
			if obj.Id != nil {
				if !yield(obj.Id) {
					return
				}
			}
		}
	})
}

type templateInput struct {
	assigneeIds []graphql.ID
	body        graphql.String
	labelIds    []graphql.ID
	title       graphql.String
}

func getTemplateInput(template *types.IssueTemplate) *templateInput {
	if template == nil {
		return nil
	}
	return &templateInput{
		assigneeIds: generateIdSlice(template.Assinees.Nodes),
		body:        template.Body,
		labelIds:    generateIdSlice(template.Labels.Nodes),
		title:       template.Title,
	}
}

func CreateIssueInput(query types.GetIssueQuery, template *types.IssueTemplate) map[string]any {
	t := getTemplateInput(template)
	if t == nil {
		t = &templateInput{
			assigneeIds: generateIdSlice(query.Repository.Issue.Assinees.Nodes),
			body:        query.Repository.Issue.Body,
			labelIds:    generateIdSlice(query.Repository.Issue.Labels.Nodes),
			title:       query.Repository.Issue.Title,
		}
	}
	if len(t.title) == 0 {
		t.title = query.Repository.Issue.Title
	}
	return map[string]any{
		"input": types.CreateIssueInput{
			AssigneeIds:   t.assigneeIds,
			Body:          t.body,
			IssueTypeId:   query.Repository.Issue.IssueType.Id,
			LabelIds:      t.labelIds,
			MilestoneId:   query.Repository.Issue.Milestone.Id,
			ParentIssueId: query.Repository.Issue.Parent.Id,
			RepositoryId:  query.Repository.Id,
			Title:         t.title,
		},
	}
}

func DeleteIssueInput(issueId graphql.ID) map[string]any {
	return map[string]any{
		"input": types.DeleteIssueInput{
			IssueId: issueId,
		},
	}
}

func AddProjectV2ItemByIdInput(projectId graphql.ID, issueId graphql.ID) map[string]any {
	return map[string]any{
		"input": types.AddProjectV2ItemByIdInput{
			ProjectId: projectId,
			ContentId: issueId,
		},
	}
}

func structureFieldValue(fieldValue types.FieldValue) any {
	switch fieldValue.Typename {
	case "ProjectV2ItemFieldDateValue":
		return struct {
			Date graphql.String `json:"date"`
		}{Date: fieldValue.ProjectV2ItemFieldDateValue.Date}
	case "ProjectV2ItemFieldIterationValue":
		return struct {
			IterationId graphql.ID `json:"iterationId"`
		}{IterationId: fieldValue.ProjectV2ItemFieldIterationValue.IterationId}
	case "ProjectV2ItemFieldNumberValue":
		return struct {
			Number graphql.Float `json:"number"`
		}{Number: fieldValue.ProjectV2ItemFieldNumberValue.Number}
	case "ProjectV2ItemFieldSingleSelectValue":
		return struct {
			SingleSelectOptionId graphql.ID `json:"singleSelectOptionId"`
		}{SingleSelectOptionId: fieldValue.ProjectV2ItemFieldSingleSelectValue.OptionId}
	case "ProjectV2ItemFieldTextValue":
		if fieldValue.ProjectV2ItemFieldValueCommon.Field.ProjectV2FieldCommon.DataType == "TEXT" {
			return struct {
				Text graphql.String `json:"text"`
			}{Text: fieldValue.ProjectV2ItemFieldTextValue.Text}
		}
	}
	return nil
}

func UpdateProjectV2ItemFieldValueInput(projectId graphql.ID, itemId graphql.ID, projectField types.FieldValue) map[string]any {
	fieldId := projectField.ProjectV2ItemFieldValueCommon.Field.ProjectV2FieldCommon.Id
	if fieldId == nil {
		return nil
	}
	value := structureFieldValue(projectField)
	if value == nil {
		return nil
	}
	return map[string]any{
		"input": types.UpdateProjectV2ItemFieldValueInput{
			ProjectId: projectId,
			ItemId:    itemId,
			FieldId:   fieldId,
			Value:     value,
		},
	}
}
