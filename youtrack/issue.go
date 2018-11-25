package youtrack

import (
	"fmt"
	"github.com/krpn/youtrack-issues-prometheus-exporter/model"
)

type apiResponse []apiIssue

type apiIssue struct {
	Project struct {
		ShortName string `json:"shortName"`
	} `json:"project"`
	Summary         string `json:"summary"`
	NumberInProject int    `json:"numberInProject"`
}

func (ai apiIssue) ID() string {
	return fmt.Sprintf("%v-%v", ai.Project.ShortName, ai.NumberInProject)
}

func (ai apiIssue) ToIssue() model.Issue {
	return model.Issue{
		ID:    ai.ID(),
		Title: ai.Summary,
	}
}
