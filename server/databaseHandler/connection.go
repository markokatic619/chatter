package databasehandler

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var dbInbox *sql.DB
var dbGroupChat *sql.DB

func init() {
	var err error

	db, err = sql.Open("mysql", "root:toor@tcp(localhost:3306)/testDb")
	if err != nil {
		panic(err)
	}
	dbInbox, err = sql.Open("mysql", "root:toor@tcp(localhost:3306)/inboxDb")
	if err != nil {
		panic(err)
	}
	dbGroupChat, err = sql.Open("mysql", "root:toor@tcp(localhost:3306)/groupChatDb")
	if err != nil {
		panic(err)
	}
}
