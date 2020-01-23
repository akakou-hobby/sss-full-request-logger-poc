package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/itslab-kyushu/sss/sss"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	chunksize := 128
	totalShares := 10
	threshold := 5

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
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
