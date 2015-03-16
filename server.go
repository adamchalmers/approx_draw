package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var urlArg = regexp.MustCompile("url=(.*)")

/**************
 * Color code *
 **************/

func (u rgbArray) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}
	return []byte(result), nil
}

func colorDist(c1, c2 color.Color) int {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	sum := math.Abs(float64(r1)-float64(r2)) + math.Abs(float64(g1)-float64(g2)) + math.Abs(float64(b1)-float64(b2))
	return int(sum)
}

/***************
 * Canvas code *
 *************/

// A HTML5 Canvas style representation of pixel data.
type flatCanvas struct {
	W, H int
	Rgb  rgbArray
}
type rgbArray []uint8

// Colors a subrect (x,y,w,h) in the canvas to color (r,g,b).
func mutate(img image.RGBA, x, y, w, h int, rgba color.RGBA) error {

	// Check the mutated region fits inside the canvas.
	if x+w > img.Bounds().Dx() || y+h > img.Bounds().Dy() {
		return fmt.Errorf("Invalid mutation size.")
	}

	// Fill in the coloured region.
	for i := x; i < w+x; i++ {
		for j := y; j < h+y; j++ {
			img.SetRGBA(i, j, rgba)
		}
	}
	return nil
}

// Returns the pixelwise distance between two canvases.
func imgDist(img1, img2 image.RGBA) (int, error) {
	// Check the two canvases are the same size
	if img1.Bounds() != img2.Bounds() {
		return 0, fmt.Errorf("Can't compare different-sized images.")
	}
	sum := 0
	for i := img1.Bounds().Min.X; i < img1.Bounds().Max.X; i++ {
		for j := img1.Bounds().Min.Y; j < img1.Bounds().Max.Y; j++ {
			sum += colorDist(img1.At(i, j), img2.At(i, j))
		}
	}
	return sum, nil
}

// Returns the pixelwise distance between this canvas with a mutation and a second canvas of the same size.
func imgDistMutated(img, other image.RGBA, cachedScore, x, y, w, h int, rgba color.RGBA) (int, error) {

	// Check the mutated region fits inside the canvas.
	if x+w > img.Bounds().Dx() || y+h > img.Bounds().Dy() {
		return 0, fmt.Errorf("Mutation won't fit.")
	}
	// Check the two canvases are the same size
	if img.Bounds() != other.Bounds() {
		return 0, fmt.Errorf("Can't compare different-sized canvases.")
	}
	return 0, nil
}

/**************
 * Server code *
 **************/

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

func testHandler(w http.ResponseWriter, r *http.Request) {
	canvas := image.NewRGBA(image.Rect(0, 0, 30, 30))
	for x := 0; x < 30; x++ {
		for y := 0; y < 30; y++ {
			//canvas.SetRGBA(x, y, color.RGBA{0, 0, 255, 255})
		}
	}
	err := png.Encode(w, canvas)
	if err != nil {
		w.Write([]byte("Encoding error."))
		fmt.Println(err)
	}
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

func decode(imgfile io.ReadCloser) *flatCanvas {
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

	data := flatCanvas{w, h, rgb}
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
	http.HandleFunc("/test/", testHandler)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
