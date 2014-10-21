package main

import (
	"net/http"
)

type CORSMiddleware struct {
	Url string
}

func NewCORSMiddleware(url string) *CORSMiddleware {
	return &CORSMiddleware{url}
}

func (m *CORSMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	w.Header().Set("Access-Control-Allow-Origin", m.Url)
	next(w, req)
}
