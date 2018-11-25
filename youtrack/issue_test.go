package youtrack

import (
	"github.com/krpn/youtrack-issues-prometheus-exporter/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApiIssue_ID(t *testing.T) {
	t.Parallel()

	type testTableData struct {
		issue    apiIssue
		expected string
	}

	testTable := []testTableData{
		{
			issue: apiIssue{
				Project: struct {
					ShortName string `json:"shortName"`
				}{ShortName: "YT"},
				Summary:         "Test issue",
				NumberInProject: 100,
			},
			expected: "YT-100",
		},
	}

	for _, testUnit := range testTable {
		assert.Equal(t, testUnit.expected, testUnit.issue.ID())
	}
}

func TestApiIssue_ToIssue(t *testing.T) {
	t.Parallel()

	type testTableData struct {
		issue    apiIssue
		expected model.Issue
	}

	testTable := []testTableData{
		{
			issue: apiIssue{
				Project: struct {
					ShortName string `json:"shortName"`
				}{ShortName: "YT"},
				Summary:         "Test issue",
				NumberInProject: 100,
			},
			expected: model.Issue{
				ID:    "YT-100",
				Title: "Test issue",
			},
		},
	}

	for _, testUnit := range testTable {
		assert.Equal(t, testUnit.expected, testUnit.issue.ToIssue())
	}
}
