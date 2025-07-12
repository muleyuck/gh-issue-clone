package main

import graphql "github.com/cli/shurcooL-graphql"

type FieldValue struct {
	Typename                      graphql.String `graphql:"__typename"`
	ProjectV2ItemFieldValueCommon struct {
		Field struct {
			ProjectV2FieldCommon struct {
				Id       graphql.ID
				Name     graphql.String
				DataType graphql.String
			} `graphql:"... on ProjectV2FieldCommon"`
		}
	} `graphql:"... on ProjectV2ItemFieldValueCommon"`
	ProjectV2ItemFieldDateValue struct {
		Date graphql.String
	} `graphql:"... on ProjectV2ItemFieldDateValue"`
	ProjectV2ItemFieldIterationValue struct {
		IterationId graphql.ID
	} `graphql:"... on ProjectV2ItemFieldIterationValue"`
	ProjectV2ItemFieldNumberValue struct {
		Number graphql.Float
	} `graphql:"... on ProjectV2ItemFieldNumberValue"`
	ProjectV2ItemFieldSingleSelectValue struct {
		OptionId graphql.ID
	} `graphql:"... on ProjectV2ItemFieldSingleSelectValue"`
	ProjectV2ItemFieldTextValue struct {
		Text graphql.String
	} `graphql:"... on ProjectV2ItemFieldTextValue"`
}

type ProjectItem struct {
	Project struct {
		Id    graphql.ID
		Title graphql.String
	}
	FieldValues struct {
		Nodes []FieldValue
	} `graphql:"fieldValues(first: 100)"`
}

type IssueTemplate struct {
	Name     graphql.String
	Title    graphql.String
	Body     graphql.String
	Assinees struct {
		Nodes []struct {
			Id graphql.ID
		}
	} `graphql:"assignees(first: 10)"`
	Labels struct {
		Nodes []struct {
			Id graphql.ID
		}
	} `graphql:"labels(first: 10)"`
}

func FindByName(templates []IssueTemplate, name string) *IssueTemplate {
	m := map[graphql.String]IssueTemplate{}
	for _, template := range templates {
		m[template.Name] = template
	}
	val, ok := m[graphql.String(name)]
	if !ok {
		return nil
	}
	return &val
}

type GetIssueQuery struct {
	Repository struct {
		Id    graphql.ID
		Issue struct {
			Title     graphql.String
			Body      graphql.String
			IssueType struct {
				Id graphql.ID
			}
			Assinees struct {
				Nodes []struct {
					Id graphql.ID
				}
			} `graphql:"assignees(first: 10)"`
			Labels struct {
				Nodes []struct {
					Id graphql.ID
				}
			} `graphql:"labels(first: 10)"`
			Parent struct {
				Id graphql.ID
			}
			ProjectItems struct {
				Nodes []ProjectItem
			} `graphql:"projectItems(first: 10,  includeArchived: false)"`
			Milestone struct {
				Id graphql.ID
			}
			SubIssues struct {
				Nodes []struct {
					Id graphql.ID
				}
			} `graphql:"subIssues(first: 10)"`
		} `graphql:"issue(number: $issueNumber)"`
		IssueTemplates []IssueTemplate
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

type GetIssueTemplateQuery struct{}

type CreateIssueInput struct {
	AssigneeIds   []graphql.ID   `json:"assigneeIds"`
	Body          any            `json:"body"`
	IssueTemplate any            `json:"issueTemplate"`
	IssueTypeId   graphql.ID     `json:"issueTypeId"`
	LabelIds      []graphql.ID   `json:"labelIds"`
	MilestoneId   graphql.ID     `json:"milestoneId"`
	ParentIssueId graphql.ID     `json:"parentIssueId"`
	RepositoryId  graphql.ID     `json:"repositoryId"`
	Title         graphql.String `json:"title"`
}

type CreateIssueMutation struct {
	CreateIssue struct {
		Issue struct {
			Id     graphql.ID
			Number graphql.Int
			Title  graphql.String
			Url    graphql.String
		}
	} `graphql:"createIssue(input: $input)"`
}

type DeleteIssueInput struct {
	IssueId graphql.ID `json:"issueId"`
}

type DeleteIssueMutation struct {
	DeleteIssue struct {
		ClientMutationId graphql.String
	} `graphql:"deleteIssue(input: $input)"`
}

type AddProjectV2ItemByIdInput struct {
	ProjectId graphql.ID `json:"projectId"`
	ContentId graphql.ID `json:"contentId"`
}

type AddProjectV2ItemByIdMutation struct {
	AddProjectV2ItemById struct {
		Item struct {
			Id graphql.ID
		}
	} `graphql:"addProjectV2ItemById(input: $input)"`
}

type UpdateProjectV2ItemFieldValueInput struct {
	ProjectId graphql.ID `json:"projectId"`
	ItemId    graphql.ID `json:"itemId"`
	FieldId   graphql.ID `json:"fieldId"`
	Value     any        `json:"value"`
}

type UpdateProjectV2ItemFieldValueMutation struct {
	UpdateProjectV2ItemFieldValue struct {
		ClientMutationId graphql.String
	} `graphql:"updateProjectV2ItemFieldValue(input: $input)"`
}
