package youtrack

import (
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/krpn/youtrack-issues-prometheus-exporter/model"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	makeRequester := NewMockmakeRequester(ctrl)

	type testTableData struct {
		tcase            string
		endpoint         string
		token            string
		expectedYouTrack *YouTrack
		expectedErr      error
	}

	testTable := []testTableData{
		{
			tcase:    "success",
			endpoint: "http://www.test.com/",
			token:    "abc",
			expectedYouTrack: &YouTrack{
				requester: makeRequester,
				url: url.URL{
					Scheme: "http",
					Host:   "www.test.com",
					Path:   apiPath,
				},
				getParams: map[string][]string{
					"fields": {"project(shortName),numberInProject,summary"},
				},
				headers: map[string]string{
					"Accept":        "application/json",
					"Content-Type":  "application/json",
					"Authorization": "Bearer abc",
				},
			},
			expectedErr: nil,
		},
		{
			tcase:            "endpoint error",
			endpoint:         "http://www test com/",
			token:            "abc",
			expectedYouTrack: nil,
			expectedErr:      &url.Error{Op: "parse", URL: "http://www test com/", Err: url.InvalidHostError(" ")},
		},
	}

	for _, testUnit := range testTable {
		youTrack, err := New(testUnit.endpoint, testUnit.token, makeRequester)

		assert.Equal(t, testUnit.expectedYouTrack, youTrack, testUnit.tcase)
		assert.Equal(t, testUnit.expectedErr, err, testUnit.tcase)
	}
}

func TestYouTrack_GetIssues(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	makeRequester := NewMockmakeRequester(ctrl)

	headers := map[string]string{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": "Bearer abc",
	}

	youTrack := &YouTrack{
		requester: makeRequester,
		url: url.URL{
			Scheme: "http",
			Host:   "www.test.com",
			Path:   apiPath,
		},
		getParams: map[string][]string{
			"fields": {"project(shortName),numberInProject,summary"},
		},
		headers: headers,
	}

	type testTableData struct {
		tcase          string
		query          string
		expectFunc     func(mr *MockmakeRequester)
		expectedIssues map[string]model.Issue
		expectedErr    error
	}

	testTable := []testTableData{
		{
			tcase: "success",
			query: "Priority: Show-Stopper #Unresolved #Unassigned",
			expectFunc: func(mr *MockmakeRequester) {
				mr.EXPECT().MakeRequest(
					"http://www.test.com/api/issues?fields=project%28shortName%29%2CnumberInProject%2Csummary&query=Priority%3A+Show-Stopper+%23Unresolved+%23Unassigned",
					headers,
				).Return([]byte(`[
    {
        "project": {
            "shortName": "YT",
            "$type": "jetbrains.charisma.persistent.Project"
        },
        "summary": "Test issue 1",
        "numberInProject": 100,
        "$type": "jetbrains.charisma.persistent.Issue"
    },
    {
        "project": {
            "shortName": "YT",
            "$type": "jetbrains.charisma.persistent.Project"
        },
        "summary": "Test issue 2",
        "numberInProject": 200,
        "$type": "jetbrains.charisma.persistent.Issue"
    }
]`), nil)
			},
			expectedIssues: map[string]model.Issue{
				"YT-100 Test issue 1": {ID: "YT-100", Title: "Test issue 1"},
				"YT-200 Test issue 2": {ID: "YT-200", Title: "Test issue 2"},
			},
			expectedErr: nil,
		},
		{
			tcase: "request error",
			query: "Priority: Show-Stopper #Unresolved #Unassigned",
			expectFunc: func(mr *MockmakeRequester) {
				mr.EXPECT().MakeRequest(
					"http://www.test.com/api/issues?fields=project%28shortName%29%2CnumberInProject%2Csummary&query=Priority%3A+Show-Stopper+%23Unresolved+%23Unassigned",
					headers,
				).Return(nil, errors.New("request error"))
			},
			expectedIssues: nil,
			expectedErr:    errors.New("request error"),
		},
		{
			tcase: "incorrect response",
			query: "Priority: Show-Stopper #Unresolved #Unassigned",
			expectFunc: func(mr *MockmakeRequester) {
				mr.EXPECT().MakeRequest(
					"http://www.test.com/api/issues?fields=project%28shortName%29%2CnumberInProject%2Csummary&query=Priority%3A+Show-Stopper+%23Unresolved+%23Unassigned",
					headers,
				).Return([]byte(``), nil)
			},
			expectedIssues: nil,
			expectedErr:    json.Unmarshal([]byte(``), nil),
		},
	}

	for _, testUnit := range testTable {
		testUnit.expectFunc(makeRequester)
		issues, err := youTrack.GetIssues(testUnit.query)
		assert.Equal(t, testUnit.expectedIssues, issues, testUnit.tcase)
		assert.Equal(t, testUnit.expectedErr, err, testUnit.tcase)
	}
}
