package devproxy

import (
	"net/http"

	"github.com/gorilla/mux"
)

type apiHandler struct {
	eventBus EventBus
	handler  http.Handler
}

func (a *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(w, r)
}

func (a *apiHandler) getBreakpoints(w http.ResponseWriter, r *http.Request) {

}

func (a *apiHandler) addBreakpoint(w http.ResponseWriter, r *http.Request) {

}

func (a *apiHandler) deleteBreakpoint(w http.ResponseWriter, r *http.Request) {

}

func NewApiHandler(eventBus EventBus) http.Handler {
	a := &apiHandler{}
	a.eventBus = eventBus
	r := mux.NewRouter()
	r.HandleFunc("/api/breakpoints", a.getBreakpoints).Methods("GET")
	r.HandleFunc("/api/breakpoints", a.addBreakpoint).Methods("POST")
	r.HandleFunc("/api/breakpoints/{id}", a.deleteBreakpoint).Methods("DELETE")
	a.handler = r
	return a
}
