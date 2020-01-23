package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/itslab-kyushu/sss/sss"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
)

type Config struct {
	Stores []string `stores`
	Threshold int `threshold`
}

var config Config
var stores []string

func handler(w http.ResponseWriter, r *http.Request) {
	chunksize := 128
	totalShares := 10
    threshold := config.Threshold

	var buffer bytes.Buffer

	r.Write(&buffer)
	shares, _ := sss.Distribute(buffer.Bytes(), chunksize, totalShares, threshold)

	for _, s := range shares {
		data, err := json.Marshal(s)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(string(data))
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
