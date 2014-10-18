package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	GENSENDGO_SELECT_ROW           = "SELECT id, maxreads, maxminutes, createdTs, expiredTs, password FROM gensendgo"
	GENSENDGO_DELETE_ROW           = "DELETE from gensendgo WHERE id=?"
	GENSENDGO_UPDATE_ROW_MAX_READS = "UPDATE gensendgo SET maxreads=? WHERE id=?"
	GENSENDGO_SELECT_ROW_WHERE_ID  = "SELECT * FROM gensendgo where id=? AND expiredTs > ?"
)

type InsertDao interface {
	InsertToken(token, password string, maxReads, maxMinutes int) error
}

type AccountingDao interface {
	UpdateMaxReadCount(aRow *GensendgoRow) error
	DeleteById(id string) error
	FetchValidRowsById(token string) (parsedRows []GensendgoRow, err error)
	FetchValidRowById(token string) (parsedRows []GensendgoRow, err error)
}

type GsgoDao struct {
	db *sql.DB
}

func NewGsgoDao(db *sql.DB) *GsgoDao {
	return &GsgoDao{db}
}

func (my *GsgoDao) UpdateMaxReadCount(aRow *GensendgoRow) (err error) {
	return dbTransactionWrapper(my.db, func(tx *sql.Tx) (err error) {
		var updateRow *sql.Stmt
		log.Printf("updating row id %q read count", aRow.Id)
		// TODO might be possible to prepare these only once
		updateRow, err = my.db.Prepare(GENSENDGO_UPDATE_ROW_MAX_READS)
		if err != nil {
			return
		}
		_, err = tx.Stmt(updateRow).Exec(aRow.MaxReads, aRow.Id)
		return
	})
}

func (my *GsgoDao) DeleteById(id string) (err error) {
	return dbTransactionWrapper(my.db, func(tx *sql.Tx) (err error) {
		var updateRow *sql.Stmt
		log.Printf("deleting row id %q expired", id)
		// TODO might be possible to prepare these only once
		updateRow, err = my.db.Prepare(GENSENDGO_DELETE_ROW)
		if err != nil {
			return
		}
		_, err = tx.Stmt(updateRow).Exec(id)
		return
	})
}

func (my *GsgoDao) FetchValidRowsById(token string) (parsedRows []GensendgoRow, err error) {
	var rows *sql.Rows
	rows, err = my.db.Query(GENSENDGO_SELECT_ROW_WHERE_ID, token, time.Now().UTC())
	if err != nil {
		return
	}
	parsedRows = []GensendgoRow{}
	defer rows.Close()
	for rows.Next() {
		aRow := GensendgoRow{}
		err = aRow.ScanRows(rows)
		if err != nil {
			return nil, err
		}
		log.Printf("%v", aRow)
		parsedRows = append(parsedRows, aRow)

	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return
}

func (my *GsgoDao) FetchValidRowById(token string) (parsedRows []GensendgoRow, err error) {
	var aRow GensendgoRow
	err = aRow.ScanRow(my.db.QueryRow(GENSENDGO_SELECT_ROW_WHERE_ID, token, time.Now().UTC()))
	if err == sql.ErrNoRows {
		return []GensendgoRow{}, nil
	}
	if err != nil {
		return nil, err
	}
	return []GensendgoRow{aRow}, nil
}

func (my *GsgoDao) InsertToken(token, password string, maxReads, maxMinutes int) (err error) {
	return dbTransactionWrapper(my.db, func(tx *sql.Tx) (err error) {
		var stmt *sql.Stmt
		stmt, err = tx.Prepare(GENSENDGO_INSERT_ROW)
		if err != nil {
			return
		}
		defer stmt.Close()
		var result sql.Result
		createdTs := time.Now().UTC()
		expiresTs := createdTs.Add(time.Minute * time.Duration(maxMinutes))
		log.Printf("%v %v %#v", createdTs, expiresTs, time.Duration(maxMinutes))
		result, err = stmt.Exec(token, maxReads, maxMinutes, createdTs, expiresTs, password)
		if err != nil {
			log.Printf("insert failed %v\n", err)
			return
		}
		log.Printf("insert success %v\n", result)
		return
	})
}
