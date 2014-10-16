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

const (
	GENSENDGO_INSERT_ROW   = "INSERT INTO gensendgo(id, maxreads, maxminutes, createdTs, expiredTs, password) VALUES(?, ?, ?, ?, ?, ?)"
	GENSENDGO_CREATE_TABLE = `
	create table gensendgo (
		id string not null primary key,
		maxreads integer not null,
		maxminutes integer not null,
		createdTs timestamp not null,
		expiredTs timestamp not null,
		password string not null
	);
	delete from gensendgo;
	`
)

// wrapper to open sql db
func dbOpen(fileName string) (db *sql.DB, err error) {
	if db, err = sql.Open("sqlite3", fileName); err != nil {
		return
	}
	if err = db.Ping(); err != nil {
		return
	}
	return db, nil
}

// helper to dump the contents of the db
func dbDumpTable(c *cli.Context) {
	db, err := dbOpenWrapper(dbDumpTableInner)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

// helper to be invoked within a dbOpenWrapper
func dbDumpTableInner(db *sql.DB) (err error) {
	var rows *sql.Rows
	rows, err = db.Query(GENSENDGO_SELECT_ROW)
	if err != nil {
		return
	}
	defer rows.Close()
	var rowCount int
	for rowCount = 0; rows.Next(); rowCount++ {
		var aRow GensendgoRow
		err = aRow.ScanRows(rows)
		if err != nil {
			return
		}
		log.Printf("%s\n", aRow.String())
	}
	log.Printf("Found %d rows in table", rowCount)
	return
}

// helper to set up a basic db with a test record
func dbSetupTestData(c *cli.Context) {
	db, err := dbOpenWrapper(dbTestDataInner)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

// helper to wrap boilerplate for transactions
func dbTransactionWrapper(db *sql.DB, fn func(tx *sql.Tx) error) (err error) {
	if err = db.Ping(); err != nil {
		return
	}
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

// helper to add a test record to a db
func dbTestDataInner(db *sql.DB) (err error) {
	err = dbTransactionWrapper(db, func(tx *sql.Tx) (err error) {
		// TODO might be possible to prepare these only once
		var stmt *sql.Stmt
		stmt, err = tx.Prepare(GENSENDGO_INSERT_ROW)
		defer stmt.Close()
		var result sql.Result
		expiresMinutes := 1 //TODO Make this configurable via CLI
		expiresTs := time.Now().Add(time.Minute * time.Duration(expiresMinutes)).UTC()
		result, err = stmt.Exec("abcdefgh", 5, expiresMinutes, time.Now().UTC(), expiresTs, "12345678")
		if err != nil {
			return
		}
		log.Printf("insert success %v\n", result)
		return
	})
	if err != nil {
		return
	}
	dbDumpTableInner(db)
	return
}

// wrap the boilerplate for opening a sqlite DB
func dbOpenWrapper(fn func(db *sql.DB) error) (db *sql.DB, err error) {
	db, err = dbOpen(dbFilename)
	if err != nil {
		return
	}
	err = fn(db)
	return
}

// helper to initialize the database
// removes the old one if found
func dbInit(c *cli.Context) {
	log.Println("Initializing sqlite db")
	_ = os.Remove(dbFilename)
	db, err := dbOpenWrapper(dbInitInner)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

// helper to create the initial database
func dbInitInner(db *sql.DB) (err error) {
	_, err = db.Exec(GENSENDGO_CREATE_TABLE)
	if err != nil {
		return errors.New(fmt.Sprintf("%q: %s\n", err, GENSENDGO_CREATE_TABLE))
	}
	log.Println("created sqlitedb " + dbFilename)
	return
}
