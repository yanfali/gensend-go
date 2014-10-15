package main

import (
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
	dao *gsgoDao
}

type JSONRecallResult struct {
	Results []GensendgoRow `json:"results"`
}

// Create a New Recall Handler with the appropriate db handle
func NewRecallHandler(dao *gsgoDao) *RecallHandler {
	return &RecallHandler{dao}
}

// business logic for the record
func (my *RecallHandler) updateRecallAccounting(aRow *GensendgoRow) (err error) {
	dao := my.dao
	aRow.MaxReads--
	if aRow.MaxReads == 0 {
		return dao.deleteById(aRow.Id)
	} else {
		return dao.updateMaxReadCount(aRow)
	}
	return
}

// HTTP Handler
func (my *RecallHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{})
	vars := mux.Vars(req)
	token := html.EscapeString(vars["token"])
	var results JSONRecallResult = JSONRecallResult{[]GensendgoRow{}}
	log.Println(token)
	if token != "" {

		parsedRows, err := my.dao.fetchValidRowById(token)
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
