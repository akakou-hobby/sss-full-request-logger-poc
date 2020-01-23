package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/itslab-kyushu/sss/sss"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
)

func main() {
	dbfiles := []string{
		"../store/data/1.db",
		"../store/data/2.db",
		"../store/data/3.db",
		"../store/data/4.db",
		"../store/data/5.db",
		"../store/data/6.db",
		"../store/data/7.db"}

	arg := os.Args[1]
	recordId, _ := strconv.Atoi(arg)

	shares := make([]sss.Share, len(dbfiles))

	for i, dbfile := range dbfiles {
		db, _ := sql.Open("sqlite3", dbfile)
		var data []byte

		row := db.QueryRow(`SELECT data FROM payload WHERE id=?`, recordId)
		row.Scan(&data)

		json.Unmarshal(data, &shares[i])

		db.Close()

	}

	secret, _ := sss.Reconstruct(shares)

	fmt.Print(string(secret))
}
