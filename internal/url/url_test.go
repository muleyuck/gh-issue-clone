package url

import (
	"testing"
)

func TestDivideURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		expected    githubURL
		expectedErr bool
	}{
		{
			name:        "Empty URL",
			url:         "",
			expectedErr: true,
		},
		{
			name:        "Non-GitHub URL",
			url:         "https://example.com/owner/repo/issues/123",
			expectedErr: true,
		},
		{
			name:        "Invalid URL format",
			url:         "https://github.com/owner/repo/pull/123",
			expectedErr: true,
		},
		{
			name: "Valid GitHub issue URL",
			url:  "https://github.com/owner/repo/issues/123",
			expected: githubURL{
				Owner:       "owner",
				Repo:        "repo",
				IssueNumber: 123,
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DivideURL(tt.url)

			if (err != nil) != tt.expectedErr {
				t.Errorf("DivideURL() error = %v, wantErr %v", err, tt.expectedErr)
				return
			}

			if err == nil {
				if got.Owner != tt.expected.Owner {
					t.Errorf("DivideURL() owner = %v, want %v", got.Owner, tt.expected.Owner)
				}
				if got.Repo != tt.expected.Repo {
					t.Errorf("DivideURL() repo = %v, want %v", got.Repo, tt.expected.Repo)
				}
				if got.IssueNumber != tt.expected.IssueNumber {
					t.Errorf("DivideURL() issueNumber = %v, want %v", got.IssueNumber, tt.expected.IssueNumber)
				}
			}
		})
	}
}
