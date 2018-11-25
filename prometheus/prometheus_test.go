package prometheus

import (
	e "errors"
	"github.com/golang/mock/gomock"
	"github.com/krpn/youtrack-issues-prometheus-exporter/model"
	"testing"
)

func TestPrometheus_ConsistentLabelCardinality(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			t.Fatal(r)
		}
	}()

	p := New()
	issue := model.Issue{
		ID:    "YT-100",
		Title: "Test issue",
	}
	queryName := "unassigned"

	p.EnableMonitoring(queryName, issue)
	p.DisableMonitoring(queryName, issue)
	p.ErrorInc(queryName, e.New("some error"))
}

func TestPrometheusMetrics_EnableMonitoring(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	issues := NewMockgaugeIniter(ctrl)
	prometheus := &Metrics{issues: issues}

	type testTableData struct {
		queryName  string
		issue      model.Issue
		expectFunc func(m *MockgaugeIniter)
	}

	testTable := []testTableData{
		{
			queryName: "test query",
			issue: model.Issue{
				ID:    "YT-100",
				Title: "Test issue",
			},
			expectFunc: func(m *MockgaugeIniter) {
				gauge := NewMockGauge(ctrl)
				m.EXPECT().WithLabelValues("test query", "YT-100", "Test issue").Return(gauge)
				gauge.EXPECT().Set(float64(1))
			},
		},
	}

	for _, testUnit := range testTable {
		testUnit.expectFunc(issues)
		prometheus.EnableMonitoring(testUnit.queryName, testUnit.issue)
	}
}

func TestPrometheusMetrics_DisableMonitoring(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	issues := NewMockgaugeIniter(ctrl)
	prometheus := &Metrics{issues: issues}

	type testTableData struct {
		queryName  string
		issue      model.Issue
		expectFunc func(gi *MockgaugeIniter)
	}

	testTable := []testTableData{
		{
			queryName: "test query",
			issue: model.Issue{
				ID:    "YT-100",
				Title: "Test issue",
			},
			expectFunc: func(gi *MockgaugeIniter) {
				gauge := NewMockGauge(ctrl)
				gi.EXPECT().WithLabelValues("test query", "YT-100", "Test issue").Return(gauge)
				gauge.EXPECT().Set(float64(0))
			},
		},
	}

	for _, testUnit := range testTable {
		testUnit.expectFunc(issues)
		prometheus.DisableMonitoring(testUnit.queryName, testUnit.issue)
	}
}

func TestPrometheusMetrics_ErrorInc(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	errors := NewMockcounterIniter(ctrl)
	prometheus := &Metrics{errors: errors}

	type testTableData struct {
		queryName  string
		error      error
		expectFunc func(ci *MockcounterIniter)
	}

	testTable := []testTableData{
		{
			queryName: "test query",
			error:     e.New("some error"),
			expectFunc: func(ci *MockcounterIniter) {
				counter := NewMockCounter(ctrl)
				ci.EXPECT().WithLabelValues("test query", "some error").Return(counter)
				counter.EXPECT().Inc()
			},
		},
	}

	for _, testUnit := range testTable {
		testUnit.expectFunc(errors)
		prometheus.ErrorInc(testUnit.queryName, testUnit.error)
	}
}
