package mantis

import (
	"net/http"
	"regexp"
)

type Route interface {
	Match(pattern string) bool
	Pattern() string
}

type route struct {
	pattern string
	regex   *regexp.Regexp
	Handle  Handler
}

func (r *route) Match(pattern string) bool {
	return false
}

type Router interface {
	Handle(rw http.ResponseWriter, r *http.Request)
	NotFound()
}

type router struct {
	routes []Route
}

func (rt *router) Handle(rw http.ResponseWriter, r *http.Request) {
	for _, route := range rt.routes {
		if r.Match(r.URL.Path) {
			return route
		}
	}
}

func (rt *router) NotFound() {

}
