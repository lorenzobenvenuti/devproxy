package devproxy

import (
	"net/http"
	"net/url"
	"testing"
)

func getRequest(urlString string) *http.Request {
	url, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}
	r := &http.Request{}
	r.URL = url
	return r
}

func testUrlRequestMatcher(t *testing.T, urlString string, pattern string, expected bool) {
	matcher := NewUrlRequestMatcher(pattern)
	if matcher.Matches(getRequest(urlString)) != expected {
		t.Errorf("Url %s: %s doesn't match!", urlString, pattern)
	}
}

func TestUrlRequestMatcher(t *testing.T) {
	testUrlRequestMatcher(t, "http://www.google.it", "http://www", true)
	testUrlRequestMatcher(t, "http://www.google.it", "http://.*oog.*", true)
	testUrlRequestMatcher(t, "http://www.google.it", "http://.*\\.google\\..*", true)
	testUrlRequestMatcher(t, "http://www.google.it", "http://.*.it", true)
	testUrlRequestMatcher(t, "http://www.google.com", "http://.*.it", false)
}

func testStatusResponseMatcher(t *testing.T, actualCode int, matcherCode int, expected bool) {
	matcher := NewStatusResponseMatcher(matcherCode)
	r := &http.Response{}
	r.StatusCode = actualCode
	if matcher.Matches(r) != expected {
		t.Errorf("Status code %d: %d doesn't match!", actualCode, matcherCode)
	}
}

func TestStatusResponseMatcher(t *testing.T) {
	testStatusResponseMatcher(t, 404, 404, true)
	testStatusResponseMatcher(t, 404, 401, false)
}
