package main

import (
	"fmt"
	"net/http"

	"gopkg.in/unrolled/render.v1"
)

type MyHttpHandler interface {
	handler(http.ResponseWriter, *http.Request) (int, error)
}

type appHandler struct {
	MyHttpHandler
}

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := ah.handler(w, r)
	if err != nil {
		r := render.New(render.Options{})
		r.JSON(w, status, JSONErrorResponse{status, fmt.Sprintf("%v", err)})
		return
	}
}
