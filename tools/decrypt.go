package main

import (
	"crypto/aes"
	"crypto/cipher"

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
		//"../store/data/5.db",
		//"../store/data/6.db",
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

	ciphertext, _ := sss.Reconstruct(shares)

	fmt.Print(ciphertext)
	fmt.Print("\n\n---\n\n")

	key := []byte("0123456789abcdef")
	block, err := aes.NewCipher(key)
	fmt.Print(err)

	fmt.Print("\n\n---\n\n")

	iv := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	//io.ReadFull(rand.Reader, iv)

	/* dectrypt */
	cipher.NewCFBDecrypter(block, iv)
	cfbdec := cipher.NewCFBDecrypter(block, iv)
	plaintextCopy := make([]byte, len(ciphertext))
	cfbdec.XORKeyStream(plaintextCopy, ciphertext)
	fmt.Print(string(plaintextCopy))

	//fmt.Print(string(secret))
}
