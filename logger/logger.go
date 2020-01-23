package main

import (
	"fmt"
	"bytes"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer

	r.Write(&buffer)
	fmt.Print(buffer.String())

	fmt.Fprintf(w, "ok")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
