package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})
	mux.HandleFunc("/store", func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	})
	mux.HandleFunc("/recall", func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	})
	mux.HandleFunc("/sweep", func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	})

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3000")
}
