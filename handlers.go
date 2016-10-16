package devproxy

import (
	"log"
	"net/http"
	"net/http/httputil"
)

type Handler interface {
	Handle(id string, r *http.Request, chain HandlerChain) *http.Response
}

type HandlerChain interface {
	Next(id string, r *http.Request) *http.Response
}

type handlerChainImpl struct {
	logger       *log.Logger
	roundTripper http.RoundTripper
	index        int
	handlers     []Handler
}

func (hc *handlerChainImpl) Next(id string, r *http.Request) *http.Response {
	hc.index++
	if hc.index == len(hc.handlers) {
		hc.logger.Printf("Forwarding request %s to %s", id, r.URL.String())
		response, err := hc.roundTripper.RoundTrip(r)
		if err != nil {
			hc.logger.Panicf("Error performing request to %s: %s", r.URL.String(), err.Error())
			panic(err)
		}
		return response
	}
	return hc.handlers[hc.index].Handle(id, r, hc)
}

func NewHandlerChain(logger *log.Logger, roundTripper http.RoundTripper, handlers []Handler) *handlerChainImpl {
	hc := &handlerChainImpl{}
	hc.logger = logger
	hc.roundTripper = roundTripper
	hc.handlers = handlers
	hc.index = -1
	return hc
}

type loggerHandler struct {
	logger *log.Logger
}

func (lh *loggerHandler) Handle(id string, r *http.Request, chain HandlerChain) *http.Response {
	lh.logger.Printf("Request %s: %s", id, r.URL.String())
	response := chain.Next(id, r)
	bytes, err := httputil.DumpResponse(response, true)
	if err != nil {
		panic(err)
	}
	lh.logger.Printf("Response to %s: %s", id, string(bytes))
	return response
}

func NewLoggerHandler(logger *log.Logger) Handler {
	return &loggerHandler{logger: logger}
}

type empty struct{}

type breakpoint struct {
	semaphore chan empty
}

type debugHandler struct {
	idGenerator           StringIdGenerator
	eventBus              EventBus
	beforeRequestMatchers map[string]RequestMatcher
	afterRequestMatchers  map[string]RequestMatcher
	afterResponseMatchers map[string]ResponseMatcher
	activeBreakpoints     map[string]breakpoint
}

func NewDebugHandler(idGenerator StringIdGenerator, eb EventBus) Handler {
	dh := &debugHandler{}
	dh.idGenerator = idGenerator
	dh.eventBus = eb
	dh.beforeRequestMatchers = make(map[string]RequestMatcher)
	dh.afterRequestMatchers = make(map[string]RequestMatcher)
	dh.afterResponseMatchers = make(map[string]ResponseMatcher)
	dh.activeBreakpoints = make(map[string]breakpoint)
	eb.Subscribe(ResumeBreakpointTopic, dh.resume)
	eb.Subscribe(AddBeforeRequestMatcherTopic, dh.addBeforeRequestMatcher)
	eb.Subscribe(AddAfterRequestMatcherTopic, dh.addAfterRequestMatcher)
	eb.Subscribe(AddAfterResponseMatcherTopic, dh.addAfterResponseMatcher)
	eb.Subscribe(RemoveBeforeRequestMatcherTopic, dh.removeBeforeRequestMatcher)
	eb.Subscribe(RemoveAfterRequestMatcherTopic, dh.removeAfterRequestMatcher)
	eb.Subscribe(RemoveAfterResponseMatcherTopic, dh.removeAfterResponseMatcher)
	return dh
}

func newBreakpoint() breakpoint {
	bp := breakpoint{}
	bp.semaphore = make(chan empty)
	return bp
}

func (dh *debugHandler) addBeforeRequestMatcher(message interface{}) {
	dh.beforeRequestMatchers[dh.idGenerator.NewId()] = message.(RequestMatcher)
}

func (dh *debugHandler) addAfterRequestMatcher(message interface{}) {
	dh.afterRequestMatchers[dh.idGenerator.NewId()] = message.(RequestMatcher)
}

func (dh *debugHandler) removeBeforeRequestMatcher(message interface{}) {
	id := message.(string)
	if _, ok := dh.beforeRequestMatchers[id]; ok {
		delete(dh.beforeRequestMatchers, id)
	}
}

func (dh *debugHandler) removeAfterRequestMatcher(message interface{}) {
	id := message.(string)
	if _, ok := dh.afterRequestMatchers[id]; ok {
		delete(dh.afterRequestMatchers, id)
	}
}

func (dh *debugHandler) addAfterResponseMatcher(message interface{}) {
	dh.afterResponseMatchers[dh.idGenerator.NewId()] = message.(ResponseMatcher)
}

func (dh *debugHandler) removeAfterResponseMatcher(message interface{}) {
	id := message.(string)
	if _, ok := dh.afterResponseMatchers[id]; ok {
		delete(dh.afterResponseMatchers, id)
	}
}

func (dh *debugHandler) resume(message interface{}) {
	uid := message.(string)
	if breakpoint, ok := dh.activeBreakpoints[uid]; ok {
		delete(dh.activeBreakpoints, uid)
		// unlock
		<-breakpoint.semaphore
		dh.eventBus.Dispatch(BreakpointResumedTopic, uid)
	}
}

func (dh *debugHandler) hit() {
	id := dh.idGenerator.NewId()
	// dispatch event
	bp := newBreakpoint()
	dh.activeBreakpoints[id] = bp
	dh.eventBus.Dispatch(BreakpointHitTopic, id)
	// lock
	bp.semaphore <- empty{}
}

func (dh *debugHandler) Handle(id string, r *http.Request, chain HandlerChain) *http.Response {
	for _, requestMatcher := range dh.beforeRequestMatchers {
		if requestMatcher.Matches(r) {
			dh.hit()
		}
	}
	resp := chain.Next(id, r)
	return resp
}

// TODO: timeout handler
// TODO: edit request (body, header)
// TODO: edit response (body, header)
// TODO: throttle
