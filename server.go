package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"regexp"
)

var urlArg = regexp.MustCompile("url=(.*)")

func imageHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the image URL from the request
	m := urlArg.FindStringSubmatch(r.URL.String())
	if m == nil {
		fmt.Println("Invalid regex.")
		return
	}

	// Open the image
	imgfile, err := os.Open("./" + m[1])
	if err != nil {
		fmt.Println("Can't find", m[1])
		return
	}
	defer imgfile.Close()

	// Decode the image
	img, _, err := image.Decode(imgfile)
	if err != nil {
		log.Fatal("Couldn't read from image.")
	}

	// Extract color data from the decoded image
	bounds := img.Bounds()
	rgba := ""
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			rgba += fmt.Sprintf("[%v, %v, %v, %v],", r, g, b, a)
		}
	}

	// Return this data.
	b := make([]byte, len(rgba))
	copy(b, rgba)
	w.Write(b)
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	port := "localhost:4000"
	fmt.Println("Running on", port)

	http.HandleFunc("/", fileHandler)
	http.HandleFunc("/image/", imageHandler)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
