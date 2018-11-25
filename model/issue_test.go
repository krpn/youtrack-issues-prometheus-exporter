package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIssue_FullID(t *testing.T) {
	t.Parallel()

	type testTableData struct {
		issue    Issue
		expected string
	}

	testTable := []testTableData{
		{
			issue: Issue{
				ID:    "YT-100",
				Title: "Test issue",
			},
			expected: "YT-100 Test issue",
		},
	}

	for _, testUnit := range testTable {
		assert.Equal(t, testUnit.expected, testUnit.issue.FullID())
	}
}
