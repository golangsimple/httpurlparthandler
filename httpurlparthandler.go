package httpurlparthandler

import (
	"net/http"
	"strings"
)

type HandlerFunc struct {
	Route string
	HandlerFunc func(part string) http.HandlerFunc
}

func NewHandlerFunc(parent string, pattern string, handlerFunc func(part string) http.HandlerFunc) *HandlerFunc {
	return &HandlerFunc{Route: parent + pattern, HandlerFunc:handlerFunc}
}

func (handler *HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	part := getNextPart(r.URL.Path, handler.Route)

	handler.HandlerFunc(part)(w,r)
}


type Handler struct {
	Route string
	Handler func(part string) http.Handler
}

func NewHandler(parent string, pattern string, handler func(part string) http.Handler) *Handler {
	return &Handler{Route: parent + pattern, Handler:handler}
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	part := getNextPart(r.URL.Path, handler.Route)

	handler.Handler(part).ServeHTTP(w,r)
}

func getNextPart(path, route string) string {
	part := strings.Replace(path, route, "", 1)
	parts := strings.Split(part, "/")
	if len(parts) > 0 {
		part = parts[0]
	}
	return part
}
