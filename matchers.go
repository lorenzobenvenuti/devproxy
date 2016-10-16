package devproxy

import (
	"net/http"
	"regexp"
)

type RequestMatcher interface {
	Matches(r *http.Request) bool
}

type ResponseMatcher interface {
	Matches(r *http.Response) bool
}

type urlRequestMatcher struct {
	re *regexp.Regexp
}

func (urm *urlRequestMatcher) Matches(r *http.Request) bool {
	return urm.re.MatchString(r.URL.String())
}

func NewUrlRequestMatcher(expr string) RequestMatcher {
	re, err := regexp.Compile(expr)
	if err != nil {
		panic(err)
	}
	return &urlRequestMatcher{re: re}
}

type statusResponseMatcher struct {
	statusCode int
}

func (srm *statusResponseMatcher) Matches(r *http.Response) bool {
	return r.StatusCode == srm.statusCode
}

func NewStatusResponseMatcher(statusCode int) ResponseMatcher {
	return &statusResponseMatcher{statusCode: statusCode}
}
