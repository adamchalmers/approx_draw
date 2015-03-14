package main

import (
	"fmt"
	"log"
	"net/http"
)

type Color [3]int
type Rect struct {
	w int32
	h int32
}

type String string

func (s String) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request) {
	fmt.Fprint(w, s)
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path, r.URL.Path[1:])
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	port := "localhost:4000"
	fmt.Println("Running on", port)

	http.HandleFunc("/", fileHandler)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
