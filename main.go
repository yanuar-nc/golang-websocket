// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	env "github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 10 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	filePeriod = 10 * time.Second
)

// openConnection to get connection to database
func openConnection() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", os.Getenv("DB_NAME"))
	if err != nil {
		return nil, err
	}

	return db, err

}

// readUserCount count total user in database realtime
func readUserCount(lastMod time.Time) ([]byte, time.Time, error) {

	var (
		countInt int
	)
	db, err := openConnection()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT count(*) FROM users")
	if err != nil {
		return nil, lastMod, err
	}

	for rows.Next() {
		if err = rows.Scan(&countInt); err != nil {
			return nil, lastMod, err
		}
	}

	rows.Close()
	db.Close()

	strCount := strconv.Itoa(countInt)
	return []byte(strCount), lastMod, nil
}

func main() {
	env.Load()

	var port = os.Getenv("PORT")
	// addr port address of application
	var addr = ":" + os.Getenv("PORT")

	fmt.Println("Listen: " + port)

	// router
	http.HandleFunc("/", ServeHome)
	http.HandleFunc("/ws", ServeWs)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
