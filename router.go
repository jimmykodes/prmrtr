package prmrtr

import (
	"context"
	"net/http"
	"path"
)

var _ http.Handler = (*Router)(nil)

type Router struct {
	parent *Router
	prefix string

	m map[string]http.Handler
	p []*pathEntry

	notFoundHandler http.Handler
}

func NewRouter(options ...Option) *Router {
	r := &Router{
		m:               make(map[string]http.Handler),
		notFoundHandler: http.HandlerFunc(defaultNotFoundHandler),
	}
	for _, option := range options {
		option.Apply(r)
	}
	return r
}

func (r *Router) SubRouter(prefix string) *Router {
	return NewRouter(withParent(r, prefix))
}

func (r *Router) Handle(route string, handler http.Handler) {
	if r.parent != nil {
		r.parent.Handle(path.Join(r.prefix, route), handler)
		return
	}
	r.m[route] = handler
	r.p = append(r.p, newPathEntry(route, checkParams))
}

func (r *Router) HandleFunc(route string, handler http.HandlerFunc) {
	r.Handle(route, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	this := newPathEntry(req.URL.Path, skipParams)
	for _, entry := range r.p {
		if entry.match(this) {
			ctx := context.WithValue(req.Context(), varCtxKey, newVarFromEntries(entry, this))
			r.m[entry.base].ServeHTTP(w, req.WithContext(ctx))
			return
		}
	}
	r.notFoundHandler.ServeHTTP(w, req)
}

var notFoundMessage = []byte("404 - Page Not Found")

func defaultNotFoundHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write(notFoundMessage)
}
