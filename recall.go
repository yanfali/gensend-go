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

type JSONRecallResult struct {
	Results []GensendgoRow `json:"results"`
}

func recallHandler(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{})
	vars := mux.Vars(req)
	token := html.EscapeString(vars["token"])
	var results JSONRecallResult = JSONRecallResult{[]GensendgoRow{}}
	log.Println(token)
	if token != "" {
		rows, err := db.Query("select * from gensendgo where id=?", token)
		if err != nil {
			http.Error(w, fmt.Sprintf("500 Internal Server Error: %v", err), http.StatusInternalServerError)
			return
		}
		parsedRows := []GensendgoRow{}
		for rows.Next() {
			aRow := GensendgoRow{}
			err = rows.Scan(&aRow.Id, &aRow.MaxReads, &aRow.MaxDays, &aRow.CreatedTs, &aRow.Password)
			parsedRows = append(parsedRows, aRow)

		}
		if len(parsedRows) == 0 {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}

		err = dbTransactionWrapper(db, func(tx *sql.Tx) (err error) {
			aRow := &parsedRows[0]
			if aRow.MaxReads == 1 {
				var updateRow *sql.Stmt
				updateRow, err = db.Prepare("DELETE from gensendgo WHERE id=?")
				if err != nil {
					return
				}
				var tx *sql.Tx
				tx, err = db.Begin()
				_, err = tx.Stmt(updateRow).Exec(aRow.Id)
				if err != nil {
					return
				}
				tx.Commit()
			} else {
				var updateRow *sql.Stmt
				updateRow, err = db.Prepare("UPDATE gensendgo SET maxreads=? WHERE id=?")
				if err != nil {
					return
				}
				var tx *sql.Tx
				tx, err = db.Begin()
				_, err = tx.Stmt(updateRow).Exec(aRow.MaxReads-1, aRow.Id)
				if err != nil {
					return
				}
				tx.Commit()
			}
			return
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("500 Internal Server Error: %v", err), http.StatusInternalServerError)
			return

		}

		for _, row := range parsedRows {
			results.Results = append(results.Results, row)
		}
	}
	r.JSON(w, http.StatusOK, results)
}
