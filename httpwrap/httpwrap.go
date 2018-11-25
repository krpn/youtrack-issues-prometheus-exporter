package httpwrap

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//go:generate mockgen -source=httpwrap.go -destination=httpwrap_mocks.go -package=httpwrap doc github.com/golang/mock/gomock

type doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// ClientWrap executes HTTP requests.
type ClientWrap struct {
	c doer
}

// New creates wrapper for http.Client for handy making requests.
func New(client *http.Client) *ClientWrap {
	return &ClientWrap{c: client}
}

// MakeRequest making request for passed parameters.
func (c *ClientWrap) MakeRequest(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("returned HTTP status: %v, body close error: %v", resp.StatusCode, resp.Body.Close())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("body read error: %v, body close error: %v", err.Error(), resp.Body.Close())
	}

	return body, resp.Body.Close()
}
