package main

import (
	"database/sql"
	"fmt"
	"time"
)

type GensendgoRow struct {
	Id         string    `json:"token"`
	MaxReads   int       `json:"maxReads"`
	MaxMinutes int       `json:"maxMinutes"`
	CreatedTs  time.Time `json:"createdTs"`
	ExpiredTs  time.Time `json:"expiredTs"`
	Password   string    `json:"password"`
}

func (my *GensendgoRow) String() string {
	return fmt.Sprintf("%q %d %d %q %q %q", my.Id, my.MaxReads, my.MaxMinutes, my.CreatedTs, my.ExpiredTs, my.Password)
}

func (my *GensendgoRow) Scan(rows *sql.Rows) (err error) {
	return rows.Scan(&my.Id, &my.MaxReads, &my.MaxMinutes, &my.CreatedTs, &my.ExpiredTs, &my.Password)
}
