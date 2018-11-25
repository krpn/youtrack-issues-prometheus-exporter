package monitoring

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/krpn/youtrack-issues-prometheus-exporter/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	issueser := NewMockgetIssueser(ctrl)
	metricser := NewMockmetricser(ctrl)

	type testTableData struct {
		queries  map[string]string
		expected *Monitoring
	}

	testTable := []testTableData{
		{
			queries: map[string]string{
				"test query 1": "#Unresolved",
				"test query 2": "#Unassigned",
			},
			expected: &Monitoring{
				issueser:  issueser,
				metricser: metricser,
				lastActiveIssues: map[string]map[string]model.Issue{
					"test query 1": {},
					"test query 2": {},
				},
				queries: map[string]string{
					"test query 1": "#Unresolved",
					"test query 2": "#Unassigned",
				},
			},
		},
	}

	for _, testUnit := range testTable {
		assert.Equal(t, testUnit.expected, New(issueser, metricser, testUnit.queries))
	}
}

func TestMonitoring_RefreshMetrics(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type testTableData struct {
		tcase                    string
		lastActiveIssues         map[string]map[string]model.Issue
		queries                  map[string]string
		expectFunc               func(i *MockgetIssueser, m *Mockmetricser)
		expectedLastActiveIssues map[string]map[string]model.Issue
	}

	testTable := []testTableData{
		{
			tcase: "1 disable, 1 new, 1 not changed, 1 renamed",
			lastActiveIssues: map[string]map[string]model.Issue{
				"test query 1": {
					"YT-100 For disable": model.Issue{
						ID:    "YT-100",
						Title: "For disable",
					},
				},
				"test query 2": {
					"YT-200 Not changed": model.Issue{
						ID:    "YT-200",
						Title: "Not changed",
					},
					"YT-300 Renamed": model.Issue{
						ID:    "YT-300",
						Title: "Renamed",
					},
				},
			},
			queries: map[string]string{
				"test query 1": "#Unresolved",
				"test query 2": "#Unassigned",
			},
			expectFunc: func(i *MockgetIssueser, m *Mockmetricser) {
				// test query 1
				i.EXPECT().GetIssues("#Unresolved").Return(
					map[string]model.Issue{
						"YT-101 New": {
							ID:    "YT-101",
							Title: "New",
						},
					},
					nil,
				)
				m.EXPECT().DisableMonitoring("test query 1", model.Issue{
					ID:    "YT-100",
					Title: "For disable",
				})
				m.EXPECT().EnableMonitoring("test query 1", model.Issue{
					ID:    "YT-101",
					Title: "New",
				})

				// test query 2
				i.EXPECT().GetIssues("#Unassigned").Return(
					map[string]model.Issue{
						"YT-200 Not changed": {
							ID:    "YT-200",
							Title: "Not changed",
						},
						"YT-300 New name": {
							ID:    "YT-300",
							Title: "New name",
						},
					},
					nil,
				)
				m.EXPECT().DisableMonitoring("test query 2", model.Issue{
					ID:    "YT-300",
					Title: "Renamed",
				})
				m.EXPECT().EnableMonitoring("test query 2", model.Issue{
					ID:    "YT-300",
					Title: "New name",
				})
			},
			expectedLastActiveIssues: map[string]map[string]model.Issue{
				"test query 1": {
					"YT-101 New": model.Issue{
						ID:    "YT-101",
						Title: "New",
					},
				},
				"test query 2": {
					"YT-200 Not changed": model.Issue{
						ID:    "YT-200",
						Title: "Not changed",
					},
					"YT-300 New name": model.Issue{
						ID:    "YT-300",
						Title: "New name",
					},
				},
			},
		},
		{
			tcase: "get issues error",
			lastActiveIssues: map[string]map[string]model.Issue{
				"test query 1": {
					"YT-100 For disable": model.Issue{
						ID:    "YT-100",
						Title: "For disable",
					},
				},
				"test query 2": {
					"YT-200 Not changed": model.Issue{
						ID:    "YT-200",
						Title: "Not changed",
					},
					"YT-300 Renamed": model.Issue{
						ID:    "YT-300",
						Title: "Renamed",
					},
				},
			},
			queries: map[string]string{
				"test query 1": "#Unresolved",
				"test query 2": "#Unassigned",
			},
			expectFunc: func(i *MockgetIssueser, m *Mockmetricser) {
				i.EXPECT().GetIssues("#Unresolved").Return(nil, errors.New("test query 1 error"))
				m.EXPECT().ErrorInc("test query 1", errors.New("test query 1 error"))
				i.EXPECT().GetIssues("#Unassigned").Return(nil, errors.New("test query 2 error"))
				m.EXPECT().ErrorInc("test query 2", errors.New("test query 2 error"))
			},
			expectedLastActiveIssues: map[string]map[string]model.Issue{
				"test query 1": {
					"YT-100 For disable": model.Issue{
						ID:    "YT-100",
						Title: "For disable",
					},
				},
				"test query 2": {
					"YT-200 Not changed": model.Issue{
						ID:    "YT-200",
						Title: "Not changed",
					},
					"YT-300 Renamed": model.Issue{
						ID:    "YT-300",
						Title: "Renamed",
					},
				},
			},
		},
	}

	for _, testUnit := range testTable {
		issueser := NewMockgetIssueser(ctrl)
		metricser := NewMockmetricser(ctrl)

		monitoring := &Monitoring{
			issueser:         issueser,
			metricser:        metricser,
			lastActiveIssues: testUnit.lastActiveIssues,
			queries:          testUnit.queries,
		}

		testUnit.expectFunc(issueser, metricser)
		monitoring.RefreshMetrics()

		assert.Equal(t, testUnit.expectedLastActiveIssues, monitoring.lastActiveIssues, testUnit.tcase)
	}
}
