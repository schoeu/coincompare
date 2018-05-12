package utils

import (
	"database/sql"
	"log"
	"regexp"
	"strings"
	"time"
)

const sqlReg = "(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\\b(select|update|and|or|delete|insert|trancate|char|into|substr|ascii|declare|exec|count|master|into|drop|execute)\\b)"

// error handler.
func ErrHandle(err error) {
	if err != nil {
		log.Println(err)
	}
}

// db opener
func OpenDb(dbType string, dbStr string) *sql.DB {
	if dbType == "" {
		log.Println("No dbType")
	}
	db, err := sql.Open(dbType, dbStr)
	ErrHandle(err)

	err = db.Ping()
	ErrHandle(err)
	return db
}

func GetCurrentDate() string {
	t := time.Now().String()
	return strings.Split(t, " ")[0]
}

// check sql string
func CheckSql(s string) string {
	match, _ := regexp.Match(sqlReg, []byte(s))
	if match {
		return ""
	}
	return s
}
