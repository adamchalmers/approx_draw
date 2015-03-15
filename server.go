package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var urlArg = regexp.MustCompile("url=(.*)")

type imgResponse struct {
	W, H int
	Rgb  rgbArray
}

type rgbArray []uint8

func (u rgbArray) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}
	return []byte(result), nil
}

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

	data := decode(imgfile)
	jsonStr, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Return this data.
	w.Write(jsonStr)
}

func decode(imgfile *os.File) *imgResponse {
	// Decode the image
	img, _, err := image.Decode(imgfile)
	if err != nil {
		log.Fatal("Couldn't read from image.")
	}

	// Extract color data from the decoded image
	bounds := img.Bounds()
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y
	rgb := make(rgbArray, 0)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rgb = append(rgb, uint8(r))
			rgb = append(rgb, uint8(g))
			rgb = append(rgb, uint8(b))
		}
	}

	data := imgResponse{w, h, rgb}
	return &data
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
