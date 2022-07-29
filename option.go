package prmrtr

import (
	"net/http"
)

type Option interface {
	Apply(*Router)
}

type OptionFunc func(*Router)

func (f OptionFunc) Apply(r *Router) {
	f(r)
}

func withParent(r *Router, prefix string) OptionFunc {
	return func(router *Router) {
		router.parent = r
		router.prefix = prefix
	}
}

func NotFoundHandlerOption(h http.Handler) OptionFunc {
	return func(router *Router) {
		router.notFoundHandler = h
	}
}
