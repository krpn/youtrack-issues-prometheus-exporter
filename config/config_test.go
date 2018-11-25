package config

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()

	type testTableData struct {
		tcase          string
		raw            []byte
		expectedConfig *Config
		expectedErr    error
	}

	testTable := []testTableData{
		{
			tcase: "success",
			raw: []byte(`
{
  "endpoint": "http://www.test.com",
  "token": "abc",
  "queries": {
    "test": "test query"
  },
  "refresh_delay_seconds": 20,
  "request_timeout_seconds": 30,
  "listen_port": 9090
}`),
			expectedConfig: &Config{
				Endpoint:              "http://www.test.com",
				Token:                 "abc",
				Queries:               map[string]string{"test": "test query"},
				RefreshDelaySeconds:   20,
				RequestTimeoutSeconds: 30,
				ListenPort:            9090,
			},
			expectedErr: nil,
		},
		{
			tcase: "empty endpoint",
			raw: []byte(`
{
  "token": "abc",
  "queries": {
    "test": "test query"
  }
}`),
			expectedConfig: nil,
			expectedErr:    errors.New("empty endpoint"),
		},
		{
			tcase: "empty token",
			raw: []byte(`
{
  "endpoint": "http://www.test.com",
  "queries": {
    "test": "test query"
  }
}`),
			expectedConfig: nil,
			expectedErr:    errors.New("empty token"),
		},
		{
			tcase: "empty queries",
			raw: []byte(`
{
  "endpoint": "http://www.test.com",
  "token": "abc"
}`),
			expectedConfig: nil,
			expectedErr:    errors.New("empty queries"),
		},
		{
			tcase: "fix default values",
			raw: []byte(`
{
  "endpoint": "http://www.test.com",
  "token": "abc",
  "queries": {
    "test": "test query"
  }
}`),
			expectedConfig: &Config{
				Endpoint:              "http://www.test.com",
				Token:                 "abc",
				Queries:               map[string]string{"test": "test query"},
				RefreshDelaySeconds:   10,
				RequestTimeoutSeconds: 10,
				ListenPort:            8080,
			},
			expectedErr: nil,
		},
		{
			tcase:          "invalid json",
			raw:            []byte(``),
			expectedConfig: nil,
			expectedErr:    json.Unmarshal([]byte(``), nil),
		},
	}

	for _, testUnit := range testTable {
		config, err := New(testUnit.raw)
		assert.Equal(t, testUnit.expectedConfig, config, testUnit.tcase)
		assert.Equal(t, testUnit.expectedErr, err, testUnit.tcase)
	}
}
