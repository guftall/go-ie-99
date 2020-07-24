package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB

func initializeDatabase() {
	db, err := sql.Open("mysql", "omid:65254585Om@tcp(192.168.8.100)/ieproj")
	if err != nil {
		log.Print("connect to mysql failed", err)
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	database = db
}

func countRecords() int {
	sql := "select count(id) as count from tblkeys"

	type count struct {
		count int
	}
	var c count
	rows, err := database.Query(sql)

	if err != nil {
		log.Print("query failed", err, sql)
		return -1
	}

	for rows.Next() {
		err := rows.Scan(&c.count)
		if err != nil {
			log.Print("scan row failed", err)
			return -1
		}

		break
	}

	return c.count
}
