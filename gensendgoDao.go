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

type gsgoDao struct {
	db *sql.DB
}

func NewGsgoDao(db *sql.DB) *gsgoDao {
	return &gsgoDao{db}
}

func (my *gsgoDao) updateMaxReadCount(aRow *GensendgoRow) (err error) {
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

func (my *gsgoDao) deleteById(id string) (err error) {
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

func (my *gsgoDao) fetchValidRowsById(token string) (parsedRows []GensendgoRow, err error) {
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

func (my *gsgoDao) fetchValidRowById(token string) (parsedRows []GensendgoRow, err error) {
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
