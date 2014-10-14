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

func dbOpen(fileName string) (db *sql.DB, err error) {
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

type GensendgoRow struct {
	Id        string    `json:"token"`
	MaxReads  int       `json:"maxReads"`
	MaxDays   int       `json:"maxDays"`
	CreatedTs time.Time `json:"createdTs"`
	Password  string    `json:"password"`
}

func (my *GensendgoRow) String() string {
	return fmt.Sprintf("%q %d %d %q %q", my.Id, my.MaxReads, my.MaxDays, my.CreatedTs, my.Password)
}

func dbDumpTableInner(db *sql.DB) (err error) {
	rows, err := db.Query("select id, maxreads, maxdays, createdTs, password from gensendgo")
	if err != nil {
		return
	}
	defer rows.Close()
	var rowCount int
	for rowCount = 0; rows.Next(); rowCount++ {
		var aRow GensendgoRow
		err = rows.Scan(&aRow.Id, &aRow.MaxReads, &aRow.MaxDays, &aRow.CreatedTs, &aRow.Password)
		if err != nil {
			return
		}
		log.Printf("%s\n", aRow.String())
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
	db, err = dbOpen(dbFilename)
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
