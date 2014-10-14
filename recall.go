package main

import (
	"database/sql"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/unrolled/render.v1"
)

// Store the db handle within the handler struct associated with the handler
// No need to pass around a global
type RecallHandler struct {
	db *sql.DB
}

type JSONRecallResult struct {
	Results []GensendgoRow `json:"results"`
}

// Create a New Recall Handler with the appropriate db handle
func NewRecallHandler(db *sql.DB) *RecallHandler {
	return &RecallHandler{db}
}

// business logic for the record
func (my *RecallHandler) updateRecallAccounting(aRow *GensendgoRow) (err error) {
	aRow.MaxReads--
	if aRow.MaxReads == 0 {
		return dbDeleteRowById(my.db, aRow.Id)
	} else {
		return dbUpdateMaxReadCount(my.db, aRow)
	}
	return
}

// HTTP Handler
func (my *RecallHandler) handler(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{})
	vars := mux.Vars(req)
	token := html.EscapeString(vars["token"])
	var results JSONRecallResult = JSONRecallResult{[]GensendgoRow{}}
	log.Println(token)
	if token != "" {

		parsedRows, err := dbFetchValidRowsById(my.db, token)
		if err != nil {
			errMsg := fmt.Sprintf("500 Internal Server Error: %v", err)
			log.Printf("%s", errMsg)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}

		if len(parsedRows) == 0 {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}
		err = my.updateRecallAccounting(&parsedRows[0])
		if err != nil {
			errMsg := fmt.Sprintf("500 Internal Server Error: %v", err)
			log.Printf("%s", errMsg)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return

		}

		if len(parsedRows) > 1 {
			log.Printf("WARN: Expected 1 row with this primary key, got %d", len(parsedRows))
		}
		results.Results = parsedRows
	}
	r.JSON(w, http.StatusOK, results)
}
