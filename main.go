package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var (
	app           *cli.App
	recallHandler *RecallHandler
	storeHandler  *StoreHandler
)

func getBaseAPIUrl() string {
	return "/api/v1"
}

func init() {
	app = cli.NewApp()
	app.Name = "gensend-go"
	app.Commands = []cli.Command{
		{
			Name:      "database",
			ShortName: "db",
			Usage:     "database commands",
			Subcommands: []cli.Command{
				{
					Name:      "initialize",
					ShortName: "init",
					Usage:     "initialize a new sqlite db",
					Action:    dbInit,
				},
				{
					Name:      "testdata",
					ShortName: "test",
					Usage:     "initialize sqlite db with test data",
					Action:    dbSetupTestData,
				},
			},
		},
	}
}

func webServer(c *cli.Context) {
	mux := mux.NewRouter()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})
	mux.Handle(getBaseAPIUrl()+"/store", storeHandler).Methods("POST")
	mux.Handle(getBaseAPIUrl()+"/recall/{token}", recallHandler).Methods("GET")
	mux.HandleFunc(getBaseAPIUrl()+"/sweep", func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	}).Methods("PUT")

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3000")
}

func main() {
	app.Action = webServer
	var err error
	var db *sql.DB
	db, err = dbOpen(dbFilename)
	// Inject dependencies into Recall Handler
	gsgoDao := NewGsgoDao(db)
	recallHandler = NewRecallHandler(gsgoDao)
	storeHandler = NewStoreHandler(gsgoDao)

	if err != nil {
		// TODO init db if it doesn't exist
		log.Fatal(err)
	}
	defer db.Close()
	app.Run(os.Args)
}
