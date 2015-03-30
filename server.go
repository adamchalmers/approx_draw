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
	"runtime"
	"strings"
)

const (
	TRIES         = 30
	MUTATIONS     = 10000
	PIXELSAMPLING = 8
	MAXSIZE       = 300
)

var urlArg = regexp.MustCompile("url=(.*)")

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

func myRGBAAt(p *image.RGBA, x, y int) color.RGBA {
	i := p.PixOffset(x, y)
	return color.RGBA{p.Pix[i+0], p.Pix[i+1], p.Pix[i+2], p.Pix[i+3]}
}

// Returns an image which approximately recreates the input image.
func approximate(target *image.RGBA, TRIES, MUTATIONS, PIXELSAMPLING int) (*image.RGBA, int) {

	NCPU := runtime.NumCPU()
	fmt.Printf("%v\n", NCPU)

	// Start with a white background.
	approx := image.NewRGBA(target.Bounds())
	imgW := approx.Bounds().Dx()
	imgH := approx.Bounds().Dy()
	start := mutation{0, 0, imgW, imgH, color.RGBA{255, 255, 255, 255}}
	colors := colorsIn(target)
	mutate(approx, start)

	score, err := imgDist(target, approx)
	if err != nil {
		log.Fatal(err)
	}

	cm := make(chan mutation, NCPU)
	cs := make(chan int, NCPU)

	for i := 0; i < TRIES; i++ {

		// Spawn NCPU goroutines, each of which does MUTATIONS/NCPU mutations.
		for ch := 0; ch < NCPU; ch++ {

			// Calculate the best mutation on this goroutine.
			go findMutation(score, NCPU, MUTATIONS, approx, target, colors, cm, cs)
		}

		// Find the best mutation amongst all the goroutines.
		bestMutation := <-cm
		bestScore := <-cs
		for ch := 1; ch < NCPU; ch++ {
			m := <-cm
			score = <-cs
			if score < bestScore {
				bestMutation = m
				bestScore = score
			}
		}
		// Apply the best mutation,
		// then restart the loop to place a new rectangle in the image.
		mutate(approx, bestMutation)
	}
	return approx, score
}

func findMutation(score, NCPU, MUTATIONS int, approx, target *image.RGBA, colors []color.RGBA, cm chan mutation, cs chan int) {
	imgW := approx.Bounds().Dx()
	imgH := approx.Bounds().Dy()
	cachedScore := score
	bestScore := cachedScore
	var bestMutation mutation

	// Try MUTATIONS different mutations and keep the best one.
	for try := 0; try < MUTATIONS/NCPU; try++ {

		// Generate a mutation
		w := rand.Intn(imgW)
		h := rand.Intn(imgH)
		x := rand.Intn(imgW - w)
		y := rand.Intn(imgH - h)
		rgb := colors[rand.Intn(len(colors))]
		m := mutation{x, y, w, h, rgb}

		// Save this mutation if it's the best.
		tryScore := imgDistMutated(approx, target, cachedScore, m, PIXELSAMPLING)
		if tryScore < bestScore {
			bestScore = tryScore
			bestMutation = m
		}

	}
	cm <- bestMutation
	cs <- bestScore
}

// Returns a slice containing all colors used in the image, and
// a map from each point/pixel in the image to its color.
// looking up colors in this map is 100x faster than using img.at again.
func colorsIn(img *image.RGBA) []color.RGBA {
	colsList := make([]color.RGBA, 1000)
	cols := make(map[color.RGBA]bool)
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			color := myRGBAAt(img, x, y)
			if _, prs := cols[color]; !prs {
				cols[color] = true
				colsList = append(colsList, color)
			}
		}
	}
	return colsList
}

// RGB distance between two colors.
func colorDist(c1, c2 color.RGBA) int {
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
			sum += colorDist(myRGBAAt(img1, i, j), myRGBAAt(img2, i, j))
		}
	}
	return sum, nil
}

// Returns the pixelwise distance between this canvas with a mutation and a second canvas of the same size.
func imgDistMutated(img, target *image.RGBA, cachedScore int, m mutation, PIXELSAMPLING int) int {
	score := cachedScore
	for i := m.x; i < m.x+m.w; i += PIXELSAMPLING {
		for j := m.y; j < m.y+m.h; j += PIXELSAMPLING {
			// Subtract the original color's score, add the mutated color's score.
			col := myRGBAAt(target, i, j)
			score -= colorDist(col, myRGBAAt(img, i, j))
			score += colorDist(col, m.rgb)
		}
	}
	return score
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
func urlParam(r *http.Request) (string, error) {
	m := urlArg.FindStringSubmatch(r.URL.String())
	if m == nil || len(m) < 1 {
		return "", fmt.Errorf("Invalid regex.", r.URL.String())
	}
	if _, err := url.Parse(m[1]); err != nil {
		return "", fmt.Errorf("Invalid url", m[1])
	}
	return m[1], nil
}

// Serves the image from the URL in the request.
// This allows us to get around CORS issues.
func remoteHandler(w http.ResponseWriter, r *http.Request) {
	img, err := getImg(r)
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}
	defer img.Close()
	io.Copy(w, img)
}

// Serves the image generated by approximate() on the URL in the request.
func approxHandler(w http.ResponseWriter, r *http.Request) {
	img, err := getImg(r)
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}
	defer img.Close()
	// read the image into target (type image.Image)
	_target, _, err := image.Decode(img)
	if err != nil {
		w.Write([]byte("err"))
		fmt.Println(err)
		return
	}
	if dx, dy := _target.Bounds().Dx(), _target.Bounds().Dy(); dx > MAXSIZE || dy > MAXSIZE {
		msg := fmt.Sprintf("Image of size %v.%v is too large (max size is %v)", dx, dy, MAXSIZE)
		io.Copy(w, strings.NewReader(msg))
	}

	target := toRGBA(_target)

	approximation, score := approximate(target, TRIES, MUTATIONS, PIXELSAMPLING)
	fmt.Println(float64(score) / 1000000)
	png.Encode(w, approximation)
}

// Serve a static file.
func fileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("%v iterations of %v mutations each, sampling 1/%v pixels.", TRIES, MUTATIONS, PIXELSAMPLING)
	io.Copy(w, strings.NewReader(msg))
}

// Returns the image from a URL.
func getImg(r *http.Request) (io.ReadCloser, error) {
	// Get the URL of the target image
	url, err := urlParam(r)
	if err != nil {
		return nil, err
	}
	// Fetch the image
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func main() {
	runtime.GOMAXPROCS(4)
	port := "localhost:4000"
	fmt.Println("Running on", port)

	http.HandleFunc("/", fileHandler)
	http.HandleFunc("/remote/", remoteHandler)
	http.HandleFunc("/approx/", approxHandler)
	http.HandleFunc("/stats/", statsHandler)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
