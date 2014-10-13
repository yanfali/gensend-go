package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/codegangsta/cli"
	_ "github.com/mattn/go-sqlite3"
)

var (
	dbFilename = "./gensendgo.db"
)

func openDb(fileName string) (db *sql.DB, err error) {
	if db, err = sql.Open("sqlite3", fileName); err != nil {
		return
	}
	return db, nil
}

func dbDumpTable(c *cli.Context) {
	db, err := dbOpenWrapper(dbDumpTableInner)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func dbDumpTableInner(db *sql.DB) (err error) {
	rows, err := db.Query("select id, maxreads, maxdays, createdTs, password from gensendgo")
	if err != nil {
		return
	}
	defer rows.Close()
	var rowCount int
	for rowCount = 0; rows.Next(); rowCount++ {
		var (
			id        string
			maxreads  int
			maxdays   int
			createdTs time.Time
			password  string
		)
		err = rows.Scan(&id, &maxreads, &maxdays, &createdTs, &password)
		if err != nil {
			return
		}
		log.Printf("%q %d %d %q %q\n", id, maxreads, maxdays, createdTs, password)
	}
	log.Printf("Found %d rows in table", rowCount)
	return
}

func dbSetupTestData(c *cli.Context) {
	db, err := dbOpenWrapper(dbTestDataInner)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func dbTransactionWrapper(db *sql.DB, fn func(tx *sql.Tx) error) (err error) {
	var tx *sql.Tx
	tx, err = db.Begin()
	if err != nil {
		return
	}
	err = fn(tx)
	if err != nil {
		return
	}
	tx.Commit()
	return
}

func dbTestDataInner(db *sql.DB) (err error) {
	dbTransactionWrapper(db, func(tx *sql.Tx) (err error) {
		stmt, err := tx.Prepare("insert into gensendgo(id, maxreads, maxdays, createdTs, password) values(?, ?, ?, ?, ?)")
		defer stmt.Close()
		result, err := stmt.Exec("abcdefgh", 5, 5, time.Now(), "12345678")
		if err != nil {
			return
		}
		log.Printf("insert success %v\n", result)
		return
	})
	dbDumpTableInner(db)
	return
}

func dbOpenWrapper(fn func(db *sql.DB) error) (db *sql.DB, err error) {
	db, err = openDb(dbFilename)
	if err != nil {
		return
	}
	err = fn(db)
	return
}

func dbInit(c *cli.Context) {
	log.Println("Initializing sqlite db")
	_ = os.Remove(dbFilename)
	db, err := dbOpenWrapper(dbInitInner)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func dbInitInner(db *sql.DB) (err error) {
	sqlStmt := `
	create table gensendgo (
		id string not null primary key,
		maxreads integer not null,
		maxdays integer not null,
		createdTs timestamp not null,
		password string not null
	);
	delete from gensendgo
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return errors.New(fmt.Sprintf("%q: %s\n", err, sqlStmt))
	}
	log.Println("created sqlitedb " + dbFilename)
	return
}
