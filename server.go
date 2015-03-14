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

func (s String) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request) {
	fmt.Fprint(w, "Hello!")
}

func main() {
	port := "localhost:4000"
	fmt.Println("Running on", port)
	err := http.ListenAndServe(port, String("Hello world"))
	if err != nil {
		log.Fatal(err)
	}
}
