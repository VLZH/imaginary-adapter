package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createImageRequest(url string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	return req
}

func TestParseRequest(t *testing.T) {
	testCases := []struct {
		name    string
		request *http.Request
		host    string
		prefix  string
		expect  *ImaginaryParameters
	}{
		{
			name:    "Default",
			request: createImageRequest("http://server.com/uploads/image.png?method=fit&width=300&height=300&quality=70"),
			host:    "http://imaginary:9000",
			prefix:  "/uploads",
			expect: &ImaginaryParameters{
				Host:    "http://imaginary:9000",
				File:    "/image.png",
				Height:  300,
				Width:   300,
				Method:  "fit",
				Quality: 70,
			},
		},
		{
			name:    "Default quality",
			request: createImageRequest("http://server.com/uploads/image.png?method=fit&width=300&height=300"),
			host:    "http://imaginary:9000",
			prefix:  "/uploads",
			expect: &ImaginaryParameters{
				Host:    "http://imaginary:9000",
				File:    "/image.png",
				Height:  300,
				Width:   300,
				Method:  "fit",
				Quality: 95,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, _ := parseRequest(tc.request, tc.host, tc.prefix)
			t.Run("check parameters", func(t *testing.T) {
				assert.Equal(t, result, tc.expect)
			})
			// TODO: check GetUrl method of ImaginaryParameters
		})
	}
}
