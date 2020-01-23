package main

import (
	"bytes"
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
}

var config Config
var stores []string

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

	fmt.Printf("--- config:\n%v\n\n", config.Stores)

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
