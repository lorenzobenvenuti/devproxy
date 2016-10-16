package devproxy

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

type ProxyHandler struct {
	idGenerator         StringIdGenerator
	handlerChainFactory HandlerChainFactory
	url                 *url.URL
}

type HandlerChainFactory interface {
	CreateHandlerChain() HandlerChain
}

type handlerChainFactoryImpl struct {
	logger       *log.Logger
	roundTripper http.RoundTripper
	handlers     []Handler
}

func (hcf *handlerChainFactoryImpl) CreateHandlerChain() HandlerChain {
	return NewHandlerChain(hcf.logger, hcf.roundTripper, hcf.handlers)
}

func NewHandlerChainFactory(logger *log.Logger, roundTripper http.RoundTripper, handlers []Handler) HandlerChainFactory {
	hcf := &handlerChainFactoryImpl{}
	hcf.logger = logger
	hcf.handlers = handlers
	hcf.roundTripper = roundTripper
	return hcf
}

func (p *ProxyHandler) transformRequest(r *http.Request) *http.Request {
	// utilizzare header "host" ?
	r.URL.Scheme = p.url.Scheme
	r.URL.Host = p.url.Host
	return r
}

func (p *ProxyHandler) copyResponse(w http.ResponseWriter, resp *http.Response) {
	// copy response
	for k, v := range resp.Header {
		for i := 0; i < len(v); i++ {
			w.Header().Add(k, v[i])
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, err := io.Copy(w, resp.Body)
	if err != nil {
		panic(err)
	}
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := p.idGenerator.NewId()
	handlerChain := p.handlerChainFactory.CreateHandlerChain()
	p.copyResponse(w, handlerChain.Next(id, p.transformRequest(r)))
}

func NewProxyHandler(idGenerator StringIdGenerator, handlerChainFactory HandlerChainFactory, url *url.URL) http.Handler {
	ph := &ProxyHandler{url: url}
	ph.idGenerator = idGenerator
	ph.url = url
	ph.handlerChainFactory = handlerChainFactory
	return ph
}
