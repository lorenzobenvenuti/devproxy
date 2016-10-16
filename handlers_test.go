package devproxy

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
)

type testHandler struct {
	t           *testing.T
	expectedId  string
	expectedReq *http.Request
	newReq      *http.Request
}

func (th *testHandler) Handle(id string, r *http.Request, chain HandlerChain) *http.Response {
	if th.expectedId != id {
		th.t.Errorf("Request id: expected %s, found %s", th.expectedId, id)
	}
	if th.expectedReq != r {
		th.t.Errorf("Unexpected request")
	}
	return chain.Next(id, th.newReq)
}

type testRoundTripper struct {
	t           *testing.T
	expectedReq *http.Request
	resp        *http.Response
}

func (rt *testRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if rt.expectedReq != r {
		rt.t.Errorf("Unexpected request")
	}
	return rt.resp, nil
}

func TestHandlerChain(t *testing.T) {
	req1 := &http.Request{}
	req1.URL, _ = url.Parse("http://www.github.com")
	req2 := &http.Request{}
	req2.URL, _ = url.Parse("http://www.google.com")
	req3 := &http.Request{}
	req3.URL, _ = url.Parse("http://golang.org")
	resp := &http.Response{}
	h1 := &testHandler{t, "id", req1, req2}
	h2 := &testHandler{t, "id", req2, req3}
	rt := &testRoundTripper{t, req3, resp}
	logger := log.New(ioutil.Discard, "", log.Ldate)
	hc := NewHandlerChain(logger, rt, []Handler{h1, h2})
	if resp != hc.Next("id", req1) {
		t.Error("Unexpected response")
	}
}
