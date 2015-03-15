package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"net/url"
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

// Returns the url query parameter
// e.g. in /remote/img?url=wwww.google.com, returns www.google.com
func urlParam(r *http.Request) string {
	m := urlArg.FindStringSubmatch(r.URL.String())
	if m == nil || len(m) < 1 {
		log.Fatal("Invalid regex.", r.URL.String())
	}
	if _, err := url.Parse(m[1]); err != nil {
		log.Println("Invalid url", m[1])
		return ""
	}
	return m[1]
}

func remoteHandler(w http.ResponseWriter, r *http.Request) {
	url := urlParam(r)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		w.Write([]byte("err"))
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the image URL from the request
	url := urlParam(r)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// Open the image
	data := decode(resp.Body)
	jsonStr, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Return this data.
	w.Write(jsonStr)
}

func decode(imgfile io.ReadCloser) *imgResponse {
	// Decode the image
	img, _, err := image.Decode(imgfile)
	if err != nil {
		log.Fatal("Couldn't read from image.")
	}

	// Extract color data from the decoded image
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
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
	http.HandleFunc("/remote/", remoteHandler)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
