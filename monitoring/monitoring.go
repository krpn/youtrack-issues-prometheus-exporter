package monitoring

import "github.com/krpn/youtrack-issues-prometheus-exporter/model"

//go:generate mockgen -source=monitoring.go -destination=monitoring_mocks.go -package=monitoring doc github.com/golang/mock/gomock

type getIssueser interface {
	GetIssues(query string) (issues map[string]model.Issue, err error)
}

type metricser interface {
	EnableMonitoring(queryName string, issue model.Issue)
	DisableMonitoring(queryName string, issue model.Issue)
	ErrorInc(queryName string, err error)
}

// Monitoring links YouTrack and Prometheus.
type Monitoring struct {
	issueser         getIssueser
	metricser        metricser
	lastActiveIssues map[string]map[string]model.Issue
	queries          map[string]string
}

// New creates Monitoring instance.
func New(issueser getIssueser, metricser metricser, queries map[string]string) *Monitoring {
	lastActiveIssues := make(map[string]map[string]model.Issue)
	for queryName := range queries {
		lastActiveIssues[queryName] = make(map[string]model.Issue)
	}

	return &Monitoring{
		issueser:         issueser,
		metricser:        metricser,
		lastActiveIssues: lastActiveIssues,
		queries:          queries,
	}
}

// RefreshMetrics gets actual issues and refreshes metrics.
func (m *Monitoring) RefreshMetrics() {
	for queryName, query := range m.queries {
		err := m.refreshMetrics(queryName, query)
		if err != nil {
			m.metricser.ErrorInc(queryName, err)
		}
	}
}

func (m *Monitoring) refreshMetrics(queryName, query string) error {
	issues, err := m.issueser.GetIssues(query)
	if err != nil {
		return err
	}

	// Disable irrelevant issues
	for key, issue := range m.lastActiveIssues[queryName] {
		if _, ok := issues[key]; !ok {
			m.metricser.DisableMonitoring(queryName, issue)
		}
	}

	// Enable relevant issues
	for key, issue := range issues {
		if _, ok := m.lastActiveIssues[queryName][key]; !ok {
			m.metricser.EnableMonitoring(queryName, issue)
		}
	}

	m.lastActiveIssues[queryName] = issues
	return nil
}
