package httpwrap

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Parallel()

	type testTableData struct {
		client   *http.Client
		expected *ClientWrap
	}

	testTable := []testTableData{
		{
			client:   &http.Client{Timeout: time.Second},
			expected: &ClientWrap{c: &http.Client{Timeout: time.Second}},
		},
	}

	for _, testUnit := range testTable {
		assert.Equal(t, testUnit.expected, New(testUnit.client))
	}
}

func TestClientWrap_MakeRequest(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	doerMock := NewMockdoer(ctrl)
	clientWrap := ClientWrap{c: doerMock}

	type testTableData struct {
		tcase        string
		url          string
		headers      map[string]string
		expectFunc   func(d *Mockdoer)
		expectedBody []byte
		expectedErr  error
	}

	testTable := []testTableData{
		{
			tcase:   "success request",
			url:     "http://www.test.com/",
			headers: map[string]string{"Authorization": "123"},
			expectFunc: func(d *Mockdoer) {
				req, _ := http.NewRequest("GET", "http://www.test.com/", nil)
				req.Header.Set("Authorization", "123")
				d.EXPECT().Do(req).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString("resp body")),
				}, nil)
			},
			expectedBody: []byte("resp body"),
			expectedErr:  nil,
		},
		{
			tcase:        "bad request url",
			url:          "http://www test com/",
			headers:      map[string]string{"Authorization": "123"},
			expectFunc:   func(d *Mockdoer) {},
			expectedBody: nil,
			expectedErr:  &url.Error{Op: "parse", URL: "http://www test com/", Err: url.InvalidHostError(" ")},
		},
		{
			tcase:   "request error",
			url:     "http://www.test.com/",
			headers: nil,
			expectFunc: func(d *Mockdoer) {
				req, _ := http.NewRequest("GET", "http://www.test.com/", nil)
				d.EXPECT().Do(req).Return(nil, errors.New("request error"))
			},
			expectedBody: nil,
			expectedErr:  errors.New("request error"),
		},
		{
			tcase:   "bad status code",
			url:     "http://www.test.com/",
			headers: nil,
			expectFunc: func(d *Mockdoer) {
				req, _ := http.NewRequest("GET", "http://www.test.com/", nil)
				d.EXPECT().Do(req).Return(&http.Response{
					StatusCode: http.StatusBadGateway,
					Body:       ioutil.NopCloser(bytes.NewBufferString("resp body")),
				}, nil)
			},
			expectedBody: nil,
			expectedErr:  errors.New("returned HTTP status: 502, body close error: <nil>"),
		},
		{
			tcase:   "body read error",
			url:     "http://www.test.com/",
			headers: nil,
			expectFunc: func(d *Mockdoer) {
				req, _ := http.NewRequest("GET", "http://www.test.com/", nil)
				d.EXPECT().Do(req).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(errorReader{}),
				}, nil)
			},
			expectedBody: nil,
			expectedErr:  errors.New("body read error: read error, body close error: <nil>"),
		},
	}

	for _, testUnit := range testTable {
		testUnit.expectFunc(doerMock)
		body, err := clientWrap.MakeRequest(testUnit.url, testUnit.headers)
		assert.Equal(t, testUnit.expectedBody, body, testUnit.tcase)
		assert.Equal(t, testUnit.expectedErr, err, testUnit.tcase)
	}
}

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}
