package main

import (
	"slices"

	graphql "github.com/cli/shurcooL-graphql"
)

func getIssueInput(owner string, repo string, issueNumber int) map[string]any {
	return map[string]any{
		"owner":       graphql.String(owner),
		"repo":        graphql.String(repo),
		"issueNumber": graphql.Int(issueNumber),
	}
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

func createIssueInput(query GetIssueQuery, template *IssueTemplate) map[string]any {
	var assigneeIds []graphql.ID
	var body graphql.String
	var labelIds []graphql.ID
	var title graphql.String
	if template != nil {
		assigneeIds = generateIdSlice(template.Assinees.Nodes)
		body = template.Body
		labelIds = generateIdSlice(template.Labels.Nodes)
		if len(template.Title) > 0 {
			title = template.Title
		} else {
			title = query.Repository.Issue.Title
		}
	} else {
		assigneeIds = generateIdSlice(query.Repository.Issue.Assinees.Nodes)
		body = query.Repository.Issue.Body
		labelIds = generateIdSlice(query.Repository.Issue.Labels.Nodes)
		title = query.Repository.Issue.Title
	}
	return map[string]any{
		"input": CreateIssueInput{
			AssigneeIds:   assigneeIds,
			Body:          body,
			IssueTypeId:   query.Repository.Issue.IssueType.Id,
			LabelIds:      labelIds,
			MilestoneId:   query.Repository.Issue.Milestone.Id,
			ParentIssueId: query.Repository.Issue.Parent.Id,
			RepositoryId:  query.Repository.Id,
			Title:         title,
		},
	}
}

func addProjectV2ItemByIdInput(projectId graphql.ID, issueId graphql.ID) map[string]any {
	return map[string]any{
		"input": AddProjectV2ItemByIdInput{
			ProjectId: projectId,
			ContentId: issueId,
		},
	}
}

func structureFieldValue(fieldValue FieldValue) any {
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

func updateProjectV2ItemFieldValueInput(projectId graphql.ID, itemId graphql.ID, projectField FieldValue) map[string]any {
	fieldId := projectField.ProjectV2ItemFieldValueCommon.Field.ProjectV2FieldCommon.Id
	if fieldId == nil {
		return nil
	}
	value := structureFieldValue(projectField)
	if value == nil {
		return nil
	}
	return map[string]any{
		"input": UpdateProjectV2ItemFieldValueInput{
			ProjectId: projectId,
			ItemId:    itemId,
			FieldId:   fieldId,
			Value:     value,
		},
	}
}
