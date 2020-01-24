package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"

	"fmt"
	"github.com/itslab-kyushu/sss/sss"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

type Config struct {
	Stores    []string `stores`
	Threshold int      `threshold`
	Password  string   `password`
}

var config Config
var stores []string

func handler(w http.ResponseWriter, r *http.Request) {
	/* set variables */
	stores = config.Stores

	chunksize := 128
	totalShares := len(stores)
	threshold := config.Threshold

	key := []byte(config.Password)

	/* read requests */
	var buffer bytes.Buffer
	r.Write(&buffer)
	plaintext := buffer.Bytes()

	/* aes encrypt */
	block, _ := aes.NewCipher(key)
	//fmt.Print(err)

	ciphertext := make([]byte, len(plaintext))

	iv := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	// TODO : sharing iv
	// io.ReadFull(rand.Reader, iv)

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext, plaintext)

	/* dectrypt */
	/*
		cfbdec := cipher.NewCFBDecrypter(block, iv)
		plaintextCopy := make([]byte, len(ciphertext))
		cfbdec.XORKeyStream(plaintextCopy, ciphertext)
		fmt.Print(string(plaintextCopy))
	*/

	/* secret sharing */
	shares, _ := sss.Distribute(ciphertext, chunksize, totalShares, threshold)

	for i, s := range shares {
		url := stores[i]
		buf, _ := json.Marshal(s)
		reader := bytes.NewReader(buf)
		http.Post(url, "text/json", reader)

		//fmt.Print(string(buf))
	}
	fmt.Fprintf(w, "ok")
}

func main() {
	buf, _ := ioutil.ReadFile("config.yaml")
	yaml.Unmarshal(buf, &config)

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
