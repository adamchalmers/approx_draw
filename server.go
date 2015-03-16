package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
)

var urlArg = regexp.MustCompile("url=(.*)")

const (
	ITERATIONS = 10
	TRIES      = 1000
)

func abs(x, y uint8) int {
	if x > y {
		return int(x - y)
	} else {
		return int(y - x)
	}
}

/**************
 * Image code *
 **************/

type mutation struct {
	x, y, w, h int
	rgb        color.RGBA
}

// Returns an image which approximately recreates the input image.
func approximate(target *image.RGBA) (*image.RGBA, int) {

	// Start with a white background.
	approx := image.NewRGBA(target.Bounds())
	imgW := approx.Bounds().Dx()
	imgH := approx.Bounds().Dy()
	start := mutation{0, 0, imgW, imgH, color.RGBA{255, 255, 255, 255}}
	colors, targetCache := colorsIn(target)
	mutate(approx, start)

	// Loop
	score, err := imgDist(target, approx)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < ITERATIONS; i++ {
		cachedScore := score
		var bestMutation mutation
		for try := 0; try < TRIES; try++ {

			// Generate a mutation
			w := rand.Intn(imgW)
			h := rand.Intn(imgH)
			x := rand.Intn(imgW - w)
			y := rand.Intn(imgH - h)
			rgb := colors[rand.Intn(len(colors))]

			// Save this mutation if it's the best.
			tryScore, err := imgDistMutated(approx, target, targetCache, cachedScore, x, y, w, h, rgb)
			if err != nil {
				log.Fatal(err)
			}
			if tryScore < score {
				score = tryScore
				bestMutation = mutation{x, y, w, h, rgb}
			}

		} // end tries
		mutate(approx, bestMutation)
	} // end iterations
	return approx, score
}

// Returns a slice containing all colors used in the image, and
// a map from each point/pixel in the image to its color.
// looking up colors in this map is 100x faster than using img.at again.
func colorsIn(img *image.RGBA) ([]color.RGBA, map[image.Point]color.RGBA) {
	colsList := make([]color.RGBA, 1000)
	cols := make(map[color.RGBA]bool)
	cache := make(map[image.Point]color.RGBA)
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			color := img.At(x, y).(color.RGBA)
			if _, prs := cols[color]; !prs {
				cols[color] = true
				colsList = append(colsList, color)
				cache[image.Point{x, y}] = color
			}
		}
	}
	return colsList, cache
}

// RGB distance between two colors.
func colorDist(_c1, _c2 color.Color) int {
	c1, c2 := _c1.(color.RGBA), _c2.(color.RGBA)
	sum := abs(c1.R, c2.R)
	sum += abs(c1.G, c2.G)
	sum += abs(c1.B, c2.B)
	return int(sum)
}

// Colors a subrect (x,y,w,h) in the canvas to color (r,g,b).
func mutate(img *image.RGBA, m mutation) error {

	// Check the mutated region fits inside the canvas.
	if m.x+m.w > img.Bounds().Dx() || m.y+m.h > img.Bounds().Dy() {
		return fmt.Errorf("Invalid mutation size.")
	}

	// Fill in the coloured region.
	for i := m.x; i < m.w+m.x; i++ {
		for j := m.y; j < m.h+m.y; j++ {
			img.SetRGBA(i, j, m.rgb)
		}
	}
	return nil
}

// Returns the pixelwise distance between two canvases.
func imgDist(img1, img2 *image.RGBA) (int, error) {
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
func imgDistMutated(img, target *image.RGBA, targetCache map[image.Point]color.RGBA, cachedScore, x, y, w, h int, rgba color.RGBA) (int, error) {
	// Check the mutated region fits inside the canvas.
	if x+w > img.Bounds().Dx() || y+h > img.Bounds().Dy() {
		return 0, fmt.Errorf("Mutation won't fit.")
	}
	// Check the two canvases are the same size
	if img.Bounds() != target.Bounds() {
		return 0, fmt.Errorf("Can't compare different-sized canvases.")
	}
	score := cachedScore
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			// Subtract the original color's score, add the mutated color's score.
			col, present := targetCache[image.Point{i, j}]
			if !present {
				col = target.At(i, j).(color.RGBA)
			}
			score -= colorDist(col, img.At(i, j))
			score += colorDist(col, rgba)
		}
	}
	return score, nil
}

func toRGBA(_target image.Image) *image.RGBA {
	target := image.NewRGBA(_target.Bounds())
	for x := target.Bounds().Min.X; x < target.Bounds().Max.X; x++ {
		for y := target.Bounds().Min.Y; y < target.Bounds().Max.Y; y++ {
			target.Set(x, y, _target.At(x, y))
		}
	}
	return target
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

func approxHandler(w http.ResponseWriter, r *http.Request) {

	// read the image data
	url := urlParam(r)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		w.Write([]byte("err"))
	}
	defer resp.Body.Close()

	// read the image into target (type image.Image)
	_target, _, err := image.Decode(resp.Body)
	if err != nil {
		w.Write([]byte("err"))
		fmt.Println(err)
		return
	}

	target := toRGBA(_target)

	approximation, score := approximate(target)
	fmt.Println(float64(score) / 1000000)
	png.Encode(w, approximation)
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	port := "localhost:4000"
	fmt.Println("Running on", port)

	http.HandleFunc("/", fileHandler)
	http.HandleFunc("/remote/", remoteHandler)
	http.HandleFunc("/approx/", approxHandler)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
