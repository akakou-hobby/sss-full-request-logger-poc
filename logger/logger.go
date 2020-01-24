package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/itslab-kyushu/sss/sss"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

type Config struct {
	Stores    []string `stores`
	Threshold int      `threshold`
}

var config Config
var stores []string

var prikeyFile = "keys/prikey.pem"
var pubkeyFile = "keys/pubkey.pem"

var prikey *rsa.PrivateKey
var pubkey *rsa.PublicKey

func handler(w http.ResponseWriter, r *http.Request) {
	stores = config.Stores

	chunksize := 128
	totalShares := len(stores)
	threshold := config.Threshold

	var buffer bytes.Buffer

	r.Write(&buffer)
	shares, _ := sss.Distribute(buffer.Bytes(), chunksize, totalShares, threshold)

	for i, s := range shares {
		url := stores[i]
		buf, _ := json.Marshal(s)
		reader := bytes.NewReader(buf)
		http.Post(url, "text/json", reader)

		fmt.Print(string(buf))
	}

	fmt.Fprintf(w, "ok")
}

func main() {
	buf, _ := ioutil.ReadFile("config.yaml")
	yaml.Unmarshal(buf, &config)

	// fmt.Printf("--- config:\n%v\n\n", config.Stores)

	bytes, _ := ioutil.ReadFile(pubkeyFile)
	block, _ := pem.Decode(bytes)
	_pubkey, _ := x509.ParsePKIXPublicKey(block.Bytes)
	pubkey = _pubkey.(*rsa.PublicKey)

	label := []byte("label")
	message := []byte("message")

	ciphertext, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubkey, message, label)

	{
		bytes, _ := ioutil.ReadFile(prikeyFile)
		block, _ := pem.Decode(bytes)
		prikey, _ = x509.ParsePKCS1PrivateKey(block.Bytes)
		plaintext, _ := rsa.DecryptOAEP(sha256.New(), rand.Reader, prikey, ciphertext, label)

		fmt.Print(string(plaintext))
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
