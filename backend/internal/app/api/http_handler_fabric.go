package api

import "net/http"

type HTTPHandlerFabric interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type fabric struct {
	targetFunc func(w http.ResponseWriter, r *http.Request)
}

func (h *fabric) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.targetFunc(w, r)
}

func NewHTTPHandler(targetFunc func(w http.ResponseWriter, r *http.Request)) HTTPHandlerFabric {
	return &fabric{targetFunc: targetFunc}
}
