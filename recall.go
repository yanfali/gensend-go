package main

import (
	"errors"
	"html"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/unrolled/render.v1"
)

// Store the db handle within the handler struct associated with the handler
// No need to pass around a global
type RecallHandler struct {
	dao AccountingDao
}

type JSONRecallResult struct {
	Results []GensendgoRow `json:"results"`
}

// Create a New Recall Handler with the appropriate db handle
func NewRecallHandler(dao AccountingDao) *RecallHandler {
	return &RecallHandler{dao}
}

// business logic for the record
func (my *RecallHandler) updateRecallAccounting(aRow *GensendgoRow) (err error) {
	dao := my.dao
	aRow.MaxReads--
	if aRow.MaxReads == 0 {
		return dao.DeleteById(aRow.Id)
	} else {
		return dao.UpdateMaxReadCount(aRow)
	}
	return
}

// HTTP Handler
func (my *RecallHandler) handler(w http.ResponseWriter, req *http.Request) (int, error) {
	vars := mux.Vars(req)
	token := html.EscapeString(vars["token"])
	log.Println(token)
	var results JSONRecallResult = JSONRecallResult{[]GensendgoRow{}}
	if token != "" {
		parsedRows, err := my.dao.FetchValidRowById(token)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		if len(parsedRows) == 0 {
			return http.StatusNotFound, errors.New("Token Not Found")
		}
		if err = my.updateRecallAccounting(&parsedRows[0]); err != nil {
			return http.StatusInternalServerError, err

		}

		if len(parsedRows) > 1 {
			log.Printf("WARN: Expected 1 row with this primary key, got %d", len(parsedRows))
		}
		results.Results = parsedRows
	}
	r := render.New(render.Options{})
	r.JSON(w, http.StatusOK, results)
	return http.StatusOK, nil
}
