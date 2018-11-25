package youtrack

import (
	"encoding/json"
	"fmt"
	"github.com/krpn/youtrack-issues-prometheus-exporter/model"
	"net/url"
)

//go:generate mockgen -source=youtrack.go -destination=youtrack_mocks.go -package=youtrack doc github.com/golang/mock/gomock

const apiPath = "api/issues"

type makeRequester interface {
	MakeRequest(url string, headers map[string]string) ([]byte, error)
}

// YouTrack describes simple YouTrack API client.
type YouTrack struct {
	requester makeRequester
	url       url.URL
	getParams url.Values
	headers   map[string]string
}

// New creates YouTrack instance.
func New(endpoint, token string, requester makeRequester) (*YouTrack, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	u.Path = apiPath

	getParams := url.Values{}
	getParams.Set("fields", "project(shortName),numberInProject,summary")

	return &YouTrack{
		requester: requester,
		url:       *u,
		getParams: getParams,
		headers: map[string]string{
			"Accept":        "application/json",
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %v", token),
		},
	}, nil
}

// GetIssues gets issues for passed query string.
func (yt *YouTrack) GetIssues(query string) (issues map[string]model.Issue, err error) {
	u := yt.getAPIURL(query)

	body, err := yt.requester.MakeRequest(u, yt.headers)
	if err != nil {
		return nil, err
	}

	response := make(apiResponse, 0)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	issues = make(map[string]model.Issue, len(response))
	for _, ai := range response {
		issue := ai.ToIssue()
		issues[issue.FullID()] = issue
	}

	return issues, nil
}

func (yt *YouTrack) getAPIURL(query string) string {
	params := yt.getParams
	params.Set("query", query)

	u := yt.url
	u.RawQuery = params.Encode()

	return u.String()
}
